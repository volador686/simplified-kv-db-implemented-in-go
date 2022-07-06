package config

// describes a shard that holds unique set of keys.
type Shard struct {
	Name         string
	Idx          int
	Address      string
	Replica_addr string
}

// config describes the sharding configuartion
type Config struct {
	Shards []Shard
}
