package bank

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
)

type Valute struct {
	XMLName  xml.Name `xml:"Valute"`
	CharCode string   `xml:"CharCode"`
	Name     string   `xml:"Name"`
	Value    string   `xml:"Value"`
}

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Date    string   `xml:"Date,attr"`
	Valutes []Valute `xml:"Valute"`
}

type Vals struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Response struct {
	Vals []Vals `json:"vals"`
}

type ExchangeRate struct {
	Code  string
	Value string
	Diff  string
}

func (r Response) ToString() string {
	var result string
	for _, v := range r.Vals {
		result += fmt.Sprintf("%s %s\n", v.Code, v.Value)
	}

	return result
}

const API_URL = "https://www.cbr.ru/scripts/XML_daily.asp"

func CreateBankResponse() string {
	var result string
	for _, exchange := range GetCurrency() {
		result += fmt.Sprintf("\n%s %s₽ (%s₽)", exchange.Code, exchange.Value, exchange.Diff)
	}

	return result
}

func GetCurrency() []ExchangeRate {
	todayCurrency, _ := getCurrency("")
	prevDate, _ := getPreviousDate(*todayCurrency)
	yesterdayCurrency, _ := getCurrency("?date_req=" + prevDate)

	return createResponse(*todayCurrency, *yesterdayCurrency) //getCurrency("?date_req=" + prevDay)
}

func getCurrency(queryParams string) (*ValCurs, error) {
	resp, err := http.Get(API_URL + queryParams) //https://www.cbr.ru/scripts/XML_daily.asp?date_req=29/04/2022
	if err != nil {
		return nil, fmt.Errorf("get currency error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code: %w", err)
	}

	var result ValCurs
	decoder := xml.NewDecoder(resp.Body)
	decoder.CharsetReader = makeCharsetReader
	err = decoder.Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("xml unmarshal error: %w", err)
	}

	return &result, nil
}

// //xml unmarshal error: xml: encoding "windows-1251" declared but Decoder.CharsetReader is nil
func makeCharsetReader(charset string, input io.Reader) (io.Reader, error) {
	if charset == "windows-1251" {
		return charmap.Windows1251.NewDecoder().Reader(input), nil
	}

	return nil, fmt.Errorf("unknown charset: %s", charset)
}

func isCurrsUsed(code string) bool {
	// to const
	currsList := []string{"USD", "EUR"}
	for _, v := range currsList {
		if code == v {
			return true
		}
	}

	return false
}

func createResponse(result, oldResult ValCurs) []ExchangeRate {
	var exchangeRate []ExchangeRate

	// @TODO refactor it
	for _, yesterdayCurrency := range oldResult.Valutes {
		for _, todayCurrency := range result.Valutes {
			if isCurrsUsed(todayCurrency.CharCode) && todayCurrency.CharCode == yesterdayCurrency.CharCode {
				exchangeRate = append(exchangeRate, ExchangeRate{
					Code:  todayCurrency.CharCode,
					Value: normalizeValue(todayCurrency.Value),
					Diff:  calculateDiff(todayCurrency.Value, yesterdayCurrency.Value),
				})
			}
		}
	}

	return exchangeRate
}

func normalizeValue(value string) string {
	value = strings.ReplaceAll(value, ",", ".")
	floatVal, _ := strconv.ParseFloat(value, 64)

	return fmt.Sprintf("%.2f", floatVal)
}

func calculateDiff(value, oldValue string) string {
	value = normalizeValue(value)
	oldValue = normalizeValue(oldValue)

	floatVal, _ := strconv.ParseFloat(value, 64)
	oldFloatVal, _ := strconv.ParseFloat(oldValue, 64)

	return fmt.Sprintf("%.2f", floatVal-oldFloatVal)

}

func getPreviousDate(result ValCurs) (string, error) {
	dataDate, err := time.Parse("02.01.2006", result.Date)
	if err != nil {
		return "", fmt.Errorf("date parse error: %w", err)
	}

	return dataDate.Add(-24 * time.Hour).Format("02/01/2006"), nil
}
