# Tiny Url

A wonderfully simple url shortener written in Go.
Converts base10 -> base62 for the shortening algorithm (https://github.com/arrwhidev/base-converter).

Uses a simple map as the database, keeping it simple.

## Usage

 ```
    ➜ curl -d '{"url":"http://google.co.uk"}' -H "Content-Type:application/json" localhost:8080
    // returns {"url":"http://localhost:8080/4z"}

    ➜ curl localhost:8080/4z
    // returns redirect to http://google.co.uk
```