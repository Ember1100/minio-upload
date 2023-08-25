package conf

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"
)

const Endpoint string = "http://1.12.42.181:9001"                         // MinIO 服务的地址
const accessKeyID string = "1hPSpBEublHyVQUm6CR0"                         // MinIO 访问密钥的 Access Key ID
const secretAccessKey string = "oRZaLMOG8NhXB9lyCwi9R4LBBy4RGJhGd9qNJkLa" // MinIO 访问密钥的 Secret Access Key
const BucketName string = "minio03"                                       // MinIO 存储桶的名称

func GetAwsS3() (*s3.S3, error) {
	// 创建 AWS 会话
	sess, err := session.NewSession(&aws.Config{
		Endpoint:         aws.String(Endpoint),
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
		Region:           aws.String("us-east-1"), // 设置所需的 AWS 区域
		S3ForcePathStyle: aws.Bool(true),          // 开启强制使用路径样式的设置
	})
	if err != nil {
		return nil, err
	}
	// 创建 S3 服务客户端
	svc := s3.New(sess)
	return svc, nil
}

func GetAwsSts() (*sts.STS, error) {
	// 创建 AWS 会话
	sess, err := session.NewSession(&aws.Config{
		Endpoint:         aws.String(Endpoint),
		Credentials:      credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
		Region:           aws.String("us-east-1"),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}
	// 创建 STS 客户端
	stsClient := sts.New(sess)
	return stsClient, nil
}
