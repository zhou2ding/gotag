package flags

import (
	"flag"
)

var (
	MeasureCnt           int
	MeasurePointInterval int64
	PulseCnt             int
	PulseInterVal        int

	//AllFlags = []any{&MeasureCnt, &MeasurePointInterval, &PulseCnt, &PulseInterVal}

	// FailReason key是flag在AllFlags切片中的索引，value是flag未指定时需要打印的原因
	//FailReason = map[int]string{
	//	0: "measure counts is not specified",
	//	1: "measure points time interval is not specified",
	//	2: "pulse counts is not specified",
	//	3: "pulse time interval is not specified",
	//}
)

func InitFlag() {
	flag.IntVar(&MeasureCnt, "n", 1, "实时上传的次数")
	flag.Int64Var(&MeasurePointInterval, "mt", 5, "测量数据中每个测点（MeasurePoint）间的时间间隔（Time），单位s")
	flag.IntVar(&PulseCnt, "pn", 1, "上传测量数据（pulse）的次数")
	flag.IntVar(&PulseInterVal, "pt", 1, "每次上传测量数据（pulse）的时间间隔")
	flag.Parse()
}
