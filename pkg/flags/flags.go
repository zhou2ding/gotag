package flags

import "flag"

var (
	Host     string
	Port     string
	User     string
	Pwd      string
	Duration string

	AllFlags = []*string{&Host, &Port, &User, &Pwd, &Duration}

	// FailReason key是flag在AllFlags切片中的索引，value是flag未指定时需要打印的原因
	FailReason = map[int]string{
		0: "host is not specified",
		1: "port is not specified",
		2: "user is not specified",
		3: "password is not specified",
		4: "time duration is not specified",
	}
)

func InitFlag() {
	flag.StringVar(&Host, "h", "", "连接rpc的IP")
	flag.StringVar(&Port, "p", "", "连接rpc的端口")
	flag.StringVar(&User, "u", "", "登录rpc的用户名")
	flag.StringVar(&Pwd, "pwd", "", "登录rpc的密码")
	flag.StringVar(&Duration, "d", "", "上传测量数据的时间间隔")
	flag.Parse()
}
