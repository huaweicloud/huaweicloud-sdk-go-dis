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

package models

//接口执行信息集合
type Result struct {
	Err        error //执行的错误信息，当Err为nil是，其他值有效
	StatusCode int   //http 返回的状态码
	ErrResponse
}

//http 返回错误信息
type ErrResponse struct {
	Code      string `json:"errorCode"` //错误信息码
	Message   string `json:"message"`   //错误信息
	RequestId string `json:"RequestId"` //本次错误请求的请求ID
	HostId    string `json:"HostId"`    //返回该消息的服务端ID
}

//创建通道请求
type CreateStreamRequest struct {
	StreamName     string `json:"stream_name"`  //通道名称
	PartitionCount int `json:"partition_count"` //分片数量
	StreamType     string `json:"stream_type"`  //通道类型：普通（对应带宽上限1MB），高级（对应带宽上限5MB）
	DataDuration   int `json:"data_duration"`   //数据保留时长，单位为小时，默认为24小时
	DataType       string `json:"data_type"`    //数据的类型，目前支持：BLOB（默认，无格式要求）、JSON、CSV格式
	DataSchema     string `json:"data_schema"`  //用户JOSN、CSV格式数据schema，用于需要数据预处理的转储任务，用avro shema描述
	Tags           []Tag   `json:"tags"`        //通道标签列表
}

//创建通道响应
type CreateStreamResult struct {
}

//通道的标签
type Tag struct {
	Key   string `json:"key"`   //键，最大长度36个unicode字符
	Value string `json:"value"` //值，最大长度为43个unicode字符
}

//获取通道详情请求
type DescribeStreamRequest struct {
	StreamName       string `json:"stream_name"`      //通道名称
	StartPartitionId string `json:"startPartitionId"` //从该分区值开始返回分区列表，返回的分区列表不包括此分区。
	LimitPartitions  int `json:"limit_partitions"`    //单次请求返回的最大分区数。最小值是1，最大值是1000；默认值是100。
}

//获取通道详情响应
type DescribeStreamResult struct {
	StreamName             string `json:"stream_name"`                             //通道名称
	StreamId               string `json:"stream_id"`                               //通道ID
	Status                 string `json:"status"`                                  //通道的当前状态：CREATING - 创建中，RUNNING - 运行中，TERMINATING - 删除中
	StreamType             string `json:"stream_type"`                             //通道类型：普通（对应带宽上限1MB），高级（对应带宽上限5MB）
	RetentionPeriod        int `json:"retention_period"`                           //数据保留时长
	CreateTime             int64 `json:"create_time"`                              //通道创建时间
	LastModifiedTime       int64 `json:"last_modified_time"`                       //通道最近修改时间
	WritablePartitionCount int `json:"writable_partition_count"`                   //可写分区(即ACTIVE状态)数量
	ReadablePartitionCount int `json:"readable_partition_count"`                   //可读分区(即ACTIVE与DELETE状态)数量
	Partitions             []PartitionResult `json:"partitions"`                   //通道的分区列表信息
	UpdatePartitionCounts  []UpdatePartitionCount `json:"update_partition_counts"` //扩缩容信息
	DataType               string `json:"data_type"`                               //数据的类型，目前支持：BLOB（默认，无格式要求）、JSON、CSV格式
	DataSchema             string `json:"data_schema"`                             //用户JOSN、CSV格式数据schema，用于需要数据预处理的转储任务，用avro shema描述
	Tags                   []Tag   `json:"tags"`                                   //通道标签列表
	HasMorePartitions      bool `json:"has_more_partitions"`                       //是否有分片未返回
}


//通道的分区信息
type PartitionResult struct {
	Status              string `json:"status"`                //分区的当前状态：CREATING - 创建中，ACTIVE - 运行中，DELETED - 已删除，EXPIRED - 已过期
	PartitionId         string `json:"partition_id"`          //分区ID
	ParentPartitionIds  string `json:"parent_partitions"`     //父分区
	HashRange           string `json:"hash_range"`            //分区的可能哈希键值范围
	SequenceNumberRange string `json:"sequence_number_range"` //分区的可能序列号范围
}

//扩缩容信息
type UpdatePartitionCount struct {
	CreateTimestamp      int64 `json:"create_timestamp"`     //创建时间
	SrcPartitionCount    int `json:"src_partition_count"`    //源分区数目
	TargetPartitionCount int `json:"target_partition_count"` //目标分区数目
	ResultCode           string `json:"result_code"`
	ResultMsg            string `json:"result_msg"`
}


//通道扩容请求
type UpdatePartitionCountRequest struct {
	StreamName           string `json:"stream_name"`         //通道名称
	TargetPartitionCount int `json:"target_partition_count"` //目标分区数目
}

//通道扩容响应
type UpdatePartitionCountResult struct {
	StreamName            string `json:"stream_name"`          //通道名称
	CurrentPartitionCount int `json:"current_partition_count"` //扩容前的分区数目
	TargetPartitionCount  int `json:"target_partition_count"`  //目标分区数目
}

//查询通道列表
type ListStreamsRequest struct {
	ExclusivetartStreamName string `json:"start_stream_name"` //The name of the stream to start the list with. Exclude this stream name.
	Limit                   int `json:"limit"`                //The maximum number of streams to list. Default value is 10.
}

//查询通道响应
type ListStreamsResult struct {
	StreamNumber   int `json:"total_number"`
	HasMoreStreams bool`json:"has_more_streams"`
	StreamInfos    []StreamInfo`json:"stream_info_list"`
	StreamNames    []string`json:"stream_names"` //废弃字段，不建议使用
}

//通道的详情列表
type StreamInfo struct {
	StreamName      string `json:"stream_name"`   //通道名称
	Status          string `json:"status"`        //通道的当前状态：CREATING - 创建中，RUNNING - 运行中，TERMINATING - 删除中
	StreamType      string `json:"stream_type"`   //通道类型：普通（对应带宽上限1MB），高级（对应带宽上限5MB）
	RetentionPeriod int `json:"retention_period"` //数据保留时长
	CreateTime      int64 `json:"create_time"`    //通道创建时间
	DataType        string `json:"data_type"`     //数据的类型，目前支持：BLOB（默认，无格式要求）、JSON、CSV格式
	PartitionCount  int `json:"partition_count"`  //通道的分区数目
	Tags            []Tag   `json:"tags"`         //通道标签列表
}

//删除通道请求
type DeleteStreamRequest struct {
	StreamName string `json:"stream_name"` //通道名称
}

//删除通道响应
type DeleteStreamResult struct {
}

//创建转储任务请求
type CreateTransferTaskRequest struct {
	StreamName                      string `json:"stream_name"`                                                        //通道名称
	DestinationType                 string `json:"destination_type"`                                                   //任务类型：OBS、MRS、DWS、DLI、CLOUDTABLE
	ObsDestinationDescriptor        *OBSDestinationDescriptorRequest `json:"obs_destination_descriptor"`               //OBS转储任务
	MrsDestinationDescriptor        *MRSDestinationDescriptorRequest `json:"mrs_destination_descriptor"`               //MRS转储任务
	DliDestinationDescriptor        *DLIDestinationDescriptorRequest `json:"dli_destination_descriptor"`               //DLI转储任务
	DwsDestinationDescriptor        *DWSDestinationDescriptorRequest `json:"dws_destination_descriptor"`               //DWS转储任务
	CloudtableDestinationDescriptor *CloudtableDestinationDescriptorRequest `json:"cloudtable_destination_descriptor"` //Cloudtable转储任务
}

//OBS转储任务请求
type OBSDestinationDescriptorRequest struct {
	TransferTaskName    string `json:"task_name"`                              //任务名称
	AgencyName          string `json:"agency_name"`                            //IAM委托名称
	FilePrefix          string `json:"file_prefix"`                            //Directory to hold files that will be dumped to OB
	PartitionFormat     string `json:"partition_format"`                       //Directory structure of the Object file written into OBS.
	ObsBucketPath       string `json:"obs_bucket_path"`                        //Name of the OBS bucket used to store data from the DIS stream
	DeliverTimeInterval int `json:"deliver_time_interval"`                     //周期转储的间隔
	ConsumerStrategy    string `json:"consumer_strategy"`                      //设置从DIS拉取数据时的初始偏移量: 默认LATEST - 从最新的记录开始读取; TRIM_HORIZON - 从最早的记录开始读取
	DestinationFileType string `json:"destination_file_type"`                  //Type of the Object file written into OBS, such as text, parquet, carbon. Default value: text.
	RecordDelimiter     string `json:"record_delimiter"`                       //记录分割符
	CarbonProperties    map[string]string `json:"carbon_properties,omitempty"` //carbon转储的属性设置
	ProcessingSchema    *ProcessingSchema `json:"processing_schema"`           //数据转换的schema配置:如支持parquet按照指定timestamp生成分区目录
}

// 转储时需要数据预处理的schema
type ProcessingSchema struct {
	TimestampName   string `json:"timestamp_name,omitempty"`   //The name of the timestamp field.
	TimestampType   string `json:"timestamp_type,omitempty"`   //The type of the timestamp field.
	TimestampFormat string `json:"timestamp_format,omitempty"` //The format of the timestamp field
}

//MRS转储任务请求
type MRSDestinationDescriptorRequest struct {
	TransferTaskName    string `json:"task_name"`                              //任务名称
	AgencyName          string `json:"agency_name"`                            //IAM委托名称
	ObsBucketPath       string `json:"obs_bucket_path"`                        //临时中转目录
	DeliverTimeInterval int `json:"deliver_time_interval"`                     //周期转储的间隔
	MrsClusterName      string `json:"mrs_cluster_name"`                       //String Name of the MRS cluster to which data in the DIS stream will be dumped
	MrsClusterId        string `json:"mrs_cluster_id"`                         //ID of the MRS cluster to which data in the DIS stream will be dumped
	MrsHdfsPath         string `json:"mrs_hdfs_path"`                          //Hadoop Distributed File System (HDFS) path of the MRS cluster to which data in the DIS stream will be dumped
	HdfsPrefixFolder    string `json:"hdfs_prefix_folder"`                     //Directory to hold files that will be dumped to MRS.
	ConsumerStrategy    string `json:"consumer_strategy"`                      //设置从DIS拉取数据时的初始偏移量: 默认LATEST - 从最新的记录开始读取; TRIM_HORIZON - 从最早的记录开始读取
	DestinationFileType string `json:"destination_file_type"`                  //Type of the Object file written into OBS, such as text, parquet, carbon. Default value: text.
	CarbonProperties    map[string]string `json:"carbon_properties,omitempty"` //carbon转储的属性设置
}

//DLI转储任务请求
type DLIDestinationDescriptorRequest struct {
	TransferTaskName    string `json:"task_name"`          //任务名称
	AgencyName          string `json:"agency_name"`        //IAM委托名称
	ObsBucketPath       string `json:"obs_bucket_path"`    //Name of the OBS bucket used to store data from the DIS stream
	DliDatabaseName     string `json:"dli_database_name"`  //Name of the DLI database to which data in the DIS stream will be dumped.
	DliTableName        string `json:"dli_table_name"`     //Name of the DLI table to which data in the DIS stream will be dumped.
	DeliverTimeInterval int `json:"deliver_time_interval"` //周期转储的间隔
	ConsumerStrategy    string `json:"consumer_strategy"`  //设置从DIS拉取数据时的初始偏移量: 默认LATEST - 从最新的记录开始读取; TRIM_HORIZON - 从最早的记录开始读取
}

//DWS转储任务请求
type DWSDestinationDescriptorRequest struct {
	TransferTaskName    string `json:"task_name"`          //任务名称
	AgencyName          string `json:"agency_name"`        //IAM委托名称
	ObsBucketPath       string `json:"obs_bucket_path"`    //Name of the OBS bucket used to store data from the DIS stream
	DwsClusterName      string `json:"dws_cluster_name"`   //Name of the DWS cluster used to store data in the DIS stream.
	DwsClusterId        string `json:"dws_cluster_id"`     //ID of the DWS cluster used to store data in the DIS stream.
	DwsSchema           string `json:"dws_schema"`         //Schema of the DWS database used to store data in the DIS stream.
	DwsDatabaseName     string `json:"dws_database_name"`  //Name of the DWS database used to store data in the DIS stream
	UserName            string `json:"user_name"`          //Username of the DWS database used to store data in the DIS stream.
	UserPassword        string `json:"user_password"`      //Password of the DWS database used to store data in the DIS stream.
	DwsTableName        string `json:"dws_table_name"`     //Name of the table in the DWS database used to store data in the DIS stream.
	DwsDelimiter        string `json:"dws_delimiter"`      //ID of the DWS cluster used to store data in the DIS stream.
	KmsUserKeyName      string `json:"kms_user_key_name"`  //Key created in Key Management Service (KMS) and used to encrypt the password of the DWS database.
	KmsUserKeyId        string `json:"kms_user_key_id"`    //ID of the key created in Key Management Service (KMS) and used to encrypt the password of the DWS database
	DeliverTimeInterval int `json:"deliver_time_interval"` //周期转储的间隔
	ConsumerStrategy    string `json:"consumer_strategy"`  //设置从DIS拉取数据时的初始偏移量: 默认LATEST - 从最新的记录开始读取; TRIM_HORIZON - 从最早的记录开始读取
}

//Cloudtable转储任务请求
type CloudtableDestinationDescriptorRequest struct {
	TransferTaskName          string `json:"task_name"`                    //任务名称
	AgencyName                string `json:"agency_name"`                  //IAM委托名称
	CloudtableClusterName     string `json:"cloudtable_cluster_name"`      //Name of the Cloudtable cluster used to store data in the DIS stream.
	CloudtableClusterId       string `json:"cloudtable_cluster_id"`        //ID of the Cloudtable cluster used to store data in the DIS stream.
	CloudtableTableName       string `json:"cloudtable_cluster_id"`        //ID of the Cloudtable cluster used to store data in the DIS stream.
	CloudtableRowkeyDelimiter string `json:"cloudtable_table_name"`        //HBase table name of the CloudTable cluster to which data will be dumped
	CloudtableSchema          *CloudtableSchema `json:"cloudtable_schema"` //Schema configuration of the CloudTable HBase data.
	OpentsdbSchema            []OpenTSDBSchema `json:"opentsdb_schema"`    //Schema configuration of the CloudTable OpenTSDB data.
	ObsBackupBucketPath       string `json:"obs_backup_bucket_path"`       //Name of the OBS bucket used to back up data that failed to be dumped to CloudTable
	BackupfilePrefix          string `json:"backup_file_prefix"`           //Name of the OBS bucket used to back up data that failed to be dumped to CloudTable
	ConsumerStrategy          string `json:"consumer_strategy"`            //设置从DIS拉取数据时的初始偏移量: 默认LATEST - 从最新的记录开始读取; TRIM_HORIZON - 从最早的记录开始读取
}

// 转储Cloudtable HBase的数据处理schema
type CloudtableSchema struct {
	RowKeySchema  []SchemaField `json:"row_key"` //HBase rowkey Schema used by the CloudTable cluster to convert JSON data into HBase rowkeys.
	ColumnsSchema []SchemaField `json:"columns"` //HBase column Schema used by the CloudTable cluster to convert JSON data into HBase columns.
}

// 转储Cloudtable Opentsdb的数据处理schema
type OpenTSDBSchema struct {
	MetricSchema    []SchemaField `json:"metric"`  //Schema configuration of the OpenTSDB data metric
	TimestampSchema SchemaField `json:"timestamp"` // Schema configuration of the OpenTSDB data timestamp
	ValueSchema     SchemaField `json:"value"`     //Schema configuration of the OpenTSDB data value
	TagsSchema      []SchemaField `json:"tags"`    //Schema configuration of the OpenTSDB data tags
}

type SchemaField struct {
	ColumnFamilyName string `json:"column_family_name"` //Name of the HBase column family to which data will be dumped.
	ColumnName       string `json:"column_name"`        //Name of the HBase column to which data will be dumped.
	Name             string `json:"name"`               //The format of the timestamp field
	Format           string `json:"format"`             //The format of the schema field.
	Value            string `json:"value"`              //The value of the schema field.
	Type             string `json:"type"`               //The type of the schema field.
}

//创建转储任务响应
type CreateTransferTaskResult struct {
}

//查看转储任务列表请求
type ListTransferTasksRquest struct {
	StreamName             string `json:"stream_name"`     //通道名称
	Limit                  int `json:"limit"`              //The maximum number of tasks to list. Default value is 10.
	ExclusiveStartTaskName string `json:"start_task_name"` //任务名称
}

//查看转储任务列表响应
type ListTransferTasksResult struct {
	TaskNumber int `json:"total_number"`           //任务数目
	Tasks      []TransferTaskResult `json:"tasks"` //任务列表
}

type TransferTaskResult struct {
	DestinationType       string `json:"destination_type"`       //The destination type of the delivery task. For Example, OBS.
	TaskName              string `json:"task_name"`              //The name of the delivery task.
	State                 string `json:"state"`                  //The transfer state of the delivery task.
	CreateTime            int64 `json:"create_time"`             //The create time of the transfer task.
	LastTransferTimeStamp int64 `json:"last_transfer_timestamp"` //The lastest transfer timeStamp of the transfer task.
}

//查看转储任务详情请求
type DescribeTransferTaskRequest struct {
	StreamName string `json:"stream_name"` //通道名称
	TaskName   string `json:"task_name"`   //任务名称
}

//查看转储任务详情响应
type DescribeTransferTaskResult struct {
	StreamName                       string `json:"stream_name"`                   //通道名称
	TaskName                         string `json:"task_name"`                     //任务名称
	DestinationType                  string `json:"destination_type"`
	State                            string `json:"state"`
	CreateTime                       int64 `json:"create_time"`                    //The create time of the transfer task.
	LastTransferTimeStamp            int64 `json:"last_transfer_timestamp"`        //The lastest transfer timeStamp of the transfer task.
	Partitions                       []PartitionTransferResult `json:"partitions"` //The transfer details of the partitions.
	ObsDestinationDescription        *OBSDestinationDescription `json:"obs_destination_description,omitempty"`
	MrsDestinationDescription        *MRSDestinationDescription `json:"mrs_destination_description,omitempty"`
	DliDestinationDescription        *DLIDestinationDescription `json:"dli_destination_description,omitempty"`
	DwsDestinationDescription        *DWSDestinationDescription `json:"dws_destination_description,omitempty"`
	CloudtableDestinationDescription *CloudtableDestinationDescription `json:"cloudtable_destination_description,omitempty"`
}

type PartitionTransferResult struct {
	State                 string `json:"state"`
	LastTransferTimeStamp int64 `json:"last_transfer_timestamp"` //The lastest transfer timeStamp of the transfer task.
	Discard               int64 `json:"discard"`                 //脏数据量
}

//OBS转储任务详情
type OBSDestinationDescription struct {
	AgencyName          string `json:"agency_name"`                            //IAM委托名称
	FilePrefix          string `json:"file_prefix"`                            //Directory to hold files that will be dumped to OB
	PartitionFormat     string `json:"partition_format"`                       //Directory structure of the Object file written into OBS.
	ObsBucketPath       string `json:"obs_bucket_path"`                        //Name of the OBS bucket used to store data from the DIS stream
	DeliverTimeInterval int `json:"deliver_time_interval"`                     //周期转储的间隔
	ConsumerStrategy    string `json:"consumer_strategy"`                      //设置从DIS拉取数据时的初始偏移量: 默认LATEST - 从最新的记录开始读取; TRIM_HORIZON - 从最早的记录开始读取
	DestinationFileType string `json:"destination_file_type"`                  //Type of the Object file written into OBS, such as text, parquet, carbon. Default value: text.
	RecordDelimiter     string `json:"record_delimiter"`                       //记录分割符
	CarbonProperties    map[string]string `json:"carbon_properties,omitempty"` //carbon转储的属性设置
	ProcessingSchema    *ProcessingSchema `json:"processing_schema,omitempty"` //数据转换的schema配置:如支持parquet按照指定timestamp生成分区目录
}

//MRS转储任务详情
type MRSDestinationDescription struct {
	AgencyName          string `json:"agency_name"`                            //IAM委托名称
	ObsBucketPath       string `json:"obs_bucket_path"`                        //临时中转目录
	DeliverTimeInterval int `json:"deliver_time_interval"`                     //周期转储的间隔
	MrsClusterName      string `json:"mrs_cluster_name"`                       //String Name of the MRS cluster to which data in the DIS stream will be dumped
	MrsClusterId        string `json:"mrs_cluster_id"`                         //ID of the MRS cluster to which data in the DIS stream will be dumped
	MrsHdfsPath         string `json:"mrs_hdfs_path"`                          //Hadoop Distributed File System (HDFS) path of the MRS cluster to which data in the DIS stream will be dumped
	HdfsPrefixFolder    string `json:"hdfs_prefix_folder"`                     //Directory to hold files that will be dumped to MRS.
	ConsumerStrategy    string `json:"consumer_strategy"`                      //设置从DIS拉取数据时的初始偏移量: 默认LATEST - 从最新的记录开始读取; TRIM_HORIZON - 从最早的记录开始读取
	DestinationFileType string `json:"destination_file_type"`                  //Type of the Object file written into OBS, such as text, parquet, carbon. Default value: text.
	CarbonProperties    map[string]string `json:"carbon_properties,omitempty"` //carbon转储的属性设置
}

//DLI转储任务详情
type DLIDestinationDescription struct {
	AgencyName          string `json:"agency_name"`        //IAM委托名称
	ObsBucketPath       string `json:"obs_bucket_path"`    //Name of the OBS bucket used to store data from the DIS stream
	DliDatabaseName     string `json:"dli_database_name"`  //Name of the DLI database to which data in the DIS stream will be dumped.
	DliTableName        string `json:"dli_table_name"`     //Name of the DLI table to which data in the DIS stream will be dumped.
	DeliverTimeInterval int `json:"deliver_time_interval"` //周期转储的间隔
	ConsumerStrategy    string `json:"consumer_strategy"`  //设置从DIS拉取数据时的初始偏移量: 默认LATEST - 从最新的记录开始读取; TRIM_HORIZON - 从最早的记录开始读取
}

//DWS转储任务详情
type DWSDestinationDescription struct {
	AgencyName          string `json:"agency_name"`        //IAM委托名称
	ObsBucketPath       string `json:"obs_bucket_path"`    //Name of the OBS bucket used to store data from the DIS stream
	DwsClusterName      string `json:"dws_cluster_name"`   //Name of the DWS cluster used to store data in the DIS stream.
	DwsClusterId        string `json:"dws_cluster_id"`     //ID of the DWS cluster used to store data in the DIS stream.
	DwsSchema           string `json:"dws_schema"`         //Schema of the DWS database used to store data in the DIS stream.
	DwsDatabaseName     string `json:"dws_database_name"`  //Name of the DWS database used to store data in the DIS stream
	UserName            string `json:"user_name"`          //Username of the DWS database used to store data in the DIS stream.
	UserPassword        string `json:"user_password"`      //Password of the DWS database used to store data in the DIS stream.
	DwsTableName        string `json:"dws_table_name"`     //Name of the table in the DWS database used to store data in the DIS stream.
	DwsDelimiter        string `json:"dws_delimiter"`      //ID of the DWS cluster used to store data in the DIS stream.
	KmsUserKeyName      string `json:"kms_user_key_name"`  //Key created in Key Management Service (KMS) and used to encrypt the password of the DWS database.
	KmsUserKeyId        string `json:"kms_user_key_id"`    //ID of the key created in Key Management Service (KMS) and used to encrypt the password of the DWS database
	DeliverTimeInterval int `json:"deliver_time_interval"` //周期转储的间隔
	ConsumerStrategy    string `json:"consumer_strategy"`  //设置从DIS拉取数据时的初始偏移量: 默认LATEST - 从最新的记录开始读取; TRIM_HORIZON - 从最早的记录开始读取
}

//Cloudtable转储任务详情
type CloudtableDestinationDescription struct {
	AgencyName                string `json:"agency_name"`                  //IAM委托名称
	CloudtableClusterName     string `json:"cloudtable_cluster_name"`      //Name of the Cloudtable cluster used to store data in the DIS stream.
	CloudtableClusterId       string `json:"cloudtable_cluster_id"`        //ID of the Cloudtable cluster used to store data in the DIS stream.
	CloudtableTableName       string `json:"cloudtable_cluster_id"`        //ID of the Cloudtable cluster used to store data in the DIS stream.
	CloudtableRowkeyDelimiter string `json:"cloudtable_table_name"`        //HBase table name of the CloudTable cluster to which data will be dumped
	CloudtableSchema          *CloudtableSchema `json:"cloudtable_schema"` //Schema configuration of the CloudTable HBase data.
	OpentsdbSchema            []OpenTSDBSchema `json:"opentsdb_schema"`    //Schema configuration of the CloudTable OpenTSDB data.
	ObsBackupBucketPath       string `json:"obs_backup_bucket_path"`       //Name of the OBS bucket used to back up data that failed to be dumped to CloudTable
	BackupfilePrefix          string `json:"backup_file_prefix"`           //Name of the OBS bucket used to back up data that failed to be dumped to CloudTable
	ConsumerStrategy          string `json:"consumer_strategy"`            //设置从DIS拉取数据时的初始偏移量: 默认LATEST - 从最新的记录开始读取; TRIM_HORIZON - 从最早的记录开始读取
}

//删除转储任务请求
type DeleteTransferTaskRequest struct {
	StreamName string `json:"stream_name"` //通道名称
	TaskName   string `json:"task_name"`   //任务名称
}

//删除转储任务响应
type DeleteTransferTaskResult struct {
}

//获取迭代器请求
type GetPartitionCursorRequest struct {
	StreamName             string `json:"stream-name"`                //通道名称
	PartitionId            string `json:"partition-id"`               //分区值
	CursorType             string   `json:"cursor-type"`              //迭代类型：AT_SEQUENCE_NUMBER - 从特定序列号所在的记录开始读取。 AFTER_SEQUENCE_NUMBER - 从特定序列号后的记录开始读取。 TRIM_HORIZON - 从分区中最时间最长的记录开始读取。 LATEST - 从分区中最新的记录开始读取。
	StartingSequenceNumber string   `json:"starting-sequence-number"` //序列号。序列号是每个记录的唯一标识符。序列号由DIS在数据生产者调用 <code>PutRecord</code> 操作以添加数据到DIS数据通道时DIS服务自动分配的。同一分区键的序列号通常会随时间变化增加。
	Timestamp              int64   `json:"timestamp"`                 //时间戳。
}

//获取迭代器响应
type GetPartitionCursorResult struct {
	PartitionCursor string `json:"partition_cursor"` //迭代器
}

//获取数据请求
type GetRecordsRequest struct {
	PartitionCursor string `json:"partition-cursor"` //迭代器
}

//获取数据响应
type GetRecordsResult struct {
	Records             []Record `json:"records"`             //下载的数据列表
	NextPartitionCursor string `json:"next_partition_cursor"` //下载的数据列表
	MillisBehindLatest  int64 `json:"millis_behind_latest"`
}


//获取数据响应
type Record struct {
	PartitionKey                string `json:"partition_key"`              //数据写入分区的分区键
	SequenceNumber              string `json:"sequence_number"`            //序列号是每个记录的唯一标识符。序列号由DIS在数据生产者调用 PutRecord 操作以添加数据到DIS数据通道时DIS服务自动分配的。同一分区键的序列号通常会随时间变化增加。
	Data                        []byte `json:"data"`                       //数据
	ApproximateArrivalTimestamp int64 `json:"ApproximateArrivalTimestamp"` //
	Timestamp                   int64 `json:"timestamp"`                   //时间戳
	TimestampType               string `json:"timestamp_type"`
}

//上传数据请求
type PutRecordsRequest struct {
	StreamName string                   `json:"stream_name"` // 通道名称
	StreamId   string                   `json:"stream_id"`   // 通道ID，用于授权访问， 与 StreamName 二选一
	Records    []PutRecordsRequestEntry `json:"records"`     // 记录列表
}

//记录列表
type PutRecordsRequestEntry struct {
	Data            []byte `json:"data"`              //数据
	PartitionId     string `json:"partition_id"`      //数据写入分区的分区值， 优先使用
	ExplicitHashKey string `json:"explicit_hash_key"` //用于明确数据需要写入分区的哈希值，此哈希值将覆盖“partition_key”的哈希值
	PartitionKey    string `json:"partition_key"`     //数据写入分区的分区键
}

//上传数据响应
type PutRecordsResult struct {
	FailedRecordCount int `json:"failed_record_count"`         //上传失败的数据数量
	Records           []PutRecordsResultEntry `json:"records"` //上传结果列表
}

//上传结果列表
type PutRecordsResultEntry struct {
	PartitionId    string `json:"partition_id"`    //数据写入分区的分区值， 优先使用
	SequenceNumber string `json:"sequence_number"` //序列号是每个记录的唯一标识符
	ErrorCode      string `json:"error_code"`      //错误码
	ErrorMessage   string `json:"error_message"`   //错误消息
}

//创建APP请求
type CreateAppRequest struct {
	AppName string `json:"app_name"` //APP名称
}

//描述APP响应
type DescribeAppResult struct {
	AppName    string `json:"app_name"`   //APP名称
	AppId      string `json:"app_id"`     //APPId
	CreateTime int64 `json:"create_time"` //APP创建时间
}

//查询APP列表请求
type ListAppsRequest struct {
	ExclusiveStartAppName string `json:"start_app_name"` //The name of the apps to start the list with. Exclude this apps name.
	Limit                 int `json:"limit"`             //The maximum number of apps to list. Default value is 10.
}

//查询APP列表响应
type ListAppsResult struct {
	HasMoreApp bool `json:"has_more_app"`
	Apps       []DescribeAppResult `json:"apps"`
}

//提交checkpoint请求
type CommitCheckpointRequest struct {
	AppName        string `json:"app_name"`        //APP名称
	CheckpointType string `json:"checkpoint_type"` //checkpoint类型
	StreamName     string `json:"stream_name"`     //通道名称
	PartitionId    string `json:"partition_id"`    //分区值
	SequenceNumber string `json:"sequence_number"` //序列号是每个记录的唯一标识符
	Metadata       string `json:"metadata"`        //用户自定义元数据信息
}

//提交checkpoint响应
type CommitCheckpointResult struct {
}

//获取checkpoint请求
type GetCheckpointRequest struct {
	StreamName     string `json:"stream_name"`     //通道名称
	PartitionId    string `json:"partition_id"`    //分区值
	AppName        string `json:"app_name"`        //APP名称
	CheckpointType string `json:"checkpoint_type"` //checkpoint类型
	Timestamp      int64 `json:"timestamp"`
}

//获取checkpoint响应
type GetCheckpointResult struct {
	SequenceNumber string `json:"sequence_number"` //序列号是每个记录的唯一标识符
	Metadata       string `json:"metadata"`        //用户自定义元数据信息
}
//删除checkpoint请求
type DeleteCheckpointRequest struct {
	StreamName     string `json:"stream_name"`     //通道名称
	PartitionId    string `json:"partition_id"`    //分区值
	AppName        string `json:"app_name"`        //APP名称
	CheckpointType string `json:"checkpoint_type"` //checkpoint类型
}

//删除checkpoint响应
type DeleteCheckpointResult struct {
}


//DIS客户端相关配置项
type DISClientConf struct {
	AK                string
	SK                string
	ProjectId         string
	Region            string
	Endpoint          string
	BodySerializeType string
	CacheResendConf   *CacheResendConf
}

func DefaultDISClientConf() *DISClientConf {
	return &DISClientConf{
		BodySerializeType: JSON,
		CacheResendConf: DefaultCacheResendConf(),
	}
}

//本地缓存重发相关配置项
type CacheResendConf struct {
	DataCacheEnable           bool
	DataCacheDir              string
	DataCacheDiskMaxSize      int64
	DataCacheArchiveSize      int64
	DataCacheArchiveLifeCycle int64
}

func DefaultCacheResendConf() *CacheResendConf {
	return &CacheResendConf{
		DataCacheEnable:false,
		DataCacheDir: PROPERTY_DATA_CACHE_DIR,
		DataCacheDiskMaxSize: PROPERTY_DATA_CACHE_DISK_MAX_SIZE,
		DataCacheArchiveSize:PROPERTY_DATA_CACHE_ARCHIVE_MAX_SIZE,
		DataCacheArchiveLifeCycle:PROPERTY_DATA_CACHE_ARCHIVE_LIFE_CYCLE,
	}
}