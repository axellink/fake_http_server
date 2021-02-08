# Fake HTTP server

## Introduction

Hello everyone, I'm just making a little fake HTTP server in GO. It will just answer 200 on every request while logging the said request. The idea came from trying to epxloit open redirection flaws.

I know there are already some websites that let you do this, but first I didn't remember their addresses, and I wanted something to continue my path on GO learning.

## How does it works 

For now, you fire it up and it listens on port 8080.

There is only one endpoint `/api/log` but there will be more.

Make (or make a server do) a request to the server and it will print a JSON representation of the request on standard output. I chose JSON because I want it to be a REST API in future.

Here is a JSON exemple of a GET Request :
```
{
   "Date":"2021-02-08 15:43:45.287150321 +0100 CET m=+8.425803986",
   "Method":"GET",
   "Url":"/api/log/test",
   "Proto":"HTTP/1.1",
   "Headers":{
      "Accept":"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
      "Accept-Encoding":"gzip, deflate",
      "Accept-Language":"fr,fr-FR;q=0.8,en-US;q=0.5,en;q=0.3",
      "Connection":"keep-alive",
      "Dnt":"1",
      "Upgrade-Insecure-Requests":"1",
      "User-Agent":"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:85.0) Gecko/20100101 Firefox/85.0"
   },
   "Params":{
      "id":"prout",
      "yolo":"swag"
   },
   "Body":""
}
```

## What's next ?

I want some more feature to be implemented.

First, a capability of having an ID to seggregate different project for you.

Make sure that an ID is only valid 24 hours with everything related to it deleted after.

A REST API to make it possible to do things programatically like create an ID, get every requests on it, maybe filter requests in an ID.

Maybe a website that uses the API for a more graphical use.

And of course enhancing the request options present.

## And that's it
I just do this for fun, if you stumble onto this in anyway, do whatever you want to do with it.
