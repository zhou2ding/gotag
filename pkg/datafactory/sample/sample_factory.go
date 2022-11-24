package sample

import (
	"gotag/fun"
	"gotag/model/sample"
	"gotag/options"
	"gotag/pkg/cache"
	"gotag/pkg/datafactory/tools"
	"reflect"
	"strconv"
	"time"
)

type SampleFactory struct {
	Type string
}

func (c *SampleFactory) MakeData(rVal reflect.Value, rType reflect.Type, fieldTag map[string]string, fields []string, opts ...*options.Options) {
	var err error
	// 从上往下依次设置字段的值
	for _, fieldName := range fields {
		field := rVal.FieldByName(fieldName)
		funcNames := tools.GetFuncName(fieldTag[fieldName])
		zdbOpts := options.MergeOptions(opts...)
		var fieldVal any
		switch field.Kind() {
		case reflect.String:
			for _, name := range funcNames {
				// 从里到外调用函数，下同
				if name == "foreign" {
					// 外键函数比较特殊，不用泛型实现，下同
					callChain, indexName := tools.GetForeignParam(fieldTag[fieldName])
					switch callChain[0] {
					case "Config":
						config := sample.Config{}
						if c.Type == "sample" {
							err = cache.GetCache().Get("gotag:sample:config", &config.Value)
							if err != nil {
								panic(err)
							}
						} else if c.Type == "sample2" {
							err = cache.GetCache().Get("gotag:sample2:config", &config.Value)
							if err != nil {
								panic(err)
							}
						}
						fieldVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), config) // foreign函数只允许有一个参数
					case "Begin":
						begin := sample.Begin{}
						if c.Type == "sample" {
							err = cache.GetCache().Get("gotag:sample:begin", &begin.Value)
							if err != nil {
								panic(err)
							}
						} else if c.Type == "sample2" {
							err = cache.GetCache().Get("gotag:sample2:begin", &begin.Value)
							if err != nil {
								panic(err)
							}
						}
						fieldVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), begin)
					case "End":
						end := sample.End{}
						if c.Type == "sample" {
							err = cache.GetCache().Get("gotag:sample:end", &end.Value)
							if err != nil {
								panic(err)
							}
						} else if c.Type == "sample2" {
							err = cache.GetCache().Get("gotag:sample2:end", &end.Value)
							if err != nil {
								panic(err)
							}
						}
						fieldVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), end)
					case "Measurement":
						measure := sample.Measure{}
						if c.Type == "sample" {
							err = cache.GetCache().Get("gotag:sample:measure", &measure.Value)
							if err != nil {
								panic(err)
							}
						} else if c.Type == "sample2" {
							err = cache.GetCache().Get("gotag:sample2:measure", &measure.Value)
							if err != nil {
								panic(err)
							}
						}
						fieldVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), measure)
					}
				} else if name == "time" {
					// time函数比较特殊，配置文件只指定了格式，具体的时间在此处设置，下同
					var params []string
					if zdbOpts.Duration != nil {
						params = append(params, time.Now().Add(*zdbOpts.Duration*time.Second).Format("2006-01-02 15:04:05"))
					}
					params = append(params, tools.GetUnFuncParam(fieldTag[fieldName], name)...)
					fieldVal = fun.StringFuncMap[name](params)
				} else {
					var params []string
					if fieldVal != nil {
						params = []string{fieldVal.(string)}
					}
					subParam := tools.GetUnFuncParam(fieldTag[fieldName], name)
					if tools.Contains(subParam, "this") {
						for _, p := range subParam {
							params = append(params, fun.This(p[5:], rVal, rType).(string))
						}
					} else {
						params = append(params, subParam...)
					}
					if name == "incrByN" && zdbOpts.N != nil {
						// N是自增了几次
						params = append(params, strconv.Itoa(*zdbOpts.N))
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
					callChain, indexName := tools.GetForeignParam(fieldTag[fieldName])
					switch callChain[0] {
					case "Config":
						config := sample.Config{}
						if c.Type == "sample" {
							err = cache.GetCache().Get("gotag:sample:config", &config.Value)
							if err != nil {
								panic(err)
							}
						} else if c.Type == "sample2" {
							err = cache.GetCache().Get("gotag:sample2:config", &config.Value)
							if err != nil {
								panic(err)
							}
						}
						fieldVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), config) // foreign函数只允许有一个参数
					case "Begin":
						begin := sample.Begin{}
						if c.Type == "sample" {
							err = cache.GetCache().Get("gotag:sample:begin", &begin.Value)
							if err != nil {
								panic(err)
							}
						} else if c.Type == "sample2" {
							err = cache.GetCache().Get("gotag:sample2:begin", &begin.Value)
							if err != nil {
								panic(err)
							}
						}
						fieldVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), begin)
					case "End":
						end := sample.End{}
						if c.Type == "sample" {
							err = cache.GetCache().Get("gotag:sample:end", &end.Value)
							if err != nil {
								panic(err)
							}
						} else if c.Type == "sample2" {
							err = cache.GetCache().Get("gotag:sample2:end", &end.Value)
							if err != nil {
								panic(err)
							}
						}
						fieldVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), end)
					case "Measurement":
						measure := sample.Measure{}
						if c.Type == "sample" {
							err = cache.GetCache().Get("gotag:sample:measure", &measure.Value)
							if err != nil {
								panic(err)
							}
						} else if c.Type == "sample2" {
							err = cache.GetCache().Get("gotag:sample2:measure", &measure.Value)
							if err != nil {
								panic(err)
							}
						}
						fieldVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), measure)
					}
				} else {
					var params []int
					if fieldVal != nil {
						params = []int{fieldVal.(int)}
					}
					subParam := tools.GetUnFuncParam(fieldTag[fieldName], name)
					if tools.Contains(subParam, "this") {
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
					if name == "incrByN" && zdbOpts.N != nil {
						params = append(params, *zdbOpts.N)
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
					callChain, indexName := tools.GetForeignParam(fieldTag[fieldName])
					switch callChain[0] {
					case "Config":
						config := sample.Config{}
						if c.Type == "sample" {
							err = cache.GetCache().Get("gotag:sample:config", &config.Value)
							if err != nil {
								panic(err)
							}
						} else if c.Type == "sample2" {
							err = cache.GetCache().Get("gotag:sample2:config", &config.Value)
							if err != nil {
								panic(err)
							}
						}
						fieldVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), config) // foreign函数只允许有一个参数
					case "Begin":
						begin := sample.Begin{}
						if c.Type == "sample" {
							err = cache.GetCache().Get("gotag:sample:begin", &begin.Value)
							if err != nil {
								panic(err)
							}
						} else if c.Type == "sample2" {
							err = cache.GetCache().Get("gotag:sample2:begin", &begin.Value)
							if err != nil {
								panic(err)
							}
						}
						fieldVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), begin)
					case "End":
						end := sample.End{}
						if c.Type == "sample" {
							err = cache.GetCache().Get("gotag:sample:end", &end.Value)
							if err != nil {
								panic(err)
							}
						} else if c.Type == "sample2" {
							err = cache.GetCache().Get("gotag:sample2:end", &end.Value)
							if err != nil {
								panic(err)
							}
						}
						fieldVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), end)
					case "Measurement":
						measure := sample.Measure{}
						if c.Type == "sample" {
							err = cache.GetCache().Get("gotag:sample:measure", &measure.Value)
							if err != nil {
								panic(err)
							}
						} else if c.Type == "sample2" {
							err = cache.GetCache().Get("gotag:sample2:measure", &measure.Value)
							if err != nil {
								panic(err)
							}
						}
						fieldVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), measure)
					}
				} else {
					var params []float32
					if fieldVal != nil {
						params = []float32{fieldVal.(float32)}
					}
					subParam := tools.GetUnFuncParam(fieldTag[fieldName], name)
					if tools.Contains(subParam, "this") {
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
					if name == "incrByN" && zdbOpts.N != nil {
						params = append(params, float32(*zdbOpts.N))
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
					callChain, indexName := tools.GetForeignParam(fieldTag[fieldName])
					switch callChain[0] {
					case "Config":
						config := sample.Config{}
						if c.Type == "sample" {
							err = cache.GetCache().Get("gotag:sample:config", &config.Value)
							if err != nil {
								panic(err)
							}
						} else if c.Type == "sample2" {
							err = cache.GetCache().Get("gotag:sample2:config", &config.Value)
							if err != nil {
								panic(err)
							}
						}
						fieldVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), config)
					case "Begin":
						begin := sample.Begin{}
						if c.Type == "sample" {
							err = cache.GetCache().Get("gotag:sample:begin", &begin.Value)
							if err != nil {
								panic(err)
							}
						} else if c.Type == "sample2" {
							err = cache.GetCache().Get("gotag:sample2:begin", &begin.Value)
							if err != nil {
								panic(err)
							}
						}
						fieldVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), begin)
					case "End":
						end := sample.End{}
						if c.Type == "sample" {
							err = cache.GetCache().Get("gotag:sample:end", &end.Value)
							if err != nil {
								panic(err)
							}
						} else if c.Type == "sample2" {
							err = cache.GetCache().Get("gotag:sample2:end", &end.Value)
							if err != nil {
								panic(err)
							}
						}
						fieldVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), end)
					case "Measurement":
						measure := sample.Measure{}
						if c.Type == "sample" {
							err = cache.GetCache().Get("gotag:sample:measure", &measure.Value)
							if err != nil {
								panic(err)
							}
						} else if c.Type == "sample2" {
							err = cache.GetCache().Get("gotag:sample2:measure", &measure.Value)
							if err != nil {
								panic(err)
							}
						}
						fieldVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), measure)
					}
				} else {
					var params []float64
					if fieldVal != nil {
						params = []float64{fieldVal.(float64)}
					}
					subParam := tools.GetUnFuncParam(fieldTag[fieldName], name)
					if tools.Contains(subParam, "this") {
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
					if name == "incrByN" && zdbOpts.N != nil {
						params = append(params, float64(*zdbOpts.N))
					}
					fieldVal = fun.Float64FuncMap[name](params)
				}
			}
			field.SetFloat(fieldVal.(float64))
		case reflect.Struct:
			subType, _ := rType.FieldByName(fieldName)
			c.MakeData(field, subType.Type, fieldTag, fields, zdbOpts)
		case reflect.Slice:
			var sliceFiledVal []string // todo 目前只有[]string类型的；后续如果有了其他切片类型，再考虑此处该如何修改
			if field.Len() > 0 {
				if field.Index(0).Kind() == reflect.Struct {
					for j := 0; j < field.Len(); j++ {
						zdbOpts.N = &j
						c.MakeData(field.Index(j), field.Index(j).Type(), fieldTag, fields, zdbOpts)
					}
				} else {
					// 非结构体数组，无需递归
					for j := 0; j < field.Len(); j++ {
						var subFiledVal any
						switch field.Index(j).Kind() {
						case reflect.String:
							for _, name := range funcNames {
								// 从里到外调用函数，下同
								if name == "foreign" {
									// 外键函数比较特殊，不用泛型实现，下同
									callChain, indexName := tools.GetForeignParam(fieldTag[fieldName])
									switch callChain[0] {
									case "Measurement":
										measure := sample.Measure{}
										err = cache.GetCache().Get("gotag:sample:measure", &measure.Value)
										if err != nil {
											panic(err)
										}
										subFiledVal = fun.Foreign(callChain, indexName, fun.This(indexName, rVal, rType), measure)
									}
								} else if name == "time" {
									// time函数比较特殊，配置文件只指定了格式，具体的时间在此处设置，下同
									var params []string
									if zdbOpts.Duration != nil {
										params = append(params, time.Now().Add(*zdbOpts.Duration*time.Second).Format("2006-01-02 15:04:05"))
									}
									params = append(params, tools.GetUnFuncParam(fieldTag[fieldName], name)...)
									subFiledVal = fun.StringFuncMap[name](params)
								} else {
									var params []string
									if subFiledVal != nil {
										params = []string{subFiledVal.(string)}
									}
									subParam := tools.GetUnFuncParam(fieldTag[fieldName], name)
									if tools.Contains(subParam, "this") {
										for _, p := range subParam {
											params = append(params, fun.This(p[5:], rVal, rType).(string))
										}
									} else {
										params = append(params, subParam...)
									}
									if name == "incrByN" && zdbOpts.N != nil {
										// N是自增了几次
										params = append(params, strconv.Itoa(*zdbOpts.N))
									}
									subFiledVal = fun.StringFuncMap[name](params)
								}
							}
							sliceFiledVal = append(sliceFiledVal, subFiledVal.(string))
							// todo 其他类型出现这些类型了再补充
						}
					}
					field.Set(reflect.ValueOf(sliceFiledVal))
				}
			}
		}
	}
}
