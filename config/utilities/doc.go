// Copyright 2026 PointerByte Contributors
// SPDX-License-Identifier: Apache-2.0

// Package utilities provides shared configuration-loading helpers for the
// config module.
//
// Its main responsibility is to load application settings from
// resources/application.yml, resources/application.yaml, or
// resources/application.json, merge .env and .env.local files, and apply
// environment-variable overrides derived from the Viper key paths.
//
// LoadEnv also accepts a directory that contains application.* directly for
// compatibility with modules that keep local example configuration beside
// their code.
//
// Main entry point:
//   - LoadEnv to load configuration files and apply override rules
package utilities
