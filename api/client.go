// Package api handles all network requests to r-a-d.io.
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"radiogo/models"
	"radiogo/proxy"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	_apiBase   = "https://r-a-d.io/api"
	_siteRoot  = "https://r-a-d.io"
	_userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
)

var (
	// Regex to extract the mrbeast gorilla security token required for POST requests from HTML.
	_csrfPattern = regexp.MustCompile(`name="gorilla.csrf.Token" value="([^"]+)"`)
	_httpClient  = &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}
)

// FetchStatus pulls current stream, checks DJ name and listener counts.
// If dj is Hanyuu-sama then you can request since its a clanker, but when a real dj comes in they usually disable requests.
func FetchStatus() (*models.RadioStatus, error) {
	req, err := http.NewRequest(http.MethodGet, _apiBase, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", _userAgent)

	resp, err := _httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	var status models.RadioStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, err
	}

	return &status, nil
}

// Search queries db for songs matching input.
func Search(query string) (*models.SearchResponse, error) {
	searchURL := fmt.Sprintf("%s/search/%s", _apiBase, url.PathEscape(query))

	req, err := http.NewRequest(http.MethodGet, searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", _userAgent)

	resp, err := _httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	var results models.SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}

	return &results, nil
}

// SubmitRequest basically gets a fresh proxy > scrape gorilla csrf token, then try to req song.
func SubmitRequest(song models.Song, query string) bool {
	client := proxy.GetRandomClient()

	token, err := getCSRFToken(client, query)
	if err != nil {
		return false
	}

	return submitWithToken(client, song, query, token)
}

// getCSRFToken scrapes r-a-d.io search page to find the token needed.
func getCSRFToken(client *http.Client, query string) (string, error) {
	searchURL := fmt.Sprintf("%s/search?q=%s", _siteRoot, url.QueryEscape(query))

	req, err := http.NewRequest(http.MethodGet, searchURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", _userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	matches := _csrfPattern.FindSubmatch(body)
	if len(matches) < 2 {
		return "", fmt.Errorf("ERR: no gorilla srf token found lols")
	}

	return string(matches[1]), nil
}

// submitWithToken makes the final POST call with correct headers to make non sus doakes browser req.
func submitWithToken(client *http.Client, song models.Song, query, token string) bool {
	postURL := fmt.Sprintf("%s/v1/request?trackid=%d&q=%s&page=1",
		_siteRoot, song.ID, url.QueryEscape(query))

	form := url.Values{}
	form.Set("gorilla.csrf.Token", token)
	form.Set("id", strconv.Itoa(song.ID))

	req, err := http.NewRequest(http.MethodPost, postURL, strings.NewReader(form.Encode()))
	if err != nil {
		return false
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", _userAgent)
	req.Header.Set("Referer", fmt.Sprintf("%s/search?q=%s", _siteRoot, url.QueryEscape(query)))
	req.Header.Set("HX-Request", "true")

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer func() { _ = resp.Body.Close() }()

	return resp.StatusCode == http.StatusOK
}
