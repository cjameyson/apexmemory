package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"apexmemory.ai/internal/app"
	"github.com/google/uuid"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		dsn       string
		email     string
		notebooks int
		facts     int
		clear     bool
	)

	flag.StringVar(&dsn, "dsn", "", "PostgreSQL DSN (or set DATABASE_URL env var)")
	flag.StringVar(&email, "email", "chrisjameyson@gmail.com", "User email to seed data for")
	flag.IntVar(&notebooks, "notebooks", 10, "Number of notebooks to create")
	flag.IntVar(&facts, "facts", 15, "Number of facts per notebook")
	flag.BoolVar(&clear, "clear", false, "Delete all existing notebooks for the user before seeding")
	flag.Parse()

	if dsn == "" {
		dsn = os.Getenv("DATABASE_URL")
	}
	if dsn == "" {
		dsn = os.Getenv("PG_APP_DSN")
	}
	if dsn == "" {
		return fmt.Errorf("missing database DSN: use -dsn flag or DATABASE_URL/PG_APP_DSN env var")
	}

	cfg := app.Config{Env: "development"}
	cfg.DB.DSN = dsn
	cfg.DB.MaxOpenConns = 5

	application := app.New(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if err := application.ConnectDB(ctx); err != nil {
		return fmt.Errorf("connect db: %w", err)
	}
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
		"notebooks", notebooks,
		"facts_per_notebook", facts,
	)

	if err := seedData(ctx, application, user.ID, notebooks, facts); err != nil {
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
		{"🧬", "Biology 202", "Cell biology, genetics, and evolution"},
		{"🇪🇸", "Spanish B2", "Intermediate Spanish vocabulary and grammar"},
		{"♾️", "Calculus", "Derivatives, integrals, and limits"},
		{"🇺🇸", "US History", "American history from colonial era to present"},
		{"🌍", "World History", "Global civilizations and major events"},
		{"🌡️", "Thermodynamics", "Heat transfer and energy systems"},
		{"🚀", "Quantum Mechanics", "Quantum mechanics and quantum field theory"},
		{"⚡", "Electrodynamics", "Electrodynamics and electromagnetism"},
		{"🧠", "Psychology", "Cognitive and behavioral psychology"},
		{"💻", "Computer Science", "Algorithms, data structures, and systems"},
		{"📐", "Linear Algebra", "Vectors, matrices, and transformations"},
		{"🧪", "Chemistry", "Organic and inorganic chemistry fundamentals"},
		{"📊", "Statistics", "Probability and statistical inference"},
		{"🏛️", "Philosophy", "Logic, ethics, and epistemology"},
		{"🎵", "Music Theory", "Harmony, rhythm, and composition"},
	}
}

type factTemplate struct {
	factType string
	content  func(index int) json.RawMessage
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
				return json.RawMessage(fmt.Sprintf(
					`{"version":1,"fields":[{"name":"front","type":"plain_text","value":"%s (#%d)"},{"name":"back","type":"plain_text","value":"%s"}]}`,
					front, idx+1, back,
				))
			},
		})
	}

	for _, s := range clozes {
		sentence := s
		templates = append(templates, factTemplate{
			factType: "cloze",
			content: func(_ int) json.RawMessage {
				return json.RawMessage(fmt.Sprintf(
					`{"version":1,"fields":[{"name":"text","type":"cloze_text","value":"%s"}]}`,
					sentence,
				))
			},
		})
	}

	return templates
}
