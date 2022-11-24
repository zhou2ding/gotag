package sample

type Measure struct {
	Domain string   `json:"domain"`
	Key    string   `json:"key"`
	Id     int      `json:"id"`
	Value  MeaValue `json:"value"`
}

type MeaValue struct {
	PlanName   string      `json:"planName"`
	BodyNum    string      `json:"bodyNum"`
	Somethings []Something `json:"somethings"` // value type is convenient to reflect
}

type Something struct {
	// sth
}
