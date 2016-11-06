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
	verbose       bool
	version       bool
	dryRun        bool
}

func init() {
	flag.StringVar(&args.inputFile, "file", "-", "input file, defaults to stdin")
	flag.BoolVar(&args.skipTLSVerify, "skip-tls-verify", false, "skips verification of the chain of certificate")
	flag.BoolVar(&args.verbose, "verbose", false, "print more verbose output")
	flag.BoolVar(&args.version, "version", false, "print version")
	flag.BoolVar(&args.dryRun, "dry-run", false, "don't send the report, this enables -verbose")
	flag.Parse()

	if args.inputFile == "-" {
		input = os.Stdin
	}

	if args.dryRun {
		args.verbose = true
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
		fmt.Println(report.String())
	}
	if !args.dryRun {
		r, err := newReporter(args.skipTLSVerify).send(report)
		if err != nil {
			fmt.Println(err)
			os.Exit(128)
		}
		fmt.Printf("Test coverage report sent\n")
		if args.verbose {
			r.Write(os.Stdout)
			io.WriteString(os.Stdout, "\n")
		}
	}
}
