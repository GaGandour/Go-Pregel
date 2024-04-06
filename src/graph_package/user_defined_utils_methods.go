package graph_package

func NewVertexIdSet() *VertexIdSet {
	return &VertexIdSet{elements: make(map[VertexIdType]struct{})}
}

func (set *VertexIdSet) Add(element VertexIdType) {
	set.elements[element] = struct{}{}
}

func (set *VertexIdSet) Remove(element VertexIdType) {
	delete(set.elements, element)
}

func (set *VertexIdSet) Contains(element VertexIdType) bool {
	_, ok := set.elements[element]
	return ok
}

func (set *VertexIdSet) Size() int {
	return len(set.elements)
}

func (set *VertexIdSet) ToSlice() []VertexIdType {
	slice := make([]VertexIdType, 0, len(set.elements))
	for element := range set.elements {
		slice = append(slice, element)
	}
	return slice
}

func VertexIdSetsAreEqual(set1 *VertexIdSet, set2 *VertexIdSet) bool {
	if set1.Size() != set2.Size() {
		return false
	}
	for element := range set1.elements {
		if !set2.Contains(element) {
			return false
		}
	}
	return true
}
