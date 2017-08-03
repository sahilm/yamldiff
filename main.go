package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/kylelemons/godebug/pretty"
	"github.com/logrusorgru/aurora"
	"gopkg.in/yaml.v2"
)

func main() {
	var opts struct {
		File1 string `long:"file1" description:"first YAML file" required:"true"`
		File2 string `long:"file2" description:"second YAML file" required:"true"`
	}

	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	errors := stat(opts.File1, opts.File2)
	failOnErr(errors...)

	yaml1, err := unmarshal(opts.File1)
	if err != nil {
		failOnErr(err)
	}
	yaml2, err := unmarshal(opts.File2)
	if err != nil {
		failOnErr(err)
	}

	diff := computeDiff(yaml1, yaml2)
	if diff != "" {
		fmt.Println(diff)
	} else {
		fmt.Println(aurora.Bold("no diff"))
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

func unmarshal(filename string) (interface{}, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var ret interface{}
	yaml.Unmarshal(contents, &ret)
	return ret, nil
}

func failOnErr(errs ...error) {
	if len(errs) > 0 {
		var errMessages []string
		for _, err := range errs {
			errMessages = append(errMessages, err.Error())
		}
		fmt.Fprintf(os.Stderr, "%v\n\n", aurora.Red(strings.Join(errMessages, "\n")))
		os.Exit(1)
	}
}

func computeDiff(a interface{}, b interface{}) string {
	var stringers []fmt.Stringer
	for _, s := range strings.Split(pretty.Compare(a, b), "\n") {
		switch {
		case strings.HasPrefix(s, "+"):
			stringers = append(stringers, aurora.Bold(aurora.Green(s)))
		case strings.HasPrefix(s, "-"):
			stringers = append(stringers, aurora.Bold(aurora.Red(s)))
		}
	}
	var s []string
	for _, stringer := range stringers {
		s = append(s, stringer.String())
	}
	return strings.Join(s, "\n")
}