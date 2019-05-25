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
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"sort"
	"strings"
	"time"
)

/**
*函数说明：构造v4鉴权的signature值
*入参：util: Util对象实例
*返回值：v4鉴权的signature值
 */
func signatureS3_v4(util *Util) string {
	v4Auth := "SDK-HMAC-SHA256"

	credenttial := getCredenttial(util)
	v4Auth += " Credential=" + credenttial

	signedHeaders := util.getSignedHeaders()
	v4Auth += ", SignedHeaders=" + signedHeaders

	reqStr := canonicalRequest(util)
	signature := hmacSHA256([]byte(signingKey(util)), stringToSign(reqStr, util))
	v4Auth += (", Signature=" + hex.EncodeToString([]byte(signature)))
	//fmt.Printf("v4Auth:%s\n", v4Auth)
	return v4Auth
}

/**
*函数说明：构造v4鉴权的canonicalRequest值
*入参：util: Util对象实例
*返回值：v4鉴权的CanonicalRequest值
 */
func canonicalRequest(util *Util) string {
	str := util.request.Method
	uri := util.getCanonicalURI()
	str += "\n" + uri + "\n"

	queryStr := util.getCanonicalQueryString()
	str += queryStr + "\n"

	headers := util.getCanonicalHeaders()
	str += headers + "\n"

	signHeader := util.getSignedHeaders()
	str += signHeader + "\n"

	contentSha256 := calculateContentHash(util)
	str += contentSha256
	return str
}

/**
*函数说明：构造v4鉴权的calculateContentHash值
*入参：util: Util对象实例
*返回值：v4鉴权的calculateContentHash值
 */
func calculateContentHash(util *Util) string {
	str := hashSHA256(util.requestBody)
	//str := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	return str
}

/**
*函数说明：构造v4鉴权的StringToSign值
*入参：util: Util对象实例
*返回值：v4鉴权的StringToSign值
 */
func stringToSign(reqStr string, util *Util) string {
	str := "SDK-HMAC-SHA256" + "\n"

	date := util.request.Header.Get("x-sdk-date")
	str += date + "\n" + getScope(date, util) + "\n"
	//str += "20170820T015124Z" + "\n" + getScope("20170820T015124Z", util) + "\n"

	//strSha256 := hmacSHA256([]byte(reqStr), "")
	//fmt.Printf("reqStr:%s\n", reqStr)
	strSha256 := hashSHA256([]byte(reqStr))
	//str += hex.EncodeToString([]byte(strSha256))
	str += strSha256
	//fmt.Printf("str:%s\n", str)
	return str
}

/**
*函数说明：构造v4鉴权的SigningKey值
*入参：util: Util对象实例
*返回值：v4鉴权的SigningKey值
 */
func signingKey(util *Util) string {
	//date := getTimesISO8601()
	//date = date[0:8]
	date := util.request.Header.Get("x-sdk-date")[0:8]

	datekey := hmacSHA256([]byte("SDK" + util.rawSK), date)
	dateRegionKey := hmacSHA256(datekey, util.region)

	dateRegionServiceKey := hmacSHA256(dateRegionKey, "dis")
	signingKey := hmacSHA256(dateRegionServiceKey, "sdk_request")
	return string(signingKey)
}

/**
*函数说明：构造v4鉴权的Scope值
*入参：util: Util对象实例
*返回值：v4鉴权的Scope值
 */
func getScope(time string, util *Util) string {
	return time[0:8] + "/" + util.region + "/dis/sdk_request"
}

/**
*函数说明：构造v4鉴权的Credential值
*入参：util: Util对象实例
*返回值：v4鉴权的Credential值
 */
func getCredenttial(util *Util) string {
	date := util.request.Header.Get("x-sdk-date")
	credenttial := util.rawAK + "/" + getScope(date, util)
	return credenttial
}

/**
*函数说明：构造v4鉴权的SignedHeaders值
*入参：util: Util对象实例
*返回值：v4鉴权的SignedHeaders值
 */
func (util *Util) getSignedHeaders() string {
	var headers string
	/*
		for k, _ := range util.request.Header {
			headers += (strings.ToLower(k) + ";")
		}
	*/

	keys := make([]string, len(util.request.Header))
	i := 0
	for k, _ := range util.request.Header {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, k := range keys {
		headers += (strings.ToLower(k) + ";")
	}

	//headers += "x-amz-content-sha256': 'e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	return headers[0 : len(headers) - 1]
}

/**
*函数说明：构造v4鉴权的Canonical Headers值
*入参：util: Util对象实例
*返回值：v4鉴权的Canonical Headers值
 */
func (util *Util) getCanonicalHeaders() string {
	var headers string
	/*
		for k, v := range util.request.Header {
			headers += string(strings.ToLower(k) + ":" + strings.Join(v, "") + "\n")
		}
	*/
	keys := make([]string, len(util.request.Header))
	i := 0
	for k, _ := range util.request.Header {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, k := range keys {
		headers += string(strings.ToLower(k) + ":" + strings.Join(util.request.Header[k], "") + "\n")
	}
	return headers
}

/**
*函数说明：构造v4鉴权的Canonical URI值
*入参：
*返回值：获取URI值
 */
func (util *Util) getCanonicalURI() string {
	uri := ""
	if util.vesion != "" {
		uri += "/" + util.vesion
	}
	if util.projectId != "" {
		uri += "/" + util.projectId
	}

	for _, v := range util.resourceMap {
		if v[0] != "" {
			uri += "/" +v[0]
		}
		if v[1] != "" {
			uri += "/" + v[1]
		}
	}

	uri += "/"
	/*
		i := 0
		for key, value := range util.pathMap {
			if i == 0 {
				uri += "?" + key
				i = 1
			} else {
				uri += "&" + key
			}
			if value != "" {
				uri += "=" + value
			}
		}
	*/
	return uri

}

/**
*函数说明：构造v4鉴权的CanonicalQueryString值
*入参：
*返回值：v4鉴权的CanonicalQueryString值
 */
func (util *Util) getCanonicalQueryString() string {
	var queryStr string
	/*
		for key, value := range util.pathMap {
			if value != "" {
				queryStr += url.QueryEscape(key) + "=" + url.QueryEscape(value) + "&"
			} else {
				queryStr += url.QueryEscape(key) + "=" + "" + "&"
			}
		}
	*/

	keys := make([]string, len(util.pathMap))
	i := 0
	for k, _ := range util.pathMap {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, k := range keys {
		if util.pathMap[k] != "" {
			queryStr += url.QueryEscape(k) + "=" + url.QueryEscape(util.pathMap[k]) + "&"
		} else {
			queryStr += url.QueryEscape(k) + "=" + "" + "&"
		}
	}

	if len(queryStr) > 0 {
		return queryStr[0 : len(queryStr) - 1]
	} else {
		return ""
	}
}

/**
*函数说明：获取当前ISO8601格式时间
*入参：
*返回值：当前ISO8601时间
 */
func getTimesISO8601() string {
	str1 := time.Now().UTC().Format(time.RFC3339)
	str2 := strings.Replace(str1, "-", "", -1)
	time8601 := strings.Replace(str2, ":", "", -1)
	return time8601
}

/**
*函数说明：hmacSHA256加密
*入参：key:加密的key值
*     content加密的字符串
*返回值：加密后的值
 */
func hmacSHA256(key []byte, content string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(content))
	return mac.Sum(nil)
}
