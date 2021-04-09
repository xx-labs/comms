///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package connect

import (
	"time"
)

// Params object for host creation
type HostParams struct {
	MaxRetries  uint32
	AuthEnabled bool

	// Toggles cool off of connections
	EnableCoolOff bool

	// Number of leaky bucket sends before it stops
	NumSendsBeforeCoolOff uint32

	// Amount of time after a cool off is triggered before allowed to send again
	CoolOffTimeout time.Duration

	// If set, metric handling will be enabled on this host
	EnableMetrics bool

	// List of sending errors that are deemed unimportant
	// Reception of these errors will not update the Metric's state
	ExcludeMetricErrors []string
}

// Get default set of host params
func GetDefaultHostParams() HostParams {
	return HostParams{
		MaxRetries:            100,
		AuthEnabled:           true,
		EnableCoolOff:         false,
		NumSendsBeforeCoolOff: 3,
		CoolOffTimeout:        60 * time.Second,
		EnableMetrics:         false,
		ExcludeMetricErrors:   make([]string, 0),
	}
}
