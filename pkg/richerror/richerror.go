package richerror

type Kind int
type Op string

const (
	KindInvalid Kind = iota + 1
	KindForbidden
	KindNotFound
	KindUnexpected
)

type RichError struct {
	wrappedError error
	operation    Op
	message      string
	kind         Kind
	meta         map[string]interface{}
}

func New(op Op) RichError {
	return RichError{operation: op}
}

func (r RichError) WithMessage(message string) RichError {
	r.message = message
	return r
}

func (r RichError) WithOp(op Op) RichError {
	r.operation = op
	return r

}

func (r RichError) WithKind(kind Kind) RichError {
	r.kind = kind
	return r
}

func (r RichError) WithErr(err error) RichError {
	r.wrappedError = err
	return r
}

func (r RichError) WithMeta(meta map[string]interface{}) RichError {
	r.meta = meta
	return r
}

func (r RichError) Error() string {
	return r.message
}

// > recursive fn
func (r RichError) Kind() Kind {
	if r.kind != 0 {
		return r.kind
	}

	re, ok := r.wrappedError.(RichError)
	if !ok {
		return 0
	}

	return re.Kind()
}

func (r RichError) Message() string {
	if r.message != "" {
		return r.message
	}

	re, ok := r.wrappedError.(RichError)
	if !ok {
		return r.wrappedError.Error()
	}

	return re.Message()
}

//! two other ways to define constructor

// func New(args ...interface{}) RichError {
// 	r := RichError{}

// 	for _, arg := range args {
// 		switch arg.(type) {
// 		case string:
// 			r.message = arg.(string)
// 		case Op:
// 			r.operation = arg.(Op)
// 		case error:
// 			r.wrappedError = arg.(error)
// 		case Kind:
// 			r.kind = arg.(Kind)
// 		case map[string]interface{}:
// 			r.meta = arg.(map[string]interface{})
// 		}
// 	}

// 	return r
// }

//!
// func New(err error, operation, message string, kind Kind, meta map[string]interface{}) RichError {
// 	return RichError{wrappedError: err, operation: operation,
// 		message: message, kind: kind, meta: meta}
// }
