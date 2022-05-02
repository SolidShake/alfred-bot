package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
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

const API_URL = "https://www.cbr.ru/scripts/XML_daily.asp"

func main() {
	result, _ := getCurrency("")
	response, _ := createResponse(*result)
	fmt.Println(string(response))

	prevDay, _ := getPreviousDate(*result)
	resultPrev, _ := getCurrency("?date_req=" + prevDay)
	responsePrev, _ := createResponse(*resultPrev)
	fmt.Println(string(responsePrev))
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

func createResponse(result ValCurs) ([]byte, error) {
	var response Response
	for _, valute := range result.Valutes {
		if isCurrsUsed(valute.CharCode) {
			response.Vals = append(
				response.Vals,
				Vals{valute.CharCode, valute.Name, valute.Value},
			)
		}
	}

	bytes, err := json.Marshal(&response)
	if err != nil {
		return nil, fmt.Errorf("marshal error %w", err)
	}

	return bytes, nil
}

func getPreviousDate(result ValCurs) (string, error) {
	dataDate, err := time.Parse("02.01.2006", result.Date)
	if err != nil {
		return "", fmt.Errorf("date parse error: %w", err)
	}

	return dataDate.Add(-24 * time.Hour).Format("02/01/2006"), nil
}
