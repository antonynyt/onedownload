package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
)

/*
Todo
- clear the code, better error handling
- better http download
- show the current status of the download like curl + estimated time (long term)
*/

func main() {
	//set the flags
	sharedUrl := flag.String("url", "", "The OneDrive Share `link`.")
	outfile := flag.String("o", "", "Set the output `filename`.")
	flag.Parse()

	//if no flags set throw an error
	if *sharedUrl == "" || *outfile == "" {
		flag.Usage()
		os.Exit(0)
	}

	printBanner()

	//format the shared OneDrive Url
	url := linkFormatter(sharedUrl)
	//Print the full link
	fmt.Printf("[-] Full link: %s\n", url)

	//create the outfile
	out, err := os.Create(*outfile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer out.Close()

	//download from the url
	resp, _ := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println()
}

func linkFormatter(sharedUrl *string) string {
	//genrate the download url
	base64 := base64.StdEncoding.EncodeToString([]byte(*sharedUrl))
	encodedUrl := "u!" + strings.Trim(base64, "=")
	encodedUrl = strings.Replace(encodedUrl, "/", "_", -1)
	encodedUrl = strings.Replace(encodedUrl, "+", "-", -1)

	url := "https://api.onedrive.com/v1.0/shares/" + encodedUrl + "/root/content"

	return url
}

func printBanner() {
	color.Cyan("--------------------")
	color.Cyan("OneDownloader v.0.1")
	color.Cyan("--------------------")
	println()
}
