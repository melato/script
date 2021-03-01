// Package script provides convenience methods for running external programs
// Features/Goals/Changes:
//  - No global Trace flag.  Use Script to wrap running commands with Trace, DryRun
//  - Use exec.Cmd whenever possible.  Avoid wrapping too much.
//  - remove Cmd.Pipe(), Cmd.PipeTo()
//
// Used by:
// go/pin
// go/lxdops
package script
