package traceID

// Indirect context object from raw pointer.
//
// Fix "checkptr: pointer arithmetic result points to invalid allocation" error in race mode.
//go:noescape
func indirectCtx(_ uintptr) *Ctx
