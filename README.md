# yamldiff
[![darwin/linux build status](https://travis-ci.org/sahilm/yamldiff.svg?branch=master)](https://travis-ci.org/sahilm/yamldiff)
[![Go Report Card](https://goreportcard.com/badge/github.com/sahilm/yamldiff)](https://goreportcard.com/report/github.com/sahilm/yamldiff)

A CLI tool to diff two YAML/JSON files.

Nothing fancy about the code, all the heavy liftin' is done by:

* [go-yaml](https://github.com/go-yaml/yaml/) - for YAML parsin'
* [godebug](https://github.com/kylelemons/godebug/) - for diffin'
* [aurora](https://github.com/logrusorgru/aurora) - for fancy printin'
* [go-isatty](https://github.com/mattn/go-isatty) - for tty detectin'
* [go-flags](https://github.com/jessevdk/go-flags) - for flaggin'
* [The Go stdlib](https://golang.org/pkg/) - for everythin'

Thanks to all the contributors of the above libraries.

## Installation

* Download a variant of `yamldiff-v$VERSION-{darwin,linux,windows}-amd64` from the [releases](https://github.com/sahilm/yamldiff/releases) page.
* Rename the downloaded file to something sane like `yamldiff` :)
* Mark the file as an executable. On *nix, `chmod +x yamldiff`.
* Put it on your `$PATH`.

## Usage

`yamldiff --file1 /path/to/yamlfile1.yml --file2 /path/to/yamlfile2.yml`. The output is colorized by default. Colors
can be suppressed by the `--no-color` flag. Colors will automatically be suppressed if `stdout` is not a `tty`, for example
when piping/redirecting the output of `yamldiff`.

## Contributing

* Pull the code: `go get -u github.com/sahilm/yamldiff`.
* Hack away!
* Send a pull request.
* Have fun.

## License

The MIT License (MIT)

Copyright (c) 2017 Sahil Muthoo

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
