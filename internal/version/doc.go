// Package version exposes build-time metadata for the Probelet binary.
//
// These variables are intended to be populated at build time using Go linker
// flags (-ldflags), allowing commands such as `probelet version` to report
// the exact version, commit, and build date of the running binary.
package version
