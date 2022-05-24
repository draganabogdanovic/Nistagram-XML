package service

import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/KristijanPill/Nishtagram/auth-service/payload"
	"github.com/KristijanPill/Nishtagram/auth-service/repository"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	PublicKey              *rsa.PublicKey
	secretKey              *rsa.PrivateKey
	hmacKey                []byte
	refreshTokenRepository *repository.RefreshTokenRepository
	credentialsRepository  *repository.CredentialsRepository
}

const accessTokenDuration = 1
const refreshTokenDuration = 43200

type LoginResponse struct {
	AccessToken         string `json:"access_token"`
	AccessTokenDuration int64  `json:"access_token_duration"`
	RefreshToken        string `json:"refresh_token"`
}

var ErrUnathorized = errors.New("unathorized")

func NewAuthService(publicKey *rsa.PublicKey, secretKey *rsa.PrivateKey, hmacKey []byte, refreshTokenRepository *repository.RefreshTokenRepository, credentialsRepository *repository.CredentialsRepository) *AuthService {
	return &AuthService{
		PublicKey:              publicKey,
		secretKey:              secretKey,
		hmacKey:                hmacKey,
		refreshTokenRepository: refreshTokenRepository,
		credentialsRepository:  credentialsRepository,
	}
}

func (service *AuthService) Login(loginRequest *payload.LoginRequest) (*LoginResponse, error) {
	credentials, err := service.credentialsRepository.FindByUsername(loginRequest.Username)

	if err != nil {

		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(credentials.Password), []byte(loginRequest.Password)); err != nil {
		return nil, err
	}

	return service.generateTokens(credentials.ID)
}

func (service *AuthService) Register(credentials *payload.Credentials) error {
	passwordBytes, err := service.hashPassword(credentials.Password)

	if err != nil {
		return err
	}

	credentials.Password = string(passwordBytes)

	return service.credentialsRepository.Create(credentials)
}

func (service *AuthService) Refresh(refreshToken string, userID uuid.UUID) (string, error) {
	token, err := service.refreshTokenRepository.FindByUserID(userID.String())

	if err != nil {
		return "", err
	}

	if token.Token != refreshToken {
		return "", ErrUnathorized
	}

	return service.generateAccessToken(userID)
}

func (service *AuthService) generateTokens(userID uuid.UUID) (*LoginResponse, error) {
	accessToken, err := service.generateAccessToken(userID)

	if err != nil {
		return nil, err
	}

	refreshToken, err := service.generateRefreshToken(userID)

	if err != nil {
		return nil, err
	}

	service.saveRefreshToken(refreshToken, userID)

	return &LoginResponse{
		AccessToken:         accessToken,
		RefreshToken:        refreshToken,
		AccessTokenDuration: time.Now().Add(accessTokenDuration * time.Minute).Unix(),
	}, nil
}

func (service *AuthService) generateAccessToken(userID uuid.UUID) (string, error) {
	claims := jwt.StandardClaims{
		Subject:   userID.String(),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(accessTokenDuration * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.GetSigningMethod("RS256"),
		claims,
	)

	return token.SignedString(service.secretKey)
}

func (service *AuthService) generateRefreshToken(userID uuid.UUID) (string, error) {
	claims := jwt.StandardClaims{
		Subject:   userID.String(),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(refreshTokenDuration * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.GetSigningMethod("HS256"),
		claims,
	)

	return token.SignedString(service.hmacKey)
}

func (service *AuthService) saveRefreshToken(refreshToken string, userID uuid.UUID) error {
	if service.refreshTokenRepository.ExistsByUserID(userID.String()) {
		token, err := service.refreshTokenRepository.FindByUserID(userID.String())

		if err != nil {
			return err
		}

		token.Token = refreshToken
		service.refreshTokenRepository.Update(&token)
	} else {
		err := service.refreshTokenRepository.Create(&payload.RefreshToken{
			UserID: userID,
			Token:  refreshToken,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (service *AuthService) hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 14)
}
