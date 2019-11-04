package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	dummycqrs "dummy_bank/dummycqrs"

	ycq "github.com/jetbasrawi/go.cqrs"
)

var (
	readModel  dummycqrs.ReadModelFacade
	dispatcher ycq.Dispatcher
)

var t = template.Must(template.ParseGlob("templates/*"))

func init() {

	readModel = dummycqrs.NewReadModel()

	listView := dummycqrs.NewEmployeeListView()

	detailView := dummycqrs.NewEmployeeDetailView()

	// Create an EventBus
	eventBus := ycq.NewInternalEventBus()
	// Register the listView as an event handler on the event bus
	// for the events specified.
	eventBus.AddHandler(listView,
		&dummycqrs.EmployeeCreated{},
	)
	// Register the detail view as an event handler on the event bus
	// for the events specified.
	eventBus.AddHandler(detailView,
		&dummycqrs.EmployeeCreated{},
		&dummycqrs.ItemsRemovedFromEmployee{},
		&dummycqrs.ItemsCheckedIntoEmployee{},
	)

	// Here we use an in memory event repository.
	repo := dummycqrs.NewInMemoryRepo(eventBus)

	EmployeeCommandHandler := dummycqrs.NewEmployeeCommandHandlers(repo)

	// Create a dispatcher
	dispatcher = ycq.NewInMemoryDispatcher()
	// Register the Employee command handlers instance as a command handler
	// for the events specified.
	err := dispatcher.RegisterHandler(EmployeeCommandHandler,
		&dummycqrs.CreateEmployee{},
		&dummycqrs.CheckInItemsToEmployee{},
		&dummycqrs.RemoveItemsFromEmployee{},
	)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	mux := setupHandlers()
	if err := http.ListenAndServe(":8088", mux); err != nil {
		log.Fatal(err)
	}

}

func setupHandlers() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Items := readModel.GetEmployees()

		err := t.ExecuteTemplate(w, "index", Items)
		if err != nil {
			log.Fatal(err)
		}
	})

	mux.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				log.Fatal(err)
			}

			id := ycq.NewUUID()
			em := ycq.NewCommandMessage(id, &dummycqrs.CreateEmployee{
				Name: r.Form.Get("name"),
			})

			err = dispatcher.Dispatch(em)
			if err != nil {
				log.Println(err)
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

		err := t.ExecuteTemplate(w, "add", nil)
		if err != nil {
			log.Fatal(err)
		}
	})

	mux.HandleFunc("/details/", func(w http.ResponseWriter, r *http.Request) {
		p := strings.Split(r.URL.Path, "/")
		id := p[len(p)-1]
		Item := readModel.GetEmployeeDetails(id)
		err := t.ExecuteTemplate(w, "details", Item)
		if err != nil {
			log.Fatal(err)
		}
	})

	mux.HandleFunc("/checkin/", func(w http.ResponseWriter, r *http.Request) {

		p := strings.Split(r.URL.Path, "/")
		id := p[len(p)-1]
		Item := readModel.GetEmployeeDetails(id)

		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				log.Fatal(err)
			}

			num, err := strconv.Atoi(r.Form.Get("number"))
			if err != nil {
				http.Error(w, "Unable to read number.", http.StatusInternalServerError)
			}

			em := ycq.NewCommandMessage(id, &dummycqrs.CheckInItemsToEmployee{Count: num})
			err = dispatcher.Dispatch(em)
			if err != nil {
				log.Println(err)
			}

			redirectURL := "/"
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		}

		err := t.ExecuteTemplate(w, "creditmoney", Item)
		if err != nil {
			log.Fatal(err)
		}

	})

	mux.HandleFunc("/remove/", func(w http.ResponseWriter, r *http.Request) {

		p := strings.Split(r.URL.Path, "/")
		id := p[len(p)-1]
		Item := readModel.GetEmployeeDetails(id)

		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				log.Fatal(err)
			}

			num, err := strconv.Atoi(r.Form.Get("number"))
			if err != nil {
				http.Error(w, "Unable to read number.", http.StatusInternalServerError)
			}

			em := ycq.NewCommandMessage(id, &dummycqrs.RemoveItemsFromEmployee{Count: num})
			err = dispatcher.Dispatch(em)
			if err != nil {
				log.Println(err)
			}
			redirectURL := "/"
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		}

		err := t.ExecuteTemplate(w, "debitmoney", Item)
		if err != nil {
			log.Fatal(err)
		}
	})

	mux.HandleFunc("/assets/", func(w http.ResponseWriter, r *http.Request) {
		staticFile := r.URL.Path[len("/assets/"):]
		if len(staticFile) != 0 {
			f, err := http.Dir("assets/").Open(staticFile)
			if err == nil {
				content := io.ReadSeeker(f)
				http.ServeContent(w, r, staticFile, time.Now(), content)
				return
			}
		}
		http.NotFound(w, r)
	})

	return mux
}
