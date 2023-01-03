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
)

/*
Todo
- Clear the code, better error handling
- Better http download
- Show the current status of the download like in curl + estimated time (content length)
- Get the http body and the filename from the head
- Don't print url if status != 200ok
- Separate code in functions
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
	fmt.Printf("[-] Full link: %s\n", url)

	//
	item, err := getItem(url)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if *outfile == "" {
		*outfile = item.Name
	}

	/*
		fmt.Println("Name:", item.Name)
		fmt.Println("Size:", item.Size)
		fmt.Println("URL:", item.URL)
	*/

	downloadItem(item.URL, *outfile)
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

	//create the outfile
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	//Write the contents of the response to the output file.
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
