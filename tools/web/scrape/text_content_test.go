package webscrapetool_test

import (
	"testing"

	webscrapetool "github.com/0xnshd/agentkit/tools/web/scrape"
	"github.com/0xnshd/testingx"
)

func Test_ScrapeTextContent(t *testing.T) {
	tests := []struct {
		name    string
		urls    []string
		wantErr bool
	}{
		{
			name:    "invalid url",
			urls:    []string{"abc"},
			wantErr: true,
		},
		{
			name:    "valid url",
			urls:    []string{"https://en.wikipedia.org/wiki/Go_(programming_language)"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := webscrapetool.ScrapeTextContent(tt.urls)
			out := <-results
			testingx.Check(t, out.Err != nil, tt.wantErr)
			if !tt.wantErr {
				testingx.Check(t, len(out.Content) > 0, true)
			}
		})
	}
}
