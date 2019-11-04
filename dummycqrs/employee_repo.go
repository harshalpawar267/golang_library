package simplecqrs

import (
	"fmt"
	"reflect"

	ycq "github.com/jetbasrawi/go.cqrs"
	goes "github.com/jetbasrawi/go.geteventstore"
)

// EmployeeRepo is a repository specialized for persistence of
// Employees.
//
// repo and a *Employee is returned from the repo.
type EmployeeRepo struct {
	repo *ycq.GetEventStoreCommonDomainRepo
}

// NewEmployeeRepo constructs a new EmployeeRepository.
func NewEmployeeRepo(eventStore *goes.Client, eventBus ycq.EventBus) (*EmployeeRepo, error) {

	r, err := ycq.NewCommonDomainRepository(eventStore, eventBus)
	if err != nil {
		return nil, err
	}

	ret := &EmployeeRepo{
		repo: r,
	}

	// An aggregate factory creates an aggregate instance given the name of an aggregate.
	aggregateFactory := ycq.NewDelegateAggregateFactory()
	aggregateFactory.RegisterDelegate(&Employee{},
		func(id string) ycq.AggregateRoot { return NewEmployee(id) })
	ret.repo.SetAggregateFactory(aggregateFactory)

	// argument and the second argument are concatenated with a hyphen.
	streamNameDelegate := ycq.NewDelegateStreamNamer()
	streamNameDelegate.RegisterDelegate(func(t string, id string) string {
		return t + "-" + id
	}, &Employee{})
	ret.repo.SetStreamNameDelegate(streamNameDelegate)

	// An event factory creates an instance of an event given the name of an event
	// as a string.
	eventFactory := ycq.NewDelegateEventFactory()
	eventFactory.RegisterDelegate(&EmployeeCreated{},
		func() interface{} { return &EmployeeCreated{} })
	eventFactory.RegisterDelegate(&ItemsRemovedFromEmployee{},
		func() interface{} { return &ItemsRemovedFromEmployee{} })
	eventFactory.RegisterDelegate(&ItemsCheckedIntoEmployee{},
		func() interface{} { return &ItemsCheckedIntoEmployee{} })
	ret.repo.SetEventFactory(eventFactory)

	return ret, nil
}

// Load loads events for an aggregate.
//
// Returns an *EmployeeAggregate.
func (r *EmployeeRepo) Load(aggregateType, id string) (*Employee, error) {
	ar, err := r.repo.Load(reflect.TypeOf(&Employee{}).Elem().Name(), id)
	if _, ok := err.(*ycq.ErrAggregateNotFound); ok {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if ret, ok := ar.(*Employee); ok {
		return ret, nil
	}

	return nil, fmt.Errorf("Could not cast aggregate returned to type of %s", reflect.TypeOf(&Employee{}).Elem().Name())
}

// Save persists an aggregate.
func (r *EmployeeRepo) Save(aggregate ycq.AggregateRoot, expectedVersion *int) error {
	return r.repo.Save(aggregate, expectedVersion)
}
