package client

import (
	"io/ioutil"
	"time"
	"com/models"
	"strings"
	"os"
	"bufio"
	"io"
	"com/utils"
	"strconv"
	"sync"
	"path/filepath"
	"com/logger"
)

//缓存归档数据文件后缀
var CACHE_ARCHIVE_DATA_FILE_SUFFIX = ".data";

//缓存归档索引文件后缀
var CACHE_ARCHIVE_INDEX_FILE_SUFFIX = ".index";

var PROPERTY_DATA_CACHE_DIR string

func AutoResendRecordsTimeTask(dis *Client) {
	interval := 10 * time.Second
	for {
		AutoResendRecords(dis)
		time.Sleep(interval)
		continue
	}
}

func AutoResendRecords(dis *Client) {
	PROPERTY_DATA_CACHE_DIR = dis.cacheResendConf.DataCacheDir

	//查询归档文件数目
	archiveFiles := checkArchiveSize()
	logger.LOG(logger.INFO, "check archive files size: %d", len(archiveFiles))

	if len(archiveFiles) == 0 {
		time.Sleep(time.Second)
		return
	}

	for i, archiveFile := range archiveFiles {
		logger.LOG(logger.INFO, "Sender-%d-%s",i,  archiveFile, " start...")
		wg := &sync.WaitGroup{}
		wg.Add(1)
		processResend(archiveFile, dis, wg)
		logger.LOG(logger.INFO, "Sender-%d-%s",i,  archiveFile, " end.")
		wg.Wait()
	}

}

// 读归档文件--> 标记index-->上传数据 --> 删除归档文件
func processResend(archiveFile string, dis *Client, wg *sync.WaitGroup) {
	fileName := PROPERTY_DATA_CACHE_DIR + string(filepath.Separator) + archiveFile
	archiveIndexname := strings.Replace(fileName, CACHE_ARCHIVE_DATA_FILE_SUFFIX, CACHE_ARCHIVE_INDEX_FILE_SUFFIX, 1);
	var lineIndex int64 = readArchiveIndex(archiveIndexname)

	for ; ; time.Sleep(time.Second) {
		line, err := ReadLine(archiveFile, lineIndex)

		if err != nil {
			//文件读取完，删除归档文件
			if err == io.EOF {
				err1 := os.Remove(fileName)
				if err1 != nil {
					logger.LOG(logger.ERROR, "remove archive file [%s] error: ", archiveFile, err1)
				}

				err2 := os.Remove(archiveIndexname)
				if err2 != nil {
					logger.LOG(logger.ERROR, "remove archive index file [%s] error: ", archiveIndexname, err2)
				}

				logger.LOG(logger.INFO, "remove archive file ", archiveFile, " success.")
			} else {
				logger.LOG(logger.DEBUG, "read line from archive file [%s] error: ", archiveFile, err)
			}
			wg.Done()
			return
		}

		if line != nil {
			putRecordsRequest := new(models.PutRecordsRequest)
			err = utils.ParseJson(line, putRecordsRequest)
			if err != nil {
				wg.Done()
				return
			}

			dis.PutRecords(putRecordsRequest)
			logger.LOG(logger.INFO, "resend put records request success, from file: ", archiveIndexname)

			//发送归档数据成功，标记index
			lineIndex += int64(len(line))
			updateArchiveIndex(archiveIndexname, lineIndex)
		}
	}
}

func ReadLine(filePth string, lineIndex int64) ([]byte, error) {
	var line []byte
	f, err := os.Open(PROPERTY_DATA_CACHE_DIR + string(filepath.Separator) + filePth)
	if err != nil {
		logger.LOG(logger.ERROR, "read data from archive file ", filePth, " error: ", err)
		return line, err
	}
	defer f.Close()

	_, err = f.Seek(lineIndex, 0)
	if err != nil {
		logger.LOG(logger.ERROR, "read data from archive file ", filePth, " error: ", err)
		return line, err
	}

	bfRd := bufio.NewReader(f)

	line, err = bfRd.ReadBytes('\n')
	if err != nil {
		if err == io.EOF {
			return line, err
		}
		logger.LOG(logger.ERROR, "read data from archive file ", filePth, " error: ", err)
		return line, err
	}

	return line, nil
}

func checkArchiveSize() []string {
	archiveFiles := make([]string, 0)
	flist, e := ioutil.ReadDir(PROPERTY_DATA_CACHE_DIR)
	if e != nil {
		logger.LOG(logger.ERROR, "check archive files ", PROPERTY_DATA_CACHE_DIR, " error: ", e)
		return archiveFiles
	}

	for _, f := range flist {
		if !f.IsDir() && strings.HasSuffix(f.Name(), CACHE_ARCHIVE_INDEX_FILE_SUFFIX) {
			archiveDataFilename := strings.Replace(f.Name(), CACHE_ARCHIVE_INDEX_FILE_SUFFIX, CACHE_ARCHIVE_DATA_FILE_SUFFIX, 1);
			archiveFiles = append(archiveFiles, archiveDataFilename)
		}
	}

	return archiveFiles
}

//修改index，标记文件读到的位置，便于恢复
func updateArchiveIndex(archiveIndexFileName string, index int64) {
	logger.LOG(logger.DEBUG, "update archive index file [%s], index: ", archiveIndexFileName, index)

	data := []byte(strconv.FormatInt(index, 10))
	err := ioutil.WriteFile(archiveIndexFileName, data, 0666)
	if err != nil {
		logger.LOG(logger.ERROR, "update archive index file [%s], index: [%d], error: %s", archiveIndexFileName, index, err)
	} else {
		logger.LOG(logger.INFO, "update archive index file [%s], index: [%d], success", archiveIndexFileName, index)
	}
}

//修改index，标记文件读到的位置，便于恢复
func readArchiveIndex(archiveIndexFileName string) int64 {
	logger.LOG(logger.DEBUG, "read archive index file ", archiveIndexFileName)
	var index int64 = 0

	var line []byte
	f, err := os.Open(archiveIndexFileName)
	if err != nil {
		logger.LOG(logger.ERROR, "read archive index file [%s] error: %s", archiveIndexFileName, err)
		return index
	}
	defer f.Close()

	bfRd := bufio.NewReader(f)

	line, _, err = bfRd.ReadLine()
	if err != nil {
		if err != io.EOF {
			logger.LOG(logger.ERROR, "read archive index file [%s] error: %s", archiveIndexFileName, err)
		}
		return index
	}

	index, err = strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		logger.LOG(logger.ERROR, "read archive index file [%s] error: %s", archiveIndexFileName, err)
		return index
	}

	logger.LOG(logger.INFO, "read archive index file [%s], index [%d] success.", archiveIndexFileName, index)
	return index

}