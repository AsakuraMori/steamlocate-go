package steamlocate

import (
	"fmt"
	"io"
	"path/filepath"
)

// ParseErrorKind represents the kind of file that failed to parse
type ParseErrorKind int

const (
	ParseErrorKindConfig ParseErrorKind = iota
	ParseErrorKindLibraryFolders
	ParseErrorKindApp
	ParseErrorKindShortcut
)

func (k ParseErrorKind) String() string {
	switch k {
	case ParseErrorKindConfig:
		return "config"
	case ParseErrorKindLibraryFolders:
		return "libraryfolders"
	case ParseErrorKindApp:
		return "app"
	case ParseErrorKindShortcut:
		return "shortcut"
	default:
		return "unknown"
	}
}

// Error is the main error type for steamlocate
type Error struct {
	Type    ErrorType
	Message string
	Path    string
	Cause   error
}

type ErrorType int

const (
	ErrorTypeLocate ErrorType = iota
	ErrorTypeValidation
	ErrorTypeIO
	ErrorTypeParse
	ErrorTypeMissingApp
)

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Cause
}

// IsLocateError checks if error is a locate error
func IsLocateError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Type == ErrorTypeLocate
	}
	return false
}

// IsValidationError checks if error is a validation error
func IsValidationError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Type == ErrorTypeValidation
	}
	return false
}

// IsIOError checks if error is an IO error
func IsIOError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Type == ErrorTypeIO
	}
	return false
}

// IsParseError checks if error is a parse error
func IsParseError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Type == ErrorTypeParse
	}
	return false
}

// Locate error helpers
func newLocateError(msg string, cause error) *Error {
	return &Error{
		Type:    ErrorTypeLocate,
		Message: fmt.Sprintf("failed locating steam dir: %s", msg),
		Cause:   cause,
	}
}

func newUnsupportedPlatformError() *Error {
	return newLocateError("unsupported platform", nil)
}

func newNoHomeError() *Error {
	return newLocateError("unable to locate user's home directory", nil)
}

func newWinRegError(cause error) *Error {
	return newLocateError("registry access failed", cause)
}

// Validation error helpers
func newValidationError(msg string) *Error {
	return &Error{
		Type:    ErrorTypeValidation,
		Message: fmt.Sprintf("failed validating steam dir: %s", msg),
	}
}

func newMissingDirError() *Error {
	return newValidationError("the steam directory either isn't a directory or doesn't exist")
}

// IO error helpers
func newIOError(path string, cause error) *Error {
	return &Error{
		Type:    ErrorTypeIO,
		Message: fmt.Sprintf("I/O error at %s: %v", path, cause),
		Path:    path,
		Cause:   cause,
	}
}

// Parse error helpers
func newParseError(kind ParseErrorKind, path string, cause error) *Error {
	return &Error{
		Type:    ErrorTypeParse,
		Message: fmt.Sprintf("failed parsing %s file at %s: %v", kind, path, cause),
		Path:    path,
		Cause:   cause,
	}
}

func newUnexpectedStructureError(kind ParseErrorKind, path string) *Error {
	return &Error{
		Type:    ErrorTypeParse,
		Message: fmt.Sprintf("file did not match expected structure: %s at %s", kind, path),
		Path:    path,
	}
}

func newMissingFileError(kind ParseErrorKind, path string) *Error {
	return &Error{
		Type:    ErrorTypeParse,
		Message: fmt.Sprintf("expected file was missing: %s at %s", kind, path),
		Path:    path,
	}
}

// MissingApp error
func newMissingAppError(appID uint32) *Error {
	return &Error{
		Type:    ErrorTypeMissingApp,
		Message: fmt.Sprintf("missing expected app with id: %d", appID),
	}
}

// IsNotExist checks if error is a "file not exist" error
func IsNotExist(err error) bool {
	if err == nil {
		return false
	}
	// Check if it's an IO error wrapping io.ErrNotExist or os.IsNotExist
	var ioErr *Error
	if As(err, &ioErr) && ioErr.Type == ErrorTypeIO {
		return Is(ioErr.Cause, io.EOF) || IsNotExist(ioErr.Cause)
	}
	return false
}

// Helper functions to avoid importing errors package directly everywhere
func Is(err, target error) bool {
	if err == nil || target == nil {
		return err == target
	}
	return err == target
}

func As(err error, target interface{}) bool {
	if err == nil || target == nil {
		return false
	}
	switch t := target.(type) {
	case **Error:
		if e, ok := err.(*Error); ok {
			*t = e
			return true
		}
	}
	return false
}

// filepath helpers for cross-platform compatibility
func joinPath(elem ...string) string {
	return filepath.Join(elem...)
}
