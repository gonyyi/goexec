// (c) Gon Y. Yi 2021-2022 <https://gonyyi.com/copyright>
// Last Update: 01/13/2022

package goexec_test

import (
	"testing"

	"github.com/gonyyi/goexec"
)

func TestNewPool(t *testing.T) {
	p := goexec.NewPool()

	{
		e1 := p.Get()
		e2 := p.Get()
		println(1, "e1", e1.ID()) // 1
		println(1, "e2", e2.ID()) // 2
		p.Put(e1)
	}

	{
		e1 := p.Get()
		e2 := p.Get()
		println(2, "e1", e1.ID()) // 1 -- reusing from 1.e1
		println(2, "e2", e2.ID()) // 3 -- new since 1.e2 didn't get returned.
		p.Put(e1)
		p.Put(e2)
	}
	println("New:", p.Count())
}
