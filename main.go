//go:generate statik -src=./assets
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/lucasepe/expose/slides"
)

const (
	banner = ` ____  _  _  ____   __   ____  ____ 
(  __)( \/ )(  _ \ /  \ / ___)(  __) v{{VERSION}}
 ) _)  )  (  ) __/(  O )\___ \ ) _) 
(____)(_/\_)(__)   \__/ (____/(____) Markdown Driven Slides Viewer`
)

var version = "0.2.1"

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() == 0 {
		usage()
	}

	printBanner()

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-stop
		cancel()
	}()

	filename := flag.Args()[0]
	_, name := filepath.Split(filename)
	nameNoExt := strings.TrimSuffix(name, filepath.Ext(name))

	exposer, err := slides.Expose(filename)
	if err != nil {
		exitOnErr(err)
	}

	fmt.Printf("Open your web browser and visit '%s'\n\n", exposer.URL())
	fmt.Printf("You can run Chrome in application mode:\n")
	fmt.Printf(" * Linux  : google-chrome --app=%s\n", exposer.URL())
	fmt.Printf(" * Windows: chrome --app=%s\n", exposer.URL())

	fmt.Printf("\nTo export a presentation as a PDF:\n")
	fmt.Printf(" * Linux  : google-chrome --headless --disable-gpu --print-to-pdf='%s.pdf' '%s'\n", nameNoExt, exposer.URL())
	fmt.Printf(" * Windows: chrome --headless --disable-gpu --print-to-pdf='%s.pdf' '%s'\n", nameNoExt, exposer.URL())

	if err := exposer.Serve(ctx); err != nil {
		exitOnErr(err)
	}
}

func usage() {
	printBanner()

	fmt.Fprintf(os.Stderr, "USAGE:\n\n")
	fmt.Fprintf(os.Stderr, "  %s /path/to/your/file.md\n\n", appName())

	fmt.Fprintf(os.Stderr, "Crafted with passion by Luca Sepe - https://github.com/lucasepe/expose\n")
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
