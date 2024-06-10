package gobacnet

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/hootrhino/gobacnet/apdus"
	"github.com/hootrhino/gobacnet/btypes"
	"github.com/hootrhino/gobacnet/btypes/ndpu"
	"github.com/hootrhino/gobacnet/datalink"
	"github.com/hootrhino/gobacnet/encoding"
	"github.com/hootrhino/gobacnet/helpers/validation"
	"github.com/hootrhino/gobacnet/tsm"
	"github.com/hootrhino/gobacnet/utsm"
	"github.com/sirupsen/logrus"
)

const mtuHeaderLength = 4
const defaultStateSize = 20
const forwardHeaderLength = 10

type client struct {
	dataLink       datalink.DataLink
	tsm            *tsm.TSM
	utsm           *utsm.Manager
	readBufferPool sync.Pool
	log            *logrus.Logger
	deviceId       uint32
	vendorId       uint32
	netWorkId      uint16
	server         *BacnetIPServer
}

type ClientBuilder struct {
	DataLink     datalink.DataLink
	Interface    string
	Ip           string
	Port         int
	SubnetCIDR   int
	MaxPDU       uint16
	NetWorkId    uint16
	DeviceId     uint32
	VendorId     uint32
	PropertyData map[uint32][2]btypes.Object
}

// NewClient creates a new client with the given interface and
func NewClient(cb *ClientBuilder) (Client, error) {
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	logger.SetLevel(logrus.InfoLevel)

	var err error
	var dataLink datalink.DataLink
	iface := cb.Interface
	ip := cb.Ip
	port := cb.Port
	maxPDU := cb.MaxPDU
	//check port
	if port == 0 {
		port = datalink.DefaultPort
	}
	ok := validation.ValidPort(port)
	if !ok {
		return nil, errors.New("invalid port")
	}
	//check adpu
	if maxPDU == 0 {
		maxPDU = btypes.MaxAPDU
	}
	//build datalink
	if iface != "" {
		dataLink, err = datalink.NewUDPDataLink(iface, port)
		if err != nil {
			return nil, err
		}
	} else {
		//check ip
		ok = validation.ValidIP(ip)
		if !ok {
			return nil, errors.New("invalid Ip")
		}
		//check subnet
		sub := cb.SubnetCIDR
		ok = validation.ValidCIDR(ip, sub)
		if !ok {
			return nil, errors.New("validate CIDR failed")
		}
		dataLink, err = datalink.NewUDPDataLinkFromIP(ip, sub, port)
		if err != nil {
			return nil, err
		}
	}
	// 1-65,534
	if cb.NetWorkId > 65534 {
		panic(fmt.Errorf("BACnet network numbers with a Range of 1-65534"))
	}
	cli := &client{
		deviceId:  cb.DeviceId,
		vendorId:  cb.VendorId,
		dataLink:  dataLink,
		netWorkId: cb.NetWorkId,
		server:    NewBacnetIPServer(cb.DeviceId, cb.VendorId, cb.NetWorkId, cb.PropertyData),
		tsm:       tsm.New(defaultStateSize),
		utsm: utsm.NewManager(
			utsm.DefaultSubscriberTimeout(time.Second*time.Duration(10)),
			utsm.DefaultSubscriberLastReceivedTimeout(time.Second*time.Duration(2)),
		),
		readBufferPool: sync.Pool{New: func() interface{} {
			return make([]byte, maxPDU)
		}},
		log: logger,
	}
	return cli, nil
}

// GetBroadcastAddress uses the given address with subnet to return the broadcast address
func (c *client) GetBroadcastAddress() *btypes.Address {
	return c.dataLink.GetBroadcastAddress()
}
func (c *client) GetMyAddress() *btypes.Address {
	return c.dataLink.GetMyAddress()
}
func (c *client) GetListener() *net.UDPConn {
	return c.dataLink.GetListener()
}
func (c *client) GetBacnetIPServer() *BacnetIPServer {
	return c.server
}

// expired
func (c *client) ClientRun() {
	c.StartPoll(context.Background())
}

/*
*
* Context
*
 */
func (c *client) StartPoll(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			{
				return
			}
		default:
			{
			}
		}
		b := c.readBufferPool.Get().([]byte)
		pduAddr, udpAddr, n, err := c.dataLink.ReceiveFrom(b)
		if err != nil {
			c.log.Error(err)
			continue
		}
		go c.handleMsg(pduAddr, udpAddr, b[:n])
	}
}

// stop
func (c *client) StopPoll() {
	c.ClientClose(true)
	c.Close()
}

func (c *client) handleMsg(src *btypes.Address, udpAddr *net.UDPAddr, b []byte) {
	var header btypes.BVLC
	var npdu btypes.NPDU
	var apdu btypes.APDU
	dec := encoding.NewDecoder(b)
	err := dec.BVLC(&header)
	if err != nil {
		c.log.Error(err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			c.log.Errorf("handleMsg recover: %v", r)
		}
	}()

	if header.Function == btypes.BacFuncBroadcast ||
		header.Function == btypes.BacFuncUnicast ||
		header.Function == btypes.BacFuncForwardedNPDU {
		// Remove the header information
		b = b[mtuHeaderLength:]
		networkList, err := dec.NPDU(&npdu)
		if err != nil {
			return
		}
		if npdu.IsNetworkLayerMessage {
			c.log.Debug("Ignored Network Layer Message")
			if npdu.NetworkLayerMessageType == ndpu.NetworkIs {
				c.utsm.Publish(int(npdu.Source.Net), npdu)
				return
			}
			if npdu.NetworkLayerMessageType == ndpu.IamRouterToNetwork {
				c.utsm.Publish(int(npdu.Source.Net), networkList)
				return
			}

		}
		// We want to keep the APDU intact, so we will get a snapshot before decoding
		send := dec.Bytes()
		err = dec.APDU(&apdu)
		if err != nil {
			c.log.Errorf("Issue decoding APDU: %v", err)
			return
		}
		switch apdu.DataType {
		case btypes.UnconfirmedServiceRequest:
			if apdu.UnconfirmedService == btypes.ServiceUnconfirmedIAm {
				dec = encoding.NewDecoder(apdu.RawData)
				var iam btypes.IAm
				err = dec.IAm(&iam)
				c.log.Debug("Received IAM Message", iam.ID)
				iam.Addr = *src

				if npdu.Source != nil {
					if npdu.Source.Net > 0 { // add in device network number
						c.log.Debug("device-network-address", npdu.Source.Net)
						iam.Addr.Net = npdu.Source.Net
					}
					if len(npdu.Source.Adr) > 0 { // add in hardware mac
						c.log.Debug("device-mstp-mac-address", npdu.Source.Adr)
						iam.Addr.Adr = npdu.Source.Adr
						iam.Addr.Len = uint8(len(iam.Addr.Adr))
					}
				}
				if err != nil {
					c.log.Error(err)
					return
				}

				c.utsm.Publish(int(iam.ID.Instance), iam)
			} else if apdu.UnconfirmedService == btypes.ServiceUnconfirmedWhoIs {
				dec := encoding.NewDecoder(apdu.RawData)
				var low, high int32
				dec.WhoIs(&low, &high)

				reply := false
				if low == -1 || high == -1 {
					reply = true
				}
				if low <= int32(c.deviceId) && high >= int32(c.deviceId) {
					reply = true
				}

				if reply {
					iamBytes, errNewIAm := apdus.NewIAm(c.deviceId, uint16(c.vendorId), c.netWorkId)
					if errNewIAm != nil {
						c.log.Errorf("New IAm failed err:%v", errNewIAm)
						return
					}
					c.log.Debug("who is from:", udpAddr.String())
					c.log.Debug("I AM to:", udpAddr.String())
					_, errWrite := c.GetListener().WriteTo(iamBytes, udpAddr)
					if errWrite != nil {
						c.log.Error("Error Write To data:", errWrite)
						return
					}
				}
			} else {
				c.log.Errorf("Unconfirmed: %d %v", apdu.UnconfirmedService, apdu.RawData)
			}
		case btypes.SimpleAck:
			c.log.Debug("Received Simple Ack")
			err := c.tsm.Send(int(apdu.InvokeId), send)
			if err != nil {
				return
			}
		case btypes.ComplexAck:
			c.log.Debug("Received Complex Ack")
			err := c.tsm.Send(int(apdu.InvokeId), send)
			if err != nil {
				return
			}
		case btypes.ConfirmedServiceRequest:
			c.log.Debug("Received Confirmed Service Request")
			if apdu.Service == btypes.ServiceConfirmedReadPropMultiple {
				Encoder := encoding.NewEncoder()
				Encoder.BVLC(btypes.BVLC{
					Type:     0x81,
					Function: 0x0A,
				})
				Encoder.NPDU(&btypes.NPDU{
					Version:     btypes.ProtocolVersion,
					Destination: src,
				})
				Decoder := encoding.NewDecoder(apdu.RawData)
				PropertyDataRequest := btypes.PropertyData{
					Object: btypes.Object{
						ID: btypes.ObjectID{
							Type:     btypes.AnalogInput,
							Instance: 0,
						},
						Properties: []btypes.Property{},
					},
				}
				errReadMultiplePropertyAck := Decoder.ReadProperty(&PropertyDataRequest)
				if errReadMultiplePropertyAck != nil {
					c.log.Error("Error sending data:", errReadMultiplePropertyAck)
					return
				}
				// ObjectInstance: 要读 ObjectInstance 代表的对象的属性表
				ObjectInstance := PropertyDataRequest.Object.ID.Instance
				// 这个地方应该分两种情况
				// 1 来自对象的必须属性（7类）
				// 2 额外属性
				// 但是当前阶段暂时支持必须属性
				ObjectProperties, errGetObjectProperties := c.server.GetObjectProperties(uint32(ObjectInstance))
				if errGetObjectProperties != nil {
					c.log.Error("Error GetObjectProperties:", errGetObjectProperties)
					return
				}
				// RequiredPropertiesResponse := apdus.NewRequiredPropertiesResponse(apdu.InvokeId, uint32(ObjectInstance))
				Encoder.PackageReadMultiplePropertyAck(apdu.InvokeId, ObjectProperties)
				_, errWrite := c.GetListener().WriteTo(Encoder.Package(), udpAddr)
				if errWrite != nil {
					c.log.Error("Error Write To data:", errWrite)
					return
				}
			}
			// 读Object Property
			if apdu.Service == btypes.ServiceConfirmedReadProperty {
				decoder := encoding.NewDecoder(apdu.RawData)
				PropertyData := btypes.PropertyData{
					Object: btypes.Object{
						Properties: []btypes.Property{},
					},
				}
				err := decoder.ReadProperty(&PropertyData)
				if err != nil {
					c.log.Errorf("decoder ReadProperty failed; %d %v err=%v", apdu.Service, apdu.RawData, err)
					return
				}
				// 返回列表: ArrayIndex 就是对象ID
				ArrayIndex := PropertyData.Object.Properties[0].ArrayIndex
				if ArrayIndex == 0 {
					// 设备id + service=readProperty + object-list(76)
					// ARRAY index = 0, 返回个数
					// array index = 1, 返回第一个object
					// array index = 2, 返回第二个object
					PropertyResponseBytes, _ := apdus.NewReadPropertyListResponse(apdu.InvokeId,
						c.deviceId, uint8(len((c.server.PropertyData))))
					_, errWrite := c.GetListener().WriteTo(PropertyResponseBytes, udpAddr)
					if errWrite != nil {
						c.log.Error("Error sending data:", errWrite)
						return
					}
				} else {
					// 遍历属性 for  1 2 3 4... N
					// ArrayIndex 是外面传进来的遍历索引，有几个 这个ArrayIndex就递增几次
					// 让 ArrayIndex 在点位表里面检索
					// if is required ..
					PropertyResponseBytes, _ := apdus.NewReadPropertyResponse(apdu.InvokeId, c.deviceId, ArrayIndex)
					_, errWrite := c.GetListener().WriteTo(PropertyResponseBytes, udpAddr)
					if errWrite != nil {
						c.log.Error("Error sending data:", errWrite)
						return
					}
				}
			} else {
				c.log.Errorf("Confimed: %d %v", apdu.Service, apdu.RawData)
			}
		case btypes.Error:
			err := fmt.Errorf("error class %s code %s", apdu.Error.Class.String(), apdu.Error.Code.String())
			err = c.tsm.Send(int(apdu.InvokeId), err)
			if err != nil {
				c.log.Debugf("unable to Send error to %d: %v", apdu.InvokeId, err)
			}
		default:
			// Ignore it
			c.log.WithFields(logrus.Fields{"raw": b}).Debug("An ignored packet went through")
		}
	}

	if header.Function == btypes.BacFuncForwardedNPDU {
		// Right now we are ignoring the NPDU data that is stored in the packet. Eventually
		// we will need to check it for any additional information we can gleam.
		// NDPU has source
		b = b[forwardHeaderLength:]
		c.log.Debugf("Ignored NDPU Forwarded:%v", b)
	}
}

type SetBroadcastType struct { //used to override the header.Function
	Set     bool
	BacFunc btypes.BacFunc
}

// Send transfers the raw apdu byte slice to the destination address.
func (c *client) Send(dest btypes.Address, npdu *btypes.NPDU, data []byte, broadcastType *SetBroadcastType) (int, error) {
	//broadcastType = &SetBroadcastType{}
	var header btypes.BVLC
	// Set packet type
	header.Type = btypes.BVLCTypeBacnetIP
	//if Adr is > 0 it must be an mst-tp device so send a UNICAST
	if len(dest.Adr) > 0 { //(aidan) not sure if this is correct, but it needs to be set to work to send (UNICAST) messages over a bacnet network
		// SET UNICAST FLAG
		// see http://www.bacnet.org/Tutorial/HMN-Overview/sld033.
		// see https://github.com/JoelBender/bacpypes/blob/9fca3f608a97a20807cd188689a2b9ff60b05085/doc/source/gettingstarted/gettingstarted001.rst#udp-communications-issues
		header.Function = btypes.BacFuncUnicast
	} else if dest.IsBroadcast() || dest.IsSubBroadcast() {
		// SET BROADCAST FLAG
		header.Function = btypes.BacFuncBroadcast
	} else {
		// SET UNICAST FLAG
		header.Function = btypes.BacFuncUnicast
	}

	if broadcastType != nil {
		if broadcastType.Set {
			header.Function = broadcastType.BacFunc
		}
	}

	header.Length = uint16(mtuHeaderLength + len(data))
	header.Data = data
	e := encoding.NewEncoder()
	err := e.BVLC(header)
	if err != nil {
		return 0, err
	}
	// use default udp type, src = network address (nil)
	return c.dataLink.Send(e.Bytes(), npdu, &dest)
}

func (c *client) ClientClose(closeLogs bool) error {
	if closeLogs {
		if f, ok := c.log.Out.(io.Closer); ok {
			return f.Close()
		}
	}
	return c.Close()
}

// Close free resources for the client. Always call this function when using NewClient
func (c *client) Close() error {
	if c.dataLink != nil {
		c.dataLink.Close()
	}

	return nil
}

// inject logger
func (c *client) SetLogger(logger *logrus.Logger) error {
	c.log = logger
	return nil
}

// inject logger
func (c *client) GetLogger() *logrus.Logger {
	return c.log
}
