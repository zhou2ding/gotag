package datafactory

import (
	"reflect"
	"testing"
)

func TestGetUnFuncParam(t *testing.T) {
	type test struct {
		input1 string
		input2 string
		want   []string
	}

	tests := map[string]*test{
		"test1": {"`bor:toStrf(incrByN(1,3),A%03d)`", "toStrf", []string{"A%03d"}},

		"test2": {`bor:funcA(funcB(funcC(funcD(d1,d2),c1,c2),b1,b2),a1,a2)`, "funcD", []string{"d1", "d2"}},
		"test3": {`tag:funcA(funcB(b1,b2),funcC(c1,c2),a1,a2)`, "funcA", []string{"a1", "a2"}},

		"test4": {`tag:funcA(funcB(b1,b2),funcC(c1,c2),a1,a2)`, "funcA", []string{"a1", "a2"}},
		"test5": {`tag:funcA(funcB(b1,b2),funcC(c1,c2),a1,a2)`, "funcB", []string{"b1", "b2"}},
		"test6": {`tag:funcA(funcB(b1,b2),funcC(c1,c2),a1,a2)`, "funcC", []string{"c1", "c2"}},
		"test7": {`bor:sum(foreign(End.Value.BodyNum),foreign(End1.Value1.PointName1))`, "sum", nil},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := getUnFuncParam(tt.input1, tt.input2)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got: %v, want: %v\n", got, tt.want)
			}
		})
	}
}

func TestGetFuncName(t *testing.T) {
	type test struct {
		input string
		want  []string
	}

	tests := map[string]*test{
		"test1": {`tag:funcA(a1)`, []string{"funcA"}},
		"test2": {`tag:funcA(funcB(b1,b2),a1,a2)`, []string{"funcB", "funcA"}},

		"test3": {`tag:funcA(funcB(funcC(c1,c2),funcD(d1,d2),b1,b2),a1,a2)`, []string{"funcC", "funcD", "funcB", "funcA"}},
		"test4": {`tag:funcA(funcB(b1,b2),funcC(c1,c2),a1,a2)`, []string{"funcB", "funcC", "funcA"}},
		"test5": {`tag:funcA(funcB(b1,b2),funcB(b1,b2),a1,a2)`, []string{"funcB", "funcB", "funcA"}},
		"test6": {`bor:sum(foreign(End.Value.BodyNum),foreign(End.Value.PointName))`, []string{"foreign", "foreign", "sum"}},
		"test7": {`bor:const(common)`, []string{"const"}},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := getFuncName(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got: %v, want: %v\n", got, tt.want)
			}
		})
	}
}

func TestGetCharacterIndex(t *testing.T) {
	type test struct {
		input1 string
		input2 rune
		want   []int
	}

	tests := map[string]*test{
		"test1": {`tag:funcA(funcB(funcC(funcD(d1,d2),c1,c2),b1,b2),a1,a2)`, '(', []int{9, 15, 21, 27}},
		"test2": {`tag:funcA(funcB(funcC(funcD(d1,d2),c1,c2),b1,b2),a1,a2)`, ')', []int{33, 40, 47, 54}},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := getCharacterIndex(tt.input1, tt.input2)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got: %v, want: %v\n", got, tt.want)
			}
		})
	}
}

func TestGetForeignParam(t *testing.T) {
	type test struct {
		input string
		want1 []string
		want2 string
	}

	tests := map[string]*test{
		"test1": {`bor:foreign(config.value.measurePoints[A2].theoreticalX)`, []string{"config", "value", "measurePoints[A2]", "theoreticalX"}, "A2"},
		"test2": {`bor:sum(foreign(config.value.measurePoints.theoreticalY),dy)`, []string{"config", "value", "measurePoints", "theoreticalY"}, ""},
		"test3": {`bor:sum(foreign(End.Value.BodyNum),rand(-1000,1000))`, []string{"End", "Value", "BodyNum"}, ""},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got1, got2 := getForeignParam(tt.input)
			if !reflect.DeepEqual(got1, tt.want1) || got2 != tt.want2 {
				t.Errorf("got1: %v, want1: %v, got2: %v, want2: %v\n", got1, tt.want1, got2, tt.want2)
			}
		})
	}
}

func TestGetFirstInNormalCharacterIndex(t *testing.T) {
	type test struct {
		input string
		want  int
	}

	tests := map[string]*test{
		"test1": {`tag:funcA(funcB`, 9},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := getFirstInNormalCharacterIndex(tt.input)
			if got != tt.want {
				t.Errorf("got: %v, want: %v\n", got, tt.want)
			}
		})
	}
}
