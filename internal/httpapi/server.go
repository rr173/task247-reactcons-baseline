package httpapi

import (
	"net/http"
	"strconv"

	"task247-reactcons/internal/service"
)

type Server struct {
	svc *service.Service
	mux *http.ServeMux
}

func New(svc *service.Service) *Server {
	server := &Server{svc: svc, mux: http.NewServeMux()}
	server.routes()
	return server
}

func (s *Server) Handler() http.Handler { return s.mux }

func (s *Server) routes() {
	s.mux.HandleFunc("GET /api/self-check", s.selfCheck)
	s.mux.HandleFunc("POST /api/networks", s.createNetwork)
	s.mux.HandleFunc("GET /api/networks", s.listNetworks)
	s.mux.HandleFunc("GET /api/networks/{id}", s.getNetwork)

	s.mux.HandleFunc("POST /api/networks/{id}/species", s.addSpecies)
	s.mux.HandleFunc("GET /api/networks/{id}/species", s.listSpecies)
	s.mux.HandleFunc("PUT /api/species/{id}", s.reviseSpecies)

	s.mux.HandleFunc("POST /api/networks/{id}/reactions", s.addReaction)
	s.mux.HandleFunc("GET /api/networks/{id}/reactions", s.listReactions)
	s.mux.HandleFunc("POST /api/reactions/{id}/exempt", s.exemptReaction)

	s.mux.HandleFunc("POST /api/networks/{id}/solve", s.solve)
	s.mux.HandleFunc("GET /api/networks/{id}/conservation", s.conservation)
	s.mux.HandleFunc("GET /api/networks/{id}/conserved-pools", s.conservedPools)
	s.mux.HandleFunc("GET /api/networks/{id}/conflicts", s.conflictSets)

	s.mux.HandleFunc("POST /api/networks/{id}/boundaries", s.markBoundary)
	s.mux.HandleFunc("GET /api/networks/{id}/boundaries", s.listBoundaries)
	s.mux.HandleFunc("DELETE /api/networks/{id}/boundaries/{sid}", s.removeBoundary)

	s.mux.HandleFunc("POST /api/networks/{id}/constraints", s.addConstraint)
	s.mux.HandleFunc("GET /api/networks/{id}/constraints", s.listConstraints)

	s.mux.HandleFunc("POST /api/networks/{id}/versions", s.publishVersion)
	s.mux.HandleFunc("GET /api/networks/{id}/versions", s.listVersions)
	s.mux.HandleFunc("GET /api/versions/{id}", s.getVersion)
	s.mux.HandleFunc("GET /api/versions/{id}/verify", s.verifyVersion)
}

func pathID(r *http.Request) (int64, error) { return strconv.ParseInt(r.PathValue("id"), 10, 64) }
func pathSID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("sid"), 10, 64)
}
