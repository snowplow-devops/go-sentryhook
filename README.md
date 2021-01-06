# Sentry hook for Logrus

[![Release][release-image]][releases]

A simple hook for Logrus to allow easy integration of Sentry.  The hook leverages the [errors](https://github.com/pkg/errors) package under the hood to allow for stacktraces to be extracted.

## How to use?

```golang
import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/getsentry/sentry-go"

	"github.com/snowplow-devops/go-sentryhook"
)

func main() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://your-dsn/1",
		AttachStacktrace: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	logrus.AddHook(sentryhook.New([]log.Level{log.PanicLevel, log.FatalLevel, log.ErrorLevel}))

	// Simple error logged
	logrus.Error("This will be sent to Sentry!")

	// Error logged with stracktrace
	// Note: The error has to be included as an extra field for the stacktrace to be extracted
	errWithStacktrace := errors.New("This will be sent to Sentry with a StrackTrace!")
	logrus.WithFields(logrus.Fields{"error": errWithStacktrace}).Error(errWithStacktrace)
}
```

[release-image]: http://img.shields.io/badge/golang-0.1.0-6ad7e5.svg?style=flat
[releases]: https://github.com/snowplow-devops/go-sentryhook/releases/
