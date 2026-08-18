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
	if b.pos%8 != 0 {
		return b.readBitsArbitrary(n)
	}
	switch n {
	case 16:
		return b.read16BitsAligned()
	default:
		// fallback to unoptimized bit-by-bit reading, but a better way would be to implement more aligned reads
		return b.readBitsArbitrary(n)
	}
}

// read16BitsAligned reads 16 bits from the stream in big-endian order when the position is byte-aligned.
func (b *BitStream) read16BitsAligned() uint64 {
	byteIdx := b.pos / 8
	val := uint64(b.data[byteIdx])<<8 | uint64(b.data[byteIdx+1])
	b.pos += 16
	return val
}

func (b *BitStream) readBitsArbitrary(n int) uint64 {
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

// SetPos is used in tests to set the current bit position in the stream.
func (b *BitStream) SetPos(pos int) {
	b.pos = pos
}
