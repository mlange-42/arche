package stats

import (
	"fmt"
	"reflect"
	"testing"
)

func TestStats(t *testing.T) {
	stats := WorldStats{
		Entities:       EntityStats{},
		ComponentCount: 1,
		ComponentTypes: []reflect.Type{reflect.TypeOf(1)},
		Locked:         false,
		Nodes: []NodeStats{
			{
				Size:           1,
				Capacity:       128,
				Components:     1,
				ComponentIDs:   []uint8{0},
				ComponentTypes: []reflect.Type{reflect.TypeOf(1)},
			},
			{
				IsActive:       true,
				Size:           1,
				Capacity:       128,
				Components:     1,
				ComponentIDs:   []uint8{0},
				ComponentTypes: []reflect.Type{reflect.TypeOf(1)},
			},
		},
	}
	fmt.Println(stats.String())

}
