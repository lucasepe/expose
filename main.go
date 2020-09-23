//go:generate statik -src=./assets
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lucasepe/expose/slides"
)

const (
	banner = ` ____  _  _  ____   __   ____  ____ 
(  __)( \/ )(  _ \ /  \ / ___)(  __) v{{VERSION}}
 ) _)  )  (  ) __/(  O )\___ \ ) _) 
(____)(_/\_)(__)   \__/ (____/(____)`
)

var version = "0.1.0"

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() == 0 {
		usage()
	}

	printBanner()
	if err := slides.Serve(flag.Args()[0]); err != nil {
		exitOnErr(err)
	}
}

func usage() {
	printBanner()

	fmt.Fprintf(os.Stderr, "Markdown Driven Slide Presentations Viewer.\n\n")

	fmt.Fprintf(os.Stderr, "USAGE:\n\n")
	fmt.Fprintf(os.Stderr, "  %s /path/to/your/file.md\n\n", appName())

	fmt.Fprintf(os.Stderr, "crafted with passion by Luca Sepe - https://github.com/lucasepe/expose\n")
	os.Exit(0)
}

func printBanner() {
	str := strings.Replace(banner, "{{VERSION}}", version, 1)
	fmt.Fprintf(os.Stderr, str)
	fmt.Fprintf(os.Stderr, "\n\n")
}

func appName() string {
	return filepath.Base(os.Args[0])
}

// exitOnErr check for an error and eventually exit
func exitOnErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}
