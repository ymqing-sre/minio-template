package utils

import (
	"context"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/viper"
	"strings"
)

func MakeBucket(c context.Context, api S3CreateBucketAPI, input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {
	return api.CreateBucket(c, input)
}

func GetAllBuckets(c context.Context, api S3ListBucketsAPI, input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
	return api.ListBuckets(c, input)
}

func PutFile(c context.Context, api S3PutObjectAPI, input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return api.PutObject(c, input)
}

func GetFile(c context.Context, api S3GetObjectAPI, input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return api.GetObject(c, input)
}

func GetObjects(c context.Context, api S3ListObjectsAPI, input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	return api.ListObjectsV2(c, input)
}

func RemoveBucket(c context.Context, api S3DeleteBucketAPI, input *s3.DeleteBucketInput) (*s3.DeleteBucketOutput, error) {
	return api.DeleteBucket(c, input)
}

func DeleteObj(c context.Context, api S3DeleteObjectAPI, input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	return api.DeleteObject(c, input)
}

func GetPresignedURL(c context.Context, api S3PresignGetObjectAPI, input *s3.GetObjectInput) (*v4.PresignedHTTPRequest, error) {
	return api.PresignGetObject(c, input)
}

func ObjectExists(c context.Context, api S3ObjectExistsAPI, input *s3.HeadObjectInput) bool {
	_, err := api.HeadObject(c, input)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return false
		}
		panic(err)
	} else {
		return true
	}
}

func BucketExists(c context.Context, api S3BucketExistsAPI, input *s3.HeadBucketInput) bool {
	_, err := api.HeadBucket(c, input)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return false
		}
		panic(err)
	} else {
		return true
	}
}

func CopyItem(c context.Context, api S3CopyObjectAPI, input *s3.CopyObjectInput) (*s3.CopyObjectOutput, error) {
	return api.CopyObject(c, input)
}

func LoadConfig(path string, configName string) (endpointURL string, key string, secret string) {
	config := viper.New()
	config.AddConfigPath(path)
	config.SetConfigName(configName)
	config.SetConfigType("yaml")
	err := config.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("找不到配置文件")
		} else {
			panic(err)
		}
	}
	endpointURL = config.GetString("endpointURL")
	key = config.GetString("credentials.accessKey")
	secret = config.GetString("credentials.secretAccessKey")
	return endpointURL, key, secret
}
