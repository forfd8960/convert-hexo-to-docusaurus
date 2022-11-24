package converter

import (
	"os"
	"reflect"
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

	testContent := `\n` + `## How to Host Static Site on AWS S3

...

### Step1

{% asset_img create_bucket_step1.jpg create bucket on s3 %}` + `\n` + "```python\n" + `
print("host static site s3")` + "```"

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
