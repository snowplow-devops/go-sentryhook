//
// Copyright (c) 2021 Snowplow Analytics Ltd. All rights reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
// The implementation is derived from https://github.com/makasim/sentryhook
//
// Copyright (c) 2020 Max Kotliar
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package sentryhook

import (
	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"reflect"
	"strings"
	"time"
)

var (
	logToSentryMap = map[logrus.Level]sentry.Level{
		logrus.TraceLevel: sentry.LevelDebug,
		logrus.DebugLevel: sentry.LevelDebug,
		logrus.InfoLevel:  sentry.LevelInfo,
		logrus.WarnLevel:  sentry.LevelWarning,
		logrus.ErrorLevel: sentry.LevelError,
		logrus.FatalLevel: sentry.LevelFatal,
		logrus.PanicLevel: sentry.LevelFatal,
	}
)

// Hook contains the structure for a Logrus hook
type Hook struct {
	hub    *sentry.Hub
	levels []logrus.Level
}

// New returns a new hook for use by Logrus
func New(levels []logrus.Level) Hook {
	return Hook{
		levels: levels,
		hub:    sentry.CurrentHub(),
	}
}

// Levels returns the levels that this hook fires on
func (hook Hook) Levels() []logrus.Level {
	return hook.levels
}

// Fire sends an event to Sentry
func (hook Hook) Fire(entry *logrus.Entry) error {
	event := sentry.NewEvent()

	event.Level = logToSentryMap[entry.Level]
	event.Message = entry.Message

	for k, v := range entry.Data {
		if k != logrus.ErrorKey {
			event.Extra[k] = v
		}
	}

	if err, ok := entry.Data[logrus.ErrorKey].(error); ok {
		// Use the final message as the error "type"
		lastMsg := strings.Split(err.Error(), ":")[0]

		// Extract the cause to set the base error type
		cause := errors.Cause(err)

		exception := sentry.Exception{
			Type:  lastMsg,
			Value: reflect.TypeOf(cause).String(),
		}

		if hook.hub.Client().Options().AttachStacktrace {
			exception.Stacktrace = sentry.ExtractStacktrace(err)
		}

		event.Exception = []sentry.Exception{exception}
	}

	hook.hub.CaptureEvent(event)

	if entry.Level == logrus.FatalLevel {
		hook.hub.Flush(2 * time.Second)
	}

	return nil
}
