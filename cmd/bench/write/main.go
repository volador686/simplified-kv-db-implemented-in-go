package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

var addr = flag.String("addr", "localhost:8080", "the HTTP host port for the instance that is banchmarked.")

func writeRand() {
	key := fmt.Sprintf("key-%d", rand.Intn(100000))
	value := fmt.Sprintf("value-%d", rand.Intn(100000))

	values := url.Values{}
	values.Set("key", key)
	values.Set("value", value)
	resp, err := http.Get("http://" + (*addr) + "/set?" + values.Encode())
	if err != nil {
		log.Fatalf("Error during set: %v", err)
	}

	defer resp.Body.Close()

}

func benchmark(function func()) time.Duration {
	start := time.Now()
	function()
	return time.Since(start)
}

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	var time_cost time.Duration = 0

	for i := 0; i < 10000; i++ {
		time_cost = time_cost + benchmark(writeRand)
	}

	fmt.Printf("10000 times average:%s\n", (time_cost / 10000))
	fmt.Printf("10000 times total:%s\n", time_cost)
}
