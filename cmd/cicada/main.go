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
var flagVersion = flag.Bool("version", false, "Show version information")
var flagHelp = flag.Bool("help", false, "Show usage information")

func main() {
	flag.Parse()

	if *flagHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *flagVersion {
		fmt.Println(cicada.Version)
		os.Exit(0)
	}

	index, err := cicada.Load(*flagUpdate)

	if err != nil {
		log.Fatal(err)
	}

	if *flagDebug {
		index.Debug = true
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
