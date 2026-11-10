/*
 * SPDX-License-Identifier: GPL-3.0
 * Vencord Installer, a cross platform gui/cli app for installing Vencord
 * Copyright (c) 2023 Vendicated and Vencord contributors
 */

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	path "path/filepath"
	"regexp"
	"strings"
)

type GithubRelease struct {
	Name    string `json:"name"`
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name        string `json:"name"`
		DownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

var ReleaseData GithubRelease
var GithubError error
var GithubDoneChan chan bool

var InstalledHash = "None"
var LatestHash = "Unknown"
var IsDevInstall bool

func GetGithubRelease(url string) (*GithubRelease, error) {
	Log.Debug("Fetching", url)

	res, err := httpGet(url)
	if err != nil {
		Log.Error("Failed to send Request", err)
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		err = errors.New(res.Status)
		Log.Error(url, "returned Non-OK status", err)
		return nil, err
	}
	var data GithubRelease

	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		Log.Error("Failed to decode GitHub JSON Response", err)
		return nil, err
	}

	return &data, nil
}

func InitGithubDownloader() {
	GithubDoneChan = make(chan bool, 1)

	IsDevInstall = os.Getenv("EQUICORD_DEV_INSTALL") == "1"
	Log.Debug("Is Dev Install: ", IsDevInstall)
	if IsDevInstall {
		GithubDoneChan <- true
		return
	}

	go func() {
		// Make sure UI updates once the request either finished or failed
		defer func() {
			GithubDoneChan <- GithubError == nil
		}()

		data, err := GetGithubRelease(ReleaseUrl)
		if err != nil {
			GithubError = err
			return
		}

		ReleaseData = *data

		i := strings.LastIndex(data.Name, " ") + 1
		LatestHash = data.Name[i:]
		Log.Debug("Finished fetching GitHub Data")
		Log.Debug("Latest hash is", LatestHash, "Local Install is", Ternary(LatestHash == InstalledHash, "up to date!", "outdated!"))
	}()

	// either .asar file or directory with main.js file (in DEV)
	EquicordFile := EquicordDirectory

	stat, err := os.Stat(EquicordFile)
	if err != nil {
		return
	}

	// dev
	if stat.IsDir() {
		EquicordFile = path.Join(EquicordFile, "main.js")
	}

	// Check hash of installed version if exists
	b, err := os.ReadFile(EquicordFile)
	if err != nil {
		return
	}

	Log.Debug("Found existing Illegalcord Install. Checking for hash...")

	re := regexp.MustCompile(`// Equicord (\w+)`)
	match := re.FindSubmatch(b)
	if match != nil {
		InstalledHash = string(match[1])
		Log.Debug("Existing hash is", InstalledHash)

	} else {
		Log.Debug("Didn't find hash")

	}
}

func installLatestBuilds() (retErr error) {
	Log.Debug("Installing latest builds...")

	if IsDevInstall {
		Log.Debug("Skipping due to dev install")
		return
	}

	downloadUrl := ""
	for _, ass := range ReleaseData.Assets {
		if ass.Name == "desktop.asar" {
			downloadUrl = ass.DownloadURL
			break
		}
	}

	if downloadUrl == "" {
		retErr = errors.New("Didn't find desktop.asar download link")
		Log.Error(retErr)
		return
	}

	Log.Debug("Downloading desktop.asar")

	res, err := httpGet(downloadUrl)
	if err != nil {
		Log.Error("Failed to download desktop.asar:", err)
		retErr = err
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		retErr = errors.New(res.Status)
		Log.Error("Failed to download desktop.asar:", retErr)
		return
	}
	out, err := os.CreateTemp(path.Dir(EquicordDirectory), ".desktop-*.asar")
	if err != nil {
		Log.Error("Failed to create", EquicordDirectory+":", err)
		retErr = err
		return
	}
	defer func() {
		_ = out.Close()
		_ = os.Remove(out.Name())
	}()
	read, err := io.Copy(out, res.Body)
	if err != nil {
		Log.Error("Failed to download to", EquicordDirectory+":", err)
		retErr = err
		return
	}
	if res.ContentLength >= 0 && read != res.ContentLength {
		retErr = fmt.Errorf("unexpected download size: expected %d bytes, read %d", res.ContentLength, read)
		Log.Error(retErr)
		return
	}
	if err := out.Chmod(0644); err != nil {
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}
	if err := os.Rename(out.Name(), EquicordDirectory); err != nil {
		return err
	}

	_ = FixOwnership(EquicordDirectory)

	InstalledHash = LatestHash
	return
}
