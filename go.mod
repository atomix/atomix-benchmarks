module github.com/atomix/atomix-benchmarks

go 1.13

require (
	github.com/atomix/atomix-go-client v0.6.1
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/onosproject/helmit v0.6.15
)

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20200229013735-71373c6105e3

replace github.com/onosproject/helmit => ../helmet
