package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {

	sharedUrl := flag.String("url", "", "The OneDrive Shared URL.")
	outfile := flag.String("o", "", "Name of the output file.")
	flag.Parse()

	if *sharedUrl == "" || *outfile == "" {
		fmt.Println("mandatory arguments -url -o")
		os.Exit(1)
	}

	base64 := base64.StdEncoding.EncodeToString([]byte(*sharedUrl))
	encodedUrl := "u!" + strings.Trim(base64, "=")
	encodedUrl = strings.Replace(encodedUrl, "/", "_", -1)
	encodedUrl = strings.Replace(encodedUrl, "+", "-", -1)

	url := "https://api.onedrive.com/v1.0/shares/" + encodedUrl + "/root/content"

	fmt.Printf("Download url: %s", url)

	out, err := os.Create(*outfile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer out.Close()

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

}
