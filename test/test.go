package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/bvisness/spall"
)

func main() {
	f, err := os.Create("test.spall")
	if err != nil {
		panic(err)
	}

	p := spall.NewProfile(f, spall.UnitMicroseconds)
	defer p.Close()
	e := p.NewEventer()
	defer e.Close()

	e.BeginNow("out")
	defer e.EndNow()

	for i := 0; i < 100; i++ {
		e.BeginNow("in loop")
		recurse(e, rand.Intn(25))
		e.EndNow()
	}
}

func recurse(e spall.Eventer, remaining int) {
	if remaining <= 0 {
		return
	}

	e.BeginNow("recurse")
	fmt.Printf("%d remanining\n", remaining)
	recurse(e, remaining-1)
	e.EndNow()
}
