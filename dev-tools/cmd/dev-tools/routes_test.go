package main

import (
	"net/http"
	"os"
	"slices"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

const openAPIPath = "../../internal/platform/httpserver/docs/openapi.yaml"

// CT-23 (ADR 001): as rotas registradas e o openapi.yaml não podem divergir.
func TestAPIRoutes_MatchOpenAPI(t *testing.T) {
	data, err := os.ReadFile(openAPIPath)
	if err != nil {
		t.Fatalf("read openapi.yaml: %v", err)
	}
	var doc struct {
		Paths map[string]map[string]yaml.Node `yaml:"paths"`
	}
	if err := yaml.Unmarshal(data, &doc); err != nil {
		t.Fatalf("parse openapi.yaml: %v", err)
	}

	var contract []string
	for path, item := range doc.Paths {
		for method := range item {
			switch method {
			case "get", "put", "post", "delete", "patch", "head", "options":
				contract = append(contract, strings.ToUpper(method)+" "+path)
			}
		}
	}

	var registered []string
	for _, r := range apiRoutes(nil, nil, nil) {
		registered = append(registered, r.Method+" "+r.Pattern)
	}

	slices.Sort(contract)
	slices.Sort(registered)
	if !slices.Equal(registered, contract) {
		t.Errorf("routes differ from openapi.yaml\nregistered: %v\ncontract:   %v", registered, contract)
	}
}

// CT-24
func TestAPIRoutes_OnlyHealthIsPublic(t *testing.T) {
	var public []string
	for _, r := range apiRoutes(nil, nil, nil) {
		if r.Public {
			public = append(public, r.Method+" "+r.Pattern)
		}
	}
	want := []string{http.MethodGet + " /api/health"}
	if !slices.Equal(public, want) {
		t.Errorf("public routes = %v, want %v", public, want)
	}
}
