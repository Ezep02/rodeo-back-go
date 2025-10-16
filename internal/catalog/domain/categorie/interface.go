package categorie

import "context"

type CategorieRepository interface {
	Create(ctx context.Context, category *Categorie) error
	Update(ctx context.Context, id uint, data *Categorie) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context) ([]Categorie, error)
	GetByID(ctx context.Context, id uint) (*Categorie, error)
}
