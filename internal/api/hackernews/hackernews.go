package hackernews

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const API_URL = "https://hn.algolia.com/api/v1/search_by_date?query=(golang,go)&tags=story&hitsPerPage=5"

type News struct {
	Hits []struct {
		Title string `json:"title"`
		URL   string `json:"url"`
	} `json:"hits"`
}

func GetResponse() string {
	news, err := GetNews()
	if err != nil {
		fmt.Println(err)
	}
	var response string

	for _, v := range news.Hits {
		response += fmt.Sprintf("\n - [%s](%s)", v.Title, v.URL)
	}

	return response
}

func GetNews() (*News, error) {
	resp, err := http.Get(API_URL)

	if err != nil {
		return nil, fmt.Errorf("get news error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code: %w", err)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("invalid response body: %w", err)
	}

	fmt.Println(string(bytes))

	var news News
	err = json.Unmarshal(bytes, &news)
	if err != nil {
		return nil, fmt.Errorf("invalid unmarshal error: %w", err)
	}

	return &news, nil
}
