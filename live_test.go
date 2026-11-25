// Copyright 2026 High Five Labs
// SPDX-License-Identifier: Apache-2.0

//go:build live

package extensions

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/highfivelabs/krm-rest-runtime/kindtest"
)

// The kinds hack/live/live.yaml lists, against the real API: create,
// update, drift repair, delete. Compiled only with -tags live, and skips,
// saying why, unless KRM_LIVE=1 and credentials are found.
// Results go to hack/live/records.json, which sets each kind's readiness.
func TestEveryKindLive(t *testing.T) {
	if os.Getenv(kindtest.EnvLive) == "1" {
		loadSandboxIntegrationEnv(t)
	}
	kindtest.Live(t, Types, Vendor, nil, liveYAML)
}

// Deletes every live-test leftover (KRM_LIVE=1 KRM_LIVE_SWEEP=1).
func TestLiveSweep(t *testing.T) {
	kindtest.LiveSweep(t, Vendor, liveYAML)
}

func loadSandboxIntegrationEnv(t *testing.T) {
	t.Helper()
	kinds := os.Getenv(kindtest.EnvLiveKinds)
	want := func(prefix string) bool {
		if kinds == "" {
			return true
		}
		for _, k := range strings.Split(kinds, ",") {
			if strings.HasPrefix(strings.TrimSpace(k), prefix) {
				return true
			}
		}
		return false
	}

	project := os.Getenv(kindtest.EnvLiveProject)
	if project == "" {
		t.Skipf("live: set $%s to the Secret Manager project holding the sandbox integrations' secrets", kindtest.EnvLiveProject)
	}

	if want("IntegrationCloudflare") {
		loadSecretToEnv(t, project, "krm-live-cloudflare", map[string]string{
			"apiToken": "DD_LIVE_CF_API_TOKEN",
			"zoneId":   "DD_LIVE_CF_ZONE_ID",
		})
	}
	if want("IntegrationConfluent") {
		loadSecretToEnv(t, project, "krm-live-confluent", map[string]string{
			"apiKey":    "DD_LIVE_CONFLUENT_API_KEY",
			"apiSecret": "DD_LIVE_CONFLUENT_API_SECRET",
		})
	}
}

func loadSecretToEnv(t *testing.T, project, secret string, keyMap map[string]string) {
	t.Helper()
	// A failed read is said, not swallowed: the integration kinds that need
	// these values fail, and the log should say why.
	if _, err := exec.LookPath("gcloud"); err != nil {
		t.Logf("live: %s not read: gcloud is not installed", secret)
		return
	}
	out, err := exec.Command("gcloud", "secrets", "versions", "access", "latest",
		"--secret="+secret, "--project="+project, "--quiet").Output()
	if err != nil {
		t.Logf("live: %s in project %s not read (the live runner may lack access to it): %v", secret, project, err)
		return
	}
	var payload map[string]string
	if err := json.Unmarshal(out, &payload); err != nil {
		t.Logf("live: %s is not a JSON object of strings: %v", secret, err)
		return
	}
	for secretKey, envVar := range keyMap {
		if v := payload[secretKey]; v != "" {
			t.Setenv(envVar, v)
		}
	}
}
