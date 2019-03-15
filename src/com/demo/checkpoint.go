package demo

import (
	"com/models"
	"fmt"
)


//提交checkpoint
func CommitCheckpoint(disInstance DISInstance) {
	dis := disInstance.Dis
	streamName := disInstance.StreamName
	appName := disInstance.AppName
	partitionId := disInstance.PartitionId
	sequenceNumber := disInstance.SequenceNumber
	input := &models.CommitCheckpointRequest{
		StreamName: streamName,
		AppName: appName,
		CheckpointType:models.LAST_READ,
		PartitionId: partitionId,
		SequenceNumber:sequenceNumber,

	}
	result, _ := dis.CommitCheckpoint(input)
	if result.Err != nil {
		fmt.Printf("Failed to commitCheckpoint, StreamName [%s], appName [%s], PartitionId [%s], SequenceNumber  [%s].\n",
			streamName, appName, partitionId, sequenceNumber);
	} else {
		fmt.Printf("Success to commitCheckpoint, StreamName [%s], appName [%s], PartitionId [%s], SequenceNumber  [%s].\n",
			streamName, appName, partitionId, sequenceNumber);
	}
}


//获取checkpoint
func GetCheckpoint(disInstance DISInstance) string {
	dis := disInstance.Dis
	streamName := disInstance.StreamName
	appName := disInstance.AppName
	partitionId := disInstance.PartitionId
	input := &models.GetCheckpointRequest{
		StreamName: streamName,
		AppName: appName,
		CheckpointType:models.LAST_READ,
		PartitionId: partitionId,

	}
	result, output := dis.GetCheckpoint(input)
	if result.Err != nil {
		fmt.Printf("Failed to getCheckpoint, StreamName [%s], appName [%s], PartitionId [%s].\n",
			streamName, appName, partitionId);
	} else {
		fmt.Printf("Success to getCheckpoint, StreamName [%s], appName [%s], PartitionId [%s], SequenceNumber  [%s].\n",
			streamName, appName, partitionId, output.SequenceNumber);
		return output.SequenceNumber
	}

	return ""
}

//删除checkpoint
func DeleteCheckpoint(disInstance DISInstance) {
	dis := disInstance.Dis
	streamName := disInstance.StreamName
	appName := disInstance.AppName
	partitionId := disInstance.PartitionId
	input := &models.DeleteCheckpointRequest{
		StreamName: streamName,
		AppName: appName,
		CheckpointType:models.LAST_READ,
		PartitionId: partitionId,
	}
	result, _ := dis.DeleteCheckpoint(input)
	if result.Err != nil {
		fmt.Printf("Failed to deleteCheckpoint, StreamName [%s], appName [%s], PartitionId [%s].\n",
			streamName, appName, partitionId);
	} else {
		fmt.Printf("Success to deleteCheckpoint, StreamName [%s], appName [%s], PartitionId [%s].\n",
			streamName, appName, partitionId);
	}
}
