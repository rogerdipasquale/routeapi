package api

import (
	"os"
	"routeapi/internal/config"

	"log/slog"
	"testing"
)

func testValidation(t *testing.T) {

	cfg := config.Load()

	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelInfo)
	opts := &slog.HandlerOptions{Level: lvl}
	log := slog.New(slog.NewTextHandler(os.Stdout, opts))

	deps := Deps{
		Config: cfg,
		Log:    log,
	}

	validReqs := [2]string{"app%3Druta%2Ccomponent%3Dbla", "app%3Druta"}
	invalidReqs := [1]string{"appruta%2Ccomponent%3Dbla"}

	for _, str := range validReqs {
		deps.ValidateLabelSelector(str)
	}

	for _, str := range invalidReqs {
		deps.ValidateLabelSelector(str)
	}

}
