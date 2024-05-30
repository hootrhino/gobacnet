package main

import (
	"github.com/hootrhino/bacnet"
	"github.com/hootrhino/bacnet/apdus"
	"github.com/hootrhino/bacnet/btypes"
)

func main() {
	// cmd.Execute()
	client, err := bacnet.NewClient(&bacnet.ClientBuilder{
		Ip:         "192.168.10.163",
		Port:       47808,
		SubnetCIDR: 24,
		DeviceId:   10,
		VendorId:   10,
		NetWorkId:  10,
		MultiplePropertyData: map[uint32][2]btypes.Object{
			1: apdus.NewAIPropertyWithRequiredFields("temp", 1, float32(3.14), "temp des"),
			2: apdus.NewAIPropertyWithRequiredFields("humi", 2, float32(77.67), "humi des"),
			3: apdus.NewAIPropertyWithRequiredFields("pres", 3, float32(101.11), "pres des"),
		},
	})
	if err != nil {
		panic(err)
	}
	println("client run success")
	client.ClientRun()
}
