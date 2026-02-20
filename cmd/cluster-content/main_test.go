package main

import "testing"

func TestIsProductionLikeURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{name: "localhost", url: "postgres://u:p@localhost:5432/db", want: false},
		{name: "loopback", url: "postgres://u:p@127.0.0.1:5432/db", want: false},
		{name: "local suffix", url: "postgres://u:p@db.local:5432/db", want: false},
		{name: "remote host", url: "postgres://u:p@prod-db.example.com:5432/db", want: true},
		{name: "invalid url", url: "not-a-url", want: true},
	}
	for _, tt := range tests {
		got := isProductionLikeURL(tt.url)
		if got != tt.want {
			t.Fatalf("%s: expected %v, got %v", tt.name, tt.want, got)
		}
	}
}

func TestValidateImportSafety(t *testing.T) {
	if err := validateImportSafety("dev", "postgres://u:p@localhost:5432/db", false); err != nil {
		t.Fatalf("dev localhost should pass: %v", err)
	}

	if err := validateImportSafety("dev", "postgres://u:p@prod-db.example.com:5432/db", false); err == nil {
		t.Fatal("expected production-like URL without allow-production to fail")
	}

	if err := validateImportSafety("dev", "postgres://u:p@prod-db.example.com:5432/db", true); err == nil {
		t.Fatal("expected production-like URL with dev env to fail")
	}

	if err := validateImportSafety("prod", "postgres://u:p@prod-db.example.com:5432/db", false); err == nil {
		t.Fatal("expected prod env without allow-production to fail")
	}

	if err := validateImportSafety("prod", "postgres://u:p@prod-db.example.com:5432/db", true); err != nil {
		t.Fatalf("prod env with allow-production should pass: %v", err)
	}

	if err := validateImportSafety("prod", "postgres://u:p@localhost:5432/db", true); err == nil {
		t.Fatal("expected prod env with local URL to fail")
	}
}
