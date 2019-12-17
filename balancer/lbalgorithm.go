package balancer

import "net/url"

// ROUND_ROBIN algorithm used for load balancing
const ROUND_ROBIN = "round-robin"

// LBAlgorithm is used for generic load balancing algorithm implementation.
type LBAlgorithm interface {
	Apply() (*url.URL, error)
}
