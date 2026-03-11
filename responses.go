package main

type powerData struct {
	Power             int     `json:"power"`
	PowerAvg          int     `json:"powerAvg"`
	AgePower          uint64  `json:"agePower"`
	PowerPhase1       int64   `json:"powerPhase1"`
	PowerPhase2       int64   `json:"powerPhase2"`
	PowerPhase3       int64   `json:"powerPhase3"`
	EnergyCounterOut  float64 `json:"energyCounterOut"`
	EnergyCounterIn   float64 `json:"energyCounterIn"`
	EnergyCounterInT1 float64 `json:"energyCounterInT1"`
	EnergyCounterInT2 float64 `json:"energyCounterInT2"`
}
