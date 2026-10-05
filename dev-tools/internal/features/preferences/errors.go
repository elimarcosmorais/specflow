package preferences

import "dev-tools/internal/apperr"

// Erros da feature (spec §3.2). A causa de ErrSaveFailed é acrescentada com Wrap.
var (
	ErrThemeInvalid    = apperr.New(apperr.Invalid, "preferences.theme_invalid").WithField("theme")
	ErrLanguageInvalid = apperr.New(apperr.Invalid, "preferences.language_invalid").WithField("language")
	ErrSaveFailed      = apperr.New(apperr.Internal, "preferences.save_failed")
)
