//
// (C) Copyright 2018 Intel Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// GOVERNMENT LICENSE RIGHTS-OPEN SOURCE SOFTWARE
// The Government's rights to use, modify, reproduce, release, perform, display,
// or disclose this software are subject to the terms of the Apache License as
// provided in Contract No. 8F-30005.
// Any reproduction of computer software, computer software documentation, or
// portions thereof marked with this legend must also reproduce the markings.
//

// Package ipmctl provides Go bindings for libipmctl Native Management API
package ipmctl

// CGO_CFLAGS & CGO_LDFLAGS env vars can be used
// to specify additional dirs.

/*
#cgo LDFLAGS: -lipmctl

#include "stdlib.h"
#include "nvm_management.h"
#include "nvm_types.h"
#include "NvmSharedDefs.h"
#include "export_api.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

//TEST_F(NvmApi_Tests, GetDeviceStatus)
//{
//  unsigned int dimm_cnt = 0;
//  nvm_get_number_of_devices(&dimm_cnt);
//  device_discovery *p_devices = (device_discovery *)malloc(sizeof(device_discovery) * dimm_cnt);
//  nvm_get_devices(p_devices, dimm_cnt);
//  device_status *p_status = (device_status *)malloc(sizeof(device_status));
//  EXPECT_EQ(nvm_get_device_status(p_devices->uid, p_status), NVM_SUCCESS);
//  free(p_status);
//  free(p_devices);
//}

// DeviceDiscovery struct represents Go equivalent of C.struct_device_discovery
// from nvm_management.h (NVM API) as reported by "go tool cgo -godefs nvm.go"
type DeviceDiscovery struct {
	All_properties_populated uint8
	Pad_cgo_0                [3]byte
	Device_handle            [4]byte
	Physical_id              uint16
	Vendor_id                uint16
	Device_id                uint16
	Revision_id              uint16
	Channel_pos              uint16
	Channel_id               uint16
	Memory_controller_id     uint16
	Socket_id                uint16
	Node_controller_id       uint16
	Pad_cgo_1                [2]byte
	Memory_type              uint32
	Dimm_sku                 uint32
	Manufacturer             [2]uint8
	Serial_number            [4]uint8
	Subsystem_vendor_id      uint16
	Subsystem_device_id      uint16
	Subsystem_revision_id    uint16
	Manufacturing_info_valid uint8
	Manufacturing_location   uint8
	Manufacturing_date       uint16
	Part_number              [21]int8
	Fw_revision              [25]int8
	Fw_api_version           [25]int8
	Pad_cgo_2                [5]byte
	Capacity                 uint64
	Interface_format_codes   [9]uint16
	Security_capabilities    _Ctype_struct_device_security_capabilities
	Device_capabilities      _Ctype_struct_device_capabilities
	Uid                      [22]int8
	Lock_state               uint32
	Manageability            uint32
	Controller_revision_id   uint16
	Reserved                 [48]uint8
	Pad_cgo_3                [6]byte
}

// GetDevices queries number of NVDIMMS and retrieves device_discovery structs
// for each.
func GetDevices() error {
	var count C.uint
	C.nvm_get_number_of_devices(&count)
	if count == 0 {
		println("no NVDIMMs found!")
		return nil
	}

	devs := make([]C.struct_device_discovery, int(count))
	println(len(devs))

	// don't need to defer free on devs as we allocated in go
	C.nvm_get_devices(&devs[0], C.NVM_UINT8(count))
	// defer C.free(unsafe.Pointer(&devs))

	// cast struct array to slice of go equivalent struct
	// (get equivalent go struct def from cgo -godefs)
	deviceSlice := (*[1 << 30]DeviceDiscovery)(unsafe.Pointer(&devs[0]))[:count:count]
	for i := 0; i < int(count); i++ {
		//item := (*DeviceDiscovery)(unsafe.Pointer(&devs[i]))
		d := deviceSlice[i]
		fmt.Printf("Device ID: %d, Memory type: %d, Fw Rev: %v, Capacity %d, ",
			d.Device_id, d.Memory_type, d.Fw_revision, d.Capacity)
		fmt.Printf("Channel Pos: %d, Channel ID: %d, Memory Ctrlr: %d, Socket ID: %d.\n",
			d.Channel_pos, d.Channel_id, d.Memory_controller_id, d.Socket_id)
	}

	fmt.Printf("%s\n", C.GoString((*C.char)(unsafe.Pointer(&devs[0].uid))))
	status := C.struct_device_status{}
	C.nvm_get_device_status((*C.char)(unsafe.Pointer(&devs[0].uid)), &status)

	uidCharPtr := (*C.char)(unsafe.Pointer(&devs[0].uid))

    //status := C.struct_device_status{}
    //C.nvm_get_device_status(uidCharPtr, &status)

	// verify api call passing in uid as param
    dev := C.struct_device_discovery{}
    C.nvm_get_device_discovery(uidCharPtr, &dev)
    dd := (*DeviceDiscovery)(unsafe.Pointer(&dev))
    fmt.Printf("Device ID: %d, Memory type: %d, Fw Rev: %v, Capacity %d, ",
               dd.Device_id, dd.Memory_type, dd.Fw_revision, dd.Capacity)


	return nil
}

// Rc2err returns an failure if rc != 0.
//
// TODO: If err is already set then it is wrapped,
// otherwise it is ignored. e.g.
// func Rc2err(label string, rc C.int, err error) error {
func Rc2err(label string, rc C.int) error {
	if rc != 0 {
		if rc < 0 {
			rc = -rc
		}
		// e := errors.Error(rc)
		return fmt.Errorf("%s: %s", label, rc) // e
	}
	return nil
}
