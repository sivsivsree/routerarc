package balancer

import "net/url"

const ROUND_ROBIN = "round-robin"

// LBAlgorithm is used to
type LBAlgorithm interface {
	Apply() (*url.URL, error)
}
