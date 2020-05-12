package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func createSession() *s3.S3 {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewSharedCredentials("", "personal"),
	})

	if err != nil {
		exitErrorf("Unable to connect to region")
	}

	svc := s3.New(sess)
	return svc
}

func createBucket(svc *s3.S3, bucketName string, region string) {

	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(region),
		},
	})

	if err != nil {
		exitErrorf("Unable to create bucket, %v", err)
	}

	fmt.Printf("Waiting for bucket %q to be created...\n", bucketName)

	err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
}

func listBucket(svc *s3.S3) {
	result, err := svc.ListBuckets(nil)
	if err != nil {
		exitErrorf("Unable to list buckets, %v", err)
	}

	fmt.Println("Buckets:")

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}
}

func createObject(svc *s3.S3, bucket string, key string) {
	_, err := svc.PutObject(&s3.PutObjectInput{
		Body:   strings.NewReader("Hello World!"),
		Bucket: &bucket,
		Key:    &key,
	})

	if err != nil {
		exitErrorf("Unable to create object %v in bucket %v", key, bucket)
	}
}

func removeObject(svc *s3.S3, bucket string, obj string) {
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucket), Key: aws.String(obj)})
	if err != nil {
		exitErrorf("Unable to delete object %q from bucket %q, %v", obj, bucket, err)
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(obj),
	})
	if err != nil {
		exitErrorf("Error occurred while waiting for object %q to be deleted, %v", obj, err)
	}

	fmt.Printf("Object %q successfully deleted\n", obj)
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
