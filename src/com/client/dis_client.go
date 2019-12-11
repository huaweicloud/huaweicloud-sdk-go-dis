/**
* Copyright 2015 Huawei Technologies Co., Ltd. All rights reserved.
* eSDK is licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package client

import (
	"com/logger"
	"com/models"
	"com/models/protobuf"
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"strconv"
	"strings"
	"com/utils"
	"sync"
	"com/utils/cache"
)

var recordsRetryLock sync.Mutex

/**
*函数原型：func (client *Client) CreateStream(input *models.CreateStreamRequest) (*models.Result, *models.CreateStreamResult)
*函数功能：创建通道
*参数说明：
*返回值：result: Result对象实例
 */
func (client *Client) CreateStream(input *models.CreateStreamRequest) (*models.Result, *models.CreateStreamResult) {
	logger.LOG(logger.INFO, "enter CreateStream...")
	dis := getRequest(POST, client.projectId)
	setDisURI("streams", "", dis)
	if (input.StreamType == "") {
		input.StreamType = models.STREAM_TYPE_COMMON
	}
	if (input.DataType == "") {
		input.DataType = models.DATA_TYPE_BLOB
	}
	if (input.DataDuration == 0) {
		input.DataDuration = models.DATA_DURATION_DEFAULT
	}

	inputbytes, err := json.Marshal(input)
	if nil != err {
		logger.LOG(logger.ERROR, "CreateStream parameter error ", err)
	}

	util, r := client.connectDisWithXml(dis, string(inputbytes))
	if r == nil {
		output := new(models.CreateStreamResult)
		if ret := client.getResponseWithOutput(util, output); ret.Err != nil {
			return ret, nil
		} else {
			return ret, output
		}
	}
	return r, nil
}

/**
*函数原型：func (client *Client) DescStream(input *models.DescribeStreamRequest) (*models.Result, *models.DescribeStreamResult)
*函数功能：描述通道详情
*参数说明：
*返回值：result: Result对象实例
 */
func (client *Client) DescStream(input *models.DescribeStreamRequest) (*models.Result, *models.DescribeStreamResult) {
	logger.LOG(logger.INFO, "enter DescStream...")
	dis := getRequest(GET, client.projectId)
	setDisURI("streams", input.StreamName, dis)
	if (input.LimitPartitions == 0) {
		input.LimitPartitions = 100
	}
	setDisPath("limit_partitions", strconv.Itoa(input.LimitPartitions), dis, false)
	setDisPath("start_partitionId", input.StartPartitionId, dis, false)

	util, r := client.connectDis(dis)
	if r == nil {
		output := new(models.DescribeStreamResult)
		if ret := client.getResponseWithOutput(util, output); ret.Err != nil {
			return ret, nil
		} else {
			return ret, output
		}
	}
	return r, nil
}


/**
*函数原型：func (client *Client) UpdatePartitionCount(input *models.UpdatePartitionCountRequest)
*函数功能：扩缩容通道
*参数说明：
*返回值：result: Result对象实例
 */
func (client *Client) UpdatePartitionCount(input *models.UpdatePartitionCountRequest) (*models.Result, *models.UpdatePartitionCountResult) {
	logger.LOG(logger.INFO, "enter UpdatePartitionCount...")
	dis := getRequest(PUT, client.projectId)
	setDisURI("streams", input.StreamName, dis)
	inputbytes, err := json.Marshal(input)
	if nil != err {
		logger.LOG(logger.ERROR, "UpdatePartitionCount parameter error ", err)
	}

	util, r := client.connectDisWithXml(dis, string(inputbytes))
	if r == nil {
		output := new(models.UpdatePartitionCountResult)
		if ret := client.getResponseWithOutput(util, output); ret.Err != nil {
			return ret, nil
		} else {
			return ret, output
		}
	}
	return r, nil
}

/**
*函数原型：func (client *Client) ListStreams(input *models.ListStreamsRequest)
*函数功能：获取流列表
*参数说明：
*返回值：result: Result对象实例
 */
func (client *Client) ListStreams(input *models.ListStreamsRequest) (*models.Result, *models.ListStreamsResult) {
	logger.LOG(logger.INFO, "enter ListStreams...")
	dis := getRequest(GET, client.projectId)
	setDisURI("streams", "", dis)
	if (input.Limit == 0) {
		input.Limit = 10
	}
	setDisPath("limit", strconv.Itoa(input.Limit), dis, false)
	setDisPath("start_stream_name", input.ExclusivetartStreamName, dis, false)

	util, r := client.connectDis(dis)
	if r == nil {
		output := new(models.ListStreamsResult)
		if ret := client.getResponseWithOutput(util, output); ret.Err != nil {
			return ret, nil
		} else {
			return ret, output
		}
	}
	return r, nil
}

/**
*函数原型：func (client *Client)  DeleteStream(input *models.DeleteStreamRequest)
*函数功能：删除通道
*参数说明：
*返回值：result: Result对象实例
 */
func (client *Client) DeleteStream(input *models.DeleteStreamRequest) (*models.Result, *models.DeleteStreamResult) {
	logger.LOG(logger.INFO, "enter DeleteStream...")
	dis := getRequest(DELETE, client.projectId)
	setDisURI("streams", input.StreamName, dis)
	util, r := client.connectDis(dis)
	if r == nil {
		output := new(models.DeleteStreamResult)
		if ret := client.getResponseWithOutput(util, output); ret.Err != nil {
			return ret, nil
		} else {
			return ret, output
		}
	}
	return r, nil
}

/**
*函数原型：func (client *Client) CreateApp(appName string) *models.Result
*函数功能：创建APP
*参数说明：
*返回值：result: Result对象实例
 */
func (client *Client) CreateApp(appName string) *models.Result {
	logger.LOG(logger.INFO, "enter CreateApp...")
	dis := getRequest(POST, client.projectId)
	setDisURI("apps", "", dis)
	input := models.CreateAppRequest{
		AppName:appName,
	};

	inputbytes, err := json.Marshal(input)
	if nil != err {
		logger.LOG(logger.ERROR, "CreateApp parameter error ", err)
	}

	util, r := client.connectDisWithXml(dis, string(inputbytes))
	if r == nil {
		ret, _ := client.getResponseWithOutputAndResposeByte(util)
		return ret
	}
	return r
}

/**
*函数原型：func (client *Client) DescApp(appName string) (*models.Result, *models.DescribeAppResult)
*函数功能：描述APP详情
*参数说明：
*返回值：result: Result对象实例
 */
func (client *Client) DescApp(appName string) (*models.Result, *models.DescribeAppResult) {
	logger.LOG(logger.INFO, "enter DescApp...")
	dis := getRequest(GET, client.projectId)
	setDisURI("apps", appName, dis)
	util, r := client.connectDis(dis)
	if r == nil {
		output := new(models.DescribeAppResult)
		if ret := client.getResponseWithOutput(util, output); ret.Err != nil {
			return ret, nil
		} else {
			return ret, output
		}
	}
	return r, nil
}

/**
*函数原型：func (client *Client) ListApps(input *models.ListAppsRequest) (*models.Result, *models.DescribeAppResult)
*函数功能：查询APP列表
*参数说明：
*返回值：result: Result对象实例
 */
func (client *Client) ListApps(input *models.ListAppsRequest) (*models.Result, *models.ListAppsResult) {
	logger.LOG(logger.INFO, "enter ListApps...")
	dis := getRequest(GET, client.projectId)
	setDisURI("apps", "", dis)
	if (input.Limit == 0) {
		input.Limit = 10
	}
	setDisPath("limit", strconv.Itoa(input.Limit), dis, false)
	setDisPath("start_app_name", input.ExclusiveStartAppName, dis, false)
	util, r := client.connectDis(dis)
	if r == nil {
		output := new(models.ListAppsResult)
		if ret := client.getResponseWithOutput(util, output); ret.Err != nil {
			return ret, nil
		} else {
			return ret, output
		}
	}
	return r, nil
}

/**
*函数原型：func (client *Client) DeleteApp(appName string) *models.Result
*函数功能：删除APP
*参数说明：
*返回值：result: Result对象实例
 */
func (client *Client) DeleteApp(appName string) *models.Result {
	logger.LOG(logger.INFO, "enter DeleteApp...")
	dis := getRequest(DELETE, client.projectId)
	setDisURI("apps", appName, dis)
	util, r := client.connectDis(dis)
	if r == nil {
		ret, _ := client.getResponseWithOutputAndResposeByte(util)
		return ret
	}
	return r
}

/**
*函数原型：func (client *Client) CreateTransferTask(input *models.CreateTransferTaskRequest) (*models.Result, *models.CreateTransferTaskResult)
*函数功能：创建转储任务
*参数说明：
*返回值：result: Result对象实例
 */
func (client *Client) CreateTransferTask(input *models.CreateTransferTaskRequest) (*models.Result, *models.CreateTransferTaskResult) {
	logger.LOG(logger.INFO, "enter CreateTransferTask...")
	dis := getRequest(POST, client.projectId)
	setDisURI("streams", input.StreamName, dis)
	setDisURI("transfer-tasks", "", dis)
	inputbytes, err := json.Marshal(input)
	if nil != err {
		logger.LOG(logger.ERROR, "CreateTransferTask parameter error ", err)
	}

	util, r := client.connectDisWithXml(dis, string(inputbytes))
	if r == nil {
		output := new(models.CreateTransferTaskResult)
		if ret := client.getResponseWithOutput(util, output); ret.Err != nil {
			return ret, nil
		} else {
			return ret, output
		}
	}
	return r, nil
}

/**
*函数原型：func (client *Client) ListTransferTasks(input *models.ListTransferTasksRquest)
*函数功能：查询通道下的转储任务列表
*参数说明：
*返回值：result: Result对象实例
 */
func (client *Client) ListTransferTasks(input *models.ListTransferTasksRquest) (*models.Result, *models.ListTransferTasksResult) {
	logger.LOG(logger.INFO, "enter ListTransferTasks...")
	dis := getRequest(GET, client.projectId)
	setDisURI("streams", input.StreamName, dis)
	setDisURI("transfer-tasks", "", dis)

	if (input.Limit == 0) {
		input.Limit = 0
	}

	setDisPath("start_task_name", input.ExclusiveStartTaskName, dis, true)
	setDisPath("limit", strconv.Itoa(input.Limit), dis, false)

	util, r := client.connectDis(dis)
	if r == nil {
		output := new(models.ListTransferTasksResult)
		if ret := client.getResponseWithOutput(util, output); ret.Err != nil {
			return ret, nil
		} else {
			return ret, output
		}
	}
	return r, nil
}

/**
*函数原型：func (client *Client) DescTransferTask(input *models.DescribeTransferTaskRequest)
*函数功能：查询通道下的转储任务详情
*参数说明：
*返回值：result: Result对象实例
 */
func (client *Client) DescTransferTask(input *models.DescribeTransferTaskRequest) (*models.Result, *models.DescribeTransferTaskResult) {
	logger.LOG(logger.INFO, "enter DescTransferTask...")
	dis := getRequest(GET, client.projectId)
	setDisURI("streams", input.StreamName, dis)
	setDisURI("transfer-tasks", input.TaskName, dis)

	util, r := client.connectDis(dis)
	if r == nil {
		output := new(models.DescribeTransferTaskResult)
		if ret := client.getResponseWithOutput(util, output); ret.Err != nil {
			return ret, nil
		} else {
			return ret, output
		}
	}
	return r, nil
}

/**
*函数原型：func (client *Client) DeleteTransferTask(input *models.DeleteTransferTaskRequest)
*函数功能：获取checkpoint
*参数说明：
*返回值：result: Result对象实例
 */
func (client *Client) DeleteTransferTask(input *models.DeleteTransferTaskRequest) (*models.Result, *models.DeleteTransferTaskResult) {
	logger.LOG(logger.INFO, "enter DeleteTransferTask...")
	dis := getRequest(DELETE, client.projectId)
	setDisURI("streams", input.StreamName, dis)
	setDisURI("transfer-tasks", input.TaskName, dis)

	util, r := client.connectDis(dis)
	if r == nil {
		output := new(models.DeleteTransferTaskResult)
		if ret := client.getResponseWithOutput(util, output); ret.Err != nil {
			return ret, nil
		} else {
			return ret, output
		}
	}
	return r, nil
}


/**
*函数原型：func (client *Client) CommitCheckpoint(input *models.CommitCheckpointRequest) (*models.Result, *models.CommitCheckpointResult)
*函数功能：提交checkpoint
*参数说明：
*返回值：result: Result对象实例
 */
func (client *Client) CommitCheckpoint(input *models.CommitCheckpointRequest) (*models.Result, *models.CommitCheckpointResult) {
	logger.LOG(logger.INFO, "enter CommitCheckpoint...")
	dis := getRequest(POST, client.projectId)
	setDisURI("checkpoints", "", dis)
	inputbytes, err := json.Marshal(input)
	if nil != err {
		logger.LOG(logger.ERROR, "CommitCheckpointRequest parameter error ", err)
	}

	util, r := client.connectDisWithXml(dis, string(inputbytes))
	if r == nil {
		output := new(models.CommitCheckpointResult)
		if ret := client.getResponseWithOutput(util, output); ret.Err != nil {
			return ret, nil
		} else {
			return ret, output
		}
	}
	return r, nil
}

/**
*函数原型：func (client *Client) GetCheckpoint(input *models.GetCheckpointRequest) (*models.Result, *models.GetCheckpointResult)
*函数功能：获取checkpoint
*参数说明：
*返回值：result: Result对象实例
 */
func (client *Client) GetCheckpoint(input *models.GetCheckpointRequest) (*models.Result, *models.GetCheckpointResult) {
	logger.LOG(logger.INFO, "enter GetCheckpoint...")
	dis := getRequest(GET, client.projectId)
	setDisURI("checkpoints", "", dis)
	setDisPath("stream_name", input.StreamName, dis, true)
	setDisPath("partition_id", input.PartitionId, dis, true)
	setDisPath("app_name", input.AppName, dis, true)
	setDisPath("checkpoint_type", input.CheckpointType, dis, true)
	setDisPath("timestamp", strconv.FormatInt(input.Timestamp, 10), dis, false)

	util, r := client.connectDis(dis)
	if r == nil {
		output := new(models.GetCheckpointResult)
		if ret := client.getResponseWithOutput(util, output); ret.Err != nil {
			return ret, nil
		} else {
			return ret, output
		}
	}
	return r, nil
}

/**
*函数原型：func (client *Client) DeleteCheckpointResult(input *models.DeleteCheckpointRequest)
*函数功能：获取checkpoint
*参数说明：
*返回值：result: Result对象实例
 */
func (client *Client) DeleteCheckpoint(input *models.DeleteCheckpointRequest) (*models.Result, *models.DeleteCheckpointResult) {
	logger.LOG(logger.INFO, "enter DeleteCheckpoint...")
	dis := getRequest(DELETE, client.projectId)
	setDisURI("checkpoints", "", dis)
	setDisPath("stream_name", input.StreamName, dis, false)
	setDisPath("partition_id", input.PartitionId, dis, false)
	setDisPath("app_name", input.AppName, dis, false)
	setDisPath("checkpoint_type", input.CheckpointType, dis, false)

	util, r := client.connectDis(dis)
	if r == nil {
		output := new(models.DeleteCheckpointResult)
		if ret := client.getResponseWithOutput(util, output); ret.Err != nil {
			return ret, nil
		} else {
			return ret, output
		}
	}
	return r, nil
}

/**
*函数原型：func (client *Client) GetPartitionCursor(input *models.GetPartitionCursorRequest)  (*models.Result, *models.GetPartitionCursorResult)
*函数功能：用户获取迭代器，根据迭代器获取一次数据和下一个迭代器。
*参数说明：GetPartitionCursorRequest 用户获取迭代器的请求参数
*返回值：获取迭代器的响应结果
 */
func (client *Client) GetPartitionCursor(input *models.GetPartitionCursorRequest) (*models.Result, *models.GetPartitionCursorResult) {
	logger.LOG(logger.INFO, "enter GetPartitionCursor...")
	dis := getRequest(GET, client.projectId)
	setDisURI("cursors", "", dis)
	setDisPath("stream-name", input.StreamName, dis, true)
	setDisPath("partition-id", input.PartitionId, dis, true)
	setDisPath("cursor-type", input.CursorType, dis, false)
	setDisPath("starting-sequence-number", input.StartingSequenceNumber, dis, false)
	setDisPath("timestamp", strconv.FormatInt(input.Timestamp, 10), dis, false)
	util, r := client.connectDis(dis)
	if r == nil {
		output := new(models.GetPartitionCursorResult)
		if ret := client.getResponseWithOutput(util, output); ret.Err != nil {
			return ret, nil
		} else {
			return ret, output
		}
	}
	return r, nil
}

/**
*函数原型：func (client *Client) GetRecords(input *models.GetRecordsRequest)  (*models.Result, *models.GetRecordsResult)
*函数功能：从DIS实例中下载数据。
*参数说明：GetRecordsRequest 用户获取数据的请求参数
*返回值：获取数据的响应结果
 */
func (client *Client) GetRecords(input *models.GetRecordsRequest) (*models.Result, *models.GetRecordsResult) {
	logger.LOG(logger.INFO, "enter GetRecordsRequest...")
	dis := getRequest(GET, client.projectId)
	setDisURI("records", "", dis)
	setDisPath("partition-cursor", input.PartitionCursor, dis, true)

	//protobuf
	if client.bodySerializeType == models.PROTOBUF {
		util, r := client.connectDis(dis)
		util.SetHeader("Content-Type", "application/x-protobuf; charset=utf-8")
		if r == nil {
			if ret, resposeByte := client.getResponseWithOutputAndResposeByte(util); ret.Err != nil {
				return ret, nil
			} else {
				protoResult := &protobuf.GetRecordsResult{}
				err := proto.Unmarshal(resposeByte, protoResult)
				if err != nil {
					logger.LOG(logger.ERROR, err.Error())
				}

				output := protobuf.ToGetRecordsResult(*protoResult)
				return ret, &output
			}
		}
		return r, nil
	} else {
		util, r := client.connectDis(dis)
		if r == nil {
			output := new(models.GetRecordsResult)
			if ret := client.getResponseWithOutput(util, output); ret.Err != nil {
				return ret, nil
			} else {
				return ret, output
			}
		}
		return r, nil
	}

}

/**
*函数原型：func (client *Client) PutRecords(input *models.PutRecordsRequest)  (*models.Result, *models.PutRecordsResult)
*函数功能：从DIS实例中下载数据。
*参数说明：GetRecordsRequest 用户上传数据的请求参数
*返回值：上传数据的响应结果
 */
func (client *Client) PutRecords(input *models.PutRecordsRequest) (*models.Result, *models.PutRecordsResult) {
	return client.innerPutRecordsSupportingCache(input)
}

func (client *Client) innerPutRecordsSupportingCache(input *models.PutRecordsRequest) (*models.Result, *models.PutRecordsResult) {
	putRecordsReq := &models.PutRecordsRequest{
		StreamId: input.StreamId,
		StreamName: input.StreamName,
		Records: input.Records,
	}
	if client.cacheResendConf.DataCacheEnable {
		result, putRecordsResult := client.innerPutRecordsWithRetry(putRecordsReq)
		if result.Err != nil {
			cache.PutToCache(putRecordsReq, client.cacheResendConf)
			return result, putRecordsResult
		} else if putRecordsResult.FailedRecordCount > 0 {
			// 过滤出上传失败的记录
			putRecordsResultEntries := putRecordsResult.Records
			failedPutRecordsRequestEntries := make([]models.PutRecordsRequestEntry, 0)
			index := 0
			for _, putRecordsResultEntry := range putRecordsResultEntries {
				if len(putRecordsResultEntry.ErrorCode) != 0 {
					failedPutRecordsRequestEntries = append(failedPutRecordsRequestEntries, putRecordsReq.Records[index])
				}
				index ++
			}
			putRecordsReq.Records = failedPutRecordsRequestEntries

			logger.LOG(logger.WARNING, "Local data cache is enabled, try to put failed records to local.")
			cache.PutToCache(putRecordsReq, client.cacheResendConf)
			return result, putRecordsResult
		}
		return result, putRecordsResult
	} else {
		return client.innerPutRecordsWithRetry(putRecordsReq)
	}
}

func (client *Client) innerPutRecordsWithRetry(input *models.PutRecordsRequest) (*models.Result, *models.PutRecordsResult) {
	var putRecordsResultEntryList []models.PutRecordsResultEntry
	var putRecordsResult *models.PutRecordsResult
	var result *models.Result
	retryPutRecordsRequest := &models.PutRecordsRequest{
		StreamId: input.StreamId,
		StreamName: input.StreamName,
		Records: input.Records,
	}

	retryCount := -1
	currentFailed := 0
	noRetryRecordsCount := 0
	var retryIndex []int
	var backOff *utils.ExponentialBackOff
	for {
		retryCount ++
		if retryCount > 0 {
			// 等待一段时间再发起重试
			if backOff == nil {
				logger.LOG(logger.INFO, "Put records retry lock.")
				recordsRetryLock.Lock()
				backOff = utils.NewExponentialBackOff(utils.DEFAULT_INITIAL_INTERVAL, utils.DEFAULT_MULTIPLIER, utils.DEFAULT_MAX_INTERVAL, utils.DEFAULT_MAX_ELAPSED_TIME)
			}

			if putRecordsResult != nil && currentFailed != len(putRecordsResult.Records) {
				// 部分失败则重置退避时间
				backOff.ResetCurrentInterval()
			}

			sleepMs := backOff.GetNextBackOff()
			if len(retryPutRecordsRequest.Records) > 0 {
				logger.LOG(logger.INFO, "Put records but failed, will re-try after backoff")
			}

			backOff.BackOff(sleepMs)
		}

		result, putRecordsResult = client.innerPutRecords(retryPutRecordsRequest)
		if result.Err != nil {
			if (putRecordsResultEntryList != nil) {
				logger.LOG(logger.ERROR, result.ErrResponse.Message)
				break;
			}
			return result, nil
		}

		if putRecordsResult != nil {
			currentFailed = putRecordsResult.FailedRecordCount

			// 第一次发送全部成功或者不需要重试，则直接返回结果
			if putRecordsResultEntryList == nil && currentFailed == 0 || models.RECORDS_RETRIES == 0 {
				return result, putRecordsResult;
			}

			// 存在发送失败的情况，需要重试，则使用数组来汇总每次请求后的结果
			if putRecordsResultEntryList == nil {
				putRecordsResultEntryList = make([]models.PutRecordsResultEntry, len(input.Records))
			}

			// 需要重试发送数据的原始下标
			retryIndexTemp := make([]int, 0)

			if currentFailed > 0 {
				// 初始化重试发送的数据请求
				retryPutRecordsRequest.Records = make([]models.PutRecordsRequestEntry, 0)
			}

			// 对每条结果分析，更新结果数据
			for i := 0; i < len(putRecordsResult.Records); i++ {
				// 获取重试数据在原始数据中的下标位置
				var originalIndex int
				if retryIndex == nil {
					originalIndex = i
				} else {
					originalIndex = retryIndex[i]
				}

				putRecordsResultEntry := &putRecordsResult.Records[i]

				if len(putRecordsResultEntry.ErrorCode) != 0 {
					if isRecordsRetriableErrorCode(putRecordsResultEntry.ErrorCode) {
						// 只对指定异常(如流控与服务端内核异常)进行重试
						retryIndexTemp = append(retryIndexTemp, originalIndex)
						retryPutRecordsRequest.Records = append(retryPutRecordsRequest.Records, input.Records[i])
					} else {
						noRetryRecordsCount++
					}
				}
				putRecordsResultEntryList[originalIndex] = *putRecordsResultEntry
			}

			if len(retryIndexTemp) > 0 {
				retryIndex = retryIndexTemp
			} else {
				retryIndex = make([]int, 0)
			}
		}

		//超过重试次数或失败数据都已重新发送成功，则中断重试
		if (retryIndex != nil && len(retryIndex) == 0) || retryCount >= models.RECORDS_RETRIES {
			break
		}
	}

	if retryCount > 0 {
		recordsRetryLock.Unlock()
		logger.LOG(logger.INFO, "Put records retry unlock.")
	}

	if retryIndex == nil {
		// 不可能存在此情况，完全没有发送出去会直接抛出异常
		putRecordsResult.FailedRecordCount = len(input.Records);
	} else {
		putRecordsResult.FailedRecordCount = len(retryIndex) + noRetryRecordsCount
		putRecordsResult.Records = putRecordsResultEntryList
	}

	return result, putRecordsResult
}

func (client *Client) innerPutRecords(input *models.PutRecordsRequest) (*models.Result, *models.PutRecordsResult) {
	dis := getRequest(POST, client.projectId)
	setDisURI("records", "", dis)
	//protobuf
	if client.bodySerializeType == models.PROTOBUF {
		protoRequest := protobuf.ToProtobufPutRecordsRequest(*input)
		inputbytes, err := proto.Marshal(&protoRequest)
		if nil != err {
			logger.LOG(logger.ERROR, "PutRecordsRequest parameter error ", err)
		}

		util, r := client.connectDisWithXml(dis, string(inputbytes))
		util.SetHeader("Content-Type", "application/x-protobuf; charset=utf-8")
		if r == nil {
			if ret, resposeByte := client.getResponseWithOutputAndResposeByte(util); ret.Err != nil {
				return ret, nil
			} else {
				protoResult := &protobuf.PutRecordsResult{}
				err := proto.Unmarshal(resposeByte, protoResult)
				if err != nil {
					logger.LOG(logger.ERROR, err.Error())
				}

				output := protobuf.ToPutRecordsResult(*protoResult)
				return ret, &output
			}
		}
		return r, nil
	} else {
		inputbytes, err := json.Marshal(input)
		if nil != err {
			logger.LOG(logger.ERROR, "PutRecordsRequest parameter error ", err)
		}

		util, r := client.connectDisWithXml(dis, string(inputbytes))
		if r == nil {
			output := new(models.PutRecordsResult)
			if ret := client.getResponseWithOutput(util, output); ret.Err != nil {
				return ret, nil
			} else {
				return ret, output
			}
		}
		return r, nil
	}
}

func isRecordsRetriableErrorCode(errorCode string) bool {
	producerRecordsRetriableErrorCode := []string{"DIS.4303", "DIS.5"}
	for _, retryCode := range producerRecordsRetriableErrorCode {
		if strings.Contains(errorCode, retryCode) {
			return true
		}
	}
	return false
}

