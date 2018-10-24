package differ

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/kylelemons/godebug/pretty"
	"github.com/logrusorgru/aurora"
	isatty "github.com/mattn/go-isatty"
	yaml "gopkg.in/yaml.v2"
)

type Differ struct {
	File1   string
	File2   string
	NoColor bool
}

func NewDiffer(file1, file2 string, noColor bool) *Differ {
	return &Differ{
		File1:   file1,
		File2:   file2,
		NoColor: noColor,
	}
}

func (differ *Differ) ComputeDiff() string {

	formatter := newFormatter(differ.NoColor)

	errors := stat(differ.File1, differ.File2)
	failOnErr(formatter, errors...)

	yaml1, err := unmarshal(differ.File1)
	if err != nil {
		failOnErr(formatter, err)
	}
	yaml2, err := unmarshal(differ.File2)
	if err != nil {
		failOnErr(formatter, err)
	}

	diffs := make([]string, 0)
	for _, s := range strings.Split(pretty.Compare(yaml1, yaml2), "\n") {
		switch {
		case strings.HasPrefix(s, "+"):
			diffs = append(diffs, formatter.Bold(formatter.Green(s)).String())
		case strings.HasPrefix(s, "-"):
			diffs = append(diffs, formatter.Bold(formatter.Red(s)).String())
		}
	}
	return strings.Join(diffs, "\n")
}

func unmarshal(filename string) (interface{}, error) {
	contents, err := ioutil.ReadFile(filename)
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
