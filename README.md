# cloner
A package for folder cloning.

## Installation
Install the package to your *$GOPATH* with the [go get](https://golang.org/cmd/go/) tool from shell:<br/>
`$ go get -u github.com/galijot/cloner`

## Usage
* Import the package like<br/>
`import "github.com/galijot/cloner"`
* Send the source and destination paths to `cloner`'s `Clone` function<br/>
`cloner.Clone("path/to/src/dir/", "path/to/dst/dir/")`
* That's it 🎉

Also, the example project is available [here](https://github.com/galijot/cloner-example).