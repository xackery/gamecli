package gog

import (
	"context"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
)

// Update will request from gog an update of local games supported
func (s *Gog) Update(ctx context.Context) error {
	return nil
	hiddenFlag := 1
	isUpdated := 0
	mediaType := 1
	page := 0
	url := fmt.Sprintf("https://www.gog.com/account/getFilteredProducts?hiddenFlag=%d&isUpdated=%d&mediaType=%d&sortBy=title&system=&page=%d", hiddenFlag, isUpdated, mediaType, page)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("newRequest: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("do: %w", err)
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("readall: %w", err)
	}
	fmt.Println("status:", resp.Status)
	err = ioutil.WriteFile("out.txt", data, fs.ModePerm)
	if err != nil {
		return fmt.Errorf("writefile: %w", err)
	}
	fmt.Println(string(data))

	return nil
}
