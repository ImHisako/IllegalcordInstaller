/*
 * SPDX-License-Identifier: GPL-3.0
 * Vencord Installer, a cross platform gui/cli app for installing Vencord
 * Copyright (c) 2023 Vendicated and Vencord contributors
 */

package main

import (
	"equilotl/buildinfo"
)

const ReleaseUrl = "https://api.github.com/repos/ImHisako/Illegalcord/releases/latest"
const InstallerReleaseUrl = "https://api.github.com/repos/ImHisako/IllegalcordInstaller/releases/latest"

var UserAgent = "IllegalcordInstaller/" + buildinfo.InstallerGitHash + " (https://github.com/ImHisako/IllegalcordInstaller)"

const SupportUrl = "https://github.com/ImHisako/IllegalcordInstaller/issues"

var LinuxDiscordNames = []string{
	"Discord",
	"DiscordPTB",
	"DiscordCanary",
	"DiscordDevelopment",
	"discord",
	"discordptb",
	"discordcanary",
	"discorddevelopment",
	"discord-ptb",
	"discord-canary",
	"discord-development",
	// Flatpak
	"com.discordapp.Discord",
	"com.discordapp.DiscordPTB",
	"com.discordapp.DiscordCanary",
	"com.discordapp.DiscordDevelopment",
}
