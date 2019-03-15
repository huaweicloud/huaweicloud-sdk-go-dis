package demo

import (
	"fmt"
	"com/models"
	"encoding/json"
	"com/logger"
)

func TransferTaskDemo(disInstance DISInstance) {
	fmt.Println("------------------------TransferTaskDemo start-------------------------");
	CreateOBSTransferTask(disInstance)
	DescTransferTask(disInstance)
	ListTransferTasks(disInstance)
	DeleteTransferTask(disInstance)
	fmt.Println("-------------------------TransferTaskDemo end--------------------------");
}

/**
*函数功能：创建目标服务为OBS的转储任务，目标文件格式为text
*说明：通道必须提前创建“源数据schema”
 */
func CreateOBSTransferTask(disInstance DISInstance) {
	dis := disInstance.Dis

	descriptor := &models.OBSDestinationDescriptorRequest{
		TransferTaskName: "feihang",
		AgencyName:"all",
		ConsumerStrategy:models.TRIM_HORIZON,
		DeliverTimeInterval:90,
		DestinationFileType:models.TEXT,
		ObsBucketPath:"dis-tests-not-delete",
		RecordDelimiter:",",
		FilePrefix:"feihang",
	}

	input := &models.CreateTransferTaskRequest{
		StreamName    :disInstance.StreamName,
		DestinationType:models.DESTINATION_OBS,
		ObsDestinationDescriptor:descriptor,
	}
	result, _ := dis.CreateTransferTask(input)
	if result.Err != nil {
		fmt.Printf("Failed to create TransferTask [%s].\n", descriptor.TransferTaskName);
	} else {
		fmt.Printf("Success to create TransferTask [%s].\n", descriptor.TransferTaskName);
	}
}

/**
*函数功能：创建目标服务为OBS的转储任务，目标文件格式为Parquet
*说明：通道必须提前创建“源数据schema”
 */
func CreateOBSParquetTransferTask(disInstance DISInstance) {
	dis := disInstance.Dis

	descriptor := &models.OBSDestinationDescriptorRequest{
		TransferTaskName: "feihang",
		AgencyName:"all",
		ConsumerStrategy:models.TRIM_HORIZON,
		DeliverTimeInterval:90,
		DestinationFileType:models.PARQUET,
		ObsBucketPath:"dis-tests-not-delete",
		FilePrefix:"feihang",
	}

	//如果需要自定义时间戳目录，则需要配置processingSchema
	processingSchema := &models.ProcessingSchema{
		TimestampName:"field_0", //源数据时间戳字段的属性名称
		TimestampType: models.TIMESTAMP_TYPE_STRING,
		TimestampFormat: models.TIMESTAMP_FORMAT_0,
	}
	descriptor.ProcessingSchema = processingSchema

	input := &models.CreateTransferTaskRequest{
		StreamName    :disInstance.StreamName,
		DestinationType:models.DESTINATION_OBS,
		ObsDestinationDescriptor:descriptor,
	}
	result, _ := dis.CreateTransferTask(input)
	if result.Err != nil {
		fmt.Printf("Failed to create TransferTask [%s].\n", descriptor.TransferTaskName);
	} else {
		fmt.Printf("Success to create TransferTask [%s].\n", descriptor.TransferTaskName);
	}
}

/**
*函数功能：创建目标服务为OBS的转储任务，目标文件格式为Carbon
*说明：通道必须提前创建“源数据schema”
 */
func CreateOBSCarbonTransferTask(disInstance DISInstance) {
	dis := disInstance.Dis

	descriptor := &models.OBSDestinationDescriptorRequest{
		TransferTaskName: "feihang",
		AgencyName:"all",
		ConsumerStrategy:models.TRIM_HORIZON,
		DeliverTimeInterval:90,
		DestinationFileType:models.CARBON,
		ObsBucketPath:"dis-tests-not-delete",
		FilePrefix:"feihang",
	}

	//如果需要配置carbon的检索属性，则需要配置carbonProperties
	carbonProperties := make(map[string]string)
	carbonProperties["sortcolumns"] = "field_0"
	descriptor.CarbonProperties = carbonProperties

	input := &models.CreateTransferTaskRequest{
		StreamName    :disInstance.StreamName,
		DestinationType:models.DESTINATION_OBS,
		ObsDestinationDescriptor:descriptor,
	}
	result, _ := dis.CreateTransferTask(input)
	if result.Err != nil {
		fmt.Printf("Failed to create TransferTask [%s].\n", descriptor.TransferTaskName);
	} else {
		fmt.Printf("Success to create TransferTask [%s].\n", descriptor.TransferTaskName);
	}
}

/**
*函数功能：创建目标服务为MRS的转储任务，目标文件格式为text
 */
func CreateMRSTransferTask(disInstance DISInstance) {
	dis := disInstance.Dis

	descriptor := &models.MRSDestinationDescriptorRequest{
		TransferTaskName: "feihang",
		AgencyName:"all",
		ConsumerStrategy:models.TRIM_HORIZON,
		DeliverTimeInterval:90,
		DestinationFileType:models.TEXT,
		ObsBucketPath:"dis-tests-not-delete",
		MrsClusterId:"xx",
		MrsClusterName:"xx",
		HdfsPrefixFolder:"/tmp",
	}

	input := &models.CreateTransferTaskRequest{
		StreamName    :disInstance.StreamName,
		DestinationType:models.DESTINATION_MRS,
		MrsDestinationDescriptor:descriptor,
	}
	result, _ := dis.CreateTransferTask(input)
	if result.Err != nil {
		fmt.Printf("Failed to create TransferTask [%s].\n", descriptor.TransferTaskName);
	} else {
		fmt.Printf("Success to create TransferTask [%s].\n", descriptor.TransferTaskName);
	}
}

/**
*函数功能：创建目标服务为MRS的转储任务，目标文件格式为text
 */
func CreateDLITransferTask(disInstance DISInstance) {
	dis := disInstance.Dis

	descriptor := &models.DLIDestinationDescriptorRequest{
		TransferTaskName: "feihang",
		AgencyName:"all",
		ConsumerStrategy:models.TRIM_HORIZON,
		DeliverTimeInterval:90,
		ObsBucketPath:"dis-tests-not-delete",
		DliDatabaseName:"A",
		DliTableName:"b",
	}

	input := &models.CreateTransferTaskRequest{
		StreamName    :disInstance.StreamName,
		DestinationType:models.DESTINATION_DLI,
		DliDestinationDescriptor:descriptor,
	}
	result, _ := dis.CreateTransferTask(input)
	if result.Err != nil {
		fmt.Printf("Failed to create TransferTask [%s].\n", descriptor.TransferTaskName);
	} else {
		fmt.Printf("Success to create TransferTask [%s].\n", descriptor.TransferTaskName);
	}
}

/**
*函数功能：创建目标服务为MRS的转储任务，目标文件格式为text
 */
func CreateDWSTransferTask(disInstance DISInstance) {
	dis := disInstance.Dis

	descriptor := &models.DWSDestinationDescriptorRequest{
		TransferTaskName: "feihang",
		AgencyName:"all",
		ConsumerStrategy:models.TRIM_HORIZON,
		DeliverTimeInterval:90,
		ObsBucketPath:"dis-tests-not-delete",
		DwsClusterId:"a",
		DwsClusterName:"b",
		DwsSchema:"c",
		DwsDelimiter:",",
	}

	input := &models.CreateTransferTaskRequest{
		StreamName    :disInstance.StreamName,
		DestinationType:models.DESTINATION_DWS,
		DwsDestinationDescriptor:descriptor,
	}
	result, _ := dis.CreateTransferTask(input)
	if result.Err != nil {
		fmt.Printf("Failed to create TransferTask [%s].\n", descriptor.TransferTaskName);
	} else {
		fmt.Printf("Success to create TransferTask [%s].\n", descriptor.TransferTaskName);
	}
}

/**
*函数功能：创建目标服务为Cloudtable HBase的转储任务
 */
func CreateCloudtablHBaseTransferTask(disInstance DISInstance) {
	dis := disInstance.Dis

	descriptor := &models.CloudtableDestinationDescriptorRequest{
		TransferTaskName: "feihang",
		AgencyName:"all",
		ConsumerStrategy:models.TRIM_HORIZON,
		CloudtableClusterId:"a",
		CloudtableClusterName:"b",
		CloudtableTableName:"c",
		CloudtableRowkeyDelimiter:".",
		ObsBackupBucketPath:"dis-tests-not-delete",
		BackupfilePrefix:"cloudtable",
	}

	rowKeySchema := []models.SchemaField{
		models.SchemaField{
			Value:"id",
			Type:models.JSON_TYPE_STRING,
		},
	}

	columnsSchema := []models.SchemaField{
		models.SchemaField{
			ColumnFamilyName:"a",
			ColumnName:"ID",
			Value:"id",
			Type:models.JSON_TYPE_STRING,
		},
	}
	cloudtableSchema := &models.CloudtableSchema{
		RowKeySchema:rowKeySchema,
		ColumnsSchema:columnsSchema,
	}
	descriptor.CloudtableSchema = cloudtableSchema

	input := &models.CreateTransferTaskRequest{
		StreamName    :disInstance.StreamName,
		DestinationType:models.DESTINATION_CLOUDTABLE,
		CloudtableDestinationDescriptor:descriptor,
	}
	result, _ := dis.CreateTransferTask(input)
	if result.Err != nil {
		fmt.Printf("Failed to create TransferTask [%s].\n", descriptor.TransferTaskName);
	} else {
		fmt.Printf("Success to create TransferTask [%s].\n", descriptor.TransferTaskName);
	}
}

/**
*函数功能：创建目标服务为Cloudtable HBase的转储任务
 */
func CreateCloudtablOpentsdbTransferTask(disInstance DISInstance) {
	dis := disInstance.Dis

	descriptor := &models.CloudtableDestinationDescriptorRequest{
		TransferTaskName: "feihang",
		AgencyName:"all",
		ConsumerStrategy:models.TRIM_HORIZON,
		CloudtableClusterId:"a",
		CloudtableClusterName:"b",
		CloudtableTableName:"c",
		ObsBackupBucketPath:"dis-tests-not-delete",
		BackupfilePrefix:"cloudtable",
	}

	metricSchema := []models.SchemaField{
		models.SchemaField{
			Value:"id",
			Type:models.JSON_TYPE_STRING,
		},
	}

	timestampSchema := models.SchemaField{
		Value:"date",
		Type:models.JSON_TYPE_STRING,
		Format:models.TIMESTAMP_FORMAT_0,
	}

	valueSchema := models.SchemaField{
		Value:"value",
		Type:models.JSON_TYPE_STRING,
	}

	tagsSchema := []models.SchemaField{
		models.SchemaField{
			Name:"ID",
			Value:"id",
			Type:models.JSON_TYPE_STRING,
		},
	}
	opentsdbSchema := models.OpenTSDBSchema{
		MetricSchema:metricSchema,
		TimestampSchema:timestampSchema,
		ValueSchema:valueSchema,
		TagsSchema:tagsSchema,
	}
	descriptor.OpentsdbSchema = []models.OpenTSDBSchema{
		opentsdbSchema,
	}

	input := &models.CreateTransferTaskRequest{
		StreamName    :disInstance.StreamName,
		DestinationType:models.DESTINATION_CLOUDTABLE,
		CloudtableDestinationDescriptor:descriptor,
	}
	result, _ := dis.CreateTransferTask(input)
	if result.Err != nil {
		fmt.Printf("Failed to create TransferTask [%s].\n", descriptor.TransferTaskName);
	} else {
		fmt.Printf("Success to create TransferTask [%s].\n", descriptor.TransferTaskName);
	}
}

//查询TransferTask详情
func DescTransferTask(disInstance DISInstance) {
	dis := disInstance.Dis
	input := &models.DescribeTransferTaskRequest{
		StreamName    :disInstance.StreamName,
		TaskName:"feihang",
	}
	result, output := dis.DescTransferTask(input)
	if result.Err != nil {
		fmt.Printf("Failed to desc transfer task [%s:%s].\n", input.StreamName, input.TaskName);
	} else {
		fmt.Printf("desc transfer task success");

		outputbytes, err := json.Marshal(output)
		if nil != err {
			logger.LOG(logger.ERROR, "CommitCheckpointRequest parameter error ", err)
		}
		fmt.Printf(": %s\n", string(outputbytes));
	}
}

//查询TransferTask列表
func ListTransferTasks(disInstance DISInstance) {
	input := &models.ListTransferTasksRquest{
		Limit:10,
		StreamName: disInstance.StreamName,
	}
	dis := disInstance.Dis
	result, output := dis.ListTransferTasks(input)
	if result.Err != nil {
		fmt.Printf("Failed to list transferTasks.\n");
	} else {
		fmt.Printf("ListTransferTasks, size [%d].\n", output.TaskNumber);
		for i, task := range output.Tasks {
			fmt.Printf("task%d name [%s], DestinationType [%s], createTime [%d].\n", i, task.TaskName, task.DestinationType, task.CreateTime);
		}
	}
}

//删除特定的TransferTask
func DeleteTransferTask(disInstance DISInstance) {
	dis := disInstance.Dis
	input := &models.DeleteTransferTaskRequest{
		StreamName    :disInstance.StreamName,
		TaskName:"feihang",
	}
	result, _ := dis.DeleteTransferTask(input)
	if result.Err != nil {
		fmt.Printf("Failed to delete transfer task [%s:%s].\n", input.StreamName, input.TaskName);
	} else {
		fmt.Printf("delete transfer task [%s:%s] success.\n", input.StreamName, input.TaskName);
	}
}