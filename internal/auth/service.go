package auth

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"

	"github.com/reecevinto/coaches-revenue-intelligences-saas/internal/platform/jwt"
	"github.com/reecevinto/coaches-revenue-intelligences-saas/internal/users"
)

type Service struct {
	usersRepo users.Repository
	jwt       *jwt.Service
}

func NewService(usersRepo users.Repository, jwt *jwt.Service) *Service {
	return &Service{
		usersRepo: usersRepo,
		jwt:       jwt,
	}
}

const (
	argonTime    = 1
	argonMemory  = 64 * 1024
	argonThreads = 4
	argonKeyLen  = 32
	saltLength   = 16
)

func generateSalt() ([]byte, error) {
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func hashPassword(password string) (string, error) {

	salt, err := generateSalt()
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		uint32(argonTime),
		uint32(argonMemory),
		uint8(argonThreads),
		uint32(argonKeyLen),
	)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf(
		"argon2id$%d$%d$%d$%s$%s",
		argonTime,
		argonMemory,
		argonThreads,
		b64Salt,
		b64Hash,
	)

	return encoded, nil
}

func comparePassword(encodedHash, password string) bool {

	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false
	}

	timeParam := uint32(argonTime)
	memoryParam := uint32(argonMemory)
	threadsParam := uint8(argonThreads)

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}

	newHash := argon2.IDKey(
		[]byte(password),
		salt,
		timeParam,
		memoryParam,
		threadsParam,
		uint32(len(hash)),
	)

	return subtle.ConstantTimeCompare(hash, newHash) == 1
}

type RegisterRequest struct {
	AccountID string `json:"account_id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Role      string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (string, error) {

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return "", err
	}

	user := &users.User{
		ID:        uuid.NewString(),
		AccountID: req.AccountID,
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      req.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.usersRepo.Create(ctx, user)
	if err != nil {
		return "", err
	}

	token, err := s.jwt.GenerateAccessToken(user.ID, user.AccountID, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (string, error) {

	user, err := s.usersRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !comparePassword(user.Password, req.Password) {
		return "", errors.New("invalid credentials")
	}

	token, err := s.jwt.GenerateAccessToken(user.ID, user.AccountID, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) Me(ctx context.Context, userID string) (*users.User, error) {

	user, err := s.usersRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
