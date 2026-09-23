package sql

import "testing"

// TestFormatColumnExtraExplicit preserves metadata from existing schema providers.
func TestFormatColumnExtraExplicit(t *testing.T) {
	column := &Column{Extra: "provider metadata", AutoIncrement: true}
	if got := FormatColumnExtra(column.Copy()); got != column.Extra {
		t.Fatalf("EXTRA = %q, want preserved metadata %q", got, column.Extra)
	}

	// Columns without an explicit value use the current engine's derived metadata.
	column.Extra = ""
	if got := FormatColumnExtra(column); got != "auto_increment" {
		t.Fatalf("EXTRA = %q, want auto_increment", got)
	}
}
