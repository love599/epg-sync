package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/epg-sync/epgsync/internal/model"
	"github.com/epg-sync/epgsync/pkg/logger"
	"github.com/epg-sync/epgsync/pkg/utils"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func (app *App) initializeDatabase() error {
	logger.Debug("Initializing database...",
		logger.String("driver", app.cfg.Database.Driver),
		logger.String("host", app.cfg.Database.Host),
		logger.Int("port", app.cfg.Database.Port),
	)
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}

	var dialector gorm.Dialector
	databaseCfg := app.cfg.Database

	switch databaseCfg.Driver {
	case "mysql":
		logger.Debug("Connecting to MySQL...",
			logger.String("host", databaseCfg.Host),
			logger.Int("port", databaseCfg.Port))

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&timeout=10s&readTimeout=10s&writeTimeout=10s",
			databaseCfg.User, databaseCfg.Password, databaseCfg.Host, databaseCfg.Port, databaseCfg.Name)
		dialector = mysql.Open(dsn)
	case "sqlite":
		logger.Debug("Connecting to SQLite...",
			logger.String("filepath", databaseCfg.Name))

		dbPath := databaseCfg.Name
		if filepath.Ext(dbPath) == "" {
			dbPath += ".db"
		}
		dir := filepath.Dir(dbPath)
		if dir != "." && dir != "/" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create database directory: %w", err)
			}
		}
		logger.Info("Connecting to SQLite...", logger.String("file", dbPath))
		dialector = sqlite.Open(dbPath)
	default:
		return fmt.Errorf("unsupported database driver: %s", databaseCfg.Driver)
	}

	db, err := gorm.Open(dialector, gormConfig)

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if app.cfg.Database.Debug {
		app.db = db.Debug()
	} else {
		app.db = db
	}

	if err := seedDefaultData(app.db); err != nil {
		logger.Error("Failed to seed default data", logger.Err(err))
	}

	logger.Debug("Database initialized successfully")

	return nil
}

func seedDefaultData(db *gorm.DB) error {
	var count int64
	db.Model(&model.User{}).Count(&count)

	if count == 0 {
		logger.Info("No users found, creating default admin user...")

		randomPassword, err := utils.GenerateRandomString(14)
		if err != nil {
			return err
		}
		hashedPassword, _ := utils.HashPassword(string(randomPassword))
		admin := model.User{
			Username: "admin",
			Password: hashedPassword,
			Role:     "admin",
			IsActive: 1,
		}
		if err := db.Create(&admin).Error; err != nil {
			return err
		}
		logger.Info(fmt.Sprintf("Default admin created. User: admin, Pass: %s", randomPassword))
	}
	return nil
}
