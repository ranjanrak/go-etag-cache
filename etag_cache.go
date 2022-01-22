package etagcache

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// AddEtag adds etag to request header
func (c *Client) AddEtag(method string, headers http.Header, url string) http.Header {
	if method == "GET" && c.exists(c.etagPath) {
		etagResult := c.readFile(c.etagPath)
		// Add etag to request header
		headers.Add("If-None-Match", etagResult[url])
	}
	return headers
}

// SaveEtag stores Etag data to etag.gob file
func (c *Client) SaveEtag(headers http.Header, url string) {
	var etag = make(map[string]string)
	if len(headers["Etag"]) > 0 {
		if c.exists(c.etagPath) {
			// If exists update the etag value
			etag = c.readFile(c.etagPath)
		} else {
			// Create client directory if it doesn't exists
			os.Mkdir(c.dirPath, os.ModePerm)
		}
		etag[url] = headers["Etag"][0]
		err := c.writeFile(c.etagPath, etag)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// WriteReadCache write/update and return response/cache data based on HTTP status code
func (c *Client) WriteReadCache(res *http.Response, url string) string {
	var response = make(map[string]string)
	// save etag from response header
	c.SaveEtag(res.Header, url)
	if res.StatusCode == 304 {
		response = c.readFile(c.cachePath)
	} else {
		// render response body for 200
		body, _ := ioutil.ReadAll(res.Body)
		if c.exists(c.cachePath) {
			// append/update the cache if file already available
			response = c.readFile(c.cachePath)
		}
		// update cache for that specific key
		response[url] = string(body)
	}
	// update the cache file with new response
	err := c.writeFile(c.cachePath, response)
	if err != nil {
		fmt.Println(err)
	}
	return response[url]
}

// readFile decodes gob encoded data from .gob file
func (c *Client) readFile(filePath string) map[string]string {
	var fileOutput map[string]string
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	} else {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(&fileOutput)
	}
	file.Close()
	return fileOutput
}

// writeFile creates and writes gob encoded data to the .gob file
func (c *Client) writeFile(filePath string, fileData map[string]string) error {
	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(fileData)
	}
	file.Close()
	return err
}

// exists returns boolean whether the given file or directory(path) exists
func (c *Client) exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
