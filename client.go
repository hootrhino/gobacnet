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

package gobacnet

import (
	"context"
	"io"

	"github.com/hootrhino/gobacnet/btypes"
	"github.com/sirupsen/logrus"
)

type Client interface {
	io.Closer
	ClientRun()
	ClientClose(closeLogs bool) error
	StartPoll(ctx context.Context)
	StopPoll()
	WhoIs(wh *WhoIsOpts) ([]btypes.Device, error)
	WhatIsNetworkNumber() []*btypes.Address
	IAm(dest btypes.Address, iam btypes.IAm) error
	WhoIsRouterToNetwork() (resp *[]btypes.Address)
	Objects(dev btypes.Device) (btypes.Device, error)
	ReadProperty(dest btypes.Device, rp btypes.PropertyData) (btypes.PropertyData, error)
	ReadMultiProperty(dev btypes.Device, rp btypes.MultiplePropertyData) (btypes.MultiplePropertyData, error)
	WriteProperty(dest btypes.Device, wp btypes.PropertyData) error
	WriteMultiProperty(dev btypes.Device, wp btypes.MultiplePropertyData) error
	SetLogger(*logrus.Logger) error
	GetLogger() *logrus.Logger
	GetBacnetIPServer() *BacnetIPServer
}
