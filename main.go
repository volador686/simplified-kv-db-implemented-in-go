package main

import (
	"flag"
	"log"
	"modules/config"
	"modules/db"
	"modules/web"
	"net/http"

	"github.com/BurntSushi/toml"
)

var (
	dbLocation = flag.String("db-location", "", "db database path")
	httpAddr   = flag.String("http-addr", "127.0.0.1:8080", "HTTP host&port")
	configFile = flag.String("config-file", "sharding.toml", "sharding config file")
	shard      = flag.String("shard", "", "choose a shard(database) to store the data")
)

func parseFlags() {
	flag.Parse()

	if *dbLocation == "" {
		log.Fatalf("Must provide db-location")
	}

	if *shard == "" {
		log.Fatalf("Must provide shard")
	}
}

func main() {
	// check flags
	parseFlags()

	var c config.Config
	if _, err := toml.DecodeFile(*configFile, &c); err != nil {
		log.Fatalf("toml.DecoedeFile(%q) : %v", *configFile, err)
	}

	var shardCount int
	var shardIndex int = -1
	var addrs = make(map[int]string)

	shardCount = len(c.Shards)
	for _, s := range c.Shards {
		addrs[s.Idx] = s.Address
		if s.Name == *shard {
			shardIndex = s.Idx
		}
	}
	// fmt.Printf("addrs: [0] = %s; [1] = %s; [2] = %s\n", addrs[0], addrs[1], addrs[2])
	if shardIndex == -1 {
		log.Fatalf("Shard %q was not found", *shard)
	}

	log.Printf("shard count: %d, current shard: %d", shardCount, shardIndex)

	// connect to database
	db, closefunc, err := db.NewDatabase(*dbLocation)
	if err != nil {
		log.Fatalf("New database(%q) : %v", *dbLocation, err)
	}

	// close function
	defer closefunc()

	// setup get&set function
	// fmt.Printf("server: [0] = %s; [1] = %s; [2] = %s\n", addrs[0], addrs[1], addrs[2])
	srv := web.NewServer(db, shardIndex, shardCount, addrs)
	http.HandleFunc("/get", srv.GetHandler)
	http.HandleFunc("/set", srv.SetHandler)
	http.HandleFunc("/replica", srv.ReplicaHandler)
	// setup http listener
	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
