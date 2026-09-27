package websearchtool

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var DEGOOG_URL = "http://localhost:4444/api/search"

func Degoog(ctx context.Context, req *mcp.CallToolRequest, in SearchInput) (*mcp.CallToolResult, SearchOutput, error) {
	res, err := querySearchApi(in.Query)
	if err != nil {
		return nil, SearchOutput{}, err
	}
	defer res.Body.Close()

	var out SearchOutput

	err = json.NewDecoder(res.Body).Decode(&out)

	if in.Top <= len(out.Results) {
		out.Results = out.Results[:in.Top]
	}

	return nil, out, nil
}

// query degoog search api
// for agent use only first page matters as this has to be quick
func querySearchApi(query string) (*http.Response, error) {
	baseUrl, _ := url.Parse(DEGOOG_URL)
	// params
	params := url.Values{}
	params.Add("q", query)
	params.Add("page", "1")
	params.Add("lang", "en")
	baseUrl.RawQuery = params.Encode()

	req, _ := http.NewRequest(http.MethodGet, baseUrl.String(), nil)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 15 * time.Second}

	return client.Do(req)
}
