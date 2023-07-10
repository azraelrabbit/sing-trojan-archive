package udpcontrol

import (
	"fmt"
	"net/netip"
	"unsafe"

	"golang.org/x/sys/windows"
)

// https://github.com/database64128/shadowsocks-go/blob/06f4c8214a0e6f80ae0319d274ade38e8ed148dd/conn/conn_windows.go

// Structure CMSGHDR from ws2def.h
type Cmsghdr struct {
	Len   uint
	Level int32
	Type  int32
}

// Structure IN_PKTINFO from ws2ipdef.h
type Inet4Pktinfo struct {
	Addr    [4]byte
	Ifindex uint32
}

// Structure IN6_PKTINFO from ws2ipdef.h
type Inet6Pktinfo struct {
	Addr    [16]byte
	Ifindex uint32
}

const (
	SizeofCmsghdr      = unsafe.Sizeof(Cmsghdr{})
	SizeofInet4Pktinfo = unsafe.Sizeof(Inet4Pktinfo{})
	SizeofInet6Pktinfo = unsafe.Sizeof(Inet6Pktinfo{})
)

const SizeofPtr = unsafe.Sizeof(uintptr(0))

// SocketControlMessageBufferSize specifies the buffer size for receiving socket control messages.
const SocketControlMessageBufferSize = SizeofCmsghdr + (SizeofInet6Pktinfo+SizeofPtr-1) & ^(SizeofPtr-1)

// ParsePktinfoCmsg parses a single socket control message of type IP_PKTINFO or IPV6_PKTINFO,
// and returns the IP address and index of the network interface the packet was received from,
// or an error.
//
// This function is only implemented for Linux, macOS and Windows. On other platforms, this is a no-op.
func ParsePktinfoCmsg(cmsg []byte) (netip.Addr, uint32, error) {
	if len(cmsg) < int(SizeofCmsghdr) {
		return netip.Addr{}, 0, fmt.Errorf("control message length %d is shorter than cmsghdr length", len(cmsg))
	}

	cmsghdr := (*Cmsghdr)(unsafe.Pointer(&cmsg[0]))

	switch {
	case cmsghdr.Level == windows.IPPROTO_IP && cmsghdr.Type == windows.IP_PKTINFO && len(cmsg) >= int(SizeofCmsghdr+SizeofInet4Pktinfo):
		pktinfo := (*Inet4Pktinfo)(unsafe.Pointer(&cmsg[SizeofCmsghdr]))
		return netip.AddrFrom4(pktinfo.Addr), pktinfo.Ifindex, nil

	case cmsghdr.Level == windows.IPPROTO_IPV6 && cmsghdr.Type == windows.IPV6_PKTINFO && len(cmsg) >= int(SizeofCmsghdr+SizeofInet6Pktinfo):
		pktinfo := (*Inet6Pktinfo)(unsafe.Pointer(&cmsg[SizeofCmsghdr]))
		return netip.AddrFrom16(pktinfo.Addr), pktinfo.Ifindex, nil

	default:
		return netip.Addr{}, 0, fmt.Errorf("unknown control message level %d type %d", cmsghdr.Level, cmsghdr.Type)
	}
}
