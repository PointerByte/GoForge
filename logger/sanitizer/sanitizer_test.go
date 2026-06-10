// Copyright 2026 PointerByte Contributors
// SPDX-License-Identifier: Apache-2.0

package sanitizer

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/PointerByte/GoForge/logger/formatter"
	viperdata "github.com/PointerByte/GoForge/logger/viperData"
	"github.com/spf13/viper"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Public   string `json:"public"`
}

func TestLogFormatRedactsSensitiveKeys(t *testing.T) {
	s := New(map[string]bool{
		"password": true,
		"email":    true,
		"phone":    true,
		"token":    false,
	})

	log := formatter.LogFormat{
		Message: "login failed password=message-secret token=public",
		Details: formatter.Details{
			System:  "api",
			Headers: http.Header{"X-Email": {"person@example.com"}, "X-Trace": {`{"phone":"555"}`}},
			Request: map[string]any{
				"password": "secret",
				"profile": map[string]any{
					"user_email": "person@example.com",
					"name":       "Ada",
				},
			},
			Response: `{"ok":false,"phone":"555-0100"}`,
		},
		Services: []formatter.Service{
			{
				System:  "auth",
				Request: loginRequest{Email: "service@example.com", Password: "service-secret", Public: "ok"},
			},
		},
	}

	got := s.LogFormat(log)
	data, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("json.Marshal(sanitized log) failed: %v", err)
	}
	out := string(data)

	for _, secret := range []string{"message-secret", "person@example.com", "secret", "555-0100", "service@example.com", "service-secret"} {
		if strings.Contains(out, secret) {
			t.Fatalf("sanitized log still contains %q: %s", secret, out)
		}
	}
	if count := strings.Count(out, RedactedValue); count < 6 {
		t.Fatalf("redaction count = %d, want at least 6 in %s", count, out)
	}
	if !strings.Contains(out, "token=public") {
		t.Fatalf("disabled key token was redacted unexpectedly: %s", out)
	}

	originalRequest := log.Details.Request.(map[string]any)
	if originalRequest["password"] != "secret" {
		t.Fatalf("original request was mutated: %#v", originalRequest)
	}
}

func TestValueSanitizesRawStrings(t *testing.T) {
	s := New(map[string]bool{"password": true, "email": true})

	tests := []struct {
		name      string
		value     string
		notWanted []string
	}{
		{
			name:      "json string",
			value:     `{"password":"secret","profile":{"email":"person@example.com"},"ok":true}`,
			notWanted: []string{"secret", "person@example.com"},
		},
		{
			name:      "key value string",
			value:     "password=secret email=person@example.com ok=true",
			notWanted: []string{"secret", "person@example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := s.Value(tt.value).(string)
			if !ok {
				t.Fatalf("Value(%q) type = %T, want string", tt.value, got)
			}
			for _, secret := range tt.notWanted {
				if strings.Contains(got, secret) {
					t.Fatalf("sanitized string still contains %q: %s", secret, got)
				}
			}
			if !strings.Contains(got, RedactedValue) {
				t.Fatalf("sanitized string missing redaction marker: %s", got)
			}
		})
	}
}

func TestFromViper(t *testing.T) {
	viper.Reset()
	viperdata.ResetViperDataSingleton()
	t.Cleanup(func() {
		viper.Reset()
		viperdata.ResetViperDataSingleton()
	})

	viper.Set(string(viperdata.LoggerSensibleKeysAtribute), map[string]any{
		"password": true,
		"token":    false,
	})

	s := FromViper()
	got := s.Value(map[string]any{
		"password": "secret",
		"token":    "visible",
	}).(map[string]any)

	if got["password"] != RedactedValue {
		t.Fatalf("password = %#v, want redacted", got["password"])
	}
	if got["token"] != "visible" {
		t.Fatalf("token = %#v, want visible", got["token"])
	}
}

func TestDisabledSanitizerReturnsOriginalValue(t *testing.T) {
	s := New(nil)
	value := map[string]any{"password": "secret"}

	got := s.Value(value)
	if got == nil || got.(map[string]any)["password"] != "secret" {
		t.Fatalf("disabled sanitizer changed value: %#v", got)
	}
}
