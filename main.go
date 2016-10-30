package main

import (
	"flag"
	"io"
	"log"
	"os"
)

var Version = "master"

var (
	input  io.Reader
	logger *log.Logger
)

var args struct {
	inputFile     string
	skipTLSVerify bool
}

func init() {
	flag.StringVar(&args.inputFile, "f", "-", "input file, defaults to stdin")
	flag.BoolVar(&args.skipTLSVerify, "S", false, "skips verification of the chain of certificate")
	flag.Parse()

	if args.inputFile == "-" {
		input = os.Stdin
	}
}

func main() {
	report, err := makeReport(input)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if err = newReporter(args.skipTLSVerify).send(report); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
