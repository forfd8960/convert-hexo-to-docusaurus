package converter

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
	Title   string   // slug title
	Date    string   // 2022-06-28
	Tags    []string //
	Content string   // blog content
}

func ReadHexoBlogs() []*HexoBlog {
	return nil
}

func extractContent(blog string) string {
	return ""
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
