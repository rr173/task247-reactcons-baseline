package httpapi

import (
	"net/http"

	"task247-reactcons/internal/model"
)

func (s *Server) publishVersion(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, model.ErrInvalidArgument)
		return
	}
	version, err := s.svc.PublishVersion(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, version)
}

func (s *Server) listVersions(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, model.ErrInvalidArgument)
		return
	}
	values, err := s.svc.ListVersions(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, values)
}

func (s *Server) getVersion(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, model.ErrInvalidArgument)
		return
	}
	value, err := s.svc.GetVersion(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, value)
}

func (s *Server) verifyVersion(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, model.ErrInvalidArgument)
		return
	}
	value, err := s.svc.VerifyVersion(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, value)
}
