// Copyright 2026 PointerByte Contributors
// SPDX-License-Identifier: Apache-2.0

package viperdata

import (
	"testing"

	"github.com/spf13/viper"
)

func TestGetViperData(t *testing.T) {
	ResetViperDataSingleton()

	viper.Reset()
	t.Cleanup(viper.Reset)

	viper.Set(string(AppAtribute), "test-app")
	viper.Set(string(LoggerModeTestAtribute), true)
	viper.Set(string(LoggerSensibleKeysAtribute), []string{"password"})
	viper.Set(string(GRPCLoggerWithConfigSkipFunctionAtribute), []string{"SayHello"})

	got := GetViperData(string(AppAtribute))
	if got != "test-app" {
		t.Errorf("app.name = %v, want %v", got, "test-app")
	}

	got = GetViperData(string(LoggerModeTestAtribute))
	if got != true {
		t.Errorf("logger.modeTest = %v, want %v", got, "true")
	}

	got = GetViperData(string(GRPCLoggerWithConfigSkipFunctionAtribute))
	if gotSlice, ok := got.([]string); !ok || len(gotSlice) != 1 || gotSlice[0] != "SayHello" {
		t.Errorf("server.grpc.LoggerWithConfig.SkipFunction = %v, want [SayHello]", got)
	}

	got = GetViperData(string(LoggerSensibleKeysAtribute))
	gotSlice, ok := got.([]string)
	if !ok || len(gotSlice) != 1 || gotSlice[0] != "password" {
		t.Errorf("logger.sensibleKeys = %v, want [password]", got)
	}
}

func TestIsIgnoredHeader(t *testing.T) {
	ResetViperDataSingleton()
	viper.Reset()
	t.Cleanup(func() {
		viper.Reset()
		ResetViperDataSingleton()
	})

	viper.Set(string(LoggerIgnoredHeadersAtribute), []string{
		"authorization",
		"COOKIE",
	})

	tests := []struct {
		name   string
		header string
		want   bool
	}{
		{name: "matches lower configured header", header: "Authorization", want: true},
		{name: "matches upper configured header", header: "Cookie", want: true},
		{name: "does not match substring", header: "X-Authorization-Token", want: false},
		{name: "does not match unrelated header", header: "Content-Type", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIgnoredHeader(tt.header); got != tt.want {
				t.Fatalf("IsIgnoredHeader(%q) = %v, want %v", tt.header, got, tt.want)
			}
		})
	}
}
