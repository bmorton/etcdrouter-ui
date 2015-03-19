package main

import (
	"net/http"
	"net/url"
	"path"

	"github.com/coreos/go-etcd/etcd"
)

type BackendsResource struct {
	etcdClient *etcd.Client
}

type Backend struct {
	Name string `json:"name"`
}

type BackendsResponse struct {
	Backends []*Backend `json:"backends"`
}

func (b *BackendsResource) Index(u *url.URL, h http.Header, req interface{}) (int, http.Header, *BackendsResponse, error) {
	response := &BackendsResponse{Backends: []*Backend{}}

	retrieved, err := b.etcdClient.Get("/haproxy/backends", true, false)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}

	for _, node := range retrieved.Node.Nodes {
		backend := &Backend{Name: path.Base(node.Key)}
		response.Backends = append(response.Backends, backend)
	}

	return http.StatusOK, nil, response, nil
}
