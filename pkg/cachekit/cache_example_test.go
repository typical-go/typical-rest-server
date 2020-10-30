package cachekit_test

import (
	"fmt"
	"log"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	"github.com/typical-go/typical-rest-server/pkg/cachekit"
)

func ExampleCache() {

	// run redis server
	server, err := miniredis.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer server.Close()

	// create redis client
	client := redis.NewClient(&redis.Options{Addr: server.Addr()})

	// define key and refresh function for your cache
	cache := cachekit.Cache{
		Client: client,
		Key:    "some-key",
		RefreshFn: func() (interface{}, error) {
			return "fresh-data", nil
		},
	}

	// execute cache to get the data
	var data string
	if err = cache.Execute(&data, pragmaWithCacheControl("")); err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(data)

	// Output:
	// fresh-data

}
