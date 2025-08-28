package emulator

import (
	"encoding/binary"
	"time"

	"sln/internal/frame"
)

func BuildTimeResponse(reqCtrl byte, reqAddr byte, reqData []byte, crcMode string, adapterAddr byte) []byte {
	respCtrl := reqCtrl | 0x80
	respAddr := reqAddr

	timeStr := time.Now().Format("2006-01-02 15:04:05")
	payload := append([]byte{0x01}, []byte(timeStr)...)

	skel := frame.BuildSkeleton(respCtrl, respAddr, payload)
	full := frame.AppendChecksum(skel, crcMode)
	return full
}

func BuildAckResponse(reqCtrl byte, reqAddr byte, reqData []byte, crcMode string, adapterAddr byte) []byte {
	respCtrl := reqCtrl | 0x80
	respAddr := reqAddr
	cmd := byte(0xFF)
	if len(reqData) > 0 {
		cmd = reqData[0]
	}
	payload := append([]byte{cmd}, []byte("OK")...)
	skel := frame.BuildSkeleton(respCtrl, respAddr, payload)
	full := frame.AppendChecksum(skel, crcMode)
	return full
}

func putUint16LE(v uint16) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, v)
	return b
}
