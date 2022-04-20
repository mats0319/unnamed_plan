module github.com/mats9693/unnamed_plan/services

go 1.17

require (
	github.com/go-pg/pg/v10 v10.10.5
	github.com/golang/mock v1.6.0
	github.com/mats9693/utils v0.0.0
	github.com/pkg/errors v0.9.1
	go.uber.org/zap v1.19.1
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.25.0
)

require (
	github.com/go-pg/zerochecker v0.2.0 // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/tmthrgd/go-hex v0.0.0-20190904060850-447a3041c3bc // indirect
	github.com/vmihailenco/bufpool v0.1.11 // indirect
	github.com/vmihailenco/msgpack/v5 v5.3.1 // indirect
	github.com/vmihailenco/tagparser v0.1.2 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b // indirect
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4 // indirect
	golang.org/x/sys v0.0.0-20210510120138-977fb7262007 // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	mellium.im/sasl v0.2.1 // indirect
)

replace github.com/mats9693/utils => ../../utils
