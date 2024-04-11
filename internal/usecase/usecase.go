package usecase

// load your adapters
type Usecase struct {
	storage Sqlc
}

func New(storage Sqlc) *Usecase {
	return &Usecase{
		storage: storage,
	}
}
