package main

import (
	"flag"
	"log"
)

func main() {
	w := flag.Int("w", 40, "board width")
	h := flag.Int("h", 23, "board height")
	f := flag.String("f", "", "pattern filename")
	d := flag.Duration("d", 100, "duration of the screen refresh period in milliseconds")
	flag.Parse()

	a, err := newApp(*w, *h, *f, *d)
	if err != nil {
		log.Fatal(err)
	}

	eventC := make(chan event)

	go a.waitEvent(eventC)
	a.doEvent(eventC)
}
