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

func ReadRand() {
	key := fmt.Sprintf("key-%d", rand.Intn(100000))
	values := url.Values{}
	values.Set("key", key)
	resp, err := http.Get("http://" + (addr) + "/get?" + values.Encode())
	if err != nil {
		log.Fatalf("Error during get: %v", err)
	}

	defer resp.Body.Close()

}

func BenchmarkRead(b *testing.B) {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ReadRand()
	}
}
