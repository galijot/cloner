package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/galijot/cloner"
)

type options struct {
	includeHiddenItems bool
}

func (o options) IncludeHidden() bool {
	return o.includeHiddenItems
}

func main() {
	src, dst := getPaths()
	err := cloner.Clone(src, dst, options{false})

	if err != nil {
		fmt.Println(err)
	}
}

// Reads the paths FROM and TO which a users wants to clone.
func getPaths() (string, string) {

	reader := bufio.NewScanner(os.Stdin)
	space := " "

	fmt.Println("Enter the path FROM which you want to clone:")
	reader.Scan()
	fromPath := strings.Trim(reader.Text(), space)

	fmt.Println("Enter the path TO which you wan to clone:")
	reader.Scan()
	toPath := strings.Trim(reader.Text(), space)

	return fromPath, toPath
}
