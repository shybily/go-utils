package hashring

import (
	"fmt"
	"hash/crc32"
	"sort"
	"sync"
)

type nodeIdx []uint32

func (idx nodeIdx) Len() int {
	return len(idx)
}

func (idx nodeIdx) Swap(i, j int) {
	idx[i], idx[j] = idx[j], idx[i]
}

func (idx nodeIdx) Less(i, j int) bool {
	return idx[i] <= idx[j]
}

type Ring struct {
	nodes        map[uint32]string
	idx          nodeIdx
	replicaCount int
	mu           sync.RWMutex
}

func NewRing(replicaCount int, nodes []string) *Ring {
	r := &Ring{
		nodes:        make(map[uint32]string, replicaCount*len(nodes)),
		idx:          make(nodeIdx, replicaCount*len(nodes)),
		replicaCount: replicaCount,
	}
	k := 0
	for _, node := range nodes {
		for i := 0; i < r.replicaCount; i++ {
			h := getHash([]byte(fmt.Sprintf("%s:%d", node, i)))

			r.idx[k] = h
			r.nodes[h] = node
			k++
		}
	}
	sort.Sort(r.idx)
	return r
}

func getHash(key []byte) uint32 {
	return crc32.ChecksumIEEE(key)
}

func (r *Ring) addNode(node string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := 0; i < r.replicaCount; i++ {
		h := getHash([]byte(fmt.Sprintf("%s:%d", node, i)))

		r.idx = append(r.idx, h)
		r.nodes[h] = node
	}

	sort.Sort(r.idx)
	return nil
}

func getKeys(m map[uint32]string) (idx nodeIdx) {
	for k := range m {
		idx = append(idx, k)
	}

	return idx
}

func (r *Ring) Locate(key string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	h := getHash([]byte(key))

	pos := sort.Search(len(r.idx), func(i int) bool { return r.idx[i] >= h })
	if pos == len(r.idx) {
		return r.nodes[r.idx[0]]
	}
	return r.nodes[r.idx[pos]]
}
