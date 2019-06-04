package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"strings"

	"github.com/kylelemons/godebug/pretty"
	"github.com/logrusorgru/aurora"
	"github.com/mattn/go-isatty"
	"gopkg.in/yaml.v2"
)

var (
	version = "latest"

	noColorFlag = flag.Bool("no-color", false, "Disable colored output")
	versionFlag = flag.Bool("version", false, "Prints version and exit")
)

func main() {
	var (
		file1 string
		file2 string
	)

	flag.Parse()

	if *versionFlag {
		fmt.Println(version)
		return
	}

	formatter := newFormatter(*noColorFlag)

	args := flag.Args()
	if len(args) < 2 {
		failOnErr(formatter, errors.New("Files must be specified"))
	}
	file1 = args[0]
	file2 = args[1]

	errors := stat(file1, file2)
	failOnErr(formatter, errors...)

	yaml1, err := unmarshal(file1)
	if err != nil {
		failOnErr(formatter, err)
	}
	yaml2, err := unmarshal(file2)
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

func unmarshal(filename string) ([]interface{}, error) {
	var (
		contents []byte
		err      error
		values   []interface{}
	)
	if filename == "-" {
		contents, err = ioutil.ReadAll(os.Stdin)
	} else {
		contents, err = ioutil.ReadFile(filename)
	}
	if err != nil {
		return nil, err
	}
	dec := yaml.NewDecoder(bytes.NewReader(contents))
	for {
		var value interface{}
		err := dec.Decode(&value)
		if err == io.EOF {
			break
		}
		values = append(values, value)
	}
	return values, nil
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
	for i, s := range strings.Split(pretty.Compare(a, b), "\n") {
		pos := formatter.Gray(6, fmt.Sprintf("%03d.", i))
		switch {
		case strings.HasPrefix(s, "+"):
			diffs = append(diffs, fmt.Sprintf("%s %s", pos, formatter.Bold(formatter.Green(s)).String()))
		case strings.HasPrefix(s, "-"):
			diffs = append(diffs, fmt.Sprintf("%s %s", pos, formatter.Bold(formatter.Red(s)).String()))
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
