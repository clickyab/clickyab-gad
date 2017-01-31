package cache

import (
	"bytes"
	"encoding/gob"
	"entity"
)

// CacheProvider is the cache backend
type CacheProvider interface {
	// Cache is called to store the cache
	Cache(entity.Cacheable) error
	// Hit called when we need to load the cache
	Hit(string, entity.Cacheable) error
}

// InterfaceToByte save interface into byte
func InterfaceToByte(in interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}

	enc := gob.NewEncoder(buf)
	err := enc.Encode(in)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ByteToInterface return object from byte
func ByteToInterface(b []byte, out interface{}) error {
	buf := bytes.NewBuffer(b)

	dnc := gob.NewDecoder(buf)
	return dnc.Decode(out)
}

var cache CacheProvider

// Cache the entity
func Cache(e entity.Cacheable, err error) error {
	if err != nil {
		return err
	}
	return cache.Cache(e)
}

// Hit the cache
func Hit(key string, out entity.Cacheable) error {
	return cache.Hit(key, out)
}
