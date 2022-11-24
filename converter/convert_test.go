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
	def
	`

	exp := `asset_img\s+(.*\.jpg)\s+`
	regp, err := regexp.Compile(exp)
	assert.Nil(t, err)

	results := regp.FindAllString(content, -1)
	expectResults := []string{
		`asset_img create_bucket_step1.jpg `,
		`asset_img create_bucket_step2.jpg `,
	}
	assert.Equal(t, expectResults, results)
	fmt.Printf("matched strings: %v\n", results)
}
