package geckecs

type Comp1Comp2Comp3Set struct {
	lastIdx int

	// owned components
	ownedComp1Store *SparseSet[Comp1]
	ownedComp2Store *SparseSet[Comp2]
	ownedComp3Store *SparseSet[Comp3]
}

func NewComp1Comp2Comp3Set(w *World) *Comp1Comp2Comp3Set {
	set := &Comp1Comp2Comp3Set{
		lastIdx: -1,

		ownedComp1Store: w.comp1Store,
		ownedComp2Store: w.comp2Store,
		ownedComp3Store: w.comp3Store,
	}
	return set
}

func (set *Comp1Comp2Comp3Set) PossibleUpdate(entities ...Entity) {
	for _, e := range entities {
		hasAllOwned := true

		if !set.ownedComp1Store.Has(e) {
			hasAllOwned = false
			break
		}

		if !set.ownedComp2Store.Has(e) {
			hasAllOwned = false
			break
		}

		if !set.ownedComp3Store.Has(e) {
			hasAllOwned = false
			break
		}

		sparseIdx := e.Index()

		if hasAllOwned {
			// swap with next after last
			set.lastIdx++

			wasSwapped := false

			if set.ownedComp1Store.ownedSetSwap(set.lastIdx, sparseIdx, false) {
				wasSwapped = true
			}

			if set.ownedComp2Store.ownedSetSwap(set.lastIdx, sparseIdx, false) {
				wasSwapped = true
			}

			if set.ownedComp3Store.ownedSetSwap(set.lastIdx, sparseIdx, false) {
				wasSwapped = true
			}

			if !wasSwapped {
				set.lastIdx--
			}
		} else {
			// swap with last
			wasSwapped := false

			if set.ownedComp1Store.ownedSetSwap(set.lastIdx, sparseIdx, true) {
				wasSwapped = true
			}

			if set.ownedComp2Store.ownedSetSwap(set.lastIdx, sparseIdx, true) {
				wasSwapped = true
			}

			if set.ownedComp3Store.ownedSetSwap(set.lastIdx, sparseIdx, true) {
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

func (set *Comp1Comp2Comp3Set) Len() int {
	return set.lastIdx + 1
}

func (set *Comp1Comp2Comp3Set) NewIterator() *Comp1Comp2Comp3SetIter {
	iter := &Comp1Comp2Comp3SetIter{set: set}
	iter.Reset()
	return iter
}

type Comp1Comp2Comp3SetIter struct {
	set     *Comp1Comp2Comp3Set
	current int
}

func (iter *Comp1Comp2Comp3SetIter) Reset() {
	iter.current = 0
}

func (iter *Comp1Comp2Comp3SetIter) HasNext() bool {
	return iter.current <= iter.set.lastIdx
}

func (iter *Comp1Comp2Comp3SetIter) Next() (
	Entity,
	Comp1,
	Comp2,
	Comp3,
) {
	e := iter.set.ownedComp1Store.dense[iter.current]
	comp0 := iter.set.ownedComp1Store.components[iter.current]
	comp1 := iter.set.ownedComp2Store.components[iter.current]
	comp2 := iter.set.ownedComp3Store.components[iter.current]
	iter.current++
	return e, comp0, comp1, comp2
}
