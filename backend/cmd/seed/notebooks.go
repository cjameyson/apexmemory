package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"apexmemory.ai/internal/app"
	"github.com/google/uuid"
)

func runNotebooksCmd(args []string) {
	fs := flag.NewFlagSet("notebooks", flag.ExitOnError)
	var (
		dsn       string
		email     string
		notebooks int
		facts     int
		clear     bool
	)
	fs.StringVar(&dsn, "dsn", "", "PostgreSQL DSN (or set DATABASE_URL env var)")
	fs.StringVar(&email, "email", "", "User email to seed data for (required)")
	fs.IntVar(&notebooks, "notebooks", 2, "Number of notebooks to create")
	fs.IntVar(&facts, "facts", 15, "Number of facts per notebook")
	fs.BoolVar(&clear, "clear", false, "Delete all existing notebooks for the user before seeding")
	fs.Usage = func() {
		fmt.Println(`seed notebooks - Seed notebooks and facts for a user

Usage:
  seed notebooks -email <email> [flags]

Flags:`)
		fs.PrintDefaults()
		fmt.Println(`
Examples:
  seed notebooks -email test-agent@apexmemory.ai
  seed notebooks -email user@example.com -notebooks 5 -facts 20
  seed notebooks -email user@example.com -clear`)
	}
	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	if email == "" {
		fmt.Fprintln(os.Stderr, "error: -email flag is required")
		fs.Usage()
		os.Exit(1)
	}

	if err := seedNotebooks(dsn, email, notebooks, facts, clear); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func seedNotebooks(flagDSN, email string, numNotebooks, factsPerNotebook int, clear bool) error {
	dsn, err := resolveDSN(flagDSN)
	if err != nil {
		return err
	}

	application, ctx, cancel, err := connectApp(dsn)
	if err != nil {
		return err
	}
	defer cancel()
	defer application.CloseDB()

	user, err := application.Queries.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("user not found for email %q: %w", email, err)
	}

	if clear {
		slog.Info("clearing existing notebooks", "user_id", user.ID, "email", email)
		_, err := application.DB.Exec(ctx, "DELETE FROM app.notebooks WHERE user_id = $1", user.ID)
		if err != nil {
			return fmt.Errorf("clear notebooks: %w", err)
		}
	}

	slog.Info("seeding data",
		"user_id", user.ID,
		"email", email,
		"notebooks", numNotebooks,
		"facts_per_notebook", factsPerNotebook,
	)

	if err := seedData(ctx, application, user.ID, numNotebooks, factsPerNotebook); err != nil {
		return fmt.Errorf("seed: %w", err)
	}

	slog.Info("seeding complete")
	return nil
}

func seedData(ctx context.Context, application *app.Application, userID uuid.UUID, numNotebooks, factsPerNotebook int) error {
	nbPool := notebookPool()
	fPool := factTemplatePool()

	for i := range numNotebooks {
		nb := nbPool[i%len(nbPool)]
		desc := nb.description
		emoji := nb.emoji
		pos := int32(i)

		name := nb.name
		if numNotebooks > len(nbPool) {
			name = fmt.Sprintf("%s %d", nb.name, i+1)
		}

		notebook, err := application.CreateNotebook(ctx, userID, app.CreateNotebookParams{
			Name:        name,
			Description: &desc,
			Emoji:       &emoji,
			Position:    &pos,
		})
		if err != nil {
			return fmt.Errorf("create notebook %d: %w", i, err)
		}

		for j := range factsPerNotebook {
			tmpl := fPool[j%len(fPool)]
			content := tmpl.content(j)
			_, _, err := application.CreateFact(ctx, userID, notebook.ID, tmpl.factType, content)
			if err != nil {
				return fmt.Errorf("create fact %d for notebook %d: %w", j, i, err)
			}
		}

		slog.Info("seeded notebook", "notebook", notebook.Name, "facts", factsPerNotebook)
	}

	return nil
}

type notebookTemplate struct {
	emoji       string
	name        string
	description string
}

func notebookPool() []notebookTemplate {
	return []notebookTemplate{
		{"", "Biology 202", "Cell biology, genetics, and evolution"},
		{"", "Spanish B2", "Intermediate Spanish vocabulary and grammar"},
		{"", "Calculus", "Derivatives, integrals, and limits"},
		{"", "US History", "American history from colonial era to present"},
		{"", "World History", "Global civilizations and major events"},
		{"", "Thermodynamics", "Heat transfer and energy systems"},
		{"", "Quantum Mechanics", "Quantum mechanics and quantum field theory"},
		{"", "Electrodynamics", "Electrodynamics and electromagnetism"},
		{"", "Psychology", "Cognitive and behavioral psychology"},
		{"", "Computer Science", "Algorithms, data structures, and systems"},
		{"", "Linear Algebra", "Vectors, matrices, and transformations"},
		{"", "Chemistry", "Organic and inorganic chemistry fundamentals"},
		{"", "Statistics", "Probability and statistical inference"},
		{"", "Philosophy", "Logic, ethics, and epistemology"},
		{"", "Music Theory", "Harmony, rhythm, and composition"},
	}
}

type factTemplate struct {
	factType string
	content  func(index int) json.RawMessage
}

type factContent struct {
	Version int         `json:"version"`
	Fields  []factField `json:"fields"`
}

type factField struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

func mustMarshal(v any) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("marshal fact content: %v", err))
	}
	return data
}

func factTemplatePool() []factTemplate {
	basics := []struct{ front, back string }{
		{"What is the powerhouse of the cell?", "The mitochondria. It generates most of the cell's supply of ATP through oxidative phosphorylation."},
		{"What is the difference between mitosis and meiosis?", "Mitosis produces two identical diploid daughter cells, while meiosis produces four genetically unique haploid gametes."},
		{"What is the central dogma of molecular biology?", "DNA is transcribed into RNA, which is translated into protein. Information flows DNA -> RNA -> Protein."},
		{"What are the four nitrogenous bases in DNA?", "Adenine (A), Thymine (T), Guanine (G), and Cytosine (C). A pairs with T, G pairs with C."},
		{"What is natural selection?", "Organisms with favorable traits are more likely to survive and reproduce, passing those traits to offspring."},
		{"What is the function of ribosomes?", "Ribosomes synthesize proteins by translating mRNA sequences into polypeptide chains of amino acids."},
		{"Define the key terminology for this field.", "The key terms establish a precise vocabulary for describing phenomena, measurements, and relationships."},
		{"How does this relate to real-world applications?", "Real-world applications include engineering, medicine, technology, and scientific research."},
		{"What is the historical significance of this concept?", "This concept transformed our understanding of the field and led to major advances in theory and practice."},
		{"What are the common misconceptions?", "A frequent misconception is oversimplifying the relationship between variables and ignoring boundary conditions."},
	}

	clozes := []string{
		`The process of {{c1::photosynthesis}} converts light energy into {{c2::chemical energy}} stored in glucose.`,
		`DNA replication is {{c1::semi-conservative}}, meaning each new double helix contains one {{c2::original}} strand and one {{c3::newly synthesized}} strand.`,
		`The {{c1::endoplasmic reticulum}} is responsible for protein folding and lipid synthesis within the cell.`,
		`In genetics, a {{c1::phenotype}} is the observable characteristic, while a {{c2::genotype}} is the genetic makeup.`,
		`The {{c1::Krebs cycle}} takes place in the {{c2::mitochondrial matrix}} and produces {{c3::NADH}} and FADH2.`,
	}

	templates := make([]factTemplate, 0, len(basics)+len(clozes))

	for _, b := range basics {
		front, back := b.front, b.back
		templates = append(templates, factTemplate{
			factType: "basic",
			content: func(idx int) json.RawMessage {
				return mustMarshal(factContent{
					Version: 1,
					Fields: []factField{
						{Name: "front", Type: "plain_text", Value: fmt.Sprintf("%s (#%d)", front, idx+1)},
						{Name: "back", Type: "plain_text", Value: back},
					},
				})
			},
		})
	}

	for _, s := range clozes {
		sentence := s
		templates = append(templates, factTemplate{
			factType: "cloze",
			content: func(_ int) json.RawMessage {
				return mustMarshal(factContent{
					Version: 1,
					Fields: []factField{
						{Name: "text", Type: "cloze_text", Value: sentence},
					},
				})
			},
		})
	}

	// Image occlusion templates - region IDs must match ^m_[a-zA-Z0-9_-]{6,24}$
	imageOcclusions := []imageOcclusionTemplate{
		{
			title:  "Cell Diagram",
			url:    "https://upload.wikimedia.org/wikipedia/commons/thumb/1/1a/Nucleus_ER.png/800px-Nucleus_ER.png",
			width:  800,
			height: 600,
			regions: []regionDef{
				{id: "m_nucleus_01", x: 200, y: 150, w: 120, h: 120},
				{id: "m_er_rough_2", x: 350, y: 200, w: 150, h: 80},
				{id: "m_er_smooth3", x: 520, y: 180, w: 140, h: 100},
			},
		},
		{
			title:  "Periodic Table Section",
			url:    "https://upload.wikimedia.org/wikipedia/commons/thumb/8/89/Colour_18-col_PT.svg/800px-Colour_18-col_PT.svg.png",
			width:  800,
			height: 500,
			regions: []regionDef{
				{id: "m_hydrogen01", x: 20, y: 30, w: 40, h: 40},
				{id: "m_helium_02", x: 740, y: 30, w: 40, h: 40},
			},
		},
		{
			title:  "Human Heart Anatomy",
			url:    "https://upload.wikimedia.org/wikipedia/commons/thumb/e/e5/Diagram_of_the_human_heart.svg/600px-Diagram_of_the_human_heart.svg.png",
			width:  600,
			height: 600,
			regions: []regionDef{
				{id: "m_leftvent_1", x: 320, y: 350, w: 100, h: 120},
				{id: "m_rightven_2", x: 180, y: 350, w: 100, h: 120},
				{id: "m_aorta_003", x: 260, y: 100, w: 80, h: 60},
				{id: "m_pulmon_04", x: 200, y: 120, w: 60, h: 50},
			},
		},
	}

	for _, img := range imageOcclusions {
		imgCopy := img
		templates = append(templates, factTemplate{
			factType: "image_occlusion",
			content: func(_ int) json.RawMessage {
				return imgCopy.toContent()
			},
		})
	}

	return templates
}

// imageOcclusionTemplate defines seed data for an image occlusion fact.
type imageOcclusionTemplate struct {
	title   string
	url     string
	width   int
	height  int
	regions []regionDef
}

type regionDef struct {
	id   string
	x, y int
	w, h int
}

func (t imageOcclusionTemplate) toContent() json.RawMessage {
	regions := make([]map[string]any, len(t.regions))
	for i, r := range t.regions {
		regions[i] = map[string]any{
			"id": r.id,
			"shape": map[string]any{
				"type":   "rect",
				"x":      r.x,
				"y":      r.y,
				"width":  r.w,
				"height": r.h,
			},
		}
	}

	content := map[string]any{
		"version": 1,
		"fields": []map[string]any{
			{"name": "title", "type": "plain_text", "value": t.title},
			{
				"name": "image",
				"type": "image_occlusion",
				"image": map[string]any{
					"url":    t.url,
					"width":  t.width,
					"height": t.height,
				},
				"regions": regions,
			},
		},
	}
	return mustMarshal(content)
}
