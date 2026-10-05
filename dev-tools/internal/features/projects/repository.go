package projects

import "context"

// Repository guarda os projetos. Os projetos entram e saem na forma guardada (Project.Path
// relativo ou absoluto, D-07); quem converte para o caminho absoluto é o Service.
//
// Contrato de erros, que toda implementação cumpre:
//   - Get, Update e Delete de um id inexistente devolvem ErrNotFound;
//   - Create e Update devolvem ErrNameTaken ou ErrPathTaken se a unicidade for violada (corrida
//     entre ExistsName/ExistsPath e a gravação);
//   - qualquer falha de infraestrutura devolve apperr.ErrDatabaseUnavailable, com a causa.
type Repository interface {
	// List devolve todos os projetos, em qualquer ordem.
	List(ctx context.Context) ([]Project, error)
	Get(ctx context.Context, id int64) (Project, error)
	// ExistsName diz se outro projeto, que não o de id ignoreID, tem essa NameKey. Na criação,
	// ignoreID é zero.
	ExistsName(ctx context.Context, nameKey string, ignoreID int64) (bool, error)
	// ExistsPath é o ExistsName para a PathKey.
	ExistsPath(ctx context.Context, pathKey string, ignoreID int64) (bool, error)
	// Create grava o projeto, ignora p.ID e devolve o id atribuído.
	Create(ctx context.Context, p Project) (int64, error)
	// Update substitui o projeto de id p.ID.
	Update(ctx context.Context, p Project) error
	Delete(ctx context.Context, id int64) error
}
