package main

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type metrics struct {
	power             *prometheus.Desc
	powerAvg          *prometheus.Desc
	powerLastUpdate   *prometheus.Desc
	powerPhase1       *prometheus.Desc
	powerPhase2       *prometheus.Desc
	powerPhase3       *prometheus.Desc
	energyCounterOut  *prometheus.Desc
	energyCounterIn   *prometheus.Desc
	energyCounterInT1 *prometheus.Desc
	energyCounterInT2 *prometheus.Desc
	up                *prometheus.Desc
}

var collectorLog = log.Sugar().Named("collector")

func newMetrics() metrics {
	return metrics{
		power: prometheus.NewDesc(
			"ecotracker_power",
			"Current power consumption or feed-in",
			[]string{"ecotracker_host"}, nil),

		powerAvg: prometheus.NewDesc(
			"ecotracker_power_avg",
			"Average power consumption or feed-in in the last minute",
			[]string{"ecotracker_host"}, nil),

		powerLastUpdate: prometheus.NewDesc(
			"ecotracker_power_last_update",
			"Last time the power information was updated at the electricity meter (in seconds)",
			[]string{"ecotracker_host"}, nil),

		powerPhase1: prometheus.NewDesc(
			"ecotracker_power_phase1",
			"Current power on phase 1 of the electricity meter",
			[]string{"ecotracker_host"}, nil),

		powerPhase2: prometheus.NewDesc(
			"ecotracker_power_phase2",
			"Current power on phase 2 of the electricity meter",
			[]string{"ecotracker_host"}, nil),

		powerPhase3: prometheus.NewDesc(
			"ecotracker_power_phase3",
			"Current power on phase 3 of the electricity meter",
			[]string{"ecotracker_host"}, nil),

		energyCounterOut: prometheus.NewDesc(
			"ecotracker_energy_counter_out",
			"Total energy counted as feed-in by the electricity meter",
			[]string{"ecotracker_host"}, nil),

		energyCounterIn: prometheus.NewDesc(
			"ecotracker_energy_counter_in",
			"Total energy consumption counted by the electricity meter",
			[]string{"ecotracker_host"}, nil),

		energyCounterInT1: prometheus.NewDesc(
			"ecotracker_energy_counter_in_t1",
			"Total energy consumption counted by the electricity meter for rate 1",
			[]string{"ecotracker_host"}, nil),

		energyCounterInT2: prometheus.NewDesc(
			"ecotracker_energy_counter_in_t2",
			"Total energy consumption counted by the electricity meter for rate 2",
			[]string{"ecotracker_host"}, nil),

		up: prometheus.NewDesc(
			"ecotracker_up",
			"The status of the exporter (0 means there was a problem querying the EcoTracker API, see logs)",
			[]string{"ecotracker_host"}, nil),
	}
}

func (m metrics) Collect(ch chan<- prometheus.Metric) {
	collectorLog.Debug("Collecting metrics …")

	powerData, error := getPowerData(*host, *port)
	if error != nil {
		collectorLog.Errorf("Error fetching power data: %s", error)

		ch <- prometheus.MustNewConstMetric(m.up, prometheus.GaugeValue, 0)

		return
	}

	lastUpdateTime := time.Now().Add(-time.Duration(powerData.AgePower) * time.Millisecond).Local().UnixMilli()

	ch <- prometheus.MustNewConstMetric(m.power, prometheus.GaugeValue, float64(powerData.Power), *host)
	ch <- prometheus.MustNewConstMetric(m.powerAvg, prometheus.GaugeValue, float64(powerData.PowerAvg), *host)
	ch <- prometheus.MustNewConstMetric(m.powerLastUpdate, prometheus.CounterValue, float64(lastUpdateTime), *host)
	ch <- prometheus.MustNewConstMetric(m.powerPhase1, prometheus.GaugeValue, float64(powerData.PowerPhase1), *host)
	ch <- prometheus.MustNewConstMetric(m.powerPhase2, prometheus.GaugeValue, float64(powerData.PowerPhase2), *host)
	ch <- prometheus.MustNewConstMetric(m.powerPhase3, prometheus.GaugeValue, float64(powerData.PowerPhase3), *host)
	ch <- prometheus.MustNewConstMetric(m.energyCounterOut, prometheus.CounterValue, float64(powerData.EnergyCounterOut), *host)
	ch <- prometheus.MustNewConstMetric(m.energyCounterIn, prometheus.CounterValue, float64(powerData.EnergyCounterIn), *host)
	ch <- prometheus.MustNewConstMetric(m.energyCounterInT1, prometheus.CounterValue, float64(powerData.EnergyCounterInT1), *host)
	ch <- prometheus.MustNewConstMetric(m.energyCounterInT2, prometheus.CounterValue, float64(powerData.EnergyCounterInT2), *host)
	ch <- prometheus.MustNewConstMetric(m.up, prometheus.GaugeValue, 1, *host)
}

func (m metrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.up
}
