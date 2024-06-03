package gobacnet

import "testing"

func Test_ReadMultiplePropertyAck(t *testing.T) {

	client, err := NewClient(&ClientBuilder{
		Ip:         "192.168.10.163",
		Port:       47808,
		SubnetCIDR: 24,
		DeviceId:   2580,
		VendorId:   2580,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("client run success")
	client.ClientRun()
}
