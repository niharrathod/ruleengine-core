package ruleenginecore

func NewSet[E comparable]() set[E] {
	var result set[E]
	result.internalMap = make(map[E]bool)
	return result
}

type set[E comparable] struct {
	internalMap map[E]bool
}

func (s *set[E]) Add(element E) {
	s.internalMap[element] = true
}

func (s *set[E]) Remove(element E) {
	delete(s.internalMap, element)
}

func (s *set[E]) Contains(element E) bool {
	if _, ok := s.internalMap[element]; !ok {
		return false
	}
	return true
}

func (s *set[E]) Elements() []E {
	elements := make([]E, len(s.internalMap))

	var index int = 0
	for key := range s.internalMap {
		elements[index] = key
		index++
	}

	return elements
}

func (s *set[E]) Size() int {
	return len(s.internalMap)
}
