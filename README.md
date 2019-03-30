# go-fakesite
A Fakesite for http test, it's a go version of [ynm3k](https://github.com/notsobad/ynm3k).

    go install github.com/notsobad/go-fakesite
    # run
    go-fakesite

    # visit http://127.0.0.1:8080/



# Static file

Visit same url, get same result

/static/$RANDOM.$EXT

* http://ynm3k.notsobad.vip/static/abc.zip
* http://ynm3k.notsobad.vip/static/xyz.html
* http://ynm3k.notsobad.vip/static/1234.js

## Dynamic url

Visit same url, get different result

/dynamic/$RANDOM.$EXT

* http://ynm3k.notsobad.vip/dynamic/abc.php
* http://ynm3k.notsobad.vip/dynamic/abc.jsp

## HTTP status code
/code/$CODE.$EXT

* http://ynm3k.notsobad.vip/code/500.php
* http://ynm3k.notsobad.vip/code/404.asp

## Specified size response
You can output a file of the specified size.

/size/$SIZE.$EXT

* http://ynm3k.notsobad.vip/size/11k.zip
* http://ynm3k.notsobad.vip/size/1m.bin
* http://ynm3k.notsobad.vip/size/1024.rar

## A server with a slow response

Visit `/slow/$SECONDS`, the url will take $SECONDS time to render.

* http://ynm3k.notsobad.vip/slow/3
* http://ynm3k.notsobad.vip/slow/4-10


## URL redirect
All kinds of url redirect method

* http://ynm3k.notsobad.vip/redirect/301?url=http://www.notsobad.vip  301
* http://ynm3k.notsobad.vip/redirect/302?url=http://www.notsobad.vip  302
* http://ynm3k.notsobad.vip/redirect/js?url=http://www.notsobad.vip javascript
* http://ynm3k.notsobad.vip/redirect/meta?url=http://www.notsobad.vip html meta

DEMO:

    curl -v 'localhost:9527/redirect/301?url=file:///etc/passwd'
    curl -v 'localhost:9527/redirect/302?url=http://www.jiasule.com'
