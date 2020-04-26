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
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"
)

/**
*函数说明：构造v2鉴权的stringToSign值
*入参：util:Util实例
*返回值：构造v2鉴权的stringToSign值
 */
func stringToSignS3(util *Util) string {
	str := util.request.Method + "\n"

	if util.request.Header.Get("Content-Md5") != "" {
		str += util.request.Header.Get("Content-Md5")
	}
	str += "\n"

	str += util.request.Header.Get("Content-Type") + "\n"

	if util.request.Header.Get("Date") != "" {
		str += util.request.Header.Get("Date")
	} else {
		str += timestampS3()
	}
	str += "\n"

	canonicalHeaders := canonicalAmzHeadersS3(util)
	if canonicalHeaders != "" {
		str += canonicalHeaders
	}

	str += canonicalResourceS3(util)

	return str
}

/**
*函数说明：构造v2鉴权的canonicalAmzHeaders值
*入参：util:Util实例
*返回值：构造v2鉴权的canonicalAmzHeaders值
 */
func canonicalAmzHeadersS3(util *Util) string {
	var headers []string

	for header := range util.request.Header {
		//TrimSpace返回将目标字符串前后所有空白字符去掉的字符串
		//ToLower将所有字幕转换成对应的小写版本
		//HasPrefix判断s是否有前缀字符串prefix
		standardized := strings.ToLower(strings.TrimSpace(header))
		if strings.HasPrefix(standardized, "x-amz") {
			headers = append(headers, standardized)
		}
	}

	sort.Strings(headers)

	for i, header := range headers {
		headers[i] = header + ":" + strings.Replace(util.request.Header.Get(header), "\n", " ", -1)
	}

	if len(headers) > 0 {
		return strings.Join(headers, "\n") + "\n"
	} else {
		return ""
	}
}

/**
*函数说明：构造v2鉴权的canonicalResource值
*入参：util:Util实例
*返回值：构造v2鉴权的canonicalResource值
 */
func canonicalResourceS3(util *Util) string {
	res := "/"
	if util.vesion != "" {
		res += "/" + util.vesion
	}
	if util.projectId != "" {
		res += util.projectId
	}

	for _, v := range util.resourceMap {
		if v[0] != "" {
			res += "/" + v[0]
		}
		if v[1] != "" {
			res += "/" + v[1]
		}
	}

	i := 0
	for _, subres := range strings.Split(SUBRESOURCES_S3, ",") {
		val, ok := util.pathMap[subres]
		if ok {
			if i == 0 {
				res += "?" + subres
				i = 1
			} else {
				res += "&" + subres
			}
			if val != "" {
				res += "=" + val
			}
		}
	}

	return res
}

/**
*函数说明：获取当前RFC时间
*入参：
*返回值：当前RFC时间
 */
func timestampS3() string {
	return now().Format(TIME_FORMAT_S3)
}

/**
*函数说明：hashSHA256加密
*入参：key:加密的key值
*     content加密的字符串
*返回值：加密后的值
 */
func hashSHA256(content []byte) string {
	h := sha256.New()
	h.Write(content)
	return fmt.Sprintf("%x", h.Sum(nil))
}

/**
*函数说明：字符串md5加密
*入参：content加密的字符串
*返回值：加密后的值
 */
func HashMD5(content []byte) string {
	h := md5.New()
	h.Write(content)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

/**
*函数说明：文件md5加密
*入参：filePath文件全路径
*返回值：加密后的值，错误信息
 */
func HashFileMD5(filePath string) (string, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0755)
	if err == nil {
		md5h := md5.New()
		io.Copy(md5h, file)
		fileMd5 := base64.StdEncoding.EncodeToString(md5h.Sum(nil))
		return fileMd5, nil
	}
	return "", err
}

/**
*函数说明：获取当前UTC时间
*返回值：UTC时间
 */
var now = func() time.Time {
	return time.Now().UTC()
}

const (
	TIME_FORMAT_S3 = time.RFC3339
	SUBRESOURCES_S3 = "acl,quota,storageinfo,deletebucket,delete,lifecycle,location,logging,notification,partNumber," +
		"policy,requestPayment,torrent,uploadId,uploads,versionId,versioning,versions,website"
)