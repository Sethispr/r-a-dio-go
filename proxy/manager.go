// Package proxy manages a pool of validated HTTP proxies to use later on.
package proxy

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	_proxyURL      = "https://raw.githubusercontent.com/monosans/proxy-list/main/proxies/http.txt"
	_checkTimeout  = 3 * time.Second
	_clientTimeout = 15 * time.Second
	_workerCount   = 50
	_maxProxyCheck = 300
)

var (
	_proxies []string
	_mu      sync.RWMutex // Protects slice during refresh.
)

// GetRandomClient returns HTTP client routed through verified proxy with its own cookie dookie pookie jar.
func GetRandomClient() *http.Client {
	_mu.RLock()
	defer _mu.RUnlock()

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Timeout: _clientTimeout,
		Jar:     jar,
	}

	if len(_proxies) > 0 {
		addr := _proxies[rand.Intn(len(_proxies))]
		proxyURL, _ := url.Parse("http://" + addr)
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	return client
}

// GetCount returns number of verified working proxies.
func GetCount() int {
	_mu.RLock()
	defer _mu.RUnlock()
	return len(_proxies)
}

// Refresh pull latest list and checks the proxy.
func Refresh() {
	fmt.Print("refreshing proxies... ")

	raw, err := fetchProxyList()
	if err != nil {
		fmt.Printf("failed: %v\n", err)
		return
	}

	live := validateProxies(raw)

	_mu.Lock()
	_proxies = live
	_mu.Unlock()

	fmt.Printf("done (%d live)\n", len(live))
}

// fetchProxyList pulls raw txt list of IP port addresses.
func fetchProxyList() ([]string, error) {
	resp, err := http.Get(_proxyURL)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var proxies []string
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			proxies = append(proxies, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return proxies, nil
}

// validateProxies filters out dead proxies.
func validateProxies(raw []string) []string {
	limit := len(raw)
	if limit > _maxProxyCheck {
		limit = _maxProxyCheck
	}

	jobs := make(chan string, limit)
	results := make(chan string, limit)
	var wg sync.WaitGroup

	for i := 0; i < _workerCount; i++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	for i := 0; i < limit; i++ {
		jobs <- raw[i]
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	var validated []string
	for addr := range results {
		validated = append(validated, addr)
	}

	return validated
}

func worker(jobs <-chan string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for addr := range jobs {
		if checkProxy(addr) {
			results <- addr
		}
	}
}

// checkProxy tries to reach Google endpoint through proxy.
func checkProxy(addr string) bool {
	proxyURL, err := url.Parse("http://" + addr)
	if err != nil {
		return false
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
		Timeout: _checkTimeout,
	}

	resp, err := client.Get("http://www.google.com/generate_204")
	if err != nil {
		return false
	}
	defer func() { _ = resp.Body.Close() }()

	return resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusOK
}
