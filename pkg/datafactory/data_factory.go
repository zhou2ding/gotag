package datafactory

import (
	"gotag/fun"
	"gotag/model"
	"gotag/options"
	"gotag/pkg/cache"
	"reflect"
	"strconv"
	"time"
)

func MakeInLineData(rVal reflect.Value, rType reflect.Type, fieldTag map[string]string, fields []string, opts ...*options.InLineOptions) {
	var err error
	// 从上往下依次设置字段的值
	for _, fieldName := range fields {
		field := rVal.FieldByName(fieldName)
		funcNames := getFuncName(fieldTag[fieldName])
		inOpts := options.MergeInlineOptions(opts...)
		var fieldVal any
		switch field.Kind() {
		case reflect.String:
			for _, name := range funcNames {
				// 从里到外调用函数，下同
				if name == "foreign" {
					// 外键函数比较特殊，不用泛型实现，下同
					callChain, indexName := getForeignParam(fieldTag[fieldName])
					switch callChain[0] {
					case "Config":
						config := model.Config{}
						err = cache.GetCache().Get("gotag:inline:config", &config)
						if err != nil {
							panic(err)
						}
						fieldVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), config) // foreign函数只允许有一个参数
					case "Begin":
						begin := model.Begin{}
						err = cache.GetCache().Get("gotag:inline:begin", &begin)
						if err != nil {
							panic(err)
						}
						fieldVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), begin)
					case "End":
						end := model.End{}
						err = cache.GetCache().Get("gotag:inline:end", &end)
						if err != nil {
							panic(err)
						}
						fieldVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), end)
					case "Measurement":
						measure := model.Measure{}
						err = cache.GetCache().Get("gotag:inline:measure", &measure)
						if err != nil {
							panic(err)
						}
						fieldVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), measure)
					}
				} else if name == "time" {
					// time函数比较特殊，配置文件只指定了格式，具体的时间在此处设置，下同
					var params []string
					if inOpts.Duration != nil {
						params = append(params, time.Now().Add(*inOpts.Duration*time.Second).Format("2006-01-02 15:04:05"))
					}
					params = append(params, getUnFuncParam(fieldTag[fieldName], name)...)
					fieldVal = fun.StringFuncMap[name](params)
				} else {
					var params []string
					if fieldVal != nil {
						params = []string{fieldVal.(string)}
					}
					subParam := getUnFuncParam(fieldTag[fieldName], name)
					if contains(subParam, "this") {
						for _, p := range subParam {
							params = append(params, fun.This(p[5:], rVal, rType).(string))
						}
					} else {
						params = append(params, subParam...)
					}
					if name == "incrByN" && inOpts.N != nil {
						// N是自增了几次
						params = append(params, strconv.Itoa(*inOpts.N))
					}
					fieldVal = fun.StringFuncMap[name](params)
				}
			}
			field.SetString(fieldVal.(string))
		case reflect.Int:
			for _, name := range funcNames {
				// 从里到外调用函数
				if name == "foreign" {
					// 外键函数比较特殊，不用泛型实现
					callChain, indexName := getForeignParam(fieldTag[fieldName])
					switch callChain[0] {
					case "Config":
						config := model.Config{}
						err = cache.GetCache().Get("gotag:inline:config", &config)
						if err != nil {
							panic(err)
						}
						fieldVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), config) // foreign函数只允许有一个参数
					case "Begin":
						begin := model.Begin{}
						err = cache.GetCache().Get("gotag:inline:begin", &begin)
						if err != nil {
							panic(err)
						}
						fieldVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), begin)
					case "End":
						end := model.End{}
						err = cache.GetCache().Get("gotag:inline:end", &end)
						if err != nil {
							panic(err)
						}
						fieldVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), end)
					case "Measurement":
						measure := model.Measure{}
						err = cache.GetCache().Get("gotag:inline:measure", &measure)
						if err != nil {
							panic(err)
						}
						fieldVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), measure)
					}
				} else {
					var params []int
					if fieldVal != nil {
						params = []int{fieldVal.(int)}
					}
					subParam := getUnFuncParam(fieldTag[fieldName], name)
					if contains(subParam, "this") {
						for _, p := range subParam {
							params = append(params, fun.This(p[5:], rVal, rType).(int))
						}
					} else {
						for _, p := range subParam {
							pInt, err := strconv.Atoi(p)
							if err != nil {
								panic(err)
							}
							params = append(params, pInt)
						}
					}
					if name == "incrByN" && inOpts.N != nil {
						params = append(params, *inOpts.N)
					}
					fieldVal = fun.IntFuncMap[name](params)
				}
			}
			field.SetInt(int64(fieldVal.(int)))
		case reflect.Float32:
			for _, name := range funcNames {
				// 从里到外调用函数
				if name == "foreign" {
					// 外键函数比较特殊，不用泛型实现
					callChain, indexName := getForeignParam(fieldTag[fieldName])
					switch callChain[0] {
					case "Config":
						config := model.Config{}
						err = cache.GetCache().Get("gotag:inline:config", &config)
						if err != nil {
							panic(err)
						}
						fieldVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), config) // foreign函数只允许有一个参数
					case "Begin":
						begin := model.Begin{}
						err = cache.GetCache().Get("gotag:inline:begin", &begin)
						if err != nil {
							panic(err)
						}
						fieldVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), begin)
					case "End":
						end := model.End{}
						err = cache.GetCache().Get("gotag:inline:end", &end)
						if err != nil {
							panic(err)
						}
						fieldVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), end)
					case "Measurement":
						measure := model.Measure{}
						err = cache.GetCache().Get("gotag:inline:measure", &measure)
						if err != nil {
							panic(err)
						}
						fieldVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), measure)
					}
				} else {
					var params []float32
					if fieldVal != nil {
						params = []float32{fieldVal.(float32)}
					}
					subParam := getUnFuncParam(fieldTag[fieldName], name)
					if contains(subParam, "this") {
						for _, p := range subParam {
							params = append(params, fun.This(p[5:], rVal, rType).(float32))
						}
					} else {
						for _, p := range subParam {
							pInt, err := strconv.ParseFloat(p, 32)
							if err != nil {
								panic(err)
							}
							params = append(params, float32(pInt))
						}
					}
					if name == "incrByN" && inOpts.N != nil {
						params = append(params, float32(*inOpts.N))
					}
					fieldVal = fun.Float32FuncMap[name](params)
				}
			}
			field.SetFloat(float64(fieldVal.(float32)))
		case reflect.Float64:
			for _, name := range funcNames {
				// 从里到外调用函数
				if name == "foreign" {
					// 外键函数比较特殊，不用泛型实现
					callChain, indexName := getForeignParam(fieldTag[fieldName])
					switch callChain[0] {
					case "Config":
						config := model.Config{}
						err = cache.GetCache().Get("gotag:inline:config", &config)
						if err != nil {
							panic(err)
						}
						fieldVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), config)
					case "Begin":
						begin := model.Begin{}
						err = cache.GetCache().Get("gotag:inline:begin", &begin)
						if err != nil {
							panic(err)
						}
						fieldVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), begin)
					case "End":
						end := model.End{}
						err = cache.GetCache().Get("gotag:inline:end", &end)
						if err != nil {
							panic(err)
						}
						fieldVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), end)
					case "Measurement":
						measure := model.Measure{}
						err = cache.GetCache().Get("gotag:inline:measure", &measure)
						if err != nil {
							panic(err)
						}
						fieldVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), measure)
					}
				} else {
					var params []float64
					if fieldVal != nil {
						params = []float64{fieldVal.(float64)}
					}
					subParam := getUnFuncParam(fieldTag[fieldName], name)
					if contains(subParam, "this") {
						for _, p := range subParam {
							params = append(params, fun.This(p[5:], rVal, rType).(float64))
						}
					} else {
						for _, p := range subParam {
							pFloat64, err := strconv.ParseFloat(p, 64)
							if err != nil {
								panic(err)
							}
							params = append(params, pFloat64)
						}
					}
					if name == "incrByN" && inOpts.N != nil {
						params = append(params, float64(*inOpts.N))
					}
					fieldVal = fun.Float64FuncMap[name](params)
				}
			}
			field.SetFloat(fieldVal.(float64))
		case reflect.Struct:
			subType, _ := rType.FieldByName(fieldName)
			MakeInLineData(field, subType.Type, fieldTag, fields, inOpts)
		case reflect.Slice:
			for j := 0; j < field.Len(); j++ {
				inOpts.N = &j
				MakeInLineData(field.Index(j), field.Index(j).Type(), fieldTag, fields, inOpts)
			}
		}
	}
}
