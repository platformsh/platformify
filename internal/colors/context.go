package colors

import (
	"context"
	"io"
)

// key is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type key int

const (
	outKey key = iota
	errKey
)

// ToContext returns a new Context that carries answers value.
func ToContext(ctx context.Context, out, err io.Writer) context.Context {
	return context.WithValue(context.WithValue(ctx, outKey, out), errKey, err)
}

// FromContext returns the out and err writers stored in ctx, if any.
func FromContext(ctx context.Context) (out, err io.Writer, ok bool) {
	out, ok = ctx.Value(outKey).(io.Writer)
	if !ok {
		return nil, nil, false
	}
	err, ok = ctx.Value(errKey).(io.Writer)
	if !ok {
		return nil, nil, false
	}
	return out, err, true
}
