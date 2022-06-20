module github.com/tikv/pd/client

go 1.16

require (
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pingcap/errors v0.11.5-0.20211224045212-9687c2b0f87c
	github.com/pingcap/failpoint v0.0.0-20210918120811-547c13e3eb00
	github.com/pingcap/kvproto v0.0.0-20220510035547-0e2f26c0a46a
	github.com/pingcap/log v0.0.0-20211215031037-e024ba4eb0ee
	github.com/prometheus/client_golang v1.11.0
	github.com/stretchr/testify v1.7.0
	github.com/tikv/pd v0.0.0-00010101000000-000000000000
	go.uber.org/goleak v1.1.12
	go.uber.org/zap v1.20.0
	google.golang.org/grpc v1.43.0
)

replace github.com/tikv/pd => ../

// reset grpc and protobuf deps in order to import client and server at the same time
replace (
	github.com/golang/protobuf v1.5.2 => github.com/golang/protobuf v1.3.4
	google.golang.org/grpc v1.43.0 => google.golang.org/grpc v1.26.0
	google.golang.org/protobuf v1.26.0 => github.com/golang/protobuf v1.3.4
)
