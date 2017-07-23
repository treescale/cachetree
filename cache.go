package cachetree

import (
	"github.com/boltdb/bolt"
	"strconv"
	"time"
)

var (
	cacheDB             *bolt.DB
	filesBucketName     = []byte("files")
	keyTimersBucketName = []byte("keyTimer")
)

func GetFile(filename string) (data []byte) {
	cacheDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(filesBucketName)
		data = b.Get([]byte(filename))
		return nil
	})

	return data
}

func PutFile(filename string, data []byte) error {
	err := cacheDB.Update(func(tx *bolt.Tx) error {
		filename_data := []byte(filename)
		err := tx.Bucket(filesBucketName).Put(filename_data, data)
		if err != nil {
			return err
		}

		err = tx.Bucket(keyTimersBucketName).Put(filename_data,
			[]byte(strconv.FormatInt(time.Now().UTC().Unix(), 10)))

		// if we have error inserting timer data, just deleting cache file
		if err != nil {
			tx.Bucket(filesBucketName).Delete(filename_data)
		}

		return err
	})

	return err
}

func startClearTimer(timeoutSeconds int) {
	for {

		cacheDB.Update(func(tx *bolt.Tx) error {
			stampBucket := tx.Bucket(keyTimersBucketName)
			fileBucket := tx.Bucket(filesBucketName)
			if stampBucket == nil || fileBucket == nil {
				tx.CreateBucket(keyTimersBucketName)
				tx.CreateBucket(filesBucketName)
				return nil
			}

			minStamp := time.Now().UTC().Unix() - int64(timeoutSeconds)
			return stampBucket.ForEach(func(k, v []byte) error {
				stamp, err := strconv.ParseInt(string(v), 10, 64)
				if err != nil {
					stampBucket.Delete(k)
					fileBucket.Delete(k)
					return nil
				}

				if stamp < minStamp {
					stampBucket.Delete(k)
					fileBucket.Delete(k)
				}

				return nil
			})
		})

		time.Sleep(time.Second * time.Duration(timeoutSeconds))
	}
}
