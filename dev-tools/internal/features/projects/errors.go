package projects

import "dev-tools/internal/apperr"

// Erros da feature (spec §3.2). Field é o campo da entrada a que o erro se refere.
var (
	ErrNameRequired = apperr.New(apperr.Invalid, "project.name_required").WithField("name")
	ErrPathRequired = apperr.New(apperr.Invalid, "project.path_required").WithField("path")
	ErrTypeRequired = apperr.New(apperr.Invalid, "project.type_required").WithField("type")
	ErrTypeInvalid  = apperr.New(apperr.Invalid, "project.type_invalid").WithField("type")
	ErrPathInvalid  = apperr.New(apperr.Invalid, "project.path_invalid").WithField("path")
	ErrPathNotFound = apperr.New(apperr.Invalid, "project.path_not_found").WithField("path")
	ErrNameTaken    = apperr.New(apperr.Conflict, "project.name_taken").WithField("name")
	ErrPathTaken    = apperr.New(apperr.Conflict, "project.path_taken").WithField("path")
	ErrNotFound     = apperr.New(apperr.NotFound, "project.not_found")
)
