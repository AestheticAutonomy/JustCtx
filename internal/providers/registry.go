package providers

import (
	"fmt"
	"sync"
)

var (
	providersMu sync.RWMutex
	registry    = make(map[string]Provider)
)

func Register(p Provider) {
	providersMu.Lock()
	defer providersMu.Unlock()
	if p == nil {
		panic("provider is nil")
	}
	name := p.Name()
	if _, dup := registry[name]; dup {
		panic("provider already registered: " + name)
	}
	registry[name] = p
}

func Get(name string) (Provider, error) {
	providersMu.RLock()
	defer providersMu.RUnlock()
	p, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("provider not found: %s", name)
	}
	return p, nil
}

func All() []Provider {
	providersMu.RLock()
	defer providersMu.RUnlock()
	list := make([]Provider, 0, len(registry))
	for _, p := range registry {
		list = append(list, p)
	}
	return list
}
