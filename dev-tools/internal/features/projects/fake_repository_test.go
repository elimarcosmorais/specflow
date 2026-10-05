package projects

import (
	"context"
	"sort"
)

// fakeRepository é um Repository em memória. Guarda o que o Service grava, para os testes
// inspecionarem a forma guardada e as chaves. List não ordena: a ordem é do Service.
type fakeRepository struct {
	items  map[int64]Project
	nextID int64
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{items: map[int64]Project{}, nextID: 1}
}

func (f *fakeRepository) List(context.Context) ([]Project, error) {
	// Ordem por id descendente: garante que a ordem por nome seja mesmo do Service.
	all := f.all()
	sort.Slice(all, func(i, j int) bool { return all[i].ID > all[j].ID })
	return all, nil
}

func (f *fakeRepository) Get(_ context.Context, id int64) (Project, error) {
	p, ok := f.items[id]
	if !ok {
		return Project{}, ErrNotFound
	}
	return p, nil
}

func (f *fakeRepository) ExistsName(_ context.Context, nameKey string, ignoreID int64) (bool, error) {
	for _, p := range f.items {
		if p.ID != ignoreID && p.NameKey == nameKey {
			return true, nil
		}
	}
	return false, nil
}

func (f *fakeRepository) ExistsPath(_ context.Context, pathKey string, ignoreID int64) (bool, error) {
	for _, p := range f.items {
		if p.ID != ignoreID && p.PathKey == pathKey {
			return true, nil
		}
	}
	return false, nil
}

func (f *fakeRepository) Create(_ context.Context, p Project) (int64, error) {
	p.ID = f.nextID
	f.nextID++
	f.items[p.ID] = p
	return p.ID, nil
}

func (f *fakeRepository) Update(_ context.Context, p Project) error {
	if _, ok := f.items[p.ID]; !ok {
		return ErrNotFound
	}
	f.items[p.ID] = p
	return nil
}

func (f *fakeRepository) Delete(_ context.Context, id int64) error {
	if _, ok := f.items[id]; !ok {
		return ErrNotFound
	}
	delete(f.items, id)
	return nil
}

// all devolve o que está guardado, ordenado por id, para a inspeção nos testes.
func (f *fakeRepository) all() []Project {
	all := make([]Project, 0, len(f.items))
	for _, p := range f.items {
		all = append(all, p)
	}
	sort.Slice(all, func(i, j int) bool { return all[i].ID < all[j].ID })
	return all
}
