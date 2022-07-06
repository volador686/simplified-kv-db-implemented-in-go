package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"testing"
	"time"
)

var addr = "localhost:8080"

func WriteRand() {
	key := fmt.Sprintf("key-%d", rand.Intn(100000))
	value := fmt.Sprintf("value-%d", rand.Intn(100000))

	values := url.Values{}
	values.Set("key", key)
	values.Set("value", value)
	resp, err := http.Get("http://" + (addr) + "/set?" + values.Encode())
	if err != nil {
		log.Fatalf("Error during set: %v", err)
	}

	defer resp.Body.Close()

}

func BenchmarkWrite(b *testing.B) {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WriteRand()
	}
}
