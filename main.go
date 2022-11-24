package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("load current dir err: ", err)
		return
	}

	blogDir := flag.String("", cwd, "hexo blog directory")
	if len(*blogDir) == 0 {
		fmt.Fprintln(os.Stderr, "blog dir is empty")
		return
	}

	fmt.Println("blog directoty: ", *blogDir)
}
