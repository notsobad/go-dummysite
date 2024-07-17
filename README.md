# Go-Fakesite

Go-Fakesite is a mock website generator for HTTP testing. It is a Go version of [ynm3k](https://github.com/notsobad/ynm3k). This tool is designed to simulate various URL addresses, static and dynamic web pages, different status codes, response sizes, and more. It is particularly useful for testing CDN and WAF systems.

## Installation

To install Go-Fakesite, use the following command:

```bash
go install github.com/notsobad/go-fakesite
```

## Usage

To run Go-Fakesite, use the following command:

```bash
go-fakesite
```

You can then access the site at [http://localhost:9527/](http://localhost:9527/).

## Features

### Static Files

Visiting the same URL will yield the same result. Use the following format:

`/static/$RANDOM.$EXT`

Examples:

* [/static/abc.zip](/static/abc.zip)
* [/static/123/abc.css](/static/123/abc.css)
* [/static/xyz.html](/static/xyz.html)
* [/static/1234.js](/static/1234.js)

### Dynamic URLs

Visiting the same URL will yield different results. Use the following format:

`/dynamic/$RANDOM.$EXT`

Examples:

* [/dynamic/abc.php](/dynamic/abc.php)
* [/dynamic/abc.jsp](/dynamic/abc.jsp)

### HTTP Status Codes

You can simulate different HTTP status codes. Use the following format:

`/code/$CODE.$EXT`

Examples:

* [/code/200.html](/code/200.html)
* [/code/403.asp](/code/403.asp)
* [/code/404.asp](/code/404.asp)
* [/code/500.php](/code/500.php)
* [/code/502.asp](/code/502.asp)

### Specified Size Response

You can output a file of a specified size. Use the following format:

`/size/$SIZE.$EXT`

Examples:

* [/size/11k.zip](/size/11k.zip)
* [/size/1m.bin](/size/1m.bin)
* [/size/1024.rar](/size/1024.rar)

### Slow Response Server

You can simulate a slow server response. Visit `/slow/$SECONDS` and the URL will take $SECONDS time to render.

Examples:

* [/slow/3](/slow/3)
* [/slow/10](/slow/10)

### URL Redirect

You can simulate various URL redirect methods. Use the following format:

`http://localhost:9527/redirect/$CODE?url=$URL`

Normal redirect:

* [/redirect/301?url=http://www.notsobad.vip](/redirect/301?url=http://www.notsobad.vip)
* [/redirect/302?url=http://www.notsobad.vip](/redirect/302?url=http://www.notsobad.vip)
* [/redirect/js?url=http://www.notsobad.vip](/redirect/js?url=http://www.notsobad.vip)
* [/redirect/meta?url=http://www.notsobad.vip](/redirect/meta?url=http://www.notsobad.vip)

You can also redirect to a 'file://' protol address:

```bash
curl -v 'localhost:9527/redirect/301?url=file:///etc/passwd'
```
### Chunked http response

You can simulate chunked http response. Use the following format, `$COUNT` is the count of chunked message. 

`http://localhost:9527/chunk/$COUNT`

Examples:

* [/chunk/3](/chunk/3)
* [/chunk/10](/chunk/10)
* [/chunk/999](/chunk/999)

### HTTP Trace
You can view raw http request and request body in the `/trace` api.

* [/trace](/trace)

in cmd line:

```
# post
curl 127.1:9527/trace -d 'a=1'
```