package httpapi

import (
	"net/http"

	"task247-reactcons/internal/model"
)

func (s *Server) addReaction(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, model.ErrInvalidArgument)
		return
	}
	var value struct {
		Equation   string `json:"equation"`
		Reversible bool   `json:"reversible"`
	}
	if err := decode(r, &value); err != nil {
		writeErr(w, model.ErrInvalidArgument)
		return
	}
	if value.Equation == "" {
		writeErr(w, model.ErrInvalidArgument)
		return
	}
	created, err := s.svc.AddReaction(r.Context(), id, value.Equation, value.Reversible)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (s *Server) listReactions(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, model.ErrInvalidArgument)
		return
	}
	values, err := s.svc.ListReactions(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, values)
}

func (s *Server) exemptReaction(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeErr(w, model.ErrInvalidArgument)
		return
	}
	var value struct {
		Reason string `json:"reason"`
	}
	_ = decode(r, &value)
	if err := s.svc.ExemptReaction(r.Context(), id, value.Reason); err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"id": id, "status": model.ReactionExempt})
}
