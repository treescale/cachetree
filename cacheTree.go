package cachetree

import "github.com/boltdb/bolt"

type CacheTreeConfig struct {
	KeyLifeTimeSec int      `json:"key_lifetime"`
	RequestTimeout int      `json:"request_timeout"`
	BlobPath       string   `json:"blob_path"`
	Targets        []string `json:"targets"`
	ServerHost     string   `json:"server_host"`
}

func StartCachingService(config CacheTreeConfig) (err error) {
	cacheDB, err = bolt.Open(config.BlobPath, 0666, nil)
	if err != nil {
		return err
	}

	err = startCacheServer(config.ServerHost)
	if err != nil {
		return err
	}

	go startClearTimer(config.KeyLifeTimeSec)
	go memberConnector(config.Targets...)
	return nil
}
