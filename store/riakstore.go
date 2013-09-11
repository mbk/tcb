package store

import (
	"bytes"
	gob "encoding/gob"
	"errors"
	cfg "github.com/mbk/tcb/config"
	riak "github.com/tpjg/goriakpbc"
	//"os"
	"strconv"
)

type riakstore struct {
	bucket *riak.Bucket
}

var store *riakstore

func (rs *riakstore) Put(name string, key string, value string) error {
	exists, err := rs.bucket.Exists(name)
	if exists {
		obj, err := rs.bucket.Get(name)
		if err != nil {
			return err
		}
		old := obj.Data
		//update case
		buf := new(bytes.Buffer)
		buf.Write(old)
		s := make(map[string]string)
		decoder := gob.NewDecoder(buf)
		err = decoder.Decode(&s)
		if err != nil {
			return errors.New("Could not decode GOB value for " + name)
		}
		s[key] = value

		newbuf := new(bytes.Buffer)
		encoder := gob.NewEncoder(newbuf)

		err = encoder.Encode(s)
		if err != nil {
			return errors.New("Could not encode GOB value for " + name)
		}
		obj.Data = newbuf.Bytes()
		return obj.Store()
	} else {
		//New key
		s := make(map[string]string)
		s[key] = value
		newbuf := new(bytes.Buffer)
		encoder := gob.NewEncoder(newbuf)

		err = encoder.Encode(s)
		if err != nil {
			err = errors.New("Could not encode GOB value for " + name)
		}
		obj := rs.bucket.NewObject(name)
		obj.ContentType = "application/binary"
		obj.Data = newbuf.Bytes()
		return obj.Store()
	}
	return nil
}

func (rs *riakstore) Get(name string, key string) (value string, err error) {

	exists, err := rs.bucket.Exists(name)
	if exists && err == nil {
		object, err := rs.bucket.Get(name)
		if err != nil {
			return "", err
		}

		val := object.Data

		if val == nil {
			err = errors.New("Value for key not found.")
			value = ""
		} else {
			buf := new(bytes.Buffer)
			buf.Write(val)
			s := make(map[string]string)

			decoder := gob.NewDecoder(buf)
			err := decoder.Decode(&s)
			if err != nil {
				err = errors.New("Could not decode keys for " + name)
				value = ""
			}

			value = s[key]
			if value == "" {
				err = errors.New("Value for key not found.")
			}
		}
	} else {
		value = ""
		err = errors.New("Key does not exist")
	}
	return value, err
}

func (rs *riakstore) Exists(name string, key string) bool {
	value, _ := rs.Get(name, key)
	return value != ""
}

func (rs *riakstore) Delete(name, key string) error {
	return rs.Put(name, key, "")
}

func (rs *riakstore) DeleteName(name string) error {
	return rs.bucket.Delete(name)
}

func EnsureRiakStore() *riakstore {
	if store == nil {
		bucketName := cfg.GetOrElse("riak_bucket_name", "tcb")
		host := cfg.GetOrElse("riak_address", "127.0.0.1:8087")
		poolsize, err := strconv.Atoi(cfg.GetOrElse("riak_pool_size", "8"))
		if err != nil {
			poolsize = 8
		}
		err = riak.ConnectClientPool(host, poolsize)
		if err != nil {
			panic(err)
		}
		bucket, err := riak.NewBucket(bucketName)
		if err != nil {
			panic(err)
		}
		store = new(riakstore)
		store.bucket = bucket
	}
	return store
}
