# go-fakesite
A Fakesite for http test, it's a go version of [ynm3k](https://github.com/notsobad/ynm3k).

    go install github.com/notsobad/go-fakesite
    # run
    go-fakesite

    # visit http://localhost:9527/
 


# Static file

Visit same url, get same result

/static/$RANDOM.$EXT

* http://localhost:9527/static/abc.zip
* http://localhost:9527/static/xyz.html
* http://localhost:9527/static/1234.js

## Dynamic url

Visit same url, get different result

/dynamic/$RANDOM.$EXT

* http://localhost:9527/dynamic/abc.php
* http://localhost:9527/dynamic/abc.jsp

## HTTP status code
/code/$CODE.$EXT

* http://localhost:9527/code/500.php
* http://localhost:9527/code/404.asp

## Specified size response
You can output a file of the specified size.

/size/$SIZE.$EXT

* http://localhost:9527/size/11k.zip
* http://localhost:9527/size/1m.bin
* http://localhost:9527/size/1024.rar

## A server with a slow response

Visit `/slow/$SECONDS`, the url will take $SECONDS time to render.

* http://localhost:9527/slow/3
* http://localhost:9527/slow/4-10


## URL redirect
All kinds of url redirect method

* http://localhost:9527/redirect/301?url=http://www.notsobad.vip  301
* http://localhost:9527/redirect/302?url=http://www.notsobad.vip  302
* http://localhost:9527/redirect/js?url=http://www.notsobad.vip javascript
* http://localhost:9527/redirect/meta?url=http://www.notsobad.vip html meta

DEMO:

    curl -v 'localhost:9527/redirect/301?url=file:///etc/passwd'
    curl -v 'localhost:9527/redirect/302?url=http://www.jiasule.com'
