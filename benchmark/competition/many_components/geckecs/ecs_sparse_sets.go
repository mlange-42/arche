package geckecs

type SparseSetSortFunc[T any] func(a, b T) bool

const sparseSetPageSize = 1024

type SparseSet[T any] struct {
	zero     T
	LessThan LessThan[T]

	sparsePages [][]int  // maps entity index to dense index
	dense       []Entity // contains the actual entities
	components  []T      // contains the actual component data

	// cache
	pagesToPossibleClear map[int]struct{}
}

func NewSparseSet[T any](lessThan LessThan[T]) *SparseSet[T] {
	return &SparseSet[T]{
		LessThan:             lessThan,
		pagesToPossibleClear: map[int]struct{}{},
	}
}

func (ss *SparseSet[T]) Has(e Entity) bool {
	_, _, _, storedEntity := ss.lookup(e)
	return storedEntity == e
}

func (ss *SparseSet[T]) Set(c T, entities ...Entity) error {
	for _, e := range entities {
		id := e.Index()
		sparsePageIdx := id / sparseSetPageSize
		sparseIdx := id % sparseSetPageSize

		// grow the sparse array if necessary
		for sparsePageIdx >= len(ss.sparsePages) {
			ss.sparsePages = append(ss.sparsePages, nil)
		}

		// fill at page if necessary
		sparsePage := ss.sparsePages[sparsePageIdx]
		if sparsePage == nil {
			sparsePage = make([]int, sparseSetPageSize)
			for i := range sparsePage {
				sparsePage[i] = DeadEntityID
			}
			ss.sparsePages[sparsePageIdx] = sparsePage
		}

		// check if the entity is already in the set
		denseIdx := ss.sparsePages[sparsePageIdx][sparseIdx]
		if denseIdx < len(ss.dense) {
			storedEntity := ss.dense[denseIdx]
			if storedEntity == e {
				ss.components[denseIdx] = c
				return nil
			}

			// the entity is not in the set, but the slot is taken
			return ErrEntityVersionMismatch
		}

		pos := len(ss.dense)
		sparsePage[sparseIdx] = pos
		ss.dense = append(ss.dense, e)
		ss.components = append(ss.components, c)
	}
	return nil
}

func (ss *SparseSet[T]) Read(e Entity) (c T, versionMatched bool) {
	_, _, denseIdx, storedEntity := ss.lookup(e)
	if storedEntity != e {
		// wrong version
		return ss.zero, false
	}
	return ss.components[denseIdx], true
}

func (ss *SparseSet[T]) Writeable(e Entity) (c *T, versionMatched bool) {
	_, _, denseIdx, storedEntity := ss.lookup(e)
	if storedEntity != e {
		// wrong version
		return nil, false
	}
	return &ss.components[denseIdx], true
}

func (ss *SparseSet[T]) Remove(entities ...Entity) {
	count := ss.Len()
	for i, e := range entities {
		removedSparsePageIdx, removedSparsePageOffset, denseIdx, storedEntity := ss.lookup(e)
		if storedEntity != e {
			return
		}

		// remove the entity from the set by moving the last element to its place
		lastIdx := count - i - 1
		lastEntity := ss.dense[lastIdx]

		lastEntityIdx := lastEntity.Index()
		lastEntityPageIdx := lastEntityIdx / sparseSetPageSize
		lastEntityPageOffset := lastEntityIdx % sparseSetPageSize

		// swap the last entity with the removed entity
		ss.sparsePages[lastEntityPageIdx][lastEntityPageOffset] = denseIdx
		ss.dense[denseIdx] = lastEntity
		ss.components[denseIdx] = ss.components[lastIdx]
		ss.dense = ss.dense[:lastIdx]
		ss.components = ss.components[:lastIdx]

		// clear the slot of the last entity
		pageRemovedFrom := ss.sparsePages[removedSparsePageIdx]
		pageRemovedFrom[removedSparsePageOffset] = DeadEntityID

		ss.pagesToPossibleClear[removedSparsePageIdx] = struct{}{}
	}

	for idx := range ss.pagesToPossibleClear {
		page := ss.sparsePages[idx]
		// check if the page is empty
		allDead := false
		for i := 0; i < sparseSetPageSize; i++ {
			if page[i] != DeadEntityID {
				allDead = false
				break
			}
		}
		if allDead {
			ss.sparsePages[idx] = nil
		}
	}
	clear(ss.pagesToPossibleClear)
}

func (ss *SparseSet[T]) ownedSetSwap(setIdx, sparseIdxToPossiblySwap int, isRemoval bool) (wasSwapped bool) {
	sparseToPossiblySwapPageIdx := sparseIdxToPossiblySwap / sparseSetPageSize
	sparseToPossiblySwapPageOffset := sparseIdxToPossiblySwap % sparseSetPageSize

	denseIdxToPossiblySwap := ss.sparsePages[sparseToPossiblySwapPageIdx][sparseToPossiblySwapPageOffset]
	if isRemoval {
		if denseIdxToPossiblySwap > setIdx {
			return false
		}
	} else {
		if denseIdxToPossiblySwap < setIdx {
			return false
		}
	}
	denseEntityToSwap := denseIdxToPossiblySwap

	denseSetEntity := ss.dense[setIdx]
	sparseIdx := denseSetEntity.Index()
	sparsePageIdx := sparseIdx / sparseSetPageSize
	sparsePageOffset := sparseIdx % sparseSetPageSize

	// swap the dense array
	ss.dense[setIdx], ss.dense[denseEntityToSwap] = ss.dense[denseEntityToSwap], ss.dense[setIdx]

	// update the sparse array
	ss.sparsePages[sparsePageIdx][sparsePageOffset], ss.sparsePages[sparseToPossiblySwapPageIdx][sparseToPossiblySwapPageOffset] = denseEntityToSwap, setIdx

	// update the component array
	ss.components[setIdx], ss.components[denseEntityToSwap] = ss.components[denseEntityToSwap], ss.components[setIdx]

	return true
}

func (ss *SparseSet[T]) IterateUnsafe(fn func(e Entity, c T)) {
	for i, e := range ss.dense {
		fn(e, ss.components[i])
	}
}

func (ss *SparseSet[T]) IterateSafe(fn func(e Entity, c T) error) error {
	for i := len(ss.dense) - 1; i >= 0; i-- {
		e := ss.dense[i]
		c := ss.components[i]
		if err := fn(e, c); err != nil {
			return err
		}
	}
	return nil
}

func (ss *SparseSet[T]) Len() int {
	return len(ss.dense)
}

func (ss *SparseSet[T]) Clear() {
	ss.sparsePages = ss.sparsePages[:0]
	ss.dense = ss.dense[:0]
	ss.components = ss.components[:0]
}

func (ss *SparseSet[T]) entities() []Entity {
	return ss.dense
}

func (ss *SparseSet[T]) lookup(e Entity) (sparsePageIdx, sparsePageOffset, denseIdx int, storedEntity Entity) {
	sparseIdx := e.Index()
	sparsePageIdx = sparseIdx / sparseSetPageSize
	if sparsePageIdx >= len(ss.sparsePages) {
		return -1, -1, -1, e.w.deadEntity
	}
	sparsePage := ss.sparsePages[sparsePageIdx]
	if sparsePage == nil {
		return -1, -1, -1, e.w.deadEntity
	}
	sparsePageOffset = sparseIdx % sparseSetPageSize

	denseIdx = sparsePage[sparsePageOffset]
	// too small or dead entity
	if denseIdx >= len(ss.dense) || denseIdx == e.w.deadEntity.Index() {
		return -1, -1, -1, e.w.deadEntity
	}
	storedEntity = ss.dense[denseIdx]
	return sparsePageIdx, sparsePageOffset, denseIdx, storedEntity
}
