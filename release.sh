#!/usr/bin/env bash
set +x
set -euo pipefail
umask 077

readonly RELEASE_DIR="./release"
readonly PROJECT_URL="https://github.com/antsankov/go-live"
readonly DOWNLOAD_BASE="$PROJECT_URL/releases/download"
readonly TAP_REPO="antsankov/homebrew-go-live"
readonly TAP_FORMULA="go-live.rb"
readonly WINGET_ID="antsankov.go-live"
readonly BOT_NAME="go-live release"
readonly BOT_EMAIL="41898282+github-actions[bot]@users.noreply.github.com"
readonly TIMESTAMP_URL="http://timestamp.digicert.com"
readonly KEYCHAIN_TIMEOUT_SECONDS=3600
readonly VERSION_PATTERN='^[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?$'
readonly PLATFORMS=(mac linux windows)
readonly ARCHES=(x64 arm64)
readonly REQUIRED_SECRETS=(
  APPLE_CERT_P12_BASE64
  APPLE_CERT_PASSWORD
  APPLE_SIGNING_IDENTITY
  APPLE_ID
  APPLE_TEAM_ID
  APPLE_APP_PASSWORD
  WINDOWS_CERT_PFX_BASE64
  WINDOWS_CERT_PASSWORD
  HOMEBREW_TAP_TOKEN
  WINGET_TOKEN
  GH_TOKEN
)
readonly REQUIRED_TOOLS=(go make git gh codesign security xcrun osslsigncode komac zip shasum base64 uuidgen jq)

die() {
  echo "error: $*" >&2
  exit 1
}

asset_name() {
  local os="$1" arch="$2"
  if [[ "$os" == windows ]]; then
    echo "go-live-$os-$arch.exe"
  else
    echo "go-live-$os-$arch"
  fi
}

sha_of() {
  awk -v f="$1" '$2 == f { print $1 }' "$RELEASE_DIR/SHA256SUMS"
}

TAG="${1:-${GITHUB_REF_NAME:-}}"
[[ -n "$TAG" ]] || die "usage: $0 <tag>  (e.g. v1.4.0)"
[[ "$TAG" == v* ]] || die "tag must start with v"
VERSION="${TAG#v}"
[[ "$VERSION" =~ $VERSION_PATTERN ]] || die "tag '$TAG' is not a version (expected vX.Y.Z)"

tag_commit="$(git rev-parse --verify "refs/tags/$TAG^{commit}")" || die "tag does not exist locally"
[[ "$(git rev-parse HEAD)" == "$tag_commit" ]] || die "HEAD does not match $TAG"
[[ -z "$(git status --porcelain --untracked-files=normal)" ]] || die "working tree differs from $TAG"

stale=()
grep -qxF "var VERSION = \"$VERSION\"" main.go || stale+=(main.go)
grep -qxF "LABEL version=\"$VERSION\"" Dockerfile || stale+=(Dockerfile)
grep -qxF "version: '$VERSION'" snap/snapcraft.yaml || stale+=(snap/snapcraft.yaml)
if ((${#stale[@]})); then
  die "version $VERSION not set in: ${stale[*]}"
fi

missing=()
for name in "${REQUIRED_SECRETS[@]}"; do
  [[ -n "${!name:-}" ]] || missing+=("env:$name")
done
for tool in "${REQUIRED_TOOLS[@]}"; do
  command -v "$tool" >/dev/null 2>&1 || missing+=("tool:$tool")
done
if ((${#missing[@]})); then
  printf 'missing: %s\n' "${missing[@]}" >&2
  die "preflight failed (${#missing[@]} missing)"
fi

if gh release view "$TAG" --repo "$PROJECT_URL" >/dev/null 2>&1; then
  die "release $TAG already exists; published assets must not be replaced"
fi

WORKDIR="$(mktemp -d)"
KEYCHAIN="$WORKDIR/release.keychain-db"
ORIGINAL_KEYCHAINS=()
KEYCHAIN_SEARCH_MODIFIED=false

cleanup() {
  if [[ "$KEYCHAIN_SEARCH_MODIFIED" == true ]]; then
    security list-keychains -d user -s "${ORIGINAL_KEYCHAINS[@]}" || true
  fi
  security delete-keychain "$KEYCHAIN" 2>/dev/null || true
  rm -rf "$WORKDIR"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

echo "==> Building go-live $VERSION"
make cross-compile VERSION="$VERSION"
for os in "${PLATFORMS[@]}"; do
  for arch in "${ARCHES[@]}"; do
    asset="$RELEASE_DIR/$(asset_name "$os" "$arch")"
    [[ -f "$asset" ]] || die "build did not produce $asset"
    printf '    %-8s %-6s %s\n' "$os" "$arch" "$asset"
  done
done

echo "==> Importing Apple signing certificate"
keychain_password="$(uuidgen)"
printf '%s' "$APPLE_CERT_P12_BASE64" | base64 --decode >"$WORKDIR/apple.p12"
keychain_list="$(security list-keychains -d user)"
while IFS= read -r keychain; do
  ORIGINAL_KEYCHAINS+=("$keychain")
done < <(printf '%s\n' "$keychain_list" | sed -e 's/^[[:space:]]*"//' -e 's/"$//')
security create-keychain -p "$keychain_password" "$KEYCHAIN"
security set-keychain-settings -lut "$KEYCHAIN_TIMEOUT_SECONDS" "$KEYCHAIN"
security unlock-keychain -p "$keychain_password" "$KEYCHAIN"
security import "$WORKDIR/apple.p12" -k "$KEYCHAIN" -P "$APPLE_CERT_PASSWORD" -T /usr/bin/codesign
security set-key-partition-list -S apple-tool:,apple:,codesign: -s -k "$keychain_password" "$KEYCHAIN" >/dev/null
KEYCHAIN_SEARCH_MODIFIED=true
security list-keychains -d user -s "$KEYCHAIN" "${ORIGINAL_KEYCHAINS[@]}"

for arch in "${ARCHES[@]}"; do
  name="$(asset_name mac "$arch")"
  stage="$WORKDIR/$name"
  mkdir -p "$stage"
  mv "$RELEASE_DIR/$name" "$stage/go-live"
  echo "==> Signing $name"
  codesign --force --options runtime --timestamp --keychain "$KEYCHAIN" \
    --sign "$APPLE_SIGNING_IDENTITY" "$stage/go-live"
  codesign --verify --strict --verbose=2 "$stage/go-live"
  zip -j -q "$RELEASE_DIR/$name.zip" "$stage/go-live"
  echo "==> Notarizing $name.zip"
  xcrun notarytool submit "$RELEASE_DIR/$name.zip" \
    --apple-id "$APPLE_ID" --team-id "$APPLE_TEAM_ID" --password "$APPLE_APP_PASSWORD" --wait \
    --output-format json >"$WORKDIR/notary-result.json"
  jq -e '.status == "Accepted"' "$WORKDIR/notary-result.json" >/dev/null || die "notarization rejected $name"
done

printf '%s' "$WINDOWS_CERT_PFX_BASE64" | base64 --decode >"$WORKDIR/windows.pfx"
printf '%s' "$WINDOWS_CERT_PASSWORD" >"$WORKDIR/windows.pass"
for arch in "${ARCHES[@]}"; do
  exe="$RELEASE_DIR/$(asset_name windows "$arch")"
  echo "==> Signing $exe"
  osslsigncode sign -pkcs12 "$WORKDIR/windows.pfx" -readpass "$WORKDIR/windows.pass" \
    -n go-live -i "$PROJECT_URL" -t "$TIMESTAMP_URL" -in "$exe" -out "$exe.signed"
  mv "$exe.signed" "$exe"
  osslsigncode verify -in "$exe"
done

echo "==> Writing checksums"
(cd "$RELEASE_DIR" && shasum -a 256 go-live-* >SHA256SUMS)

echo "==> Publishing GitHub release $TAG"
gh release create "$TAG" "$RELEASE_DIR"/* --repo "$PROJECT_URL" \
  --title "go-live $TAG" --generate-notes --verify-tag

echo "==> Updating Homebrew tap $TAP_REPO"
tap="$WORKDIR/tap"
cat >"$WORKDIR/git-askpass" <<'ASKPASS'
#!/usr/bin/env bash
case "$1" in
  *Username*) printf '%s\n' x-access-token ;;
  *Password*) printf '%s\n' "$HOMEBREW_TAP_TOKEN" ;;
  *) exit 1 ;;
esac
ASKPASS
chmod 700 "$WORKDIR/git-askpass"
GIT_ASKPASS="$WORKDIR/git-askpass" GIT_TERMINAL_PROMPT=0 \
  git -c credential.helper= clone --quiet --depth 1 "https://github.com/$TAP_REPO.git" "$tap"
cat >"$tap/$TAP_FORMULA" <<FORMULA
class GoLive < Formula
  desc "Portable Go utility that serves files, HTML or otherwise, over HTTP"
  homepage "$PROJECT_URL"
  version "$VERSION"
  license "Apache-2.0"

  on_macos do
    if Hardware::CPU.intel?
      url "$DOWNLOAD_BASE/$TAG/go-live-mac-x64.zip"
      sha256 "$(sha_of go-live-mac-x64.zip)"
    else
      url "$DOWNLOAD_BASE/$TAG/go-live-mac-arm64.zip"
      sha256 "$(sha_of go-live-mac-arm64.zip)"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "$DOWNLOAD_BASE/$TAG/go-live-linux-x64"
      sha256 "$(sha_of go-live-linux-x64)"
    else
      url "$DOWNLOAD_BASE/$TAG/go-live-linux-arm64"
      sha256 "$(sha_of go-live-linux-arm64)"
    end
  end

  def install
    bin.install Dir["go-live*"].first => "go-live"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/go-live --version")
  end
end
FORMULA
if grep -q 'sha256 ""' "$tap/$TAP_FORMULA"; then
  die "missing checksum while writing $TAP_FORMULA"
fi
if git -C "$tap" diff --quiet; then
  echo "    tap already at $VERSION"
else
  git -C "$tap" -c user.name="$BOT_NAME" -c user.email="$BOT_EMAIL" commit -q -am "go-live $VERSION"
  GIT_ASKPASS="$WORKDIR/git-askpass" GIT_TERMINAL_PROMPT=0 \
    git -C "$tap" -c credential.helper= push -q
fi

echo "==> Submitting winget manifest for $WINGET_ID"
komac update "$WINGET_ID" --version "$VERSION" \
  --urls "$DOWNLOAD_BASE/$TAG/$(asset_name windows x64)" "$DOWNLOAD_BASE/$TAG/$(asset_name windows arm64)" \
  --submit --token "$WINGET_TOKEN"
