package etagcache

import "path/filepath"

// ClientParam represent input cache dir path from user
type ClientParam struct {
	DirPath string
}

type Client struct {
	dirPath   string
	etagPath  string
	cachePath string
}

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
