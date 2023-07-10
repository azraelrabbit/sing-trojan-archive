package udpnat

import (
	"context"
	"io"
	"net"
	"net/netip"
	"os"
	"time"

	"github.com/sagernet/sing/common"
	M "github.com/sagernet/sing/common/metadata"
	"github.com/sagernet/sing/tcpip/buffer"
	N "github.com/sagernet/sing/tcpip/network"
)

var (
	_ net.PacketConn = (*NATConn)(nil)
	_ N.PacketConn   = (*NATConn)(nil)
)

type NATConn struct {
	ctx         context.Context
	cancel      common.ContextCancelCauseFunc
	source      netip.AddrPort
	destination netip.AddrPort
	recvChan    chan *buffer.PacketBuffer
	writer      N.PacketWriter
}

func (c *NATConn) RecvChan() chan *buffer.PacketBuffer {
	return c.recvChan
}

func (c *NATConn) ReadPacket() (*buffer.PacketBuffer, error) {
	select {
	case packet := <-c.recvChan:
		return packet, nil
	case <-c.ctx.Done():
		return nil, io.ErrClosedPipe
	}
}

func (c *NATConn) WritePacket(packet *buffer.PacketBuffer) error {
	return c.writer.WritePacket(packet)
}

func (c *NATConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	select {
	case packet := <-c.recvChan:
		n = copy(p, packet.AsSlice())
		addr = M.SocksaddrFromNetIP(packet.Destination).UDPAddr()
		packet.Release()
		return n, addr, nil
	case <-c.ctx.Done():
		return 0, nil, io.ErrClosedPipe
	}
}

func (c *NATConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	packet := buffer.From(p)
	packet.Destination = M.AddrPortFromNet(addr)
	err = c.writer.WritePacket(packet)
	if err == nil {
		n = len(p)
	}
	return
}

func (c *NATConn) Read(p []byte) (n int, err error) {
	select {
	case packet := <-c.recvChan:
		n = copy(p, packet.AsSlice())
		packet.Release()
		return
	case <-c.ctx.Done():
		return 0, io.ErrClosedPipe
	}
}

func (c *NATConn) Write(b []byte) (n int, err error) {
	packet := buffer.From(b)
	packet.Destination = c.destination
	err = c.writer.WritePacket(packet)
	if err == nil {
		n = len(b)
	}
	return
}

func (c *NATConn) Close() error {
	c.cancel(io.ErrClosedPipe)
	return nil
}

func (c *NATConn) LocalAddr() net.Addr {
	return M.SocksaddrFromNetIP(c.destination).UDPAddr()
}

func (c *NATConn) RemoteAddr() net.Addr {
	return M.SocksaddrFromNetIP(c.source).UDPAddr()
}

func (c *NATConn) SetDeadline(t time.Time) error {
	return os.ErrInvalid
}

func (c *NATConn) SetReadDeadline(t time.Time) error {
	return os.ErrInvalid
}

func (c *NATConn) SetWriteDeadline(t time.Time) error {
	return os.ErrInvalid
}

func (c *NATConn) Upstream() any {
	return c.writer
}
