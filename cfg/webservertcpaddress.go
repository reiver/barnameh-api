package cfg

import (
	"fmt"

	"github.com/reiver/barnameh-api/env"
)

func WebServerTCPAddress() string {
	return fmt.Sprintf(":%s", env.TcpPort)
}
