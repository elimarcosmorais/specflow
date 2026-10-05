package projects

import (
	"net/http"
	"strconv"
	"time"

	"dev-tools/internal/platform/httpserver"
)

// projectDTO é o Project da API (spec §3.1).
type projectDTO struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Path         string `json:"path"`
	Type         string `json:"type"`
	FolderExists bool   `json:"folderExists"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

// listDTO é a resposta de GET /api/projects.
type listDTO struct {
	Items []projectDTO `json:"items"`
}

// inputDTO é o corpo de POST e PUT. Campo ausente no JSON fica como string vazia.
type inputDTO struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"`
}

func (in inputDTO) input() Input {
	return Input{Name: in.Name, Path: in.Path, Type: Type(in.Type)}
}

func toDTO(v View) projectDTO {
	return projectDTO{
		ID:           v.ID,
		Name:         v.Name,
		Path:         v.Path,
		Type:         string(v.Type),
		FolderExists: v.FolderExists,
		CreatedAt:    v.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:    v.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

// Routes devolve as rotas da feature: GET e POST /api/projects, PUT e DELETE
// /api/projects/{id}, todas protegidas por token.
func Routes(s *Service) []httpserver.Route {
	return []httpserver.Route{
		{Method: http.MethodGet, Pattern: "/api/projects", Handler: func(w http.ResponseWriter, r *http.Request) error {
			views, err := s.List(r.Context())
			if err != nil {
				return err
			}
			items := make([]projectDTO, 0, len(views))
			for _, v := range views {
				items = append(items, toDTO(v))
			}
			httpserver.WriteJSON(w, http.StatusOK, listDTO{Items: items})
			return nil
		}},
		{Method: http.MethodPost, Pattern: "/api/projects", Handler: func(w http.ResponseWriter, r *http.Request) error {
			var in inputDTO
			if err := httpserver.DecodeJSON(r, &in); err != nil {
				return err
			}
			v, err := s.Create(r.Context(), in.input())
			if err != nil {
				return err
			}
			httpserver.WriteJSON(w, http.StatusCreated, toDTO(v))
			return nil
		}},
		{Method: http.MethodPut, Pattern: "/api/projects/{id}", Handler: func(w http.ResponseWriter, r *http.Request) error {
			id, err := parseID(r)
			if err != nil {
				return err
			}
			var in inputDTO
			if err := httpserver.DecodeJSON(r, &in); err != nil {
				return err
			}
			v, err := s.Update(r.Context(), id, in.input())
			if err != nil {
				return err
			}
			httpserver.WriteJSON(w, http.StatusOK, toDTO(v))
			return nil
		}},
		{Method: http.MethodDelete, Pattern: "/api/projects/{id}", Handler: func(w http.ResponseWriter, r *http.Request) error {
			id, err := parseID(r)
			if err != nil {
				return err
			}
			if err := s.Delete(r.Context(), id); err != nil {
				return err
			}
			w.WriteHeader(http.StatusNoContent)
			return nil
		}},
	}
}

// parseID lê {id} da rota. Só um inteiro positivo, sem sinal, vale; qualquer outra coisa é um
// projeto que não existe (spec §3.1).
func parseID(r *http.Request) (int64, error) {
	n, err := strconv.ParseUint(r.PathValue("id"), 10, 63)
	if err != nil || n == 0 {
		return 0, ErrNotFound
	}
	return int64(n), nil
}
