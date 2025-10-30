package data_representation

// BitStream is used for reading bits from a byte slice and decoding them into values.
type BitStream struct {
	data []byte
	pos  int // bit position
}

// NewBitStream instantiates a new BitStream for the given byte slice.
func NewBitStream(data []byte) *BitStream {
	return &BitStream{data: data, pos: 0}
}

// ReadBits reads n bits from the stream and returns them as a uint64.
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

// ReadSignedBits reads n bits from the stream and returns them as an int64.
func (b *BitStream) ReadSignedBits(n int) int64 {
	val := b.ReadBits(n)
	if (val & (1 << (n - 1))) != 0 {
		val |= ^uint64(0) << (n)
	}
	return int64(val)
}

// Pos returns the current bit position in the stream.
func (b *BitStream) Pos() int {
	return b.pos
}
