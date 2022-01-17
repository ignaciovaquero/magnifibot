package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/igvaquero18/magnifibot/archimadrid"
)

func main() {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": "localhost:6379", // TODO: Parameterize
		},
	})
	ca := cache.New(&cache.Options{
		Redis:      ring,
		LocalCache: cache.NewTinyLFU(10, 24*time.Hour), // TODO: Parameterize
	})
	client := archimadrid.NewClient(archimadrid.SetCache(ca))
	gospel, err := client.GetGospel(time.Now())
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(gospel.Day)
	fmt.Println(gospel.Title)
	fmt.Println(gospel.Reference)
	fmt.Println(gospel.Content)
}
