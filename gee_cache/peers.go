package gee_cache

import (
	pb "gee_cache/geecachepb"
)

type PeerGetter interface {
	Get(*pb.Request,*pb.Response) error
}

type PeerPicker interface {
	PickPeer(string) (PeerGetter,bool)
}