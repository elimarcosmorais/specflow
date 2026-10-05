package main

import (
	"dev-tools/internal/features/preferences"
	"dev-tools/internal/features/projects"
	"dev-tools/internal/features/system"
	"dev-tools/internal/platform/httpserver"
)

// apiRoutes junta as rotas de todas as features. Só monta os handlers: nenhum serviço é chamado
// aqui, então o teste de contrato (CT-23) pode passar serviços nulos.
func apiRoutes(p *projects.Service, prefs *preferences.Service, sys *system.Service) []httpserver.Route {
	var routes []httpserver.Route
	routes = append(routes, projects.Routes(p)...)
	routes = append(routes, preferences.Routes(prefs)...)
	routes = append(routes, system.Routes(sys)...)
	return routes
}
