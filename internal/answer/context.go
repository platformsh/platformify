package answer

import (
	"context"
)

// key is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type key int

// answerKey is the key for answers values in Contexts. It is unexported.
// Clients use answer.NewContext and answer.FromContext instead of using this key directly.
var answerKey key

// ToContext returns a new Context that carries answers value.
func ToContext(ctx context.Context, answers *Answers) context.Context {
	return context.WithValue(ctx, answerKey, answers)
}

// FromContext returns the answers value stored in ctx, if any.
func FromContext(ctx context.Context) (*Answers, bool) {
	answers, ok := ctx.Value(answerKey).(*Answers)
	return answers, ok
}
