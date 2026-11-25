// Copyright 2026 High Five Labs
// SPDX-License-Identifier: Apache-2.0

package extensions

import (
	"net/http"

	"github.com/highfivelabs/krm-rest-runtime/rest"
)

// Vendor is how the manager reaches Datadog.
var Vendor = &rest.Vendor{
	Name:           "datadog",
	DefaultBaseURL: "https://api.datadoghq.com",
	CredentialKeys: []string{"apiKey", "appKey"},
	Authorize: func(r *http.Request, c map[string]string) {
		r.Header.Set("DD-API-KEY", c["apiKey"])
		r.Header.Set("DD-APPLICATION-KEY", c["appKey"])
	},
	// The Secret's optional "site" key (datadoghq.eu, us5.datadoghq.com, ...)
	// picks the region when the Connection names no base URL.
	BaseURL: func(c map[string]string) string {
		if s := c["site"]; s != "" {
			return "https://api." + s
		}
		return ""
	},
}
