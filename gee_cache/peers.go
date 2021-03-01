package gee_cache

import (
	pb "gee_cache/geecachepb"
)

// PeerGetter 从对应的group查找缓存值
type PeerGetter interface {
	Get(*pb.Request,*pb.Response) error
}

// PeerPicker 根据key查找对应的节点
type PeerPicker interface {
	PickPeer(string) (PeerGetter,bool)
}