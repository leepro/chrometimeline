package main

import (
	"flag"
	"log"

	"github.com/leepro/chrometimeline/timeline"
)

func main() {
	fn := flag.String("s", "./perf.json", "json file")
	output := flag.String("o", "./timeline.png", "output PNG")
	flag.Parse()

	err := timeline.Render(*fn, *output)
	if err != nil {
		log.Fatal(err)
	}
}
