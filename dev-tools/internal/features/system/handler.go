package system

import (
	"net/http"

	"dev-tools/internal/platform/httpserver"
)

// Routes devolve as rotas da feature: GET /api/health (pública) e
// POST /api/system/open-logs-folder (protegida por token).
func Routes(s *Service) []httpserver.Route {
	return []httpserver.Route{
		{Method: http.MethodGet, Pattern: "/api/health", Public: true, Handler: func(w http.ResponseWriter, r *http.Request) error {
			httpserver.WriteJSON(w, http.StatusOK, s.Health(r.Context()))
			return nil
		}},
		{Method: http.MethodPost, Pattern: "/api/system/open-logs-folder", Handler: func(w http.ResponseWriter, r *http.Request) error {
			if err := s.OpenLogsFolder(); err != nil {
				return err
			}
			w.WriteHeader(http.StatusNoContent)
			return nil
		}},
	}
}
