package qcollector_docker_events

import (
	"github.com/moby/moby/api/types"
	"sync"
)

// Simple inventory to use here

type Inventory struct {
	mu sync.Mutex
	data map[string]types.ContainerJSON
}

func NewInventory() Inventory {
	return Inventory{
		data: map[string]types.ContainerJSON{},
	}
}

func (i *Inventory) SetItem(id string, cnt types.ContainerJSON) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.data[id] = cnt
}

func (i *Inventory) GetItem(id string) (cnt types.ContainerJSON, err error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	if cnt, ok := i.data[id]; !ok {
		return cnt, err
	}
	return
}