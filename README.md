# mossgow
mossgow is a simple project which cleans up students assignments and submits them into moss.
### Build
`make build`
### Install
`make install` 
### Run
```
detect software similarity 

Usage:
mossgow detect [flags]

Flags:
-b, --base string         To define common code file
-h, --help                help for detect
-i, --input string        To define input zip file (default "uploads.zip")
-l, --languages strings   To define supported languages (default [.go,.py,.java,.c,.cpp,.cs,.js])
-m, --moss string         To define path to moss (default "moss")
-p, --pathlayers int      To define path layers (default 3)
```