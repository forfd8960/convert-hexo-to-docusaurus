package converter

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseHeader(t *testing.T) {
	type args struct {
		header string
	}
	tests := []struct {
		name string
		args args
		want *HexoBlog
	}{
		{
			name: "happy path",
			args: args{
				header: `title: host-static-blog-on-aws-s3
				date: 2022-06-28 09:16:48
				tags: AWS, S3, Hexo`,
			},
			want: &HexoBlog{
				Title: "host-static-blog-on-aws-s3",
				Date:  "2022-06-28 09:16:48",
				Tags:  []string{"AWS", "S3", "Hexo"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseHeader(tt.args.header); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadHexoBlogs(t *testing.T) {
	type args struct {
		dir string
	}

	testContent := `
## How to Host Static Site on AWS S3

...

### Step1

{% asset_img create_bucket_step1.jpg create bucket on s3 %}` + "\n\n" + "```python" + `
print("host static site s3")
` + "```\n"

	testPath := os.Getenv("HOME") + "/Documents/projects/convert-hexo-to-docusaurus/mockhexoblogs"

	tests := []struct {
		name    string
		args    args
		blogs   []*HexoBlog
		imgs    []string
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				dir: testPath,
			},
			blogs: []*HexoBlog{
				{
					SlugTitle: "host-static-site-on-aws-s3.md",
					Title:     "host static blog on aws s3",
					Date:      "2022-06-28 09:16:48",
					Tags:      []string{"AWS", "S3", "Hexo"},
					Content:   testContent,
					Imgs:      []string{"create_bucket_step1.jpg"},
				},
			},
			imgs: []string{
				testPath + "/create_bucket_step1.jpg",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ReadHexoBlogs(tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadHexoBlogs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.blogs, got)
			assert.Equal(t, tt.imgs, got1)
		})
	}
}

func TestRegexp(t *testing.T) {
	content := `abc
	{% asset_img create_bucket_step1.jpg create bucket on s3 %}
	{% asset_img create_bucket_step2.jpg create bucket on s3 %}
	def
	`

	exp := `{%\s*asset_img\s*(.+\.jpg).*%}`
	regp, err := regexp.Compile(exp)
	assert.Nil(t, err)

	results := regp.FindAllString(content, -1)
	expectResults := []string{
		`{% asset_img create_bucket_step1.jpg create bucket on s3 %}`,
		`{% asset_img create_bucket_step2.jpg create bucket on s3 %}`,
	}
	assert.Equal(t, expectResults, results)
	fmt.Printf("matched strings: %v\n", results)
}

func TestRegexp1(t *testing.T) {
	content := `abc
	{% asset_img create_bucket_step1.jpg create bucket on s3 %}
	{% asset_img create_bucket_step2.jpg create bucket on s3 %}
	{% asset_img test_img.png create bucket on s3 %}
	{% asset_img test_img1.jpeg create bucket on s3 %}
	def
	`

	exp := `asset_img\s+(.*(\.jpg|\.png|\.jpeg))\s+`
	regp, err := regexp.Compile(exp)
	assert.Nil(t, err)

	results := regp.FindAllString(content, -1)
	expectResults := []string{
		`asset_img create_bucket_step1.jpg `,
		`asset_img create_bucket_step2.jpg `,
		`asset_img test_img.png `,
		`asset_img test_img1.jpeg `,
	}
	assert.Equal(t, expectResults, results)
	fmt.Printf("matched strings: %v\n", results)
}

func TestRegexp2(t *testing.T) {
	content := `abc
	{% asset_img create_bucket_step1.jpg create bucket on s3 %}
	{% asset_img create_bucket_step2.jpg create bucket on s3 %}
	{% asset_img test_img.png test img %}
	{% asset_img test_img1.jpeg test img1 %}
	def
	`

	exp := `{%\s+asset_img\s+(.*(\.jpg|\.png|\.jpeg)).*%}`
	regp, err := regexp.Compile(exp)
	assert.Nil(t, err)

	results := regp.FindAllString(content, -1)
	expectResults := []string{
		`{% asset_img create_bucket_step1.jpg create bucket on s3 %}`,
		`{% asset_img create_bucket_step2.jpg create bucket on s3 %}`,
		`{% asset_img test_img.png test img %}`,
		`{% asset_img test_img1.jpeg test img1 %}`,
	}
	assert.Equal(t, expectResults, results)
	fmt.Printf("matched strings: %v\n", results)

	replaceMap := map[string]string{
		results[0]: "![0](./create_bucket_step1.jpg)",
		results[1]: "![1](./create_bucket_step2.jpg)",
		results[2]: "![2](./test_img.png)",
		results[3]: "![3](./test_img1.jpeg)",
	}

	replaced := regp.ReplaceAllStringFunc(content, func(s string) string {
		return replaceMap[s]
	})
	expectRepalced := `abc
	![0](./create_bucket_step1.jpg)
	![1](./create_bucket_step2.jpg)
	![2](./test_img.png)
	![3](./test_img1.jpeg)
	def
	`
	fmt.Printf("replaced: %s\n", replaced)
	assert.Equal(t, expectRepalced, replaced)
}

func Test_replaceImg(t *testing.T) {
	type args struct {
		content string
	}

	testContent := `abc
	{% asset_img create_bucket_step1.jpg create bucket on s3 %}
	{% asset_img create_bucket_step2.jpg create bucket on s3 %}
	{% asset_img test_img.png test img %}
	{% asset_img test_img1.jpeg test img1 %}
	def
	`
	expectRepalced := `abc
	![create_bucket_step1](./create_bucket_step1.jpg)
	![create_bucket_step2](./create_bucket_step2.jpg)
	![test_img](./test_img.png)
	![test_img1](./test_img1.jpeg)
	def
	`

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				content: testContent,
			},
			want:    expectRepalced,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := replaceImg(tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("replaceImg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("replaceImg() = %v, want %v", got, tt.want)
			}
		})
	}
}
