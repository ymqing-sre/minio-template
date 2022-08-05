# minio-template
基于aws-sdk-go-v2版本。
适用于Ceph、Minio等支持S3协议的对象存储。
## Build：
```shell
go build -o mc main.go
```
## 配置
配置文件`config.yaml`在可执行文件同级目录的`./conf`目录下
```yaml
# 对象存储的endpoint地址
endpointURL: "http://fs.minio.cloud"
# 凭证信息
credentials:
  accessKey: "minio"
  secretAccessKey: "Minio123456"
```
## 命令
传递动作命令`./mc -a <ACTION>` 指示需要进行的对象存储相关操作：

Action:
- listBuckets:列出对象存储上所有桶
- createBucket：创建桶
- copyObj：复制对象
- putFile：上传文件
- getFile：下载文件到本地
- deleteObj：删除对象
- listFiles：列出指定桶中的所有对象
- removeBucket：删除桶
- getPresignedURL：获取对象的预签名URL
- objExists：判断某个桶中的对象是否存在
- bktExists：判断某个桶是否存在

不同命令，需要接收不同传参。
### 桶操作
#### 列出桶
```shell
./mc -a listBuckets
```
#### 创建桶
```shell
./mc -a createBucket -b BUCKET_NAME
```
#### 删除桶
```shell
./mc -a removeBucket -b BUCKET_NAME
```
#### 判断桶是否存在
```shell
./mc -a bktExists -b BUCKET_NAME
```

### 对象操作
#### 列出指定桶中对象
```shell
./mc -a listFiles -b BUCKET_NAME
```
#### 上传对象
```shell
./mc -a putFile -b BUCKET_NAME -p PATH -f FILE
```
`-p PATH`可设定文件上传到桶下的路径，不设置即上传文件到根路径下。
> Note:
> 上传后，对象所在的路径和文件名称构成了这个对象的KEY值。
> 
> 例如将名为testfile的文件上传到，某个桶中的a文件夹下，这个对象的KEY值为a/testfile,代表这个文件在桶中的唯一位置。如果重复上传，桶中的对象会被后上传的文件覆盖。
#### 下载对象
```shell
./mc -a getFile -b BUCKET_NAME -o OBJECT -f FILENAME
```
OBJECT为对象的KEY值，FILENAME为下载到本地后的文件名。 下载目录为当前目录。
#### 复制对象
用于将对象存储上的对象在桶内或桶间复制
```shell
./mc -a copyObj -s SRC_BUCKET_NAME -so SRC_OBJECT -d DST_BUCKET_NAME -do DST_OBJECT
```
- SRC_BUCKET_NAME:被复制的对象所在的桶 
- SRC_OBJECT：被复制的对象KEY值 
- DST_BUCKET_NAME：目标桶 
- DST_OBJECT：复制后在目标桶的KEY值
#### 删除对象
```shell
./mc -a deleteObj -b BUCKET_NAME -o OBJECT
```
#### 获取对象的预签名URL
获取预签名URL后返回一个下载地址，可通过下载地址直接下载文件。
```shell
./mc -a getPresignedURL -b BUCKET_NAME -o OBJECT
```
#### 判断对象是否存在
```shell
./mc -a objExists -b BUCKET_NAME -o OBJECT
```