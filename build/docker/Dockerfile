FROM alpine:3.8

RUN apk upgrade --update --no-cache && apk add libc6-compat

USER nobody

ADD build/_output/atomix-benchmarks /usr/local/bin/atomix-benchmarks

ENTRYPOINT ["atomix-benchmarks"]
