package backends

import (
	"../../config"
	"github.com/ncw/swift"
	"io"
	"os"
)

type swiftBackend struct {
	conn *swift.Connection
}

var swifter *swiftBackend
var containerName = config.GetOrElse("swift_container_name", "tcbtest")

func (sw *swiftBackend) GetReader(name string) (io.ReadCloser, error) {
	objOF, _, err := sw.conn.ObjectOpen(containerName, name, false, nil)
	if err != nil {
		return nil, err
	} else {
		return objOF, nil
	}
}
func (sw *swiftBackend) StoreObject(name string, source *os.File, metadata map[string]string) error {
	objCF, err := sw.conn.ObjectCreate(containerName, name, false, "", "application/binary", nil)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(objCF, source)
	objCF.Close()
	if err != nil {
		panic(err)
	}
	return err
}
func (sw *swiftBackend) DeleteObject(name string) error {
	return sw.conn.ObjectDelete(containerName, name)
}

func GetSwiftBackend() *swiftBackend {

	if swifter == nil {
		swifter = new(swiftBackend)
		swiftUserName := config.GetOrElse("swift_user", "user")
		swiftURL := config.GetOrElse("swift_url", "url")
		swiftAPIKey := config.GetOrElse("swift_api_key", "key")
		c := swift.Connection{
			UserName: swiftUserName,
			ApiKey:   swiftAPIKey,
			AuthUrl:  swiftURL,
		}
		err := c.Authenticate()
		if err != nil {
			panic(err)
		}

		swifter.conn = &c
	}
	return swifter
}
