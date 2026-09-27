package websearchtool

type SearchInput struct {
	Query string `json:"query" jsonschema:"user query to search"`
	Top   int    `json:"top" jsonschema:"no of top results to return"`
}

type SearchResult struct {
	Url   string  `json:"url" jsonschema:"url of the search result"`
	Score float64 `json:"score" jsonschema:"score given by search engine to result"`
}

type SearchOutput struct {
	Results []SearchResult `json:"results" jsonschema:"top or less results retutned by search engine"`
}
