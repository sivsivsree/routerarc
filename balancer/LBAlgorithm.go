package balancer

import "net/url"

// LBAlgorithm is used to
type LBAlgorithm interface {
	Apply() (*url.URL, error)
}
