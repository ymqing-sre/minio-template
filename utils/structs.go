package utils

import (
	"context"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3BucketExistsAPI 定义桶存在判断看接口
type S3BucketExistsAPI interface {
	HeadBucket(ctx context.Context,
		params *s3.HeadBucketInput,
		optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error)
}

// S3CreateBucketAPI 定义创建桶接口
type S3CreateBucketAPI interface {
	CreateBucket(ctx context.Context,
		params *s3.CreateBucketInput,
		optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error)
}

// S3ListBucketsAPI 定义列出桶接口
type S3ListBucketsAPI interface {
	ListBuckets(ctx context.Context,
		params *s3.ListBucketsInput,
		optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
}

// S3DeleteBucketAPI 定义删除桶接口
type S3DeleteBucketAPI interface {
	DeleteBucket(ctx context.Context,
		params *s3.DeleteBucketInput,
		optFns ...func(*s3.Options)) (*s3.DeleteBucketOutput, error)
}

// S3ObjectExistsAPI 定义对象存在判断接口
type S3ObjectExistsAPI interface {
	HeadObject(ctx context.Context,
		params *s3.HeadObjectInput,
		optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
}

// S3PresignGetObjectAPI 定义获取对象预签名URL接口
type S3PresignGetObjectAPI interface {
	PresignGetObject(
		ctx context.Context,
		params *s3.GetObjectInput,
		optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

// S3PresignPutObjectAPI 定义上传对象预签名URL接口
type S3PresignPutObjectAPI interface {
	PresignPutObject(
		ctx context.Context,
		params *s3.PutObjectInput,
		optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

// S3PutObjectAPI 定义上传对象接口
type S3PutObjectAPI interface {
	PutObject(ctx context.Context,
		params *s3.PutObjectInput,
		optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

// S3ListObjectsAPI 定义获取对象列表接口
type S3ListObjectsAPI interface {
	ListObjectsV2(ctx context.Context,
		params *s3.ListObjectsV2Input,
		optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

// S3DeleteObjectAPI 定义删除对象接口
type S3DeleteObjectAPI interface {
	DeleteObject(ctx context.Context,
		params *s3.DeleteObjectInput,
		optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
}

// S3CopyObjectAPI 定义复制对象接口
type S3CopyObjectAPI interface {
	CopyObject(ctx context.Context,
		params *s3.CopyObjectInput,
		optFns ...func(*s3.Options)) (*s3.CopyObjectOutput, error)
}

// S3GetObjectAPI 定义下载对象接口
type S3GetObjectAPI interface {
	GetObject(ctx context.Context,
		params *s3.GetObjectInput,
		optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}
