package encoding

import (
	"fmt"

	"github.com/hootrhino/gobacnet/btypes"
)

func isValidObjectType(idType btypes.ObjectType) error {
	if idType > MaxObject {
		return fmt.Errorf("Object btypes is %d which must be less then %d", idType, MaxObject)
	}
	return nil
}

func isValidPropertyType(propType uint32) error {
	if propType > MaxPropertyID {
		return fmt.Errorf("Object btypes is %d which must be less then %d", propType, MaxPropertyID)
	}
	return nil
}
