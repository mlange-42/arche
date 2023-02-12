package ecs

import (
	"github.com/mlange-42/arche/internal/base"
)

// Eid is the entity identifier/index type
type eid uint32

// ID is the component identifier type
type ID = base.ID

// bitMask is a bitmask.
type bitMask = base.BitMask

// maskTotalBits is the size of Mask in bits.
//
// It is the maximum number of component types that may exist in any [World].
const maskTotalBits = base.MaskTotalBits

// Component is a Component ID/Component pointer pair
type Component struct {
	ID
	Component interface{}
}
