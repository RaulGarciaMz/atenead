package modelos

import (
	"sync"
)

type IntentosConcurrente struct {
	sync.RWMutex
	r map[int32]int
}

func (e *IntentosConcurrente) Set(eq int32, val int) {
	e.Lock()
	defer e.Unlock()

	e.r[eq] = val
}

func (e *IntentosConcurrente) Read(eq int32) (int, bool) {
	e.Lock()
	defer e.Unlock()

	v, ok := e.r[eq]

	if ok {
		return v, ok
	}

	return 0, ok
}

func (e *IntentosConcurrente) Inicializa() {
	e.Lock()
	defer e.Unlock()

	e.r = make(map[int32]int)
}
