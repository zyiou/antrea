// Copyright 2020 Antrea Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package traceflow

import (
	"errors"
	"fmt"
	"time"

	"github.com/contiv/libOpenflow/openflow13"
	"github.com/contiv/ofnet/ofctrl"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/vmware-tanzu/antrea/pkg/apis/traceflow/v1"
)

func (c *Controller) ReceivePacketIn(stopCh <-chan struct{}) {
	ch := make(chan *ofctrl.PacketIn)
	c.ofClient.SubscribePacketIn(1, ch)

	wait.PollUntil(time.Second, func() (done bool, err error) {
		pktIn := <-ch
		_, err = parsePacketIn(pktIn, c)
		return false, err
	}, stopCh)
}

func parsePacketIn(pktIn *ofctrl.PacketIn, q Querier) (*v1.Traceflow, error) {
	matchers := pktIn.GetMatches()
	var match *ofctrl.MatchField

	// get cross node tag
	if match = getMatchField(matchers, 0); match == nil {
		return nil, errors.New("traceflow cross node tag not found")
	}
	rngTag := openflow13.NewNXRange(28, 31)
	tag, err := getInfoInReg(matchers, match, rngTag)
	if err != nil {
		return nil, err
	}

	// get CRD
	tf, err := q.GetTraceflowCRD(uint8(tag))
	if err != nil {
		return nil, err
	}

	var ob *v1.Observation
	ifSender := ifSender(tag)
	tableID := pktIn.TableId
	if ifSender {
		ob = &tf.Status.NodeSender[tableIDToObIdx[tableID]]
	} else {
		ob = &tf.Status.NodeReceiver[tableIDToObIdx[tableID]]
	}
	ob.ComponentName = tableIDToComponentName[tableID]

	// get drop table
	if match = getMatchField(matchers, 3); match != nil {
		rngDrop := openflow13.NewNXRange(0, 7)
		_, err := getInfoInReg(matchers, match, rngDrop)
		if err != nil {
			return nil, err
		}
		if ifSender {
			ob.ResourceType = v1.DROPPED
		} else {
			ob.ResourceType = v1.RECEIVED
		}
		ob.ComponentType = v1.DFW
		tf.Status.Phase = v1.SUCCESS
	}

	// get output table
	if match = getMatchField(matchers, 3); match != nil {
		rngOut := openflow13.NewNXRange(8, 15)
		_, err := getInfoInReg(matchers, match, rngOut)
		if err != nil {
			return nil, err
		}
		ob.ResourceType = v1.FORWARDED
		ob.ComponentType = v1.FORWARDING
	}

	// ingress Rule, DFW
	if match = getMatchField(matchers, 4); match != nil {
		ingress, err := getInfoInReg(matchers, match, nil)
		if err != nil {
			return nil, err
		}
		ob.RuleID = fmt.Sprint(ingress)
		ob.ComponentType = v1.DFW
	}

	// egress Rule, DFW
	if match = getMatchField(matchers, 5); match != nil {
		egress, err := getInfoInReg(matchers, match, nil)
		if err != nil {
			return nil, err
		}
		ob.RuleID = fmt.Sprint(egress)
		ob.ComponentType = v1.DFW
	}

	return tf, nil
}

func ifSender(tag uint32) bool {
	if _, ok := currentTags[tag]; ok {
		delete(currentTags, tag)
		return true
	}
	return false
}
