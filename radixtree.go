package redirector

const (
	// IPV4TYP ipv4 type
	IPV4TYP = iota
	// IPV6TYP ipv6 type
	IPV6TYP
	// DOMAINTYP domain type
	DOMAINTYP
)

type imapNode struct {
	key int
	val int
}

//
type imap struct {
	size int
	arr  []imapNode
}

func (i *imap) get(b int) (val int, ok bool) {
	size := len(i.arr)
	for factor := 0; factor < size; factor++ {
		if i.arr[(b+factor)%size].key == b {
			return i.arr[(b+factor)%size].val, true
		}
	}
	return -1, false
}

func (i *imap) exist(b int) (idx int, ok bool) {
	size := len(i.arr)
	for factor := 0; factor < size; factor++ {
		if i.arr[(b+factor)%size].key == b {
			return (b + factor) % size, true
		}
	}
	return -1, false
}

func (i *imap) set(b int, idx int) {
	if index, ok := i.exist(b); ok {
		i.arr[idx].val = index
		return
	}
	size := len(i.arr)
	load := float32(i.size) / float32(size)
	if load < 0.8 {
		for factor := 0; factor < size; factor++ {
			if i.arr[(b+factor)%size].key == -1 {
				i.arr[(b+factor)%size].key = b
				i.arr[(b+factor)%size].val = idx
				i.size++
				return
			}
		}
	} else {
		tmp := i.arr
		i.arr = make([]imapNode, 2*size, 2*size)
		for k := 0; k < 2*size; k++ {
			i.arr[k].key = -1
		}
		for j := 0; j < size; j++ {
			if tmp[j].key != -1 {
				i.set(tmp[j].key, tmp[j].val)
			}
		}
		return
	}
}

func newimap() *imap {
	i := &imap{
		size: 0,
		arr:  make([]imapNode, 8),
	}
	for j := 0; j < 8; j++ {
		i.arr[j].key = -1
	}
	return i
}

// RadixTree compressed radix tree
type RadixTree interface {
	AddRoute(string, interface{}) error
	FindRoute(string) (interface{}, bool)
}

const (
	mapOnThresHold = 32
)

// RadixNode compressed radix tree
type RadixNode struct {
	path    string
	child   []*RadixNode
	indices []byte
	handle  interface{}
	m       []int
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func (r *RadixNode) insertChild(path string) {
	r.path = path
}

// AddRoute add route binded with hanle to the tree,
// acceptable parameter type: net.IP, []byte, string, uint32
func (r *RadixNode) AddRoute(addr string, handle interface{}) error {
	if len(r.path) > 0 || len(r.child) > 0 {
	walk:
		for {
			i := 0
			m := min(len(r.path), len(addr))
			for i < m && r.path[i] == addr[i] {
				i++
			}
			// split edge
			if i < len(r.path) {
				child := &RadixNode{
					path:    r.path[i:],
					handle:  r.handle,
					child:   r.child,
					indices: r.indices,
				}
				r.child = []*RadixNode{child}
				r.indices = []byte{r.path[i]}
				r.path = r.path[:i]
				r.handle = nil
			}

			if i < len(addr) {
				addr = addr[i:]
				indice := addr[0]
				length := len(r.indices)
				for i := 0; i < length; i++ {
					if indice == r.indices[i] {
						r = r.child[i]
						continue walk
					}
				}
				r.indices = append(r.indices, indice)
				child := &RadixNode{}
				r.child = append(r.child, child)
				r = child
				r.insertChild(addr)
				return nil
			} else if i == len(addr) {
				if r.handle != nil {
					panic("handler already registered")
				}
				r.handle = handle
			}
			return nil
		}
	} else {
		r.insertChild(addr)
	}
	return nil
}

// FindRoute find route in radix tree, returns handler bind to the path and
// whether the path exists
func (r *RadixNode) FindRoute(path string) (handler interface{}, ok bool) {
walk:
	for {
		if len(path) > len(r.path) {
			if path[:len(r.path)] == r.path {
				path = path[len(r.path):]
				c := path[0]
				length := len(r.indices)
				for i := 0; i < length; i++ {
					if c == r.indices[i] {
						r = r.child[i]
						continue walk
					}
				}

				return r.handle, false
			}
		} else if path == r.path {
			return r.handle, true
		}
		return nil, false
	}
}

// NewRadixTree returns an empty radix tree node
func NewRadixTree() *RadixNode {
	return &RadixNode{
		m: make([]int, 256),
	}
}
