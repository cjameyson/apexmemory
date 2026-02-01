package app

import (
	"net/http/httptest"
	"testing"
)

func TestParseSort(t *testing.T) {
	tests := []struct {
		name      string
		query     string
		allowed   []string
		wantField string
		wantAsc   bool
		wantErr   bool
	}{
		{
			name:      "empty returns zero value",
			query:     "",
			allowed:   []string{"created", "updated"},
			wantField: "",
			wantAsc:   false,
			wantErr:   false,
		},
		{
			name:      "ascending field",
			query:     "?sort=created",
			allowed:   []string{"created", "updated"},
			wantField: "created",
			wantAsc:   true,
			wantErr:   false,
		},
		{
			name:      "descending with dash prefix",
			query:     "?sort=-updated",
			allowed:   []string{"created", "updated"},
			wantField: "updated",
			wantAsc:   false,
			wantErr:   false,
		},
		{
			name:    "invalid field returns error",
			query:   "?sort=name",
			allowed: []string{"created", "updated"},
			wantErr: true,
		},
		{
			name:    "invalid field with dash prefix returns error",
			query:   "?sort=-name",
			allowed: []string{"created", "updated"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/"+tt.query, nil)
			got, err := parseSort(r, tt.allowed...)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Field != tt.wantField {
				t.Errorf("Field = %q, want %q", got.Field, tt.wantField)
			}
			if got.Asc != tt.wantAsc {
				t.Errorf("Asc = %v, want %v", got.Asc, tt.wantAsc)
			}
		})
	}
}
