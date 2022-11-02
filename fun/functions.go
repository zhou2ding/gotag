package fun

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"go.uber.org/zap"
	"gotag/model"
	"gotag/pkg/l"
	"reflect"
	"strconv"
	"time"
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

func Sum[T ZdbType](p []T) T {
	var sum T
	for _, v := range p {
		sum += v
	}
	return sum
}

func Const[T ZdbType](p []T) T {
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

func RandStr(p []string) string {
	min, _ := strconv.Atoi(p[0])
	max, _ := strconv.Atoi(p[1])
	return strconv.Itoa(gofakeit.IntRange(min, max))
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

func Global[T ZdbType](p T) {
	// do nothing
}

// Time 形参切片的第一个元素为指定时间转换的字符串，第二个元素为需要的时间格式，eg: Time(time.Now().Format("2006-01-02 15:04:05"),"2006/01/02")
func Time(p []string) string {
	t, _ := time.Parse("2006-01-02 15:04:05", p[0])
	return t.Format(p[1])
}

// Foreign calls是调用链的字符串，indexName是根据索引在数组中查找的索引名字，indexVal是根据索引在数组中查找的索引的值，origin是根对象
func Foreign(calls []string, indexName string, indexVal, origin any) any {
	switch calls[0] {
	case "Config":
		config := origin.(model.Config)
		ret := getRefVal(calls[len(calls)-1], indexName, indexVal, reflect.ValueOf(config), reflect.TypeOf(config))
		if ret != nil {
			return ret
		}
	case "Begin":
		begin := origin.(model.Begin)
		ret := getRefVal(calls[len(calls)-1], indexName, indexVal, reflect.ValueOf(begin), reflect.TypeOf(begin))
		if ret != nil {
			return ret
		}
	case "End":
		end := origin.(model.End)
		ret := getRefVal(calls[len(calls)-1], indexName, indexVal, reflect.ValueOf(end), reflect.TypeOf(end))
		if ret != nil {
			return ret
		}
	case "Measurement":
		measure := origin.(model.Measure)
		ret := getRefVal(calls[len(calls)-1], indexName, indexVal, reflect.ValueOf(measure), reflect.TypeOf(measure))
		if ret != nil {
			return ret
		}
	}
	return nil
}

func This(field string, origin reflect.Value, rType reflect.Type) any {
	// this函数目前只支持单级调用，不支持根据索引在数组中查找
	ret := getRefVal(field, "", nil, origin, rType)
	return ret
}

// index是根据数组匹配时用到的索引，eg：Config.Value.MeasurePoints[PointName].TheoreticalX，index就是PointName
func getRefVal(callName, indexName string, indexVal any, rVal reflect.Value, rType reflect.Type) any {
	found := false
	for i := 0; i < rVal.NumField(); i++ {
		fieldVal := rVal.Field(i)
		l.GetLogger().Debug("getRefVal",
			zap.String("fieldName", rType.Field(i).Name),
			zap.Any("fieldVal", fieldVal),
			zap.String("indexName", indexName),
			zap.Any("indexVal", indexVal),
			zap.String("callName", callName),
			zap.Bool("found", found),
			zap.Bool("name equal", rType.Field(i).Name == indexName),
			zap.Bool("val equal", fieldVal.Interface() == indexVal),
		)
		switch fieldVal.Kind() {
		case reflect.String:
			if indexVal != nil {
				// 调用链中有根据索引在数组中查找的条件
				if !found && rType.Field(i).Name == indexName && fieldVal.Interface() == indexVal {
					found = true
				} else if found && rType.Field(i).Name == callName {
					return fieldVal.String()
				}
			} else {
				// 调用链中没有根据索引在数组中查找的条件
				if rType.Field(i).Name == callName {
					return fieldVal.String()
				}
			}
		case reflect.Int:
			if indexVal != nil {
				if !found && rType.Field(i).Name == indexName && fieldVal.Interface() == indexVal {
					found = true
				} else if found && rType.Field(i).Name == callName {
					return int(fieldVal.Int())
				}
			} else {
				if rType.Field(i).Name == callName {
					return int(fieldVal.Int())
				}
			}
		case reflect.Float32:
			if indexVal != nil {
				if !found && rType.Field(i).Name == indexName && fieldVal.Interface() == indexVal {
					found = true
				} else if found && rType.Field(i).Name == callName {
					return float32(fieldVal.Float())
				}
			} else {
				if rType.Field(i).Name == callName {
					return float32(fieldVal.Float())
				}
			}
		case reflect.Float64:
			if indexVal != nil {
				if !found && rType.Field(i).Name == indexName && fieldVal.Interface() == indexVal {
					found = true
				} else if found && rType.Field(i).Name == callName {
					return fieldVal.Float()
				}
			} else {
				if rType.Field(i).Name == callName {
					return fieldVal.Float()
				}
			}
		case reflect.Struct:
			return getRefVal(callName, indexName, indexVal, fieldVal, rType.Field(i).Type)
		case reflect.Slice:
			for j := 0; j < fieldVal.Len(); j++ {
				val := getRefVal(callName, indexName, indexVal, fieldVal.Index(j), fieldVal.Index(j).Type())
				if val != nil {
					return val
				}
			}
		}
	}
	return nil
}
