package logruskit_test

import (
	"strings"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/typical-go/typical-rest-server/pkg/logruskit"
)

func TestEchoLogrus(t *testing.T) {
	patch := monkey.Patch(time.Now, func() time.Time {
		t, _ := time.Parse(time.RFC3339, "2014-11-12T11:45:26Z")
		return t
	})
	defer patch.Unpatch()

	var out strings.Builder
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(&out)

	log := logruskit.EchoLogger(logger)

	log.Printj(map[string]interface{}{"Printj": "1"})
	log.Debugj(map[string]interface{}{"Debugj": "1"})
	log.Infoj(map[string]interface{}{"Infoj": "1"})
	log.Warnj(map[string]interface{}{"Warnj": "1"})
	log.Errorj(map[string]interface{}{"Errorj": "1"})

	// NOTE: cause breaking test
	// log.Fatalj(map[string]interface{}{"Fatalj": "1"})
	// log.Panicj(map[string]interface{}{"Panicj": "1"})

	require.Equal(t, `{"Printj":"1","level":"info","msg":"","time":"2014-11-12T11:45:26Z"}
{"Infoj":"1","level":"info","msg":"","time":"2014-11-12T11:45:26Z"}
{"Warnj":"1","level":"warning","msg":"","time":"2014-11-12T11:45:26Z"}
{"Errorj":"1","level":"error","msg":"","time":"2014-11-12T11:45:26Z"}
`, out.String())

	log.SetPrefix("prefix01")
	require.Equal(t, "prefix01", log.Prefix())
	require.Equal(t, &out, log.Output())

	// NOTE: nothing to verify
	log.SetHeader("")

}

func TestEchoLogrus_SetLevel(t *testing.T) {
	testcases := []struct {
		TestName string
		lvl      log.Lvl
		expected logrus.Level
	}{
		{lvl: log.INFO, expected: logrus.InfoLevel},
		{lvl: log.WARN, expected: logrus.WarnLevel},
		{lvl: log.ERROR, expected: logrus.ErrorLevel},
		{lvl: log.DEBUG, expected: logrus.DebugLevel},
		{lvl: 11, expected: logrus.WarnLevel},
	}
	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			logger := logrus.New()
			log := logruskit.EchoLogger(logger)
			log.SetLevel(tt.lvl)
			require.Equal(t, tt.expected, logger.GetLevel())
		})
	}
}

func TestEchoLogrus_Level(t *testing.T) {
	testcases := []struct {
		TestName string
		lvl      logrus.Level
		expected log.Lvl
	}{
		{lvl: logrus.DebugLevel, expected: log.DEBUG},
		{lvl: logrus.WarnLevel, expected: log.WARN},
		{lvl: logrus.ErrorLevel, expected: log.ERROR},
		{lvl: logrus.InfoLevel, expected: log.INFO},
		{lvl: 22, expected: log.WARN},
	}
	for _, tt := range testcases {
		t.Run(tt.TestName, func(t *testing.T) {
			logger := logrus.New()
			logger.SetLevel(tt.lvl)
			log := logruskit.EchoLogger(logger)
			require.Equal(t, tt.expected, log.Level())
		})
	}
}
