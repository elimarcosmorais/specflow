package preferences

import (
	"net/http"

	"dev-tools/internal/platform/httpserver"
)

// Routes devolve as rotas da feature: PUT /api/preferences, protegida por token.
func Routes(s *Service) []httpserver.Route {
	return []httpserver.Route{
		{Method: http.MethodPut, Pattern: "/api/preferences", Handler: func(w http.ResponseWriter, r *http.Request) error {
			var p Preferences
			if err := httpserver.DecodeJSON(r, &p); err != nil {
				return err
			}
			if err := s.Update(p); err != nil {
				return err
			}
			httpserver.WriteJSON(w, http.StatusOK, p)
			return nil
		}},
	}
}
