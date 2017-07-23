package cachetree

type CacheTreeConfig struct {
	KeyLifeTimeSec int    `json:"key_lifetime"`
	RequestTimeout int    `json:"request_timeout"`
	BlobPath       string `json:"blob_path"`
}
