package demo

import (
	"fmt"
	"com/models"
	"time"
)


//读取数据
func GetRecords(disInstance DISInstance) {
	fmt.Println("-------------------------GetRecords start---------------------------");
	//读取迭代器
	partitionCursor := GetPartitionCursor(disInstance)
	dis := disInstance.Dis
	for i := 0; i < 2; i++ {
		input := &models.GetRecordsRequest{PartitionCursor: partitionCursor}
		_, output := dis.GetRecords(input)
		if output != nil {
			for _, record := range output.Records {
				fmt.Printf("Get record [%s], partitionKey [%s], sequenceNumber [%s], Timestamp [%d], TimestampType [%s].\n",
					string(record.Data),
					record.PartitionKey,
					record.SequenceNumber,
					record.Timestamp, record.TimestampType);
			}

			if len(output.Records) == 0 {
				fmt.Printf("No record!\n");
				time.Sleep(1000 * time.Millisecond)
			}

			partitionCursor = output.NextPartitionCursor
		}
	}
	fmt.Println("--------------------------GetRecords end----------------------------");
}

//获取迭代器:返回partitionCursor，用于读取数据
func GetPartitionCursor(disInstance DISInstance) string {
	dis := disInstance.Dis
	streamName := disInstance.StreamName
	partitionId := disInstance.PartitionId
	sequenceNumber := disInstance.SequenceNumber
	cursorType := disInstance.CursorType
	timestamp := disInstance.Timestamp
	input := &models.GetPartitionCursorRequest{
		StreamName: streamName,
		PartitionId: partitionId,
		StartingSequenceNumber:sequenceNumber,
		CursorType:cursorType,
		Timestamp:timestamp,
	}
	result, output := dis.GetPartitionCursor(input)
	if result.Err == nil {
		fmt.Printf("Get partition cursor [%s], partitionId [%s], cursor [%s] success.\n", streamName, input.PartitionId, output.PartitionCursor);
		return output.PartitionCursor
	}
	return ""
}

//上传数据
func GetRecordsWithCheckpoint(disInstance DISInstance) {
	fmt.Println("-------------------GetRecordsWithCheckpoint start-------------------");

	sequenceNumber := GetCheckpoint(disInstance) //获取checkpoint

	disInstance.CursorType = models.AT_SEQUENCE_NUMBER
	disInstance.SequenceNumber = sequenceNumber
	GetRecords(disInstance) //上传数据

	DeleteCheckpoint(disInstance) //删除checkpoint

	DeleteApp(disInstance) //删除app

	fmt.Println("--------------------GetRecordsWithCheckpoint end--------------------");
}
