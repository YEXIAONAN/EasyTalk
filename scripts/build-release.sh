#!/usr/bin/env bash
#
# EasyTalk 发布构建脚本。
# 从源码完整构建并打包六个平台：Windows / Linux / macOS × amd64 / arm64。
#
# 用法:
#   ./scripts/build-release.sh              # 版本缺省为 dev
#   VERSION=v0.1.0 ./scripts/build-release.sh
#   ./scripts/build-release.sh v0.1.0
#
# 输出:
#   release/easytalk-<version>-<os>-<arch>.zip / .tar.gz
#   release/SHA256SUMS.txt
#
# 本脚本只构建打包，不创建 Git Tag、不发布到 GitHub。
set -euo pipefail

cd "$(dirname "$0")/.."

VERSION="${1:-${VERSION:-}}"
if [ -z "$VERSION" ]; then
  VERSION="dev"
fi

echo "==> Version: $VERSION"

# 1. 前端构建（go:embed 需要 web/dist）
echo "==> Building frontend"
(
  cd web
  npm ci
  npm run build
)

# 2. 后端校验
echo "==> Running go test / vet"
go test ./...
go vet ./...

# 3. 交叉编译 + 打包
RELEASE_DIR="release"
rm -rf "$RELEASE_DIR"
mkdir -p "$RELEASE_DIR"

LDFLAGS="-s -w -X easytalk/internal/version.Version=${VERSION}"

package() {
  local goos="$1" goarch="$2" ext="${3:-}"
  local bin="easytalk${ext}"
  local staging="$RELEASE_DIR/.staging-${goos}-${goarch}"
  local archive="easytalk-${VERSION}-${goos}-${goarch}"

  mkdir -p "$staging"
  CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" \
    go build -trimpath -ldflags "$LDFLAGS" -o "$staging/$bin" ./cmd/easytalk

  cp config.example.json .env.example README.md "$staging/"

  if [ "$goos" = "windows" ]; then
    (cd "$staging" && zip -q -r "../$archive.zip" .)
  else
    tar -czf "$RELEASE_DIR/$archive.tar.gz" -C "$staging" .
  fi

  rm -rf "$staging"
}

package windows amd64 .exe
package windows arm64 .exe
package linux amd64
package linux arm64
package darwin amd64
package darwin arm64

# 4. SHA256 校验文件（兼容 macOS / Linux：优先 sha256sum，回退 shasum）
checksum() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1"
  elif command -v shasum >/dev/null 2>&1; then
    shasum -a 256 "$1" | awk '{print $1"  "$2}'
  else
    echo "error: neither sha256sum nor shasum found" >&2
    exit 1
  fi
}

(
  cd "$RELEASE_DIR"
  : > SHA256SUMS.txt
  for f in easytalk-"$VERSION"-*.zip easytalk-"$VERSION"-*.tar.gz; do
    checksum "$f" >> SHA256SUMS.txt
  done
)

# 5. 校验六个平台包完整
count="$(find "$RELEASE_DIR" -maxdepth 1 -type f \( -name '*.zip' -o -name '*.tar.gz' \) | wc -l | tr -d ' ')"
if [ "$count" -ne 6 ]; then
  echo "error: expected 6 archives, got $count" >&2
  exit 1
fi

echo "==> Release artifacts:"
ls -lh "$RELEASE_DIR"