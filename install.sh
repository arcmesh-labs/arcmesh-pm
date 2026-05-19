#!/bin/sh
# Installs apm (arcmesh-pm) for the current user. No admin rights required.
set -e

# ── OS ────────────────────────────────────────────────────────────────────────

os=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$os" in
    linux)  os="linux" ;;
    darwin) os="darwin" ;;
    *)
        echo "Unsupported OS: $os" >&2
        exit 1
        ;;
esac

# ── architecture ──────────────────────────────────────────────────────────────

machine=$(uname -m)
case "$machine" in
    x86_64)          arch="amd64" ;;
    amd64)           arch="amd64" ;;
    aarch64 | arm64) arch="arm64" ;;
    *)
        echo "Unsupported architecture: $machine" >&2
        exit 1
        ;;
esac

# ── paths ─────────────────────────────────────────────────────────────────────

install_dir="$HOME/.local/bin"
bin_path="$install_dir/apm"
tarball_url="https://github.com/arcmesh-labs/arcmesh-pm/releases/latest/download/apm-${os}-${arch}.tar.gz"
tmp_tar="${TMPDIR:-/tmp}/apm-${os}-${arch}.tar.gz"

# ── download ──────────────────────────────────────────────────────────────────

echo "Detected OS: $os"
echo "Detected architecture: $arch"
echo "Downloading $tarball_url ..."

if ! curl -fsSL "$tarball_url" -o "$tmp_tar"; then
    echo "Download failed." >&2
    exit 1
fi

echo "Download complete."

# ── extract ───────────────────────────────────────────────────────────────────

if [ ! -d "$install_dir" ]; then
    mkdir -p "$install_dir"
    echo "Created $install_dir"
fi

echo "Extracting to $install_dir ..."

if ! tar -xzf "$tmp_tar" -C "$install_dir" --strip-components=0 apm 2>/dev/null; then
    # Fallback: extract all and move apm if strip-components isn't enough
    tmp_dir="${TMPDIR:-/tmp}/apm-extract-$$"
    mkdir -p "$tmp_dir"
    tar -xzf "$tmp_tar" -C "$tmp_dir"
    found=$(find "$tmp_dir" -name "apm" -type f | head -n 1)
    if [ -z "$found" ]; then
        echo "Extraction failed: apm binary not found in archive." >&2
        rm -rf "$tmp_dir"
        rm -f "$tmp_tar"
        exit 1
    fi
    mv "$found" "$bin_path"
    rm -rf "$tmp_dir"
fi

rm -f "$tmp_tar"

if [ ! -f "$bin_path" ]; then
    echo "Extraction failed: apm binary not found in $install_dir." >&2
    exit 1
fi

chmod +x "$bin_path"
echo "Installed apm to $bin_path"

# ── PATH ──────────────────────────────────────────────────────────────────────

add_to_path() {
    rc_file="$1"
    if [ -f "$rc_file" ]; then
        if grep -qF "$install_dir" "$rc_file" 2>/dev/null; then
            echo "$install_dir is already in PATH in $rc_file"
        else
            printf '\nexport PATH="%s:$PATH"\n' "$install_dir" >> "$rc_file"
            echo "Added $install_dir to PATH in $rc_file"
        fi
    fi
}

add_to_path "$HOME/.bashrc"
add_to_path "$HOME/.zshrc"

# ── done ──────────────────────────────────────────────────────────────────────

echo ""
echo "apm installed successfully."
echo "Open a new terminal for PATH changes to take effect, then run: apm --help"
