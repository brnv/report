package main

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
)

func requestAppstore(
	vendor string, date string, token string,
) ([]byte, error) {
	request, err := http.NewRequest(
		"GET",
		"https://api.appstoreconnect.apple.com/v1/salesReports",
		nil,
	)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	request.Header.Set("Accept", "application/a-gzip")

	query := request.URL.Query()
	query.Add("filter[vendorNumber]", vendor)
	query.Add("filter[reportDate]", date)
	query.Add("filter[frequency]", "DAILY")
	query.Add("filter[reportSubType]", "SUMMARY")
	query.Add("filter[reportType]", "SALES")

	request.URL.RawQuery = query.Encode()

	client := http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	reader := response.Body

	if response.Header.Get("Content-Encoding") == "agzip" {
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			return nil, err
		}
	}

	output, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return output, nil
}
