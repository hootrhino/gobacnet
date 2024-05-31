package main

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/hootrhino/gobacnet"
	"github.com/hootrhino/gobacnet/apdus"
	"github.com/hootrhino/gobacnet/btypes"
)

func main() {
	// cmd.Execute()
	client, err := bacnet.NewClient(&bacnet.ClientBuilder{
		Ip:         "192.168.10.163",
		Port:       47808,
		SubnetCIDR: 24,
		DeviceId:   0,
		VendorId:   0,
		NetWorkId:  0,
		PropertyData: map[uint32][2]btypes.Object{
			0: apdus.NewAIPropertyWithRequiredFields("altitude", 0, float32(8848.46), ""),
			1: apdus.NewAIPropertyWithRequiredFields("temp", 1, float32(3.14), ""),
			2: apdus.NewAIPropertyWithRequiredFields("humi", 2, float32(77.67), ""),
			3: apdus.NewAIPropertyWithRequiredFields("pres", 3, float32(101.11), ""),
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("client run success")
	go func() {
		for {
			for i := 1; i <= 3; i++ {
				newValue := rand.Float32()
				client.GetBacnetIPServer().UpdateAIPropertyValue(uint32(i), newValue)
				fmt.Println("Update Value: ", i, ", ", newValue)
			}
			time.Sleep(3 * time.Second)
		}
	}()
	client.ClientRun()
}
