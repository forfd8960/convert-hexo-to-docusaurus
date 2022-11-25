package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/forfd8960/convert-hexo-to-docusaurus/converter"
)

func main() {
	hexoBlogDir := flag.String("hexo", "", "hexo blog directory")
	docBlogDir := flag.String("docusaurus", "", "docusaurus blog directory")
	author := flag.String("author", "", "blog author name")
	flag.Parse()

	if len(*hexoBlogDir) == 0 || len(*docBlogDir) == 0 {
		fmt.Fprintln(os.Stderr, "blog dir is empty")
		return
	}
	if len(*author) == 0 {
		fmt.Fprintln(os.Stderr, "author is empty")
		return
	}

	hexoBlogs, allImgs, err := converter.ReadHexoBlogs(*hexoBlogDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	docBlogs := converter.GenerateDocusaurusBlogs(*author, hexoBlogs)
	if err := converter.ExportDocusaurusBlogs(*docBlogDir, docBlogs, allImgs); err != nil {
		fmt.Fprintln(os.Stderr, "export DocusaurusBlogs has err: "+err.Error())
		return
	}

	fmt.Println("Convert Hexo Blogs to Docusaurus Blogs Success!")
}
