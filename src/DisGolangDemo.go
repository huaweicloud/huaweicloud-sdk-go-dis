package main

import (
	"com/client"
	"com/models"
	"com/logger"
	"com/demo"
	"time"
)

var dis *client.Client = nil

func main() {
	//初始化日志，如果不想记录日志，该操作省略
	logger.InitLog("D:/SDK/dis-go", "dis_go", logger.INFO, logger.RECORD_CONSOLE)
	client.SetTimeOut(10, 30)
	client.SetReconnectNum(3)

	disConf := models.DefaultDISClientConf();

	//必选配置
	disConf.AK = "7KWMSJR4GGH666GONBKM" //用户Access Key ID，可在公有云“我的凭证”页面下载生成
	disConf.SK = "bECrDtlEGy6Gi3uSuIlTPkcnHALoZMg4FEcJJFvA" //用户Secret Access Key，可在公有云“我的凭证”页面下载生成
	disConf.ProjectId = "43eec61d4b514f359c97008a4f8bfb02" //用户ProjectID
	disConf.Region = "southchina" //中国华北区1，默认不用修改
	disConf.Endpoint = "https://10.61.121.201:8443"

	//可选配置
	disConf.BodySerializeType = models.JSON; //默认JSON
	cacheResendConf := &models.CacheResendConf{DataCacheEnable:true,
		DataCacheDir: "D:\\tmp\\dis1",
		DataCacheDiskMaxSize: 10,
		DataCacheArchiveSize:1,
		DataCacheArchiveLifeCycle:60, }
	disConf.CacheResendConf = cacheResendConf        //本地缓存重发

	dis = client.FactoryEx(disConf) //初始化服务类工厂

	//初始化测试用例的输入参数
	disInstance := demo.DISInstance{
		Dis: dis,
		StreamName:"feihang",
		PartitionId : "0",
		CursorType:models.AT_SEQUENCE_NUMBER,
		SequenceNumber:"0",
		Timestamp:1534142781293,
		AppName:"feihang1",
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
