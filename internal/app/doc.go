// Package app contains Probelet application logic shared by multiple interfaces.
//
// It defines the public behavior of the service independently from transport
// concerns such as HTTP handlers or CLI commands.
//
// The package is responsible for:
//   - validating user inputs
//   - enforcing Probelet safety limits and defaults
//   - orchestrating calls to the underlying engine package
//
// This layer allows different frontends, such as the HTTP API and the CLI,
// to share the same rules and behavior without duplicating validation logic.
package app
