#!/bin/sh
# Install the latest kicad-init release from GitHub.
set -eu

repo='https://github.com/avycado13/kicad-init'
install_dir=${INSTALL_DIR:-"$HOME/.local/bin"}

case "$(uname -s)" in
  Darwin) os=Darwin ;;
  Linux) os=Linux ;;
  *) echo 'Unsupported OS (use the GitHub Releases page instead).' >&2; exit 1 ;;
esac
case "$(uname -m)" in
  x86_64|amd64) arch=x86_64 ;;
  arm64|aarch64) arch=arm64 ;;
  i386|i486|i586|i686) arch=i386 ;;
  *) echo 'Unsupported architecture.' >&2; exit 1 ;;
esac

archive="kicad-init_${os}_${arch}.tar.gz"
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT HUP INT TERM

fetch() {
  if command -v curl >/dev/null 2>&1; then
    curl -fLsS --retry 3 "$1" -o "$2"
  elif command -v wget >/dev/null 2>&1; then
    wget -q "$1" -O "$2"
  else
    echo 'curl or wget is required.' >&2
    exit 1
  fi
}

fetch "https://api.github.com/repos/avycado13/kicad-init/releases/latest" "$tmp/release.json"
tag=$(sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' "$tmp/release.json" | head -n 1)
case "$tag" in
  ''|*[!a-zA-Z0-9._-]*) echo 'Could not determine the latest release tag.' >&2; exit 1 ;;
esac
version=${tag#v}
release="$repo/releases/download/$tag"
fetch "$release/$archive" "$tmp/$archive"
fetch "$release/kicad-init_${version}_checksums.txt" "$tmp/checksums.txt"
expected=$(awk -v name="$archive" '$2 == name { print $1 }' "$tmp/checksums.txt")
if [ -z "$expected" ]; then
  echo "No checksum found for $archive" >&2
  exit 1
fi
if command -v shasum >/dev/null 2>&1; then
  actual=$(shasum -a 256 "$tmp/$archive" | awk '{ print $1 }')
elif command -v sha256sum >/dev/null 2>&1; then
  actual=$(sha256sum "$tmp/$archive" | awk '{ print $1 }')
else
  echo 'shasum or sha256sum is required to verify the download.' >&2
  exit 1
fi
if [ "$actual" != "$expected" ]; then
  echo 'Checksum mismatch; refusing to install.' >&2
  exit 1
fi

tar -xzf "$tmp/$archive" -C "$tmp" kicad-init
mkdir -p "$install_dir"
install -m 755 "$tmp/kicad-init" "$install_dir/kicad-init"
printf 'Installed kicad-init to %s/kicad-init\n' "$install_dir"
case ":$PATH:" in
  *":$install_dir:"*) ;;
  *) printf 'Add %s to your PATH to run kicad-init.\n' "$install_dir" ;;
esac
