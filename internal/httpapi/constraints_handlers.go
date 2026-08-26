package httpapi

import (
	"net/http"

	"task247-reactcons/internal/model"
)

func (s *Server) addConstraint(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, model.ErrInvalidArgument)
		return
	}
	var closure model.MoietyClosure
	if err := decode(r, &closure); err != nil {
		writeErr(w, model.ErrInvalidArgument)
		return
	}
	created, err := s.svc.AddConstraint(r.Context(), id, closure)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (s *Server) listConstraints(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, model.ErrInvalidArgument)
		return
	}
	values, err := s.svc.ListConstraints(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, values)
}
