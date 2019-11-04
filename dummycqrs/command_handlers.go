package simplecqrs

import (
	"log"
	"reflect"

	ycq "github.com/jetbasrawi/go.cqrs"
)

type EmployeeRepository interface {
	Load(string, string) (*Employee, error)
	Save(ycq.AggregateRoot, *int) error
}

// EmployeeCommandHandlers provides methods for processing commands related
// to Employee items.
type EmployeeCommandHandlers struct {
	repo EmployeeRepository
}

// NewEmployeeCommandHandlers contructs a new EmployeeCommandHandlers
func NewEmployeeCommandHandlers(repo EmployeeRepository) *EmployeeCommandHandlers {
	return &EmployeeCommandHandlers{
		repo: repo,
	}
}

// Handle processes Employee item commands.
func (h *EmployeeCommandHandlers) Handle(message ycq.CommandMessage) error {

	var item *Employee

	switch cmd := message.Command().(type) {

	case *CreateEmployee:

		item = NewEmployee(message.AggregateID())
		if err := item.Create(cmd.Name); err != nil {
			return &ycq.ErrCommandExecution{Command: message, Reason: err.Error()}
		}
		return h.repo.Save(item, ycq.Int(item.OriginalVersion()))

		//debit amount
	case *RemoveItemsFromEmployee:

		item, _ = h.repo.Load(reflect.TypeOf(&Employee{}).Elem().Name(), message.AggregateID())
		item.Remove(cmd.Count)
		return h.repo.Save(item, ycq.Int(item.OriginalVersion()))

		//credit amount
	case *CheckInItemsToEmployee:

		item, _ = h.repo.Load(reflect.TypeOf(&Employee{}).Elem().Name(), message.AggregateID())
		item.CheckIn(cmd.Count)
		return h.repo.Save(item, ycq.Int(item.OriginalVersion()))

	default:
		log.Fatalf("EmployeeCommandHandlers has received a command that it is does not know how to handle, %#v", cmd)
	}

	return nil
}
