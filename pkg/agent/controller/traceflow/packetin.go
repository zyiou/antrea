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
		_, err = c.ParsePacketIn(pktIn)
		return false, err
	}, stopCh)
}

func (c *Controller) ParsePacketIn(pktIn *ofctrl.PacketIn) (*v1.Traceflow, error) {
	matchers := pktIn.GetMatches()
	var match *ofctrl.MatchField

	// get cross node tag
	if match = getMatchField(matchers, 9); match == nil {
		return nil, errors.New("traceflow cross node tag not found")
	}
	rngTag := openflow13.NewNXRange(28, 31)
	tag, err := getInfoInReg(matchers, match, rngTag)
	if err != nil {
		return nil, err
	}

	// get CRD
	tf, err := c.GetTraceflowCRD(uint8(tag))
	if err != nil {
		return nil, err
	}

	var obs []v1.Observation
	ifSender := ifSender(uint8(tag), c.senders)
	tableID := pktIn.TableId
	if ifSender {
		obs = tf.NodeSender
	} else {
		obs = tf.NodeReceiver
	}

	// get drop table
	if match = getMatchField(matchers, 3); match != nil {
		rngDrop := openflow13.NewNXRange(0, 7)
		_, err := getInfoInReg(matchers, match, rngDrop)
		if err != nil {
			return nil, err
		}
		ob := new(v1.Observation)
		if ifSender {
			ob.ResourceType = v1.DROPPED
		} else {
			ob.ResourceType = v1.RECEIVED
		}
		ob.ComponentType = v1.DFW
		obs = append(obs, *ob)
		tf.Phase = v1.SUCCESS
	}

	// ingress Rule, DFW
	if match = getMatchField(matchers, 4); match != nil {
		ingress, err := getInfoInReg(matchers, match, nil)
		if err != nil {
			return nil, err
		}
		ob := new(v1.Observation)
		ob.RuleID = fmt.Sprint(ingress)
		ob.ComponentType = v1.DFW
		obs = append(obs, *ob)
	}

	// egress Rule, DFW
	if match = getMatchField(matchers, 5); match != nil {
		egress, err := getInfoInReg(matchers, match, nil)
		if err != nil {
			return nil, err
		}
		ob := new(v1.Observation)
		ob.RuleID = fmt.Sprint(egress)
		ob.ComponentType = v1.DFW
		obs = append(obs, *ob)
	}

	// get output table
	if tableID == 110 {
		rngOut := openflow13.NewNXRange(8, 15)
		_, err := getInfoInReg(matchers, match, rngOut)
		if err != nil {
			return nil, err
		}
		ob := new(v1.Observation)
		ob.ComponentName = tableIDToComponentName[tableID]
		ob.ResourceType = v1.FORWARDED
		ob.ComponentType = v1.FORWARDING
		obs = append(obs, *ob)
	}
	
	return tf, nil
}

func ifSender(tag uint8, senders map[uint8]*v1.Traceflow) bool {
	if _, ok := senders[tag]; ok {
		delete(senders, tag)
		return true
	}
	return false
}

// GetTraceflowCRD returns traceflow CRD with cross node tag.
func (c *Controller) GetTraceflowCRD(tag uint8) (*v1.Traceflow, error) {
	if tf, ok := c.running[tag]; ok {
		return tf, nil
	}
	return nil, errors.New("traceflow with the cross tag node doesn't exist")
}
