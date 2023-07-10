package network

import (
	"github.com/sagernet/sing/tcpip/buffer"
)

type PacketConn interface {
	AbstractConn
	PacketReader
	PacketWriter
}

type PacketReader interface {
	ReadPacket() (*buffer.PacketBuffer, error)
}

type PacketWriter interface {
	WritePacket(packet *buffer.PacketBuffer) error
}
