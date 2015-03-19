FROM golang
MAINTAINER Brian Morton "bmorton@yammer-inc.com"

# Install etcd client
ENV ETCD_VERSION v2.0.0-rc.1
ADD https://github.com/coreos/etcd/releases/download/${ETCD_VERSION}/etcd-${ETCD_VERSION}-linux-amd64.tar.gz /tmp/
RUN tar -xvzf /tmp/etcd-${ETCD_VERSION}-linux-amd64.tar.gz -C /tmp
RUN mv /tmp/etcd-${ETCD_VERSION}-linux-amd64/etcdctl /bin/etcdctl
RUN rm -rf /tmp/*

ENV REPO github.com/bmorton/etcdrouter-ui
ADD . /go/src/$REPO
RUN mv /go/src/$REPO/env /go/bin/env
RUN chmod +x /go/bin/env
RUN cd /go/src/$REPO && go get -v -d
RUN go install $REPO

ENV STATIC_ROOT /go/src/$REPO

ENTRYPOINT ["/go/bin/env"]
CMD ["/go/bin/etcdrouter-ui"]
EXPOSE 3000
