package generic

import "reflect"

// Comp is an alias for component types.
type Comp reflect.Type

// T provides a component type from a generic type argument.
func T[A any]() Comp {
	return Comp(typeOf[A]())
}

// T1 creates a component type list for one component types.
func T1[A any]() []Comp {
	return []Comp{typeOf[A]()}
}

// T2 creates a component type list for two component types.
func T2[A any, B any]() []Comp {
	return []Comp{typeOf[A](), typeOf[B]()}
}

// T3 creates a component type list for three component types.
func T3[A any, B any, C any]() []Comp {
	return []Comp{
		typeOf[A](), typeOf[B](), typeOf[C](),
	}
}

// T4 creates a component type list for four component types.
func T4[A any, B any, C any, D any]() []Comp {
	return []Comp{
		typeOf[A](), typeOf[B](), typeOf[C](), typeOf[D](),
	}
}

// T5 creates a component type list for five component types.
func T5[A any, B any, C any, D any, E any]() []Comp {
	return []Comp{
		typeOf[A](), typeOf[B](), typeOf[C](), typeOf[D](),
		typeOf[E](),
	}
}

// T6 creates a component type list for six component types.
func T6[A any, B any, C any, D any, E any, F any]() []Comp {
	return []Comp{
		typeOf[A](), typeOf[B](), typeOf[C](), typeOf[D](),
		typeOf[E](), typeOf[F](),
	}
}

// T7 creates a component type list for seven component types.
func T7[A any, B any, C any, D any, E any, F any, G any]() []Comp {
	return []Comp{
		typeOf[A](), typeOf[B](), typeOf[C](), typeOf[D](),
		typeOf[E](), typeOf[F](), typeOf[G](),
	}
}

// T8 creates a component type list for eight component types.
func T8[A any, B any, C any, D any, E any, F any, G any, H any]() []Comp {
	return []Comp{
		typeOf[A](), typeOf[B](), typeOf[C](), typeOf[D](),
		typeOf[E](), typeOf[F](), typeOf[G](), typeOf[H](),
	}
}

// T9 creates a component type list for nine component types.
func T9[A any, B any, C any, D any, E any, F any, G any, H any, I any]() []Comp {
	return []Comp{
		typeOf[A](), typeOf[B](), typeOf[C](), typeOf[D](),
		typeOf[E](), typeOf[F](), typeOf[G](), typeOf[H](),
		typeOf[I](),
	}
}

// T10 creates a component type list for ten component types.
func T10[A any, B any, C any, D any, E any, F any, G any, H any, I any, J any]() []Comp {
	return []Comp{
		typeOf[A](), typeOf[B](), typeOf[C](), typeOf[D](),
		typeOf[E](), typeOf[F](), typeOf[G](), typeOf[H](),
		typeOf[I](), typeOf[J](),
	}
}

// T11 creates a component type list for eleven component types.
func T11[A any, B any, C any, D any, E any, F any, G any, H any, I any, J any, K any]() []Comp {
	return []Comp{
		typeOf[A](), typeOf[B](), typeOf[C](), typeOf[D](),
		typeOf[E](), typeOf[F](), typeOf[G](), typeOf[H](),
		typeOf[I](), typeOf[J](), typeOf[K](),
	}
}

// T12 creates a component type list for twelve component types.
func T12[A any, B any, C any, D any, E any, F any, G any, H any, I any, J any, K any, L any]() []Comp {
	return []Comp{
		typeOf[A](), typeOf[B](), typeOf[C](), typeOf[D](),
		typeOf[E](), typeOf[F](), typeOf[G](), typeOf[H](),
		typeOf[I](), typeOf[J](), typeOf[K](), typeOf[L](),
	}
}
