package main

import (
	"fmt"
	"os"
)

// build time variables
var (
	AppName      string
	BuildDate    string
	BuildVersion string
	BuildSHA     string
)

func main() {
	code, err := run()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s: %s\n", AppName, err)
	}

	os.Exit(code)
}

func run() (int, error) {

	cli := setup(start)
	err := cli.Execute()
	if err != nil {
		return 1, err
	}

	return 0, nil
}
