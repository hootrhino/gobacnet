package main

import (
	"math/rand/v2"
	"os"
	"time"

	"github.com/hootrhino/gobacnet"
	"github.com/hootrhino/gobacnet/apdus"
	"github.com/hootrhino/gobacnet/btypes"
	"github.com/sirupsen/logrus"
)

func main() {
	if len(os.Args) < 2 {
		panic("Missing Ip")
	}
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	client, err := bacnet.NewClient(&bacnet.ClientBuilder{
		Ip:         os.Args[1],
		Port:       47808,
		SubnetCIDR: 24,
		DeviceId:   0x00_01_02_03,
		VendorId:   0x00_01_02_03,
		NetWorkId:  1000, // 1-65,534
		PropertyData: map[uint32][2]btypes.Object{
			1:  apdus.NewAIPropertyWithRequiredFields("temp", 1, float32(3.14), "-/-"),
			2:  apdus.NewAIPropertyWithRequiredFields("humi", 2, float32(77.67), "-/-"),
			3:  apdus.NewAIPropertyWithRequiredFields("pres", 3, float32(101.11), "-/-"),
			4:  apdus.NewAIPropertyWithRequiredFields("altitude", 4, float32(0), "-/-"),
			5:  apdus.NewAIPropertyWithRequiredFields("altitude", 5, float32(0), "-/-"),
			6:  apdus.NewAIPropertyWithRequiredFields("altitude", 6, float32(0), "-/-"),
			7:  apdus.NewAIPropertyWithRequiredFields("altitude", 7, float32(0), "-/-"),
			8:  apdus.NewAIPropertyWithRequiredFields("altitude", 8, float32(0), "-/-"),
			9:  apdus.NewAIPropertyWithRequiredFields("altitude", 9, float32(0), "-/-"),
			10: apdus.NewAIPropertyWithRequiredFields("altitude", 10, float32(0), "-/-"),
		},
	})
	client.SetLogger(logger)
	if err != nil {
		panic(err)
	}
	logger.Debug("client run success")
	go func() {
		for {
			for i := 1; i <= 10; i++ {
				newValue := rand.Float32()
				client.GetBacnetIPServer().UpdateAIPropertyValue(uint32(i), newValue)
				logger.Debug("Update Value: ", i, ", ", newValue)
			}
			time.Sleep(3 * time.Second)
		}
	}()
	client.ClientRun()
}
