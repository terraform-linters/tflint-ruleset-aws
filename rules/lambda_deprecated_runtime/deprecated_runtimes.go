package lambda_deprecated_runtime

import (
	_ "embed"
	"encoding/json"
	"time"
)

// deprecatedRuntimesJSON contains lifecycle dates for Lambda runtimes that have
// already reached end of support. Supported runtimes are intentionally excluded
// because their deprecation dates are subject to change.
//
//go:embed deprecated_runtimes.json
var deprecatedRuntimesJSON []byte

type deprecatedRuntime struct {
	EndOfSupportDate time.Time  `json:"end_of_support_date"`
	BlockCreateDate  *time.Time `json:"block_create_date,omitempty"`
	BlockUpdateDate  *time.Time `json:"block_update_date,omitempty"`
}

var deprecatedRuntimes map[string]deprecatedRuntime

func init() {
	if err := json.Unmarshal(deprecatedRuntimesJSON, &deprecatedRuntimes); err != nil {
		panic(err)
	}
}
