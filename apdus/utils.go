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

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// IntToIndex
func IntToIndex(n uint32) ([]byte, error) {
	if n > 4194302 {
		return nil, fmt.Errorf(" InstanceId must be in the range of 0 to 4194302")
	}
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, n)
	if err != nil {
		return nil, err
	}

	return []byte{buf.Bytes()[1], buf.Bytes()[2], buf.Bytes()[3]}, nil
}

// IntToInstanceId
func IntToInstanceId(n uint32) ([]byte, error) {
	if n > 4194302 {
		return nil, fmt.Errorf(" InstanceId must be in the range of 0 to 4194302")
	}
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, n)
	if err != nil {
		return nil, err
	}

	return []byte{buf.Bytes()[1], buf.Bytes()[2], buf.Bytes()[3]}, nil
}

// IntToVendorId
func IntToVendorId(n uint16) ([]byte, error) {
	if n > 65535 {
		return nil, fmt.Errorf(" VendorId must be in the range of 0 to 65535")
	}
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, n)
	if err != nil {
		return nil, err
	}
	return []byte{buf.Bytes()[0], buf.Bytes()[1]}, nil
}

// IntToNetworkId
func IntToNetworkId(n uint16) ([]byte, error) {
	if n > 65535 {
		return nil, fmt.Errorf(" NetworkId must be in the range of 0 to 65535")
	}
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, n)
	if err != nil {
		return nil, err
	}
	return []byte{buf.Bytes()[0], buf.Bytes()[1]}, nil
}

// IntToUint16
func IntToBVLCLen(n uint16) ([]byte, error) {
	if n > 65535 {
		return nil, fmt.Errorf(" VendorId must be in the range of 0 to 65535")
	}
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, n)
	if err != nil {
		return nil, err
	}
	return []byte{buf.Bytes()[0], buf.Bytes()[1]}, nil
}
