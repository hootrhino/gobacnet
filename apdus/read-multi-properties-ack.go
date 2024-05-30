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

import "github.com/hootrhino/gobacnet/btypes"

// 返回列表
// Analog Input Object:
//     Object_Name: A character string that identifies the analog input.
//     Object_Identifier:
//     Object_Type: “analogInput”.
//     Present_Value: 1.23.
//     Status_Flags: Flags  (e.g., in alarm, overridden).
//     Event_State: The current state .
//     Out_of_Service: A boolean indicating whether the analog input is operational.
/*
*
* 这些属性是一个点位必备的
*
 */
func NewRequiredPropertiesResponse(InvokeId uint8, ObjectInstance uint32) btypes.MultiplePropertyData {
	return btypes.MultiplePropertyData{
		Objects: []btypes.Object{
			{
				ID: btypes.ObjectID{
					Type:     btypes.AnalogInput,
					Instance: btypes.ObjectInstance(ObjectInstance),
				},
				Properties: []btypes.Property{
					{
						Type: btypes.PROP_PRESENT_VALUE,
						Data: float32(3.14),
					},
					{
						Type: btypes.PROP_OBJECT_IDENTIFIER,
						Data: ObjectInstance,
					},
					{
						Type: btypes.PROP_OBJECT_NAME,
						Data: "PROP_OBJECT_NAME",
					},
					{
						Type: btypes.PROP_OUT_OF_SERVICE,
						Data: uint16(0),
					},
					{
						Type: btypes.PROP_STATUS_FLAGS,
						Data: uint16(0),
					},
					{
						Type: btypes.PROP_EVENT_STATE,
						Data: uint16(0),
					},
					{
						Type: btypes.PropUnits,
						Data: uint16(95),
					},
				},
			},
			{
				ID: btypes.ObjectID{
					Type:     btypes.AnalogInput,
					Instance: btypes.ObjectInstance(ObjectInstance),
				},
				Properties: []btypes.Property{
					{
						Type: btypes.PROP_DESCRIPTION,
						Data: "PROP_DESCRIPTION",
					},
				},
			},
		},
	}
}
