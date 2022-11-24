package converter

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	errInvalidBlogFormat = fmt.Errorf("invalid blog format")
)

const docusaurusTmpl = `
---
slug: %s
title: %s
authors: [%s]
tags: %v
---

%s
`

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

// DocusaurusBlog ...
// Use Date-SlugTitle as the blog dir name
type DocusaurusBlog struct {
	SlugTitle string
	Date      string
	Content   string // rendered by template
	Imgs      []string
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
	Imgs      []string // imgs used by this blog
}

// ReadHexoBlogs read blogs and collect imgs under the blog
func ReadHexoBlogs(dir string) ([]*HexoBlog, []string, error) {
	fs, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, err
	}

	var blogFiles []string
	var blogDirs []string
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), ".DS_Store") {
			continue
		}

		if f.IsDir() {
			blogDirs = append(blogDirs, dir+"/"+f.Name())
			continue
		}

		if strings.HasSuffix(f.Name(), ".md") {
			blogFiles = append(blogFiles, f.Name())
		}
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

	if len(errs) > 0 {
		return nil, nil, fmt.Errorf("%s", strings.Join(errs, ", "))
	}

	imgs, err := collectImgsFromBlogDir(dir, blogDirs)
	if err != nil {
		fmt.Printf("collectImgsFromBlogDir err: %v\n", err)
		return nil, nil, err
	}

	return hexoBlogs, imgs, nil
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
	blog.Imgs = extractImgFromContent(blog.Content)
	return blog, nil
}

func extractImgFromContent(content string) []string {
	exp := `asset_img\s+(.*\.jpg)\s+`
	regp, err := regexp.Compile(exp)
	if err != nil {
		fmt.Printf("compile %s err: %v\n", exp, err)
		return nil
	}

	var imgs []string
	results := regp.FindAllString(content, -1)
	for _, rs := range results {
		// asset_img create_bucket_step2.jpg
		img := strings.Split(strings.TrimSpace(rs), " ")[1]
		imgs = append(imgs, img)
	}

	return imgs
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
	fmt.Printf("headers: %v\n", hds)

	for _, hd := range hds {
		hd = strings.TrimSpace(hd)

		switch {
		case strings.Contains(hd, "title"):
			title = strings.Split(hd, ": ")[1]
		case strings.Contains(hd, "date"):
			date = strings.Split(hd, ": ")[1]
		case strings.Contains(hd, "tags"):
			tagsStr := strings.Split(hd, ": ")[1]
			tags = strings.Split(tagsStr, ", ")
		}
	}

	return &HexoBlog{
		Title: title,
		Date:  date,
		Tags:  tags,
	}
}

func collectImgsFromBlogDir(blogPath string, blogDirs []string) ([]string, error) {
	var images = []string{}
	for _, dir := range blogDirs {
		fs, err := os.ReadDir(dir)
		if err != nil {
			return nil, err
		}

		for _, f := range fs {
			if strings.HasSuffix(f.Name(), ".DS_Store") {
				continue
			}

			if !f.IsDir() {
				images = append(images, blogPath+"/"+f.Name())
			}
		}
	}

	return images, nil
}

func GenerateDocusaurusBlogs(blogs []*HexoBlog) []*DocusaurusBlog {
	var results []*DocusaurusBlog
	for _, hb := range blogs {
		dblog := generateDocusaurusBlog(hb)
		results = append(results, dblog)
	}
	return results
}

func getImgName(p string) string {
	imgs := extractImgFromContent(p)
	if len(imgs) > 0 {
		return imgs[0]
	}

	return ""
}

func replaceImg(content string) (string, error) {
	exp := `{%\s+asset_img\s+(.*\.jpg).*%}`
	regp, err := regexp.Compile(exp)
	if err != nil {
		return "", err
	}

	results := regp.FindAllString(content, -1)
	fmt.Printf("matched strings: %v\n", results)

	replaceMap := map[string]string{}
	for _, rs := range results {
		imgName := getImgName(rs)
		replaceMap[rs] = fmt.Sprintf("![%s](./%s)", imgName, imgName)
	}

	replaced := regp.ReplaceAllStringFunc(content, func(s string) string {
		return replaceMap[s]
	})
	return replaced, nil
}

func generateDocusaurusBlog(hblog *HexoBlog) *DocusaurusBlog {
	dblog := &DocusaurusBlog{
		SlugTitle: hblog.SlugTitle,
		Date:      strings.Split(hblog.Date, " ")[0], // only need year-month-day
		Content:   hblog.Content,
	}

	if len(hblog.Imgs) == 0 {
		return dblog
	}

	var err error
	dblog.Content, err = replaceImg(dblog.Content)
	if err != nil {
		fmt.Printf("replaceImg err: %v\n", err)
		return dblog
	}

	return dblog
}

func WriteDocusaurusBlogs([]*DocusaurusBlog) error {
	return nil
}
