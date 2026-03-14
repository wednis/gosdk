package main

import (
	"fmt"

	"github.com/wednis/gosdk"
)

type dep_1 struct{}

func (dep *dep_1) say() {
	fmt.Println("dep 1")
}

type dep_2 struct{}

func (dep *dep_2) say() {
	fmt.Println("dep 2")
}

type dep_3 struct {
	a *dep_1
	b *dep_2
}

func (dep *dep_3) say() {
	fmt.Println("dep 3")
}

func fn_1(d1 *dep_1) {
	d1.say()
}

func fn_2(d2 *dep_2) {
	d2.say()
}

func fn_3(d1 *dep_1, d2 *dep_2) error {
	return nil
}

func fn_4(d1 *dep_1, d2 *dep_2) (*dep_3, error) {
	d3 := &dep_3{}
	d3.say()
	return d3, nil
}

func fn_5(d3 *dep_3) {
}

func test_depinject() {
	d1 := &dep_1{}
	d2 := &dep_2{}
	dic := gosdk.Inject(d1, d2, fn_1, fn_2, fn_3, fn_4, fn_5)
	if dic.Err() != nil {
		fmt.Println(dic.Err().Error())
		return
	}
	dic.Invoke(func(d1 *dep_1) {
		fmt.Println("invoke dep 1")
		d1.say()
	}, func(d1 *dep_1, d2 *dep_2, d3 *dep_3) {
		d3.a = d1
		d3.b = d2
		fmt.Println(d1 == d3.a)
		fmt.Println(d2 == d3.b)
	})
}

type config struct {
	Name string
	dev  bool
}

func main() {
	cfg := &config{}
	fmt.Println(cfg)
	gosdk.BindConfig(`/home/wednis/dev/project/go1.26.0/gosdk/test/test.json`, cfg)
	fmt.Println(cfg)
}
