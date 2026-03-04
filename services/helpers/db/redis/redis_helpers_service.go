package redis_helpers_service

import (
	"context"
	"fmt"
	"strings"
	"time"

	errors_constants "go_boilerplate_project/constants/errors"
	redis_models "go_boilerplate_project/models/databases/redis"
)

func (s *service) resolveContext(ctx context.Context, cancelFunc context.CancelFunc) (context.Context, context.CancelFunc) {
	if ctx != nil && cancelFunc != nil {
		return ctx, cancelFunc
	}
	return s.Input.Services.Context.GetContext()
}

func quoteKeys(keys []string) string {
	quoted := make([]string, len(keys))
	for i, k := range keys {
		quoted[i] = fmt.Sprintf("%q", k)
	}
	return strings.Join(quoted, " ")
}

func expirationSuffix(d time.Duration) string {
	if d <= 0 {
		return ""
	}
	secs := int64(d.Seconds())
	if secs > 0 {
		return fmt.Sprintf(" EX %d", secs)
	}
	return fmt.Sprintf(" PX %d", d.Milliseconds())
}

// ──────────────────────────────────────────────
// Set (Insert)
// ──────────────────────────────────────────────

func (s *service) Set(input *redis_models.SetInput) error {
	if input == nil {
		return errors_constants.ErrRedisSetInputNil
	}
	if input.Key == "" {
		return errors_constants.ErrRedisKeyRequired
	}
	if input.Value == nil {
		return errors_constants.ErrRedisValueRequired
	}

	ctx, cancel := s.resolveContext(input.Context, input.CancelFunc)
	defer cancel()

	query := fmt.Sprintf("SET %q %q%s",
		input.Key, fmt.Sprint(input.Value), expirationSuffix(input.Expiration))

	s.Input.Logger.Debugw("Redis SET",
		"query", query,
		"key", input.Key,
		"expiration", input.Expiration,
	)

	err := s.Input.Client.RedisClient.Set(ctx, input.Key, input.Value, input.Expiration).Err()
	if err != nil {
		s.Input.Logger.Errorw("Redis SET failed",
			"query", query,
			"error", err,
		)
		return err
	}

	s.Input.Logger.Debugw("Redis SET success", "key", input.Key)
	return nil
}

// ──────────────────────────────────────────────
// SetNX (Insert If Not Exists)
// ──────────────────────────────────────────────

func (s *service) SetNX(input *redis_models.SetNXInput) (*redis_models.SetNXOutput, error) {
	if input == nil {
		return nil, errors_constants.ErrRedisSetNXInputNil
	}
	if input.Key == "" {
		return nil, errors_constants.ErrRedisKeyRequired
	}
	if input.Value == nil {
		return nil, errors_constants.ErrRedisValueRequired
	}

	ctx, cancel := s.resolveContext(input.Context, input.CancelFunc)
	defer cancel()

	query := fmt.Sprintf("SET %q %q NX%s",
		input.Key, fmt.Sprint(input.Value), expirationSuffix(input.Expiration))

	s.Input.Logger.Debugw("Redis SET NX",
		"query", query,
		"key", input.Key,
		"expiration", input.Expiration,
	)

	ok, err := s.Input.Client.RedisClient.SetNX(ctx, input.Key, input.Value, input.Expiration).Result()
	if err != nil {
		s.Input.Logger.Errorw("Redis SET NX failed",
			"query", query,
			"error", err,
		)
		return nil, err
	}

	s.Input.Logger.Debugw("Redis SET NX completed",
		"key", input.Key,
		"wasSet", ok,
	)
	return &redis_models.SetNXOutput{Success: ok}, nil
}

// ──────────────────────────────────────────────
// Get (value + remaining TTL)
// ──────────────────────────────────────────────

func (s *service) Get(input *redis_models.GetInput) (*redis_models.GetOutput, error) {
	if input == nil {
		return nil, errors_constants.ErrRedisGetInputNil
	}
	if input.Key == "" {
		return nil, errors_constants.ErrRedisKeyRequired
	}

	ctx, cancel := s.resolveContext(input.Context, input.CancelFunc)
	defer cancel()

	query := fmt.Sprintf("GET %q", input.Key)

	s.Input.Logger.Debugw("Redis GET + TTL",
		"query", query,
		"key", input.Key,
	)

	val, err := s.Input.Client.RedisClient.Get(ctx, input.Key).Result()
	if err != nil {
		s.Input.Logger.Errorw("Redis GET failed",
			"query", query,
			"error", err,
		)
		return nil, err
	}

	ttl, err := s.Input.Client.RedisClient.TTL(ctx, input.Key).Result()
	if err != nil {
		s.Input.Logger.Errorw("Redis TTL failed",
			"query", fmt.Sprintf("TTL %q", input.Key),
			"error", err,
		)
		return nil, err
	}

	s.Input.Logger.Debugw("Redis GET success",
		"key", input.Key,
		"ttl", ttl,
	)

	return &redis_models.GetOutput{
		Value: val,
		TTL:   ttl,
	}, nil
}

// ──────────────────────────────────────────────
// Update (SET XX — only if key already exists)
// ──────────────────────────────────────────────

func (s *service) Update(input *redis_models.UpdateInput) (*redis_models.UpdateOutput, error) {
	if input == nil {
		return nil, errors_constants.ErrRedisUpdateInputNil
	}
	if input.Key == "" {
		return nil, errors_constants.ErrRedisKeyRequired
	}
	if input.Value == nil {
		return nil, errors_constants.ErrRedisValueRequired
	}

	ctx, cancel := s.resolveContext(input.Context, input.CancelFunc)
	defer cancel()

	expiration := input.Expiration
	if input.KeepTTL {
		ttl, err := s.Input.Client.RedisClient.TTL(ctx, input.Key).Result()
		if err != nil {
			s.Input.Logger.Errorw("Redis TTL lookup for KeepTTL failed",
				"key", input.Key,
				"error", err,
			)
			return nil, err
		}
		if ttl > 0 {
			expiration = ttl
		}
	}

	suffix := " XX"
	if input.KeepTTL {
		suffix += " KEEPTTL"
	}
	suffix += expirationSuffix(expiration)
	query := fmt.Sprintf("SET %q %q%s",
		input.Key, fmt.Sprint(input.Value), suffix)

	s.Input.Logger.Debugw("Redis SET XX (Update)",
		"query", query,
		"key", input.Key,
		"keepTTL", input.KeepTTL,
		"expiration", expiration,
	)

	ok, err := s.Input.Client.RedisClient.SetXX(ctx, input.Key, input.Value, expiration).Result()
	if err != nil {
		s.Input.Logger.Errorw("Redis SET XX failed",
			"query", query,
			"error", err,
		)
		return nil, err
	}

	s.Input.Logger.Debugw("Redis SET XX completed",
		"key", input.Key,
		"updated", ok,
	)
	return &redis_models.UpdateOutput{Updated: ok}, nil
}

// ──────────────────────────────────────────────
// Delete
// ──────────────────────────────────────────────

func (s *service) Delete(input *redis_models.DeleteInput) (*redis_models.DeleteOutput, error) {
	if input == nil {
		return nil, errors_constants.ErrRedisDeleteInputNil
	}
	if len(input.Keys) == 0 {
		return nil, errors_constants.ErrRedisKeysRequired
	}

	ctx, cancel := s.resolveContext(input.Context, input.CancelFunc)
	defer cancel()

	query := fmt.Sprintf("DEL %s", quoteKeys(input.Keys))

	s.Input.Logger.Debugw("Redis DEL",
		"query", query,
		"keys", input.Keys,
	)

	count, err := s.Input.Client.RedisClient.Del(ctx, input.Keys...).Result()
	if err != nil {
		s.Input.Logger.Errorw("Redis DEL failed",
			"query", query,
			"error", err,
		)
		return nil, err
	}

	s.Input.Logger.Debugw("Redis DEL success",
		"keys", input.Keys,
		"deletedCount", count,
	)
	return &redis_models.DeleteOutput{DeletedCount: count}, nil
}

// ──────────────────────────────────────────────
// Exists
// ──────────────────────────────────────────────

func (s *service) Exists(input *redis_models.ExistsInput) (*redis_models.ExistsOutput, error) {
	if input == nil {
		return nil, errors_constants.ErrRedisExistsInputNil
	}
	if len(input.Keys) == 0 {
		return nil, errors_constants.ErrRedisKeysRequired
	}

	ctx, cancel := s.resolveContext(input.Context, input.CancelFunc)
	defer cancel()

	query := fmt.Sprintf("EXISTS %s", quoteKeys(input.Keys))

	s.Input.Logger.Debugw("Redis EXISTS",
		"query", query,
		"keys", input.Keys,
	)

	count, err := s.Input.Client.RedisClient.Exists(ctx, input.Keys...).Result()
	if err != nil {
		s.Input.Logger.Errorw("Redis EXISTS failed",
			"query", query,
			"error", err,
		)
		return nil, err
	}

	s.Input.Logger.Debugw("Redis EXISTS success",
		"keys", input.Keys,
		"existCount", count,
	)
	return &redis_models.ExistsOutput{Count: count}, nil
}

// ──────────────────────────────────────────────
// Expire (set / update TTL)
// ──────────────────────────────────────────────

func (s *service) Expire(input *redis_models.ExpireInput) (*redis_models.ExpireOutput, error) {
	if input == nil {
		return nil, errors_constants.ErrRedisExpireInputNil
	}
	if input.Key == "" {
		return nil, errors_constants.ErrRedisKeyRequired
	}

	ctx, cancel := s.resolveContext(input.Context, input.CancelFunc)
	defer cancel()

	query := fmt.Sprintf("EXPIRE %q %d", input.Key, int64(input.Expiration.Seconds()))

	s.Input.Logger.Debugw("Redis EXPIRE",
		"query", query,
		"key", input.Key,
		"expiration", input.Expiration,
	)

	ok, err := s.Input.Client.RedisClient.Expire(ctx, input.Key, input.Expiration).Result()
	if err != nil {
		s.Input.Logger.Errorw("Redis EXPIRE failed",
			"query", query,
			"error", err,
		)
		return nil, err
	}

	s.Input.Logger.Debugw("Redis EXPIRE completed",
		"key", input.Key,
		"applied", ok,
	)
	return &redis_models.ExpireOutput{Success: ok}, nil
}
