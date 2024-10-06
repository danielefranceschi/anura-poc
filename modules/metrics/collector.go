// Copyright 2018 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package metrics

import (
	activities_model "code.gitea.io/gitea/models/activities"
	"code.gitea.io/gitea/models/db"

	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "gitea_"

// Collector implements the prometheus.Collector interface and
// exposes gitea metrics for prometheus
type Collector struct {
	Accesses     *prometheus.Desc
	Attachments  *prometheus.Desc
	LoginSources *prometheus.Desc
	Oauths       *prometheus.Desc
	Users        *prometheus.Desc
}

// NewCollector returns a new Collector with all prometheus.Desc initialized
func NewCollector() Collector {
	return Collector{
		Accesses: prometheus.NewDesc(
			namespace+"accesses",
			"Number of Accesses",
			nil, nil,
		),
		Attachments: prometheus.NewDesc(
			namespace+"attachments",
			"Number of Attachments",
			nil, nil,
		),
		LoginSources: prometheus.NewDesc(
			namespace+"loginsources",
			"Number of LoginSources",
			nil, nil,
		),
		Oauths: prometheus.NewDesc(
			namespace+"oauths",
			"Number of Oauths",
			nil, nil,
		),
		Users: prometheus.NewDesc(
			namespace+"users",
			"Number of Users",
			nil, nil,
		),
	}
}

// Describe returns all possible prometheus.Desc
func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Accesses
	ch <- c.Attachments
	ch <- c.LoginSources
	ch <- c.Oauths
	ch <- c.Users
}

// Collect returns the metrics with values
func (c Collector) Collect(ch chan<- prometheus.Metric) {
	stats := activities_model.GetStatistic(db.DefaultContext)

	ch <- prometheus.MustNewConstMetric(
		c.LoginSources,
		prometheus.GaugeValue,
		float64(stats.Counter.AuthSource),
	)
	ch <- prometheus.MustNewConstMetric(
		c.Oauths,
		prometheus.GaugeValue,
		float64(stats.Counter.Oauth),
	)
	ch <- prometheus.MustNewConstMetric(
		c.Users,
		prometheus.GaugeValue,
		float64(stats.Counter.User),
	)
}
