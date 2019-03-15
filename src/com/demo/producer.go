package demo

import (
	"fmt"
	"com/models"
	"time"
	"strconv"
	"math/rand"
	"bytes"
)

//上传数据
func PutRecords(disInstance DISInstance) (string, string) {
	//fmt.Println("--------------------------PutRecords start--------------------------");

	// 配置上传的数据
	//message := "hello world."
	var buf bytes.Buffer
	message := "b"
	for a := 0; a < 1; a++ {
		buf.WriteString(message)
	}

	var putRecordsRequestEntrys = make([]models.PutRecordsRequestEntry, 0, 1)

	// 上传10条
	rand.Seed(time.Now().Unix())
	for i := 0; i < 10; i++ {
		// PartitionKey为随机值可使数据均匀分布到所有分区中
		putRecordsRequestEntry := models.PutRecordsRequestEntry{
			Data: []byte(message),
			PartitionKey:strconv.Itoa(rand.Intn(10000)),
		}
		putRecordsRequestEntrys = append(putRecordsRequestEntrys, putRecordsRequestEntry)
	}

	input := &models.PutRecordsRequest{
		StreamName: disInstance.StreamName,
		Records: putRecordsRequestEntrys,
	}

	dis := disInstance.Dis
	dis.PutRecords(input)

	var partitionId, sequenceNumber string
	putResult, output := dis.PutRecords(input)
	if putResult.Err == nil {
		fmt.Printf("Put %d [%d successful / %d failed] records.\n",
			len(input.Records),
			len(input.Records) - output.FailedRecordCount,
			output.FailedRecordCount);

		for j := 0; j < len(input.Records); j++ {
			putRecordsResultEntry := output.Records[j];
			if putRecordsResultEntry.ErrorCode != "" {
				// 上传失败
				fmt.Printf("[%s] put failed, errorCode [%s], errorMessage [%s]\n",
					string(putRecordsRequestEntrys[j].Data),
					putRecordsResultEntry.ErrorCode,
					putRecordsResultEntry.ErrorMessage);
			} else
			{
				// 上传成功
				partitionId = putRecordsResultEntry.PartitionId
				sequenceNumber = putRecordsResultEntry.SequenceNumber
				fmt.Printf("[%s] put success, partitionId [%s], partitionKey [%s], sequenceNumber [%s]\n",
					string(putRecordsRequestEntrys[j].Data),
					putRecordsResultEntry.PartitionId,
					putRecordsRequestEntrys[j].PartitionKey,
					putRecordsResultEntry.SequenceNumber);

			}
		}
	}
	fmt.Println("---------------------------PutRecords end---------------------------");
	return partitionId, sequenceNumber
}


//上传数据
func PutRecordsAndCheckpoint(disInstance DISInstance) {
	fmt.Println("--------------------PutRecordsAndCheckpoint start-------------------");

	partitionId, sequenceNumber := PutRecords(disInstance)        //上传数据

	CreateApp(disInstance) //创建APP

	disInstance.PartitionId = partitionId
	disInstance.SequenceNumber = sequenceNumber
	CommitCheckpoint(disInstance) //提交checkpoint

	fmt.Println("---------------------PutRecordsAndCheckpoint end--------------------");
}