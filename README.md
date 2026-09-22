# Illegalcord Installer

Install [Illegalcord, a Discord Desktop client mod](https://github.com/ImHisako/Illegalcord)

Based on [Equilotl](https://github.com/Equicord/Equilotl). Installer executables and release packages are named `Illegalotl*` (`Illegalotl.exe` for the Windows GUI and `IllegalotlCli.exe` for the CLI).

## Usage

Windows

| | X64 | ARM64 |
| --- | --- | --- |
| GUI | [Download](https://github.com/ImHisako/IllegalcordInstaller/releases/latest/download/Illegalotl.exe) | [Download](https://github.com/ImHisako/IllegalcordInstaller/releases/latest/download/Illegalotl-arm64.exe) |
| CLI | [Download](https://github.com/ImHisako/IllegalcordInstaller/releases/latest/download/IllegalotlCli.exe) | [Download](https://github.com/ImHisako/IllegalcordInstaller/releases/latest/download/IllegalotlCli-arm64.exe) |

MacOS

| | Universal | X64 | ARM64 |
| --- | --- | --- | --- |
| GUI | [Download](https://github.com/ImHisako/IllegalcordInstaller/releases/latest/download/Illegalotl.dmg) | [Download](https://github.com/ImHisako/IllegalcordInstaller/releases/latest/download/Illegalotl-x64.dmg) | [Download](https://github.com/ImHisako/IllegalcordInstaller/releases/latest/download/Illegalotl-arm64.dmg) |
| CLI | [Download](https://github.com/ImHisako/IllegalcordInstaller/releases/latest/download/IllegalotlCli-universal) | [Download](https://github.com/ImHisako/IllegalcordInstaller/releases/latest/download/IllegalotlCli-x64) | [Download](https://github.com/ImHisako/IllegalcordInstaller/releases/latest/download/IllegalotlCli-arm64) |

macOS builds are published when Apple signing and notarization credentials are configured in the release workflow.

The CLI builds are plain binaries, so run `chmod +x <file>` after downloading.

Linux

| | X64 | ARM64 |
| --- | --- | --- |
| Combined GUI | [Download](https://github.com/ImHisako/IllegalcordInstaller/releases/latest/download/Illegalotl) | [Download](https://github.com/ImHisako/IllegalcordInstaller/releases/latest/download/Illegalotl-arm64) |
| X11 GUI | [Download](https://github.com/ImHisako/IllegalcordInstaller/releases/latest/download/Illegalotl-x11) | [Download](https://github.com/ImHisako/IllegalcordInstaller/releases/latest/download/Illegalotl-x11-arm64) |
| Wayland GUI | [Download](https://github.com/ImHisako/IllegalcordInstaller/releases/latest/download/Illegalotl-wayland) | [Download](https://github.com/ImHisako/IllegalcordInstaller/releases/latest/download/Illegalotl-wayland-arm64) |
| CLI | [Download](https://github.com/ImHisako/IllegalcordInstaller/releases/latest/download/IllegalotlCli-linux) | [Download](https://github.com/ImHisako/IllegalcordInstaller/releases/latest/download/IllegalotlCli-linux-arm64) |


## Building from source

### Prerequisites

You need to install the [Go programming language](https://go.dev/doc/install) and GCC, the GNU Compiler Collection (MinGW on Windows)

<details>
<summary>Additionally, if you're using Linux, you have to install some additional dependencies:</summary>

#### Base dependencies

```sh
apt install -y pkg-config libsdl2-dev libglx-dev libgl1-mesa-dev
dnf install pkg-config libGL-devel libXxf86vm-devel
```

#### X11 dependencies

```sh
apt install -y xorg-dev
dnf install libXcursor-devel libXi-devel libXinerama-devel libXrandr-devel
```

#### Wayland dependencies

```sh
apt install -y libwayland-dev libxkbcommon-dev wayland-protocols extra-cmake-modules
dnf install wayland-devel libxkbcommon-devel wayland-protocols-devel extra-cmake-modules
```

</details>

### Building

#### Install dependencies

```sh
go mod tidy
```

#### Build the GUI

##### Windows / Mac / Linux X11

```sh
make GUI=1
```

##### Linux Wayland

```sh
make GUI=1 WAYLAND=1
```

#### Build the CLI

```
make GUI=0
```

You might want to pass some flags to this command to get a better build.
See [the GitHub workflow](https://github.com/ImHisako/IllegalcordInstaller/blob/main/.github/workflows/release.yml) for what flags I pass or if you want more precise instructions

## Compatibility

Illegalcord currently uses the upstream `EQUICORD_USER_DATA_DIR`, `EQUICORD_DIRECTORY`, and `EQUICORD_DEV_INSTALL` variables, the `Equicord` data directory, and `equicord.asar`. These names are intentionally preserved. Client builds and installer updates are downloaded from the ImHisako repositories.

When upgrading from an installer named `Equilotl*`, download the new installer manually. Development scripts that download the installer must also use the `Illegalotl*` asset names.
