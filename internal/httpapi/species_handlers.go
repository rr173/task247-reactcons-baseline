package httpapi

import (
	"net/http"

	"task247-reactcons/internal/model"
)

func (s *Server) addSpecies(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, model.ErrInvalidArgument)
		return
	}
	var value struct {
		Symbol string `json:"symbol"`
		Name   string `json:"name"`
	}
	if err := decode(r, &value); err != nil {
		writeErr(w, model.ErrInvalidArgument)
		return
	}
	if value.Symbol == "" {
		writeErr(w, model.ErrInvalidArgument)
		return
	}
	created, err := s.svc.AddSpecies(r.Context(), id, value.Symbol, value.Name)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (s *Server) listSpecies(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, model.ErrInvalidArgument)
		return
	}
	values, err := s.svc.ListSpecies(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, values)
}

func (s *Server) reviseSpecies(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, model.ErrInvalidArgument)
		return
	}
	var value struct {
		Status string `json:"status"`
		Note   string `json:"note"`
	}
	if err := decode(r, &value); err != nil {
		writeErr(w, model.ErrInvalidArgument)
		return
	}
	if err := s.svc.ReviseSpecies(r.Context(), id, model.SpeciesStatus(value.Status), value.Note); err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"id": id, "status": value.Status})
}
