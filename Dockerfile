FROM ckeyer/go AS building

COPY . /go/src/github.com/ckeyer/attack
WORKDIR /go/src/github.com/ckeyer/attack

RUN make test || exit 1 ;\
	make build

FROM alpine:edge

COPY --from=building /go/src/github.com/ckeyer/attack/bundles/attack /usr/local/bin/attack

ENTRYPOINT ["/usr/local/bin/attack"]
