package main

import (
	"encoding/base64"
	"encoding/json"
	"gotag/model/sample"
	"gotag/options"
	"gotag/pkg/cache"
	_ "gotag/pkg/cache"
	"gotag/pkg/datafactory"
	"gotag/pkg/fieldtag"
	"gotag/pkg/flags"
	"gotag/pkg/l"
	_ "gotag/pkg/l"
	"gotag/pkg/rpcclient"
	"gotag/pkg/v"
	_ "gotag/pkg/v"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	flags.InitFlag() // 直接调用初始化而非通过导入_的方式初始化，是为了避免单元测试函数报错

	// 连接rpc服务端
	cli := rpcclient.CreateRPCClient(l.GetLogger())
	defer cli.Close()
	err := cli.Open(v.GetViper().GetString("rpcserver.host")+":"+v.GetViper().GetString("rpcserver.port"), 10000)
	if err != nil {
		l.GetLogger().Error("rpc open failed", zap.Error(err))
		return
	}
	serial := v.GetViper().GetString("sample.station_num")
	user, _ := json.Marshal(&rpcclient.UserInfo{
		Name: v.GetViper().GetString("sample.factory_name"),
		No:   serial,
	})

	random := "123456"
	userName := "sample%" + random + "%" + base64.StdEncoding.EncodeToString(user)

	err = cli.Login(userName, rpcclient.GeneratePwd(serial, random), 2000)
	if err != nil {
		l.GetLogger().Error("rpc login failed", zap.Error(err))
		return
	}

	// 测量配置
	configTag, configFields, err := fieldtag.GetFiledTagAndFields("data/sample/config")
	if err != nil {
		l.GetLogger().Error("GetFiledTagAndFields of config failed", zap.Error(err))
		return
	}
	configCnt, err := fieldtag.GetArrayCnts("data/sample/config", reflect.Struct)
	if err != nil {
		l.GetLogger().Error("GetArrayCnts of config failed", zap.Error(err))
		return
	}
	// 构造数据
	config := sample.Config{}
	config.Value.Somethings = make([]sample.ConfigSomething, configCnt["Somethings"])

	factory := datafactory.GetFactory("sample")
	factory.MakeData(reflect.ValueOf(&config).Elem(), reflect.TypeOf(config), configTag, configFields)

	configBytes, err := json.Marshal(&config.Value)
	if err != nil {
		l.GetLogger().Error("Marshal config failed", zap.Error(err))
		return
	}
	// 测量配置数据写入缓存，供其他数据使用
	err = cache.GetCache().Set("gotag:sample:config", configBytes, 3000)
	if err != nil {
		l.GetLogger().Error("SetCache of config failed", zap.Error(err))
		return
	}
	l.GetLogger().Debug("make config done", zap.Any("config", config))

	reply, err := cli.Call(config.Domain, config.Key, configBytes, nil, 2000)
	if err != nil {
		l.GetLogger().Error("rpc call failed", zap.String("domain", config.Domain), zap.String("key", config.Key), zap.Error(err))
	}
	l.GetLogger().Info("upload config done", zap.Any("reply", reply))

	for i := 0; i < flags.MeasureCnt; i++ {
		// 开始测量
		beginTag, beginFields, err := fieldtag.GetFiledTagAndFields("data/sample/begin")
		if err != nil {
			l.GetLogger().Error("GetFiledTagAndFields of begin failed", zap.Error(err))
			return
		}

		interval := time.Duration(flags.MeasurePointInterval)
		opts := options.Options{Duration: &interval}
		begin := sample.Begin{}

		factory.MakeData(reflect.ValueOf(&begin).Elem(), reflect.TypeOf(begin), beginTag, beginFields, &opts)

		beginBytes, err := json.Marshal(&begin.Value)
		if err != nil {
			l.GetLogger().Error("Marshal begin failed", zap.Error(err))
			return
		}
		err = cache.GetCache().Set("gotag:sample:begin", beginBytes, 3000)
		if err != nil {
			l.GetLogger().Error("SetCache of begin failed", zap.Error(err))
			return
		}
		l.GetLogger().Debug("make begin done", zap.Any("begin", begin))

		reply, err = cli.Call(begin.Domain, begin.Key, beginBytes, nil, 2000)
		if err != nil {
			l.GetLogger().Error("rpc call failed", zap.String("domain", begin.Domain), zap.String("key", begin.Key), zap.Error(err))
		}
		l.GetLogger().Info("upload begin done", zap.Any("reply", reply))

		// 测量数据
		meaTag, meaFields, err := fieldtag.GetFiledTagAndFields("data/sample/measure")
		if err != nil {
			l.GetLogger().Error("GetFiledTagAndFields of measure failed", zap.Error(err))
			return
		}
		measureCnt, err := fieldtag.GetArrayCnts("data/sample/measure", reflect.Struct)
		if err != nil {
			l.GetLogger().Error("GetArrayCnts failed", zap.Error(err))
			return
		}
		for j := 0; j < flags.PulseCnt; j++ {
			measure := sample.Measure{}
			measure.Value.Somethings = make([]sample.Something, measureCnt["Somethings"])

			factory.MakeData(reflect.ValueOf(&measure).Elem(), reflect.TypeOf(measure), meaTag, meaFields, &opts)

			meaBytes, err := json.Marshal(&measure.Value)
			if err != nil {
				l.GetLogger().Error("Marshal measure failed", zap.Error(err))
				return
			}
			err = cache.GetCache().Set("gotag:sample:measure", meaBytes, 3000)
			if err != nil {
				l.GetLogger().Error("SetCache of measure failed", zap.Error(err))
				return
			}
			l.GetLogger().Debug("make measure done", zap.Any("measure", measure))

			reply, err = cli.Call(measure.Domain, measure.Key, meaBytes, nil, 2000)
			if err != nil {
				l.GetLogger().Error("rpc call failed", zap.String("domain", measure.Domain), zap.String("key", measure.Key), zap.Error(err))
			}
			l.GetLogger().Info("upload pulse done", zap.Any("reply", reply))
			time.Sleep(time.Duration(flags.PulseInterVal) * time.Second)
		}

		// 结束测量
		endTag, endFields, err := fieldtag.GetFiledTagAndFields("data/sample/end")
		if err != nil {
			l.GetLogger().Error("GetFiledTagAndFields of end failed", zap.Error(err))
			return
		}

		end := sample.End{}
		factory.MakeData(reflect.ValueOf(&end).Elem(), reflect.TypeOf(end), endTag, endFields, &opts)

		endBytes, err := json.Marshal(&end.Value)
		if err != nil {
			l.GetLogger().Error("Marshal end failed", zap.Error(err))
			return
		}
		err = cache.GetCache().Set("gotag:sample:end", endBytes, 3000)
		if err != nil {
			l.GetLogger().Error("SetCache of end failed", zap.Error(err))
			return
		}
		l.GetLogger().Debug("make end done", zap.Any("end", end))

		reply, err = cli.Call(end.Domain, end.Key, endBytes, nil, 2000)
		if err != nil {
			l.GetLogger().Error("rpc call failed", zap.String("domain", end.Domain), zap.String("key", end.Key), zap.Error(err))
		}
		l.GetLogger().Info("upload end done", zap.Any("reply", reply))
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	select {
	case <-c:
		l.GetLogger().Info("receive interrupt signal from console")
	}
}
