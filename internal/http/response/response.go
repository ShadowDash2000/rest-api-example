package response

import (
	"github.com/go-chi/render"
	"net/http"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

const (
	statusError   = "error"
	statusSuccess = "success"
)

func RenderError(w http.ResponseWriter, r *http.Request, status int, message string) {
	render.Status(r, status)
	render.JSON(w, r, &Response{
		Status:  statusError,
		Message: message,
	})
}

func RenderSuccess(w http.ResponseWriter, r *http.Request, status int, message string) {
	render.Status(r, status)
	render.JSON(w, r, &Response{
		Status:  statusSuccess,
		Message: message,
	})
}
