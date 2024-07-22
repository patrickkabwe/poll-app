package utils

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/a-h/templ"
)

func Render(templ templ.Component, ctx context.Context, w http.ResponseWriter) error {
	return templ.Render(ctx, w)
}

func Redirect(w http.ResponseWriter, r *http.Request, redirectTo string) {
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", redirectTo)
		w.WriteHeader(http.StatusSeeOther)
	} else {
		http.Redirect(w, r, redirectTo, http.StatusSeeOther)
	}
}

type ToastType string

const (
	INFO    ToastType = "info"
	SUCCESS ToastType = "success"
	WARNING ToastType = "warning"
	ERROR   ToastType = "error"
)

type Toast struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

func ShowToast(w http.ResponseWriter, ttype ToastType, message string) {
	toastRes := map[string]Toast{}
	toastRes["showToast"] = Toast{Message: message, Type: string(ttype)}
	jsonData, err := json.Marshal(toastRes)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("HX-Trigger", string(jsonData))
}
