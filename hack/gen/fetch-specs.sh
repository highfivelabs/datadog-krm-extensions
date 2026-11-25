#!/usr/bin/env bash
# Copyright 2026 High Five Labs
# SPDX-License-Identifier: Apache-2.0
# Download the vendor's OpenAPI specs at the commits vendor.yaml pins.
set -euo pipefail
cd "$(dirname "$0")"
mkdir -p specs
[ -s specs/datadog-v1.yaml ] || curl -fsSL -o specs/datadog-v1.yaml "https://raw.githubusercontent.com/DataDog/datadog-api-client-go/8f7c1eede09a7216d358551a261b3286687e500d/.generator/schemas/v1/openapi.yaml"
[ -s specs/datadog-v2.yaml ] || curl -fsSL -o specs/datadog-v2.yaml "https://raw.githubusercontent.com/DataDog/datadog-api-client-go/8f7c1eede09a7216d358551a261b3286687e500d/.generator/schemas/v2/openapi.yaml"
