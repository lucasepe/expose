package slides

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rakyll/statik/fs"

	// execute statik init() function
	_ "github.com/lucasepe/expose/statik"
)

// Exposer defines the slideshow server.
type Exposer struct {
	filename string
	port     int
	url      string
	server   *http.Server
}

// URL returns the slideshow link.
func (ex *Exposer) URL() string { return ex.url }

// Expose create a slideshow server.
func Expose(filename string) (*Exposer, error) {
	workDir, err := filepath.Abs(filepath.Dir(filename))
	if err != nil {
		return nil, err
	}
	extDir := filepath.Join(workDir, "/assets")

	sfs, err := fs.New()
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.Handle("/", MarkdownHandler(filename))
	mux.Handle("/internal/", http.StripPrefix("/internal/", http.FileServer(sfs)))
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(extDir))))

	// we want a free, available port selected by the system
	port, err := FreeTCPPort()
	if err != nil {
		return nil, err
	}

	origin := &url.URL{Scheme: "http"}
	origin.Host = net.JoinHostPort(GetOutboundIP(), strconv.Itoa(port))

	res := &Exposer{filename: filename, port: port, url: origin.String()}
	res.server = &http.Server{Addr: fmt.Sprintf(":%d", res.port), Handler: mux}

	return res, nil
}

// Serve starts the HTTP server and serves the slideshow.
func (ex *Exposer) Serve(ctx context.Context) error {
	go func() error {
		if err := ex.server.ListenAndServe(); err != nil {
			return err
		}
		return nil
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := ex.server.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
