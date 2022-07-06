This file is created to record the function of every part of this project.

cmd/bench:
    encompasses two parts
    read is for testing the speed of random read
    write is for testing the speed of random write
    not official benchmark test
config:
    the structs used amid decoding the sharding.toml
db:
    functions uesd to create & setup database
web:
    redirection
    copy http body
crane.db, whale.db, volador.db & panda.db:
    bolt database file
    contain key-value pairs
go.mod & go.sum:
    go modules files
launch.sh:
    shell process to activate four database
main.go:
    receive flags and create database, activate http functions, like: get, set, listenandserve
populate.sh:
    random input shell process
benchmark:
1.write_test.go
result:
Running tool: /usr/local/go/bin/go test -benchmem -run=^$ -bench ^BenchmarkWrite$ modules/benchmark

goos: linux
goarch: amd64
pkg: modules/benchmark
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkWrite-12    	      49	  24769037 ns/op	   14965 B/op	     116 allocs/op
PASS
ok  	modules/benchmark	2.163s
2.read_test.go
result:
Running tool: /usr/local/go/bin/go test -benchmem -run=^$ -bench ^BenchmarkRead$ modules/benchmark

goos: linux
goarch: amd64
pkg: modules/benchmark
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkRead-12    	    4815	    233616 ns/op	   14687 B/op	     112 allocs/op
PASS
ok  	modules/benchmark	1.156s

about replication:
in my opinion, replications is used to cope with the situations like database deleted accidently/disk failure...
I designed the replication mechanism in another way:
(replication with get operation should be designed into a cache, not a replication.)
first, create a extra bucket to store all the operations applied to the main database
second, create func() { apply all the operations to the replica database }
