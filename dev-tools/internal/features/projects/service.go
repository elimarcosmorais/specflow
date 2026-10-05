package projects

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

// Service aplica as regras de Projetos (RN-12 a RN-18). Só lê o disco (os.Stat): nunca cria,
// move, copia nem apaga arquivo ou pasta (RN-16).
type Service struct {
	repo        Repository
	installRoot string
	logger      *slog.Logger
	now         func() time.Time
}

// NewService cria o serviço. installRoot é a raiz da instalação (D-07); now dá as datas de criação
// e de alteração.
func NewService(repo Repository, installRoot string, logger *slog.Logger, now func() time.Time) *Service {
	return &Service{repo: repo, installRoot: filepath.Clean(installRoot), logger: logger, now: now}
}

// List devolve os projetos ordenados por nome, sem diferenciar maiúsculas, com o caminho absoluto
// e FolderExists calculado agora (RN-18).
func (s *Service) List(ctx context.Context) ([]View, error) {
	all, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	slices.SortFunc(all, func(a, b Project) int { return strings.Compare(a.NameKey, b.NameKey) })

	views := make([]View, 0, len(all))
	for _, p := range all {
		views = append(views, s.view(p))
	}
	return views, nil
}

// Create valida a entrada (RN-12 a RN-15), grava o projeto e devolve a View.
func (s *Service) Create(ctx context.Context, in Input) (View, error) {
	p, err := s.build(ctx, in, 0)
	if err != nil {
		return View{}, err
	}
	p.CreatedAt = s.now().UTC()
	p.UpdatedAt = p.CreatedAt

	id, err := s.repo.Create(ctx, p)
	if err != nil {
		return View{}, err
	}
	p.ID = id

	v := s.view(p)
	s.logger.Info("project.created", "id", v.ID, "name", v.Name, "path", v.Path)
	return v, nil
}

// Update substitui os dados do projeto id, com as mesmas regras do Create. A data de criação é
// mantida. O próprio projeto é ignorado na checagem de unicidade (RN-13).
func (s *Service) Update(ctx context.Context, id int64, in Input) (View, error) {
	current, err := s.repo.Get(ctx, id)
	if err != nil {
		return View{}, err
	}
	p, err := s.build(ctx, in, id)
	if err != nil {
		return View{}, err
	}
	p.ID = id
	p.CreatedAt = current.CreatedAt
	p.UpdatedAt = s.now().UTC()

	if err := s.repo.Update(ctx, p); err != nil {
		return View{}, err
	}

	v := s.view(p)
	s.logger.Info("project.updated", "id", v.ID, "name", v.Name, "path", v.Path)
	return v, nil
}

// Delete remove o projeto do repositório. A pasta no disco não é tocada (RN-16).
func (s *Service) Delete(ctx context.Context, id int64) error {
	p, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.logger.Info("project.deleted", "id", p.ID, "name", p.Name, "path", s.absolute(p.Path))
	return nil
}

// build valida a entrada na ordem da RN-12 e devolve o Project pronto para gravar, sem id nem
// datas. ignoreID é o projeto a ignorar na unicidade (zero na criação).
func (s *Service) build(ctx context.Context, in Input, ignoreID int64) (Project, error) {
	name := strings.TrimSpace(in.Name)
	path := strings.TrimSpace(in.Path)
	typ := Type(strings.TrimSpace(string(in.Type)))

	switch {
	case name == "":
		return Project{}, ErrNameRequired
	case path == "":
		return Project{}, ErrPathRequired
	case typ == "":
		return Project{}, ErrTypeRequired
	case !typ.valid():
		return Project{}, ErrTypeInvalid
	case !filepath.IsAbs(path):
		return Project{}, ErrPathInvalid
	}

	clean := filepath.Clean(path)
	if info, err := os.Stat(clean); err != nil || !info.IsDir() {
		return Project{}, ErrPathNotFound
	}

	stored := s.stored(clean)
	p := Project{
		Name:    name,
		NameKey: NameKey(name),
		Path:    stored,
		PathKey: PathKey(stored),
		Type:    typ,
	}

	taken, err := s.repo.ExistsName(ctx, p.NameKey, ignoreID)
	if err != nil {
		return Project{}, err
	}
	if taken {
		return Project{}, ErrNameTaken
	}
	taken, err = s.repo.ExistsPath(ctx, p.PathKey, ignoreID)
	if err != nil {
		return Project{}, err
	}
	if taken {
		return Project{}, ErrPathTaken
	}
	return p, nil
}

// stored devolve a forma guardada de um caminho absoluto e limpo (D-07): relativo à raiz da
// instalação se a pasta está dentro dela (sem diferenciar maiúsculas), senão o próprio absoluto.
// O relativo é o que filepath.Rel(installRoot, clean) devolveria; o prefixo é comparado à mão
// porque Rel só ignora maiúsculas no Windows.
func (s *Service) stored(clean string) string {
	prefix := strings.TrimSuffix(s.installRoot, string(filepath.Separator))
	n := len(prefix)
	if len(clean) > n+1 && strings.EqualFold(clean[:n], prefix) && clean[n] == filepath.Separator {
		return clean[n+1:]
	}
	return clean
}

// absolute é o inverso de stored.
func (s *Service) absolute(stored string) string {
	if filepath.IsAbs(stored) {
		return stored
	}
	return filepath.Join(s.installRoot, stored)
}

func (s *Service) view(p Project) View {
	abs := s.absolute(p.Path)
	info, err := os.Stat(abs)
	return View{
		ID:           p.ID,
		Name:         p.Name,
		Path:         abs,
		Type:         p.Type,
		FolderExists: err == nil && info.IsDir(),
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}
