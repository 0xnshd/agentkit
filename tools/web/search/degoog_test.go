package websearchtool_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	websearchtool "github.com/0xnshd/agentkit/tools/web/search"
	"github.com/0xnshd/testingx"
)

func Test_Degoog_Success(t *testing.T) {
	fakeserver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"results": [{"url": "url", "score": 1}, {"url": "url", "score": 1}, {"url": "url", "score": 1}]}`))
	}))

	temp := websearchtool.DEGOOG_URL
	websearchtool.DEGOOG_URL = fakeserver.URL

	defer func() {
		websearchtool.DEGOOG_URL = temp
		fakeserver.Close()
	}()

	tests := []struct {
		name  string
		input websearchtool.SearchInput
		want  int
	}{
		{
			name: "top < len(results)",
			input: websearchtool.SearchInput{
				Query: "test",
				Top:   2,
			},
			want: 2,
		},
		{
			name: "top > len(results)",
			input: websearchtool.SearchInput{
				Query: "test",
				Top:   4,
			},
			want: 3,
		},
		{
			name: "top == len(results)",
			input: websearchtool.SearchInput{
				Query: "test",
				Top:   3,
			},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := websearchtool.Degoog(tt.input)
			testingx.CheckWErr(t, len(out.Results), tt.want, err, nil)
		})
	}
}

func Test_Degoog_int(t *testing.T) {
	out, err := websearchtool.Degoog(websearchtool.SearchInput{Query: "k8s", Top: 2})
	fmt.Println(out, err)
}
