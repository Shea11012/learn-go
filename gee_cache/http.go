package gee_cache

import (
	"fmt"
	"gee_cache/consistenthash"
	pb "gee_cache/geecachepb"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/golang/protobuf/proto"
)

const (
	defaultBasePath = "/_geecache/"
	defaultReplicas = 50
)

type HTTPPool struct {
	self     string	// 记录自己的地址
	basePath string	// 节点通讯地址前缀，用以区分其他HTTP接口
	mu sync.Mutex
	peers *consistenthash.Map
	httpGetters map[string]*httpGetter
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// Set 实例化一致性哈希算法，添加传入节点
func (h *HTTPPool) Set(peers ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.peers = consistenthash.New(defaultReplicas,nil)
	h.peers.Add(peers...)
	h.httpGetters = make(map[string]*httpGetter,len(peers))
	for _,peer := range peers {
		h.httpGetters[peer] = &httpGetter{baseURL: peer + h.basePath}
	}
}

// PickPeer 根据key返回对应的 http客户端
func (h *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if peer := h.peers.Get(key); peer != "" && peer != h.self {
		h.Log("pick peer %s",peer)
		return h.httpGetters[peer],true
	}

	return nil, false
}

func (h *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", h.self, fmt.Sprintf(format, v...))
}

// ServeHTTP url 格式: /basepath/groupname/key 
// 判断url是否含有basepath前缀，根据不同的groupname得到group实例,查询对应的key
func (h *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, h.basePath) {
		panic("HTTPPool serving unexpected path:" + r.URL.Path)
	}

	h.Log("%s %s", r.Method, r.URL.Path)
	// _basepath/group_name/key
	parts := strings.SplitN(r.URL.Path[len(h.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName, key := parts[0], parts[1]
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such gorup:"+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body,err := proto.Marshal(&pb.Response{Value: view.ByteSlice()})
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(body)
}

type httpGetter struct {
	baseURL string	// 要访问的远程地址
}

// Get 从远程节点获取值
func (h *httpGetter) Get(req *pb.Request, res *pb.Response) error {
	u := fmt.Sprintf(
		"%v%v/%v",
		h.baseURL,
		url.QueryEscape(req.Group),
		url.QueryEscape(req.Key),
	)

	response,err := http.Get(u)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %v",response.Status)
	}

	bytes,err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v",err)
	}

	if err = proto.Unmarshal(bytes,res); err != nil {
		return fmt.Errorf("decoding response body: %v",err)
	}

	return nil
}

var _ PeerGetter = (*httpGetter)(nil)
var _ PeerPicker = (*HTTPPool)(nil)
