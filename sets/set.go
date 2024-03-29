package sets

// Set is a set of elements
type Set[T comparable] map[T]struct{}

// NewSet returns a set of elements with assigned type
func NewSet[T comparable](es ...T) Set[T] {
	s := Set[T]{}
	for _, e := range es {
		s.Add(e)
	}
	return s
}

// Len report the elements number of s
func (s *Set[T]) Len() int {
	return len(*s)
}

// IsEmpty report wether s is empty
func (s *Set[T]) IsEmpty() bool {
	return s.Len() == 0
}

// Add add elements to set s
// if element is already in s this has no effect
func (s *Set[T]) Add(es ...T) {
	for _, e := range es {
		(*s)[e] = struct{}{}
	}
}

// Contains report wether v is in s
func (s *Set[T]) Contains(e T) bool {
	_, ok := (*s)[e]
	return ok
}

// Clone create a new set with the same elements as s
func (s *Set[T]) Clone() Set[T] {
	r := Set[T]{}
	r.Add(s.ToSlice()...)
	return r
}

// ToSlice transform set to slice
func (s *Set[T]) ToSlice() []T {
	r := make([]T, 0, s.Len())

	for e := range *s {
		r = append(r, e)
	}

	return r
}
