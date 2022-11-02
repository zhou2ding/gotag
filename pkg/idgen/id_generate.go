package idgen

import "sync"

type IdGenerator struct {
	id  int
	mtx sync.Mutex
}

var gIdGenerator *IdGenerator = &IdGenerator{
	id: 0,
}

func GetIdGenerator() *IdGenerator {
	return gIdGenerator
}

func (g *IdGenerator) GetId() int {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	g.id++
	if g.id < 0 {
		g.id = 1
	}
	return g.id
}
