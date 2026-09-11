// Package version holds the single source of truth for the CLI version.
// Bump it here when cutting a release; the client User-Agent, `mobileops version`
// and the agent manifest all read from it.
package version

const Version = "0.3.0"
