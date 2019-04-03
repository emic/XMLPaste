package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var indent = "    "
var version string

const (
	HelpTextTemplate = `Usage: xmlpaste [-h] [-n] [-o] [-p] [-v]

Description:
    XMLPaste is the command line tool to paste XML text of
    FileMaker clipboard objects.

    The following options are available:

    -h, --help
    Print a brief help message.

    -n
    Do not print the trailing newline character. (macOS only)

    -o, --output <file>
    Write output to <file> instead of stdout.

    -p, --pretty
    Format the XML content.

    -v, --version
    Display version information.

Author:
    Emic Corporation <https://www.emic.co.jp/>

License:
    This software is distributed under MIT License.
`
)

type cli struct {
	outStream, errStream io.Writer
}

func main() {
	cli := &cli{outStream: os.Stdout, errStream: os.Stderr}
	os.Exit(cli.Run(os.Args))
}

func (c *cli) Run(args []string) int {
	helpFlag := false
	nFlag := false
	filePath := ""
	prettyFlag := false
	versionFlag := false

	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flags.Usage = func() {}

	flags.BoolVar(&helpFlag, "h", false, "Print help pages.")
	flags.BoolVar(&helpFlag, "help", false, "Print help pages.")
	flags.BoolVar(&nFlag, "n", false, "Do not print the trailing newline character.")
	flags.StringVar(&filePath, "o", "", "Write output to <file> instead of stdout.")
	flags.StringVar(&filePath, "output", "", "Write output to <file> instead of stdout.")
	flags.BoolVar(&prettyFlag, "p", false, "Format the XML content.")
	flags.BoolVar(&prettyFlag, "pretty", false, "Format the XML content.")
	flags.BoolVar(&versionFlag, "v", false, "Print version information.")
	flags.BoolVar(&versionFlag, "version", false, "Print version information.")

	buf := &bytes.Buffer{}
	flags.SetOutput(buf)

	err := flags.Parse(args[1:])
	if err != nil {
		return 1
	}

	if helpFlag == true {
		fmt.Fprint(c.outStream, HelpTextTemplate)
	} else if versionFlag == true {
		fmt.Fprintln(c.outStream, "XMLPaste "+version)
	} else {
		str, err := getClipboard()
		if err != nil {
			return 1
		}

		if str != "" {
			if (filePath != "" || runtime.GOOS == "darwin") && strings.HasPrefix(str, "<?xml version=\"1.0\" encoding=\"utf-16\"?>") {
				str = strings.Replace(str, "<?xml version=\"1.0\" encoding=\"utf-16\"?>", "<?xml version=\"1.0\" encoding=\"UTF-8\"?>", 1)
			}

			if prettyFlag == true {
				str, _ = formatXML(str)
			}

			if filePath != "" {
				path, err := filepath.Abs(filePath)
				if err != nil {
					panic(err)
				}
				file, err := os.Create(path)
				if err != nil {
					panic(err)
				}
				defer file.Close()
				file.Write(([]byte)(str))
			} else {
				if nFlag == true {
					fmt.Fprint(c.outStream, str)
				} else {
					fmt.Fprintln(c.outStream, str)
				}
			}
		}
	}

	return 0
}

func formatXML(str string) (string, error) {
	re := regexp.MustCompile(`>\s+<`)
	src := re.ReplaceAllString(str, "><")
	reg := regexp.MustCompile(`<([/!]?)([^>]+?)(/?)>`)

	newLine := "\n"
	if runtime.GOOS == "windows" {
		newLine = "\r\n"
	}

	return (reg.ReplaceAllStringFunc(src, replace(newLine))), nil
}

func replace(newLine string) func(string) string {
	indentLevel := 0
	firstFlag := true
	newLineFlag := false
	return func(str string) string {
		defer func() {
			firstFlag = false
		}()

		if strings.HasPrefix(str, "<!") {
			return str
		}

		if strings.HasPrefix(str, "</") {
			if indentLevel > 0 {
				indentLevel--
			}
			if newLineFlag == true {
				return newLine + strings.Repeat(indent, indentLevel) + str
			}
			newLineFlag = true
			return str
		}

		if strings.HasPrefix(str, "<") || strings.HasSuffix(str, "/>") {
			newLineFlag = true
		}

		defer func() {
			if !strings.HasPrefix(str, "<?xml") && !strings.HasSuffix(str, "/>") {
				indentLevel++
			}
			if strings.HasSuffix(str, "/>") {
				newLineFlag = true
			}
		}()

		if newLineFlag == false {
			return strings.Repeat(indent, indentLevel) + str
		} else {
			newLineFlag = false
			if firstFlag == true {
				return str
			}
			return newLine + strings.Repeat(indent, indentLevel) + str
		}
	}
}
