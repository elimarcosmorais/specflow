package preferences

import (
	"bytes"
	"errors"
	"log/slog"
	"strings"
	"sync"
	"testing"

	"dev-tools/internal/apperr"
)

// fakeRepo é um Repository em memória que registra cada Save.
type fakeRepo struct {
	loaded  Preferences
	loadErr error
	saveErr error
	saves   []Preferences
}

func (f *fakeRepo) Load() (Preferences, error) { return f.loaded, f.loadErr }

func (f *fakeRepo) Save(p Preferences) error {
	f.saves = append(f.saves, p)
	return f.saveErr
}

// testConfigPath é o caminho do config.json passado ao serviço nos testes.
const testConfigPath = "config.json"

// newTestService monta o serviço com um logger de texto que grava no buffer devolvido e um
// onChange que anota cada chamada em *changes.
func newTestService(repo Repository, detected string, changes *[]Preferences) (*Service, *bytes.Buffer) {
	logs := &bytes.Buffer{}
	logger := slog.New(slog.NewTextHandler(logs, nil))
	onChange := func(p Preferences) { *changes = append(*changes, p) }
	return NewService(repo, testConfigPath, func() string { return detected }, logger, onChange), logs
}

// CT-29
func TestService_Load_FirstRun(t *testing.T) {
	for _, detected := range []string{"pt-br", "en"} {
		t.Run(detected, func(t *testing.T) {
			repo := &fakeRepo{loadErr: ErrNotFound}
			var changes []Preferences
			svc, logs := newTestService(repo, detected, &changes)

			if err := svc.Load(); err != nil {
				t.Fatalf("Load() = %v, quer nil", err)
			}

			want := Preferences{Theme: "system", Language: detected}
			if got := svc.Get(); got != want {
				t.Errorf("Get() = %+v, quer %+v", got, want)
			}
			if len(repo.saves) != 1 || repo.saves[0] != want {
				t.Errorf("saves = %+v, quer [%+v]", repo.saves, want)
			}
			if !strings.Contains(logs.String(), "level=INFO msg=config.created theme=system language="+detected) {
				t.Errorf("log sem config.created com os padrões:\n%s", logs)
			}
			if strings.Contains(logs.String(), "config.invalid") {
				t.Errorf("primeira execução não deveria logar config.invalid:\n%s", logs)
			}
			if len(changes) != 0 {
				t.Errorf("onChange chamado %d vez(es) no Load, quer 0", len(changes))
			}
		})
	}
}

func TestService_Load_DetectorWithUnsupportedLanguageFallsBackToEn(t *testing.T) {
	repo := &fakeRepo{loadErr: ErrNotFound}
	svc, _ := newTestService(repo, "fr", new([]Preferences))

	if err := svc.Load(); err != nil {
		t.Fatalf("Load() = %v, quer nil", err)
	}

	if got, want := svc.Get(), (Preferences{Theme: "system", Language: "en"}); got != want {
		t.Errorf("Get() = %+v, quer %+v", got, want)
	}
}

func TestService_Load_ValidFileIsKeptWithoutSaveOrLog(t *testing.T) {
	repo := &fakeRepo{loaded: Preferences{Theme: "dark", Language: "en"}}
	svc, logs := newTestService(repo, "pt-br", new([]Preferences))

	if err := svc.Load(); err != nil {
		t.Fatalf("Load() = %v, quer nil", err)
	}

	if got, want := svc.Get(), (Preferences{Theme: "dark", Language: "en"}); got != want {
		t.Errorf("Get() = %+v, quer %+v", got, want)
	}
	if len(repo.saves) != 0 {
		t.Errorf("saves = %+v, quer nenhum", repo.saves)
	}
	if logs.Len() != 0 {
		t.Errorf("arquivo válido não deveria logar nada:\n%s", logs)
	}
}

// CT-30
func TestService_Load_InvalidFile(t *testing.T) {
	tests := []struct {
		name     string
		repo     *fakeRepo
		detected string
		want     Preferences
	}{
		{
			name:     "JSON inválido volta aos dois padrões",
			repo:     &fakeRepo{loadErr: errors.New("parse config.json: invalid character 'x'")},
			detected: "pt-br",
			want:     Preferences{Theme: "system", Language: "pt-br"},
		},
		{
			name:     "tema inválido: só o tema volta ao padrão",
			repo:     &fakeRepo{loaded: Preferences{Theme: "roxo", Language: "en"}},
			detected: "pt-br",
			want:     Preferences{Theme: "system", Language: "en"},
		},
		{
			name:     "idioma inválido: só o idioma volta ao padrão",
			repo:     &fakeRepo{loaded: Preferences{Theme: "dark", Language: "fr"}},
			detected: "pt-br",
			want:     Preferences{Theme: "dark", Language: "pt-br"},
		},
		{
			name:     "campos ausentes (strings vazias) voltam aos padrões",
			repo:     &fakeRepo{},
			detected: "en",
			want:     Preferences{Theme: "system", Language: "en"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, logs := newTestService(tt.repo, tt.detected, new([]Preferences))

			if err := svc.Load(); err != nil {
				t.Fatalf("Load() = %v, quer nil", err)
			}

			if got := svc.Get(); got != tt.want {
				t.Errorf("Get() = %+v, quer %+v", got, tt.want)
			}
			if len(tt.repo.saves) != 1 || tt.repo.saves[0] != tt.want {
				t.Errorf("saves = %+v, quer [%+v]", tt.repo.saves, tt.want)
			}
			if !strings.Contains(logs.String(), "level=WARN msg=config.invalid path="+testConfigPath+" err=") {
				t.Errorf("log sem config.invalid (WARN) com path=%s:\n%s", testConfigPath, logs)
			}
			if strings.Contains(logs.String(), "config.created") {
				t.Errorf("arquivo inválido não deveria logar config.created:\n%s", logs)
			}
		})
	}
}

func TestService_Load_SaveFailureKeepsDefaultsAndReturnsSaveFailed(t *testing.T) {
	cause := errors.New("disk full")
	repo := &fakeRepo{loadErr: ErrNotFound, saveErr: cause}
	svc, _ := newTestService(repo, "en", new([]Preferences))

	err := svc.Load()

	if !errors.Is(err, ErrSaveFailed) {
		t.Errorf("Load() = %v, quer errors.Is(err, ErrSaveFailed)", err)
	}
	if !errors.Is(err, cause) {
		t.Errorf("Load() = %v, quer a causa embrulhada", err)
	}
	if got, want := svc.Get(), (Preferences{Theme: "system", Language: "en"}); got != want {
		t.Errorf("Get() = %+v, quer os padrões %+v", got, want)
	}
}

// CT-31
func TestService_Update(t *testing.T) {
	repo := &fakeRepo{loaded: Preferences{Theme: "light", Language: "pt-br"}}
	var changes []Preferences
	svc, logs := newTestService(repo, "pt-br", &changes)
	if err := svc.Load(); err != nil {
		t.Fatal(err)
	}
	want := Preferences{Theme: "dark", Language: "en"}

	if err := svc.Update(want); err != nil {
		t.Fatalf("Update() = %v, quer nil", err)
	}

	if len(repo.saves) != 1 || repo.saves[0] != want {
		t.Errorf("saves = %+v, quer [%+v]", repo.saves, want)
	}
	if got := svc.Get(); got != want {
		t.Errorf("Get() = %+v, quer %+v", got, want)
	}
	if len(changes) != 1 || changes[0] != want {
		t.Errorf("onChange = %+v, quer uma chamada com %+v", changes, want)
	}
	if !strings.Contains(logs.String(), "level=INFO msg=preferences.updated theme=dark language=en") {
		t.Errorf("log sem preferences.updated:\n%s", logs)
	}
}

func TestService_Update_SaveFailure(t *testing.T) {
	cause := errors.New("access denied")
	repo := &fakeRepo{loaded: Preferences{Theme: "light", Language: "pt-br"}}
	var changes []Preferences
	svc, logs := newTestService(repo, "pt-br", &changes)
	if err := svc.Load(); err != nil {
		t.Fatal(err)
	}
	repo.saveErr = cause

	err := svc.Update(Preferences{Theme: "dark", Language: "en"})

	var ae *apperr.Error
	if !errors.As(err, &ae) || ae.Code != "preferences.save_failed" || ae.Kind != apperr.Internal {
		t.Fatalf("Update() = %v, quer preferences.save_failed (Internal)", err)
	}
	if !errors.Is(err, cause) {
		t.Errorf("Update() = %v, quer a causa embrulhada", err)
	}
	if got, want := svc.Get(), (Preferences{Theme: "light", Language: "pt-br"}); got != want {
		t.Errorf("Get() = %+v, quer %+v (inalterado)", got, want)
	}
	if len(changes) != 0 {
		t.Errorf("onChange chamado %d vez(es), quer 0", len(changes))
	}
	if strings.Contains(logs.String(), "preferences.updated") {
		t.Errorf("falha não deveria logar preferences.updated:\n%s", logs)
	}
}

// CT-32
func TestService_Update_Invalid(t *testing.T) {
	tests := []struct {
		name      string
		input     Preferences
		wantCode  string
		wantField string
	}{
		{"tema inválido", Preferences{Theme: "roxo", Language: "en"}, "preferences.theme_invalid", "theme"},
		{"idioma inválido", Preferences{Theme: "dark", Language: "fr"}, "preferences.language_invalid", "language"},
		{"ambos vazios: o tema vem primeiro", Preferences{}, "preferences.theme_invalid", "theme"},
		{"tema e idioma inválidos: o tema vem primeiro", Preferences{Theme: "roxo", Language: "fr"}, "preferences.theme_invalid", "theme"},
		{"idioma com maiúsculas é inválido", Preferences{Theme: "dark", Language: "PT-BR"}, "preferences.language_invalid", "language"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initial := Preferences{Theme: "light", Language: "pt-br"}
			repo := &fakeRepo{loaded: initial}
			var changes []Preferences
			svc, _ := newTestService(repo, "pt-br", &changes)
			if err := svc.Load(); err != nil {
				t.Fatal(err)
			}

			err := svc.Update(tt.input)

			var ae *apperr.Error
			if !errors.As(err, &ae) {
				t.Fatalf("Update() = %v, quer *apperr.Error", err)
			}
			if ae.Code != tt.wantCode || ae.Field != tt.wantField || ae.Kind != apperr.Invalid {
				t.Errorf("erro = {Kind: %s, Code: %s, Field: %s}, quer {Invalid, %s, %s}",
					ae.Kind, ae.Code, ae.Field, tt.wantCode, tt.wantField)
			}
			if len(repo.saves) != 0 {
				t.Errorf("saves = %+v, quer nenhum", repo.saves)
			}
			if got := svc.Get(); got != initial {
				t.Errorf("Get() = %+v, quer %+v (inalterado)", got, initial)
			}
			if len(changes) != 0 {
				t.Errorf("onChange chamado %d vez(es), quer 0", len(changes))
			}
		})
	}
}

func TestService_Update_AcceptsEveryValidValue(t *testing.T) {
	for _, theme := range []string{"light", "dark", "system"} {
		for _, language := range []string{"pt-br", "en"} {
			t.Run(theme+"/"+language, func(t *testing.T) {
				svc, _ := newTestService(&fakeRepo{}, "en", new([]Preferences))
				want := Preferences{Theme: theme, Language: language}

				if err := svc.Update(want); err != nil {
					t.Fatalf("Update(%+v) = %v, quer nil", want, err)
				}
				if got := svc.Get(); got != want {
					t.Errorf("Get() = %+v, quer %+v", got, want)
				}
			})
		}
	}
}

// Get roda em goroutines de requisição enquanto Update roda em outra: o valor lido nunca pode
// ser uma mistura de dois Updates. Só é conclusivo com -race.
func TestService_GetDuringUpdate(t *testing.T) {
	svc, _ := newTestService(&fakeRepo{}, "en", new([]Preferences))
	if err := svc.Load(); err != nil {
		t.Fatal(err)
	}
	a := Preferences{Theme: "dark", Language: "en"}
	b := Preferences{Theme: "light", Language: "pt-br"}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 200 {
			for _, p := range []Preferences{a, b} {
				if err := svc.Update(p); err != nil {
					t.Error(err)
					return
				}
			}
		}
	}()
	for range 1000 {
		if got := svc.Get(); got != a && got != b && got != (Preferences{Theme: "system", Language: "en"}) {
			t.Errorf("Get() = %+v: mistura de valores", got)
			break
		}
	}
	wg.Wait()
}
