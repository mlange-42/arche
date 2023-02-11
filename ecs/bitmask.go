package ecs

import "internal/base"

// bitMask is a bitmask.
type bitMask = base.BitMask

// MaskTotalBits is the size of Mask in bits.
//
// It is the maximum number of component types that may exist in any [World].
const MaskTotalBits = 64
