package main

import (
	"io"
	"log"
	"os"

	"github.com/ogier/pflag"
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
	pflag.StringVarP(&args.inputFile, "input-file", "f", "-", "input file, defaults to stdin")
	pflag.BoolVarP(&args.skipTLSVerify, "insecure-skip-tls-verify", "S", false, "skips verification of the chain of certificate")
	pflag.Parse()

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
	if err = newReporter().send(report); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
