package entities

import "io"

type Resource interface {
	io.Closer

	Migrate() error
}

type Repository struct {
	Resource
	// example
	// ResourceName interface {
	// 	Resource
	// 	Get(id string) Entity
	// 	Save(entity Entity) error
	// 	Delete(id string) (ok bool, err error)
	// }
	// use repo.ResourceName.Get(id)
}

// repository close
// errors
//   - failed to close: failed to close
func (r Repository) Close() error {
	for _, processor := range []struct {
		name string
		fn   func() error
	}{
		// {"repository", r.Resource.Close},
	} {
		if err := processor.fn(); err != nil {
			if v, ok := err.(*StackTraceableError); ok {
				return v.Addf(FormatErrorFailedMigrate, processor.name)
			}
			return NewStackTrace().
				SetIdAutomatically().
				Add(err).
				Addf(FormatErrorFailedMigrate, processor.name)
		}
	}
	return nil
}

//	repository migration
//
// errors
//   - failed to migrate: failed to migrate
func (r Repository) Migrate() error {
	for _, processor := range []struct {
		name string
		fn   func() error
	}{
		// {"repository", r.Resource.Migrate},
	} {
		if err := processor.fn(); err != nil {
			if v, ok := err.(*StackTraceableError); ok {
				return v.Addf(FormatErrorFailedMigrate, processor.name)
			}
			return NewStackTrace().
				SetIdAutomatically().
				Add(err).
				Addf(FormatErrorFailedMigrate, processor.name)
		}
	}
	return nil
}
