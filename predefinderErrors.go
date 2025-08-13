package errors

// errors...
var (
	ErrBadRequest           = New("bad request")           // HTTP 400
	ErrUnauthorized         = New("user unauthorized")     // HTTP 401
	ErrRegistrationRequired = New("registration required") // HTTP 401
	ErrPaymentError         = New("payment error")         // HTTP 402
	ErrForbiddenAction      = New("forbidden")             // HTTP 403
	ErrNotFound             = New("entity not found")      // HTTP 404
	ErrConflict             = New("conflict request")      // HTTP 409
	ErrPreconditionFailed   = New("precondition failed")   // HTTP 412
	ErrValidation           = New("validation failed")     // HTTP 422
	ErrInternalServerError  = New("internal server error") // HTTP 500
)
