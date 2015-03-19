package main

import (
	"github.com/coreos/go-etcd/etcd"
	"github.com/rcrowley/go-tigertonic"
)

func main() {
	mux := tigertonic.NewTrieServeMux()
	server := tigertonic.NewServer(":3000", tigertonic.ApacheLogged(mux))

	etcdClient := etcd.NewClient([]string{"http://127.0.0.1:4001"})
	backendsResource := &BackendsResource{etcdClient: etcdClient}
	mux.Handle("GET", "/api/backends", tigertonic.Marshaled(backendsResource.Index))

	hostsResource := &HostsResource{etcdClient: etcdClient}
	mux.Handle("GET", "/api/hosts", tigertonic.Marshaled(hostsResource.Index))

	server.ListenAndServe()
}
