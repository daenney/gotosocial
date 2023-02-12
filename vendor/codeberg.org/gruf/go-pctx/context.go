package pctx

import (
	"context"
	"sync"
)

// ctxkey is a package private context
// key type in order to protect access
// to directly fetching the persistctx.
type ctxkey string

// persistkey is the key under which
// the persistctx itself is accessible.
var persistkey = ctxkey("ctx")

// persistctx wraps a base context.Context in
// order to provide persistent key-value storage
// across both directions of a callstack.
type PersistCtx struct {
	// embedded base ctx.
	context.Context

	// persisted vals.
	vs []struct{ k, v any }
	mu sync.Mutex
}

// Get will wrap the provided context in a PersistCtx, or
// return an existing copy if found in the context stack.
func Get(ctx context.Context) *PersistCtx {
	if ctx == nil {
		// Ensure a non-nil ctx.
		ctx = context.Background()
	}

	// Look for an existing persistctx in stack.
	pctx, ok := ctx.Value(persistkey).(*PersistCtx)

	if !ok {
		// Alloc new persistctx.
		pctx = new(PersistCtx)
		pctx.Context = ctx
	}

	return pctx
}

// Persist will store the key-value pair in provided context. Either wrapping it in a persistent
// context that supports multiple value storage, or digging through the context stack for an existing
// persistent context into which it will store the key-value pair. This allows storing and accessing
// context values across multiple layers of a callstack, both forwards and backwards.
func Persist(ctx context.Context, key any, value any) context.Context {
	if ctx == nil {
		// Ensure a non-nil ctx.
		ctx = context.Background()
	}

	// Look for an existing persistctx in stack.
	pctx, ok := ctx.Value(persistkey).(*PersistCtx)

	if !ok {
		// Alloc new persistctx.
		pctx = new(PersistCtx)
		pctx.Context = ctx

		// Set return ctx.
		ctx = pctx
	}

	// Set kv in context.
	pctx.Set(key, value)

	return ctx
}

// Set will safely set the key-value pair in context "map".
func (ctx *PersistCtx) Set(key any, value any) {
	ctx.mu.Lock()

	// Check for existing key-value.
	for i := range ctx.vs {
		if key == ctx.vs[i].k {
			ctx.vs[i].v = value
			ctx.mu.Unlock()
			return
		}
	}

	// Append new key-value pair.
	ctx.vs = append(ctx.vs, struct {
		k, v any
	}{k: key, v: value})

	ctx.mu.Unlock()
}

// Value implements context.Context's Value().
func (ctx *PersistCtx) Value(key any) any {
	if key == persistkey {
		// Return ourselves.
		return ctx
	}

	var (
		val any
		ok  bool
	)

	// Check "map" for key.
	ctx.mu.Lock()
	for i := range ctx.vs {
		if key == ctx.vs[i].k {
			val = ctx.vs[i].v
			ok = true
			break
		}
	}
	ctx.mu.Unlock()

	if ok {
		// Return persisted value.
		return val
	}

	// Check base context for value.
	return ctx.Context.Value(key)
}
