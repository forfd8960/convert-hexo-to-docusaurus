package converter

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var (
	errInvalidBlogFormat = fmt.Errorf("invalid blog format")
)

const docusaurusTmpl = `---
slug: %s
title: %s
authors: [%s]
tags: [%s]
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
func ReadHexoBlogs(blogsDir string) ([]*HexoBlog, []string, error) {
	fs, err := os.ReadDir(blogsDir)
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
			blogDirs = append(blogDirs, blogsDir+"/"+f.Name())
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
		lastIdx := strings.LastIndex(name, ".")
		slugTitle := name[:lastIdx]
		blog, err := extractHexoBlog(slugTitle, blogsDir+"/"+name)
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

	imgs, err := collectImgsFromBlogDir(blogDirs)
	if err != nil {
		fmt.Printf("collectImgsFromBlogDir err: %v\n", err)
		return nil, nil, err
	}

	return hexoBlogs, imgs, nil
}

func extractHexoBlog(slugTitle, blogPath string) (*HexoBlog, error) {
	fileContentBs, err := os.ReadFile(blogPath)
	if err != nil {
		return nil, err
	}

	fileContent := string(fileContentBs)
	blog, err := extractHeaderAndContent(fileContent)
	if err != nil {
		return nil, err
	}
	blog.SlugTitle = slugTitle
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
	exp := `asset_img\s+(.*(\.jpg|\.png|\.jpeg))\s+`
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

func collectImgsFromBlogDir(blogDirs []string) ([]string, error) {
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
				images = append(images, dir+"/"+f.Name())
			}
		}
	}

	return images, nil
}

func GenerateDocusaurusBlogs(author string, blogs []*HexoBlog) []*DocusaurusBlog {
	var results []*DocusaurusBlog
	for _, hb := range blogs {
		dblog := generateDocusaurusBlog(author, hb)
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
	exp := `{%\s+asset_img\s+(.*(\.jpg|\.png|\.jpeg)).*%}`
	regp, err := regexp.Compile(exp)
	if err != nil {
		return "", err
	}

	results := regp.FindAllString(content, -1)
	fmt.Printf("matched strings: %v\n", results)

	replaceMap := map[string]string{}
	for _, rs := range results {
		imgName := getImgName(rs)
		replaceMap[rs] = fmt.Sprintf("![%s](./%s)", strings.Split(imgName, ".")[0], imgName)
	}

	replaced := regp.ReplaceAllStringFunc(content, func(s string) string {
		return replaceMap[s]
	})
	return replaced, nil
}

func generateDocusaurusBlog(author string, hblog *HexoBlog) *DocusaurusBlog {
	dblog := &DocusaurusBlog{
		SlugTitle: hblog.SlugTitle,
		Date:      strings.Split(hblog.Date, " ")[0], // only need year-month-day
		Content:   hblog.Content,
		Imgs:      hblog.Imgs,
	}
	dblog.Content = fmt.Sprintf(docusaurusTmpl,
		hblog.SlugTitle,
		hblog.Title,
		author,
		strings.Join(hblog.Tags, ", "),
		hblog.Content,
	)

	if len(hblog.Imgs) == 0 {
		return dblog
	}

	var err error
	dblog.Content, err = replaceImg(dblog.Content)
	if err != nil {
		fmt.Printf("[generateDocusaurusBlog] replaceImg err: %v\n", err)
		return dblog
	}

	return dblog
}

func ExportDocusaurusBlogs(docBlogPath string, docBlogs []*DocusaurusBlog, allImgs []string) error {
	errs := []string{}
	for _, blog := range docBlogs {
		if err := generateDocusaurusBlogDirs(docBlogPath, blog, allImgs); err != nil {
			fmt.Printf("generateDocusaurusBlogDirs err: %v\n", err)
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, ", "))
	}

	return nil
}

func generateDocusaurusBlogDirs(docBlogPath string, docBlog *DocusaurusBlog, allImgs []string) error {
	dirName := docBlog.Date + "-" + docBlog.SlugTitle
	dirPath := docBlogPath + "/" + dirName
	if err := os.Mkdir(dirPath, os.ModePerm); err != nil {
		fmt.Printf("mkdir err: %v\n", err)
		return err
	}

	indexMd := dirPath + "/" + "index.md"
	if err := ioutil.WriteFile(indexMd, []byte(docBlog.Content), fs.ModePerm); err != nil {
		return err
	}
	if len(docBlog.Imgs) <= 0 {
		return nil
	}

	blogSrcImgs := filterImgs(allImgs, docBlog.Imgs)
	if len(blogSrcImgs) == 0 {
		fmt.Printf("filterImgs not found src imgs, docBlog.Imgs: %v\n", docBlog.Imgs)
		return nil
	}

	copyImgs(blogSrcImgs, dirPath)
	return nil
}

func filterImgs(allSrcImgs []string, imgNames []string) []string {
	imgSet := make(map[string]struct{}, len(imgNames))
	for _, img := range imgNames {
		imgSet[img] = struct{}{}
	}

	filteredImgs := []string{}
	for _, imgPath := range allSrcImgs {
		_, imgFile := filepath.Split(imgPath)
		_, ok := imgSet[imgFile]
		if ok {
			filteredImgs = append(filteredImgs, imgPath)
		}
	}

	return filteredImgs
}

func copyImgs(blogSrcImgs []string, dst string) {
	var wg sync.WaitGroup
	for _, srcImg := range blogSrcImgs {
		wg.Add(1)

		imgPath := srcImg
		go func() {
			defer wg.Done()
			_, imgFile := filepath.Split(imgPath)
			if err := copyImg(imgPath, dst+"/"+imgFile); err != nil {
				fmt.Printf("copyImg err: %b\n", err)
			}
		}()
	}
	wg.Wait()
}

func copyImg(src, dst string) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		fmt.Printf("read file err: %v, src: %s\n", err, src)
		return err
	}

	// Write data to dst
	if err := ioutil.WriteFile(dst, data, 0644); err != nil {
		fmt.Printf("write file err: %v, dst: %s\n", err, dst)
		return err
	}

	return nil
}
