# Huawei Cloud DIS SDK for GO

DIS GO SDK是数据接入服务（DIS）提供的一个sdk，第三方应用程序直接调用DIS SDK中的接口即可实现获取DIS系统的业务能力。
## 1. 接口主要功能
|接口类型         |描述                            |
|----------------|--------------------------------|
|上传数据         |上传数据到DIS通道中。            |
|获取迭代器       |用户获取迭代器，根据迭代器获取一次数据和下一个迭代器。迭代器指定了读取数据记录分片位置的开始顺序。|
|下载数据         |从DIS通道中下载数据。             |

接口细节参考数据接入服务新手指南：http://support.huaweicloud.com/api-dis/zh-cn_topic_0058207012.html

## 2. 开发环境
DIS GO SDK目前支持go1.9及以上版本的开发环境。

### 2.1 Windows下GO环境的搭建：
#### 2.1.1 Go环境的安装
下载并安装go1.9.2.windows-386.msi ，可至官网https://golang.org/dl/下载对应版本的安装包。
#### 2.1.2 Go开发工具的安装
Go语言开发工具比较多，根据自己需要和习惯选择，此处以liteIDE作为示例。可至官网http://www.golangtc.com/download/liteide下载。

1. 首先下载并解压liteidex20.1.windows.7z，进入liteIDE\bin目录。
2. 双击liteIDE.exe启动。
3. 单击查看，选择环境变量。
4. 设置GOROOT为GO的安装路径。
5. 新建一个用于存放工程的目录。
6. 单击查看，选择管理GOPATH。
7. 单击添加项目，导入创建的目录。

#### 2.1.3 为golang安装protobuf，建议protoc-3.3.0及以上版本

##### 1. 下载protobuf的编译器protoc，[下载地址](https://github.com/google/protobuf/releases)

window：下载protoc-3.3.0-win32.zip解压，把bin目录下的protoc.exe复制到GOROOT/bin下，GOPATH/bin加入环境变量。

linux：下载protoc-3.3.0-linux-x86_64.zip或protoc-3.3.0-linux-x86_32.zip解压，把bin目录下的protoc复制到GOPATH/bin下，GOPATH/bin加入环境变量。

##### 2. 获取protobuf的编译器插件protoc-gen-go

进入GOPATH目录运行: go get -u github.com/golang/protobuf/protoc-gen-go， 如果成功会在GOPATH/bin下生成protoc-gen-go.exe文件。

#### 2.1.4 使用代理

##### 1. liteIDE下设置代理服务器, 点击`工具`-`编辑当前环境`, 添加如下两行
``` shell
http_proxy="http://{username}:{password}@{host}:{port}"
https_proxy="http://{username}:{password}@{host}:{port}"
# {username}, {password}, {host}, {port} 都需要使用URL编码
```

##### 2. windows设置代理服务器
``` shell
set http_proxy=http://{username}:{password}@{host}:{port}
set https_proxy=http://{username}:{password}@{host}:{port}
```

##### 3. linux设置代理服务器
``` shell
export http_proxy=http://{username}:{password}@{host}:{port}
export https_proxy=http://{username}:{password}@{host}:{port}
```

