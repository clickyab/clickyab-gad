package store

import (
	"sync"

	"github.com/Sirupsen/logrus"
)

type dataStruct struct {
	data string
	w    *sync.WaitGroup
}

var (
	dataLock = &sync.RWMutex{}
	data     = make(map[string]*dataStruct)
)

// Reserve try to reserve a response, the lock is gone after 30 sec
func Reserve(key string) {
	dataLock.Lock()
	defer dataLock.Unlock()

	tmp := dataStruct{
		w: &sync.WaitGroup{},
	}
	tmp.w.Add(1)
	data[key] = &tmp
}

// Set the data in already reserved key, unlock the key after that
func Set(key string, v string) {
	dataLock.Lock()
	defer dataLock.Unlock()

	d, ok := data[key]
	if !ok {
		logrus.Panic("Not reserved key")
	}

	if d.data != "" {
		logrus.Panic("key set twice")
	}

	d.data = v
	d.w.Done()
}

// Get the key from the system
func Get(key string) (string, bool) {
	dataLock.RLock()

	d, ok := data[key]
	if !ok {
		return "", false
	}

	dataLock.RUnlock()

	d.w.Wait()
	res := d.data
	go func() {
		dataLock.Lock()
		defer dataLock.Unlock()
		delete(data, key)
	}()
	return res, true
}
