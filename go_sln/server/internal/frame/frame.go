package frame

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func ExtractFrame(buf *bytes.Buffer) ([]byte, bool) {
	b := buf.Bytes()
	start := bytes.IndexByte(b, 0x68)
	if start < 0 {
		return nil, false
	}
	if len(b) <= start+2 {
		return nil, false
	}
	if b[start+2] != 0x68 {
		buf.Next(start + 1)
		return nil, false
	}
	lenByte := int(b[start+1])
	minEnd := start + 3 + lenByte + 1 + 1
	if len(b) < minEnd {
		return nil, false
	}
	endIdx1 := start + 3 + lenByte + 1
	if endIdx1 < len(b) && b[endIdx1+1] == 0x16 {
		frame := make([]byte, endIdx1-start+2)
		copy(frame, b[start:endIdx1+2])
		buf.Next(endIdx1 + 2)
		return frame, true
	}
	endIdx2 := start + 3 + lenByte + 2
	if endIdx2 < len(b) && b[endIdx2+1] == 0x16 {
		frame := make([]byte, endIdx2-start+2)
		copy(frame, b[start:endIdx2+2])
		buf.Next(endIdx2 + 2)
		return frame, true
	}
	return nil, false
}

func PayloadData(frame []byte) []byte {
	if len(frame) <= 5 {
		return nil
	}
	lenByte := int(frame[1])
	if lenByte < 2 {
		return nil
	}
	dataLen := lenByte - 2
	dataStart := 5
	if dataStart+dataLen > len(frame)-3 {
		if dataStart >= len(frame)-3 {
			return nil
		}
		return frame[dataStart : len(frame)-3]
	}
	return frame[dataStart : dataStart+dataLen]
}

func ControlAddrFromFrame(frame []byte) (control byte, addr byte, err error) {
	if len(frame) < 6 {
		return 0, 0, fmt.Errorf("frame too short")
	}
	return frame[3], frame[4], nil
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
