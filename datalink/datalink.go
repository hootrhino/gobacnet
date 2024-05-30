package datalink

import (
	"net"

	"github.com/hootrhino/gobacnet/btypes"
)

type DataLink interface {
	GetMyAddress() *btypes.Address
	GetListener() *net.UDPConn
	GetBroadcastAddress() *btypes.Address
	Send(data []byte, npdu *btypes.NPDU, dest *btypes.Address) (int, error)
	Receive(data []byte) (*btypes.Address, int, error)
	ReceiveFrom(data []byte) (*btypes.Address, *net.UDPAddr, int, error)
	Close() error
}
