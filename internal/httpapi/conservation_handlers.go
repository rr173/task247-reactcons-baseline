package httpapi

import (
	"net/http"

	"task247-reactcons/internal/model"
)

func (s *Server) solve(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, model.ErrInvalidArgument)
		return
	}
	result, err := s.svc.Solve(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) conservation(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, model.ErrInvalidArgument)
		return
	}
	result, err := s.svc.Conservation(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) conservedPools(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, model.ErrInvalidArgument)
		return
	}
	pools, err := s.svc.ConservedPools(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, pools)
}

func (s *Server) conflictSets(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, model.ErrInvalidArgument)
		return
	}
	sets, err := s.svc.ConflictSets(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, sets)
}
