XMLPaste
=========
XMLPaste is a command line tool to paste XML text of FileMaker clipboard objects.


Usage
-----
The following options are available:
```
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
```


Example of using XMLPaste from AppleScript
-----
```
set xml to do shell script "/usr/local/bin/xmlpaste" as «class utf8»
set the clipboard to (xml as text)
```


System Requirements
-----
- macOS version   : macOS Big Sur 11, macOS Monterey 12, macOS Ventura 13 or macOS Sonoma 14
- Windows version : Windows 10 Version 22H2


Download
-----
Download from the [latest release page](https://github.com/emic/xmlpaste/releases/latest).


Author
-----
Emic Corporation <https://www.emic.co.jp/>


License
-----
This software is distributed under the MIT License.