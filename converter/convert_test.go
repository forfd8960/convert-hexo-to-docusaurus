package converter

import (
	"reflect"
	"testing"
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
