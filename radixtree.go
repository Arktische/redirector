package redirector

const (
	// IPV4TYP ipv4 type
	IPV4TYP = iota
	// IPV6TYP ipv6 type
	IPV6TYP
	// DOMAINTYP domain type
	DOMAINTYP
)

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
	return &RadixNode{}
}
