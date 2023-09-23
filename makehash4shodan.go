package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/twmb/murmur3"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Replace this with the website URL you want to fetch the favicon from

	var url *string
	var faviconURL string
	var err error
	url = flag.String("u", "https://example.com", "URL")
	flag.Parse()
	faviconURL = *url
	if *url == "https://example.com" {
		fmt.Println("[!] Error!")
		fmt.Printf("[-] Use: makehash4shodan -u http://example.com/favicon.ico\n")
		fmt.Println("[i] Get all hosts with the same favicon!")
		os.Exit(1)

	}
	if !strings.HasSuffix(*url, "/favicon.ico") {
		faviconURL, err = getFaviconURL(*url)
		if err != nil {
			fmt.Println("Error fetching favicon URL:", err)
			return
		}
	}
	// Fetch the website's favicon

	// Fetch the favicon file
	//	fmt.Println(faviconURL)
	if !strings.HasPrefix(faviconURL, "https://") {
		if !strings.HasPrefix(faviconURL, "http://") {
			faviconURL = fmt.Sprintf(*url+"%s", faviconURL)

		}

	}
	//	fmt.Println(faviconURL)
	faviconBytes := fetchFavicon(faviconURL)
	if err != nil {
		fmt.Println("Error fetching favicon:", err)
		return
	}

	// Calculate the Shodan favicon hash

	//fmt.Println("[!] Shodan Favicon Hash:", faviconBytes)
	fmt.Println("[!] http.favicon.hash:", faviconBytes)
	fmt.Printf("[*] View Results:\n> https://www.shodan.io/search?query=http.favicon.hash%%3A%d\n", faviconBytes)

}

// getFaviconURL fetches the website's HTML and extracts the favicon URL.
func getFaviconURL(websiteURL string) (string, error) {
	// Allow insecure TLS connections for websites without SSL/TLS
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	response, err := client.Get(websiteURL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	// Search for the favicon URL in the HTML
	faviconURL := findFaviconURL(body)

	return faviconURL, nil
}

// findFaviconURL extracts the favicon URL from the HTML.
func findFaviconURL(html []byte) string {
	// Search for the favicon link tag in the HTML
	const faviconTag = `<link rel="icon" href="`
	startIndex := bytes.Index(html, []byte(faviconTag))
	if startIndex == -1 {
		return ""
	}

	// Find the end of the URL
	endIndex := bytes.Index(html[startIndex+len(faviconTag):], []byte(`"`))
	if endIndex == -1 {
		return ""
	}

	// Extract the URL
	faviconURL := html[startIndex+len(faviconTag) : startIndex+len(faviconTag)+endIndex]

	// Convert to string
	return string(faviconURL)
}

// fetchFavicon fetches the favicon file using the given URL.
func fetchFavicon(url string) int32 {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	final := ""
	fix := 76
	s := make([]string, 0)

	f, errs := http.Get(url)
	if errs != nil {
		fmt.Println("[-] Error On Loading Page ", errs)
	}
	content, err := ioutil.ReadAll(f.Body)
	if err != nil {
		fmt.Println("[-] Error On Loading Page Contents", errs)
	}
	str := base64.StdEncoding.EncodeToString(content)

	// slice up string
	for i := 0; i*fix+fix < len(str); i++ {
		it := str[i*fix : i*fix+fix]
		s = append(s, it)
	}

	// find last piece of string
	findlen := len(s) * fix
	last := str[findlen:] + "\n"

	// put it all together
	for _, s := range s {
		final = final + s + "\n"
	}
	str = final + last

	// do murmurhash3 stuff
	mm3 := murmur3.StringSum32(str)

	// convert uint32 to int32
	return int32(mm3)
}
