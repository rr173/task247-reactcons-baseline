package httpapi

import (
	"net/http"

	"task247-reactcons/internal/model"
)

func (s *Server) markBoundary(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, model.ErrInvalidArgument)
		return
	}
	var value struct {
		SpeciesID int64  `json:"species_id"`
		Note      string `json:"note"`
	}
	if err := decode(r, &value); err != nil {
		writeErr(w, model.ErrInvalidArgument)
		return
	}
	created, err := s.svc.MarkBoundary(r.Context(), id, value.SpeciesID, value.Note)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (s *Server) listBoundaries(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, model.ErrInvalidArgument)
		return
	}
	values, err := s.svc.ListBoundaries(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, values)
}

func (s *Server) removeBoundary(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, model.ErrInvalidArgument)
		return
	}
	sid, err := pathSID(r)
	if err != nil {
		writeErr(w, model.ErrInvalidArgument)
		return
	}
	if err := s.svc.RemoveBoundary(r.Context(), id, sid); err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"network_id": id, "species_id": sid, "removed": true})
}
