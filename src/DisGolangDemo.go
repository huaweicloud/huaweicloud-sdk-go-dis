package main

import (
	"com/client"
	"com/demo"
	"com/logger"
	"com/models"
	"time"
)

var dis *client.Client = nil

func main() {
	// 初始化日志，如果不想记录日志，该操作省略
	logger.InitLog("D:/SDK/dis-go", "dis_go", logger.INFO, logger.RECORD_CONSOLE)
	client.SetTimeOut(10, 30)
	client.SetReconnectNum(3)

	disConf := models.DefaultDISClientConf()

	// 必选配置
	disConf.AK = "YOUR_AK"                                        // 用户Access Key ID，可在公有云“我的凭证”页面下载生成
	disConf.SK = "YOUR_SK"                                        // 用户Secret Access Key，可在公有云“我的凭证”页面下载生成
	disConf.ProjectId = "YOUR_PROJECT_ID"                         // 用户ProjectID
	disConf.Region = "cn-north-1"                                 // 中国华北区北京一，根据通道所属Region选择
	disConf.Endpoint = "https://dis.cn-north-1.myhuaweicloud.com" // 根据通道所属Region选择Endpoint

	dis = client.FactoryEx(disConf) // 初始化客户端

	// 初始化测试用例的输入参数
	disInstance := demo.DISInstance{
		Dis:            dis,
		StreamName:     "YOUR_STREAM_NAME", // 通道名称
		StreamId:       "",                 // 通道ID，访问授权通道时需要指定通道ID
		PartitionId:    "0",
		CursorType:     models.AT_SEQUENCE_NUMBER, // 游标类型，用于确定从分区的什么位置开始获取数据
		SequenceNumber: "0",
		Timestamp:      1534142781293,
		AppName:        "YOUR_APP_NAME",
	}

	//demo.StreamDemo(disInstance) //通道的创建删除等
	//demo.DeleteStream(disInstance)
	//demo.CreateStream(disInstance)
	//demo.CreateTransferTask(disInstance)

	for i := 0; i < 1; i++ {
		getInfo(disInstance) //上传数据
	}
	//c = make(chan int)
	//go getInfo(disInstance)//上传数据
	//go getInfo(disInstance) //上传数据
	//go getInfo(disInstance) //上传数据
	//go getInfo(disInstance) //上传数据
	//<-c
	//<-c
	//<-c
	//<-c
	//demo.DeleteStream(disInstance)
	//demo.GetRecords(disInstance) //下载数据
	//
	//demo.AppDemo(disInstance); //APP的创建删除等
	//
	//demo.PutRecordsAndCheckpoint(disInstance) //上传数据并checkpoint
	//demo.GetRecordsWithCheckpoint(disInstance) //获取checkpoint并下载数据

	for ; ; time.Sleep(1 * time.Second) {

	}
	//demo.GetRecords(disInstance) //下载数据
	//
	//demo.AppDemo(disInstance); //APP的创建删除等
	//
	//demo.PutRecordsAndCheckpoint(disInstance) //上传数据并checkpoint
	//demo.GetRecordsWithCheckpoint(disInstance) //获取checkpoint并下载数据
}

var c chan int

func getInfo(disInstance demo.DISInstance) {
	demo.PutRecords(disInstance)
	//c <- 0
}
