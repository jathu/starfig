package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/sirupsen/logrus"
)

func SetupLogger() {
	_, inDebug := os.LookupEnv(starfigDebugEnvvar)
	if inDebug {
		logrus.SetLevel(logrus.TraceLevel)
	}

	logrus.SetFormatter(new(starfigLoggerFormatter))
	logrus.SetReportCaller(inDebug)
}

// Private ---------------------------------------------------------------------

const starfigDebugEnvvar string = "STARFIG_DEBUG"

type starfigLoggerFormatter struct{}

func (f *starfigLoggerFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var components []string

	components = append(components,
		aurora.Faint(entry.Time.Format("15:04:05.000")).String())

	switch entry.Level {
	case logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel:
		components = append(components, aurora.Red(entry.Level.String()).String())
	case logrus.WarnLevel:
		components = append(components, aurora.Yellow(" warn").String())
	case logrus.InfoLevel:
		components = append(components, aurora.Cyan(" info").String())
	case logrus.DebugLevel:
		components = append(components, aurora.Green("debug").String())
	case logrus.TraceLevel:
		components = append(components, aurora.White("trace").String())
	}

	if entry.HasCaller() {
		components = append(components,
			aurora.Faint(
				fmt.Sprintf(
					"%s:%d:", filepath.Base(entry.Caller.File), entry.Caller.Line)).String())
	}

	components = append(components, entry.Message)
	components = append(components, "\n")

	return []byte(strings.Join(components[:], " ")), nil
}
