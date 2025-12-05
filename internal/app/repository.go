package app

import (
	mysql "github.com/epg-sync/epgsync/internal/repository/db"
	"github.com/epg-sync/epgsync/pkg/logger"
)

func (app *App) initializeRepositories() error {
	logger.Debug("Initializing repositories...")

	app.repos = &Repositories{
		Channel:         mysql.NewChannelRepository(app.db),
		Program:         mysql.NewProgramRepository(app.db),
		ChannelMappings: mysql.NewChannelMappingsRepository(app.db),
		Timezone:        mysql.NewTimezoneRepository(app.db),
		User:            mysql.NewUserRepository(app.db),
	}

	return nil

}
