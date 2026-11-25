#!/usr/bin/env bash
set -euo pipefail

# cimgui-go only ships Linux/Windows archives for x64. Build the same sources
# from the Go module verified by go.sum for the ARM64 release jobs.
target="${1:?Usage: build-arm64-gui-libs.sh linux|windows}"
module_dir="$(go list -m -f '{{.Dir}}' github.com/AllenDang/cimgui-go)"
mkdir -p build/native-arm64
native_dir="$(cd build/native-arm64 && pwd)"
cmake_args=(-G Ninja -DCMAKE_BUILD_TYPE=Release)

case "$target" in
  linux)
    if [[ "$(go env GOHOSTARCH)" != arm64 ]]; then
      echo 'The Linux native libraries must be built on an ARM64 runner.' >&2
      exit 1
    fi
    system_libs='-ldl -lGL -lX11 -lm -lpthread'
    ;;
  windows)
    module_dir="$(cygpath -u "$module_dir")"
    cmake_args+=(
      -DCMAKE_SYSTEM_NAME=Windows
      -DCMAKE_SYSTEM_PROCESSOR=ARM64
      -DCMAKE_C_COMPILER=aarch64-w64-mingw32-clang
      -DCMAKE_CXX_COMPILER=aarch64-w64-mingw32-clang++
      -DCMAKE_RC_COMPILER=aarch64-w64-mingw32-windres
      -DCMAKE_AR=llvm-ar
      -DCMAKE_RANLIB=llvm-ranlib
    )
    system_libs='-lgdi32 -lopengl32 -limm32'
    ;;
  *)
    echo "Unsupported target: $target" >&2
    exit 1
    ;;
esac

if [[ ! -f "$native_dir/glfw/src/libglfw3.a" ]]; then
  cmake -S "$module_dir/thirdparty/glfw" -B "$native_dir/glfw" \
    "${cmake_args[@]}" -DBUILD_SHARED_LIBS=OFF \
    -DGLFW_BUILD_EXAMPLES=OFF -DGLFW_BUILD_TESTS=OFF -DGLFW_BUILD_DOCS=OFF
  cmake --build "$native_dir/glfw" --parallel 2
fi

if [[ ! -f "$native_dir/cimgui/cimgui.a" ]]; then
  cmake -S "$module_dir/lib" -B "$native_dir/cimgui" \
    "${cmake_args[@]}" -DIMGUI_STATIC=ON -DIMGUI_FREETYPE=OFF -DCIMGUI_TEST=OFF
  cmake --build "$native_dir/cimgui" --parallel 2
fi

if [[ "$target" == windows ]]; then
  native_dir="$(cygpath -m "$native_dir")"
fi

link_flags="$native_dir/cimgui/cimgui.a $native_dir/glfw/src/libglfw3.a $system_libs"
if [[ -n "${GITHUB_ENV:-}" ]]; then
  echo "CGO_LDFLAGS=$link_flags" >> "$GITHUB_ENV"
else
  echo "CGO_LDFLAGS=$link_flags"
fi
