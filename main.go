package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/alexkappa/errors"
)

var Version = "master"

var input io.Reader

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

	errors.PrintTrace = false
}

func main() {
	report, err := makeReport(input)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err = newReporter(args.skipTLSVerify).send(report); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Test coverage data sent")
}
