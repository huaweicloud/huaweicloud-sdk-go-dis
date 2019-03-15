package cache

import (
	"os"
	"time"
	"com/models"
	"path/filepath"
	"strings"
	"encoding/json"
	"strconv"
	"io/ioutil"
	"sync"
	"com/logger"
)

var CACHE_FILE_PREFIX = "dis-cache-data-"

//缓存临时文件后缀
var CACHE_TMP_FILE_SUFFIX = ".tmp";

//缓存归档数据文件后缀
var CACHE_ARCHIVE_DATA_FILE_SUFFIX = ".data";

//缓存归档索引文件后缀
var CACHE_ARCHIVE_INDEX_FILE_SUFFIX = ".index";

var MB int64 = 1024 * 1024;

var c CacheManager
var l sync.Mutex

type CacheManager struct {
	DataTmpFileName   string
	TmpFileCreateTime int64
}

func PutDataToCache(putRecordsRequest *models.PutRecordsRequest) {
	initCacheManager()

	bytes, _ := json.Marshal(putRecordsRequest)
	data := string(bytes)
	l.Lock()
	if needToArchive(data) {
		archive()
	}

	if hasEnoughSpace(data) {
		WriteToFile(data)
	} else {
		logger.LOG(logger.INFO, "Put to cache failed, cache space is not enough, configured max dir size: %d", PROPERTY_DATA_CACHE_DISK_MAX_SIZE)
	}
	l.Unlock()
}

func hasEnoughSpace(data string) bool {
	dataSize := getDataSize(data)

	cacheDir, err := getCacheDir()
	if err != nil {
		logger.LOG(logger.WARNING, "get cache dir [%s] error: %s", cacheDir, err)
		return false
	}

	logger.LOG(logger.DEBUG, "dataSize: [%d], cacheDir size: [%d]", dataSize, readDirSize(cacheDir))
	if (dataSize + readDirSize(cacheDir)) / MB > PROPERTY_DATA_CACHE_DISK_MAX_SIZE {
		return false
	}

	return true
}

func needToArchive(data string) bool {
	dataSize := getDataSize(data)

	fileInfo, err := os.Stat(c.DataTmpFileName)
	if err != nil {
		return false
	}
	fileSize := fileInfo.Size()

	if (dataSize + fileSize) / MB > PROPERTY_DATA_CACHE_ARCHIVE_MAX_SIZE || (time.Now().UTC().Unix() - c.TmpFileCreateTime / 1000) > PROPERTY_DATA_CACHE_ARCHIVE_LIFE_CYCLE {
		return true
	}

	return false
}

func getDataSize(data string) int64 {
	return int64(len([]byte(data)))
}

func archive() {
	logger.LOG(logger.INFO, "need to archive tmp file: ", c.DataTmpFileName)
	if len(c.DataTmpFileName) == 0 {
		return
	}

	archiveDataFilename := strings.Replace(c.DataTmpFileName, CACHE_TMP_FILE_SUFFIX, "", 1) + CACHE_ARCHIVE_DATA_FILE_SUFFIX;
	archiveIndexFilename := strings.Replace(c.DataTmpFileName, CACHE_TMP_FILE_SUFFIX, "", 1) + CACHE_ARCHIVE_INDEX_FILE_SUFFIX;
	dataTmpFileName := c.DataTmpFileName

	err := os.Rename(dataTmpFileName, archiveDataFilename)
	if err != nil {
		logger.LOG(logger.INFO, "archive tmp file [%s] error: %s", c.DataTmpFileName, err)
		return
	}

	f, err := os.Create(archiveIndexFilename)
	if err != nil {
		logger.LOG(logger.ERROR, "create archiveIndex [%s] error: %s", archiveIndexFilename, err)
		os.Remove(archiveDataFilename)
		return
	}
	f.Close()

	// 重置缓存临时文件
	reset();
	initCacheManager();
	logger.LOG(logger.INFO, "archive tmp file [%s] success", c.DataTmpFileName)
}

func initCacheManager() {
	if len(c.DataTmpFileName) == 0 {
		// 生成缓存数据文件和缓存索引文件
		cacheDir, err := getCacheDir()
		if err != nil {
			return
		}

		//检查是否有tmp文件，存在的话则在旧tmp文件上继续写入
		cacheTmpDataFileName := checkTmpFileExists();
		timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)[:13]

		if cacheTmpDataFileName == "" {
			cacheTmpDataFileName := CACHE_FILE_PREFIX + timestamp + CACHE_TMP_FILE_SUFFIX
			c.DataTmpFileName = cacheDir + string(filepath.Separator) + cacheTmpDataFileName
			f, err := os.Create(c.DataTmpFileName)
			if err != nil {
				logger.LOG(logger.ERROR, "failed to create data tmp file [%s] error", cacheTmpDataFileName, err)
				return
			}
			f.Close()
		} else {
			c.DataTmpFileName = cacheDir + string(filepath.Separator) + cacheTmpDataFileName
		}
		createTime, _ := strconv.ParseInt(timestamp, 10, 64)
		c.TmpFileCreateTime = createTime

	}
}

func reset() {
	c.DataTmpFileName = ""
}


/**
    * 获取配置的缓存目录路径
    * @return 存放缓存文件的目录
    */
func getCacheDir() (string, error) {
	dataCacheDir := PROPERTY_DATA_CACHE_DIR

	fi, err := os.Stat(dataCacheDir)
	if err != nil || !fi.IsDir() {
		err1 := os.MkdirAll(dataCacheDir, 0755)
		if err1 != nil {
			logger.LOG(logger.ERROR, "mkdir dataCacheDir [%s] error: %s", dataCacheDir, err1)
			return dataCacheDir, err1
		}

		return dataCacheDir, nil
	}

	return dataCacheDir, nil
}

/**
 * 写入缓存文件
 * @param data 待写入缓存文件的数据
 */
func WriteToFile(data string) error {
	// 以只写的模式，打开文件
	logger.LOG(logger.DEBUG, "write data to tmp file [%s]", c.DataTmpFileName)
	f, err := os.OpenFile(c.DataTmpFileName, os.O_RDWR | os.O_APPEND, 0666)
	if err != nil {
		logger.LOG(logger.ERROR, "write data to tmp file [%s] error: ", c.DataTmpFileName, err)
	} else {
		// 从末尾的偏移量开始写入内容
		_, err = f.Write([]byte(data))
		_, err = f.Write([]byte(models.LINE_SEPARATOR))
	}
	defer f.Close()
	return err
}

func readDirSize(dirPath string) int64 {
	var dirSize int64 = 0
	flist, e := ioutil.ReadDir(dirPath)
	if e != nil {
		logger.LOG(logger.WARNING, "read cache dir [%s], error: ", dirPath, e)
		return dirSize
	}

	for _, f := range flist {
		if f.IsDir() {
			dirSize += readDirSize(dirPath + string(filepath.Separator) + f.Name())
		} else {
			dirSize += f.Size()
		}
	}

	return dirSize
}

func checkTmpFileExists() string {
	fileExistName := ""
	flist, e := ioutil.ReadDir(PROPERTY_DATA_CACHE_DIR)
	if e != nil {
		logger.LOG(logger.WARNING, "check tmp file [%s] exists error: %s", PROPERTY_DATA_CACHE_DIR, e)
		return fileExistName
	}

	for _, f := range flist {
		if !f.IsDir() && strings.HasSuffix(f.Name(), CACHE_TMP_FILE_SUFFIX)&& strings.HasPrefix(f.Name(), CACHE_FILE_PREFIX) {
			logger.LOG(logger.INFO, "check tmp file [%s] exists and continue to write data to the tmp file.", f.Name())
			return f.Name()
		}
	}

	return fileExistName
}