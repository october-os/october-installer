package arch_chroot

import "fmt"

// ArchChrootError represents an error that occurred during
// the execution of a command with arch-chroot.
//
// It wraps the original error along with the STDERR output for better debugging.
type ArchChrootError struct {
	command string
	stdErr  string
	err     error
}

// Error returns a formatted error message including the content of STDERR
// and the original error message.
func (e ArchChrootError) Error() string {
	return fmt.Sprintf("Error running arch-chroot: command=%s, stderr=%q, error=%v", e.command, e.stdErr, e.err.Error())
}

// Unwrap returns the original error wrapped inside.
func (e ArchChrootError) Unwrap() error {
	return e.err
}
