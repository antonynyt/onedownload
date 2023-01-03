package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

/*
Todo
- clear the sharedUrl of any ?parameters
- Don't download if file already exists
- Don't print url if status != 200ok
- Download Folder?
*/

const (
	apiURL = "https://api.onedrive.com/v1.0/"
)

type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Size int64  `json:"size"`
	URL  string `json:"@content.downloadUrl"`
}

func main() {
	//set the flags
	onlyURL := flag.Bool("nodownload", false, "Optional: generate only the download link.")
	outfile := flag.String("o", "", "Set the output `filename`. Example: (image.png)")
	sharedUrl := flag.String("url", "", "The OneDrive Share `link`.")
	flag.Parse()

	//if no flags set throw an error
	if !*onlyURL {
		if *sharedUrl == "" {
			flag.Usage()
			os.Exit(1)
		}
	} else if *onlyURL && *sharedUrl != "" {
		fmt.Printf("[-] Full link: %s\n", formatURL(sharedUrl))
		return
	} else {
		flag.Usage()
		os.Exit(1)
	}

	//format the shared OneDrive Url
	url := formatURL(sharedUrl)

	//Print the full link
	//fmt.Printf("[-] Full link: %s\n", url)

	//
	item, err := getItem(url)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if *outfile == "" {
		*outfile = item.Name
	}

	err = downloadItem(item.URL, *outfile)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func formatURL(sharedUrl *string) string {
	//genrate the download url from OneDrive API
	base64 := base64.StdEncoding.EncodeToString([]byte(*sharedUrl))
	encodedUrl := "u!" + strings.Trim(base64, "=")
	encodedUrl = strings.Replace(encodedUrl, "/", "_", -1)
	encodedUrl = strings.Replace(encodedUrl, "+", "-", -1)

	//only for files /DriveItem
	//?$expand=children show the childrens
	url := apiURL + "shares/" + encodedUrl + "/driveItem"

	return url
}

func getItem(url string) (*Item, error) {
	//HTTP GET
	// Set up the request.
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Send the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the status code.
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get item: %s", resp.Status)
	}

	// Decode the response body.
	var item Item
	err = json.NewDecoder(resp.Body).Decode(&item)
	if err != nil {
		return nil, fmt.Errorf("failed to decode json: %s", resp.Status)
	}

	return &item, nil
}

func downloadItem(url, filename string) error {

	// Download the file.
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the status code.
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: %s", resp.Status)
	}

	//total size of the file
	totalBytes := resp.ContentLength

	//create the outfile
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	go downloadProgress(filename, totalBytes)

	//Write the contents of the response to the output file.
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func downloadProgress(path string, totalSize int64) error {

	// Size of the progress bar
	const barWidth = 40

	// Set up a timer to update the progress every 100ms.
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	// Open the file.
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	for {

		// Get the size of the file.
		fi, err := file.Stat()
		if err != nil {
			return err
		}
		size := fi.Size()

		// Update the progress. downloading [=======--------] 50%
		select {
		case <-ticker.C:
			percent := float64(size) / float64(totalSize) * 100
			barSize := int(percent) * barWidth / 100
			bar := strings.Repeat("=", barSize) + strings.Repeat("-", 40-barSize)
			fmt.Printf("Downloading [%s] %.0f%%\r", bar, percent)
		default:
		}

		// Break out of the loop when the file is fully downloaded.
		if size == totalSize {
			break
		}
	}

	fmt.Println()

	return nil
}
