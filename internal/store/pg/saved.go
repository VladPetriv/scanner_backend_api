package pg

type SavedRepo struct {
	db *DB
}

func NewSavedRepo(db *DB) *SavedRepo {
	return &SavedRepo{db: db}
}
