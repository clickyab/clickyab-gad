package core

import (
	"assert"
	"context"
	"entity"
	"sync"
	"time"
)

var allProviders []providerData

type providerData struct {
	name         string
	provider     entity.AdProvider
	timeout      time.Duration
	lowerTimeout bool
}

func (p *providerData) watch(imp entity.Impression, timeout time.Duration) map[int][]entity.Advertise {
	ctx, cnl := context.WithCancel(context.Background())
	defer cnl()
	t := p.timeout
	if p.lowerTimeout {
		t = timeout
	}
	chn := p.provider.Provide(ctx, imp)
	select {
	case <-time.After(t):
		return nil
	case data := <-chn:
		return data
	}
}

// Register is used to handle new layer in system
func Register(name string, provider entity.AdProvider, timeout time.Duration, callOnLowerTimeOut bool) {

	for i := range allProviders {
		assert.True(allProviders[i].name != name, "[BUG] same name registered twice")
	}

	allProviders = append(
		allProviders,
		providerData{
			name:         name,
			provider:     provider,
			timeout:      timeout,
			lowerTimeout: callOnLowerTimeOut,
		},
	)
}

// Call is for getting the current ads for this imp
func Call(imp entity.Impression, timeout time.Duration) map[int][]entity.Advertise {
	wg := sync.WaitGroup{}
	l := len(allProviders)
	wg.Add(l)
	allRes := make(chan map[int][]entity.Advertise, l)
	for i := range allProviders {
		if allProviders[i].lowerTimeout || allProviders[i].timeout <= timeout {
			go func(inner int) {
				defer wg.Done()
				res := allProviders[inner].watch(imp, timeout)
				if res != nil {
					allRes <- res
				}
			}(i)
		}
	}

	wg.Wait()
	// The close is essential here.
	close(allRes)
	res := make(map[int][]entity.Advertise)
	for provided := range allRes {
		for j := range provided {
			res[j] = append(res[j], provided[j]...)
		}
	}

	return res
}
