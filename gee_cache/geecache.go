package gee_cache

import (
	"errors"
	"gee_cache/singleflight"
	"log"
	"sync"
	pb "gee_cache/geecachepb"
)

type Group struct {
	name string
	getter Getter
	mainCache cache
	peers PeerPicker
	loader *singleflight.Group
}

var (
	mu sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string,cacheBytes int64,getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}

	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name: name,
		getter: getter,
		mainCache: cache{
			cacheBytes: cacheBytes,
		},
		loader: &singleflight.Group{},
	}

	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("registerPeerPicker called more than once")
	}

	g.peers = peers
}

func (g *Group) load(key string) (ByteView, error) {
	bv,err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				var value ByteView
				var err error
				if value, err = g.getFromPeer(peer, key); err != nil {
					return value, nil
				}
				log.Println("[GeeCache] Failed to get from peer", err)
			}
		}

		return g.getLocally(key)
	})

	if err != nil {
		return ByteView{},err
	}

	return bv.(ByteView),nil
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, errors.New("key is required")
	}
	
	if v,ok := g.mainCache.get(key);ok {
		log.Println("hit")
		return v,nil
	}
	
	return g.load(key)
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	req := &pb.Request{
		Group: g.name,
		Key: key,
	}
	res := &pb.Response{}
	err := peer.Get(req,res)
	if err != nil {
		return ByteView{},err
	}

	return ByteView{b:res.Value},nil
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes,err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key,value)
	return value,nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key,value)
}
