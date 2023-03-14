package ruleenginecore

import (
	"testing"
)

func Test_Set_NewSet(t *testing.T) {
	set := NewSet[int]()
	t.Run("Set Create operation", func(t *testing.T) {
		if len(set.internalMap) != 0 {
			t.Errorf("Set internalMap is not empty")
		}
	})
	t.Run("Set Create operation", func(t *testing.T) {
		if set.internalMap == nil {
			t.Errorf("Set internalMap is nil")
		}
	})
}

func Test_Set_Add(t *testing.T) {
	var element int = 10
	set := NewSet[int]()
	set.Add(element)
	t.Run("Set Add operation test", func(t *testing.T) {
		if !set.Contains(element) {
			t.Errorf("Set Add failed, could not find element")
		}
	})
}

func Test_Set_multi_Add(t *testing.T) {
	var element1 int = 10
	var element2 int = 20
	set := NewSet[int]()
	set.Add(element1)
	set.Add(element2)
	t.Run("Set Add operation test", func(t *testing.T) {
		if !set.Contains(element1) {
			t.Errorf("Set Add failed, could not find element")
		}
		if !set.Contains(element2) {
			t.Errorf("Set Add failed, could not find element")
		}
	})
}

func Test_Set_multi_Add_SameElement(t *testing.T) {
	var element1 int = 10
	set := NewSet[int]()
	set.Add(element1)
	set.Add(element1)
	t.Run("Set Add operation test", func(t *testing.T) {
		if !set.Contains(element1) {
			t.Errorf("Set Add failed, could not find element")
		}

		if set.Size() != 1 {
			t.Errorf("Set Add failed, adding same element multiple times resulted as multiple element in set")
		}
	})
}
func Test_Set_Remove(t *testing.T) {
	var element int = 10
	set := NewSet[int]()
	set.Add(element)

	set.Remove(element)
	t.Run("Set Remove operation check", func(t *testing.T) {
		if set.Contains(element) {
			t.Errorf("Set Remove failed, element exist after removal")
		}
	})
}

func Test_Set_Size(t *testing.T) {
	set := NewSet[int]()
	set.Add(10)
	var expectedSize = 1
	t.Run("Set Size operation check", func(t *testing.T) {
		if set.Size() != expectedSize {
			t.Errorf("Set Size failed, got %v, expected %v", set.Size(), expectedSize)
		}
	})
}

func Test_Set_Size2(t *testing.T) {
	set := NewSet[int]()
	set.Add(10)
	set.Add(20)
	var expectedSize = 2
	t.Run("Set Size operation check", func(t *testing.T) {
		if set.Size() != expectedSize {
			t.Errorf("Set Size failed, got %v, expected %v", set.Size(), expectedSize)
		}
	})
}

func Test_Set_Elements(t *testing.T) {
	var element1 int = 10
	var element2 int = 20
	set := NewSet[int]()
	set.Add(element1)
	set.Add(element2)

	var elements = set.Elements()

	t.Run("Set Elements operation check", func(t *testing.T) {
		if len(elements) != set.Size() {
			t.Errorf("Set Elements failed, got %v, expected %v", len(elements), set.Size())
		}

		for _, val := range elements {
			if val != element1 && val != element2 {
				t.Errorf("Set Elements failed, got %v invalid element", val)
			}
		}
	})
}
