package kafka

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/VladPetriv/scanner_backend_api/internal/model"
	"github.com/VladPetriv/scanner_backend_api/internal/service"
	"github.com/VladPetriv/scanner_backend_api/internal/store/pg"
	"github.com/VladPetriv/scanner_backend_api/pkg/config"
	"github.com/VladPetriv/scanner_backend_api/pkg/logger"
	"go.uber.org/zap"
)

func createWorker(addr string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	conn, err := sarama.NewConsumer([]string{addr}, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}

	return conn, nil
}

func GetChannelFromQueue(srvManager *service.Manager, cfg *config.Config, log *logger.Logger) {
	worker, err := createWorker(cfg.KafkaAddr)
	if err != nil {
		log.Error("failed to create kafka worker", zap.Error(err))
	}

	consumer, err := worker.ConsumePartition("channels.get", 0, sarama.OffsetOldest)
	if err != nil {
		log.Error("failed to create consumer", zap.Error(err))
	}

	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				log.Error("failed to get channels from kafka queue", zap.Error(err))

				return

			case data := <-consumer.Messages():
				channel := model.ChannelDTO{}

				err := json.Unmarshal(data.Value, &channel)
				if err != nil {
					log.Error("unmarshal error", zap.Error(err))
				}

				candidate, err := srvManager.Channel.GetChannelByName(channel.Name)
				if err != nil && !errors.Is(err, pg.ErrChannelNotFound) {
					log.Error("get channel by name error", zap.Error(err))
				}

				if candidate == nil {
					err := srvManager.Channel.CreateChannel(&channel)
					if err != nil {
						log.Error("create channel error", zap.Error(err))
					}
				} else {
					log.Info(fmt.Sprintf("channel with name %s is exist", channel.Name))
				}
			}
		}
	}()
}

func GetDataFromQueue(srvManager *service.Manager, cfg *config.Config, log *logger.Logger) {
	worker, err := createWorker(cfg.KafkaAddr)
	if err != nil {
		log.Error("failed to create kafka worker", zap.Error(err))
	}

	consumer, err := worker.ConsumePartition("messages.get", 0, sarama.OffsetOldest)
	if err != nil {
		log.Error("failed to create consumer", zap.Error(err))
	}

	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				log.Error("failed to get data from kafka queue", zap.Error(err))

				return
			case data := <-consumer.Messages():
				telegramMessage := model.TgMessage{}

				err := json.Unmarshal(data.Value, &telegramMessage)
				if err != nil {
					log.Error("unmarshal error", zap.Error(err))
				}

				channel, err := srvManager.Channel.GetChannelByName(telegramMessage.PeerID.Username)
				if err != nil {
					log.Error("get channel by name error", zap.Error(err))
				}

				userID, err := srvManager.User.CreateUser(&model.UserDTO{
					Username: telegramMessage.FromID.Username,
					Fullname: telegramMessage.FromID.Fullname,
					ImageURL: telegramMessage.FromID.ImageURL,
				})
				if err != nil {
					log.Error("create user error", zap.Error(err))
				}

				messageID, err := srvManager.Message.CreateMessage(&model.MessageDTO{
					ChannelID:  channel.ID,
					UserID:     userID,
					Title:      telegramMessage.Message,
					MessageURL: telegramMessage.MessageURL,
					ImageURL:   telegramMessage.ImageURL,
				})
				if err != nil {
					log.Error("create message error", zap.Error(err))
				}

				for _, replie := range telegramMessage.Replies.Messages {
					userID, err := srvManager.User.CreateUser(&model.UserDTO{
						Username: replie.FromID.Username,
						Fullname: replie.FromID.Fullname,
						ImageURL: replie.FromID.ImageURL,
					})
					if err != nil {
						log.Error("create user for replie error", zap.Error(err))
					}

					err = srvManager.Replie.CreateReplie(&model.ReplieDTO{
						MessageID: messageID,
						UserID:    userID,
						Title:     replie.Message,
						ImageURL:  replie.ImageURL,
					})
					if err != nil {
						log.Error("create replie error", zap.Error(err))
					}
				}
			}
		}
	}()
}
