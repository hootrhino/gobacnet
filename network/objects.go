package network

import (
	"fmt"

	"github.com/hootrhino/bacnet"
	"github.com/hootrhino/bacnet/btypes"
	"github.com/hootrhino/bacnet/helpers/data"
	log "github.com/sirupsen/logrus"
)

// DeviceObjects get device objects
func (device *Device) DeviceObjects(deviceID btypes.ObjectInstance, checkAPDU bool) (objectList []btypes.ObjectID, err error) {
	if checkAPDU { //check set the maxADPU and Segmentation
		whoIs, err := device.Whois(&bacnet.WhoIsOpts{
			High:            int(deviceID),
			Low:             int(deviceID),
			GlobalBroadcast: true,
		})
		if err != nil {
			return nil, err
		}
		for _, dev := range whoIs {
			if dev.ID.Instance == deviceID {
				device.MaxApdu = dev.MaxApdu
				device.Segmentation = uint32(dev.Segmentation)
			}
		}
		log.Infoln("bacnet.DeviceObjects() do whois on deviceID:", deviceID, " maxADPU:", device.MaxApdu, " Segmentation:", device.Segmentation)
	}

	//get object list
	obj := &Object{
		ObjectID:   deviceID,
		ObjectType: btypes.DeviceType,
		Prop:       btypes.PropObjectList,
		ArrayIndex: btypes.ArrayAll, //btypes.ArrayAll

	}
	out, err := device.Read(obj)
	if err != nil { //this is a device that would have a low maxADPU
		fmt.Println("DeviceObjects, now read here", err)
		return device.deviceObjectsBuilder(deviceID)

	}
	if len(out.Object.Properties) == 0 {
		fmt.Println("No value returned")
		return nil, nil
	}
	_, ids := data.ToArr(out)
	for _, id := range ids {
		objectID := id.(btypes.ObjectID)
		objectList = append(objectList, objectID)
	}
	return objectList, nil
}

// DeviceObjectsBuilder this is used when a device can't send the object list in the fully ArrayIndex
// it first reads the size of the object list and then loops the list to build an object list
func (device *Device) deviceObjectsBuilder(deviceID btypes.ObjectInstance) (objectList []btypes.ObjectID, err error) {
	//get object list
	obj := &Object{
		ObjectID:   deviceID,
		ObjectType: btypes.DeviceType,
		Prop:       btypes.PropObjectList,
		ArrayIndex: 0, //start at 0 and then loop through
	}
	out, err := device.Read(obj)
	if err != nil {
		log.Errorln("failed to read object list in deviceObjectsBuilder() err:", err)
		return nil, err
	}
	_, o := data.ToUint32(out)
	log.Println("size of object-list", o)
	var listLen = int(o)
	for i := 1; i <= listLen; i++ {
		obj.ArrayIndex = uint32(i)
		out, _ := device.Read(obj)
		objectID := out.Object.Properties[0].Data.(btypes.ObjectID)
		objectList = append(objectList, objectID)
	}
	return objectList, err

}
