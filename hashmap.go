package redirector

// IHashMap interface of hashmap
type IHashMap interface {
	Get(key int) (val int, ok bool)
	Set(key int, val int)
}

type hashMapNode struct {
	key  int
	val  int
	next *hashMapNode
}

// HashMap linked address hash map
type HashMap struct {
	size      int
	cap       int
	arr       [2][]*hashMapNode
	used      [2]int
	rehashIdx int
}

func (h *HashMap) hash(key int) int {
	key = ^key + (key << 15) // key = (key << 15) - key - 1;
	key = key ^ (key >> 12)
	key = key + (key << 2)
	key = key ^ (key >> 4)
	key = key * 2057 // key = (key + (key << 3)) + (key << 11);
	key = key ^ (key >> 16)
	return key
}

func (h *HashMap) exist(key int, arrIdx int) (hashIdx int, idx int, val int, ok bool) {
	length := len(h.arr[arrIdx])
	hashIdx = h.hash(key) % length
	if h.arr[arrIdx][hashIdx] != nil {
		head := h.arr[arrIdx][hashIdx]
		for i := 0; head != nil; i++ {
			if head.key == key {
				ok = true
				val = head.val
				idx = i
				return
			}
			head = head.next
		}
	}
	return
}

// Get returns value and ok if key exists, otherwise returns (0,false)
func (h *HashMap) Get(key int) (val int, ok bool) {
	_, _, val, ok = h.exist(key, 0)
	if ok {
		return
	}
	if h.rehashIdx == 1 {
		_, _, val, ok = h.exist(key, 1)
	}
	return
}

func (h *HashMap) transfer() {
	length := len(h.arr[0])
	for i := 0; i < length; i++ {
		head := h.arr[0][i]
		if head != nil {
			h.set(head.key, head.val, 1)
			h.used[0]--
			h.used[1]++
			h.arr[0][i] = head.next
			return
		}
	}
}

// set insert key,val into the specific internal array in HashMap
func (h *HashMap) set(key int, val int, arrIdx int) {
	hashIdx, idx, _, ok := h.exist(key, arrIdx)
	if ok {
		head := h.arr[arrIdx][hashIdx]
		for i := 0; i < idx; i++ {
			head = head.next
		}
		head.val = val
		return
	}
	head := h.arr[arrIdx][hashIdx]
	newNode := &hashMapNode{key: key, val: val}
	if head == nil {
		h.arr[arrIdx][hashIdx] = newNode
		return
	}
	for head.next != nil {
		head = head.next
	}
	head.next = newNode
	return
}

// Set insert key,val into the map, if key already exists, only modify the old val
func (h *HashMap) Set(key int, val int) {
	if h.size >= h.cap-1 {
		h.arr[1] = make([]*hashMapNode, 2*h.cap)
		h.rehashIdx = 1
		h.cap *= 2
	}
	hashIdx, idx, _, ok := h.exist(key, 0)
	if ok {
		head := h.arr[0][hashIdx]
		for i := 0; i < idx; i++ {
			head = head.next
		}
		head.val = val
		return
	}
	if h.rehashIdx == 1 {
		hashIdx, idx, _, ok = h.exist(key, 1)
		if ok {
			head := h.arr[1][hashIdx]
			for i := 0; i < idx; i++ {
				head = head.next
			}
			head.val = val
			return
		}
		if h.used[0] > 0 {
			h.transfer()
		} else {
			h.rehashIdx = 0
			h.arr[0] = h.arr[1]
			h.used[0] = h.used[1]
			h.used[1] = 0
			h.arr[1] = h.arr[1][0:0]
		}
	}
	// insert new key-val pair according to h.rehashIdx
	newNode := &hashMapNode{key: key, val: val}
	head := h.arr[h.rehashIdx][hashIdx]
	if head == nil {
		h.arr[h.rehashIdx][hashIdx] = newNode
		h.size++
		h.used[h.rehashIdx]++
		return
	}
	for head.next != nil {
		head = head.next
	}
	head.next = newNode
	h.size++
	h.used[h.rehashIdx]++
	return
}

// NewHashMap returns a new hashmap
func NewHashMap() *HashMap {
	m := &HashMap{
		size: 0,
		cap:  8,
	}
	m.arr[0] = make([]*hashMapNode, 8)
	return m
}
