package frame

import (
	"bytes"
	"encoding/binary"
)

func ExtractFrame(buf *bytes.Buffer) ([]byte, bool) {
	b := buf.Bytes()

	start := bytes.IndexByte(b, 0x68)
	if start < 0 {
		return nil, false
	}

	if len(b) < start+3 {
		return nil, false
	}

	if b[start+2] != 0x68 {
		buf.Next(start + 1)
		return nil, false
	}

	lenByte := int(b[start+1])
	if lenByte < 0 {
		buf.Next(start + 1)
		return nil, false
	}

	payloadStart := start + 3
	payloadEnd := payloadStart + lenByte

	if len(b) < payloadEnd+2 {
		return nil, false
	}

	endIdx1 := payloadEnd + 1
	if endIdx1 < len(b) && b[endIdx1] == 0x16 {
		frame := make([]byte, endIdx1-start+1)
		copy(frame, b[start:endIdx1+1])
		buf.Next(endIdx1 + 1)
		return frame, true
	}

	endIdx2 := payloadEnd + 2
	if endIdx2 < len(b) && b[endIdx2] == 0x16 {
		frame := make([]byte, endIdx2-start+1)
		copy(frame, b[start:endIdx2+1])
		buf.Next(endIdx2 + 1)
		return frame, true
	}

	return nil, false
}

func PayloadData(frame []byte) []byte {
	if len(frame) < 7 {
		return nil
	}
	lenByte := int(frame[1])
	if lenByte < 2 {
		return nil
	}
	payloadStart := 3
	payloadEnd := payloadStart + lenByte

	if payloadEnd > len(frame)-2 {
		return nil
	}

	dataStart := payloadStart + 2
	if dataStart > payloadEnd {
		return nil
	}
	return frame[dataStart:payloadEnd]
}

func BuildSkeleton(control byte, addr byte, data []byte) []byte {
	lenByte := byte(2 + len(data))
	var b bytes.Buffer
	b.WriteByte(0x68)
	b.WriteByte(lenByte)
	b.WriteByte(0x68)
	b.WriteByte(control)
	b.WriteByte(addr)
	b.Write(data)
	return b.Bytes()
}

func AppendChecksum(frameSoFar []byte, crcMode string) []byte {
	if crcMode == "crc16" {
		crc := ComputeCRC16(frameSoFar[3:])
		tmp := make([]byte, 2)
		binary.LittleEndian.PutUint16(tmp, crc)
		return append(frameSoFar, append(tmp, 0x16)...)
	}
	sum := ComputeSum(frameSoFar[3:])
	return append(frameSoFar, append([]byte{sum}, 0x16)...)
}
