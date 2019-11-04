package simplecqrs

// CreateEmployee create a new Employee item
type CreateEmployee struct {
	Name string
}

// CheckInItemsToEmployee adds items to Employee
type CheckInItemsToEmployee struct {
	OriginalVersion int
	Count           int
}

// RemoveItemsFromEmployee removes items from Employee
type RemoveItemsFromEmployee struct {
	OriginalVersion int
	Count           int
}
