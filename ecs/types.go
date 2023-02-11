package ecs

import (
	"internal/base"
)

// eid is the entity identifier/index type
type eid = base.Eid

// ID is the component identifier type
type ID = base.ID

// Entity identifier.
// Holds an entity ID and it's generation for recycling.
//
// Entities should only be created via the [World], using [World.NewEntity].
type Entity = base.Entity

// bitMask is a bitmask.
type bitMask = base.BitMask

// MaskTotalBits is the size of Mask in bits.
//
// It is the maximum number of component types that may exist in any [World].
const MaskTotalBits = base.MaskTotalBits

// Mask is a mask for a combination of components.
type Mask = base.Mask

// Component is a Component ID/Component pointer pair
type Component base.Component
