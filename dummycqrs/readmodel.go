package simplecqrs

import (
	"errors"
	"log"

	ycq "github.com/jetbasrawi/go.cqrs"
)

var bullShitDatabase *BullShitDatabase

// ReadModelFacade is an interface for the readmodel facade
type ReadModelFacade interface {
	GetEmployees() []*EmployeeListDto
	GetEmployeeDetails(uuid string) *EmployeeDetailsDto
}

// EmployeeDetailsDto holds details for an Employee item.
type EmployeeDetailsDto struct {
	ID           string
	Name         string
	CurrentCount int
	Version      int
}

// EmployeeListDto provides a lightweight lookup view of an Employee item
type EmployeeListDto struct {
	ID   string
	Name string
}

// ReadModel is an implementation of the ReadModelFacade interface.
//
// ReadModel provides an in memory read model.
type ReadModel struct {
}

// NewReadModel constructs a new read model
func NewReadModel() *ReadModel {
	if bullShitDatabase == nil {
		bullShitDatabase = NewBullShitDatabase()
	}

	return &ReadModel{}
}

// GetEmployees returns a slice of all Employee items
func (m *ReadModel) GetEmployees() []*EmployeeListDto {
	return bullShitDatabase.List
}

// GetEmployeeDetails gets an EmployeeDetailsDto by ID
func (m *ReadModel) GetEmployeeDetails(uuid string) *EmployeeDetailsDto {
	if i, ok := bullShitDatabase.Details[uuid]; ok {
		return i
	}
	return nil
}

// EmployeeListView handles messages related to Employee and builds an
// in memory read model of Employee item summaries in a list.
type EmployeeListView struct {
}

// NewEmployeeListView constructs a new EmployeeListView
func NewEmployeeListView() *EmployeeListView {
	if bullShitDatabase == nil {
		bullShitDatabase = NewBullShitDatabase()
	}

	return &EmployeeListView{}
}

// Handle processes events related to Employee and builds an in memory read model
func (v *EmployeeListView) Handle(message ycq.EventMessage) {

	switch event := message.Event().(type) {

	case *EmployeeCreated:

		bullShitDatabase.List = append(bullShitDatabase.List, &EmployeeListDto{
			ID:   message.AggregateID(),
			Name: event.Name,
		})

		// case *EmployeeRenamed:

		// 	for _, v := range bullShitDatabase.List {
		// 		if v.ID == message.AggregateID() {
		// 			v.Name = event.NewName
		// 			break
		// 		}
		// 	}

		// case *EmployeeDeactivated:
		// 	i := -1
		// 	for k, v := range bullShitDatabase.List {
		// 		if v.ID == message.AggregateID() {
		// 			i = k
		// 			break
		// 		}
		// 	}

		// if i >= 0 {
		// 	bullShitDatabase.List = append(
		// 		bullShitDatabase.List[:i],
		// 		bullShitDatabase.List[i+1:]...,
		// 	)
		// }
	}
}

// EmployeeDetailView handles messages related to Employee and builds an
// in memory read model of Employee item details.
type EmployeeDetailView struct {
}

// NewEmployeeDetailView constructs a new EmployeeDetailView
func NewEmployeeDetailView() *EmployeeDetailView {
	if bullShitDatabase == nil {
		bullShitDatabase = NewBullShitDatabase()
	}

	return &EmployeeDetailView{}
}

// Handle handles events and build the projection
func (v *EmployeeDetailView) Handle(message ycq.EventMessage) {

	switch event := message.Event().(type) {

	case *EmployeeCreated:

		bullShitDatabase.Details[message.AggregateID()] = &EmployeeDetailsDto{
			ID:      message.AggregateID(),
			Name:    event.Name,
			Version: 0,
		}

	case *ItemsRemovedFromEmployee:

		d, err := v.GetDetailsItem(message.AggregateID())
		if err != nil {
			log.Fatal(err)
		}
		d.CurrentCount -= event.Count

	case *ItemsCheckedIntoEmployee:

		d, err := v.GetDetailsItem(message.AggregateID())
		if err != nil {
			log.Fatal(err)
		}
		d.CurrentCount += event.Count

	}
}

// GetDetailsItem gets an EmployeeDetailsDto by ID
func (v *EmployeeDetailView) GetDetailsItem(id string) (*EmployeeDetailsDto, error) {

	d, ok := bullShitDatabase.Details[id]
	if !ok {
		return nil, errors.New("did not find the original Employee this shouldn't not happen")
	}

	return d, nil
}

// BullShitDatabase is a simple in memory repository
type BullShitDatabase struct {
	Details map[string]*EmployeeDetailsDto
	List    []*EmployeeListDto
}

// NewBullShitDatabase constructs a new BullShitDatabase
func NewBullShitDatabase() *BullShitDatabase {
	return &BullShitDatabase{
		Details: make(map[string]*EmployeeDetailsDto),
	}
}
