package main

import (
	"flag"
	"log"
)

func main() {
	w := flag.Int("w", 40, "board width")
	h := flag.Int("h", 23, "board height")
	f := flag.String("f", "", "pattern filename")
	p := flag.Duration("p", 100, "screen refresh period in milliseconds")
	flag.Parse()

	a, err := newApp(*w, *h, *f, *p)
	if err != nil {
		log.Fatal(err)
	}

	eventC := make(chan event)

	go a.waitEvent(eventC)
	a.doEvent(eventC)
}
