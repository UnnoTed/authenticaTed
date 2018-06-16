package users

import (
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/UnnoTed/authenticaTed/secret"
)

type Cfg struct {
	TokenExpirationTime time.Duration
	TokenSecret         []byte

	EncryptionLevel int
	EncryptionKey   string
}

var Config = &Cfg{
	TokenExpirationTime: 7 * 24 * time.Hour, // a week
	EncryptionLevel:     15,

	TokenSecret: secret.TokenSecret,

	EncryptionKey: secret.EncryptionKey,
}

var logger = log.New()

func init() {
	if isTest {
		Config.EncryptionLevel = 1
		logger.Level = log.DebugLevel
	}

	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})

	// Output to stderr instead of stdout, could also be a file.
	//log.SetOutput(os.Stderr)

	// Only log the warning severity or above.
	logger.Formatter = &log.TextFormatter{ForceColors: true}

	// use the line below to see the db log
	// export UPPERIO_DB_DEBUG=1
}
