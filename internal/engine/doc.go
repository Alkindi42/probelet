// Package engine provides the low-level execution primitives used by Probelet.
//
// It contains the core runtime operations that perform work such as CPU or
// memory stress, as well as simple in-memory state used by the service.
//
// The engine package is intentionally transport-agnostic: it does not know
// about HTTP, CLI commands, or request parsing.
package engine
