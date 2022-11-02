package model

type Measure struct {
	Domain        string
	Key           string
	Id            int
	MeasurePoints []MeasurePoint // value type is convenient to reflect
}

type MeasurePoint struct {
	PointName     string  `json:"pointName"`
	Date          string  `json:"date"`
	Time          string  `json:"time"`
	DataType      int     `json:"dataType"`
	DX            float64 `json:"dx"`
	DY            float64 `json:"dy"`
	DZ            float64 `json:"dz"`
	D             float64 `json:"d"`
	X             float64 `json:"x"`
	Y             float64 `json:"y"`
	Z             float64 `json:"z"`
	DistributionX int     `json:"distributionX"`
	DistributionY int     `json:"distributionY"`
	DistributionZ int     `json:"distributionZ"`
	DistributionD int     `json:"distributionD"`
	OffsetX       float64 `json:"offsetX"`
	OffsetY       float64 `json:"offsetY"`
	OffsetZ       float64 `json:"offsetZ"`
	ThermalComX   float64 `json:"thermalComX,omitempty"`
	ThermalComY   float64 `json:"thermalComY,omitempty"`
	ThermalComZ   float64 `json:"thermalComZ,omitempty"`
	MeasureMode   int     `json:"measureMode,omitempty"`
}
