package unittest_repositories

type MockRepositoryImpl[Tentity any] struct {
	GetAllFn func(*Tentity, []string) ([]*Tentity, error)
	GetFn    func(*Tentity, []string) (*Tentity, error)
	AddFn    func(*Tentity) error
	UpdateFn func(*Tentity) error
	DeleteFn func(*Tentity) error
}

//	func newRepository[Tentity any]() *mockRepositoryImpl[Tentity] {
//		return &mockRepositoryImpl[Tentity]{}
//	}
func (r *MockRepositoryImpl[Tentity]) GetAll(filters *Tentity, preload []string) ([]*Tentity, error) {
	var entities []*Tentity
	if r.GetAllFn != nil {
		return r.GetAllFn(filters, preload)
	}
	return entities, nil
}
func (r *MockRepositoryImpl[Tentity]) Get(filters *Tentity, preload []string) (*Tentity, error) {
	var entity *Tentity
	if r.GetFn != nil {
		return r.GetFn(filters, preload)
	}

	return entity, nil
}
func (r *MockRepositoryImpl[Tentity]) Add(entity *Tentity) error {
	if r.AddFn != nil {
		return r.AddFn(entity)
	}
	return nil
}
func (r *MockRepositoryImpl[Tentity]) Update(entity *Tentity) error {
	if r.UpdateFn != nil {
		return r.UpdateFn(entity)
	}

	return nil
}
func (r *MockRepositoryImpl[Tentity]) Delete(entity *Tentity) error {
	if r.DeleteFn != nil {
		return r.DeleteFn(entity)
	}
	return nil
}
