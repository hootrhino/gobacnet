// Copyright (C) 2024 wwhai
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package apdus

// [129 10 0 31 1 8 0 20 6 1 0 0 0 0 0 48 1 12 12 2 0 0 1 25 76 41 0 62 33 4 63]
func NewReadPropertyListResponse(InvokeId uint8, deviceId uint32, propertiesCount uint8) ([]byte, error) {
	DeviceId, err1 := IntToInstanceId(deviceId)
	if err1 != nil {
		return nil, err1
	}
	// len = 31
	response := []byte{
		// BACnet Virtual Link Control------------------------------------------------
		0x81, 0x0a, 0x00, 0x00, /*这里会补长度位*/
		// NPDU-----------------------------------------------------------------------
		0x01, 0x00, //0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// APDU-----------------------------------------------------------------------
		0x30,                                              // Type: Complex-ACK (3)
		InvokeId,                                          // Invoke ID
		0x0c,                                              // Service Choice: readProperty (12)
		0x0c, 0x02, DeviceId[0], DeviceId[1], DeviceId[2], // ObjectIdentifier: device, id
		0x19, 0x4c, // Property Identifier: object-list (76)
		0x29,
		0x00,                                          // property Array Index (Unsigned) 0
		0x3e, 0x21, propertiesCount /*返回的属性数量*/, 0x3f, // [3]{...}[3]
	}
	Len, _ := IntToBVLCLen(uint16(len(response)))
	response[2] = Len[0]
	response[3] = Len[1]
	return response, nil
}

// 810a0023012400ff06aabbccddeeffff30260c0c02000001194c29013ec4000000003f
func NewReadPropertyResponse(InvokeId uint8, deviceId uint32, propertyId uint32) ([]byte, error) {

	DeviceId, err1 := IntToInstanceId(deviceId)
	if err1 != nil {
		return nil, err1
	}
	PropertyId, err1 := IntToInstanceId(propertyId)
	if err1 != nil {
		return nil, err1
	}

	// 810a0023
	// 01 24 00 ff 06 aabbccddeeff ff30260c0c02000001194c29013ec4000000003f
	response := []byte{
		// BACnet Virtual Link Control------------------------------------------------
		0x81, 0x0a, 0x00, 0x00, /*这里会补长度位*/
		// NPDU-----------------------------------------------------------------------
		0x01, 0x00, // 0x00, 0x01, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// APDU-----------------------------------------------------------------------
		0x30,                                              // Type: Complex-ACK (3)
		InvokeId,                                          // Invoke ID
		0x0c,                                              // Service Choice: readProperty (12)
		0x0c, 0x02, DeviceId[0], DeviceId[1], DeviceId[2], // ObjectIdentifier: device, 1314
		0x19, 0x4c, // Property Identifier: object-list (76)
		0x29,
		0,                                                                 // property Array Index (Unsigned) 0
		0x3e, 0xC4, 00, PropertyId[0], PropertyId[1], PropertyId[2], 0x3f, // [3]{...}[3]
	}
	Len, _ := IntToBVLCLen(uint16(len(response)))
	response[2] = Len[0]
	response[3] = Len[1]
	return response, nil
}
