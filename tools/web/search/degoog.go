package websearchtool

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

var DEGOOG_URL = "http://localhost:4444/api/search"

func Degoog(in SearchInput) (SearchOutput, error) {
	res, err := querySearchApi(in.Query)
	if err != nil {
		return SearchOutput{}, err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	var out SearchOutput

	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return SearchOutput{}, err
	}

	if in.Top <= len(out.Results) {
		out.Results = out.Results[:in.Top]
	}

	return out, nil
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
