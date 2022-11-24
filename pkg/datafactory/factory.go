package datafactory

import (
	"gotag/options"
	"gotag/pkg/datafactory/sample"
	"reflect"
)

type Factory interface {
	MakeData(rVal reflect.Value, rType reflect.Type, fieldTag map[string]string, fields []string, opts ...*options.Options)
}

func GetFactory(name string) Factory {
	var f Factory
	switch name {
	case "sample", "sample2":
		f = &sample.SampleFactory{Type: name}
	}
	return f
}
