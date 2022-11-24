package fun

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"gotag/model/sample"
	"reflect"
	"strconv"
	"time"
	"unicode"
)

var StringFuncMap = map[string]func(p []string) string{
	"sum":     Sum[string],
	"const":   Const[string],
	"toStrf":  ToStrf,
	"time":    Time,
	"rand":    RandStr,
	"incrByN": IncrByNStr,
}

var IntFuncMap = map[string]func(p []int) int{
	"sum":     Sum[int],
	"const":   Const[int],
	"incrByN": IncrByN[int],
	"rand":    Rand[int],
}

var Int64FuncMap = map[string]func(p []int64) int64{
	"sum":     Sum[int64],
	"const":   Const[int64],
	"incrByN": IncrByN[int64],
	"rand":    Rand[int64],
}

var Float32FuncMap = map[string]func(p []float32) float32{
	"sum":     Sum[float32],
	"const":   Const[float32],
	"incrByN": IncrByN[float32],
	"rand":    Rand[float32],
}

var Float64FuncMap = map[string]func(p []float64) float64{
	"sum":     Sum[float64],
	"const":   Const[float64],
	"incrByN": IncrByN[float64],
	"rand":    Rand[float64],
}

func Sum[T ZdbBaseType](p []T) T {
	var sum T
	for _, v := range p {
		sum += v
	}
	return sum
}

func Const[T ZdbBaseType](p []T) T {
	return p[0]
}

// IncrByN 形参的第一个元素为要执行增加的数字，第二个元素是该数字增加多少，第三个元素是已经加了几次
func IncrByN[T ZdbNumber](p []T) T {
	ret := p[0] + p[1]
	if len(p) > 2 {
		ret += p[2] * p[1]
	}
	return ret
}

// IncrByNStr 形参的第一个元素为要执行增加的数字，第二个元素是该数字增加多少，第三个元素是已经加了几次
func IncrByNStr(p []string) string {
	num, _ := strconv.Atoi(p[0])
	step, _ := strconv.Atoi(p[1])
	ret := num + step
	if len(p) > 2 {
		idx, _ := strconv.Atoi(p[2])
		ret += idx * step
	}
	return strconv.Itoa(ret)
}

// Rand 形参的第一个元素为随机范围的下限，第二个元素为随机范围的上限
func Rand[T ZdbNumber](p []T) T {
	i := interface{}(p[0])
	switch i.(type) {
	case int:
		return T(gofakeit.IntRange(int(p[0]), int(p[1])))
	case float32:
		return T(gofakeit.Float32Range(float32(p[0]), float32(p[1])))
	case float64:
		return T(gofakeit.Float64Range(float64(p[0]), float64(p[1])))
	default:
		return 0
	}
}

// RandStr 形参的第一个元素（只能是数字或字母）为随机范围的下限，第二个元素（只能是数字或字母）为随机范围的上限
func RandStr(p []string) string {
	if unicode.IsLetter(rune(p[0][0])) {
		var s string
		for i := 0; i < len(p[0]); i++ {
			s += gofakeit.RandomString(getLetters(string(p[0][i]), string(p[1][i])))
		}
		return s
	} else if unicode.IsDigit(rune(p[0][0])) {
		min, _ := strconv.Atoi(p[0])
		max, _ := strconv.Atoi(p[1])
		return strconv.Itoa(gofakeit.IntRange(min, max))
	}
	return ""
}

// ToStrf 形参切片的最后一个元素为结果字符串的格式，剩余的元素为要转成指定格式的字符串
func ToStrf(p []string) string {
	strs := p[:len(p)-1]
	var s string
	for _, str := range strs {
		s += str
	}
	return fmt.Sprintf(p[len(p)-1], s)
}

func Array[T ZdbArray](n int, p T) {

}

// Time 形参切片的第一个元素为指定时间转换的字符串，第二个元素为需要的时间格式，eg: Time(time.Now().Format("2006-01-02 15:04:05"),"2006/01/02")
func Time(p []string) string {
	t, _ := time.Parse("2006-01-02 15:04:05", p[0])
	return t.Format(p[1])
}

// Foreign calls是调用链的字符串，indexName是根据索引在数组中查找的索引名字，indexVal是根据索引在数组中查找的索引的值，origin是根对象
func Foreign(calls []string, indexName string, indexVal, origin any) any {
	switch calls[0] {
	case "Sample1":
		switch origin.(type) {
		case sample.Config:
			config := origin.(sample.Config)
			ret := getRefVal(calls[len(calls)-1], indexName, indexVal, reflect.ValueOf(config), reflect.TypeOf(config))
			if ret != nil {
				return ret
			}
		}
	case "Sample2":
		switch origin.(type) {
		case sample.Begin:
			begin := origin.(sample.Begin)
			ret := getRefVal(calls[len(calls)-1], indexName, indexVal, reflect.ValueOf(begin), reflect.TypeOf(begin))
			if ret != nil {
				return ret
			}
		}
	case "End":
		switch origin.(type) {
		case sample.End:
			end := origin.(sample.End)
			ret := getRefVal(calls[len(calls)-1], indexName, indexVal, reflect.ValueOf(end), reflect.TypeOf(end))
			if ret != nil {
				return ret
			}
		}
	case "Measurement":
		switch origin.(type) {
		case sample.Measure:
			measure := origin.(sample.Measure)
			ret := getRefVal(calls[len(calls)-1], indexName, indexVal, reflect.ValueOf(measure), reflect.TypeOf(measure))
			if ret != nil {
				return ret
			}
		}
	}
	return nil
}

func This(field string, origin reflect.Value, rType reflect.Type) any {
	// this函数目前只支持单级调用，不支持根据索引在数组中查找
	ret := getRefVal(field, "", nil, origin, rType)
	return ret
}
