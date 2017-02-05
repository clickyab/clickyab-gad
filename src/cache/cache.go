package cache

import (
	"encoding/gob"
	"entity"
	"io"
	"time"
)

// CacheProvider is the cache backend
type CacheProvider interface {
	// Cache is called to store the cache
	Cache(entity.Cacheable, time.Duration) error
	// Hit called when we need to load the cache
	Hit(string, entity.Cacheable) error
}

// Wrapper is a provider with support for inner entity
type Wrapper interface {
	entity.Cacheable
	Entity() interface{}
}

type cachable struct {
	entity interface{}
	key    string
}

var cache CacheProvider

// Decode try to decode cookie profile into gob
func (cp *cachable) Decode(w io.Writer) error {
	enc := gob.NewEncoder(w)
	return enc.Encode(cp.entity)
}

// Encode try to encode cookie profile from gob
func (cp *cachable) Encode(i io.Reader) error {
	dnc := gob.NewDecoder(i)
	return dnc.Decode(cp.entity)
}

func (cp *cachable) String() string {
	return cp.key
}

func (cp *cachable) Entity() interface{} {
	return cp.entity
}

// Cache the entity
func Cache(e entity.Cacheable, t time.Duration, err error) error {
	if err != nil {
		return err
	}
	return cache.Cache(e, t)
}

// Hit the cache
func Hit(key string, out entity.Cacheable) error {
	return cache.Hit(key, out)
}

// CreateCacheWrapper return an cachable object for this ntt
func CreateCacheWrapper(key string, ntt interface{}) Wrapper {
	return &cachable{
		key:    key,
		entity: ntt,
	}
}

// Register a new cache provider
func Register(p CacheProvider) {
	cache = p
}
