# cloner
A package for folder cloning.

## Installation
Install the package to your *$GOPATH* with the [go get](https://golang.org/cmd/go/) tool from shell:<br/>
`$ go get -u github.com/galijot/cloner`

## Usage
* Import the package like<br/>
`import "github.com/galijot/cloner"`
* Create a `struct` & implement `CloneOptions` interface.

```
type options struct {
	includeHiddenItems bool
}

func (o options) IncludeHidden() bool {
	return o.includeHiddenItems
}
```

* Send the source and destination paths along options to `cloner`'s `Clone` function<br/>
`cloner.Clone("path/to/src/dir/", "path/to/dst/dir/", options{false})`
* That's it 🎉

Also, the example project is available [here](https://github.com/galijot/cloner-example).

## Sub-packages
cloner includes few additional standalone sub-packages;

* [cper](https://github.com/galijot/cloner/tree/master/cper) - used for directory & file copying
* [flogger](https://github.com/galijot/cloner/tree/master/flogger) - used for logs, which are written to a file at provided path