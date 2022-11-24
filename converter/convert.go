package converter

import (
	"fmt"
	"os"
	"strings"
)

var (
	errInvalidBlogFormat = fmt.Errorf("invalid blog format")
)

/*
Docusaurus Blog Format

---
slug: welcome
title: Welcome
authors: [slorber, yangshun]
tags: [facebook, hello, docusaurus]
---

The blog post date can be extracted from filenames, such as:

- `2019-05-30-welcome.md`
- `2019-05-30-welcome/index.md`

A blog post folder can be convenient to co-locate blog post images:

![Docusaurus Plushie](./docusaurus-plushie-banner.jpeg)

*/
type DocusaurusBlog struct {
}

/*
Hexo Blog Format

---
title: host-static-blog-on-aws-s3
date: 2022-06-28 09:16:48
tags: AWS, S3, Hexo
---

## Second step: Create bucket on aws s3

### create a bucket on s3

{% asset_img create_bucket_step1.jpg create bucket on s3 %}
{% asset_img create_bucket_step2.jpg create bucket on s3 %}

*/
type HexoBlog struct {
	SlugTitle string
	Title     string
	Date      string   // 2022-06-28
	Tags      []string //
	Content   string   // blog content
}

func ReadHexoBlogs(dir string) ([]*HexoBlog, Imgs []string, error) {
	fs, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var blogFiles []string
	var blogDirs []string
	for _, f := range fs {
		if f.IsDir() {
			blogDirs = append(blogDirs, f.Name())
			continue
		}

		blogFiles = append(blogFiles, f.Name())
	}

	fmt.Printf("blogDirs: %v\n", blogDirs)
	fmt.Printf("blogFiles: %v\n", blogFiles)

	var hexoBlogs []*HexoBlog
	var errs []string
	for _, name := range blogFiles {
		blog, err := extractHexoBlog(name, dir+"/"+name)
		if err != nil {
			fmt.Printf("extractHexoBlog err: %v, name: %s", err, name)
			errs = append(errs, err.Error())
			continue
		}
		hexoBlogs = append(hexoBlogs, blog)
	}

	//todo: collect all imgs
	return hexoBlogs, nil, nil
}

func extractHexoBlog(name, blogPath string) (*HexoBlog, error) {
	fileContentBs, err := os.ReadFile(blogPath)
	if err != nil {
		return nil, err
	}

	fileContent := string(fileContentBs)
	blog, err := extractHeaderAndContent(fileContent)
	if err != nil {
		return nil, err
	}
	blog.SlugTitle = name
	return blog, nil
}

func extractHeaderAndContent(fileContent string) (*HexoBlog, error) {
	sep := "---"
	headerIdx := strings.Index(fileContent, sep)
	if headerIdx == -1 {
		return nil, errInvalidBlogFormat
	}
	lastHeaderIdx := strings.LastIndex(fileContent, sep)
	if lastHeaderIdx == -1 {
		return nil, errInvalidBlogFormat
	}

	start := headerIdx + len(sep)
	header := fileContent[start:lastHeaderIdx]
	blog := parseHeader(header)

	start = lastHeaderIdx + len(sep) + len("\n")
	blog.Content = fileContent[start:]
	return blog, nil
}

func parseHeader(header string) *HexoBlog {
	hds := strings.Split(header, "\n")
	var title, date string
	var tags []string
	/*
		---
		title: host-static-blog-on-aws-s3
		date: 2022-06-28 09:16:48
		tags: AWS, S3, Hexo
		---
	*/
	for _, hd := range hds {
		hd = strings.TrimSpace(hd)

		if strings.HasPrefix(hd, "title") {
			title = strings.Split(hd, ":")[1]
		}
		if strings.HasPrefix(hd, "date") {
			title = strings.Split(hd, ":")[1]
		}
		if strings.HasPrefix(hd, "tags") {
			tagsStr := strings.Split(hd, ":")[1]
			tags = strings.Split(tagsStr, ", ")
		}
	}

	return &HexoBlog{
		Title: title,
		Date:  date,
		Tags:  tags,
	}
}

func collectImgsFromBlogDir(blogDirs []string) ([]string, error) {
	return nil, nil
}

func replaceImg(concent string) string {
	return ""
}

func generateDocusaurusBlog(hblog *HexoBlog) *DocusaurusBlog {
	return nil
}

func WriteDocusaurusBlogs([]*DocusaurusBlog) error {
	return nil
}
