package redirector

// HashTrieNode hash trie tree node (also could be root)
type HashTrieNode struct {
	m       map[byte]int
	handler interface{}
	child   []*HashTrieNode
}

// AddRoute add route into HashTrieNode
func (h *HashTrieNode) AddRoute(addr string, handler interface{}) {
	length := len(addr)
	for i := 0; i < length; i++ {
		idx, ok := h.m[addr[i]]
		if ok {
			h = h.child[idx]
		} else {
			h.child = append(h.child, &HashTrieNode{
				m: make(map[byte]int),
			})
			h.m[addr[i]] = len(h.child) - 1
			h = h.child[len(h.child)-1]
		}
	}
	h.handler = handler
	return
}

// FindRoute find route in HashTrieNode
func (h *HashTrieNode) FindRoute(path string) (handler interface{}, ok bool) {
	length := len(path)
	for i := 0; i < length; i++ {
		idx, exist := h.m[path[i]]
		if exist {
			h = h.child[idx]
		} else {
			return
		}
	}
	handler = h.handler
	ok = true
	return
}

// NewHashTrieTree returns root HashTrieNode
func NewHashTrieTree() *HashTrieNode {
	return &HashTrieNode{
		m: make(map[byte]int),
	}
}
