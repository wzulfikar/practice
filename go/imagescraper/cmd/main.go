package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/wzulfikar/lab/go/imagescraper"
)

// scraper file:///Users/strawhat/Desktop/Faculties.webarchive ".card__thumbnail img"
// output: {image title|alt}_{filename}.jpg
func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: scraper [url] [selector] [dir]")
		return
	}

	url := os.Args[1]
	selector := os.Args[2]
	dir := os.Args[3]

	const scraperConcurrency = 5
	countImages := 0
	defer func() func() {
		start := time.Now()
		return func() {
			fmt.Printf("[DONE] Images scraped: %d\n", countImages)
			fmt.Println("Time elapsed:", time.Since(start))
		}
	}()()

	if newUrl, from, to, pageOk := getUrlPage(url); pageOk {
		concurrency := 5
		sem := make(chan bool, concurrency)

		for i := from; i <= to; i++ {
			sem <- true
			go func(i int) {
				defer func() { <-sem }()

				targetUrl := newUrl + strconv.Itoa(i)
				fmt.Printf("[START] %s\n", targetUrl)

				images := imagescraper.Scrape(targetUrl, selector, dir, scraperConcurrency)

				countImages += len(images)

				fmt.Println("[DONE]", targetUrl)
			}(i)
		}
		for i := 0; i < cap(sem); i++ {
			sem <- true
		}
		return
	}

	fmt.Println("Scraping images from", url)
	images := imagescraper.Scrape(url, selector, dir, scraperConcurrency)
	countImages = len(images)
}

func getUrlPage(url string) (string, int, int, bool) {
	param := strBetween(url, "[", "]")
	if param == "" || !strings.Contains(param, "-") {
		return "", 0, 0, false
	}

	page := strings.Split(param, "-")
	if len(page) != 2 {
		return "", 0, 0, false
	}
	page1, err := strconv.Atoi(page[0])
	if err != nil {
		return "", 0, 0, false
	}
	page2, err := strconv.Atoi(page[1])
	if err != nil {
		return "", 0, 0, false
	}

	if page1 >= page2 {
		log.Fatal("ERROR: URL contains invalid page. `page1` must be smaller than `page2`.")
	}

	url = strings.Replace(url, "["+param+"]", "", -1)
	return url, page1, page2, true
}

func strBetween(str string, a string, b string) string {
	// Get substring between two strings.
	posFirst := strings.Index(str, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(str, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return str[posFirstAdjusted:posLast]
}
