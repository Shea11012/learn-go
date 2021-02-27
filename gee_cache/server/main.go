package main

import (
	"flag"
	"fmt"
	"gee_cache"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":"640",
	"Jack":"589",
	"Sam":"378",
}

func createGroup() *gee_cache.Group {
	return gee_cache.NewGroup("scores",2<<10,gee_cache.GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key",key)
		if v,ok := db[key];ok {
			return []byte(v),nil
		}

		return nil,fmt.Errorf("%s not exist",key)
	}))
}

func startCacheServer(addr string, addrs []string, gee *gee_cache.Group) {
	peers := gee_cache.NewHTTPPool(addr)
	peers.Set(addrs...)
	gee.RegisterPeers(peers)
	log.Println("geecache is running at",addr)
	log.Fatal(http.ListenAndServe(addr[7:],peers))
}

func startAPIServer(apiAddr string, gee *gee_cache.Group) {
	http.Handle("/api",http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		vb,err := gee.Get(key)
		if err != nil {
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type","application/octet-stream")
		w.Write(vb.ByteSlice())
	}))

	log.Println("fronted server is running at",apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:],nil))
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port,"port",8001,"geecache server port")
	flag.BoolVar(&api,"api",false,"start a api server?")
	flag.Parse()
	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001:"http://localhost:8001",
		8002:"http://localhost:8002",
		8003:"http://localhost:8003",
	}

	var addrs []string
	for _,v := range addrMap {
		addrs = append(addrs,v)
	}
	gee := createGroup()
	if api {
		go startAPIServer(apiAddr,gee)
	}

	startCacheServer(addrMap[port],addrs,gee)
}
