package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
)

/*
Todo
- Clear the code, better error handling
- Better http download
- Show the current status of the download like in curl + estimated time (content length)
- Get the http body and the filename from the head
- Don't print url if status != 200ok
- Separate code in functions
*/

func main() {
	//set the flags
	sharedUrl := flag.String("url", "", "The OneDrive Share `link`.")
	outfile := flag.String("o", "", "Set the output `filename`. E.g. video.mkv.")
	flag.Parse()

	//if no flags set throw an error
	if *sharedUrl == "" || *outfile == "" {
		flag.Usage()
		os.Exit(0)
	}

	printBanner()

	//format the shared OneDrive Url
	url := linkFormatter(sharedUrl)

	//create the outfile
	out, err := os.Create(*outfile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer out.Close()

	//HTTP GET

	//GET THE HTTP HEADER
	headResp, err := http.Head(url)
	if err != nil {
		panic(err)
	}
	if headResp.StatusCode != 200 {
		panic("stop")
	}

	//GET THE CONTENT
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	//Print the full link
	fmt.Printf("[-] Full link: %s\n", url)

	//COPY THE CONTENT OF THE GET INTO THE FILE

	// size, err = io.Copy(out, resp.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	fmt.Println()
}

func linkFormatter(sharedUrl *string) string {
	//genrate the download url from OneDrive API
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

/*

&{200 OK 200 HTTP/2.0 2 0 map[Accept-Ranges:[bytes] Cache-Control:[public] Content-Disposition:[attachment; filename="1-221.png"] Content-Length:[9635415] Content-Location:[REDACTED] Content-Type:[image/png] Ctag:[] Date:[Tue, 06 Sep 2022 20:43:54 GMT] Etag:[] Expires:[Mon, 05 Dec 2022 20:43:53 GMT] Last-Modified:[Thu, 10 Jun 2021 19:14:03 GMT] Ms-Cv:[.0] P3p:[CP="BUS CUR CONo FIN IVDo ONL OUR PHY SAMo TELo"] Strict-Transport-Security:[max-age=31536000; includeSubDomains] X-Asmversion:[UNKNOWN; 19...2005] X-Cache:[CONFIG_NOCACHE] X-Content-Type-Options:[nosniff] X-Msedge-Ref:[Ref A:  Ref B:  Ref C: 2022-09-06T20:43:51Z] X-Msnserver:[] X-Preauthinfo:[rv;poba;] X-Sqldataorigin:[S] X-Streamorigin:[X]] {} 9635415 [] false false map[]  }*/
