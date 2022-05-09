package binance

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const API_URL = "https://api.binance.com/api/v3/avgPrice"

type Price struct {
	Price float64 `json:",string"`
}

func GetResponse() string {
	vals := []string{"BTC", "ETH"}

	var response string
	for _, v := range vals {
		curs, _ := GetCurrs(v)
		response += fmt.Sprintf("\n%s %s$", v, curs)
	}
	return response
}

func GetCurrs(val string) (string, error) {
	if val == "" {
		return "", errors.New("empty parameter")
	}

	resp, err := http.Get(API_URL + "?symbol=" + val + "USDT")

	if err != nil {
		return "", fmt.Errorf("get currency error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid status code: %w", err)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("invalid response body: %w", err)
	}

	var price Price
	err = json.Unmarshal(bytes, &price)
	if err != nil {
		return "", fmt.Errorf("invalid unmarshal error: %w", err)
	}

	return fmt.Sprintf("%.2f", price.Price), nil
}
