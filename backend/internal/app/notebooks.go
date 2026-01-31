package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"apexmemory.ai/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// FSRSSettings represents the FSRS algorithm configuration for a notebook.
// Go code is the source of truth for default values.
type FSRSSettings struct {
	DesiredRetention float64   `json:"desired_retention"`
	Version          string    `json:"version"`
	Params           []float64 `json:"params"`
	LearningSteps    []int     `json:"learning_steps"`
	RelearningSteps  []int     `json:"relearning_steps"`
	MaximumInterval  int       `json:"maximum_interval"`
	EnableFuzzing    bool      `json:"enable_fuzzing"`
}

// DefaultFSRSSettings returns the default FSRS v6 settings.
// This is the canonical source of truth for FSRS defaults.
func DefaultFSRSSettings() FSRSSettings {
	return FSRSSettings{
		DesiredRetention: 0.9,
		Version:          "6",
		Params: []float64{
			0.212, 1.2931, 2.3065, 8.2956, 6.4133, 0.8334, 3.0194, 0.001,
			1.8722, 0.1666, 0.796, 1.4835, 0.0614, 0.2629, 1.6483, 0.6014,
			1.8729, 0.5425, 0.0912, 0.0658, 0.1542,
		},
		LearningSteps:   []int{60, 600},
		RelearningSteps: []int{600},
		MaximumInterval: 36500,
		EnableFuzzing:   true,
	}
}

// NotebookResponse is the API response representation of a notebook.
type NotebookResponse struct {
	ID           uuid.UUID    `json:"id"`
	Name         string       `json:"name"`
	Description  *string      `json:"description"`
	Emoji        *string      `json:"emoji"`
	Color        *string      `json:"color"`
	FSRSSettings FSRSSettings `json:"fsrs_settings"`
	Position     int32        `json:"position"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

// toNotebookResponse converts a db.Notebook to a NotebookResponse.
func toNotebookResponse(n db.Notebook) NotebookResponse {
	resp := NotebookResponse{
		ID:        n.ID,
		Name:      n.Name,
		Position:  n.Position,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}

	if n.Description.Valid {
		resp.Description = &n.Description.String
	}
	if n.Emoji.Valid {
		resp.Emoji = &n.Emoji.String
	}
	if n.Color.Valid {
		resp.Color = &n.Color.String
	}

	var settings FSRSSettings
	if err := json.Unmarshal(n.FsrsSettings, &settings); err != nil {
		settings = DefaultFSRSSettings()
	}
	resp.FSRSSettings = settings

	return resp
}

// CreateNotebookParams holds the parameters for creating a notebook.
type CreateNotebookParams struct {
	Name        string
	Description *string
	Emoji       *string
	Color       *string
	Position    *int32
}

// CreateNotebook creates a new notebook for the given user.
func (app *Application) CreateNotebook(ctx context.Context, userID uuid.UUID, params CreateNotebookParams) (db.Notebook, error) {
	// Generate default FSRS settings (Go is source of truth)
	fsrsJSON, err := json.Marshal(DefaultFSRSSettings())
	if err != nil {
		return db.Notebook{}, err
	}

	dbParams := db.CreateNotebookParams{
		UserID:       userID,
		Name:         params.Name,
		FsrsSettings: fsrsJSON,
	}

	if params.Description != nil {
		dbParams.Description = pgtype.Text{String: *params.Description, Valid: true}
	}

	if params.Emoji != nil {
		dbParams.Emoji = pgtype.Text{String: *params.Emoji, Valid: true}
	}

	if params.Color != nil {
		dbParams.Color = pgtype.Text{String: *params.Color, Valid: true}
	}

	if params.Position != nil {
		dbParams.Position = *params.Position
	}

	return app.Queries.CreateNotebook(ctx, dbParams)
}

// GetNotebook retrieves a notebook by ID for the given user.
// Returns errNotebookNotFound if the notebook doesn't exist or belongs to another user.
func (app *Application) GetNotebook(ctx context.Context, userID, notebookID uuid.UUID) (db.Notebook, error) {
	notebook, err := app.Queries.GetNotebook(ctx, db.GetNotebookParams{
		UserID: userID,
		ID:     notebookID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Notebook{}, errNotebookNotFound
		}
		return db.Notebook{}, err
	}
	return notebook, nil
}

// ListNotebooks retrieves all non-archived notebooks for the given user.
func (app *Application) ListNotebooks(ctx context.Context, userID uuid.UUID) ([]db.Notebook, error) {
	return app.Queries.ListNotebooks(ctx, userID)
}

// UpdateNotebookParams holds the optional fields for updating a notebook.
type UpdateNotebookParams struct {
	Name             *string
	Description      OptionalString // Supports clearing via explicit null
	Emoji            OptionalString // Supports clearing via explicit null
	Color            OptionalString // Supports clearing via explicit null
	Position         *int32
	DesiredRetention *float64
}

// UpdateNotebook updates a notebook with the provided fields.
// All updates happen atomically in a single query.
func (app *Application) UpdateNotebook(ctx context.Context, userID, notebookID uuid.UUID, params UpdateNotebookParams) (db.Notebook, error) {
	dbParams := db.UpdateNotebookParams{
		UserID: userID,
		ID:     notebookID,
	}

	if params.Name != nil {
		dbParams.Name = pgtype.Text{String: *params.Name, Valid: true}
	}

	// Handle description: can be updated or cleared
	if params.Description.Set {
		if params.Description.Value == nil {
			dbParams.ClearDescription = true
		} else {
			dbParams.Description = pgtype.Text{String: *params.Description.Value, Valid: true}
		}
	}

	if params.Emoji.Set {
		if params.Emoji.Value == nil {
			dbParams.ClearEmoji = true
		} else {
			dbParams.Emoji = pgtype.Text{String: *params.Emoji.Value, Valid: true}
		}
	}

	if params.Color.Set {
		if params.Color.Value == nil {
			dbParams.ClearColor = true
		} else {
			dbParams.Color = pgtype.Text{String: *params.Color.Value, Valid: true}
		}
	}

	if params.Position != nil {
		dbParams.Position = pgtype.Int4{Int32: *params.Position, Valid: true}
	}

	if params.DesiredRetention != nil {
		dbParams.UpdateRetention = true
		dbParams.DesiredRetention = *params.DesiredRetention
	}

	notebook, err := app.Queries.UpdateNotebook(ctx, dbParams)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Notebook{}, errNotebookNotFound
		}
		return db.Notebook{}, err
	}

	return notebook, nil
}

// ArchiveNotebook archives a notebook (soft delete).
// This operation is idempotent: archiving an already-archived notebook succeeds.
// Returns errNotebookNotFound only if the notebook doesn't exist.
func (app *Application) ArchiveNotebook(ctx context.Context, userID, notebookID uuid.UUID) error {
	rowsAffected, err := app.Queries.ArchiveNotebook(ctx, db.ArchiveNotebookParams{
		UserID: userID,
		ID:     notebookID,
	})
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		// Check if notebook exists but is already archived (idempotent case)
		status, err := app.Queries.IsNotebookArchived(ctx, db.IsNotebookArchivedParams{
			UserID: userID,
			ID:     notebookID,
		})
		if err != nil {
			return err
		}
		if !status.Exists {
			return errNotebookNotFound
		}
		// Already archived - idempotent success
	}
	return nil
}

// errNotebookNotFound is a sentinel error for notebook not found conditions.
var errNotebookNotFound = errors.New("notebook not found")

// TODO(seed): Remove this function once auto-seed script is created.
// seedExampleNotebooks creates example notebooks and seeds facts for new users.
func (app *Application) seedExampleNotebooks(ctx context.Context, userID uuid.UUID) ([]db.Notebook, error) {
	examples := []struct {
		emoji       string
		name        string
		description string
		position    int32
	}{
		{"🧬", "Biology 202", "Cell biology, genetics, and evolution fundamentals", 0},
		{"🇪🇸", "Spanish B2", "Intermediate Spanish vocabulary and grammar", 1},
		{"♾️", "Calculus", "Derivatives, integrals, and limits", 2},
		{"🇺🇸", "US History", "American history from colonial era to present", 3},
		{"🌍", "World History", "World history from colonial era to present", 4},
		{"🇺🇸", "Civil War", "American Civil War history", 5},
		{"🌡️", "Thermodynamics", "Thermodynamics and heat transfer", 6},
		// create
		{"🚀", "Quantum Mechanics", "Quantum mechanics and quantum field theory", 7},
		{"⚡", "Electrodynamics", "Electrodynamics and electromagnetism", 8},
		{"⚓", "Nautical Archaeology", "Nautical archaeology and shipwreck history", 9},
	}

	notebooks := make([]db.Notebook, 0, len(examples))
	for _, ex := range examples {
		desc := ex.description
		emoji := ex.emoji
		pos := ex.position
		nb, err := app.CreateNotebook(ctx, userID, CreateNotebookParams{
			Name:        ex.name,
			Description: &desc,
			Emoji:       &emoji,
			Position:    &pos,
		})
		if err != nil {
			return nil, err
		}
		notebooks = append(notebooks, nb)
	}
	if err := app.seedExampleFacts(ctx, userID, notebooks); err != nil {
		return notebooks, err
	}

	return notebooks, nil
}

// seedExampleFacts populates Biology 202 and Quantum Mechanics notebooks with example facts.
func (app *Application) seedExampleFacts(ctx context.Context, userID uuid.UUID, notebooks []db.Notebook) error {
	// Find target notebooks by name
	var bioID, qmID uuid.UUID
	for _, nb := range notebooks {
		switch nb.Name {
		case "Biology 202":
			bioID = nb.ID
		case "Quantum Mechanics":
			qmID = nb.ID
		}
	}
	if bioID == uuid.Nil || qmID == uuid.Nil {
		return nil // notebooks not found, skip seeding
	}

	type seedFact struct {
		notebookID uuid.UUID
		factType   string
		content    string
	}

	facts := []seedFact{
		// Biology 202 - Basic (8)
		{bioID, "basic", `{"version":1,"fields":[{"name":"front","type":"plain_text","value":"What is the powerhouse of the cell?"},{"name":"back","type":"plain_text","value":"The mitochondria. It generates most of the cell's supply of ATP through oxidative phosphorylation."}]}`},
		{bioID, "basic", `{"version":1,"fields":[{"name":"front","type":"plain_text","value":"What is the difference between mitosis and meiosis?"},{"name":"back","type":"plain_text","value":"Mitosis produces two identical diploid daughter cells, while meiosis produces four genetically unique haploid gametes through two rounds of division."}]}`},
		{bioID, "basic", `{"version":1,"fields":[{"name":"front","type":"plain_text","value":"What is the central dogma of molecular biology?"},{"name":"back","type":"plain_text","value":"DNA is transcribed into RNA, which is translated into protein. Information flows DNA → RNA → Protein."}]}`},
		{bioID, "basic", `{"version":1,"fields":[{"name":"front","type":"plain_text","value":"What are the four nitrogenous bases in DNA?"},{"name":"back","type":"plain_text","value":"Adenine (A), Thymine (T), Guanine (G), and Cytosine (C). A pairs with T, G pairs with C."}]}`},
		{bioID, "basic", `{"version":1,"fields":[{"name":"front","type":"plain_text","value":"What is natural selection?"},{"name":"back","type":"plain_text","value":"The process where organisms with favorable traits are more likely to survive and reproduce, passing those traits to offspring. It is the primary mechanism of evolution."}]}`},
		{bioID, "basic", `{"version":1,"fields":[{"name":"front","type":"plain_text","value":"What is the function of ribosomes?"},{"name":"back","type":"plain_text","value":"Ribosomes synthesize proteins by translating mRNA sequences into polypeptide chains of amino acids."}]}`},
		{bioID, "basic", `{"version":1,"fields":[{"name":"front","type":"plain_text","value":"What is the difference between prokaryotic and eukaryotic cells?"},{"name":"back","type":"plain_text","value":"Prokaryotic cells lack a membrane-bound nucleus and organelles (bacteria, archaea). Eukaryotic cells have a nucleus and membrane-bound organelles (animals, plants, fungi)."}]}`},
		{bioID, "basic", `{"version":1,"fields":[{"name":"front","type":"plain_text","value":"What is the role of ATP in cellular metabolism?"},{"name":"back","type":"plain_text","value":"ATP (adenosine triphosphate) is the primary energy currency of the cell. It stores and transfers chemical energy for cellular processes like muscle contraction, active transport, and biosynthesis."}]}`},

		// Biology 202 - Cloze (5)
		{bioID, "cloze", `{"version":1,"fields":[{"name":"text","type":"cloze_text","cloze_text":"The process of {{c1::photosynthesis}} converts light energy into {{c2::chemical energy}} stored in glucose."}]}`},
		{bioID, "cloze", `{"version":1,"fields":[{"name":"text","type":"cloze_text","cloze_text":"DNA replication is {{c1::semi-conservative}}, meaning each new double helix contains one {{c2::original}} strand and one {{c3::newly synthesized}} strand."}]}`},
		{bioID, "cloze", `{"version":1,"fields":[{"name":"text","type":"cloze_text","cloze_text":"The {{c1::endoplasmic reticulum}} is responsible for protein folding and lipid synthesis within the cell."}]}`},
		{bioID, "cloze", `{"version":1,"fields":[{"name":"text","type":"cloze_text","cloze_text":"In genetics, a {{c1::phenotype}} is the observable characteristic, while a {{c2::genotype}} is the genetic makeup."}]}`},
		{bioID, "cloze", `{"version":1,"fields":[{"name":"text","type":"cloze_text","cloze_text":"The {{c1::Krebs cycle}} (citric acid cycle) takes place in the {{c2::mitochondrial matrix}} and produces {{c3::NADH}} and FADH2."}]}`},

		// Biology 202 - Image Occlusion (2)
		{bioID, "image_occlusion", `{"version":1,"fields":[{"name":"title","type":"plain_text","value":"Animal Cell Diagram"},{"name":"image","type":"image_occlusion","image":{"url":"https://placeholder.apexmemory.ai/cell-diagram.png","width":800,"height":600},"regions":[{"id":"m_region01nucleus","shape":{"type":"rect","x":300,"y":200,"width":200,"height":200}},{"id":"m_region02mito","shape":{"type":"rect","x":100,"y":350,"width":120,"height":80}}]}]}`},
		{bioID, "image_occlusion", `{"version":1,"fields":[{"name":"title","type":"plain_text","value":"DNA Double Helix Structure"},{"name":"image","type":"image_occlusion","image":{"url":"https://placeholder.apexmemory.ai/dna-structure.png","width":800,"height":600},"regions":[{"id":"m_region03helix","shape":{"type":"rect","x":200,"y":100,"width":400,"height":400}},{"id":"m_region04bases","shape":{"type":"rect","x":350,"y":250,"width":100,"height":100}}]}]}`},

		// Quantum Mechanics - Basic (5)
		{qmID, "basic", `{"version":1,"fields":[{"name":"front","type":"plain_text","value":"What is wave-particle duality?"},{"name":"back","type":"plain_text","value":"Wave-particle duality is the concept that quantum entities (photons, electrons) exhibit both wave-like and particle-like behavior depending on the experimental setup."}]}`},
		{qmID, "basic", `{"version":1,"fields":[{"name":"front","type":"plain_text","value":"State Heisenberg's uncertainty principle."},{"name":"back","type":"plain_text","value":"It is impossible to simultaneously know both the exact position and exact momentum of a particle. Mathematically: Δx·Δp ≥ ℏ/2, where ℏ is the reduced Planck constant."}]}`},
		{qmID, "basic", `{"version":1,"fields":[{"name":"front","type":"plain_text","value":"What is quantum superposition?"},{"name":"back","type":"plain_text","value":"A quantum system can exist in multiple states simultaneously until measured. Upon measurement, the system collapses into one definite state."}]}`},
		{qmID, "basic", `{"version":1,"fields":[{"name":"front","type":"plain_text","value":"What is quantum entanglement?"},{"name":"back","type":"plain_text","value":"When two particles become entangled, measuring the state of one instantly determines the state of the other, regardless of distance. Einstein called it 'spooky action at a distance.'"}]}`},
		{qmID, "basic", `{"version":1,"fields":[{"name":"front","type":"plain_text","value":"What is the Schrödinger equation used for?"},{"name":"back","type":"plain_text","value":"The Schrödinger equation describes how the quantum state (wave function) of a physical system changes over time. It is the fundamental equation of non-relativistic quantum mechanics."}]}`},

		// Quantum Mechanics - Cloze (3)
		{qmID, "cloze", `{"version":1,"fields":[{"name":"text","type":"cloze_text","cloze_text":"The energy of a photon is given by E = {{c1::hf}}, where h is Planck's constant and f is the {{c2::frequency}}."}]}`},
		{qmID, "cloze", `{"version":1,"fields":[{"name":"text","type":"cloze_text","cloze_text":"In the double-slit experiment, single particles create an {{c1::interference pattern}}, demonstrating {{c2::wave-particle duality}}."}]}`},
		{qmID, "cloze", `{"version":1,"fields":[{"name":"text","type":"cloze_text","cloze_text":"The {{c1::Pauli exclusion principle}} states that no two {{c2::fermions}} can occupy the same quantum state simultaneously."}]}`},

		// Quantum Mechanics - Image Occlusion (2)
		{qmID, "image_occlusion", `{"version":1,"fields":[{"name":"title","type":"plain_text","value":"Hydrogen Energy Levels"},{"name":"image","type":"image_occlusion","image":{"url":"https://placeholder.apexmemory.ai/energy-levels.png","width":800,"height":600},"regions":[{"id":"m_region05ground","shape":{"type":"rect","x":100,"y":450,"width":600,"height":50}},{"id":"m_region06excited","shape":{"type":"rect","x":100,"y":200,"width":600,"height":50}}]}]}`},
		{qmID, "image_occlusion", `{"version":1,"fields":[{"name":"title","type":"plain_text","value":"Atomic Orbital Shapes"},{"name":"image","type":"image_occlusion","image":{"url":"https://placeholder.apexmemory.ai/orbital-shapes.png","width":800,"height":600},"regions":[{"id":"m_region07sorbital","shape":{"type":"rect","x":50,"y":200,"width":150,"height":150}},{"id":"m_region08porbital","shape":{"type":"rect","x":300,"y":200,"width":200,"height":150}}]}]}`},
	}

	for _, f := range facts {
		_, _, err := app.CreateFact(ctx, userID, f.notebookID, f.factType, json.RawMessage(f.content))
		if err != nil {
			return fmt.Errorf("seed fact (%s/%s): %w", f.factType, f.notebookID, err)
		}
	}

	return nil
}
