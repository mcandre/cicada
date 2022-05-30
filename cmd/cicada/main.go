package main

import (
	"github.com/mcandre/cicada"

	"flag"
	"fmt"
	"log"
	"os"
)

var flagDebug = flag.Bool("debug", false, "Enable additional logging")
var flagUpdate = flag.Bool("update", false, "Force LTS index cache update")

func main() {
	flag.Parse()

	index, err := cicada.Load(*flagUpdate)

	if err != nil {
		log.Fatal(err)
	}

	if *flagDebug {
		log.Printf("Index: %v\n", index)
	}

	warnings, err := index.Scan()

	if err != nil {
		log.Fatal(err)
	}

	for _, warning := range warnings {
		fmt.Printf("warning: %v\n", warning)
	}

	if len(warnings) > 0 {
		os.Exit(1)
	}
}
