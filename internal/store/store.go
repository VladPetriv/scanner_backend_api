package store

import (
	"time"

	"go.uber.org/zap"

	"github.com/VladPetriv/scanner_backend_api/internal/store/pg"
	"github.com/VladPetriv/scanner_backend_api/pkg/config"
	"github.com/VladPetriv/scanner_backend_api/pkg/logger"
)

const KeepAlivePollPeriod = 5

type Store struct {
	db  *pg.DB
	log *logger.Logger

	Channel ChannelRepo
	User    UserRepo
	Message MessageRepo
	Replie  ReplieRepo
	WebUser WebUserRepo
	Saved   SavedRepo
}

func New(cfg *config.Config, log *logger.Logger) (*Store, error) {
	db, err := pg.Dial(cfg)
	if err != nil {
		return nil, err
	}

	if db != nil {
		log.Info("running migrations...")

		err := runMigrations(cfg)
		if err != nil {
			return nil, err
		}
	}

	var store Store
	store.log = log

	if db != nil {
		store.db = db

		store.Channel = pg.NewChannelRepo(store.db)
		store.Message = pg.NewMessageRepo(store.db)
		store.Replie = pg.NewReplieRepo(store.db)
		store.User = pg.NewUserRepo(store.db)
		store.WebUser = pg.NewWebUserRepo(store.db)
		store.Saved = pg.NewSavedRepo(store.db)

		go store.keepAliveDB(cfg)
	}

	return &store, nil
}

func (s *Store) keepAliveDB(cfg *config.Config) {
	var err error

	for {
		time.Sleep(time.Second * KeepAlivePollPeriod)

		lostConnection := false
		if s.db == nil {
			lostConnection = true
		} else if _, err = s.db.Exec("SELECT 1;"); err != nil {
			lostConnection = true
		}

		if !lostConnection {
			continue
		}

		s.log.Debug("[store.KeepAliveDB] Lost db connection. Restoring...")

		s.db, err = pg.Dial(cfg)
		if err != nil {
			s.log.Error("failed to connect", zap.Error(err))

			continue
		}

		s.log.Debug("[store.KeepAliveDB] DB reconnected")
	}
}
