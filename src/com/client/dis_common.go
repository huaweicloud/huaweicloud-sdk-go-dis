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
	"bytes"
	"com/logger"
	"com/models"
	"com/utils"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"net/http"
)

const (
	PUT = "PUT"
	POST = "POST"
	DELETE = "DELETE"
	GET = "GET"
	HEAD = "HEAD"
)

type disRequest struct {
	method      string
	projectId   string
	resourceMap [][]string
	pathMap     map[string]string //prefix,marker....
	headerMap   map[string]string //date,hots.....
}

func getRequest(method, projectId string) *disRequest {
	dis := &disRequest{method: method, projectId: projectId}
	dis.resourceMap = make([][]string, 0)
	dis.pathMap = make(map[string]string)
	dis.headerMap = make(map[string]string)
	return dis
}

func setDisURI(k, v string, dis *disRequest) {
	if k != "" {
		dis.resourceMap = append(dis.resourceMap, []string{k, v})
	}
}

//设置请求结构体中headerMap
func setDisHeader(k, v string, dis *disRequest) {
	if v != "" {
		dis.headerMap [k] = v
	}
}

func setDisPath(k, v string, dis *disRequest, isMust bool) {
	if isMust {
		dis.pathMap[k] = v
	} else {
		if v != "" {
			dis.pathMap[k] = v
		}
	}
}

func getXmlRead(xml string, util *utils.Util) *strings.Reader {
	util.BackupRequestBody([]byte(xml), "")
	return strings.NewReader(xml)
}

func readXml(r io.Reader) ([]byte, error) {
	xml, err := ioutil.ReadAll(r)
	if err != nil {
		logger.LOG(logger.ERROR, err.Error())
		return nil, err
	}
	return xml, nil
}

func (c *Client) initDis(dis *disRequest, ioread io.Reader) (*utils.Util, *models.Result) {
	util := utils.NewUtil(c.ak, c.sk, c.region, c.authType, c.rawAK, c.rawSK, c.endpoint, c.pathStyle) //初始化Util(工具)实例
	for k, v := range dis.pathMap {
		util.SetPath(k, v)
	}
	if err := util.InitConect(dis.method, dis.projectId, dis.resourceMap, ioread); err != nil {
		r := new(models.Result)
		r.Err = err
		return util, r
	}
	for k, v := range dis.headerMap {
		util.SetHeader(k, v)
	}

	util.SetHeader("Content-Type", "application/json;charset=UTF-8")
	return util, nil
}

/**
*函数原型:func (c *Client) connectDisWithXml(dis *disRequest, xml string) (*utils.Util, *models.Result)
*函数功能:链接DIS带请求消息元素
*参数说明:input: obsRequest对象实例
       : xml  : 构造请求消息元素xml
*返回值：result: Result对象实例
*/
func (c *Client) connectDisWithXml(dis *disRequest, xml string) (*utils.Util, *models.Result) {
	//构造Http请求消息长度,不包含消息体
	setDisHeader("Content-Length", strconv.Itoa(len(xml)), dis)
	u, r := c.initDis(dis, strings.NewReader(xml))
	if r == nil {
		u.BackupRequestBody([]byte(xml), "")
	}
	return u, r
}

func (c *Client) connectDis(dis *disRequest) (*utils.Util, *models.Result) {
	return c.initDis(dis, nil)
}

func (c *Client) connectDisWithFile(dis *disRequest, body []byte, filePath string) (*utils.Util, *models.Result) {
	var ioRead io.Reader = nil
	if body != nil {
		ioRead = bytes.NewReader(body)
		setDisHeader("Content-Length", strconv.Itoa(len(body)), dis)
	} else if filePath != "" {
		var err error
		var fi os.FileInfo
		result := new(models.Result)
		ioRead, err = os.Open(filePath)
		if err != nil {
			logger.LOG(logger.ERROR, err.Error())
			result.Err = err
			return nil, result
		}
		fi, err = os.Stat(filePath)
		if err != nil {
			logger.LOG(logger.ERROR, err.Error())
			result.Err = err
			return nil, result
		}
		setDisHeader("Content-Length", strconv.FormatInt(int64(fi.Size()), 10), dis)
	}
	u, r := c.initDis(dis, ioRead)
	if r == nil {
		if body != nil {
			u.BackupRequestBody(body, "")
		} else if filePath != "" {
			u.BackupRequestBody(nil, filePath)
		}
	}
	return u, r
}

func (c *Client) getCommonResponse(u *utils.Util) *models.Result {
	_, r := u.DoExec()
	u.Close()
	logger.LOG(logger.DEBUG, "exec result:statusCode:%d,code:%s,message:%s,requestId:%s,hostId:%s",
		r.StatusCode, r.Code, r.Message, r.RequestId, r.HostId)
	return r
}

/**
适用于带有"响应消息元素"的操作
*/
func (c *Client) getResponseWithOutput(u *utils.Util, out interface{}) *models.Result {
	r, _ := c.getResponseWithOutputAndStream(u, out)
	return r
}

func (c *Client) getResponseWithOutputAndStream(u *utils.Util, out interface{}) (*models.Result, *http.Response) {
	res, r := u.DoExec()
	if r.Err != nil {
		return r, nil
	}
	resposeByte, err := readXml(res.Body)
	u.Close()
	if err != nil {
		r.Err = err
		return r, nil
	}
	err = utils.ParseJson(resposeByte, out)
	if err != nil {
		r.Err = err
		return r, nil
	}
	logger.LOG(logger.DEBUG, "exec result:statusCode:%d,code:%s,message:%s,requestId:%s,hostId:%s",
		r.StatusCode, r.Code, r.Message, r.RequestId, r.HostId)
	return r, res
}

func (c *Client) getResponseWithOutputAndResposeByte(u *utils.Util) (*models.Result, []byte) {
	res, r := u.DoExec()
	if r.Err != nil {
		return r, nil
	}

	resposeByte, err := readXml(res.Body)
	u.Close()
	if err != nil {
		r.Err = err
		return r, nil
	}

	logger.LOG(logger.INFO, "exec result:statusCode:%d,code:%s,message:%s,requestId:%s,hostId:%s",
		r.StatusCode, r.Code, r.Message, r.RequestId, r.HostId)

	return r, resposeByte
}


