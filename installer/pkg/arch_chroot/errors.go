package arch_chroot

import "fmt"

// ArchChrootError represents an error that occurred during
// the execution of a command with arch-chroot.
//
// It wraps the original error along with the STDERR output for better debugging.
type ArchChrootError struct {
	Command string
	StdErr  string
	Err     error
}

// Error returns a formatted error message including the content of STDERR
// and the original error message.
func (e ArchChrootError) Error() string {
	return fmt.Sprintf("arch-chroot failed: Command=%s, STDERR=%q, error=%v", e.Command, e.StdErr, e.Err.Error())
}

// Unwrap returns the underlying error for error chaining.
func (e ArchChrootError) Unwrap() error {
	return e.Err
}
