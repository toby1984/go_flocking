package list

import "testing"

func TestListIsEmpty(t *testing.T) {
	var list List[int]
	if !list.IsEmpty() {
		t.Error("list should be empty")
	}
}

type intslice []int

func (s intslice) contains(x int) bool {
	return contains(s, x)
}

func contains(slice []int, value int) bool {
	for _, a := range slice {
		if a == value {
			return true
		}
	}
	return false
}

func intList() List[int] {
	var list List[int]
	list.Equality = func(a int, b int) bool { return a == b }
	return list
}

func TestListAddOne(t *testing.T) {
	list := intList()
	list.Add(1)
	if list.IsEmpty() {
		t.Error("list should not be empty")
	}
	if list.Size() != 1 {
		t.Error("list should have size 1")
	}
	slice := list.ToSlice()
	if len(slice) != 1 {
		t.Error("Expected slice to have size 1")
	}
	if !contains(slice, 1) {
		t.Error("Slice should contain 1")
	}
	if !list.Contains(1) {
		t.Error("Slice should contain 1")
	}
}

func TestListAddTwo(t *testing.T) {
	list := intList()
	list.Add(1)
	list.Add(2)
	if list.IsEmpty() {
		t.Error("list should not be empty")
	}
	if list.Size() != 2 {
		t.Error("list should have size 2")
	}
	var slice intslice = list.ToSlice()
	if len(slice) != 2 {
		t.Error("Expected slice to have size 2")
	}
	if !slice.contains(1) {
		t.Error("Slice should contain 1")
	}
	if !slice.contains(2) {
		t.Error("Slice should contain 2")
	}
	if !list.Contains(1) {

	}
}

func TestListAddAndRemove(t *testing.T) {
	list := intList()
	list.Add(1)
	list.Add(2)
	list.Add(3)
	if list.IsEmpty() {
		t.Error("list should not be empty")
	}
	if list.Size() != 3 {
		t.Error("list should have size 3")
	}
	var slice intslice = list.ToSlice()
	if len(slice) != 3 {
		t.Error("Expected slice to have size 3")
	}
	if !slice.contains(1) {
		t.Error("Slice should contain 1")
	}
	if !slice.contains(2) {
		t.Error("Slice should contain 2")
	}
	if !slice.contains(3) {
		t.Error("Slice should contain 3")
	}

	// remove
	list.RemoveValue(3)
	slice = list.ToSlice()
	if list.Contains(3) {
		t.Error("Failed to remove 3")
	}
	if list.Size() != 2 {
		t.Error("Size not updated ?")
	}
	list.RemoveValue(2)
	if list.Contains(2) {
		t.Error("Failed to remove 2")
	}
	if list.Size() != 1 {
		t.Error("Size not updated ?")
	}
	list.RemoveValue(1)
	if list.Contains(1) {
		t.Error("Failed to remove 1")
	}
	if list.Size() != 0 || !list.IsEmpty() {
		t.Error("Size not updated ?")
	}
}
