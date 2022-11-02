package options

import "time"

type InLineOptions struct {
	Duration *time.Duration
	N        *int // 数组中的第n个元素；不是数组元素时此值为-1
}

func MergeInlineOptions(opts ...*InLineOptions) *InLineOptions {
	io := Inline()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.Duration != nil {
			io.Duration = opt.Duration
		}
		if opt.N != nil {
			io.N = opt.N
		}
	}
	return io
}

func Inline() *InLineOptions {
	return &InLineOptions{}
}
