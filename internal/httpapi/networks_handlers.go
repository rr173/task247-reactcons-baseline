package httpapi

import (
	"net/http"

	"task247-reactcons/internal/model"
)

func (s *Server) selfCheck(w http.ResponseWriter, r *http.Request) {
	value, err := s.svc.SelfCheck(r.Context())
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, value)
}

func (s *Server) createNetwork(w http.ResponseWriter, r *http.Request) {
	var value model.Network
	if err := decode(r, &value); err != nil {
		writeErr(w, model.ErrInvalidArgument)
		return
	}
	created, err := s.svc.CreateNetwork(r.Context(), value)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (s *Server) listNetworks(w http.ResponseWriter, r *http.Request) {
	values, err := s.svc.ListNetworks(r.Context())
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, values)
}

func (s *Server) getNetwork(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, model.ErrInvalidArgument)
		return
	}
	value, err := s.svc.GetNetwork(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, value)
}
