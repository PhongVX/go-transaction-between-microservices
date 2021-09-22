package redisx

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"orchestrator/internal/transactioncache"
	"orchestrator/pkg/constx"
	"orchestrator/pkg/msgx"
	"time"
)

type ServiceI interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, keys ...string) error
	SubscribeTransactionActions(ctx context.Context, topic string, tx *sql.Tx)
	PublishTransactionActions(ctx context.Context, topic string, action string) error
}

func NewRedisSrv(redisClient *redis.Client, txCacheSrv transactioncache.ServiceI) *Service {
	return &Service{
		redisClient: redisClient,
		txCacheSrv: txCacheSrv,
	}
}

func (s *Service) SubscribeTransactionActions(ctx context.Context, topic string, tx *sql.Tx)  {
	t := s.redisClient.Subscribe(ctx, topic)
	go func() {
		for  {
			// Checking time live of a topic
			//TODO read time alive of topic from config file
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute * 10)
			defer cancel()
			msg, err := t.ReceiveMessage(ctx)
			if err != nil {
				t.Close()
				fmt.Println(err)
				//Delete from redis cache and memory cache
				s.txCacheSrv.Remove(topic)
				s.redisClient.Del(context.Background(), topic)
				return
			}
			var txInfo msgx.TransactionInfo
			if err := json.Unmarshal([]byte(msg.Payload), &txInfo); err != nil {
				t.Close()
				s.txCacheSrv.Remove(topic)
				s.redisClient.Del(context.Background(), topic)
				return
			}
			switch txInfo.Action {
			case constx.Commit:
				tx.Commit()
				t.Close()
				s.txCacheSrv.Remove(topic)
				s.redisClient.Del(context.Background(), topic)
				return
			case constx.RollBack:
				tx.Rollback()
				t.Close()
				s.txCacheSrv.Remove(topic)
				s.redisClient.Del(context.Background(), topic)
				return
			}
		}
	}()
}

func (s *Service) PublishTransactionActions(ctx context.Context, topic string, action string) error  {
	txInfo := &msgx.TransactionInfo{CorrelationID: topic, Action: action}
	txInfoBytes, err := json.Marshal(txInfo)
	if err != nil {
		return err
	}
	if err := s.redisClient.Publish(ctx, topic, txInfoBytes).Err(); err != nil {
		return err
	}
	return nil
}

