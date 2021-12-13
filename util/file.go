package util

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func UpLoadFile(filename, path string) (string, string, error) {
	file, err := os.Open(filename)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", err)
	}

	defer file.Close()

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(AwsRegion),
		Credentials: credentials.NewStaticCredentials(AccessKeyID, SecretAccessKey, ""),
	},
	)

	uploader := s3manager.NewUploader(sess)

	// Upload the file's body to S3 bucket as an object with the key being the
	// same as the filename.
	s3path := filepath.Join(path, filepath.Base(filename))
	fulls3path := fmt.Sprintf("http://drpiggy.s3.amazonaws.com/%s", s3path)
	cdnpath := fmt.Sprintf("http://d2xl03dasqf9l8.cloudfront.net/%s", s3path)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("drpiggy"),

		// Can also use the `filepath` standard library package to modify the
		// filename as need for an S3 object key. Such as turning absolute path
		// to a relative path.

		Key: aws.String(s3path),
		ACL: aws.String("public-read"),

		// The file to be uploaded. io.ReadSeeker is preferred as the Uploader
		// will be able to optimize memory when uploading large content. io.Reader
		// is supported, but will require buffering of the reader's bytes for
		// each part.
		Body: file,
	})
	if err == nil {
		fmt.Printf("Successfully uploaded %q to %q\n", filename, "drpiggy")
	}
	return fulls3path, cdnpath, err

}

//UpLoadFileFromByte 通过字节上传图片
func UpLoadFileFromByte(filename, path string, reader *io.Reader) (string, string, error) {

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(AwsRegion),
		Credentials: credentials.NewStaticCredentials(AccessKeyID, SecretAccessKey, ""),
	},
	)

	uploader := s3manager.NewUploader(sess)

	// Upload the file's body to S3 bucket as an object with the key being the
	// same as the filename.
	s3path := filepath.Join(path, filepath.Base(filename))
	fulls3path := fmt.Sprintf("http://drpiggy.s3.amazonaws.com/%s", s3path)
	cdnpath := fmt.Sprintf("http://d2xl03dasqf9l8.cloudfront.net/%s", s3path)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("drpiggy"),

		// Can also use the `filepath` standard library package to modify the
		// filename as need for an S3 object key. Such as turning absolute path
		// to a relative path.

		Key: aws.String(s3path),
		ACL: aws.String("public-read"),

		// The file to be uploaded. io.ReadSeeker is preferred as the Uploader
		// will be able to optimize memory when uploading large content. io.Reader
		// is supported, but will require buffering of the reader's bytes for
		// each part.
		Body: *reader,
	})
	if err == nil {
		fmt.Printf("Successfully uploaded %q to %q\n", filename, "drpiggy")
	}
	return fulls3path, cdnpath, err
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

// Exist 文件是否存在
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
