// Copyright 2026 PointerByte Contributors
// SPDX-License-Identifier: Apache-2.0

package builder

import (
	"context"
	"encoding/json"
	"io"
	"maps"
	"sync"
	"time"

	"log/slog"

	"github.com/PointerByte/GoForge/logger/formatter"
	"github.com/PointerByte/GoForge/logger/sanitizer"
	viperdata "github.com/PointerByte/GoForge/logger/viperData"
)

type jsonHandler struct {
	level    slog.Level
	w        io.Writer
	mux      *sync.Mutex
	attrs    []slog.Attr
	groups   []string
	handlers []slog.Handler
}

func newHandler(level slog.Level, w io.Writer, handlers ...slog.Handler) *jsonHandler {
	return &jsonHandler{level: level, w: w, mux: &sync.Mutex{}, handlers: handlers}
}

func (h *jsonHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *jsonHandler) Handle(ctx context.Context, r slog.Record) error {
	data := make(map[string]any)
	ctxLogger := New(ctx)
	maps.Copy(data, ctxLogger.customLogFormat())
	layout := viperdata.GetViperData(string(viperdata.LoggerFormatDateAtribute)).(string)
	data[string(timestampAtribute)] = time.Now().Format(layout)
	data[string(loggerMessage)] = r.Message
	data[string(levelAtribute)] = r.Level.String()

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	var logObj formatter.LogFormat
	if err = json.Unmarshal(jsonBytes, &logObj); err != nil {
		panic(err)
	}

	formatterAtribute := viperdata.GetViperData(string(viperdata.LoggerFormatterAtribute)).(string)
	formatter := formatter.New(formatterAtribute)
	logObj.Latency = ctxLogger.GetLatency()
	logObj = sanitizer.FromViper().LogFormat(logObj)
	jsonBytes, err = formatter.Format(logObj)
	if err != nil {
		panic(err)
	}
	// The formatter now owns a snapshot of completed traces. Clear only after
	// that final payload has been built, so TraceEnd entries cannot disappear
	// before serialization.
	ctxLogger.clearProcesses(len(logObj.Process))
	var jsonMap any
	if err := json.Unmarshal(jsonBytes, &jsonMap); err != nil {
		if err = h.writeData(jsonBytes); err != nil {
			panic(err)
		}
		return nil
	}
	if err = h.writeData(jsonBytes); err != nil {
		panic(err)
	}
	return nil
}

func (h *jsonHandler) writeData(jsonBytes []byte) error {
	line := make([]byte, 0, len(jsonBytes)+1)
	line = append(line, jsonBytes...)
	line = append(line, '\n')

	if h.mux != nil {
		h.mux.Lock()
		defer h.mux.Unlock()
	}

	_, err := h.w.Write(line)
	if err != nil {
		return err
	}
	return nil
}

func (h *jsonHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	clone := *h
	clone.attrs = append(clone.attrs, attrs...)
	return &clone
}

func (h *jsonHandler) WithGroup(name string) slog.Handler {
	clone := *h
	clone.groups = append(clone.groups, name)
	return &clone
}
