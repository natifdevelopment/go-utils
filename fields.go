package utils

import (
	"github.com/natifdevelopment/go-types"
)

// Re-export common field types from go-types so services can import a single
// utils package for shared helpers and types.
type DateField = types.DateField
type YearMonthField = types.YearMonthField
type EncryptedField = types.EncryptedField
