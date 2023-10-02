package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/r3labs/diff/v3"

	"github.com/jessevdk/go-flags"
	"github.com/logrusorgru/aurora"
	"github.com/mattn/go-isatty"
	"gopkg.in/yaml.v2"
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

	diff := computeDiff(formatter, yaml1, yaml2)
	if diff != "" {
		fmt.Println(diff)
	}
}

func stat(filenames ...string) []error {
	var errs []error
	for _, filename := range filenames {
		if filename == "-" {
			continue
		}
		_, err := os.Stat(filename)
		if err != nil {
			errs = append(errs, fmt.Errorf("cannot find file: %v. Does it exist?", filename))
		}
	}
	return errs
}

func unmarshal(filename string) (interface{}, error) {
	var contents []byte
	var err error
	if filename == "-" {
		contents, err = ioutil.ReadAll(os.Stdin)
	} else {
		contents, err = ioutil.ReadFile(filename)
	}
	if err != nil {
		return nil, err
	}
	var ret interface{}
	err = yaml.Unmarshal(contents, &ret)
	if err != nil {
		return nil, err
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

func computeDiff(formatter aurora.Aurora, a interface{}, b interface{}) string {
	diffs := make([]string, 0)
	differ, err := diff.NewDiffer(diff.AllowTypeMismatch(true))
	if err != nil {
		return err.Error()
	}
	changelog, err := differ.Diff(a, b)
	if err != nil {
		return err.Error()
	}
	for _, s := range changelog {
		pathStr := formatter.Brown(strings.Join(s.Path, "."))
		fromStr := formatter.Red(fmt.Sprintf("- %v", s.From))
		toStr := formatter.Green(fmt.Sprintf("+ %v", s.To))
		chunk := fmt.Sprintf("%s:\n%s\n%s\n", pathStr, fromStr, toStr)
		diffs = insert(diffs, chunk)
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

func insert(ss []string, s string) []string {
	i := sort.SearchStrings(ss, s)
	ss = append(ss, "")
	copy(ss[i+1:], ss[i:])
	ss[i] = s
	return ss
}
