package crypto

// https://exchange-docs.crypto.com/exchange/v1/rest-ws/index.html#introduction
import (
	"tri-arb/internal/platforms"
)

func Crypto(sandbox bool) platforms.Platform {
	var host platforms.Host
	if !sandbox {
		host = platforms.Host{
			User:       "stream.crypto.com",
			UserPath:   "/exchange/v1/user",
			Market:     "stream.crypto.com",
			MarketPath: "/exchange/v1/market",
			Scheme:     "wss",
			// TLS:        false,
		}
	} else {
		host = platforms.Host{
			User:       "uat-stream.3ona.co",
			UserPath:   "/exchange/v1/user",
			Market:     "uat-stream.3ona.co",
			MarketPath: "/exchange/v1/market",
			Scheme:     "wss",
			// TLS:        false,
		}
	}

	return platforms.Platform{
		Host: host,
	}

}
