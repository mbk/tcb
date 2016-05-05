package backends

import (
	"github.com/mbk/tcb/config"
	"github.com/mbk/tcb/store"
	"io"
	aws "launchpad.net/goamz/aws"
	s3 "launchpad.net/goamz/s3"
	"os"
	"strconv"
)

type S3Backend struct {
	bucket *s3.Bucket
}

//TBD implement StorageBackend methods

func (s3b *S3Backend) GetReader(name string) (io.ReadCloser, error) {
	return s3b.bucket.GetReader(name)
}

func (s3b *S3Backend) StoreObject(name string, source *os.File, path string, metadata store.Store) error {
	strlength, err := metadata.Get(path, "length")
	length, err := strconv.ParseInt(strlength, 10, 64)
	defer source.Close()

	if err != nil {
		return err
	} else {
		//data, err := ioutil.ReadAll(source)
		if err != nil {
			return err
		}
		//err = s3b.bucket.Put(name, data, "application/binary", s3.Private)
		err = s3b.bucket.PutReader(name, source, length, "application/binary", s3.Private)
		if err != nil {
			return err
		}
	}
	return err
}
func (s3b *S3Backend) DeleteObject(name string) error {
	return s3b.bucket.Del(name)
}

func GetS3Backend() *S3Backend {
	var awsId = config.GetOrElse("aws_id", "nope")
	var secretKey = config.GetOrElse("aws_secret_key", "nope")
	var auth = aws.Auth{AccessKey: awsId, SecretKey: secretKey}
	var region = aws.Region{S3Endpoint: config.GetOrElse("aws_s3_region", "https://s3-eu-west-1.amazonaws.com")}
	var es3 = s3.New(auth, region)
	aBucket := es3.Bucket(config.GetOrElse("aws_s3_bucket", "nope"))
	return &S3Backend{bucket: aBucket}
}
