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

// The BACnet instance number is a 22-bit value (valid range from 0 – 4194302, 4194303 is reserved as wildcard and cannot be used).
// Combined with the object type (10-bit), this information creates the unique address of the device (the Device Object Identifier).
func NewIAm(deviceId uint32, vendorId uint16, networkId uint16) ([]byte, error) {
	DeviceInstanceId, err1 := IntToInstanceId(deviceId)
	if err1 != nil {
		return nil, err1
	}
	VendorId, err2 := IntToVendorId(vendorId)
	if err2 != nil {
		return nil, err2
	}
	NetworkId, err2 := IntToNetworkId(networkId)
	if err2 != nil {
		return nil, err2
	}
	iAm := [25]byte{
		// BACnet Virtual Link Control------------------------------------------------
		0x81, 0x0a, 0x00, 0x1d,
		// NPDU-----------------------------------------------------------------------
		0x01, 0x20, NetworkId[0], NetworkId[1], 00, 0xFF, /*Hop*/
		// APDU-----------------------------------------------------------------------
		0x10,                                                                      // APDU Type
		0x00,                                                                      // Service Choice: 00 IAM
		0xc4, 0x02, DeviceInstanceId[0], DeviceInstanceId[1], DeviceInstanceId[2], // DeviceID
		0x22, 0x05, 0xc4, // MAX APDU LENGTH
		0x91, 0x00, // Segment
		0x22, VendorId[0], VendorId[1], // VendorId 0x22 2Byte
	}
	Len, _ := IntToBVLCLen(uint16(len(iAm)))
	iAm[2] = Len[0]
	iAm[3] = Len[1]
	return iAm[:], nil
}
