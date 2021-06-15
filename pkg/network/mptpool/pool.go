package mptpool

import (
	"sync"

	"github.com/nspcc-dev/neo-go/pkg/util"
)

// Pool stores unknown MPT nodes along with the corresponding paths.
type Pool struct {
	lock   sync.RWMutex
	hashes map[util.Uint256]*Item
}

// Item stores MPT node's path.
type Item struct {
	Path []byte
}

// New returns new MPT node hashes pool using provided chain.
func New() *Pool {
	return &Pool{
		hashes: make(map[util.Uint256]*Item),
	}
}

// ContainsKey checks if an MPT node hash is in the Pool.
func (mp *Pool) ContainsKey(hash util.Uint256) bool {
	mp.lock.RLock()
	defer mp.lock.RUnlock()

	return mp.containsKey(hash)
}

func (mp *Pool) containsKey(hash util.Uint256) bool {
	_, ok := mp.hashes[hash]
	return ok
}

// TryGet returns MPT path for the specified HashNode.
func (mp *Pool) TryGet(hash util.Uint256) (*Item, bool) {
	mp.lock.RLock()
	defer mp.lock.RUnlock()

	itm, ok := mp.hashes[hash]
	return itm, ok
}

// Remove removes item from the pool by the specified hash.
func (mp *Pool) Remove(hash util.Uint256) {
	mp.lock.Lock()
	defer mp.lock.Unlock()

	mp.remove(hash)
}

func (mp *Pool) remove(hash util.Uint256) {
	if mp.containsKey(hash) {
		delete(mp.hashes, hash)
	}
}

// Add adds item to the pool.
func (mp *Pool) Add(hash util.Uint256, item *Item) {
	mp.lock.Lock()
	defer mp.lock.Unlock()

	mp.hashes[hash] = item
}

// Update removes/adds specified items from/to the pool.
func (mp *Pool) Update(remove map[util.Uint256]bool, add map[util.Uint256]*Item) {
	mp.lock.Lock()
	defer mp.lock.Unlock()

	for h := range remove {
		mp.remove(h)
	}
	for h, itm := range add {
		mp.hashes[h] = itm
	}
}

// Count returns the number of items in the pool.
func (mp *Pool) Count() int {
	mp.lock.RLock()
	defer mp.lock.RUnlock()

	return len(mp.hashes)
}
