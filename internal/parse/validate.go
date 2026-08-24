package parse

import (
	"fmt"
	"strings"

	"pat-batch/internal/stats"
)

// ValidationError describes an issue found during data validation.
type ValidationError struct {
	Row     int
	Column  string
	Message string
}

// ValidateMeasurements checks measurement data for common issues:
// missing values, duplicate readings, negative times, etc.
func ValidateMeasurements(ms []stats.Measurement) []ValidationError {
	var errors []ValidationError

	seen := make(map[string]int) // "batch|param|value" -> row
	for i, m := range ms {
		if m.Batch == "" {
			errors = append(errors, ValidationError{
				Row:     i + 2,
				Column:  "batch",
				Message: "empty batch identifier",
			})
		}
		if m.Parameter == "" {
			errors = append(errors, ValidationError{
				Row:     i + 2,
				Column:  "parameter",
				Message: "empty parameter name",
			})
		}

		// Check for exact duplicates.
		key := fmt.Sprintf("%s|%s|%f", m.Batch, m.Parameter, m.Value)
		if prevRow, ok := seen[key]; ok {
			errors = append(errors, ValidationError{
				Row:     i + 2,
				Column:  "value",
				Message: fmt.Sprintf("duplicate of row %d", prevRow),
			})
		}
		seen[key] = i + 2
	}
	return errors
}

// ValidateSpecs checks specification limits for consistency.
func ValidateSpecs(specs map[string]stats.Spec) []ValidationError {
	var errors []ValidationError
	for name, sp := range specs {
		if sp.Low >= sp.High {
			errors = append(errors, ValidationError{
				Column:  "spec:" + name,
				Message: fmt.Sprintf("low (%.3f) >= high (%.3f)", sp.Low, sp.High),
			})
		}
		if sp.Target < sp.Low || sp.Target > sp.High {
			errors = append(errors, ValidationError{
				Column:  "spec:" + name,
				Message: fmt.Sprintf("target (%.3f) outside [%.3f, %.3f]", sp.Target, sp.Low, sp.High),
			})
		}
	}
	return errors
}

// CrossValidate checks that all measured parameters have corresponding specs.
func CrossValidate(ms []stats.Measurement, specs map[string]stats.Spec) []ValidationError {
	var errors []ValidationError
	missing := make(map[string]bool)
	for _, m := range ms {
		if _, ok := specs[m.Parameter]; !ok {
			if !missing[m.Parameter] {
				missing[m.Parameter] = true
				errors = append(errors, ValidationError{
					Column:  "parameter",
					Message: fmt.Sprintf("parameter %q has no specification", m.Parameter),
				})
			}
		}
	}
	return errors
}

// FormatErrors returns a human-readable summary of validation errors.
func FormatErrors(errors []ValidationError) string {
	if len(errors) == 0 {
		return ""
	}
	var b strings.Builder
	for _, e := range errors {
		if e.Row > 0 {
			fmt.Fprintf(&b, "  row %d [%s]: %s\n", e.Row, e.Column, e.Message)
		} else {
			fmt.Fprintf(&b, "  [%s]: %s\n", e.Column, e.Message)
		}
	}
	return b.String()
}
