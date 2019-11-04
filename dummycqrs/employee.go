package simplecqrs

import (
	"errors"

	ycq "github.com/jetbasrawi/go.cqrs"
)

// Employee is the aggregate for an Employee item.
type Employee struct {
	*ycq.AggregateBase
	activated bool
	count     int
}

// NewEmployee constructs a new Employee item aggregate.
//
// Importantly it embeds a new AggregateBase.
func NewEmployee(id string) *Employee {
	i := &Employee{
		AggregateBase: ycq.NewAggregateBase(id),
	}

	return i
}

// Create raises EmployeeCreatedEvent
func (a *Employee) Create(name string) error {
	if name == "" {
		return errors.New("the name can not be empty")
	}

	a.Apply(ycq.NewEventMessage(a.AggregateID(),
		&EmployeeCreated{ID: a.AggregateID(), Name: name},
		ycq.Int(a.CurrentVersion())), true)

	return nil
}

// Remove removes items from Employee.
func (a *Employee) Remove(count int) error {
	if count <= 0 {
		return errors.New("can't remove negative count from Employee")
	}

	if a.count-count < 0 {
		return errors.New("can't remove more items from Employee than the number of items in Employee")
	}

	a.Apply(ycq.NewEventMessage(a.AggregateID(),
		&ItemsRemovedFromEmployee{ID: a.AggregateID(), Count: count},
		ycq.Int(a.CurrentVersion())), true)

	return nil
}

// CheckIn adds items to Employee.
func (a *Employee) CheckIn(count int) error {
	if count <= 0 {
		return errors.New("must have a count greater than 0 to add to Employee")
	}

	a.Apply(ycq.NewEventMessage(a.AggregateID(),
		&ItemsCheckedIntoEmployee{ID: a.AggregateID(), Count: count},
		ycq.Int(a.CurrentVersion())), true)

	return nil
}

// Apply handles the logic of events on the aggregate.
func (a *Employee) Apply(message ycq.EventMessage, isNew bool) {
	if isNew {
		a.TrackChange(message)
	}

	switch ev := message.Event().(type) {

	case *EmployeeCreated:
		a.activated = true

	case *EmployeeDeactivated:
		a.activated = false

	case *ItemsRemovedFromEmployee:
		a.count -= ev.Count

	case *ItemsCheckedIntoEmployee:
		a.count += ev.Count

	}

}
