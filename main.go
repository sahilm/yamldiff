package main

import (
	"fmt"
	"io"
	"os"

	"strings"

	"github.com/go-yaml/yaml"
	"github.com/jessevdk/go-flags"
	"github.com/kylelemons/godebug/pretty"
	"github.com/logrusorgru/aurora"
	"github.com/mattn/go-isatty"
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

	formatter := newFormatter(opts.NoColor)

	errors := stat(opts.File1, opts.File2)
	failOnErr(formatter, errors...)

	yaml1, err := unmarshal(opts.File1)
	if err != nil {
		failOnErr(formatter, err)
	}
	yaml2, err := unmarshal(opts.File2)
	if err != nil {
		failOnErr(formatter, err)
	}

	diff := computeDiffs(formatter, yaml1, yaml2)
	if diff != "" {
		fmt.Println(diff)
	}
}

func stat(filenames ...string) []error {
	var errs []error
	for _, filename := range filenames {
		_, err := os.Stat(filename)
		if err != nil {
			errs = append(errs, fmt.Errorf("cannot find file: %v. Does it exist?", filename))
		}
	}
	return errs
}

func unmarshal(filename string) ([]interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	d := yaml.NewDecoder(f)
	var ret []interface{}
	for {
		var r interface{}
		err = d.Decode(&r)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		ret = append(ret, r)
	}
	return ret, nil
}

func failOnErr(formatter aurora.Aurora, errs ...error) {
	if len(errs) == 0 {
		return
	}
	var errMessages []string
	for _, err := range errs {
		errMessages = append(errMessages, err.Error())
	}
	fmt.Fprintf(os.Stderr, "%v\n\n", formatter.Red(strings.Join(errMessages, "\n")))
	os.Exit(1)
}

func computeDiffs(formatter aurora.Aurora, a []interface{}, b []interface{}) string {
	diffs := make([]string, 0)
	if len(a) > len(b) {
		//extend - https://github.com/golang/go/wiki/SliceTricks
		b = append(b, make([]interface{}, len(a)-len(b))...)
	} else if len(b) > len(a) {
		a = append(a, make([]interface{}, len(b)-len(a))...)
	}

	for i := range a {
		diffs = append(diffs, computeDiff(formatter, a[i], b[i]), "---")
	}
	return strings.Join(diffs, "\n")
}

func computeDiff(formatter aurora.Aurora, a interface{}, b interface{}) string {
	diffs := make([]string, 0)
	for _, s := range strings.Split(pretty.Compare(a, b), "\n") {
		switch {
		//hack
		case (s == "+nil" || s == "-nil"):
			//drop this
		case strings.HasPrefix(s, "+"):
			diffs = append(diffs, formatter.Bold(formatter.Green(s)).String())
		case strings.HasPrefix(s, "-"):
			diffs = append(diffs, formatter.Bold(formatter.Red(s)).String())
		}
	}
	return strings.Join(diffs, "\n")
}

func newFormatter(noColor bool) aurora.Aurora {
	var formatter aurora.Aurora
	if noColor || !isTerminal() {
		formatter = aurora.NewAurora(false)
	} else {
		formatter = aurora.NewAurora(true)
	}
	return formatter
}

func isTerminal() bool {
	fd := os.Stdout.Fd()
	switch {
	case isatty.IsTerminal(fd):
		return true
	case isatty.IsCygwinTerminal(fd):
		return true
	default:
		return false
	}
}
