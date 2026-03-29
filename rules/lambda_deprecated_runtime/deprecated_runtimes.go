package lambda_deprecated_runtime

import (
	_ "embed"
	"encoding/json"
	"time"
)

// runtimesJSON contains lifecycle dates for Lambda runtimes, generated from
// the AWS Lambda runtimes documentation. Dates that were in the future at
// generation time are speculative and may have shifted since.
//
//go:embed deprecated_runtimes.json
var runtimesJSON []byte

type runtimeLifecycle struct {
	EndOfSupportDate time.Time  `json:"end_of_support_date"`
	BlockCreateDate  *time.Time `json:"block_create_date,omitempty"`
	BlockUpdateDate  *time.Time `json:"block_update_date,omitempty"`
}

type runtimesData struct {
	UpdatedAt time.Time                    `json:"updated_at"`
	Runtimes  map[string]runtimeLifecycle  `json:"runtimes"`
}

var runtimes runtimesData

func init() {
	if err := json.Unmarshal(runtimesJSON, &runtimes); err != nil {
		panic(err)
	}
}
