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
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
)

type vltLocalInfo struct {
	UnitID float64 `json:"unit-id"`
	Role   string
}

type VltLocalInfoProber struct {
	Info *vltLocalInfo `json:"dell-vlt:local-info"`

	infoGauge *prometheus.GaugeVec
}

func (p *VltLocalInfoProber) GetCmdEndpoint() string {
	return "/data/dell-node-management:topology-oper-data=1/dell-vlt:vlt-domain/dell-vlt:local-info"
}

func (p *VltLocalInfoProber) Register(registry *prometheus.Registry) {
	p.infoGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "dellos10_vlt_unit_id",
		Help: "Contains VLT unit id",
	}, []string{"role"})

	registry.MustRegister(p.infoGauge)
}

func (p *VltLocalInfoProber) Handler(logger log.Logger) {
	p.infoGauge.WithLabelValues(p.Info.Role).Set(p.Info.UnitID)
}
