package pdv

const (
	headerMagic  = "Playdate VID"
	headerLength = 28
)

type Header struct {
	Magic     [12]byte
	Reserved1 uint32
	NumFrames uint16
	Reserved2 uint16
	FrameRate float32
	Width     uint16
	Height    uint16
}

// alternate header with uint32 NumFrames
// type Header struct {
// 	Magic     [12]byte
// 	Reserved1 uint32
// 	NumFrames uint32
// 	FrameRate float32
// 	Width     uint16
// 	Height    uint16
// }

// NewHeader builds an initial header structure, omitting the number of frames.
func NewHeader(frameRate float32, width, height uint16) *Header {
	var magicBytes [12]byte
	copy(magicBytes[:], headerMagic)
	return &Header{
		Magic:     magicBytes,
		FrameRate: frameRate,
		Width:     width,
		Height:    height,
	}
}

// NewHeaderWithNumFrames builds a header structure.
func NewHeaderWithNumFrames(frameRate float32, width, height, numFrames uint16) *Header {
	var magicBytes [12]byte
	copy(magicBytes[:], headerMagic)
	return &Header{
		Magic:     magicBytes,
		NumFrames: numFrames,
		FrameRate: frameRate,
		Width:     width,
		Height:    height,
	}
}

func (h Header) ValidateMagic() bool {
	return string(h.Magic[:]) == headerMagic
}

type FrameType int

const (
	FrameTypeEmpty FrameType = iota
	FrameTypeIFrame
	FrameTypePFrame
	FrameTypeIPFrame
)

type FrameTableEntry uint32

func NewFrameTableEntry(offset uint32, typ FrameType) FrameTableEntry {
	return FrameTableEntry(offset<<2 | uint32(typ)&0x3)
}

func (e FrameTableEntry) Offset() uint32 {
	return uint32(e >> 2)
}

func (e FrameTableEntry) Type() FrameType {
	return FrameType(e & 0x3)
}

func ExpandBits(b byte) [8]byte {
	var bs [8]byte
	for i := 7; i >= 0; i-- {
		bit := b >> i & 0b00000001
		bs[7-i] = bit // when using a paletted image
		// bs[7-i] = bit * 0xFF // when using an rgb image
	}
	return bs
}
