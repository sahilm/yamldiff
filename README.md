# yamldiff

A CLI tool to diff two YAML/JSON files.

Nothing fancy about the code, all the heavy liftin' is done by:

* [go-yaml](https://github.com/go-yaml/yaml/) - for YAML parsin'
* [godebug](https://github.com/kylelemons/godebug/) - for diffin'
* [aurora](https://github.com/logrusorgru/aurora) - for fancy printin'
* [go-flags](https://github.com/jessevdk/go-flags) - for flaggin'
* [The Go stdlib](https://golang.org/pkg/) - for everythin'

Thanks to all the contributors of the above libraries.

## Installation

`go get -u github.com/sahilm/yamldiff`

## Usage

`yamldiff --file1 /path/to/yamlfile1.yml --file2 /path/to/yamlfile2.yml`

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
