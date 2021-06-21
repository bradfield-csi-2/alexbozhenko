package main

/*
#cgo CFLAGS: -I/usr/include
#cgo LDFLAGS:	 -l leveldb
#include <leveldb/c.h>
#include <stdlib.h>
#include <stdio.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type level struct {
	db           *C.leveldb_t
	options      *C.leveldb_options_t
	readOptions  *C.leveldb_readoptions_t
	writeOptions *C.leveldb_writeoptions_t
}

func db_init(db_name string) (level, error) {
	options := C.leveldb_options_create()
	C.leveldb_options_set_create_if_missing(options, 1)
	writeOptions := C.leveldb_writeoptions_create()
	readOptions := C.leveldb_readoptions_create()
	dbName := C.CString(db_name)
	defer C.free(unsafe.Pointer(dbName))

	errptr := (*C.char)(nil)
	defer C.free(unsafe.Pointer(errptr))
	db := C.leveldb_open(options, dbName, &errptr)
	if unsafe.Pointer(errptr) != C.NULL {
		return level{}, fmt.Errorf("error opening a database %s. %s",
			db_name, C.GoString(errptr))
	}
	return level{db, options, readOptions, writeOptions}, nil

}

func get(db level, key string) (string, error) {
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))
	vallen := (C.size_t(0))

	errptr := (*C.char)(nil)
	defer C.free(unsafe.Pointer(errptr))

	get_result := C.leveldb_get(db.db, db.readOptions, k,
		C.size_t(len(key)), &vallen,
		&errptr)

	if unsafe.Pointer(errptr) != C.NULL {
		return "", fmt.Errorf("error getting key %s: %s", key,
			C.GoString(errptr))
	}
	if unsafe.Pointer(get_result) == C.NULL {
		return "", fmt.Errorf("key %s not found", key)
	} else {
		return C.GoString(get_result), nil
	}

}

func put(db level, key string, value string) error {
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))
	v := C.CString(value)
	defer C.free(unsafe.Pointer(v))

	errptr := (*C.char)(nil)

	defer C.free(unsafe.Pointer(errptr))
	C.leveldb_put(db.db, db.writeOptions, k, C.size_t(len(key)),
		v, C.size_t(len(value)),
		&errptr)

	if unsafe.Pointer(errptr) != C.NULL {
		return fmt.Errorf("error putting key %s: %s", key,
			C.GoString(errptr))
	}
	return nil
}

func delete(db level, key string) error {
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))

	errptr := (*C.char)(nil)
	defer C.free(unsafe.Pointer(errptr))
	C.leveldb_delete(db.db, db.writeOptions, k, C.size_t(len(key)),
		&errptr)

	if unsafe.Pointer(errptr) != C.NULL {
		return fmt.Errorf("error deleting key %s: %s", key,
			C.GoString(errptr))
	}
	return nil
}

func main() {
	db, err := db_init("/tmp/leveldb_test")
	if err != nil {
		panic(err)
	}
	put(db, "Crimea", "is ours")

	if val, err := get(db, "Crimea"); err == nil {
		fmt.Println(val)
	} else {
		fmt.Println(err)
	}

	if err = delete(db, "Crimea"); err != nil {
		fmt.Println(err)
	}

	if val, err := get(db, "Crimea"); err == nil {
		fmt.Println(val)
	} else {
		fmt.Println(err)
	}

	if err = delete(db, "Crimea"); err != nil {
		fmt.Println(err)
	}

}
