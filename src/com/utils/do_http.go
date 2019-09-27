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

package utils

import (
	"bytes"
	"com/logger"
	"com/models"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	V2 = "v2"
)

type Util struct {
	pathMap           map[string]string
	request           *http.Request
	Response          *http.Response
	vesion            string
	projectId         string
	resourceMap       [][]string
	ak                string
	sk                string
	rawAK             string
	rawSK             string
	region            string
	authType          string
	endpoint          string
	pathStyle         bool
	requestBody       []byte
	requestSourceFile string
}

var connectTimeOut uint = 30  //连接超时时间,默认30s
var responseTimeOut uint = 30 //返回超时时间,默认30s
var reConnectNum uint = 3     //网络重连次数,默认3次
var caCrt []byte = nil        //https数字证书

/**
*函数说明：初始化Util实例
*入参：AK：用户的AccessKeyID
*	  SK：用户的SecretAccessKeyID
*     Region: V4鉴权使用的region
*     AuthType: 鉴权类型
*	  Endpoint：服务器地址，如（https://129.7.182.2:443）
*     PathStyle：请求方式是否为绝对路径方式，取值True 或 False
*返回值：初始化的Util实例
 */
func NewUtil(ak, sk, reg, auth, raw_ak, raw_sk, endpoint string, pathStyle bool) *Util {
	util := &Util{ak: ak, sk: sk, region: reg, authType: auth, rawAK: raw_ak, rawSK: raw_sk, endpoint: endpoint, pathStyle: pathStyle}
	util.pathMap = make(map[string]string)
	return util
}

/**
*函数说明：初始化连接
*入参：mothed：请求方法
*	  bucket：桶名
*	  object：对象名
*     ioread：待发送的数据流
*返回值：执行失败值
 */
func (util *Util) InitConect(mothed, projectId string, resourceMap [][]string, ioread io.Reader) error {
	murl, er := url.Parse(V2)
	if er == nil {
		util.vesion = murl.String()
	} else {
		util.vesion = V2
	}

	murl, er = url.Parse(projectId)
	if er == nil {
		util.projectId = murl.String()
	} else {
		util.projectId = projectId
	}

	util.resourceMap = resourceMap
	erro, path := util.getPath()
	if erro != nil {
		return erro
	}
	var err error
	util.request, err = http.NewRequest(mothed, path, ioread)
	if err != nil {
		logger.LOG(logger.ERROR, err.Error())
		return err
	}

	host := strings.Split(strings.Split(util.endpoint, "//")[1], "/")[0]

	util.request.Header.Set("host", host)
	util.request.Header.Set("x-sdk-date", getTimesISO8601())
	//util.request.Header.Set("Date", "Sun, 20 Aug 2017 01:51:24 GMT")
	if "V4" == util.authType {
		util.request.Header.Set("x-sdk-content-sha256", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	}
	return nil
}

/**
*函数说明：向服务端发送请求
*入参：
*返回值：http.Response实例，Result执行结构信息
 */
func (util *Util) DoExec() (*http.Response, *models.Result) {
	if util.request.Header.Get("Authorization") == "" {
		if "V2" == util.authType {
			util.request.Header.Set("Authorization", signatureS3(util))
		}
		if "V4" == util.authType {
			util.request.Header.Set("Authorization", signatureS3_v4(util))
		}
	}
	sourceRange := util.request.Header.Get("X-Amz-Copy-Source-Range")
	if "" != sourceRange {
		util.request.Header["x-amz-copy-source-range"] = []string{sourceRange}
		util.request.Header.Del("X-Amz-Copy-Source-Range")
	}
	b, err := strconv.ParseInt(util.request.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		//fmt.Println("string to int64 failed.")
	}
	util.request.ContentLength = b

	result := &models.Result{}
	if err := util.doHttp(); err != nil {
		logger.LOG(logger.ERROR, err.Error())
		result.Err = err
		return nil, result
	}
	if util.Response == nil {
		logger.LOG(logger.ERROR, "Failed to send request with response")
		result.Err = errors.New("Failed to send request with response")
		return nil, result
	}
	ret := util.getCommonErrResponse()
	return util.Response, ret
}

/**
*函数说明：设置请求URI
*入参：key：uri对应key值
*     value：uri的key对应值
*返回值：
 */
func (util *Util) SetPath(key, value string) {
	mapKey := key
	mapValue := value
	murl, er := url.Parse(key)
	if er == nil {
		mapKey = murl.String()
		mapKey = strings.Replace(mapKey, "%20", " ", -1)
	}
	murl, er = url.Parse(value)
	if er == nil {
		mapValue = murl.String()
		mapValue = strings.Replace(mapValue, "%20", " ", -1)
	}
	util.pathMap[mapKey] = mapValue
}

/**
*函数说明：设置请求头域
*入参：key：head头对应key值
*     value：head头的key对应值
*返回值：
 */
func (util *Util) SetHeader(key string, value string) {
	headKey := key
	headValue := value
	murl, er := url.Parse(key)
	if er == nil {
		headKey = murl.String()
		headKey = strings.Replace(headKey, "%20", " ", -1)
	}
	murl, er = url.Parse(value)
	if er == nil {
		headValue = murl.String()
		headValue = strings.Replace(headValue, "%20", " ", -1)
	}
	util.request.Header.Set(headKey, headValue)
}

/**
*函数说明：关闭http数据流
*入参：
*返回值：
 */
func (util *Util) Close() {
	if util.Response != nil {
		defer util.Response.Body.Close()
	}
}

/**
*函数说明：备份request请求消息体，重连时需要
*入参：body：请求消息体，和sourceFile不能同时存在
*     sourceFile:上传的文件路径，和body不能同时存在
*返回值：获取URI值
 */
func (util *Util) BackupRequestBody(body []byte, sourceFile string) {
	util.requestBody = body
	util.requestSourceFile = sourceFile
}

/**
*函数说明：获取URI
*入参：
*返回值：获取URI值
 */
func (util *Util) getPath() (error, string) {
	ser := strings.Split(util.endpoint, "//")
	if len(ser) != 2 || (ser[0] != "https:" && ser[0] != "http:") {
		logger.LOG(logger.ERROR, "the server address is err:"+util.endpoint)
		err := errors.New("the server address is err:" + util.endpoint)
		return err, ""
	}
	path := ser[0] + "//" + ser[1]
	if util.vesion != "" {
		path += "/" + util.vesion
	}
	if util.projectId != "" {
		path += "/" + util.projectId
	}

	for _, v := range util.resourceMap {
		if v[0] != "" {
			path += "/" + v[0]
		}
		if v[1] != "" {
			path += "/" + v[1]
		}
	}

	i := 0
	for key, value := range util.pathMap {
		if i == 0 {
			path += "?" + key
			i = 1
		} else {
			path += "&" + key
		}
		if value != "" {
			path += "=" + value
		}
	}
	return nil, path
}

func (util *Util) refreshRequest() error {
	var ioread io.Reader = nil
	var err error = nil
	if util.requestBody != nil {
		ioread = bytes.NewReader(util.requestBody)
	} else if util.requestSourceFile != "" {
		ioread, err = os.Open(util.requestSourceFile)
		if err != nil {
			return err
		}
	}
	method := util.request.Method
	_, path := util.getPath()

	headers := util.request.Header
	util.request, err = http.NewRequest(method, path, ioread)
	util.request.Header = headers
	util.request.Header.Del("x-sdk-date")
	util.request.Header.Del("Authorization")
	util.request.Header.Set("x-sdk-date", getTimesISO8601())
	util.request.Header.Set("Authorization", signatureS3_v4(util))
	return err
}

/**
*函数说明：http连接设置
*入参：
*返回值：
 */
func (util *Util) getTransport() *http.Transport {
	dial := func(netw, addr string) (net.Conn, error) {
		con, err := net.DialTimeout(netw, addr, time.Second*time.Duration(connectTimeOut))
		if err != nil {
			return nil, err
		}
		tcp_conn := con.(*net.TCPConn)
		tcp_conn.SetKeepAlive(false)
		return tcp_conn, nil
	}

	if strings.Split(util.endpoint, ":")[0] == "https" {
		skipVerify := true
		if caCrt != nil {
			skipVerify = false
		}
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(caCrt)

		transport := &http.Transport{
			TLSClientConfig:       &tls.Config{RootCAs: pool, InsecureSkipVerify: skipVerify},
			DisableCompression:    true,
			Dial:                  dial,
			ResponseHeaderTimeout: time.Second * time.Duration(responseTimeOut),
		}

		if os.Getenv("https_proxy") != "" {
			proxyUrl, _ := url.Parse(os.Getenv("https_proxy"))
			transport.Proxy = http.ProxyURL(proxyUrl)
		}
		return transport

	} else {
		transport := &http.Transport{
			Dial:                  dial,
			ResponseHeaderTimeout: time.Second * time.Duration(responseTimeOut),
		}
		if os.Getenv("http_proxy") != "" {

			proxyUrl, _ := url.Parse(os.Getenv("http_proxy"))
			transport.Proxy = http.ProxyURL(proxyUrl)
		}
		return transport
	}
}

func (util *Util) doHttp() error {
	var conn *http.Client = &http.Client{Transport: util.getTransport(), Timeout: 60 * time.Second}
	var err error = nil
	_, path := util.getPath()

	request_headers := make(map[string][]string)
	for k, v := range util.request.Header {
		request_headers[k] = v
	}
	logger.LOG(logger.DEBUG, "request header msg:method:%s,path:%s",
		util.request.Method, path)

	for i := 0; i < int(reConnectNum); i++ {
		util.Response, err = conn.Do(util.request)
		if err == nil && util.Response.StatusCode < 500 {
			break
		}
		if i == int(reConnectNum)-1 {
			break
		}
		if err != nil {
			logger.LOG(logger.ERROR, "Failed to send request,err:%v,send again", err)
		} else if util.Response.StatusCode >= 500 {
			logger.LOG(logger.ERROR, "Failed to send request with response: %d", util.Response.StatusCode)
		}
		err = util.refreshRequest()
		if err != nil {
			logger.LOG(logger.ERROR, "Failed to refeshRequest with err: %v", err)
		}
		time.Sleep(time.Duration(2 * rand.Float64() * float64(time.Second)))
	}
	return err
}

func (util *Util) getCommonErrResponse() *models.Result {
	result := new(models.Result)
	statusCode := util.Response.StatusCode
	result.StatusCode = statusCode
	result.RequestId = util.Response.Header.Get("X-Amz-Request-Id")
	logger.LOG(logger.DEBUG, "response header msg:statusCode:%d,header:%v",
		result.StatusCode, util.Response.Header)

	if statusCode >= 300 {
		body, err := ioutil.ReadAll(util.Response.Body)
		if err != nil {
			logger.LOG(logger.ERROR, "response msg:statusCode:%d", result.StatusCode, err.Error())
			result.Err = err
			return result
		} else {
			errMsg := string(body)
			result.Err = errors.New(errMsg)
			logger.LOG(logger.ERROR, "response msg:statusCode:%d, errMsg:%s", result.StatusCode, errMsg)
			if statusCode == 400 {
				errResponse := new(models.ErrResponse)
				ParseJson(body, errResponse)
				result.ErrResponse = *errResponse
			}
			return result
		}
		//err = ParseJson(body, result)
		//if err != nil {
		//	logger.LOG(logger.ERROR, err.Error())
		//	result.Err = err
		//	return result
		//}
	}

	return result
}

/**
*函数说明：解析json
*入参：body:待解析的字符串
* 	 obj：解析的格式，传入对象指针
*返回值：
 */
func ParseJson(body []byte, obj interface{}) error {
	if len(body) == 0 {
		return nil
	}
	logger.LOG(logger.DEBUG, "receive msg:%s", string(body))
	err := json.Unmarshal(body, obj)
	if err != nil {
		logger.LOG(logger.ERROR, err.Error())
	}
	return err
}

/**
*函数说明：设置网络超时
*入参：second：超时时间，秒
*返回值：
 */
func SetTimeOut(connectTime, responseTime uint) {
	connectTimeOut = connectTime
	responseTimeOut = responseTime
}

/**
*函数说明：设置重连次数
*入参：num 重连次数
*返回值：
 */
func SetReconnectNum(num uint) {
	reConnectNum = num
}

/**
*函数说明：设置https数字证书
*入参：isCrt 是否验证证书，false：不验证，true：验证
*	  crtPath 证书路径，当isCrt为true是有效
*返回值：错误信息
 */
func SetCaCertificate(isCrt bool, crtPath string) error {
	if !isCrt {
		caCrt = nil
		return nil
	} else {
		var err error = nil
		caCrt, err = ioutil.ReadFile(crtPath)
		if err != nil {
			caCrt = nil
			logger.LOG(logger.ERROR, "read the SSL certificate failed,error:"+err.Error())
			return err
		}
	}
	return nil
}
