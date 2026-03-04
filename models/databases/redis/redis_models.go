package redis_models

import (
	"context"
	"time"
)

// --- Set (Insert) ---

type SetInput struct {
	Context    context.Context
	CancelFunc context.CancelFunc
	Key        string
	Value      any
	Expiration time.Duration // 0 = no expiration
}

// --- SetNX (Insert If Not Exists) ---

type SetNXInput struct {
	Context    context.Context
	CancelFunc context.CancelFunc
	Key        string
	Value      any
	Expiration time.Duration
}

type SetNXOutput struct {
	Success bool // true if key was set (did not already exist)
}

// --- Get ---

type GetInput struct {
	Context    context.Context
	CancelFunc context.CancelFunc
	Key        string
}

type GetOutput struct {
	Value string
	TTL   time.Duration // remaining time-to-live; -1 = no expiry, -2 = key missing
}

// --- Update ---

type UpdateInput struct {
	Context    context.Context
	CancelFunc context.CancelFunc
	Key        string
	Value      any
	Expiration time.Duration // new expiration; ignored when KeepTTL is true
	KeepTTL    bool          // preserve the existing TTL instead of setting a new one
}

type UpdateOutput struct {
	Updated bool // false when the key did not exist
}

// --- Delete ---

type DeleteInput struct {
	Context    context.Context
	CancelFunc context.CancelFunc
	Keys       []string
}

type DeleteOutput struct {
	DeletedCount int64
}

// --- Exists ---

type ExistsInput struct {
	Context    context.Context
	CancelFunc context.CancelFunc
	Keys       []string
}

type ExistsOutput struct {
	Count int64 // number of keys that exist
}

// --- Expire ---

type ExpireInput struct {
	Context    context.Context
	CancelFunc context.CancelFunc
	Key        string
	Expiration time.Duration
}

type ExpireOutput struct {
	Success bool // false if the key does not exist
}
