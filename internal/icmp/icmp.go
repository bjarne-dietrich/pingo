package icmp

type ICMPPacket interface {
	Raw() []byte
}
