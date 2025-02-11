package main

import (
	"math/rand"
	"testing"

	"github.com/mlange-42/arche/ecs"
	"github.com/mlange-42/arche/generic"
)

// SpeciesParams component for species (relation targets).
type SpeciesParams struct {
	GrowthRate float64
}

// Species relation component for tree individuals.
type Species struct {
	ecs.Relation
}

// Biomass component for tree individuals.
type Biomass struct {
	BM float64
}

func TestMain(t *testing.T) {
	world := ecs.NewWorld()

	speciesBuilder := generic.NewMap1[SpeciesParams](&world)
	treeBuilder := generic.NewMap2[Biomass, Species](&world, generic.T[Species]())

	// Create 10 species.
	for s := 0; s < 10; s++ {
		species := speciesBuilder.NewWith(
			&SpeciesParams{GrowthRate: rand.Float64()},
		)

		// Create 100 trees per species. Biomass is zero.
		for t := 0; t < 100; t++ {
			treeBuilder.New(species)
		}
	}

	speciesFilter := generic.NewFilter1[SpeciesParams]()
	treeFilter := generic.NewFilter1[Biomass](). // We want to access biomass.
							With(generic.T[Species]()).        // We want this, but will not access it
							WithRelation(generic.T[Species]()) // Finally, the relation.

	// Time loop.
	for tick := 0; tick < 100; tick++ {
		// Query and iterate species.
		speciesQuery := speciesFilter.Query(&world)
		for speciesQuery.Next() {
			// Get species params and entity.
			params := speciesQuery.Get()
			species := speciesQuery.Entity()

			// Query and iterate trees for the current species.
			treeQuery := treeFilter.Query(&world, species)
			for treeQuery.Next() {
				bm := treeQuery.Get()
				// Increase biomass by the species' growth rate.
				bm.BM += params.GrowthRate
			}
		}
	}

	_ = world
}
