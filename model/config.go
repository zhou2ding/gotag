package model

type Config struct {
	Domain string      `json:"domain"`
	Key    string      `json:"key"`
	Id     int         `json:"id"`
	Value  ConfigValue `json:"value"`
}

type ConfigValue struct {
	PlanName      string               `json:"planName"`
	MeasurePoints []ConfigMeasurePoint `json:"measurePoints"` // value type is convenient to reflect
}

type ConfigMeasurePoint struct {
	PlanID       int     `json:"planId,omitempty"`
	PointName    string  `json:"pointName"`
	TrajId       int     `json:"trajId,omitempty"`
	PMode        int     `json:"pMode"`
	PointId      int     `json:"pointId"`
	Factory      string  `json:"factory"`
	Station      string  `json:"station"`
	TheoreticalX float64 `json:"theoreticalX"`
	TheoreticalY float64 `json:"theoreticalY"`
	TheoreticalZ float64 `json:"theoreticalZ"`
	I            float64 `json:"i"`
	J            float64 `json:"j"`
	K            float64 `json:"k"`

	MonitorD int `json:"monitorD"`
	MonitorX int `json:"monitorX"`
	MonitorY int `json:"monitorY"`
	MonitorZ int `json:"monitorZ"`

	PointType   int `json:"pointType"`
	RepeatPoint int `json:"repeatPoint"`

	Tol1minusD float64 `json:"tol1minusD"`
	Tol1minusX float64 `json:"tol1minusX"`
	Tol1minusY float64 `json:"tol1minusY"`
	Tol1minusZ float64 `json:"tol1minusZ"`
	Tol1plusD  float64 `json:"tol1plusD"`
	Tol1plusX  float64 `json:"tol1plusX"`
	Tol1plusY  float64 `json:"tol1plusY"`
	Tol1plusZ  float64 `json:"tol1plusZ"`
	Tol2minusD float64 `json:"tol2minusD"`
	Tol2minusX float64 `json:"tol2minusX"`
	Tol2minusY float64 `json:"tol2minusY"`
	Tol2minusZ float64 `json:"tol2minusZ"`
	Tol2plusD  float64 `json:"tol2plusD"`
	Tol2plusX  float64 `json:"tol2plusX"`
	Tol2plusY  float64 `json:"tol2plusY"`
	Tol2plusZ  float64 `json:"tol2plusZ"`
	Tol3minusD float64 `json:"tol3minusD"`
	Tol3minusX float64 `json:"tol3minusX"`
	Tol3minusY float64 `json:"tol3minusY"`
	Tol3minusZ float64 `json:"tol3minusZ"`
	Tol3plusD  float64 `json:"tol3plusD"`
	Tol3plusX  float64 `json:"tol3plusX"`
	Tol3plusY  float64 `json:"tol3plusY"`
	Tol3plusZ  float64 `json:"tol3plusZ"`

	OffsetX float64 `json:"offsetX"`
	OffsetY float64 `json:"offsetY"`
	OffsetZ float64 `json:"offsetZ"`

	P0  int `json:"p0"`
	P1  int `json:"p1"`
	P2  int `json:"p2"`
	P3  int `json:"p3"`
	P4  int `json:"p4"`
	P5  int `json:"p5"`
	P6  int `json:"p6"`
	P7  int `json:"p7"`
	P8  int `json:"p8"`
	P9  int `json:"p9"`
	P10 int `json:"p10"`
}
