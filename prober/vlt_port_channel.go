// Copyright 2024 Andrew E. Bruno
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prober

import (
	"fmt"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type peer struct {
	UnitID          float64 `json:"unit-id"`
	LagID           string  `json:"lag-id"`
	LagStatus       string  `json:"lag-status"`
	ConfiguredPorts float64 `json:"configured-ports"`
	ActivePorts     float64 `json:"active-ports"`
}

type vltPortChannel struct {
	ID    int
	Peers []*peer
}

type VltPortChannelProber struct {
	VltPortChannels []*vltPortChannel `json:"dell-vlt:vlt-port-info"`

	vltPortChannelGauge *prometheus.GaugeVec
}

func (p *VltPortChannelProber) GetCmdEndpoint() string {
	return "/data/dell-node-management:topology-oper-data=1/dell-vlt:vlt-domain/vlt-port-info"
}

func (p *VltPortChannelProber) Register(registry *prometheus.Registry) {
	p.vltPortChannelGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "dellos10_vlt_portchannel_ports",
		Help: "Contains VLT port channel ports by state",
	}, []string{"interface", "state"})

	registry.MustRegister(p.vltPortChannelGauge)
}

func (p *VltPortChannelProber) Handler(logger log.Logger) {
	for _, pc := range p.VltPortChannels {
		active := float64(0)
		configured := float64(0)
		id := fmt.Sprintf("%d", pc.ID)
		for _, peer := range pc.Peers {
			active += peer.ActivePorts
			configured += peer.ConfiguredPorts
			// XXX assume same port-channel name across both peers
			id = peer.LagID
		}

		p.vltPortChannelGauge.WithLabelValues(id, "active").Set(active)
		p.vltPortChannelGauge.WithLabelValues(id, "inactive").Set(configured - active)
	}
}
