package store

import (
	"time"

	"github.com/VladPetriv/scanner_backend_api/internal/store/pg"
	"github.com/VladPetriv/scanner_backend_api/pkg/logger"
	"go.uber.org/zap"
)

const KeepAlivePollPeriod = 5

type Store struct {
	DB  *pg.DB
	log *logger.Logger

	Channel ChannelRepo
	User    UserRepo
	Message MessageRepo
	Replie  ReplieRepo
	WebUser WebUserRepo
	Saved   SavedRepo
}

func New(log *logger.Logger) (*Store, error) {
	db, err := pg.Dial()
	if err != nil {
		return nil, err
	}

	if db != nil {
		log.Info("running migrations...")

		err := runMigrations()
		if err != nil {
			return nil, err
		}
	}

	var store Store
	store.log = log

	if db != nil {
		store.DB = db

		go store.KeepAliveDB()
	}

	return &store, nil
}

func (s *Store) KeepAliveDB() {
	var err error

	for {
		time.Sleep(time.Second * KeepAlivePollPeriod)

		lostConnection := false
		if s.DB == nil {
			lostConnection = true
		} else if _, err := s.DB.Exec("SELECT 1;"); err != nil {
			lostConnection = true
		}

		if !lostConnection {
			continue
		}

		s.log.Debug("[store.KeepAliveDB] Lost db connection. Restoring...")

		s.DB, err = pg.Dial()
		if err != nil {
			s.log.Error("failed to connect", zap.Error(err))

			continue
		}

		s.log.Debug("[store.KeepAliveDB] DB reconnected")
	}
}
