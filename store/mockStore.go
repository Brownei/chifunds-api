package store

func NewMockStore() Store {
	return Store{
		Users: &UserStore{},
	}
}
