package main

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
)

func main() {

	runtime.GOMAXPROCS(1)
	c := container{}
	c.instances = new(sync.Map)
	m := sync.Map{}
	m.Store("kk", "vv")
	c.RegisterInstance(&m)
	var wg sync.WaitGroup

	f := func() {
		defer wg.Done()
		//wg.Add(1)
		fmt.Println("started f")
		for {
			// fmt.Printf("%s\n",value)
			m.Store("kkk", "vvv")
			m.Load("kkk")
		}
	}
	/*f2 := func() {
		defer wg.Done()
		wg.Add(1)
		fmt.Println("started f2")
		for {
			m3 := sync.Map{}
			c.MutableInstance(&m3)
			m3.Load("k1")
			m3.Store("kk", "vv")
		}
	}*/

	for i := 0; i < 10; i++ {
		go f()
		//go f2()
		wg.Add(1)
	}

	mm := sync.Map{}
	c.MutableInstance(&mm)
	mm.Load("kkl1")
	mm.Store("kk", "vv")
	wg.Wait()
	fmt.Println("end")

}

// container
type container struct {
	// instances
	instances *sync.Map
}

// RegisterInstance
func (c *container) RegisterInstance(val any) (err error) {
	valR := reflect.ValueOf(val)
	typPtr := valR.Type()

	if typPtr.Kind() != reflect.Ptr {
		return fmt.Errorf("instance must ptr type, error type:%s", typPtr)
	}

	typ := typPtr.Elem()

	if _, has := c.instances.Load(typ); has {
		return fmt.Errorf("register duplicate instance type:%s", typ)
	}

	c.instances.Store(typ, valR.Elem().Interface())
	return
}

// MutableInstance get object from container
func (c *container) MutableInstance(val any) bool {
	valr := reflect.ValueOf(val)
	ins, has := c.instances.Load(valr.Type().Elem())
	if has && ins != nil {
		valr.Elem().Set(reflect.ValueOf(ins))
	}

	return has
}
