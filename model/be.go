package model

type Begin struct {
	Domain string  `json:"domain"`
	Key    string  `json:"key"`
	Id     int     `json:"id"`
	Value  BEValue `json:"value"` // value type is convenient to reflect
}

type End struct {
	Domain string  `json:"domain"`
	Key    string  `json:"key"`
	Id     int     `json:"id"`
	Value  BEValue `json:"value"` // value type is convenient to reflect
}

type BEValue struct {
	PlanName string `json:"planName"`
	BodyNum  string `json:"bodyNum"`
	Date     string `json:"data"`
	Time     string `json:"time"`
	Operator string `json:"operator,omitempty"`
}
