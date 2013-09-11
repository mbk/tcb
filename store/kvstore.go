package store

import (
	"bytes"
	gob "encoding/gob"
	"errors"
	kv "github.com/cznic/kv"
	cfg "github.com/mbk/tcb/config"
	"os"
)

type KVStore struct {
	kv *kv.DB
}

var kvStore = new(KVStore)

func (kv *KVStore) Put(name string, key string, value string) (e error) {

	updater := func(akey, old []byte) (newval []byte, write bool, err error) {
		if old == nil {
			s := make(map[string]string)
			s[key] = value

			newbuf := new(bytes.Buffer)
			encoder := gob.NewEncoder(newbuf)

			err = encoder.Encode(s)
			if err != nil {
				err = errors.New("Could not encode GOB value for " + name)
			}
			err = nil
			write = true
			newval = newbuf.Bytes()

		} else {
			buf := new(bytes.Buffer)
			buf.Write(old)
			s := make(map[string]string)
			decoder := gob.NewDecoder(buf)
			err := decoder.Decode(&s)
			if err != nil {
				err = errors.New("Could not decode GOB value for " + name)
			}
			s[key] = value

			newbuf := new(bytes.Buffer)
			encoder := gob.NewEncoder(newbuf)

			err = encoder.Encode(s)
			if err != nil {
				err = errors.New("Could not encode GOB value for " + name)
			}
			err = nil
			write = true
			newval = newbuf.Bytes()
		}
		return newval, write, err
	}

	_, written, err := kv.kv.Put(nil, []byte(name), updater)

	if err != nil {
		e = err
	}
	if written {
		e = nil
	} else {
		e = errors.New("Could not write name/key/value pair")
	}
	return e
}

func (kv *KVStore) Get(name string, key string) (value string, err error) {
	val := make([]byte, 16384)
	_, e := kv.kv.Get(val, []byte(name))

	if e != nil {
		err = errors.New("Could not execute Get on Key value store for key " + name)
		value = ""
	} else {
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
	}
	return value, err
}

func (kv *KVStore) Exists(name string, key string) bool {
	value, _ := kv.Get(name, key)
	return value != ""
}

func (kv *KVStore) DeleteName(name string) (err error) {
	e := kv.kv.Delete([]byte(name))
	if e != nil {
		err = errors.New("Could not delete key " + name)
	} else {
		err = nil
	}
	return
}

func (kv *KVStore) Delete(name, key string) (err error) {
	kv.Put(name, key, "")
	return
}

func EnsureKVStore() Store {
	if kvStore.kv == nil {
		//TBD: Set the compare function to one that compares to gob encoded bytes
		options := new(kv.Options)
		dbname := cfg.GetOrElse("file_metadata_path", "./tcb.db")
		_, dbexists := os.Stat(dbname)
		var db *kv.DB
		var err error
		//Create the db as file
		if dbexists != nil {
			db, err = kv.Create(dbname, options)
		} else {
			db, err = kv.Open(dbname, options)
		}
		if err != nil {
			panic(err)
		}

		kvStore.kv = db
	}
	return kvStore
}

func EnsureKVMemoryStore() Store {
	if kvStore.kv == nil {
		options := new(kv.Options)
		db, err := kv.CreateMem(options)
		if err != nil {
			panic(err)
		}
		kvStore.kv = db
	}
	return kvStore
}
