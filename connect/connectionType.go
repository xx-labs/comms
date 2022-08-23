///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

//go:build !js || !wasm
// +build !js !wasm

package connect

// This file is only compiled for all architectures except WebAssembly.

// GetDefaultConnectionType returns Grpc as the default connection type when
// compiling for all architectures except WebAssembly.
func GetDefaultConnectionType() ConnectionType {
	return Grpc
}
