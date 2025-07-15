// All credits to Lekuruu for the original code, I added on some functions for Tencho specifically though
// https://github.com/Lekuruu/ubisoft-game-service/blob/main/common/logging.go
package common

import (
	"encoding/binary"
	"math"
)

type IOStream struct {
	data     []byte
	position int
	endian   binary.ByteOrder
}

func NewIOStream(data []byte, endian binary.ByteOrder) *IOStream {
	return &IOStream{
		data:     data,
		position: 0,
		endian:   endian,
	}
}

func (stream *IOStream) Push(data []byte) {
	stream.data = append(stream.data, data...)
}

func (stream *IOStream) Get() []byte {
	return stream.data
}

func (stream *IOStream) Len() int {
	return len(stream.data)
}

func (stream *IOStream) Available() int {
	return stream.Len() - stream.position
}

func (stream *IOStream) Tell() int {
	return stream.position
}

func (stream *IOStream) Seek(position int) {
	stream.position = position
}

func (stream *IOStream) Skip(offset int) {
	stream.position += offset
}

func (stream *IOStream) Eof() bool {
	return stream.position >= stream.Len()
}

func (stream *IOStream) Read(size int) []byte {
	if stream.Eof() {
		return []byte{}
	}

	if stream.Available() < size {
		size = stream.Available()
	}

	data := stream.data[stream.position : stream.position+size]
	stream.position += size

	return data
}

func (stream *IOStream) ReadAll() []byte {
	return stream.Read(stream.Available())
}

func (stream *IOStream) ReadU8() uint8 {
	return stream.Read(1)[0]
}

func (stream *IOStream) ReadU16() uint16 {
	return stream.endian.Uint16(stream.Read(2))
}

func (stream *IOStream) ReadU32() uint32 {
	return stream.endian.Uint32(stream.Read(4))
}

func (stream *IOStream) ReadU64() uint64 {
	return stream.endian.Uint64(stream.Read(8))
}

func (stream *IOStream) ReadULEB128() int {
	var value int
	var shift uint

	for {
		b := stream.ReadU8()
		value |= int(b&0x7F) << shift
		if b&0x80 == 0 {
			break
		}
		shift += 7
	}

	return value
}

func (stream *IOStream) ReadString() string {
	if stream.ReadU8() != 0x0B {
		return ""
	}

	length := stream.ReadULEB128()
	return string(stream.Read(length))
}

func (stream *IOStream) ReadF32() float32 {
	bits := stream.ReadU32()
	return math.Float32frombits(bits)
}

func (stream *IOStream) ReadF64() float64 {
	bits := stream.ReadU64()
	return math.Float64frombits(bits)
}

func (stream *IOStream) Write(data []byte) {
	stream.Push(data)
}

func (stream *IOStream) WriteU8(value uint8) {
	stream.Write([]byte{value})
}

func (stream *IOStream) WriteU16(value uint16) {
	data := make([]byte, 2)
	stream.endian.PutUint16(data, value)
	stream.Write(data)
}

func (stream *IOStream) WriteU32(value uint32) {
	data := make([]byte, 4)
	stream.endian.PutUint32(data, value)
	stream.Write(data)
}

func (stream *IOStream) WriteU64(value uint64) {
	data := make([]byte, 8)
	stream.endian.PutUint64(data, value)
	stream.Write(data)
}

func (stream *IOStream) WriteF32(value float32) {
	bits := math.Float32bits(value)
	stream.WriteU32(bits)
}

func (stream *IOStream) WriteF64(value float64) {
	bits := math.Float64bits(value)
	stream.WriteU64(bits)
}

func (stream *IOStream) WriteULEB128(value int) {
	for {
		b := uint8(value & 0x7F)
		value >>= 7
		if value == 0 {
			stream.WriteU8(b)
			break
		}
		stream.WriteU8(b | 0x80)
	}
}

func (stream *IOStream) WriteString(value string) {
	if value == "" {
		stream.WriteU8(0)
		return
	}

	stream.WriteU8(0x0B) // String type identifier
	stream.WriteULEB128(len(value))
	stream.Write([]byte(value))
}
