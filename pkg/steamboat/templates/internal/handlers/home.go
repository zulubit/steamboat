package handlers

import (
	"net/http"
	"<<!.ProjectName!>>/internal/views/pages"
)

func (h *Handlers) HomeHandler(w http.ResponseWriter, r *http.Request) {
	component := pages.Home()
	err := component.Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
