package tunnel

import (
	"context"
	"net/http"
	"os"

	"server/internal/log"

	"golang.ngrok.com/ngrok/v2"
)

// NgrokTunnel represents an ngrok tunnel service
type NgrokTunnel struct {
	endpoint ngrok.EndpointListener
	logger   log.Logger
}

// New creates a new NgrokTunnel instance
func New(logger log.Logger) *NgrokTunnel {
	return &NgrokTunnel{logger: logger}
}

// Start starts the ngrok tunnel if environment variables are set.
func (n *NgrokTunnel) Start(handler http.Handler) {
	domain := os.Getenv("NGROK_URL")
	auth := os.Getenv("NGROK_AUTHTOKEN")
	if domain == "" || auth == "" {
		n.logger.Info("NGROK_URL or NGROK_AUTHTOKEN not set, skipping ngrok tunnel")
		return
	}

	go func() {
		ctx := context.Background()

		ln, err := ngrok.Listen(ctx, ngrok.WithURL(domain))
		if err != nil {
			n.logger.Error("Failed to start ngrok endpoint: %v", err)
			return
		}
		n.endpoint = ln

		n.logger.Info("Ngrok endpoint online at %s", ln.URL())

		// Serve HTTP via the ngrok endpoint
		if err := http.Serve(ln, handler); err != nil && ctx.Err() == nil {
			n.logger.Error("Ngrok endpoint server error: %v", err)
		}
	}()
}

func (n *NgrokTunnel) Shutdown(ctx context.Context) {
	if n.endpoint != nil {
		n.logger.Info("Closing ngrok endpoint")
		_ = n.endpoint.CloseWithContext(ctx)
	}
}

func (n *NgrokTunnel) IsActive() bool {
	return n.endpoint != nil
}
