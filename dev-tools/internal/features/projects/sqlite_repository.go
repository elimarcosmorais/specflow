package projects

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"modernc.org/sqlite"
	sqlitelib "modernc.org/sqlite/lib"

	"dev-tools/internal/apperr"
)

// errNoDatabase é a causa registrada quando o app subiu sem banco (D-06).
var errNoDatabase = errors.New("database is not available")

const projectColumns = "id, name, name_key, path, path_key, type, created_at, updated_at"

// SQLiteRepository guarda os projetos na tabela projects (spec §4.1). Cumpre o contrato de erros
// de Repository: unicidade violada vira ErrNameTaken ou ErrPathTaken, linha ausente vira
// ErrNotFound e qualquer outra falha do SQLite vira apperr.ErrDatabaseUnavailable com a causa.
type SQLiteRepository struct {
	db *sql.DB
}

// NewSQLiteRepository cria o repositório sobre db. Um db nulo (banco que não abriu) é aceito:
// toda operação devolve apperr.ErrDatabaseUnavailable.
func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

func (r *SQLiteRepository) List(ctx context.Context) ([]Project, error) {
	if r.db == nil {
		return nil, unavailable(errNoDatabase)
	}
	rows, err := r.db.QueryContext(ctx, "SELECT "+projectColumns+" FROM projects")
	if err != nil {
		return nil, mapError(err)
	}
	defer rows.Close()

	var all []Project
	for rows.Next() {
		p, err := scanProject(rows)
		if err != nil {
			return nil, mapError(err)
		}
		all = append(all, p)
	}
	if err := rows.Err(); err != nil {
		return nil, mapError(err)
	}
	return all, nil
}

func (r *SQLiteRepository) Get(ctx context.Context, id int64) (Project, error) {
	if r.db == nil {
		return Project{}, unavailable(errNoDatabase)
	}
	row := r.db.QueryRowContext(ctx, "SELECT "+projectColumns+" FROM projects WHERE id = ?", id)
	p, err := scanProject(row)
	if err != nil {
		return Project{}, mapError(err)
	}
	return p, nil
}

func (r *SQLiteRepository) ExistsName(ctx context.Context, nameKey string, ignoreID int64) (bool, error) {
	return r.exists(ctx, "name_key", nameKey, ignoreID)
}

func (r *SQLiteRepository) ExistsPath(ctx context.Context, pathKey string, ignoreID int64) (bool, error) {
	return r.exists(ctx, "path_key", pathKey, ignoreID)
}

// exists consulta a coluna de chave (valor fixo no código, nunca vindo da entrada).
func (r *SQLiteRepository) exists(ctx context.Context, keyColumn, key string, ignoreID int64) (bool, error) {
	if r.db == nil {
		return false, unavailable(errNoDatabase)
	}
	var found bool
	err := r.db.QueryRowContext(ctx,
		"SELECT EXISTS (SELECT 1 FROM projects WHERE "+keyColumn+" = ? AND id <> ?)", key, ignoreID,
	).Scan(&found)
	if err != nil {
		return false, mapError(err)
	}
	return found, nil
}

func (r *SQLiteRepository) Create(ctx context.Context, p Project) (int64, error) {
	if r.db == nil {
		return 0, unavailable(errNoDatabase)
	}
	res, err := r.db.ExecContext(ctx,
		"INSERT INTO projects (name, name_key, path, path_key, type, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		p.Name, p.NameKey, p.Path, p.PathKey, string(p.Type), formatTime(p.CreatedAt), formatTime(p.UpdatedAt))
	if err != nil {
		return 0, mapError(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, mapError(err)
	}
	return id, nil
}

// Update não altera created_at: a data de criação não muda numa edição.
func (r *SQLiteRepository) Update(ctx context.Context, p Project) error {
	if r.db == nil {
		return unavailable(errNoDatabase)
	}
	res, err := r.db.ExecContext(ctx,
		"UPDATE projects SET name = ?, name_key = ?, path = ?, path_key = ?, type = ?, updated_at = ? WHERE id = ?",
		p.Name, p.NameKey, p.Path, p.PathKey, string(p.Type), formatTime(p.UpdatedAt), p.ID)
	if err != nil {
		return mapError(err)
	}
	return requireAffected(res)
}

func (r *SQLiteRepository) Delete(ctx context.Context, id int64) error {
	if r.db == nil {
		return unavailable(errNoDatabase)
	}
	res, err := r.db.ExecContext(ctx, "DELETE FROM projects WHERE id = ?", id)
	if err != nil {
		return mapError(err)
	}
	return requireAffected(res)
}

// requireAffected devolve ErrNotFound quando o comando não atingiu nenhuma linha.
func requireAffected(res sql.Result) error {
	n, err := res.RowsAffected()
	if err != nil {
		return mapError(err)
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// scanner é o que *sql.Row e *sql.Rows têm em comum.
type scanner interface {
	Scan(dest ...any) error
}

func scanProject(s scanner) (Project, error) {
	var (
		p                    Project
		typ                  string
		createdAt, updatedAt string
	)
	if err := s.Scan(&p.ID, &p.Name, &p.NameKey, &p.Path, &p.PathKey, &typ, &createdAt, &updatedAt); err != nil {
		return Project{}, err
	}
	p.Type = Type(typ)

	var err error
	if p.CreatedAt, err = time.Parse(time.RFC3339, createdAt); err != nil {
		return Project{}, err
	}
	if p.UpdatedAt, err = time.Parse(time.RFC3339, updatedAt); err != nil {
		return Project{}, err
	}
	return p, nil
}

// formatTime grava a data em RFC 3339, UTC. Mantém as frações de segundo, para ler de volta
// exatamente o que foi gravado.
func formatTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339Nano)
}

// mapError traduz o erro do SQLite para o contrato do Repository (spec §5, D-06).
func mapError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	var sqliteErr *sqlite.Error
	if errors.As(err, &sqliteErr) && sqliteErr.Code() == sqlitelib.SQLITE_CONSTRAINT_UNIQUE {
		switch msg := sqliteErr.Error(); {
		case strings.Contains(msg, "projects.name_key"):
			return ErrNameTaken
		case strings.Contains(msg, "projects.path_key"):
			return ErrPathTaken
		}
	}
	return unavailable(err)
}

func unavailable(cause error) error {
	return apperr.ErrDatabaseUnavailable.Wrap(cause)
}
