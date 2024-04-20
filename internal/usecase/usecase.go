package usecase

// load your adapters
type Usecase struct {
	storage ClickhouseRepo
}

func New(storage ClickhouseRepo) *Usecase {
	return &Usecase{
		storage: storage,
	}
}
