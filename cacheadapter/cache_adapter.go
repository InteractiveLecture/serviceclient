package cacheadapter

import (
	"github.com/InteractiveLecture/servicecache"
	"time"
)

type CacheAdapter struct {
	intervall     time.Duration
	timeout       time.Duration
	consulAddress string
	maxRetries    int
}

func New(address string, intervall time.Duration, timeout time.Duration, maxRetries int) *CacheAdapter {
	return &CacheAdapter{intervall, timeout, address, maxRetries}
}

func (adapter *CacheAdapter) Resolve(name string) (string, error) {
	return servicecache.GetServiceAddress(name)
}

func (adapter *CacheAdapter) Configure(services ...string) error {
	_, _ = servicecache.Configure(adapter.consulAddress, adapter.intervall, services...)
	return servicecache.Start(adapter.maxRetries, adapter.timeout)
}

func (adapter *CacheAdapter) Refresh() error {
	return servicecache.RefreshAndRestart()
}
