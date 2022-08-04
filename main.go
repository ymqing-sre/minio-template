package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ymqing-sre/minio-template/utils"
	"os"
)

// createBucket,创建桶
/*func createBucket(client *s3.Client, bktName string) {
	bktInput := s3.CreateBucketInput{
		Bucket: &bktName,
	}
	_, err := client.CreateBucket(context.TODO(), &bktInput)
	if err != nil {
		panic(err)
	}
}*/

// listBuckets，列出桶
/*func listBuckets(client *s3.Client) {
	out, err := client.ListBuckets(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Print("[")
	for i := range out.Buckets {
		if i > 0 {
			fmt.Print(",")
		}
		bkt := out.Buckets[i]
		fmt.Print(*bkt.Name)
	}
	fmt.Print("]\n")
}*/

func main() {
	// 定义Minio API地址
	endpointURL := "http://fs.quanxiang.cloud:30145"
	// 定义EndpointResolver
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: endpointURL,
		}, nil
	})
	// 定义CredentialsProvider
	customCredentials := credentials.NewStaticCredentialsProvider("minio", "Minio123456", "")
	// 加载配置
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(customCredentials),
	)
	if err != nil {
		panic(err)
	}
	// 创建S3客户端
	client := s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.UsePathStyle = true
	})
	bktName := "minio-test"
	filepath := "ababab"
	filename := "testfile"
	keyname := filepath + "/" + filename

	// 创建Bucket
	//input := &s3.CreateBucketInput{
	//	Bucket: &bktName,
	//}
	//
	//_, err = utils.MakeBucket(context.TODO(), client, input)

	// 上传文件
	file, err := os.Open("testfile")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	input := &s3.PutObjectInput{
		Bucket: &bktName,
		Key:    &keyname,
		Body:   file,
	}
	_, err = utils.PutFile(context.TODO(), client, input)

	if err != nil {
		panic(err)
	}

	// 创建presign客户端
	//psClient := s3.NewPresignClient(client)

	//fmt.Println("创建桶")
	//createBucket(client, "minio-test")
	//fmt.Println("列出桶")
	//listBuckets(client)

}
