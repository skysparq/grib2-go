package data_representation

// BitStream for reading bits from a byte slice
type BitStream struct {
	data []byte
	pos  int // bit position
}

func NewBitStream(data []byte) *BitStream {
	return &BitStream{data: data, pos: 0}
}

func (b *BitStream) ReadBits(n int) uint64 {
	var val uint64
	for i := 0; i < n; i++ {
		byteIdx := b.pos / 8
		bitIdx := 7 - (b.pos % 8)
		bit := (b.data[byteIdx] >> bitIdx) & 1
		val = (val << 1) | uint64(bit)
		b.pos++
	}
	return val
}

func (b *BitStream) ReadSignedBits(n int) int64 {
	val := b.ReadBits(n)
	if (val & (1 << (n - 1))) != 0 {
		val |= ^uint64(0) << (n)
	}
	return int64(val)
}

func (b *BitStream) Pos() int {
	return b.pos
}
