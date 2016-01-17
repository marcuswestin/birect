package birect

import "github.com/marcuswestin/go-errs"

type Error errs.Err
type ErrorInfo errs.Info

var (
	NewError  = errs.New
	WrapError = errs.Wrap

	DefaultPublicErrorMessage = "Oops! Something went wrong - please try again."
)
