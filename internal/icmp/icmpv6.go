package icmp

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	ICMPv6TypeDestinationUnreachable                       ICMPv6Type = 1
	ICMPv6TypePacketTooBig                                 ICMPv6Type = 2
	ICMPv6TypeTimeExceeded                                 ICMPv6Type = 3
	ICMPv6TypeParameterProblem                             ICMPv6Type = 4
	ICMPv6TypeEchoRequest                                  ICMPv6Type = 128
	ICMPv6TypeEchoReply                                    ICMPv6Type = 129
	ICMPv6MulticastListenerQuery                           ICMPv6Type = 130
	ICMPv6MulticastListenerReport                          ICMPv6Type = 131
	ICMPv6MulticastListenerDone                            ICMPv6Type = 132
	ICMPv6TypeRouterSolicitation                           ICMPv6Type = 133
	ICMPv6TypeRouterAdvertisement                          ICMPv6Type = 134
	ICMPv6TypeNeighborSolicitation                         ICMPv6Type = 135
	ICMPv6TypeNeighborAdvertisement                        ICMPv6Type = 136
	ICMPv6TypeRedirectMessage                              ICMPv6Type = 137
	ICMPv6TypeRouterRenumbering                            ICMPv6Type = 138
	ICMPv6TypeICMPNodeInformationQuery                     ICMPv6Type = 139
	ICMPv6TypeICMPNodeInformationResponse                  ICMPv6Type = 140
	ICMPv6TypeInverseNeighborDiscoverySolicitationMessage  ICMPv6Type = 141
	ICMPv6TypeInverseNeighborDiscoveryAdvertisementMessage ICMPv6Type = 142
	ICMPv6TypeMulticastListenerDiscoveryReports            ICMPv6Type = 143
	ICMPv6TypeHomeAgentAddressDiscoveryRequestMessage      ICMPv6Type = 144
	ICMPv6TypeHomeAgentAddressDiscoveryReplyMessage        ICMPv6Type = 145
	ICMPv6TypeMobilePrefixSolicitation                     ICMPv6Type = 146
	ICMPv6TypeMobilePrefixAdvertisement                    ICMPv6Type = 147
	ICMPv6TypeCertificationPathSolicitation                ICMPv6Type = 148
	ICMPv6TypeCertificationPathAdvertisement               ICMPv6Type = 149
	ICMPv6TypeMulticastRouterAdvertisement                 ICMPv6Type = 151
	ICMPv6TypeMulticastRouterSolicitation                  ICMPv6Type = 152
	ICMPv6TypeMulticastRouterTermination                   ICMPv6Type = 153
	ICMPv6TypeRPLControlMessage                            ICMPv6Type = 155
)

type ICMPv6Type uint8

type ICMPv6Packet struct {
	messageType ICMPv6Type
	messageCode uint8
	checksum    uint16
	data        []byte

	raw []byte
}

func (p *ICMPv6Packet) Raw() []byte {
	return p.raw
}

func NewICMPv6Packet(messageType ICMPv6Type, messageCode uint8, data []byte) *ICMPv6Packet {

	// Build ICMP header
	raw := []byte{uint8(messageType), messageCode, 0, 0}

	// Add ICMP data
	raw = append(raw, data...)

	raw, checksum, err := InternetChecksum(raw, 2)
	if err != nil {
		panic(err)
	}

	return &ICMPv6Packet{messageType: messageType, messageCode: messageCode, checksum: binary.LittleEndian.Uint16(checksum), data: data, raw: raw}

}

func NewICMPv6EchoRequestPacket(identifier uint16, sequenceNumber uint16, messageLength uint16, payloadContent []byte) (packet *ICMPv6Packet, err error) {
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

	packet = NewICMPv6Packet(ICMPv6TypeEchoRequest, 0, payload)

	return packet, nil
}
