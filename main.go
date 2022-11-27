package main

import (
	"flag"
	"log"
	"os"

	"github.com/lamber92/go-cover/internal/convert"
)

func usage() {
	log.Printf("Usage:\n\n\tgo-convert command [arguments]\n\n")
	log.Printf("The command is:\n\n")
	log.Printf("\tconvert ${cover.profile} --report ${report-mode}\n")
	log.Printf("\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	command := ""
	if flag.NArg() > 0 {
		command = flag.Arg(0)
		switch command {
		case "convert":
			if flag.NArg() <= 1 {
				log.Printf("Missing convert profile\n")
				os.Exit(1)
			}
			if err := convert.Convert(flag.Args()[1]); err != nil {
				log.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		default:
			log.Printf("Unknown command: %#q\n\n", command)
			usage()
		}
	} else {
		usage()
	}
}
