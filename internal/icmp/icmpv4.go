package icmp

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	ICMPv4TypeEchoReply              ICMPv4Type = 0
	ICMPv4TypeDestinationUnreachable ICMPv4Type = 3
	ICMPv4TypeRedirectMessage        ICMPv4Type = 5
	ICMPv4TypeEchoRequest            ICMPv4Type = 8
	ICMPv4TypeRouterAdvertisement    ICMPv4Type = 9
	ICMPv4TypeRouterSolicitation     ICMPv4Type = 10
	ICMPv4TypeTimeExceeded           ICMPv4Type = 11
	ICMPv4TypeParameterProblem       ICMPv4Type = 12
	ICMPv4TypeTimestamp              ICMPv4Type = 13
	ICMPv4TypeTimestampReply         ICMPv4Type = 14
	ICMPv4TypeInformationRequest     ICMPv4Type = 15
	ICMPv4TypeInformationReply       ICMPv4Type = 16
	ICMPv4TypeAddressMaskRequest     ICMPv4Type = 17
	ICMPv4TypeAddressMaskReply       ICMPv4Type = 18
	ICMPv4TypeExtendedEchoRequest    ICMPv4Type = 42
	ICMPv4TypeExtendedEchoReply      ICMPv4Type = 43
)

type ICMPv4Type uint8

type ICMPv4Packet struct {
	messageType ICMPv4Type
	messageCode uint8
	checksum    uint16
	data        []byte

	raw []byte
}

func (p *ICMPv4Packet) Raw() []byte {
	return p.raw
}

func NewICMPv4Packet(messageType ICMPv4Type, messageCode uint8, data []byte) *ICMPv4Packet {

	// Build ICMP header
	raw := []byte{uint8(messageType), messageCode, 0, 0}

	// Add ICMP data
	raw = append(raw, data...)

	raw, checksum, err := InternetChecksum(raw, 2)
	if err != nil {
		panic(err)
	}

	return &ICMPv4Packet{messageType: messageType, messageCode: messageCode, checksum: binary.LittleEndian.Uint16(checksum), data: data, raw: raw}

}

func NewICMPv4EchoRequestPacket(identifier uint16, sequenceNumber uint16, messageLength uint16, payloadContent []byte) (packet *ICMPv4Packet, err error) {
	if messageLength < 8 {
		return nil, fmt.Errorf("minimum message length is 8 bytes. %d bytes were requested", messageLength)
	}
	payload := make([]byte, messageLength-4)

	if messageLength > 8 {
		// For every byte thats unset copy the payload content
		copy(payload[4:], bytes.Repeat(payloadContent, int((messageLength-8)/uint16(len(payloadContent)))+1))
	}

	payload[0] = uint8(identifier >> 8)
	payload[1] = uint8(identifier & 0xff)
	payload[2] = uint8(sequenceNumber >> 8)
	payload[3] = uint8(sequenceNumber & 0xff)

	packet = NewICMPv4Packet(ICMPv4TypeEchoRequest, 0, payload)

	return packet, nil
}
