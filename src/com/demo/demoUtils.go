package demo

import "com/client"


//dis测试用例的输入参数
type DISInstance struct {
	Dis            *client.Client //初始化服务类工厂

	StreamId       string 		  // 通道ID
	StreamName     string         // 配置流名称
	PartitionId    string         //获取迭代器时需指定分区ID

	CursorType     string         //下载数据方式
				      // AT_SEQUENCE_NUMBER 从指定的sequenceNumber开始获取，需要设置StartingSequenceNumber
				      // AFTER_SEQUENCE_NUMBER 从指定的sequenceNumber之后开始获取，需要设置StartingSequenceNumber
				      // TRIM_HORIZON 从最旧的记录开始获取
				      // LATEST 从最新的记录开始获取
				      // AT_TIMESTAMP 从指定的时间戳(13位)开始获取，需要设置Timestamp
	SequenceNumber string         //cursorType := models.AT_SEQUENCE_NUMBER时需指定下载数据序列号
	Timestamp      int64          //cursorType := models.AT_TIMESTAMP时需指定下载数据的时间戳

	AppName        string         //设置app名称，APP用于管理checkpoint
}
