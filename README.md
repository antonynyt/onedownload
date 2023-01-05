# OneDownload

CLI to download a file from a shared OneDrive link without an API key.
This can be useful for a server whitout a GUI. You just need a way to paste the shared url (use ssh).

## Requirements

You need Go installed to compile the code `go build main.go`

## How to use

1. Create a shared url from you onedrive account.
2. `./onedownload -url 'https://1drv.ms/x/s!ABCDEFG1234'` (escape the `!` with `\` if you use double quotes).
3. Wait for it...

## Future improvements

- [ ] Change the path (already possible with `-o` but need the full filename.png)
- [ ] Don't start if the filename already exists.
- [ ] Download a folder and create a zip file.
- [ ] Less ressource intensive
