package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/coreos/go-etcd/etcd"
	"github.com/k0kubun/pp"
	"github.com/rcrowley/go-tigertonic"
)

var staticRoot string

func main() {
	staticRoot = os.Getenv("STATIC_ROOT")
	if staticRoot != "" {
		staticRoot = fmt.Sprintf("%s/", staticRoot)
	}
	etcdHosts := os.Getenv("ETCDCTL_PEERS")
	if etcdHosts == "" {
		etcdHosts = "http://127.0.0.1:4001"
	}

	mux := tigertonic.NewTrieServeMux()
	server := tigertonic.NewServer(":3000", tigertonic.ApacheLogged(mux))

	etcdClient := etcd.NewClient([]string{etcdHosts})
	backendsResource := &BackendsResource{etcdClient: etcdClient}
	mux.Handle("GET", "/api/backends", tigertonic.Marshaled(backendsResource.Index))

	hostsResource := &HostsResource{etcdClient: etcdClient}
	mux.Handle("GET", "/api/hosts", tigertonic.Marshaled(hostsResource.Index))

	mux.HandleNamespace("", http.HandlerFunc(indexHandler))
	mux.HandleNamespace("static", http.HandlerFunc(staticHandler))

	server.ListenAndServe()
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	filename := fmt.Sprintf("%sstatic/%s", staticRoot, r.URL.Path[1:])
	pp.Println(filename)

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	} else {
		http.ServeFile(w, r, filename)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	filename := fmt.Sprintf("%sstatic/index.html", staticRoot)
	pp.Println(filename)
	http.ServeFile(w, r, filename)
}
