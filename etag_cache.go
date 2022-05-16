package etagcache

import (
	"encoding/gob"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// ClientParam represent input cache dir path from user
type ClientParam struct {
	DirPath string
}

// Client represents interface for etag cache
type Client struct {
	dirPath   string
	etagPath  string
	cachePath string
}

// RequestParam represents input http request params
type RequestParam struct {
	ReqClient *http.Request
	Headers   map[string]string
}

// Creates new cache client
func New(userParam ClientParam) *Client {
	// Set default dir path
	if userParam.DirPath == "" {
		userParam.DirPath = "./goapp"
	}
	// Set etag file path
	etagPath := filepath.Join(userParam.DirPath, "/etag.gob")
	// Set cache file path
	cachePath := filepath.Join(userParam.DirPath, "/cache.gob")

	return &Client{
		dirPath:   userParam.DirPath,
		etagPath:  etagPath,
		cachePath: cachePath,
	}
}

// AddEtag adds etag to request header
func (c *Client) AddEtag(reqParam RequestParam) *http.Request {
	req := reqParam.ReqClient
	// Add all headers to http request
	for name, value := range reqParam.Headers {
		req.Header.Add(name, value)
	}

	if req.Method == "GET" && c.exists(c.etagPath) {
		etagResult := c.readFile(c.etagPath)
		// Add etag to request header
		req.Header.Add("If-None-Match", etagResult[req.URL.String()])
	}
	return req
}

// SaveEtag stores Etag data to etag.gob file
func (c *Client) SaveEtag(headers http.Header, url string) {
	var etag = make(map[string]string)
	if c.exists(c.etagPath) {
		// If older etag file exists update the etag value
		etag = c.readFile(c.etagPath)
	} else {
		// Create client directory if it doesn't exists
		os.Mkdir(c.dirPath, os.ModePerm)
	}
	if len(headers["Etag"]) > 0 {
		// update etag
		etag[url] = headers["Etag"][0]
		err := c.writeFile(c.etagPath, etag)
		if err != nil {
			log.Println(err)
		}
	}
}

// HandleCache write/update and return response/cache data based on HTTP status code
func (c *Client) HandleEtagCache(reqParam RequestParam, url string) (string, error) {
	var response = make(map[string]string)
	// add etag to request header
	req := c.AddEtag(reqParam)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()

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
	return response[url], err
}

// readFile decodes gob encoded data from .gob file
func (c *Client) readFile(filePath string) map[string]string {
	var fileOutput map[string]string
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
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
