package sample

type Config struct {
	Domain string      `json:"domain"`
	Key    string      `json:"key"`
	Id     int         `json:"id"`
	Value  ConfigValue `json:"value"`
}

type ConfigValue struct {
	ConfigName string            `json:"configName"`
	Somethings []ConfigSomething `json:"somethings"` // value type is convenient to reflect
}

type ConfigSomething struct {
	// sth
}
