package simplecqrs

// Events are just plain structs

// EmployeeCreated event
type EmployeeCreated struct {
	ID   string
	Name string
}

// EmployeeRenamed event
type EmployeeRenamed struct {
	ID      string
	NewName string
}

// EmployeeDeactivated event
type EmployeeDeactivated struct {
	ID string
}

// ItemsRemovedFromEmployee event
type ItemsRemovedFromEmployee struct {
	ID    string
	Count int
}

// ItemsCheckedIntoEmployee event
type ItemsCheckedIntoEmployee struct {
	ID    string
	Count int
}
