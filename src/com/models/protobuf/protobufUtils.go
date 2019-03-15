package protobuf

import (
	"com/models"
)

/**
*函数原型：
*函数功能：将protobuf类型的上传数据响应，转换为标准的响应类型。
*参数说明：result protobuf类型的上传数据响应
*返回值：上传数据的标准响应
 */
func ToPutRecordsResult(result PutRecordsResult) models.PutRecordsResult {
	var putRecordsResultEntrys = make([]models.PutRecordsResultEntry, 0, 1)

	for _, resultEntry := range result.Records {
		putRecordsResultEntry := models.PutRecordsResultEntry{
		}

		if resultEntry.SequenceNumber != nil {
			putRecordsResultEntry.SequenceNumber = *resultEntry.SequenceNumber
		}

		if resultEntry.ShardId != nil {
			putRecordsResultEntry.PartitionId = *resultEntry.ShardId
		}

		if resultEntry.ErrorCode != nil {
			putRecordsResultEntry.ErrorCode = *resultEntry.ErrorCode
		}

		if resultEntry.ErrorMessage != nil {
			putRecordsResultEntry.ErrorMessage = *resultEntry.ErrorMessage
		}

		putRecordsResultEntrys = append(putRecordsResultEntrys, putRecordsResultEntry)
	}

	putRecordsResult := models.PutRecordsResult{
		FailedRecordCount: int(*result.FailedRecordCount),
		Records: putRecordsResultEntrys,
	}

	return putRecordsResult
}


/**
*函数原型：
*函数功能：将标准请求类型的对象转换为protobuf的请求参数类型
*参数说明：request 标准类型的上传数据请求
*返回值：protobuf类型的上传数据请求
 */
func ToProtobufPutRecordsRequest(request models.PutRecordsRequest) PutRecordsRequest {
	var putRecordsRequestEntrys = []*PutRecordsRequestEntry{};

	for _, requestEntry := range request.Records {
		putRecordsRequestEntry := &PutRecordsRequestEntry{
			Data: requestEntry.Data,
			PartitionKey:&requestEntry.PartitionKey,
			PartitionId:&requestEntry.PartitionId,
			ExplicitHashKey:&requestEntry.ExplicitHashKey,
		}

		putRecordsRequestEntrys = append(putRecordsRequestEntrys, putRecordsRequestEntry)
	}

	protoPutRecordsRequest := PutRecordsRequest{
		StreamName: &request.StreamName,
		Records: putRecordsRequestEntrys,
	}

	return protoPutRecordsRequest
}

/**
*函数原型：
*函数功能：将protobuf类型的下载数据响应，转换为标准的响应类型。
*参数说明：result protobuf类型的下载数据响应
*返回值：下载数据的标准响应
 */
func ToGetRecordsResult(result GetRecordsResult) models.GetRecordsResult {
	var records = make([]models.Record, 0, 1)

	for _, resultEntry := range result.Records {
		record := models.Record{
		}

		if resultEntry.Data != nil {
			record.Data = resultEntry.Data
		}

		if resultEntry.PartitionKey != nil {
			record.PartitionKey = *resultEntry.PartitionKey
		}

		if resultEntry.SequenceNumber != nil {
			record.SequenceNumber = *resultEntry.SequenceNumber
		}

		if resultEntry.Timestamp != nil {
			record.Timestamp = *resultEntry.Timestamp
		}

		if resultEntry.TimestampType != nil {
			record.TimestampType = *resultEntry.TimestampType
		}

		records = append(records, record)
	}

	getRecordsResult := models.GetRecordsResult{
		Records: records,
	}

	if result.NextShardIterator != nil {
		getRecordsResult.NextPartitionCursor = *result.NextShardIterator
	}

	return getRecordsResult
}