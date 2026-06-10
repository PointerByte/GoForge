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
	viper.Set(string(LoggerSensibleKeysAtribute), map[string]any{"password": true, "token": false})
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
	gotMap, ok := got.(map[string]bool)
	if !ok || !gotMap["password"] || gotMap["token"] {
		t.Errorf("logger.sensibleKeys = %v, want password=true token=false", got)
	}
}

func TestGetStringMapBool(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	if got := getStringMapBool("missing"); len(got) != 0 {
		t.Fatalf("getStringMapBool(missing) = %#v, want empty map", got)
	}

	viper.Set("sensible", map[string]any{
		" password ": true,
		"token":      "false",
		"email":      "true",
		"phone":      "not-bool",
		"":           true,
		"count":      1,
	})

	got := getStringMapBool("sensible")
	if !got["password"] {
		t.Fatalf("password = false, want true in %#v", got)
	}
	if got["token"] {
		t.Fatalf("token = true, want false in %#v", got)
	}
	if !got["email"] {
		t.Fatalf("email = false, want true in %#v", got)
	}
	if got["phone"] || got["count"] {
		t.Fatalf("invalid values should be false in %#v", got)
	}
	if _, ok := got[""]; ok {
		t.Fatalf("blank key was not ignored in %#v", got)
	}
}
