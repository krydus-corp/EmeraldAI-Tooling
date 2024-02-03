# Emerald Tooling

## Packages

### CLI [pkg/cli]

Navigate to /cmd/cli and run `go run main.go` to see CLI help screen.

---

These are stand-alone packages that provide some additional, stand alone functionality. These tools are what the Emerald CLI call into for testing various functionality of the system. The current tooling is as follows:

### Fetch [pkg/fetch]

This tool provides fetching utility for various types of content. Currently, the _Fetcher Tool_ only includes support for images.

### Download [pkg/download]

This tool provides functionality for downloading a given URL to the filesystem or to a byte stream.

### Image [pkg/image]

This tool provides functionality for processing images. Current support features include:

- Re-format images to JPEG or PNG formats
- Resize images

## Tooling Examples

These commands to be run from the the cli cmd pkg.

> Note: The following examples require [jq](https://stedolan.github.io/jq/)

Execute a batched search on 'boats' and output JSON formatted metadata to a file.

```bash
./emld-cli fetch "boats" -t images | \
    jq -r '.results' > \
    search-results.json
```

Pipe search-results.json into the download command

```bash
cat search-results.json | \
    jq -r 'map(.img_src_b64) | join(",")' | \
    xargs ./emld-cli download -p ~/Desktop/images -b -w 4 -vv
```

Format a directory of downloaded content.

```bash
ls -d ~/Desktop/images/* | \
    ./emld-cli image --stream --replace -f "image/png" -x 200 -w 10

```

Execute the full pipeline.

```bash
./emld-cli fetch "boats" -t images --stream --pages 5 | \
    jq -r '.img_src_b64' | \
    ./emld-cli download --stream -b -w 4 -p ~/Desktop/images | \
    ./emld-cli image --stream --replace -f "image/jpeg" -x 200
```
