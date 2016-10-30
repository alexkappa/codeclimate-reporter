package main

import (
	"encoding/json"
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
	verbose       bool
	version       bool
}

func init() {
	flag.StringVar(&args.inputFile, "f", "-", "input file, defaults to stdin")
	flag.BoolVar(&args.skipTLSVerify, "S", false, "skips verification of the chain of certificate")
	flag.BoolVar(&args.verbose, "v", false, "print more verbose output")
	flag.BoolVar(&args.version, "V", false, "print version")
	flag.Parse()

	if args.inputFile == "-" {
		input = os.Stdin
	}

	errors.PrintTrace = false
}

func main() {
	if args.version {
		fmt.Print(Version)
		os.Exit(0)
	}
	report, err := makeReport(input)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if args.verbose {
		b, err := json.MarshalIndent(report, "  ", "  ")
		if err != nil {
			fmt.Println(err)
			os.Exit(128)
		}
		fmt.Printf("Test coverage report:\n%s\n", b)
	}
	if err = newReporter(args.skipTLSVerify).send(report); err != nil {
		fmt.Println(err)
		os.Exit(128)
	}
	fmt.Println("Test coverage report sent")
}
