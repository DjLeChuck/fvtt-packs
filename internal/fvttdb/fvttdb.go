package fvttdb

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"log"
)

type FvttDb struct {
	db *leveldb.DB
}

func Open(path string) (*FvttDb, error) {
	db, err := leveldb.OpenFile(path, &opt.Options{
		ErrorIfMissing: true,
		ReadOnly:       true,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot open db \"%s\": %s\n", path, err)
	}

	return &FvttDb{db: db}, nil
}

func (fvttDb *FvttDb) Close() {
	if err := fvttDb.db.Close(); err != nil {
		log.Fatalf("cannot close DB: %s", err)
	}
}

func (fvttDb *FvttDb) IterateAll(fn func(iter iterator.Iterator) error) error {
	iter := fvttDb.db.NewIterator(nil, nil)
	for iter.Next() {
		if itErr := fn(iter); itErr != nil {
			return fmt.Errorf("cannot iterate key %s: %s\n", iter.Key(), itErr)
		}
	}
	iter.Release()

	return iter.Error()
}

func (fvttDb *FvttDb) Get(key string) ([]byte, error) {
	v, err := fvttDb.db.Get([]byte(key), nil)
	if err != nil {
		return nil, fmt.Errorf("cannot get entry %s: %s\n", key, err)
	}

	return v, nil
}
