package fun

type ZdbType interface {
	ZdbNumber | ~string
}

type ZdbNumber interface {
	~int | ~float32 | ~float64
}
