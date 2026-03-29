#!/usr/bin/env bash
# ABOUTME: Release script for sanitize. Creates a git tag, GitHub release,
# and updates the Homebrew formula in the tap repo.
set -euo pipefail
IFS=$'\n\t'

VERSION="$1"
CURRENT_VERSION="$2"
REPO="$3"
TAP_REPO="$4"
BINARY="$5"

# Determine new version
if [ "${VERSION}" = "${CURRENT_VERSION}" ]; then
    NEW_VERSION=$(echo "${CURRENT_VERSION}" | awk -F. '{printf "%s.%s.%s", $1, $2+1, $3}')
else
    NEW_VERSION="${VERSION}"
fi

echo "Releasing v${NEW_VERSION}..."

# Update VERSION file and build
echo "${NEW_VERSION}" > VERSION
go build -ldflags "-X main.version=${NEW_VERSION}" -o "${BINARY}" ./cmd/sanitize/

# Commit, tag, push
git add VERSION
git commit -m "release: v${NEW_VERSION}"
git tag "v${NEW_VERSION}"
git push
git push --tags

# Create GitHub release
echo "Creating GitHub release..."
gh release create "v${NEW_VERSION}" \
    --repo "${REPO}" \
    --title "v${NEW_VERSION}" \
    --generate-notes

# Download tarball and compute SHA256
echo "Computing SHA256..."
TARBALL_URL="https://github.com/${REPO}/archive/refs/tags/v${NEW_VERSION}.tar.gz"
TMPFILE=$(mktemp)
cleanup() {
    rm -f "${TMPFILE}"
}
trap cleanup EXIT

# Retry download — GitHub may take a moment to make the archive available
for i in 1 2 3 4 5; do
    if curl -sfL "${TARBALL_URL}" -o "${TMPFILE}"; then
        break
    fi
    echo "Waiting for archive to be available (attempt ${i})..."
    sleep 2
done

if [ ! -s "${TMPFILE}" ]; then
    echo "Error: failed to download tarball from ${TARBALL_URL}"
    exit 1
fi

SHA256=$(shasum -a 256 "${TMPFILE}" | awk '{print $1}')
echo "SHA256: ${SHA256}"

# Clone tap, update formula, push
echo "Updating Homebrew formula..."
TAP_DIR=$(mktemp -d)
trap 'rm -rf "${TMPFILE}" "${TAP_DIR}"' EXIT

gh repo clone "${TAP_REPO}" "${TAP_DIR}" -- --depth 1

cat > "${TAP_DIR}/Formula/sanitize.rb" << FORMULA
class Sanitize < Formula
  desc "Fast CLI tool for converting English text to Oxford (OED) spelling"
  homepage "https://github.com/${REPO}"
  url "${TARBALL_URL}"
  sha256 "${SHA256}"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", "-ldflags", "-X main.version=${NEW_VERSION}", "-o", bin/"${BINARY}", "./cmd/sanitize/"
  end

  test do
    assert_match "sanitize", shell_output("#{bin}/${BINARY} --version")
    output = pipe_output("#{bin}/${BINARY} oed -q", "organise the center")
    assert_equal "organize the centre", output.strip
  end
end
FORMULA

cd "${TAP_DIR}"
git add Formula/sanitize.rb
git commit -m "Update sanitize formula to v${NEW_VERSION}"
git push

echo "Released v${NEW_VERSION}"
