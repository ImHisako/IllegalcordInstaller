package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestUnpatchLegacyDirectoryLoader(t *testing.T) {
	for _, manifest := range []string{`{"name":"discord","main":"index.js"}`, PackageJson} {
		t.Run(manifest, func(t *testing.T) {
			dir := t.TempDir()
			loader := filepath.Join(dir, "app.asar")
			if err := os.Mkdir(loader, 0755); err != nil {
				t.Fatal(err)
			}
			writeTestFile(t, filepath.Join(loader, "package.json"), []byte(manifest))
			writeTestFile(t, filepath.Join(loader, "index.js"), []byte(`require("C:\\Users\\Example\\AppData\\Roaming\\Equicord\\equicord.asar\\patcher.js");`))
			original := []byte("original Discord archive")
			writeTestFile(t, filepath.Join(dir, "_app.asar"), original)

			if err := unpatchAppAsar(dir, false); err != nil {
				t.Fatal(err)
			}
			assertTestFile(t, loader, original)
			for _, name := range []string{"_app.asar", "app.asar.tmp"} {
				if _, err := os.Stat(filepath.Join(dir, name)); !os.IsNotExist(err) {
					t.Fatalf("%s was not removed: %v", name, err)
				}
			}
			// Repair must be able to replace the legacy loader with the current format.
			if err := patchAppAsar(dir, false); err != nil {
				t.Fatal(err)
			}
			if ok, err := isEquicordLoaderAppAsar(loader); !ok || err != nil {
				t.Fatalf("new loader not recognized: %v", err)
			}
			assertTestFile(t, filepath.Join(dir, "_app.asar"), original)
		})
	}
}

func TestUnpatchUnknownDirectoryPreservesBackup(t *testing.T) {
	for _, manifest := range []string{"", "invalid json", `{"name":"discord","main":"other.js"}`, `{"name":"discord","main":"index.js"}`} {
		t.Run(manifest, func(t *testing.T) {
			dir := t.TempDir()
			loader := filepath.Join(dir, "app.asar")
			if err := os.Mkdir(loader, 0755); err != nil {
				t.Fatal(err)
			}
			if manifest != "" {
				writeTestFile(t, filepath.Join(loader, "package.json"), []byte(manifest))
			}
			index := []byte("// Not a recognized loader")
			writeTestFile(t, filepath.Join(loader, "index.js"), index)
			original := []byte("original Discord archive")
			writeTestFile(t, filepath.Join(dir, "_app.asar"), original)

			if err := unpatchAppAsar(dir, false); err == nil {
				t.Fatal("unrecognized directory must not be treated as a Discord update")
			}
			assertTestFile(t, filepath.Join(dir, "_app.asar"), original)
			assertTestFile(t, filepath.Join(loader, "index.js"), index)
		})
	}
}

func writeTestFile(t *testing.T, name string, data []byte) {
	t.Helper()
	if err := os.WriteFile(name, data, 0644); err != nil {
		t.Fatal(err)
	}
}

func assertTestFile(t *testing.T, name string, want []byte) {
	t.Helper()
	got, err := os.ReadFile(name)
	if err != nil || !bytes.Equal(got, want) {
		t.Fatalf("%s: got %q, want %q, error: %v", name, got, want, err)
	}
}
