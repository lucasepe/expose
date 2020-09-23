package slides

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/rakyll/statik/fs"

	// execute statik init() function
	_ "github.com/lucasepe/expose/statik"
)

// Serve create the Remark HTML slides and the HTTP server.
func Serve(filename string) error {
	sfs, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", Handler(filename))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(sfs)))

	// we want a free, available port selected by the system
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return err
	}

	port := listener.Addr().(*net.TCPAddr).Port

	origin := &url.URL{Scheme: "http"}
	origin.Host = net.JoinHostPort(getOutboundIP(), strconv.Itoa(port))

	fmt.Printf("Open your web browser and visit '%s'\n\n", origin)
	fmt.Printf("You can run Chrome in application mode:\n")
	fmt.Printf(" * Linux  : google-chrome --app=%s\n", origin)
	fmt.Printf(" * Windows: chrome --app=%s\n", origin)

	server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: mux}

	go func() error {
		if err := server.Serve(listener); err != nil {
			return err
		}
		return nil
	}()

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (pkill -2)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

// Get preferred outbound ip of this machine
func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

// Handler returns a handler that serves HTTP requests
// with the contents of the Markdown slides using Remark JS.
func Handler(filename string) http.Handler {
	return &slidesHandler{filename}
}

type slidesHandler struct {
	fileName string
}

func (h *slidesHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	doc, err := FromFile(h.fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := doc.Render(w); err != nil {
		msg := fmt.Sprintf("error: %s while rendering slide", err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
}
