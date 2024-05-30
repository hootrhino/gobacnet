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

package bacnet

import (
	"fmt"
	"math"

	"github.com/hootrhino/bacnet/btypes"
)

// Server Ip模式
type BacnetIPServer struct {
	deviceId     uint32
	vendorId     uint32
	netWorkId    uint16
	PropertyData map[uint32][2]btypes.Object // 点位, 后续对外提供的服务就是这个
}

func NewBacnetIPServer(
	deviceId uint32,
	vendorId uint32,
	netWorkId uint16,
	PropertyData map[uint32][2]btypes.Object,
) *BacnetIPServer {
	return &BacnetIPServer{
		deviceId:     deviceId,
		vendorId:     vendorId,
		netWorkId:    netWorkId,
		PropertyData: PropertyData,
	}
}

// 根据Bacnet请求查点位里面的数据
func (s *BacnetIPServer) GetObjectProperties(ObjectInstanceId uint32) (btypes.MultiplePropertyData, error) {
	Objects, ok := s.PropertyData[ObjectInstanceId]
	if ok {
		return btypes.MultiplePropertyData{
			Objects: Objects[:],
		}, nil
	}
	return btypes.MultiplePropertyData{}, fmt.Errorf(" Object InstanceId not exists")
}

// 更新值
func (s *BacnetIPServer) UpdateAIPropertyValue(InstanceId uint32, value interface{}) error {
	// 这个地方应该分两种情况
	// 1 来自对象的必须属性（7类）
	// 2 额外属性
	// 但是当前阶段暂时支持必须属性
	ObjectProperties, errGetObjectProperties := s.GetObjectProperties(uint32(InstanceId))
	if errGetObjectProperties != nil {
		return errGetObjectProperties
	}
	switch T := value.(type) {
	case int32:
		ObjectProperties.Objects[0].Properties[0].Data = T
	case uint32:
		ObjectProperties.Objects[0].Properties[0].Data = T
	case float32: // float32 需要IEEE浮点数算法
		ObjectProperties.Objects[0].Properties[0].Data = math.Float32bits(T)
	default:
		// 不允许出现别的类型, 从根源上直接崩溃
		panic(fmt.Errorf(" Unsupported Type: %v", T))
	}
	return nil
}
