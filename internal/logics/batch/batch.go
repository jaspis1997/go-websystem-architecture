package batch

import "sync"

type Batch struct {
	once sync.Once
}

func NewBatch() (batch *Batch) {
	return &Batch{}
}

func (b *Batch) Start() (err error) {
	if err := b.registerSchedule(); err != nil {
		return err
	}
	return
}

func (b *Batch) Stop() (err error) {
	return
}

func (b *Batch) registerSchedule() (err error) {
	b.once.Do(func() {

	})
	return
}
