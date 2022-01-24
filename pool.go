// (c) Gon Y. Yi 2021-2022 <https://gonyyi.com/copyright>
// Last Update: 01/13/2022

package goexec

import (
	"sync"
)

const poolBaseID = 1000

func NewPool() *Pool {
	b := &Pool{
		count: poolBaseID,
	}
	b.pool.New = func() interface{} {
		b.count += 1
		return &Exec{
			id: b.count,
		}
	}
	return b
}

// **************************************************
// Exec Pool
// **************************************************

type Pool struct {
	pool      sync.Pool
	count     int
	AutoReset bool
}

// Count returns how many has been created
func (p *Pool) Count() int {
	return p.count - poolBaseID
}

// Get will take *Exec from the pool
func (p *Pool) Get() *Exec {
	e := p.pool.Get().(*Exec)
	if p.AutoReset {
		e.Reset()
	}
	return e
}

// Put release *Exec to the pool
func (p *Pool) Put(e *Exec) {
	p.pool.Put(e)
}
