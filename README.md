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

No MacOS Version because i don't like Mac. 
Jk i don't have the developer cert of apple, because their ToS and Policy is a joke. so install illegalcord from the dev build.

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

The Linux GUI uses the x86-64 GLFW and ImGui libraries bundled with cimgui-go.

```sh
apt install -y libgl-dev libx11-dev
dnf install libglvnd-devel libX11-devel
```

</details>

### Building

#### Install dependencies

```sh
go mod tidy
```

#### Build the GUI

##### Windows

Install [MSYS2](https://www.msys2.org), open the **MSYS2 MINGW64** shell and run:

```sh
pacman -S --needed git make mingw-w64-x86_64-gcc mingw-w64-x86_64-go
go install github.com/tc-hib/go-winres@v0.3.3
export PATH="$PATH:$(go env GOPATH)/bin"
make GUI=1
make GUI=0
```

This writes `build/Illegalotl.exe` and `build/IllegalotlCli.exe`. The GUI dependencies bundle native libraries for Windows x64, Linux x64, and macOS x64/ARM64. The release workflow builds the additional native libraries needed for Windows and Linux ARM64 GUIs from the verified Go module sources. Local ARM64 GUI builds need these libraries too; the CLI can be built for ARM64 with `CGO_ENABLED=0`.

If you'd rather not use MSYS2, any MinGW toolchain works by pointing cgo at it, for example LLVM's clang with a MinGW sysroot:

```sh
CGO_LDFLAGS=-lwinpthread CC="clang --target=x86_64-w64-windows-gnu" CXX="clang++ --target=x86_64-w64-windows-gnu" \
  CGO_ENABLED=1 go build -tags "static gui" -ldflags "-s -w -H=windowsgui -extldflags=-static" -o Illegalotl.exe
```

##### Mac / Linux

On Linux the same binary selects X11 or Wayland at startup.

```sh
make GUI=1
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
