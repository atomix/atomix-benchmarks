FROM alpine:3.8

RUN apk upgrade --update --no-cache && apk add libc6-compat

USER nobody

ADD build/_output/k8s-benchmarks /usr/local/bin/k8s-benchmarks

ENTRYPOINT ["k8s-benchmarks"]
