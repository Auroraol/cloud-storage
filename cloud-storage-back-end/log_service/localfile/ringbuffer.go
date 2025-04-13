package localfile

// RingBuffer 实现一个简单的环形缓冲区，用于存储最近的几行
type RingBuffer struct {
	data  []string
	size  int
	start int
	count int
}

// NewRingBuffer 创建一个新的环形缓冲区
func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		data:  make([]string, size),
		size:  size,
		start: 0,
		count: 0,
	}
}

// Add 添加一个元素到环形缓冲区
func (rb *RingBuffer) Add(item string) {
	pos := (rb.start + rb.count) % rb.size

	if rb.count < rb.size {
		rb.data[pos] = item
		rb.count++
	} else {
		rb.data[rb.start] = item
		rb.start = (rb.start + 1) % rb.size
	}
}

// GetAll 获取环形缓冲区中的所有元素，按添加顺序排列
func (rb *RingBuffer) GetAll() []string {
	result := make([]string, rb.count)

	for i := 0; i < rb.count; i++ {
		pos := (rb.start + i) % rb.size
		result[i] = rb.data[pos]
	}

	return result
}
