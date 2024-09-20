package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// SharedData represents the data structure where scraped data will be stored.
type SharedData struct {
	mu   sync.Mutex
	data map[string]float64
}

// ScrapeWebsite simulates scraping financial data from a website.
func ScrapeWebsite(url string, shared *SharedData, wg *sync.WaitGroup, r *rand.Rand) {
	// Decrement the WaitGroup counter when the goroutine completes.
	defer wg.Done()

	// Simulate data scraping with random delay
	time.Sleep(time.Duration(r.Intn(1000)) * time.Millisecond)

	// Simulate the data scraped (random stock price)
	scrapedData := r.Float64() * 1000

	// Use Mutex to safely write to the shared data structure
	shared.mu.Lock()
	shared.data[url] = scrapedData
	shared.mu.Unlock()

	fmt.Printf("Scraped data from %s: %f\n", url, scrapedData)
}

func main() {
	// Seed random number generator for simulation
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	// Shared data structure for storing scraped financial data
	sharedData := &SharedData{
		data: make(map[string]float64),
	}

	// List of websites to scrape financial data from
	websites := []string{
		"https://finance.yahoo.com/",
		"https://www.investing.com/",
		"https://www.alphavantage.co/",
		"https://www.google.com/finance/",
		"https://www.nasdaq.com/",
		"https://www.bloomberg.com/",
		"https://www.morningstar.com/",
		"https://coinmarketcap.com/",
		"https://data.worldbank.org/",
		"https://www.quandl.com/",
	}

	// Create a WaitGroup to manage goroutines
	var wg sync.WaitGroup

	// Start scraping each website concurrently
	for _, url := range websites {
		wg.Add(1) // Increment the WaitGroup counter
		go ScrapeWebsite(url, sharedData, &wg, r)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Display the collected financial data
	fmt.Println("Collected Financial Data:")
	sharedData.mu.Lock()
	for site, value := range sharedData.data {
		fmt.Printf("%s: %f\n", site, value)
	}
	sharedData.mu.Unlock()
}
