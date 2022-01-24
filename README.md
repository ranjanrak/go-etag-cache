# go-etag-cache
Go package for caching HTTP response based on etag.

## Concept
Store gob encoded response header etag into `etag` and response body into `cache` gob files.<br>
Add `If-None-Match` header for all `GET` request.</br> 
Update `etag` and `cache` for required request url, when none 304 http status is received i.e if response is modified.</br>
Return `cache` data for non-modified response.</br>

Sample `GET` request header with `If-None-Match`:
```
map[Authorization:[token api_key:access_token] 
If-None-Match:[W/"sBBM4eIM4NDlAf4R"] 
User-Agent:[gokiteconnect/4.0.2] 
X-Kite-Version:[3]]
```

## Installation
```
go get -u github.com/ranjanrak/go-etag-cache
```

## Usage
```go
package main

import (
	"fmt"
	"net/http"

	etagcache "github.com/ranjanrak/go-etag-cache"
)

func main() {
    // Create new cache instance
    cacheClient := etagcache.New(etagcache.ClientParam{DirPath: ""})
    url := "https://api.kite.trade/orders"
    
    req, _ := http.NewRequest("GET", url, nil)
    
    req.Header.Add("User-Agent", "gokiteconnect/4.0.2")
    req.Header.Add("x-kite-version", "3")
    req.Header.Add("authorization", "token api_key:access_token")
    
    // Add etag to request header
    req = cacheClient.AddEtag(req)
    
    res, _ := http.DefaultClient.Do(req)
    defer res.Body.Close()
    
    // Store and fetch response/cache data
    response := cacheClient.HandleCache(res, url)
    fmt.Println(response)
}
```

## Response
```
&{Status:304 Not Modified StatusCode:304 Proto:HTTP/2.0 ProtoMajor:2 ProtoMinor:0 
Header:map[Cf-Cache-Status:[DYNAMIC] Cf-Ray:[6d179aa4dcf131d5-BOM] Content-Type:[application/json] 
Date:[Sat, 22 Jan 2022 09:02:33 GMT] Expect-Ct:[max-age=604800, 
report-uri="https://report-uri.cloudflare.com/cdn-cgi/beacon/expect-ct"] Server:
[cloudflare] Strict-Transport-Security:[max-age=15552000; includeSubDomains]] 
Body:{Reader:0xc00009aa80} ContentLength:0 TransferEncoding:[] 
Close:false Uncompressed:false Trailer:map[] Request:0xc000108000 TLS:0xc00030e9a0}

{"status":"success","data":[{"order_id":"XXXXXX","exchange_order_id":null,"parent_order_id":null,
"status":"CANCELLED AMO","status_message":null,"status_message_raw":null,"order_timestamp":"2022-01-22 
07:18:35","exchange_update_timestamp":null,"exchange_timestamp":null,"variety":"amo","exchange":"NSE",
"tradingsymbol":"BHEL","instrument_token":112129,"order_type":"MARKET","transaction_type":"BUY",
"validity":"DAY","product":"CNC","quantity":1,"disclosed_quantity":0,...}