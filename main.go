package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ymqing-sre/minio-template/utils"
	"io"
	"net/url"
	"os"
)

var (
	action    string
	bktName   string
	filename  string
	path      string
	objName   string
	srcObj    string
	dstObj    string
	srcBucket string
	dstBucket string
)

func init() {
	flag.StringVar(&action, "a", "", "对对象存储进行的操作指令")
	flag.StringVar(&bktName, "b", "", "桶名称")
	flag.StringVar(&filename, "f", "", "文件名")
	flag.StringVar(&path, "p", "", "文件路径")
	flag.StringVar(&objName, "o", "", "对象名称")
	flag.StringVar(&srcBucket, "s", "", "源Bucket名称")
	flag.StringVar(&srcObj, "so", "", "源对象")
	flag.StringVar(&dstBucket, "d", "", "目的Bucket名称")
	flag.StringVar(&dstObj, "do", "", "目的对象")
}

func main() {
	flag.Parse()
	if action != "" {
		// 读取配置文件
		endpointURL, key, secret := utils.LoadConfig("./conf/", "config")
		// 定义EndpointResolver
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL: endpointURL,
			}, nil
		})
		// 定义CredentialsProvider
		customCredentials := credentials.NewStaticCredentialsProvider(key, secret, "")
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

		switch action {
		case "listBuckets":
			input := &s3.ListBucketsInput{}
			res, err := utils.GetAllBuckets(context.TODO(), client, input)
			if err != nil {
				panic(err)
			}
			for _, bkt := range res.Buckets {
				fmt.Println(*bkt.Name + ": " + bkt.CreationDate.Format("2006-01-02 15:01:05"))
			}
		case "createBucket":
			if bktName != "" {
				input := &s3.HeadBucketInput{
					Bucket: &bktName,
				}
				if utils.BucketExists(context.TODO(), client, input) {
					fmt.Println("Bucket", bktName, "is already exists")
				} else {
					input := &s3.CreateBucketInput{
						Bucket: &bktName,
					}
					_, err := utils.MakeBucket(context.TODO(), client, input)
					if err != nil {
						panic(err)
					}
					fmt.Println("Bucket " + bktName + " created.")
				}

			} else {
				fmt.Println("You must supply a bucket name (-b BUCKET)")
			}
		case "copyObj":
			if srcBucket != "" && dstBucket != "" && dstObj != "" && srcObj != "" {
				// 判断源对象和目的桶是否存在
				inputSrc := &s3.HeadObjectInput{
					Bucket: &srcBucket,
					Key:    &srcObj,
				}
				inputDst := &s3.HeadBucketInput{
					Bucket: &dstBucket,
				}
				if utils.ObjectExists(context.TODO(), client, inputSrc) && utils.BucketExists(context.TODO(), client, inputDst) {
					copySource := srcBucket + "/" + srcObj
					input := &s3.CopyObjectInput{
						Bucket:     &dstBucket,
						CopySource: aws.String(url.PathEscape(copySource)),
						Key:        &dstObj,
					}
					_, err := utils.CopyItem(context.TODO(), client, input)
					if err != nil {
						panic(err)
					}
					fmt.Println("Copied " + srcObj + " from " + srcBucket + " to " + dstBucket)
				} else {
					fmt.Println(srcBucket+"/"+srcObj, "or", dstBucket, "may not exists")
				}
			} else {
				fmt.Println("You must supply the bucket to copy from (-s BUCKET), to (-d BUCKET), and object to copy (-o OBJECT")
			}
		case "putFile":
			if bktName != "" && filename != "" {
				input := &s3.HeadBucketInput{
					Bucket: &bktName,
				}
				if utils.BucketExists(context.TODO(), client, input) {
					file, err := os.Open(filename)
					if err != nil {
						panic(err)
					}
					defer func(file *os.File) {
						err := file.Close()
						if err != nil {
							panic(err)
						}
					}(file)
					var keyname string
					if path != "" {
						keyname = path + "/" + filename
					} else {
						keyname = filename
					}
					input := &s3.PutObjectInput{
						Bucket: &bktName,
						Key:    &keyname,
						Body:   file,
					}
					_, err = utils.PutFile(context.TODO(), client, input)
					if err != nil {
						panic(err)
					}
					fmt.Println("Upload File " + filename)
				} else {
					fmt.Println("Bucket", bktName, "is not founded")
				}
			} else {
				fmt.Println("You must supply a bucket name (-b BUCKET) and file name (-f FILE)")
			}
		case "getFile":
			if bktName != "" && objName != "" && filename != "" {
				// 先判断对象是否存在
				input := &s3.HeadObjectInput{
					Bucket: &bktName,
					Key:    &objName,
				}
				if utils.ObjectExists(context.TODO(), client, input) {
					input := &s3.GetObjectInput{
						Bucket: &bktName,
						Key:    &objName,
					}
					res, err := utils.GetFile(context.TODO(), client, input)
					if err != nil {
						panic(err)
					}
					defer func(Body io.ReadCloser) {
						err := Body.Close()
						if err != nil {

						}
					}(res.Body)
					localFile, err := os.Create(filename)
					if err != nil {
						panic(err)
					}
					defer func(localFile *os.File) {
						err := localFile.Close()
						if err != nil {

						}
					}(localFile)
					_, err = io.CopyN(localFile, res.Body, res.ContentLength)
					if err != nil {
						panic(err)
					}
				} else {
					fmt.Println("Object", objName, "is not found")
				}
			} else {
				fmt.Println("You must supply a bucket name (-b BUCKET), object name (-o OBJECT) and local file name (-f FILE)")
			}
		case "deleteObj":
			if bktName != "" && objName != "" {
				// 先判断对象是否存在
				input := &s3.HeadObjectInput{
					Bucket: &bktName,
					Key:    &objName,
				}
				if utils.ObjectExists(context.TODO(), client, input) {
					input := &s3.DeleteObjectInput{
						Bucket: &bktName,
						Key:    &objName,
					}
					_, err := utils.DeleteObj(context.TODO(), client, input)
					if err != nil {
						panic(err)
					}
					fmt.Println("Object " + objName + " Deleted")
				} else {
					fmt.Println("Object", objName, "is not found")
				}
			} else {
				fmt.Println("You must supply a bucket name (-b BUCKET) and object name (-o OBJECT)")
			}
		case "listFiles":
			if bktName != "" {
				input := &s3.ListObjectsV2Input{
					Bucket: &bktName,
				}
				res, err := utils.GetObjects(context.TODO(), client, input)
				if err != nil {
					panic(err)
				}
				fmt.Println("Objects in " + bktName + ":")

				for _, item := range res.Contents {
					fmt.Println("Name:          ", *item.Key)
					fmt.Println("Last modified: ", *item.LastModified)
					fmt.Println("Size:          ", item.Size)
					fmt.Println("Storage class: ", item.StorageClass)
					fmt.Println("")
				}
				fmt.Println("Found", len(res.Contents), "items in bucket", bktName)
			} else {
				fmt.Println("You must supply the name of a bucket (-b BUCKET)")
			}
		case "removeBucket":
			if bktName != "" {
				input := &s3.HeadBucketInput{
					Bucket: &bktName,
				}
				if utils.BucketExists(context.TODO(), client, input) {
					input := &s3.DeleteBucketInput{
						Bucket: &bktName,
					}
					_, err := utils.RemoveBucket(context.TODO(), client, input)
					if err != nil {
						panic(err)
					}
					fmt.Println("Bucket " + bktName + " removed.")
				} else {
					fmt.Println("Bucket", bktName, "is not founded")
				}
			} else {
				fmt.Println("You must supply the name of a bucket (-b BUCKET)")
			}
		case "putPresignedURL":
			if bktName != "" && objName != "" {
				input := &s3.HeadBucketInput{
					Bucket: &bktName,
				}
				if utils.BucketExists(context.TODO(), client, input) {
					input := &s3.PutObjectInput{
						Bucket: &bktName,
						Key:    &objName,
					}
					psClient := s3.NewPresignClient(client)
					res, err := utils.PutPresignedURL(context.TODO(), psClient, input)
					if err != nil {
						panic(err)
					}
					fmt.Println("The Object Put URL:")
					fmt.Println(res.URL)
				}
			} else {
				fmt.Println("Bucket", bktName, "is not founded")
			}
		case "getPresignedURL":
			if bktName != "" && objName != "" {
				// 先判断对象是否存在
				input := &s3.HeadObjectInput{
					Bucket: &bktName,
					Key:    &objName,
				}
				if utils.ObjectExists(context.TODO(), client, input) {
					input := &s3.GetObjectInput{
						Bucket: &bktName,
						Key:    &objName,
					}
					psClient := s3.NewPresignClient(client)
					res, err := utils.GetPresignedURL(context.TODO(), psClient, input)
					if err != nil {
						panic(err)
					}
					fmt.Println("The Object Get URL:")
					fmt.Println(res.URL)
				} else {
					fmt.Println("Object", objName, "is not found")
				}
			} else {
				fmt.Println("You must supply a bucket name (-b BUCKET) and object name (-o OBJECT)")
			}
		case "objExists":
			if bktName != "" && objName != "" {
				input := &s3.HeadObjectInput{
					Bucket: &bktName,
					Key:    &objName,
				}
				res := utils.ObjectExists(context.TODO(), client, input)
				fmt.Println(res)
			} else {
				fmt.Println("You must supply a bucket name (-b BUCKET) and object name (-o OBJECT)")
			}
		case "bktExists":
			if bktName != "" {
				input := &s3.HeadBucketInput{
					Bucket: &bktName,
				}
				res := utils.BucketExists(context.TODO(), client, input)
				fmt.Println(res)
			} else {
				fmt.Println("You must supply a bucket name (-b BUCKET)")
			}
		default:
			fmt.Println("The action must be one of " +
				"[listBuckets, createBucket, copyObj, putFile, getFile, deleteObj, listFiles, removeBucket," +
				" putPresignedURL, getPresignedURL, objExists, bktExists]")
		}
	}
}
