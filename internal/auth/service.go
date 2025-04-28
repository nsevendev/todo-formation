package auth

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nsevenpack/ginresponse"
	"github.com/nsevenpack/logger/v2/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userService struct {
	userRepo userRepoInterface
	jwtKey   []byte
}

type UserServiceInterface interface {
	Register(ctx context.Context, userCreateDto UserCreateDto) (*User, error)
	Login(ctx context.Context, userLoginDto UserLoginDto) (string, error)
	ValidateToken(tokenString string) (*tokenClaims, error)
	GetProfilCurrentUser(ctx context.Context, id primitive.ObjectID) (*User, error)
    GetIdUserInContext(ctx *gin.Context) primitive.ObjectID
    DeleteOneByUser(ctx context.Context, id primitive.ObjectID) error
    DeleteByAdmin(ctx context.Context, ids []primitive.ObjectID) (int, error)
}

func NewUserService(userRepo userRepoInterface, jwtKey string) UserServiceInterface {
	return &userService{
		userRepo: userRepo,
		jwtKey:   []byte(jwtKey),
	}
}

func (s *userService) Register(ctx context.Context, userCreateDto UserCreateDto) (*User, error) {
    roleDefault := "user"

    user := &User{
        Email:    userCreateDto.Email,
        Password: userCreateDto.Password,
        Username: userCreateDto.Username,
        Role:     userCreateDto.Role,
    }

    if user.Role == "" {
        user.Role = roleDefault
    }

    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }

    return user, nil
}

func (s *userService) Login(ctx context.Context, userLoginDto UserLoginDto) (string, error) {
    user, err := s.userRepo.FindByEmail(ctx, userLoginDto.Email)
    if err != nil {
        return "", err
    }
    if user == nil {
        return "", errors.New("identifiants invalides")
    }

    if !user.CheckPassword(userLoginDto.Password) {
        return "", errors.New("identifiants invalides")
    }

    token, err := s.generateToken(user)
    if err != nil {
        return "", err
    }

    return token, nil
}

func (s *userService) generateToken(user *User) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour) // change here for time expiration
    claims := &tokenClaims{
        IdUser: user.ID.Hex(),
        Email:  user.Email,
        Role:   user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    
    tokenString, err := token.SignedString(s.jwtKey)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

func (s *userService) ValidateToken(tokenString string) (*tokenClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
        return s.jwtKey, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*tokenClaims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, errors.New("le token est invalide")
}

func (s *userService) GetProfilCurrentUser(ctx context.Context, id primitive.ObjectID) (*User, error) {
    return s.userRepo.FindByID(ctx, id)
}

func (s *userService) GetIdUserInContext(ctx *gin.Context) primitive.ObjectID {
    idUser, exists := ctx.Get("id_user")
	if !exists {
		logger.Ef("Erreur d'authentification : ID utilisateur non trouvé dans le contexte")
        ginresponse.Unauthorized(ctx, "Erreur d'authentification", "Vous n'avez pas les droits pour effectuer cette action")
        ctx.Abort()
	}

    return idUser.(primitive.ObjectID)
}

func (s *userService) DeleteOneByUser(ctx context.Context, id primitive.ObjectID) error {
    if _, err := s.userRepo.Delete(ctx, id); err != nil {
        return err
    }

    return nil
}

func (s *userService) DeleteByAdmin(ctx context.Context, ids []primitive.ObjectID) (int, error) {
    deletedCount := 0

    for _, id := range ids {
        d, err := s.userRepo.Delete(ctx, id)
        if err != nil {
            continue
        }

        if d == 1 {
            deletedCount++
        }
    }

    return deletedCount, nil
}