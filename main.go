package main

import (
		"flag"
		"fmt"
		"encoding/base64"
		"strings"
)

func main(){

sharedUrl := flag.String("url", "", "The OneDrive Shared URL.")
flag.Parse()

base64 := base64.StdEncoding.EncodeToString([]byte(*sharedUrl))
encodedUrl := "u!" + strings.Trim(base64, "=")
encodedUrl = strings.Replace(encodedUrl, "/", "_", -1)
encodedUrl = strings.Replace(encodedUrl, "+", "-", -1)

url := "https://api.onedrive.com/v1.0/shares/" + encodedUrl + "/root/content"

fmt.Println(url)

}

/*

$sharedUrl = urldecode($_GET['url']);
$base64 = base64_encode($sharedUrl);
$encodedUrl = "u!" . rtrim($base64, '=');
$encodedUrl = str_replace('/', '_', $encodedUrl);
$encodedUrl = str_replace('+', '-', $encodedUrl);

$final = sprintf('https://api.onedrive.com/v1.0/shares/%s/root/content', $encodedUrl);
header('Location:' . $final, true, 302);

*/
