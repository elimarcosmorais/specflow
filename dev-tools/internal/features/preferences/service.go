package preferences

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"
)

// Service guarda as preferências atuais em memória e as persiste pelo Repository.
type Service struct {
	repo           Repository
	configPath     string
	detectLanguage func() string
	logger         *slog.Logger
	onChange       func(Preferences)

	writeMu sync.Mutex // serializa Load e Update, para a ordem de gravação ser a da memória
	mu      sync.RWMutex
	current Preferences
}

// NewService cria o serviço. configPath é o caminho do config.json, usado só no atributo path do
// log config.invalid (a leitura e a gravação são do repo); detectLanguage devolve o idioma do
// Windows para a primeira execução; onChange é chamado com as novas preferências a cada Update
// bem-sucedido.
func NewService(repo Repository, configPath string, detectLanguage func() string, logger *slog.Logger, onChange func(Preferences)) *Service {
	return &Service{repo: repo, configPath: configPath, detectLanguage: detectLanguage, logger: logger, onChange: onChange}
}

// Load carrega as preferências do repositório (spec §4.3). Arquivo ausente, conteúdo inválido ou
// campo inválido usam os padrões (só no campo inválido, neste último caso) e regravam o
// repositório. O valor em memória fica definido mesmo quando a regravação falha: nesse caso o
// erro devolvido é ErrSaveFailed, e o app pode seguir com os padrões.
func (s *Service) Load() error {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()

	loaded, err := s.repo.Load()
	if err == nil && validTheme(loaded.Theme) && validLanguage(loaded.Language) {
		s.set(loaded)
		return nil
	}

	var prefs Preferences
	switch {
	case errors.Is(err, ErrNotFound):
		prefs = s.defaults()
		s.logger.Info("config.created", "theme", prefs.Theme, "language", prefs.Language)
	case err != nil:
		prefs = s.defaults()
		s.logger.Warn("config.invalid", "path", s.configPath, "err", err)
	default:
		prefs = loaded
		var problems []error
		if !validTheme(prefs.Theme) {
			problems = append(problems, fmt.Errorf("theme %q is invalid", prefs.Theme))
			prefs.Theme = DefaultTheme
		}
		if !validLanguage(prefs.Language) {
			problems = append(problems, fmt.Errorf("language %q is invalid", prefs.Language))
			prefs.Language = s.defaultLanguage()
		}
		s.logger.Warn("config.invalid", "path", s.configPath, "err", errors.Join(problems...))
	}
	s.set(prefs)

	if err := s.repo.Save(prefs); err != nil {
		return ErrSaveFailed.Wrap(err)
	}
	return nil
}

// Get devolve as preferências atuais. Pode ser chamado de qualquer goroutine.
func (s *Service) Get() Preferences {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.current
}

// Update valida (tema antes do idioma), grava e passa a usar as novas preferências. Com valor
// inválido ou falha na gravação, nada muda e onChange não é chamado.
func (s *Service) Update(p Preferences) error {
	if !validTheme(p.Theme) {
		return ErrThemeInvalid
	}
	if !validLanguage(p.Language) {
		return ErrLanguageInvalid
	}

	s.writeMu.Lock()
	defer s.writeMu.Unlock()

	if err := s.repo.Save(p); err != nil {
		return ErrSaveFailed.Wrap(err)
	}
	s.set(p)
	s.logger.Info("preferences.updated", "theme", p.Theme, "language", p.Language)
	if s.onChange != nil {
		s.onChange(p)
	}
	return nil
}

func (s *Service) set(p Preferences) {
	s.mu.Lock()
	s.current = p
	s.mu.Unlock()
}

func (s *Service) defaults() Preferences {
	return Preferences{Theme: DefaultTheme, Language: s.defaultLanguage()}
}

// defaultLanguage é o idioma detectado; um valor fora dos suportados vira en (RN-09).
func (s *Service) defaultLanguage() string {
	if language := s.detectLanguage(); validLanguage(language) {
		return language
	}
	return LanguageEn
}
