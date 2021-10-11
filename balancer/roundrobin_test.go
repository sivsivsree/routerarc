package balancer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewRoundRobin is test of NewRoundRobin()
func TestNewRoundRobin(t *testing.T) {
	endpoints := []string{"http://aa.aa", "http://bb.bb", "http://cc.cc", "http://dd.dd"}
	rr := NewRoundRobin(endpoints).(*roundRobin)
	for i, endpoint := range rr.endpoints {
		assert.Equal(t, endpoints[i], endpoint)
	}
}

// TestApply is test of Apply() with a roundrobin algorithm.
func TestApply(t *testing.T) {
	endpoints := []string{"http://aa.aa", "http://bb.bb", "http://cc.cc", "http://dd.dd"}
	rr := NewRoundRobin(endpoints).(*roundRobin)
	prevIndex := rr.index
	url, err := rr.Apply()
	assert.Nil(t, err)
	assert.Equal(t, endpoints[(prevIndex+1)%len(endpoints)], url.String())
}
