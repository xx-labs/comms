///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package connect

// This file is only compiled for WebAssembly.

// GetDefaultConnectionType returns Web as the default connection type when
// compiling for WebAssembly.
func GetDefaultConnectionType() ConnectionType {
	return Web
}
