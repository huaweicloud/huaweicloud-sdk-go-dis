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

package client

import (
	"com/utils"
	"net/url"
	"strings"
	"com/models"
)

var region = "CHINA" //区域,默认CHINA
var authType = "V4"  //鉴权方式，默认V4
var rawAK string = ""
var rawSK string = ""

type Client struct {
	projectId         string                  //连接数据接入服务的用户projectId
	ak                string                  //连接对象存储服务的AK
	sk                string                  //鉴权使用的SK,可用于字符串的签名
	region            string                  //V4鉴权时使用的region
	authType          string                  //鉴权类型
	rawAK             string                  //未加密的AK
	rawSK             string                  //未加密的SK
	endpoint          string                  //服务器地址
	pathStyle         bool                    //连接请求的格式,true(路径方式)，false(子域名方式)
	bodySerializeType string                  //通信格式:默认json, protobuf
	cacheResendConf   *models.CacheResendConf //本地缓存重发相关配置项
}

/**
*函数原型：func Factory(ak, sk, endpoint string, pathStyle bool)  *Client
*函数功能：初始化Client类实例
*参数说明：AK：用户的AccessKeyID
*		 SK：用户的SecretAccessKeyID
*		 Endpoint：服务器地址，如（https://129.7.182.2:443）
*		 PathStyle：请求方式是否为绝对路径方式，取值True 或 False
*返回值：Client实例化对象
 */
func Factory(projectId, ak, sk, endpoint, bodySerializeType string, pathStyle bool) *Client {
	en_ak := string(encrypto([]byte(ak), "ak"))
	en_sk := string(encrypto([]byte(sk), "sk"))
	SetTimeOut(300, 300)
	u, err := url.Parse(endpoint)
	if err != nil {
		panic(err)
	}
	//endpoint 同时包括主机名和端口信息，如过端口存在的话，使用 strings.Split() 从 Host 中手动提取端口。
	h := strings.Split(u.Host, ":")
	endpoint = u.Scheme + "://" + h[0]
	return &Client{projectId:projectId, ak: en_ak, sk: en_sk, region: region, authType: authType, rawAK: ak, rawSK: sk, endpoint: endpoint, pathStyle: pathStyle, bodySerializeType:bodySerializeType}
}

/**
*函数原型：func FactoryEx(ak, sk, reg, auth, endpoint string, pathStyle bool) *Client
*函数功能：初始化Client类实例
*参数说明：ak：用户的AccessKeyID
*		 sk：用户的SecretAccessKeyID
*        reg:V4鉴权使用的region
*        auth:鉴权方式
*		 endpoint：服务器地址，如（https://129.7.182.2:443）
*		 pathStyle：请求方式是否为绝对路径方式，取值True 或 False
*返回值：Client实例化对象
 */
func FactoryEx(disConf *models.DISClientConf) *Client {
	en_ak := string(encrypto([]byte(disConf.AK), "ak"))
	en_sk := string(encrypto([]byte(disConf.SK), "sk"))
	SetTimeOut(300, 300)
	//u, err := url.Parse(endpoint)
	//if err != nil {
	//	panic(err)
	//}

	reg := disConf.Region
	if "" == reg {
		reg = "southchina"
	}

	auth := "V4"

	bodySerializeType := disConf.BodySerializeType
	if models.PROTOBUF != bodySerializeType {
		bodySerializeType = models.JSON
	}

	cacheResendConf := disConf.CacheResendConf
	if cacheResendConf == nil {
		cacheResendConf = models.DefaultCacheResendConf()
	} else if cacheResendConf.DataCacheDir == "" {
		cacheResendConf.DataCacheDir = models.PROPERTY_DATA_CACHE_DIR
	} else if cacheResendConf.DataCacheDiskMaxSize == 0 {
		cacheResendConf.DataCacheDiskMaxSize = models.PROPERTY_DATA_CACHE_DISK_MAX_SIZE
	} else if cacheResendConf.DataCacheArchiveSize == 0 {
		cacheResendConf.DataCacheArchiveSize = models.PROPERTY_DATA_CACHE_ARCHIVE_MAX_SIZE
	} else if cacheResendConf.DataCacheArchiveLifeCycle == 0 {
		cacheResendConf.DataCacheArchiveLifeCycle = models.PROPERTY_DATA_CACHE_ARCHIVE_LIFE_CYCLE
	}

	dis := &Client{projectId: disConf.ProjectId, ak: en_ak, sk: en_sk, region: reg, authType: auth, rawAK: disConf.AK, rawSK: disConf.SK, endpoint: disConf.Endpoint, pathStyle: true, bodySerializeType:bodySerializeType, cacheResendConf: cacheResendConf}

	if cacheResendConf.DataCacheEnable {
		go AutoResendRecordsTimeTask(dis)
	}
	return dis
}

/**
*函数原型：func SetAuthentication(auth string)
*函数功能：设置鉴权方式
*参数说明：authType鉴权类型
*返回值：无
 */
func SetAuthentication(auth string) {
	if "" != auth {
		authType = auth
	}
}

/**
*函数原型：func SetRegion(reg string)
*函数功能：设置鉴权方式
*参数说明：authType鉴权类型
*返回值：无
 */
func SetRegion(reg string) {
	if "" != reg {
		region = reg
	}
}

/**
*函数原型：func SetTimeOut(connectTime,responseTime uint)
*函数功能：设置网络超时
*参数说明：connectTime: 连接超时（秒），responseTime：响应超时
*返回值：
 */
func SetTimeOut(connectTime, responseTime uint) {
	utils.SetTimeOut(connectTime, responseTime)
}

/**
*函数说明：设置重连次数
*入参：num 重连次数
*返回值：
 */
func SetReconnectNum(num uint) {
	utils.SetReconnectNum(num)
}

/**
*函数说明：设置https证书验证
*入参：isCrt 是否验证证书，false：不验证，true：验证
*	  crtPath 证书路径，当isCrt为true是有效
*返回值：错误信息
 */
func SetCaCertificate(isCrt bool, crtPath string) error {
	return utils.SetCaCertificate(isCrt, crtPath)
}

func encrypto(plaintext []byte, refer string) []byte {
	return utils.Encrypto(plaintext, refer)
}
