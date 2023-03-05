package icmp

import (
	"encoding/binary"
	"fmt"
)

func InternetChecksum(data []byte, checksumFieldByteOffset uint) (packet []byte, checksum []byte, err error) {

	if len(data) < 2 || len(data) < int(checksumFieldByteOffset+2) {
		return nil, nil, fmt.Errorf("argument error")
	}

	// Clear Checksum Field
	data[checksumFieldByteOffset] = 0
	data[checksumFieldByteOffset+1] = 0

	var sum uint32 = 0
	count := len(data)

	// Calculate sum
	for i := 0; i < count; i += 2 {
		sum += uint32(binary.LittleEndian.Uint16(data[i : i+2]))
	}

	// Add left-over byte
	if count%2 != 0 {
		sum += uint32(data[count-1])
	}

	// Fold 32-bit sum to 16 bits
	for sum>>16 != 0 {
		sum = (sum & 0xffff) + (sum >> 16)
	}

	checksum = make([]byte, 2)

	binary.LittleEndian.PutUint16(checksum, ^uint16(sum))

	data[checksumFieldByteOffset] = checksum[0]
	data[checksumFieldByteOffset+1] = checksum[1]

	packet = data

	return data, checksum, nil
}
