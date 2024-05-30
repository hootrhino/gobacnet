package bacnet

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/hootrhino/gobacnet/btypes"
	"github.com/hootrhino/gobacnet/datalink"
	"github.com/hootrhino/gobacnet/encoding"
)

const interfaceName = "eth0"
const testServer = -1

// TestMain are general test
func TestUdpDataLink(t *testing.T) {
	c, _ := NewClient(&ClientBuilder{Interface: interfaceName})
	c.Close()

	_, err := datalink.NewUDPDataLink("pizzainterfacenotreal", 0)
	if err == nil {
		t.Fatal("Successfully passed a false interface.")
	}
}

func TestMac(t *testing.T) {
	var mac []byte
	json.Unmarshal([]byte("\"ChQAzLrA\""), &mac)
	l := len(mac)
	p := uint16(mac[l-1])<<8 | uint16(mac[l-1])
	log.Printf("%d", p)
}

func TestServices(t *testing.T) {
	c, _ := NewClient(&ClientBuilder{Interface: "以太网"})
	defer c.Close()
	go c.ClientRun()

	t.Run("Read Property", func(t *testing.T) {
		testReadPropertyService(c, t)
	})

	t.Run("Who Is", func(t *testing.T) {
		testWhoIs(c, t)
	})

	t.Run("WriteProperty", func(t *testing.T) {
		testWritePropertyService(c, t)
	})

}

func testReadPropertyService(c Client, t *testing.T) {
	read := btypes.PropertyData{
		Object: btypes.Object{
			ID: btypes.ObjectID{
				Type:     btypes.AnalogValue,
				Instance: 1,
			},
			Properties: []btypes.Property{
				btypes.Property{
					Type:       btypes.PropDescription, // Present value
					ArrayIndex: ArrayAll,
				},
			},
		},
	}

	mac := make([]byte, 6)
	fmt.Sscanf("192.168.0.197", "%d.%d.%d.%d", &mac[0], &mac[1], &mac[2], &mac[3])
	port := uint16(47808)
	mac[4] = byte(port >> 8)
	mac[5] = byte(port & 0x00FF)
	remoteDev := btypes.Device{
		Addr: btypes.Address{
			MacLen: 6,
			Mac:    mac,
		},
	}

	objects, err2 := c.Objects(remoteDev)
	if err2 != nil {
		t.Fatal(err2)
	}
	t.Logf("%v", objects)
	resp, err := c.ReadProperty(remoteDev, read)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Response: %v", resp.Object.Properties[0].Data)
}

func testWhoIs(c Client, t *testing.T) {
	wh := &WhoIsOpts{
		GlobalBroadcast: true,
		NetworkNumber:   0,
	}
	wh.Low = testServer - 1
	wh.High = testServer + 1
	dev, err := c.WhoIs(wh)
	if err != nil {
		t.Fatal(err)
	}
	if len(dev) == 0 {
		t.Fatalf("Unable to find device id %d", testServer)
	}
}

// This test will first cconver the name of an analogue sensor to a different
// value, read the property to make sure the name was changed, revert back, and
// ensure that the revert was successful
func testWritePropertyService(c Client, t *testing.T) {
	const targetName = "Hotdog"
	wh := &WhoIsOpts{
		GlobalBroadcast: false,
		NetworkNumber:   0,
	}
	wh.Low = testServer
	wh.High = testServer
	dev, err := c.WhoIs(wh)
	wp := btypes.PropertyData{
		Object: btypes.Object{
			ID: btypes.ObjectID{
				Type:     btypes.AnalogValue,
				Instance: 1,
			},
			Properties: []btypes.Property{
				btypes.Property{
					Type:       btypes.PropObjectName, // Present value
					ArrayIndex: ArrayAll,
					Priority:   btypes.Normal,
				},
			},
		},
	}

	if len(dev) == 0 {
		t.Fatalf("Unable to find device id %d", testServer)
	}
	resp, err := c.ReadProperty(dev[0], wp)
	if err != nil {
		t.Fatal(err)
	}
	// Store the original response since we plan to put it back in after
	org := resp.Object.Properties[0].Data
	t.Logf("original name is: %d", org)

	wp.Object.Properties[0].Data = targetName
	err = c.WriteProperty(dev[0], wp)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = c.ReadProperty(dev[0], wp)
	if err != nil {
		t.Fatal(err)
	}

	d := resp.Object.Properties[0].Data
	s, ok := d.(string)
	if !ok {
		log.Fatalf("unexpected return type %T", d)
	}

	if s != targetName {
		log.Fatalf("write to name %s did not successed, name was %s", targetName, s)
	}

	// Revert Changes
	wp.Object.Properties[0].Data = org
	err = c.WriteProperty(dev[0], wp)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = c.ReadProperty(dev[0], wp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Object.Properties[0].Data != org {
		t.Fatalf("unable to revert name back to original value %v: name is %v", org, resp.Object.Properties[0].Data)
	}
}

func TestDeviceClient(t *testing.T) {
	c, _ := NewClient(&ClientBuilder{Interface: interfaceName})
	go c.ClientRun()
	wh := &WhoIsOpts{
		GlobalBroadcast: false,
		NetworkNumber:   0,
	}
	wh.Low = testServer - 1
	wh.High = testServer - 1
	devs, err := c.WhoIs(wh)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", devs)
	//	c.Objects(devs[0])

	prop, err := c.ReadProperty(
		devs[0],
		btypes.PropertyData{
			Object: btypes.Object{
				ID: btypes.ObjectID{
					Type:     btypes.AnalogInput,
					Instance: 0,
				},
				Properties: []btypes.Property{{
					Type:       85,
					ArrayIndex: encoding.ArrayAll,
				}},
			},
			ErrorClass: 0,
			ErrorCode:  0,
		})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(prop.Object.Properties)

	props, err := c.ReadMultiProperty(devs[0], btypes.MultiplePropertyData{Objects: []btypes.Object{
		{
			ID: btypes.ObjectID{
				Type:     btypes.AnalogInput,
				Instance: 0,
			},
			Properties: []btypes.Property{
				{
					Type:       8,
					ArrayIndex: encoding.ArrayAll,
				},
				/*	{
					Type:       85,
					ArrayIndex: encoding.ArrayAll,
				},*/
			},
		},
	}})

	fmt.Println(props)
	if err != nil {
		fmt.Println(err)
		return
	}
}
