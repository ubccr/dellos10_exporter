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

type powerSupply struct {
	ID     int `json:"psu-id"`
	Status string
}

type systemInfo struct {
	NodeID        int            `json:"node-id"`
	PowerSupplies []*powerSupply `json:"power-supply"`
}

type SystemProber struct {
	Info *systemInfo `json:"dell-equipment:node"`

	powerGuage *prometheus.GaugeVec
}

func (p *SystemProber) GetCmdEndpoint() string {
	return "/data/dell-equipment:system/node"
}

func (p *SystemProber) Register(registry *prometheus.Registry) {
	p.powerGuage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "dellos10_power_supply_state",
		Help: "Contains Power Supply state",
	}, []string{"powerSupply", "state"})

	registry.MustRegister(p.powerGuage)
}

func (p *SystemProber) Handler(logger log.Logger) {
	for _, ps := range p.Info.PowerSupplies {
		if ps.Status == "up" {
			p.powerGuage.WithLabelValues(fmt.Sprintf("%d", ps.ID), ps.Status).Set(1)
		} else {
			p.powerGuage.WithLabelValues(fmt.Sprintf("%d", ps.ID), ps.Status).Set(0)
		}
	}
}
