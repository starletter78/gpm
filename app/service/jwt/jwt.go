package jwt

import (
	"errors"
	"fmt"
	"gpm/global"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt/v5"
)

// AccessClaims 定义访问令牌的声明
type AccessClaims struct {
	UserClaims
	Type string `json:"type"`
	jwt.RegisteredClaims
}

type UserClaims struct {
	Id   string   `json:"id"`
	Role []string `json:"role"`
}

// RefreshClaims 定义刷新令牌的声明
type RefreshClaims struct {
	Id   string `json:"id"`
	Type string `json:"type"`
	jwt.RegisteredClaims
}

// BaseClaims 基础声明，包含用户ID
type BaseClaims struct {
	Id uint `json:"id"`
}

// JWT 配置结构
type JWT struct {
}

// NewJWT 创建新的JWT实例
func NewJWT() *JWT {
	return &JWT{}
}

// generateAccessToken 生成访问令牌
func (j *JWT) generateAccessToken(id string, roles []string) (string, error) {
	claims := AccessClaims{
		UserClaims: UserClaims{
			Id:   id,
			Role: roles,
		},
		Type: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(global.Config.Jwt.AccessExpire) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    global.Config.Jwt.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(global.Config.Jwt.AccessSecret))
}

// generateRefreshToken 生成刷新令牌
func (j *JWT) generateRefreshToken(id string) (string, error) {
	claims := RefreshClaims{
		Id:   id,
		Type: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(global.Config.Jwt.RefreshExpire) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    global.Config.Jwt.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(global.Config.Jwt.RefreshSecret))
}

// ParseAccessToken 解析访问令牌
func (j *JWT) ParseAccessToken(tokenString string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		return []byte(global.Config.Jwt.AccessSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("访问令牌已过期")
		}
		return nil, fmt.Errorf("解析访问令牌失败: %w", err)
	}

	if claims, ok := token.Claims.(*AccessClaims); ok && token.Valid && claims.Type == "access" {
		return claims, nil
	}
	return nil, errors.New("无效的访问令牌")
}

// ParseRefreshToken 解析刷新令牌
func (j *JWT) ParseRefreshToken(tokenString string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		return []byte(global.Config.Jwt.RefreshSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("刷新令牌已过期")
		}
		return nil, fmt.Errorf("解析刷新令牌失败: %w", err)
	}

	if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid && claims.Type == "refresh" {
		return claims, nil
	}

	return nil, errors.New("无效的刷新令牌")
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// RefreshTokens 使用刷新令牌获取新的令牌对
func (j *JWT) RefreshTokens(refreshTokenString string, tenant string) (*TokenPair, error) {
	// 先解析刷新令牌
	claims, err := j.ParseRefreshToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("刷新令牌无效: %w", err)
	}
	roles := global.CasbinEnforcer.GetRolesForUserInDomain(claims.ID, tenant)
	accessToken, err := j.generateAccessToken(claims.Id, roles)
	if err != nil {
		return nil, err
	}
	if time.Now().Sub(claims.ExpiresAt.Time) < 24*3*time.Hour {
		refreshTokenString, err = j.generateRefreshToken(claims.Id)
		if err != nil {
			return nil, err
		}
	} else {
		refreshTokenString = ""
	}
	// 使用解析出的用户ID生成新的令牌对
	return &TokenPair{AccessToken: accessToken, RefreshToken: refreshTokenString}, nil
}

func (j *JWT) GenPairToken(tenant string, userId string) (*TokenPair, error) {
	roles := global.CasbinEnforcer.GetRolesForUserInDomain(userId, tenant)
	refreshToken, err := j.generateRefreshToken(userId)
	if err != nil {
		return nil, err
	}
	accessToken, err := j.generateAccessToken(userId, roles)
	if err != nil {
		return nil, err
	}
	return &TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func GetClaimsByGin(c *gin.Context) (*UserClaims, error) {
	if user, ok := c.Get("user"); ok {
		return user.(*UserClaims), nil
	} else {
		return nil, errors.New("user not found in context")
	}
}
