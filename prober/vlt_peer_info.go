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

type vltPeerInfo struct {
	UnitID int `json:"unit-id"`
	Role   string
	Status string `json:"peer-status"`
}

type VltPeerInfoProber struct {
	Peers []*vltPeerInfo `json:"dell-vlt:peer-info"`

	infoGauge *prometheus.GaugeVec
}

func (p *VltPeerInfoProber) GetCmdEndpoint() string {
	return "/data/dell-node-management:topology-oper-data=1/dell-vlt:vlt-domain/dell-vlt:peer-info"
}

func (p *VltPeerInfoProber) Register(registry *prometheus.Registry) {
	p.infoGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "dellos10_vlt_peer_info",
		Help: "Contains VLT peer info",
	}, []string{"role", "status", "unitId"})

	registry.MustRegister(p.infoGauge)
}

func (p *VltPeerInfoProber) Handler(logger log.Logger) {
	for _, peer := range p.Peers {
		if peer.Status == "up" {
			p.infoGauge.WithLabelValues(peer.Role, peer.Status, fmt.Sprintf("%d", peer.UnitID)).Set(1)
		} else {
			p.infoGauge.WithLabelValues(peer.Role, peer.Status, fmt.Sprintf("%d", peer.UnitID)).Set(0)
		}
	}
}
