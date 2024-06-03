package main

import (
	"math/rand/v2"
	"os"
	"time"

	bacnet "github.com/hootrhino/gobacnet"
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
		DeviceId:   163,
		VendorId:   163,
		NetWorkId:  1000, // 1-65,534
		PropertyData: map[uint32][2]btypes.Object{
			1:  apdus.NewAIPropertyWithRequiredFields("property:1", 1, float32(3.14), "-/-"),
			2:  apdus.NewAIPropertyWithRequiredFields("property:2", 2, float32(77.67), "-/-"),
			3:  apdus.NewAIPropertyWithRequiredFields("property:3", 3, float32(101.11), "-/-"),
			4:  apdus.NewAIPropertyWithRequiredFields("property:4", 4, float32(0), "-/-"),
			5:  apdus.NewAIPropertyWithRequiredFields("property:5", 5, float32(0), "-/-"),
			6:  apdus.NewAIPropertyWithRequiredFields("property:6", 6, float32(0), "-/-"),
			7:  apdus.NewAIPropertyWithRequiredFields("property:7", 7, float32(0), "-/-"),
			8:  apdus.NewAIPropertyWithRequiredFields("property:8", 8, float32(0), "-/-"),
			9:  apdus.NewAIPropertyWithRequiredFields("property:9", 9, float32(0), "-/-"),
			10: apdus.NewAIPropertyWithRequiredFields("property:10", 10, float32(0), "-/-"),
			11: apdus.NewAIPropertyWithRequiredFields("property:11", 11, float32(0), "-/-"),
			12: apdus.NewAIPropertyWithRequiredFields("property:12", 12, float32(0), "-/-"),
			13: apdus.NewAIPropertyWithRequiredFields("property:13", 13, float32(0), "-/-"),
			14: apdus.NewAIPropertyWithRequiredFields("property:14", 14, float32(0), "-/-"),
			15: apdus.NewAIPropertyWithRequiredFields("property:15", 15, float32(0), "-/-"),
			16: apdus.NewAIPropertyWithRequiredFields("property:16", 16, float32(0), "-/-"),
			17: apdus.NewAIPropertyWithRequiredFields("property:17", 17, float32(0), "-/-"),
			18: apdus.NewAIPropertyWithRequiredFields("property:18", 18, float32(0), "-/-"),
			19: apdus.NewAIPropertyWithRequiredFields("property:19", 19, float32(0), "-/-"),
			20: apdus.NewAIPropertyWithRequiredFields("property:20", 20, float32(0), "-/-"),
			21: apdus.NewAIPropertyWithRequiredFields("property:21", 21, float32(0), "-/-"),
		},
	})
	if err != nil {
		panic(err)
	}
	client.SetLogger(logger)
	logger.Debug("Bacnet client run success")
	go func() {
		for {
			for i := 1; i <= 20; i++ {
				newValue := rand.Float32() * 10
				client.GetBacnetIPServer().UpdateAIPropertyValue(uint32(i), newValue)
				logger.Debug("Update Value: ", i, ", ", newValue)
			}
			time.Sleep(3 * time.Second)
		}
	}()
	client.ClientRun()
}
