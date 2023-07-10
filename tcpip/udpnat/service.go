package udpnat

import (
	"context"
	"net/netip"
	"time"

	"github.com/sagernet/sing/common"
	"github.com/sagernet/sing/common/cache"
	"github.com/sagernet/sing/common/log"
	"github.com/sagernet/sing/tcpip/buffer"
	N "github.com/sagernet/sing/tcpip/network"
)

type Service struct {
	table    *cache.LruCache[netip.AddrPort, *NATConn]
	capacity int
	handler  Handler
	logger   log.Logger
}

type Options struct {
	Timeout  time.Duration
	Capacity int
	Handler  Handler
	Logger   log.Logger
}

type Handler = func(ctx context.Context, conn *NATConn)

func New(options Options) *Service {
	return &Service{
		table: cache.New(
			cache.WithAge[netip.AddrPort, *NATConn](int64(options.Timeout.Seconds())),
			cache.WithUpdateAgeOnGet[netip.AddrPort, *NATConn](),
			cache.WithEvict[netip.AddrPort, *NATConn](func(key netip.AddrPort, conn *NATConn) {
				conn.Close()
			}),
		),
		capacity: options.Capacity,
		handler:  options.Handler,
		logger:   options.Logger,
	}
}

func (s *Service) HandlePacket(source netip.AddrPort, destination netip.AddrPort, packet *buffer.PacketBuffer, constructor func(natConn *NATConn) (context.Context, N.PacketWriter)) {
	conn, loaded := s.table.LoadOrStore(source, createNatConn)
	if !loaded {
		conn.ctx, conn.writer = constructor(conn)
		conn.ctx, conn.cancel = common.ContextWithCancelCause(conn.ctx)
		*conn = NATConn{
			source:      source,
			destination: destination,
			recvChan:    make(chan *buffer.PacketBuffer, s.capacity),
		}
	} else if common.Done(conn.ctx) {
		s.table.Delete(source)
		s.HandlePacket(source, destination, packet, constructor)
		return
	}
	select {
	case conn.recvChan <- packet:
	default:
		packet.Release()
		if s.logger != nil {
			s.logger.Warn(conn.ctx, "dropped packet from ", source, " due to full receive channel")
		}
	}
}

func createNatConn() *NATConn {
	return &NATConn{}
}
