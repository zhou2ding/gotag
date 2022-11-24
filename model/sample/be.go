package sample

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
	// sth
}
