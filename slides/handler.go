package slides

import (
	"fmt"
	"net/http"
)

// MarkdownHandler serves HTTP requests with the
// content of the Markdown file rendered as HTML slides.
func MarkdownHandler(filename string) http.Handler {
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
