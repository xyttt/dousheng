package MinIO

import (
	"bytes"
	"context"
	"dousheng/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

var Client *minio.Client
var err error

func createBucket() {
	ctx := context.Background()

	bucketName := config.BucketName
	location := config.Location

	err = Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		exists, errBucketExists := Client.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %v\n", bucketName)
		} else {
			log.Printf("failed with makeBucket , error : %v", errBucketExists)
		}

	} else {
		log.Printf("Successfully created %v\n", bucketName)
	}
}

func Init() {
	Client, err = minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		log.Printf("failed with minio connecting , error : %v", err)
	} else {
		log.Print("connect minio successfully.\n")
	}

	createBucket()
}

func FileUpLoader(buf *bytes.Buffer, saveName string, conType string) (string, error) {
	bucketName := config.BucketName
	objectName := saveName
	fileStream := buf
	objectSize := int64(buf.Len())

	_, err := Client.PutObject(context.Background(),
		bucketName,
		objectName,
		fileStream,
		objectSize,
		minio.PutObjectOptions{
			ContentType: conType,
		})

	if err != nil {
		log.Printf("upload %v failed, %v", saveName, err)
		return "", err
	}
	log.Printf("upload %v successfully", saveName)

	fileURL := "http://" + config.Endpoint + "/" + bucketName + "/" + objectName
	log.Printf("videoURL:  %v ", fileURL)

	return fileURL, err
}
