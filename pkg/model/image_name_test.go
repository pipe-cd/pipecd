package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImageNameString(t *testing.T) {
	testcases := []struct {
		name   string
		domain string
		repo   string
		want   string
	}{
		{
			name: "empty repo",
			want: "",
		},
		{
			name: "domain omitted",
			repo: "repo",
			want: "repo",
		},
		{
			name:   "with domain",
			domain: "domain",
			repo:   "repo",
			want:   "domain/repo",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			i := ImageName{
				Repo:   tc.repo,
				Domain: tc.domain,
			}
			got := i.String()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestImageRefString(t *testing.T) {
	testcases := []struct {
		name      string
		imageName ImageName
		tag       string
		want      string
	}{
		{
			name: "empty repo",
			want: "",
		},
		{
			name: "tag omitted",
			imageName: ImageName{
				Domain: "domain",
				Repo:   "repo",
			},
			want: "domain/repo",
		},
		{
			name: "with tag",
			imageName: ImageName{
				Domain: "domain",
				Repo:   "repo",
			},
			tag:  "tag",
			want: "domain/repo:tag",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			i := ImageRef{
				ImageName: tc.imageName,
				Tag:       tc.tag,
			}
			got := i.String()
			assert.Equal(t, tc.want, got)
		})
	}
}
