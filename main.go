package main

import (
	"encoding/base64"
	"encoding/json"
	"gotag/model"
	"gotag/options"
	"gotag/pkg/cache"
	_ "gotag/pkg/cache"
	"gotag/pkg/datafactory"
	"gotag/pkg/flags"
	"gotag/pkg/l"
	_ "gotag/pkg/l"
	"gotag/pkg/rpcclient"
	"gotag/pkg/v"
	_ "gotag/pkg/v"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	flags.InitFlag() // 直接调用初始化而非通过导入_的方式初始化，是为了避免单元测试函数报错

	// 校验参数
	ok := true
	for i, flag := range flags.AllFlags {
		if *flag == "" {
			l.GetLogger().Warn(flags.FailReason[i])
			ok = false
		}
	}
	if !ok {
		return
	}
	d, err := strconv.ParseInt(flags.Duration, 10, 64)
	if err != nil {
		l.GetLogger().Error("parse time duration failed", zap.Error(err))
		return
	}

	// 测量配置
	configTag, configFields, err := datafactory.GetFiledTagAndFields("data/config")
	if err != nil {
		l.GetLogger().Error("GetFiledTagAndFields of config failed", zap.Error(err))
		return
	}
	configCnt, err := datafactory.GetArrayCnt("data/config")
	if err != nil {
		l.GetLogger().Error("GetArrayCnt of config failed", zap.Error(err))
		return
	}
	// 构造数据
	config := model.Config{}
	config.Value.MeasurePoints = make([]model.ConfigMeasurePoint, configCnt)
	datafactory.MakeInLineData(reflect.ValueOf(&config).Elem(), reflect.TypeOf(config), configTag, configFields)

	configBytes, err := json.Marshal(&config)
	if err != nil {
		l.GetLogger().Error("Marshal config failed", zap.Error(err))
		return
	}
	// 测量配置数据写入缓存，供其他数据使用
	err = cache.GetCache().Set("gotag:inline:config", configBytes, 300)
	if err != nil {
		l.GetLogger().Error("SetCache of config failed", zap.Error(err))
		return
	}

	// 开始测量
	beginTag, beginFields, err := datafactory.GetFiledTagAndFields("data/begin")
	if err != nil {
		l.GetLogger().Error("GetFiledTagAndFields of begin failed", zap.Error(err))
		return
	}
	du := time.Duration(d)
	opts := options.InLineOptions{Duration: &du}

	begin := model.Begin{}
	datafactory.MakeInLineData(reflect.ValueOf(&begin).Elem(), reflect.TypeOf(begin), beginTag, beginFields, &opts)

	beginBytes, err := json.Marshal(&begin)
	if err != nil {
		l.GetLogger().Error("Marshal begin failed", zap.Error(err))
		return
	}
	err = cache.GetCache().Set("gotag:inline:begin", beginBytes, 300)
	if err != nil {
		l.GetLogger().Error("SetCache of begin failed", zap.Error(err))
		return
	}

	// 测量数据
	meaTag, meaFields, err := datafactory.GetFiledTagAndFields("data/measure")
	if err != nil {
		l.GetLogger().Error("GetFiledTagAndFields of measure failed", zap.Error(err))
		return
	}
	measureCnt, err := datafactory.GetArrayCnt("data/measure")
	if err != nil {
		l.GetLogger().Error("GetArrayCnt failed", zap.Error(err))
		return
	}

	measure := model.Measure{}
	measure.MeasurePoints = make([]model.MeasurePoint, measureCnt)
	datafactory.MakeInLineData(reflect.ValueOf(&measure).Elem(), reflect.TypeOf(measure), meaTag, meaFields, &opts)

	meaBytes, err := json.Marshal(&measure)
	if err != nil {
		l.GetLogger().Error("Marshal measure failed", zap.Error(err))
		return
	}
	err = cache.GetCache().Set("gotag:inline:measure", meaBytes, 300)
	if err != nil {
		l.GetLogger().Error("SetCache of measure failed", zap.Error(err))
		return
	}

	// 结束测量
	endTag, endFields, err := datafactory.GetFiledTagAndFields("data/end")
	if err != nil {
		l.GetLogger().Error("GetFiledTagAndFields of end failed", zap.Error(err))
		return
	}

	end := model.End{}
	datafactory.MakeInLineData(reflect.ValueOf(&end).Elem(), reflect.TypeOf(end), endTag, endFields, &opts)

	endBytes, err := json.Marshal(&end)
	if err != nil {
		l.GetLogger().Error("Marshal end failed", zap.Error(err))
		return
	}
	err = cache.GetCache().Set("gotag:inline:end", endBytes, 300)
	if err != nil {
		l.GetLogger().Error("SetCache of end failed", zap.Error(err))
		return
	}

	//fmt.Printf("config: %#v\nbegin: %#v\nmeasurement: %#v\nend: %#v\n", config, begin, measure, end)

	// 连接rpc服务端
	cli := rpcclient.CreateRPCClient(l.GetLogger())
	defer cli.Close()
	err = cli.Open(v.GetViper().GetString("rpcserver.host")+":"+v.GetViper().GetString("rpcserver.port"), 10000)
	if err != nil {
		l.GetLogger().Error("rpc open failed", zap.Error(err))
		return
	}
	serial := v.GetViper().GetString("inline.station_num")
	user, _ := json.Marshal(&rpcclient.UserInfo{
		FactoryName: v.GetViper().GetString("inline.factory_name"),
		StationName: v.GetViper().GetString("inline.station_name"),
		SerialNum:   serial,
	})
	random := "123456"

	userName := "inline%" + random + "%" + base64.StdEncoding.EncodeToString(user)
	err = cli.Login(userName, rpcclient.GeneratePwd(serial, random), 2000)
	if err != nil {
		l.GetLogger().Error("rpc login failed", zap.Error(err))
		return
	}

	reply, err := cli.Call(config.Domain, config.Key, configBytes, nil, 2000)
	if err != nil {
		l.GetLogger().Error("rpc call failed", zap.String("domain", config.Domain), zap.String("key", config.Key), zap.Error(err))
	}
	l.GetLogger().Info("rpc call done", zap.Any("reply", reply))

	consoleCh := make(chan os.Signal)
	signal.Notify(consoleCh, os.Interrupt, syscall.SIGTERM)
	select {
	case <-consoleCh:
		l.GetLogger().Info("receive interrupt signal from console")
	}
}
