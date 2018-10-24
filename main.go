package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/thadc23/yamldiff/differ"
)

var version = "latest"

func main() {
	var opts struct {
		File1   string `long:"file1" description:"first YAML file" required:"true"`
		File2   string `long:"file2" description:"second YAML file" required:"true"`
		NoColor bool   `long:"no-color" description:"disable colored output" required:"false"`
		Version func() `long:"version" description:"print version and exit"`
	}

	opts.Version = func() {
		fmt.Fprintf(os.Stderr, "%v\n", version)
		os.Exit(0)
	}

	_, err := flags.Parse(&opts)
	if err != nil {
		if err.(*flags.Error).Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	d := differ.NewDiffer(opts.File1, opts.File2, opts.NoColor)

	diff := d.ComputeDiff()
	if diff != "" {
		fmt.Println(diff)
	}
}
