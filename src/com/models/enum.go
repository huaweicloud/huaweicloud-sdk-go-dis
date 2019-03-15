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

//PartitionCursorTypeEnum
const (
	AT_SEQUENCE_NUMBER = "AT_SEQUENCE_NUMBER"
	AFTER_SEQUENCE_NUMBER = "AFTER_SEQUENCE_NUMBER"
	AT_TIMESTAMP = "AT_TIMESTAMP"
	TRIM_HORIZON = "TRIM_HORIZON"
	LATEST = "LATEST"
)

//BodySerializeType
const (
	PROTOBUF = "protobuf"
	JSON = "json"
)

//CheckpointTypeEnum
const (
	LAST_READ = "LAST_READ"
)

//通道类型
const (
	STREAM_TYPE_COMMON = "COMMON"
	STREAM_TYPE_ADVANCED = "ADVANCED"
)

//数据类型
const (
	DATA_TYPE_BLOB = "BLOB"
	DATA_TYPE_JSON = "JSON"
	DATA_TYPE_CSV = "CSV"
)

//通道的生命周期
const (
	DATA_DURATION_DEFAULT = 24
)

//转储任务的类型
const (
	DESTINATION_DEFAULT = "NOWHERE"
	DESTINATION_OBS = "OBS"
	DESTINATION_MRS = "MRS"
	DESTINATION_DLI = "DLI"
	DESTINATION_CLOUDTABLE = "CLOUDTABLE"
	DESTINATION_OPENTSDB = "OPENTSDB"
	DESTINATION_DWS = "DWS"
)

//转储OBS/MRS的目标文件格式
const (
	TEXT = "text"
	PARQUET = "parquet"
	CARBON = "carbon"
)

//时间戳的类型
const (
	TIMESTAMP_TYPE_STRING = "STRING"
	TIMESTAMP_TYPE_TIMESTAMP = "TIMESTAMP"
)

//数据类型
const (
	JSON_TYPE_BIGINT = "bigint"
	JSON_TYPE_DOUBLE = "double"
	JSON_TYPE_BOOLEAN = "boolean"
	JSON_TYPE_STRING = "string"
	JSON_TYPE_DECIMAL = "decimal"
	JSON_TYPE_CONSTANT = "constant"
)

//时间戳的数据格式
const (
	TIMESTAMP_FORMAT_0 = "yyyy/MM/dd HH:mm:ss"
	TIMESTAMP_FORMAT_1 = "MM/dd/yyyy HH:mm:ss"
	TIMESTAMP_FORMAT_2 = "dd/MM/yyyy HH:mm:ss"
	TIMESTAMP_FORMAT_3 = "yyyy-MM-dd HH:mm:ss"
	TIMESTAMP_FORMAT_4 = "MM-dd-yyyy HH:mm:ss"
	TIMESTAMP_FORMAT_5 = "dd-MM-yyyy HH:mm:ss"
)

//本地缓存重发
const (
	RECORDS_RETRIES = 3
	PROPERTY_DATA_CACHE_DIR = "/data/dis"
	LINE_SEPARATOR = "\n"
	PROPERTY_DATA_CACHE_DISK_MAX_SIZE = 2048
	PROPERTY_DATA_CACHE_ARCHIVE_MAX_SIZE = 512
	PROPERTY_DATA_CACHE_ARCHIVE_LIFE_CYCLE = 24 * 3600
)