// Package buildinfo holds build-time metadata. Version and Commit are overridden
// at release time via -ldflags "-X easytalk/internal/buildinfo.Version=...".
package buildinfo

// Name is the product name.
const Name = "EasyTalk"

// Version is the semantic version, defaulting to "dev" for local builds.
var Version = "dev"

// Commit is the short git commit SHA, defaulting to "unknown".
var Commit = "unknown"

// Repository is the canonical source repository URL.
const Repository = "https://github.com/YEXIAONAN/EasyTalk"