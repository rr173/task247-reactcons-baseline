package model

import "encoding/json"

// Network lifecycle.
type NetworkStatus string

const (
	NetworkDraft       NetworkStatus = "draft"
	NetworkSolving     NetworkStatus = "solving"
	NetworkConflict    NetworkStatus = "conflict"
	NetworkPublishable NetworkStatus = "publishable"
	NetworkSealed      NetworkStatus = "sealed"
)

// Species lifecycle.
type SpeciesStatus string

const (
	SpeciesPending       SpeciesStatus = "pending"
	SpeciesValid         SpeciesStatus = "valid"
	SpeciesMissingComp   SpeciesStatus = "missing_composition"
	SpeciesConflict      SpeciesStatus = "conflict"
)

// Reaction lifecycle.
type ReactionStatus string

const (
	ReactionCandidate ReactionStatus = "candidate"
	ReactionConserved ReactionStatus = "conserved"
	ReactionViolates  ReactionStatus = "violates"
	ReactionExempt    ReactionStatus = "exempt"
)

// Version lifecycle.
type VersionStatus string

const (
	VersionDraft      VersionStatus = "draft"
	VersionPublished  VersionStatus = "published"
	VersionSuperseded VersionStatus = "superseded"
)

type Network struct {
	ID          int64         `json:"id"`
	ExtKey      string        `json:"ext_key"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Status      NetworkStatus `json:"status"`
	CreatedAt   string        `json:"created_at"`
	UpdatedAt   string        `json:"updated_at"`
}

type Species struct {
	ID          int64         `json:"id"`
	NetworkID   int64         `json:"network_id"`
	Symbol      string        `json:"symbol"`
	Name        string        `json:"name"`
	Charge      int           `json:"charge"`
	Composition map[string]int `json:"composition"`
	Status      SpeciesStatus `json:"status"`
	Note        string        `json:"note"`
	CreatedAt   string        `json:"created_at"`
}

type Participant struct {
	SpeciesID int64   `json:"species_id"`
	Symbol    string  `json:"symbol"`
	Role      string  `json:"role"` // "reactant" | "product"
	Coeff     float64 `json:"coeff"`
}

type Reaction struct {
	ID           int64        `json:"id"`
	NetworkID    int64        `json:"network_id"`
	Equation     string       `json:"equation"`
	Reversible   bool         `json:"reversible"`
	Status       ReactionStatus `json:"status"`
	Note         string       `json:"note"`
	Participants []Participant `json:"participants,omitempty"`
	CreatedAt    string       `json:"created_at"`
}

type Violation struct {
	Kind   string  `json:"kind"`  // "element" | "charge"
	Target string  `json:"target"` // element symbol or "charge"
	Net    float64 `json:"net"`
}

type ReactionConservation struct {
	ReactionID int64       `json:"reaction_id"`
	Equation   string      `json:"equation"`
	Conserved  bool        `json:"conserved"`
	Violations []Violation `json:"violations,omitempty"`
}

type PoolMember struct {
	SpeciesID int64   `json:"species_id"`
	Symbol    string  `json:"symbol"`
	Coeff     string  `json:"coeff"` // rational string
}

type ConservedPool struct {
	ID      int64        `json:"id"`
	Label   string       `json:"label"`
	Members []PoolMember `json:"members"`
}

type ConflictSet struct {
	ID         int64   `json:"id"`
	Kind       string  `json:"kind"` // "element" | "charge" | "constraint"
	Target     string  `json:"target"`
	ReactionIDs []int64 `json:"reaction_ids"`
}

type Boundary struct {
	ID        int64  `json:"id"`
	NetworkID int64  `json:"network_id"`
	SpeciesID int64  `json:"species_id"`
	Symbol    string `json:"symbol"`
	Note      string `json:"note"`
	CreatedAt string `json:"created_at"`
}

type MoietyMember struct {
	Symbol string  `json:"symbol"`
	Coeff  float64 `json:"coeff"`
}

type MoietyClosure struct {
	Members []MoietyMember `json:"members"`
}

type Constraint struct {
	ID          int64           `json:"id"`
	NetworkID   int64           `json:"network_id"`
	Kind        string          `json:"kind"` // "moietyclosure"
	Status      string          `json:"status"`
	Description string          `json:"description"`
	Payload     json.RawMessage `json:"payload"`
	CreatedAt   string          `json:"created_at"`
}

type NetworkVersion struct {
	ID          int64           `json:"id"`
	NetworkID   int64           `json:"network_id"`
	Status      VersionStatus   `json:"status"`
	ContentHash string          `json:"content_hash"`
	Payload     json.RawMessage `json:"payload"`
	CreatedAt   string          `json:"created_at"`
	PublishedAt string          `json:"published_at,omitempty"`
}

type SolveResult struct {
	NetworkID          int64                 `json:"network_id"`
	Status             NetworkStatus         `json:"status"`
	ReactionChecks     []ReactionConservation `json:"reaction_checks"`
	ConservedPools     []ConservedPool       `json:"conserved_pools"`
	ConflictSets       []ConflictSet         `json:"conflict_sets"`
	ConstraintChecks   []Constraint          `json:"constraint_checks"`
	ConflictReactionIDs []int64              `json:"conflict_reaction_ids"`
}

// Reference types passed to the matrix/diagnose layer (decoupled from storage rows).
type SpeciesRef struct {
	ID          int64
	Symbol      string
	Composition map[string]int
	Charge      int
	Open        bool
}

type ParticipantRef struct {
	SpeciesID int64
	Role      string
	Coeff     float64
}

type ReactionRef struct {
	ID           int64
	Equation     string
	Participants []ParticipantRef
}
