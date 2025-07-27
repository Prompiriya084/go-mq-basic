package ports_repositories

type repository[Tentity any] interface {
	GetAll(filters *Tentity, preload []string) ([]Tentity, error)
	Get(filters *Tentity, preload []string) (*Tentity, error)
	Add(entity *Tentity) error
	Update(entity *Tentity) error
	Delete(entity *Tentity) error
}
