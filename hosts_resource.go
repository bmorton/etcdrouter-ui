package main

import (
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/coreos/go-etcd/etcd"
)

type HostsResource struct {
	etcdClient *etcd.Client
}

type Host struct {
	Name      string      `json:"name"`
	Locations []*Location `json:"locations"`
}

type Location struct {
	Name    string   `json:"name"`
	Path    string   `json:"path"`
	Backend *Backend `json:"backend"`
}

type HostsResponse struct {
	Hosts []*Host `json:"hosts"`
}

func (hr *HostsResource) Index(u *url.URL, h http.Header, req interface{}) (int, http.Header, *HostsResponse, error) {
	response := &HostsResponse{Hosts: []*Host{}}

	rawHosts, err := hr.etcdClient.Get("/vulcand/hosts", true, false)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}

	for _, node := range rawHosts.Node.Nodes {
		host := &Host{
			Name:      path.Base(node.Key),
			Locations: []*Location{},
		}
		rawLocations, err := hr.etcdClient.Get(fmt.Sprintf("/vulcand/hosts/%s/locations", host.Name), true, false)
		if err != nil {
			return http.StatusInternalServerError, nil, nil, err
		}

		for _, locNode := range rawLocations.Node.Nodes {
			location := &Location{Name: path.Base(locNode.Key)}
			rawPath, err := hr.etcdClient.Get(fmt.Sprintf("/vulcand/hosts/%s/locations/%s/path", host.Name, location.Name), false, false)
			if err != nil {
				return http.StatusInternalServerError, nil, nil, err
			}
			location.Path = rawPath.Node.Value
			rawBackend, err := hr.etcdClient.Get(fmt.Sprintf("/vulcand/hosts/%s/locations/%s/upstream", host.Name, location.Name), false, false)
			if err != nil {
				return http.StatusInternalServerError, nil, nil, err
			}
			location.Backend = &Backend{Name: rawBackend.Node.Value}

			host.Locations = append(host.Locations, location)
		}
		response.Hosts = append(response.Hosts, host)
	}

	return http.StatusOK, nil, response, nil
}
