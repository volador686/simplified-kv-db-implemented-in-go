#!/bin/bash


go run main.go -db-location=volador.db -http-addr=127.0.0.1:8080 -config-file=sharding.toml -shard=volador &
# go run main.go -db-location=volador-r.db -http-addr=127.0.0.1:8180 -config-file=sharding.toml -shard=volador -replica &

go run main.go -db-location=crane.db -http-addr=127.0.0.1:8081 -config-file=sharding.toml -shard=crane &
# go run main.go -db-location=crane-r.db -http-addr=127.0.0.1:8181 -config-file=sharding.toml -shard=crane -replica &

go run main.go -db-location=panda.db -http-addr=127.0.0.1:8082 -config-file=sharding.toml -shard=panda &
# go run main.go -db-location=panda-r.db -htpp-addr=127.0.0.1:8182 -config-file=sharding.toml -shard=panda -replica &

go run main.go -db-location=whale.db -http-addr=127.0.0.1:8083 -config-file=sharding.toml -shard=whale &
# go run main.go -dn-location=whale-r.db -http-addr=127.0.0.1:8183 -config-file=sharding.toml -shard=whale -replica &

wait