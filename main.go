// Package main is the entry point for the radiogo CLI.
package main

import (
	"bufio"
	"fmt"
	"os"
	"radiogo/api"
	"radiogo/models"
	"radiogo/proxy"
	"strconv"
	"strings"
	"time"
)

var cart []models.Song

func main() {
	proxy.Refresh()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		showStatus()
		fmt.Printf("\nProxies: %d | Cart: %d\n", proxy.GetCount(), len(cart))
		fmt.Print("[1] search [2] cart [3] refresh [4] exit\n> ")

		if !scanner.Scan() {
			break
		}

		switch strings.TrimSpace(scanner.Text()) {
		case "1":
			handleSearch(scanner)
		case "2":
			handleCart(scanner)
		case "3":
			proxy.Refresh()
		case "4":
			return
		}
	}
}

// showStatus fetches then displays the current r-a-d.io station state.
func showStatus() {
	status, err := api.FetchStatus()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		time.Sleep(time.Second)
		return
	}

	fmt.Printf("now playing: %s\n", status.Main.NowPlaying)
	fmt.Printf("dj: %s | listeners: %d\n", status.Main.DJ.Name, status.Main.Listeners)

	renderQueue(status.Main.Queue)
}

func renderQueue(queue []models.Track) {
	if len(queue) == 0 {
		return
	}
	fmt.Println("\nqueue:")
	for i, track := range queue {
		if i >= 5 {
			break
		}
		marker := "auto"
		if track.Type == 1 {
			marker = "req"
		}
		fmt.Printf("  %d. [%s] %s\n", i+1, marker, track.MetaData)
	}
}

func handleSearch(scanner *bufio.Scanner) {
	fmt.Print("search: ")
	if !scanner.Scan() {
		return
	}

	query := strings.TrimSpace(scanner.Text())
	if query == "" {
		return
	}

	results, err := api.Search(query)
	if err != nil || results.Total == 0 {
		handleSearchError(err, results)
		return
	}

	displaySearchResults(results)

	fmt.Print("\n[number] request | [a number] add to cart | [b] back\n> ")
	if !scanner.Scan() {
		return
	}

	processSearchSelection(scanner.Text(), results, query)
}

func handleSearchError(err error, results *models.SearchResponse) {
	if err != nil {
		fmt.Printf("search failed: %v\n", err)
	} else {
		fmt.Println("no results")
	}
	time.Sleep(time.Second)
}

func displaySearchResults(results *models.SearchResponse) {
	limit := results.Total
	if limit > 10 {
		limit = 10
	}

	fmt.Printf("\nfound %d songs:\n", results.Total)
	for i := 0; i < limit; i++ {
		song := results.Data[i]
		status := "ok"
		if !song.Requestable {
			status = "cooldown"
		}
		fmt.Printf("[%d] %s - %s [%s]\n", i+1, song.Artist, song.Title, status)
	}
}

func processSearchSelection(input string, results *models.SearchResponse, query string) {
	input = strings.TrimSpace(input)
	if input == "b" || input == "" {
		return
	}

	parts := strings.Fields(input)
	if len(parts) == 2 && parts[0] == "a" {
		addToCart(parts[1], results.Data)
		return
	}

	performImmediateRequest(input, results.Data, query)
}

func addToCart(indexStr string, songs []models.Song) {
	idx, err := strconv.Atoi(indexStr)
	if err != nil || idx < 1 || idx > len(songs) {
		fmt.Println("invalid number")
	} else {
		cart = append(cart, songs[idx-1])
		fmt.Println("added to cart")
	}
	time.Sleep(time.Second)
}

func performImmediateRequest(indexStr string, songs []models.Song, query string) {
	idx, err := strconv.Atoi(indexStr)
	if err != nil || idx < 1 || idx > len(songs) {
		fmt.Println("invalid number")
		time.Sleep(time.Second)
		return
	}

	song := songs[idx-1]
	fmt.Printf("requesting %s - %s...\n", song.Artist, song.Title)

	if api.SubmitRequest(song, query) {
		fmt.Println("✓ sent")
	} else {
		fmt.Println("✗ failed")
	}
	time.Sleep(time.Second)
}

func handleCart(scanner *bufio.Scanner) {
	if len(cart) == 0 {
		fmt.Println("cart empty")
		time.Sleep(time.Second)
		return
	}

	fmt.Printf("\ncart (%d songs):\n", len(cart))
	for i, song := range cart {
		fmt.Printf("%d. %s - %s\n", i+1, song.Artist, song.Title)
	}

	fmt.Print("\n[s] send all | [c] clear | [b] back\n> ")
	if !scanner.Scan() {
		return
	}

	switch strings.TrimSpace(scanner.Text()) {
	case "s":
		sendAll()
	case "c":
		cart = nil
		fmt.Println("cleared")
		time.Sleep(time.Second)
	}
}

// sendAll reads the cart then submits requests and delays to prevent IP flagging.
func sendAll() {
	total := len(cart)
	sent := 0

	fmt.Printf("\nsending %d requests...\n", total)
	for i, song := range cart {
		fmt.Printf("[%d/%d] %s - %s ", i+1, total, song.Artist, song.Title)

		if api.SubmitRequest(song, song.Title) {
			fmt.Println("[DONE]")
			sent++
		} else {
			fmt.Println("[ERR]")
		}
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Printf("\ncompleted: %d/%d\n", sent, total)
	cart = nil
	time.Sleep(2 * time.Second)
}
