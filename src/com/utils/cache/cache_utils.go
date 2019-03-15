package cache

import (
	"com/models"
	"com/utils/pool"
	"sync"
)

var cachePool pool.Scheduler
var DEFAULT_THREAD_POOL_SIZE = 100

var PROPERTY_DATA_CACHE_DIR string
var PROPERTY_DATA_CACHE_DISK_MAX_SIZE int64
var PROPERTY_DATA_CACHE_ARCHIVE_MAX_SIZE int64
var PROPERTY_DATA_CACHE_ARCHIVE_LIFE_CYCLE int64

func init() {
	cachePool, _ = pool.NewBlocking(DEFAULT_THREAD_POOL_SIZE)
}

func PutToCache(putRecordsRequest *models.PutRecordsRequest, cacheResendConf *models.CacheResendConf) {
	PROPERTY_DATA_CACHE_DIR = cacheResendConf.DataCacheDir
	PROPERTY_DATA_CACHE_DISK_MAX_SIZE = cacheResendConf.DataCacheDiskMaxSize
	PROPERTY_DATA_CACHE_ARCHIVE_MAX_SIZE = cacheResendConf.DataCacheArchiveSize
	PROPERTY_DATA_CACHE_ARCHIVE_LIFE_CYCLE = cacheResendConf.DataCacheArchiveLifeCycle

	wg := &sync.WaitGroup{}
	wg.Add(1)
	cachePool.Add(pool.Job{
		Run: func(args ...interface{}) {
			if putRecordsRequest, ok := args[0].(*models.PutRecordsRequest); ok {
				PutDataToCache(putRecordsRequest)
			}
			if wg, ok := args[1].(*sync.WaitGroup); ok {
				wg.Done()
			}
		},
		Args: []interface{}{putRecordsRequest, wg},
	})
	wg.Wait()
}