package compute

import "fmt"

// ParsingError represents a syntax error caught at the parsing level.
type ParsingError struct {
	Message string
}

func (e *ParsingError) Error() string {
	return fmt.Sprintf("Parsing Error: %s", e.Message)
}

func newParsingError(format string, args ...interface{}) *ParsingError {
	return &ParsingError{fmt.Sprintf(format, args...)}
}
