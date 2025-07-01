package unittest_repositories

type MockRepositoryImpl[Tentity any] struct {
}

//	func newRepository[Tentity any]() *mockRepositoryImpl[Tentity] {
//		return &mockRepositoryImpl[Tentity]{}
//	}
func (r *MockRepositoryImpl[Tentity]) GetAll(filters *Tentity, preload []string) ([]Tentity, error) {
	var entities []Tentity

	return entities, nil
}
func (r *MockRepositoryImpl[Tentity]) Get(filters *Tentity, preload []string) (*Tentity, error) {
	var entity *Tentity

	return entity, nil
}
func (r *MockRepositoryImpl[Tentity]) Add(entity *Tentity) error {
	return nil
}
func (r *MockRepositoryImpl[Tentity]) Update(entity *Tentity) error {

	return nil
}
func (r *MockRepositoryImpl[Tentity]) Delete(entity *Tentity) error {
	return nil
}
