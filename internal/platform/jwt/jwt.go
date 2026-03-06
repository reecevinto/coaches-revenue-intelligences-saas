package jwt

import (
	"crypto/rsa"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

type Claims struct {
	UserID    string `json:"user_id"`
	AccountID string `json:"account_id"`
	Role      string `json:"role"`
	jwt.RegisteredClaims
}

func NewService(privateKeyPath, publicKeyPath string) (*Service, error) {

	privBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	pubBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privBytes)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
	if err != nil {
		return nil, err
	}

	return &Service{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

func (s *Service) GenerateAccessToken(userID, accountID, role string) (string, error) {

	claims := Claims{
		UserID:    userID,
		AccountID: accountID,
		Role:      role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(s.privateKey)
}

func (s *Service) ValidateToken(tokenStr string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return s.publicKey, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims := token.Claims.(*Claims)

	return claims, nil
}
