package users

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/conflux-tech/fiber-rest-boilerplate/configs"
	"github.com/go-redis/redis/v8"
)

// UseCase holds the user usecase struct definition
type UseCase struct {
	pgRepo Repository
	redis  redis.Client
}

// NewUseCase initializes a new user usecase
func NewUseCase(pgRepo Repository, redis redis.Client) *UseCase {
	return &UseCase{
		pgRepo: pgRepo,
		redis:  redis,
	}
}

// List func is the implementation of listing users from data source
func (u *UseCase) List(ctx context.Context, page, limit int) ([]*User, error) {
	users, err := u.pgRepo.List(page, limit)
	if err != nil {
		return []*User{}, err
	}
	return users, nil
}

// Create func is the implemenation of creating user data into various data sources
func (u *UseCase) Create(ctx context.Context, user *User) (*User, error) {
	user, err := u.pgRepo.Create(user)
	if err != nil {
		return &User{}, err
	}
	cacheKey := fmt.Sprintf("%s:%s", "user", strconv.Itoa(int(user.ID)))
	cacheTTL := configs.Get().Redis.UserDataTTL
	marshalledUser, _ := json.Marshal(user)
	if rErr := u.redis.Set(ctx, cacheKey, marshalledUser, cacheTTL); rErr.Err() != nil {
		panic(rErr)
	}
	return user, nil
}

// Get func is the implementation of getting user details from various data sources
func (u *UseCase) Get(ctx context.Context, id int) (*User, error) {
	cacheKey := fmt.Sprintf("%s:%s", "user", strconv.Itoa(id))
	cache, _ := u.redis.Get(ctx, cacheKey).Result()
	if cache != "" {
		cachedUser := User{}
		json.Unmarshal([]byte(cache), &cachedUser)
		return &cachedUser, nil
	}

	user, err := u.pgRepo.Get(id)
	if err != nil {
		return &User{}, err
	}

	if user.ID != 0 {
		marshalledUser, _ := json.Marshal(user)
		cacheTTL := configs.Get().Redis.UserDataTTL
		if rErr := u.redis.Set(ctx, cacheKey, marshalledUser, cacheTTL); rErr.Err() != nil {
			panic(rErr)
		}
	}

	return user, nil
}

// Update func is the implementation of updating user details in various data sources
func (u *UseCase) Update(ctx context.Context, id int, user *User) (*User, error) {
	updatedUser, err := u.pgRepo.Update(id, user)
	if err != nil {
		return &User{}, err
	}
	cacheKey := fmt.Sprintf("%s:%s", "user", strconv.Itoa(int(id)))
	cacheTTL := configs.Get().Redis.UserDataTTL
	marshalledUser, _ := json.Marshal(user)

	u.redis.Del(ctx, cacheKey)
	if rErr := u.redis.Set(ctx, cacheKey, marshalledUser, cacheTTL); rErr.Err() != nil {
		panic(rErr)
	}
	return updatedUser, nil
}

// Delete func is the implementation of deleiting user details in various data sources
func (u *UseCase) Delete(ctx context.Context, id int) (bool, error) {
	cacheKey := fmt.Sprintf("%s:%s", "user", strconv.Itoa(int(id)))
	res, err := u.pgRepo.Delete(id)
	if err != nil {
		return false, err
	}
	u.redis.Del(ctx, cacheKey)
	return res, nil
}
