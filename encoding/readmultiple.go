package encoding

import (
	"fmt"
	"github.com/hootrhino/gobacnet/btypes"
)

func (e *Encoder) ReadMultipleProperty(invokeID uint8, data btypes.MultiplePropertyData) error {
	a := btypes.APDU{
		DataType:         btypes.ConfirmedServiceRequest,
		Service:          btypes.ServiceConfirmedReadPropMultiple,
		MaxSegs:          0,
		MaxApdu:          MaxAPDU,
		InvokeId:         invokeID,
		SegmentedMessage: false,
	}
	e.APDU(a)
	err := e.objects(data.Objects, false)
	if err != nil {
		return err
	}

	return e.Error()
}

func (d *Decoder) ReadMultipleProperty(data *btypes.MultiplePropertyData) error {
	// Must have at least 7 bytes
	if d.buff.Len() < 7 {
		return fmt.Errorf("missing parameters")
	}

	data.Objects = make([]btypes.Object, 0)
	objectIdx := 0
	newObject := false
	var tag uint8
	var meta tagMeta
	for {
		if d.buff.Len() == 0 {
			break
		}
		if newObject == false {
			// first object identifier
			tag, meta = d.tagNumber() // tag number = 0
		}
		newObject = false
		fmt.Printf("isSpecifitag=%v tagNumber=%v len=%v\n", meta.isContextSpecific(), tag, d.value(meta))
		objectType, instance := d.objectId()
		fmt.Println(objectType, instance)
		data.Objects = append(data.Objects, btypes.Object{
			ID: btypes.ObjectID{
				Type:     objectType,
				Instance: instance,
			},
			Properties: []btypes.Property{},
		})

		// list of property references
		for {
			if d.buff.Len() == 0 {
				break
			}
			tag, meta := d.tagNumber() // tag number = 1
			if tag == 0 {
				// 另一个object
				newObject = true
				break
			}
			fmt.Printf("isSpecifitag=%v tagNumber=%v len=%v\n", meta.isContextSpecific(), tag, d.value(meta))

			tag, meta = d.tagNumber() // tag number = 0
			fmt.Printf("isSpecifitag=%v tagNumber=%v len=%v\n", meta.isContextSpecific(), tag, d.value(meta))

			var prop btypes.Property
			prop.Type = btypes.PropertyType(d.enumerated(int(d.value(meta))))
			fmt.Printf("property=%v\n", prop)
			data.Objects[objectIdx].Properties = append(data.Objects[objectIdx].Properties, prop)

			tag, meta = d.tagNumber() // tag number = 1
			fmt.Printf("isSpecifitag=%v tagNumber=%v len=%v\n", meta.isContextSpecific(), tag, d.value(meta))
		}

		objectIdx++
	}

	return nil
}
