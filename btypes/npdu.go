package btypes

import "github.com/hootrhino/bacnet/btypes/ndpu"

type NPDUPriority byte

const ProtocolVersion uint8 = 1
const DefaultHopCount uint8 = 255

const (
	LifeSafety        NPDUPriority = 3
	CriticalEquipment NPDUPriority = 2
	Urgent            NPDUPriority = 1
	Normal            NPDUPriority = 0
)

type NPDU struct {
	Version uint8

	// Destination (optional)
	Destination *Address

	// Source (optional)
	Source *Address

	VendorId uint16

	IsNetworkLayerMessage   bool
	NetworkLayerMessageType ndpu.NetworkMessageType
	ExpectingReply          bool
	Priority                NPDUPriority
	HopCount                uint8
}

func NewNPDU() NPDU {

	return NPDU{
		Version:                 1,
		Destination:             &Address{},
		Source:                  &Address{},
		VendorId:                0,
		IsNetworkLayerMessage:   false,
		NetworkLayerMessageType: ndpu.NetworkMessageType(None),
	}
}
