package generic

import "reflect"

// Mask1 creates a component type list for one component type.
func Mask1[A any]() []reflect.Type {
	return []reflect.Type{typeOf[A]()}
}

// Mask2 creates a component type list for two component types.
func Mask2[A any, B any]() []reflect.Type {
	return []reflect.Type{typeOf[A](), typeOf[B]()}
}

// Mask3 creates a component type list for three component types.
func Mask3[A any, B any, C any]() []reflect.Type {
	return []reflect.Type{
		typeOf[A](), typeOf[B](), typeOf[C](),
	}
}

// Mask4 creates a component type list for four component types.
func Mask4[A any, B any, C any, D any]() []reflect.Type {
	return []reflect.Type{
		typeOf[A](), typeOf[B](), typeOf[C](), typeOf[D](),
	}
}

// Mask5 creates a component type list for five component types.
func Mask5[A any, B any, C any, D any, E any]() []reflect.Type {
	return []reflect.Type{
		typeOf[A](), typeOf[B](), typeOf[C](), typeOf[D](),
		typeOf[E](),
	}
}

// Mask6 creates a component type list for six component types.
func Mask6[A any, B any, C any, D any, E any, F any]() []reflect.Type {
	return []reflect.Type{
		typeOf[A](), typeOf[B](), typeOf[C](), typeOf[D](),
		typeOf[E](), typeOf[F](),
	}
}

// Mask7 creates a component type list for seven component types.
func Mask7[A any, B any, C any, D any, E any, F any, G any]() []reflect.Type {
	return []reflect.Type{
		typeOf[A](), typeOf[B](), typeOf[C](), typeOf[D](),
		typeOf[E](), typeOf[F](), typeOf[G](),
	}
}

// Mask8 creates a component type list for eight component types.
func Mask8[A any, B any, C any, D any, E any, F any, G any, H any]() []reflect.Type {
	return []reflect.Type{
		typeOf[A](), typeOf[B](), typeOf[C](), typeOf[D](),
		typeOf[E](), typeOf[F](), typeOf[G](), typeOf[H](),
	}
}
