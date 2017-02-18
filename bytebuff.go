package gotiny

type byteBuf struct {
	buf    []byte //编码目的数组
	offset int    //从buf[offset]开始写数据
	index  int    //下一个要写的下标
	length int    //len(buf)
	pool   Pool   //byte切片池
}

func NewByteBuf(p Pool, off, minlen int) *byteBuf {
	b := &byteBuf{
		pool:   p,
		offset: off,
		index:  off,
	}
	b.setNBytes(off + minlen)
	return b
}

func (b *byteBuf) Reset(off, minlen int) {
	b.setNBytes(off + minlen)
	b.offset = off
	b.index = off
}

func (b *byteBuf) Bytes() []byte {
	return b.buf[:b.index]
}

func (b *byteBuf) Len() int {
	return b.length
}

func (b *byteBuf) Write(p []byte) {
	b.index += copy(b.buf[b.index:], p)
}

func (b *byteBuf) WriteString(s string) {
	b.index += copy(b.buf[b.index:], s)
}

func (b *byteBuf) WriteByte(c byte) {
	b.buf[b.index] = c
	b.index++
}

func (b *byteBuf) SetBit(pos int, bit byte) {
	b.buf[pos] |= bit
}

func (b *byteBuf) Expand(n int) {
	b.setNBytes(b.index + n)
}

func (b *byteBuf) setNBytes(length int) {
	if len(b.buf) < length {
		buf := b.pool.Get(length)
		b.pool.Put(b.buf)
		b.buf = buf[:cap(buf)]
	}
	b.length = length
}
