package geckecs

type Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet struct {
	lastIdx int

	// owned components
	ownedPositionsStore *SparseSet[Position]
	ownedComp1SStore    *SparseSet[Comp1]
	ownedComp2SStore    *SparseSet[Comp2]
	ownedComp3SStore    *SparseSet[Comp3]
	ownedComp4SStore    *SparseSet[Comp4]
	ownedComp5SStore    *SparseSet[Comp5]
	ownedComp6SStore    *SparseSet[Comp6]
	ownedComp7SStore    *SparseSet[Comp7]
	ownedComp8SStore    *SparseSet[Comp8]
	ownedComp9SStore    *SparseSet[Comp9]
}

func NewComp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet(w *World) *Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet {
	set := &Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet{
		lastIdx: -1,

		ownedPositionsStore: w.positionsStore,
		ownedComp1SStore:    w.comp1SStore,
		ownedComp2SStore:    w.comp2SStore,
		ownedComp3SStore:    w.comp3SStore,
		ownedComp4SStore:    w.comp4SStore,
		ownedComp5SStore:    w.comp5SStore,
		ownedComp6SStore:    w.comp6SStore,
		ownedComp7SStore:    w.comp7SStore,
		ownedComp8SStore:    w.comp8SStore,
		ownedComp9SStore:    w.comp9SStore,
	}
	return set
}

func (set *Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet) PossibleUpdate(entities ...Entity) {
	for _, e := range entities {
		hasAllOwned := true

		if !set.ownedPositionsStore.Has(e) {
			hasAllOwned = false
			break
		}

		if !set.ownedComp1SStore.Has(e) {
			hasAllOwned = false
			break
		}

		if !set.ownedComp2SStore.Has(e) {
			hasAllOwned = false
			break
		}

		if !set.ownedComp3SStore.Has(e) {
			hasAllOwned = false
			break
		}

		if !set.ownedComp4SStore.Has(e) {
			hasAllOwned = false
			break
		}

		if !set.ownedComp5SStore.Has(e) {
			hasAllOwned = false
			break
		}

		if !set.ownedComp6SStore.Has(e) {
			hasAllOwned = false
			break
		}

		if !set.ownedComp7SStore.Has(e) {
			hasAllOwned = false
			break
		}

		if !set.ownedComp8SStore.Has(e) {
			hasAllOwned = false
			break
		}

		if !set.ownedComp9SStore.Has(e) {
			hasAllOwned = false
			break
		}

		sparseIdx := e.Index()

		if hasAllOwned {
			// swap with next after last
			set.lastIdx++

			wasSwapped := false

			if set.ownedPositionsStore.ownedSetSwap(set.lastIdx, sparseIdx, false) {
				wasSwapped = true
			}

			if set.ownedComp1SStore.ownedSetSwap(set.lastIdx, sparseIdx, false) {
				wasSwapped = true
			}

			if set.ownedComp2SStore.ownedSetSwap(set.lastIdx, sparseIdx, false) {
				wasSwapped = true
			}

			if set.ownedComp3SStore.ownedSetSwap(set.lastIdx, sparseIdx, false) {
				wasSwapped = true
			}

			if set.ownedComp4SStore.ownedSetSwap(set.lastIdx, sparseIdx, false) {
				wasSwapped = true
			}

			if set.ownedComp5SStore.ownedSetSwap(set.lastIdx, sparseIdx, false) {
				wasSwapped = true
			}

			if set.ownedComp6SStore.ownedSetSwap(set.lastIdx, sparseIdx, false) {
				wasSwapped = true
			}

			if set.ownedComp7SStore.ownedSetSwap(set.lastIdx, sparseIdx, false) {
				wasSwapped = true
			}

			if set.ownedComp8SStore.ownedSetSwap(set.lastIdx, sparseIdx, false) {
				wasSwapped = true
			}

			if set.ownedComp9SStore.ownedSetSwap(set.lastIdx, sparseIdx, false) {
				wasSwapped = true
			}

			if !wasSwapped {
				set.lastIdx--
			}
		} else {
			// swap with last
			wasSwapped := false

			if set.ownedPositionsStore.ownedSetSwap(set.lastIdx, sparseIdx, true) {
				wasSwapped = true
			}

			if set.ownedComp1SStore.ownedSetSwap(set.lastIdx, sparseIdx, true) {
				wasSwapped = true
			}

			if set.ownedComp2SStore.ownedSetSwap(set.lastIdx, sparseIdx, true) {
				wasSwapped = true
			}

			if set.ownedComp3SStore.ownedSetSwap(set.lastIdx, sparseIdx, true) {
				wasSwapped = true
			}

			if set.ownedComp4SStore.ownedSetSwap(set.lastIdx, sparseIdx, true) {
				wasSwapped = true
			}

			if set.ownedComp5SStore.ownedSetSwap(set.lastIdx, sparseIdx, true) {
				wasSwapped = true
			}

			if set.ownedComp6SStore.ownedSetSwap(set.lastIdx, sparseIdx, true) {
				wasSwapped = true
			}

			if set.ownedComp7SStore.ownedSetSwap(set.lastIdx, sparseIdx, true) {
				wasSwapped = true
			}

			if set.ownedComp8SStore.ownedSetSwap(set.lastIdx, sparseIdx, true) {
				wasSwapped = true
			}

			if set.ownedComp9SStore.ownedSetSwap(set.lastIdx, sparseIdx, true) {
				wasSwapped = true
			}

			if wasSwapped {
				set.lastIdx--
			}
		}

		// do something with
		// hasAllBorrowed := true

	}
}

func (set *Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet) Len() int {
	return set.lastIdx + 1
}

func (set *Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet) NewIterator() *Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSetIter {
	iter := &Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSetIter{set: set}
	iter.Reset()
	return iter
}

type Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSetIter struct {
	set     *Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSet
	current int
}

func (iter *Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSetIter) Reset() {
	iter.current = 0
}

func (iter *Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSetIter) HasNext() bool {
	return iter.current <= iter.set.lastIdx
}

func (iter *Comp1Comp2Comp3Comp4Comp5Comp6Comp7Comp8Comp9PositionSetIter) Next() (
	Entity,
	Position,
	Comp1,
	Comp2,
	Comp3,
	Comp4,
	Comp5,
	Comp6,
	Comp7,
	Comp8,
	Comp9,
) {
	e := iter.set.ownedPositionsStore.dense[iter.current]
	comp0 := iter.set.ownedPositionsStore.components[iter.current]
	comp1 := iter.set.ownedComp1SStore.components[iter.current]
	comp2 := iter.set.ownedComp2SStore.components[iter.current]
	comp3 := iter.set.ownedComp3SStore.components[iter.current]
	comp4 := iter.set.ownedComp4SStore.components[iter.current]
	comp5 := iter.set.ownedComp5SStore.components[iter.current]
	comp6 := iter.set.ownedComp6SStore.components[iter.current]
	comp7 := iter.set.ownedComp7SStore.components[iter.current]
	comp8 := iter.set.ownedComp8SStore.components[iter.current]
	comp9 := iter.set.ownedComp9SStore.components[iter.current]
	iter.current++
	return e, comp0, comp1, comp2, comp3, comp4, comp5, comp6, comp7, comp8, comp9
}
