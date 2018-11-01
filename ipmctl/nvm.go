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

func GetNumDevices() error {
	var count C.uint
	C.nvm_get_number_of_devices(&count)
	println(count)

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
