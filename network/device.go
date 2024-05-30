package network

import (
	"fmt"

	"github.com/hootrhino/bacnet"
	"github.com/hootrhino/gobacnet/btypes"
)

type Device struct {
	Ip            string
	Port          int
	DeviceID      int
	NetworkNumber int
	MacMSTP       int
	MaxApdu       uint32
	Segmentation  uint32
	StoreID       string
	DeviceName    string
	VendorName    string
	dev           btypes.Device
	network       bacnet.Client
}

// NewDevice returns a new instance of ta bacnet device
func NewDevice(net *Network, device *Device) (*Device, error) {
	var err error
	if net == nil {
		fmt.Println("network can not be nil")
		return nil, err
	}
	dev := &btypes.Device{
		Ip:            device.Ip,
		DeviceID:      device.DeviceID,
		NetworkNumber: device.NetworkNumber,
		MacMSTP:       device.MacMSTP,
		MaxApdu:       device.MaxApdu,
		Segmentation:  btypes.Enumerated(device.Segmentation),
	}
	dev, err = btypes.NewDevice(dev)
	if err != nil {
		return nil, err
	}
	if dev == nil {
		fmt.Println("dev is nil")
		return nil, err
	}
	device.network = net.bacnet
	device.dev = *dev
	if net.Store.BacStore != nil {
		net.Store.BacStore.Set(device.StoreID, device, -1)
	}
	return device, nil
}
