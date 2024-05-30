package main

import "github.com/hootrhino/bacnet"

func main() {
	// cmd.Execute()
	client, err := bacnet.NewClient(&bacnet.ClientBuilder{
		Ip:         "192.168.10.163",
		Port:       47808,
		SubnetCIDR: 24,
		DeviceId:   10,
		VendorId:   10,
		NetWorkId:  10,
	})
	if err != nil {
		panic(err)
	}
	println("client run success")
	client.ClientRun()
}
