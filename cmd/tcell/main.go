package main

import (
	"flag"
	"log"
)

func main() {
	f := flag.String("f", "", "pattern filename")
	d := flag.Duration("d", 100, "duration of the screen refresh period in milliseconds")
	flag.Parse()

	a, err := newApp(*f, *d, 2)
	if err != nil {
		log.Fatal(err)
	}

	e := make(chan event)
	ep := make(chan eventPoint)

	go a.waitEvent(e, ep)
	a.doEvent(e, ep)
}
