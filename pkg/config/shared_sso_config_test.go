package config

import (
	"encoding/json"
	"testing"
)

func TestSharedSSOConfigMalformed(t *testing.T) {
	for name, body := range map[string]string{
		"provider missing":  `{"name":"github"}`,
		"provider a number": `{"name":"github","provider":42}`,
		"name a number":     `{"name":42,"provider":"GITHUB"}`,
	} {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("panicked on user config: %v", r)
				}
			}()
			var c SharedSSOConfig
			_ = json.Unmarshal([]byte(body), &c)
		})
	}
}
