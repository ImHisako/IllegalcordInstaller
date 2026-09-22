package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstallerDownloadLinks(t *testing.T) {
	for _, tc := range []struct {
		os, arch string
		cli      bool
		asset    string
	}{
		{"windows", "amd64", false, "Equilotl.exe"},
		{"windows", "amd64", true, "EquilotlCli.exe"},
		{"windows", "arm64", false, "Equilotl-arm64.exe"},
		{"windows", "arm64", true, "EquilotlCli-arm64.exe"},
		{"linux", "amd64", false, "Equilotl"},
		{"linux", "amd64", true, "EquilotlCli-linux"},
		{"linux", "arm64", false, "Equilotl-arm64"},
		{"linux", "arm64", true, "EquilotlCli-linux-arm64"},
		{"darwin", "amd64", false, "Equilotl-x64.dmg"},
		{"darwin", "amd64", true, "EquilotlCli-x64"},
		{"darwin", "arm64", false, "Equilotl-arm64.dmg"},
		{"darwin", "arm64", true, "EquilotlCli-arm64"},
		{"linux", "386", true, ""},
		{"freebsd", "amd64", true, ""},
	} {
		t.Run(fmt.Sprintf("%s/%s/cli=%t", tc.os, tc.arch, tc.cli), func(t *testing.T) {
			want := ""
			if tc.asset != "" {
				want = "https://github.com/ImHisako/IllegalcordInstaller/releases/latest/download/" + tc.asset
			}
			if got := installerDownloadLink(tc.os, tc.arch, tc.cli); got != want {
				t.Fatalf("got %q, want %q", got, want)
			}
		})
	}
}

func TestReleaseHTTPFailures(t *testing.T) {
	for _, status := range []int{403, 404, 429, 500} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(status)
				fmt.Fprint(w, `{"message":"release unavailable"}`)
			}))
			defer server.Close()
			if _, err := GetGithubRelease(server.URL); err == nil {
				t.Fatal("HTTP error was accepted as a release")
			}
		})
	}
}

func TestDownloadAndPatchRoundTrip(t *testing.T) {
	oldRelease, oldDir, oldDev := ReleaseData, EquicordDirectory, IsDevInstall
	oldInstalled, oldLatest := InstalledHash, LatestHash
	t.Cleanup(func() {
		ReleaseData, EquicordDirectory, IsDevInstall = oldRelease, oldDir, oldDev
		InstalledHash, LatestHash = oldInstalled, oldLatest
	})
	dir := t.TempDir()
	EquicordDirectory = filepath.Join(dir, "equicord.asar")
	IsDevInstall = false
	LatestHash = "abc123"
	payload := "// Equicord abc123\n// Illegalcord desktop fixture"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/release":
			if r.UserAgent() != UserAgent || !strings.HasPrefix(r.UserAgent(), "IllegalcordInstaller/") {
				t.Error("release request is missing Illegalcord identification")
			}
			json.NewEncoder(w).Encode(map[string]any{
				"name": "Build abc123", "tag_name": "dev",
				"assets": []map[string]string{{"name": "desktop.asar", "browser_download_url": "http://" + r.Host + "/desktop.asar"}},
			})
		case "/desktop.asar":
			// Exercise a chunked response without Content-Length.
			w.(http.Flusher).Flush()
			fmt.Fprint(w, payload)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	release, err := GetGithubRelease(server.URL + "/release")
	if err != nil {
		t.Fatal(err)
	}
	ReleaseData = *release
	if err := installLatestBuilds(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(EquicordDirectory)
	if err != nil || string(data) != payload || InstalledHash != LatestHash {
		t.Fatalf("download mismatch: %q, %v, hash %q", data, err, InstalledHash)
	}
	original := []byte("Original Discord archive fixture")
	appAsar := filepath.Join(dir, "app.asar")
	if err := os.WriteFile(appAsar, original, 0644); err != nil {
		t.Fatal(err)
	}
	if err := patchAppAsar(dir, false); err != nil {
		t.Fatal(err)
	}
	if loader, err := isEquicordLoaderAppAsar(appAsar); err != nil || !loader {
		t.Fatalf("Illegalcord loader was not recognized: %v", err)
	}
	data, err = os.ReadFile(appAsar)
	if err != nil || !bytes.Contains(data, []byte("equicord.asar")) {
		t.Fatalf("loader does not reference the downloaded archive: %v", err)
	}
	if err := unpatchAppAsar(dir, false); err != nil {
		t.Fatal(err)
	}
	data, err = os.ReadFile(appAsar)
	if err != nil || !bytes.Equal(data, original) {
		t.Fatalf("uninstall did not restore Discord: %q, %v", data, err)
	}
	ReleaseData.Assets[0].DownloadURL = server.URL + "/missing"
	if err := installLatestBuilds(); err == nil {
		t.Fatal("missing asset was accepted")
	}
	data, err = os.ReadFile(EquicordDirectory)
	if err != nil || string(data) != payload {
		t.Fatal("HTTP error overwrote the installed client")
	}
}
