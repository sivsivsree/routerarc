package balancer

import (
	"math/rand"
	"net/url"
	"sync"
	"time"
)

//type random struct {
//	rand *rand.Rand
//}
//
//func newRandom(endpoints []string) *random {
//	t := time.Now().UnixNano()
//	return &random{rand: rand.New(rand.NewSource(t))}
//}

type roundRobin struct {
	mx        sync.Mutex
	index     int
	endpoints []string
}

func NewRoundRobin(endpoints []string) LBAlgorithm {
	i := time.Now().UnixNano()
	rand.Seed(i)
	return &roundRobin{
		index:     rand.Intn(len(endpoints)),
		endpoints: endpoints,
	}
}

// Apply implements routing.LBAlgorithm with a roundrobin algorithm.
func (r *roundRobin) Apply() (*url.URL, error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	r.index = (r.index + 1) % len(r.endpoints)
	return url.Parse(r.endpoints[r.index])
}
