# :globe_with_meridians: Go-Fakesite

Go-Fakesite is a mock website generator for HTTP testing. It is a Go version of [ynm3k](https://github.com/notsobad/ynm3k). This tool is designed to simulate various URL addresses, static and dynamic web pages, different status codes, response sizes, and more. It is particularly useful for testing CDN and WAF systems.

## :wrench: Installation

To install Go-Fakesite, use the following command:

```bash
go install github.com/notsobad/go-fakesite
```

## :computer: Usage

To run Go-Fakesite, use the following command:

```bash
go-fakesite
```

You can then access the site at `http://localhost:9527/`.

## :star2: Features

### :file_folder: Static Files

Visiting the same URL will yield the same result. Use the following format:

`/static/$RANDOM.$EXT`

Examples:

* http://localhost:9527/static/abc.zip
* http://localhost:9527/static/xyz.html
* http://localhost:9527/static/1234.js

### :cyclone: Dynamic URLs

Visiting the same URL will yield different results. Use the following format:

`/dynamic/$RANDOM.$EXT`

Examples:

* http://localhost:9527/dynamic/abc.php
* http://localhost:9527/dynamic/abc.jsp

### :warning: HTTP Status Codes

You can simulate different HTTP status codes. Use the following format:

`/code/$CODE.$EXT`

Examples:

* http://localhost:9527/code/500.php
* http://localhost:9527/code/404.asp

### :chart_with_upwards_trend: Specified Size Response

You can output a file of a specified size. Use the following format:

`/size/$SIZE.$EXT`

Examples:

* http://localhost:9527/size/11k.zip
* http://localhost:9527/size/1m.bin
* http://localhost:9527/size/1024.rar

### :snail: Slow Response Server

You can simulate a slow server response. Visit `/slow/$SECONDS` and the URL will take $SECONDS time to render.

Examples:

* http://localhost:9527/slow/3
* http://localhost:9527/slow/4-10

### :arrows_counterclockwise: URL Redirect

You can simulate various URL redirect methods. Use the following format:

`http://localhost:9527/redirect/$CODE?url=$URL`

Normal redirect:

* http://localhost:9527/redirect/301?url=http://www.notsobad.vip
* http://localhost:9527/redirect/302?url=http://www.notsobad.vip
* http://localhost:9527/redirect/js?url=http://www.notsobad.vip
* http://localhost:9527/redirect/meta?url=http://www.notsobad.vip

You can also redirect to a 'file://' protol address:

```bash
curl -v 'localhost:9527/redirect/301?url=file:///etc/passwd'
```