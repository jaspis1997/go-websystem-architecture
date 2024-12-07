package logics

import (
	"playground/internal/entities"
	"sync"
)

var (
	once  sync.Once
	mutex sync.Mutex
)

var (
	Log entities.Logger
)

type logic struct {
	Name string
	Init func(repo *entities.Repository) error
}

var (
	logics     []logic
	repository *entities.Repository
)

func RegisterLogics(name string, init func(repo *entities.Repository) error) {
	mutex.Lock()
	defer mutex.Unlock()
	logics = append(logics, logic{Name: name, Init: init})
}

func Initialize(repo *entities.Repository) error {
	if repo == nil {
		panic("Initialize: repo == nil")
	}
	var (
		initRun   bool
		initError error
	)
	once.Do(func() {
		initRun = true
		for _, logic := range logics {
			err := logic.Init(repo)
			if err == nil {
				continue
			}
			e, ok := err.(*entities.StackTraceableError)
			if !ok {
				initError = entities.NewStackTrace().
					SetIdAutomatically().
					Add(err).
					Add(entities.ErrorInitializeFailed)
				return
			}
			initError = e.Add(entities.ErrorInitializeFailed)
			return
		}
		repository = repo
	})
	if Log != nil && initRun {
		Log.Info("logics initialized")
	}
	return initError
}

var (
	waitGroup sync.WaitGroup
)

func Done() {
	waitGroup.Wait()
}
