// Package resolver implements an SRV resolver for gRPC.
package resolver

import (
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc/attributes"

	"google.golang.org/grpc/resolver"
)

// Builder implements resolver.Builder.
type Builder struct {
}

// Build implements resolver.Builder.Build.
func (b *Builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &Resolver{
		target:    target,
		cc:        cc,
		rn:        make(chan struct{}, 1),
		close:     make(chan struct{}, 0),
		closeOnce: sync.Once{},
	}
	go r.watch()
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}

// Scheme implements resolver.Builder.Scheme.
func (b *Builder) Scheme() string {
	return "dns-srv"
}

// Resolver implements resolver.Resolve.
type Resolver struct {
	target    resolver.Target
	cc        resolver.ClientConn
	rn        chan struct{}
	close     chan struct{}
	closeOnce sync.Once
}

func (r *Resolver) watch() {
	for {
		select {
		case <-time.After(time.Minute):
		case <-r.rn:
		case <-r.close:
			return
		}
		s := strings.Split(r.target.Endpoint, "|")
		service, proto, name := s[0], s[1], s[2]
		_, srvs, err := net.LookupSRV(service, proto, name)
		if err != nil {
			r.cc.ReportError(err)
		} else {
			addrs := make([]resolver.Address, 0, len(srvs))
			for _, srv := range srvs {
				addrs = append(addrs, resolver.Address{
					Addr: net.JoinHostPort(srv.Target, strconv.Itoa(int(srv.Port))),
					Attributes: attributes.New(
						"priority", srv.Priority,
						"weight", srv.Weight,
					),
				})
			}
			r.cc.UpdateState(resolver.State{
				Addresses: addrs,
			})
		}
	}
}

// ResolveNow implements resolver.Resolve.ResolveNow.
func (r *Resolver) ResolveNow(opts resolver.ResolveNowOptions) {
	select {
	case r.rn <- struct{}{}:
	default:
	}
}

// Close implements resolver.Resolve.Close.
func (r *Resolver) Close() {
	r.closeOnce.Do(func() { close(r.close) })
}

func init() {
	resolver.Register(&Builder{})
}
