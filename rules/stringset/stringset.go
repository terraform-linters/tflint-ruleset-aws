// Package stringset provides a membership test decoded from a JSON array of
// strings, the shape the generated AWS data files hold.
package stringset

import "encoding/json"

// Set reports membership over the strings it was decoded from.
type Set map[string]bool

func (s *Set) UnmarshalJSON(data []byte) error {
	var values []string
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}

	*s = make(Set, len(values))
	for _, value := range values {
		(*s)[value] = true
	}

	return nil
}
