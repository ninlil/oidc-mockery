package main

import (
	"github.com/ninlil/oidc-mockery/internal/config"
	"github.com/ninlil/oidc-mockery/internal/handlers"

	"github.com/ninlil/butler"
	"github.com/ninlil/butler/log"
	"github.com/ninlil/butler/router"
)

func main() {
	defer butler.Cleanup(nil)

	// Load configuration
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	// Get routes with config
	routes := handlers.GetRoutes(cfg)

	// Start server using router.Serve
	log.Info().Int("port", cfg.Server.Port).Msg("Starting OIDC Mockery server")
	err = router.Serve(routes,
		router.WithPort(cfg.Server.Port),
		router.WithStrictSlash(false),
		router.WithExposedErrors(),
		router.Without204(), // Return 200 instead of 204 for empty responses
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}

	butler.Run()
}
