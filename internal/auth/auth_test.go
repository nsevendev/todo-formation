package auth

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"todof/internal/config"
	initializer "todof/internal/init"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var r UserRepoInterface
var s UserServiceInterface
var c *mongo.Collection
var ctx context.Context
var cancelCtx context.Context
var cancelFunc context.CancelFunc
var router *gin.Engine
var middleware AuthMiddlewareInterface
var ids []primitive.ObjectID
var users []*User
var tokenString string

//service
func TestMain(m *testing.M) {
	c = initializer.Db.Collection("users")
	r = NewUserRepo(initializer.Db)
	s = NewUserService(r, config.Get("JWT_SECRET"))
	ctx := context.Background()

	cancelCtx, cancelFunc = context.WithCancel(ctx)
	cancelFunc()

	router = gin.New()
	middleware = NewAuthMiddleware(s)

	if _, err := c.DeleteMany(ctx, bson.M{}); err != nil {
		log.Fatalf("Erreur lors du nettoyage de la collection users : %v", err)
	}

	code := m.Run()

	os.Exit(code)
}

func TestRegister(t *testing.T){
	tests := []struct {
		name	 string
		email    string
		password string
		username string
		role     string
		isUser	 bool
		isErr    bool
	}{
		{"test avec user admin valid", "admin@gmail.com", "password", "admin", "admin", true, false},
		{"test avec email existant", "admin@gmail.com", "password", "admin2", "admin", false, true},
		{"test avec propriete role vide", "admin2@gmail.com", "password", "admin2", "", true, false},
	}

	for _, tt := range tests {
		userCreateDto := UserCreateDto{
			Email:    tt.email,
			Password: tt.password,
			Username: tt.username,
			Role:     tt.role,
		}

		user, err := s.Register(ctx, userCreateDto)

		if (err != nil) != tt.isErr {
			t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
		}

		if user != nil {
			ids = append(ids, user.ID)
			users = append(users, user)
		}
	}
}

func TestLogin(t *testing.T) {
	serviceWithJWTEmpty := NewUserService(r, config.Get(""))

	tests := []struct {
		name     string
		email    string
		password string
		isToken  bool
		isErr    bool
	}{
		{"test success", "admin@gmail.com", "password", true, false},
		{"test avec email inexistant", "fail@gmail.com", "password", false, true},
		{"test avec mauvais password", "admin@gmail.com", "invalid password", false, true},
		{"test avec JWT_SECRET vide", "admin@gmail.com", "password", false, true},
		{"test echec mongo", "admin@gmail.com", "password", false, true},
	}

	for _, tt := range tests {
		userLoginDto := UserLoginDto{
			Email:    tt.email,
			Password: tt.password,
		}
		var token string
		var err error

		switch tt.name {
		case "test echec mongo":
			token, err = s.Login(cancelCtx, userLoginDto)

			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}

			if (token == "") == tt.isToken {
				t.Errorf("%s: login fail with email %v and password %v", tt.name, tt.email, tt.password)
			}

		case "test avec JWT_SECRET vide":
			token, err = serviceWithJWTEmpty.Login(ctx, userLoginDto)

			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}

			if (token == "") == tt.isToken {
				t.Errorf("%s: login fail with email %v and password %v", tt.name, tt.email, tt.password)
			}

		default:
			token, err = s.Login(ctx, userLoginDto)

			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}

			if (token == "") == tt.isToken {
				t.Errorf("%s: login fail with email %v and password %v", tt.name, tt.email, tt.password)
			}

			if token != "" {
				tokenString = token
			}
		}
	}
}

func TestValidateToken(t *testing.T){
	tests := []struct {
		name string
		tokenString string
		isTokenClaims bool
		isErr bool
	}{
		{"test success", tokenString, true, false},
	}

	for _, tt := range tests {
		token, err := s.ValidateToken(tt.tokenString)

		if (err != nil) != tt.isErr {
			t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
		}

		if (token == nil) == tt.isTokenClaims {
			t.Errorf("%s: got %v, expect token claims", tt.name, token)
		}
	}
}

func TestGetProfilCurrentUser(t *testing.T){
	tests := []struct {
		name string
		id primitive.ObjectID
		isUser bool
		isErr bool
	}{
		{"test success", ids[0], true, false},
		{"test id valide mais inexistant", primitive.NewObjectID(), false, false},
		{"test echec mongo", ids[0], false, true},
	}

	for _, tt := range tests {
		switch tt.name {
		case "test echec mongo":
			user, err := s.GetProfilCurrentUser(cancelCtx, tt.id)

			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}

			if (user == nil) == tt.isUser {
				t.Errorf("%s: got %v, expect user", tt.name, user)
			}

		default:
			user, err := s.GetProfilCurrentUser(ctx, tt.id)

			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}

			if (user == nil) == tt.isUser {
				t.Errorf("%s: got %v, expect user", tt.name, user)
			}
		}
	}
}

func TestGetIdUserInContext(t *testing.T) {
    tests := []struct {
        name      string
        setup     func(c *gin.Context)
        isAborted bool
        expectID  primitive.ObjectID
    }{
        {
            name: "test success",
            setup: func(c *gin.Context) {
                c.Set("id_user", ids[0])
            },
            isAborted: false,
            expectID:  ids[0],
        },
        {
            name: "test sans id_user dans le contexte",
            setup: func(c *gin.Context) {
            },
            isAborted: true,
        },
    }

    for _, tt := range tests {
        w := httptest.NewRecorder()
        c, _ := gin.CreateTestContext(w)
        tt.setup(c)

        func() {
            defer func() {
                if r := recover(); r != nil {
                    if !tt.isAborted {
                        t.Errorf("%s: panique inattendue : %v", tt.name, r)
                    }
                }
            }()
            
            id := s.GetIdUserInContext(c)

            if c.IsAborted() != tt.isAborted {
                t.Errorf("%s: expected abort %v, got aborted %v", tt.name, tt.isAborted, c.IsAborted())
            }

            if !tt.isAborted && id != tt.expectID {
                t.Errorf("%s: expected id %v, got %v", tt.name, tt.expectID, id)
            }
        }()
    }
}

func TestDeleteOneByUser(t *testing.T) {
	tests := []struct {
		name    string
		id      primitive.ObjectID
		isErr   bool
	}{
		{"test succès avec un utilisateur existant", ids[0], false},
		{"test avec un ID valide mais inexistant dans la base", primitive.NewObjectID(), false},
		{"test avec un ID invalide", primitive.NilObjectID, false},
		{"test echec mongo", ids[0], true},
	}

	for _, tt := range tests {
		switch tt.name {
		case "test echec mongo":
			err := s.DeleteOneByUser(cancelCtx, tt.id)

			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}

		default:
			err := s.DeleteOneByUser(ctx, tt.id)

			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}
		}
	}
}

func TestDeleteByAdmin(t *testing.T){
	tests := []struct {
		name string
		deletedCount int
		isErr bool
	}{
		{"test delete success", 1, false},
		{"test echec mongo", 0, false},
	}

	for _, tt := range tests {
		switch tt.name {
		case "test echec mongo":
			deletedCount, err := s.DeleteByAdmin(cancelCtx, ids)

			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}

			if deletedCount != tt.deletedCount {
				t.Errorf("%s: got deletedCount : %v, expect deletedCount : %v", tt.name, deletedCount, tt.deletedCount)
			}

		default:
			deletedCount, err := s.DeleteByAdmin(ctx, ids)

			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}

			if deletedCount != tt.deletedCount {
				t.Errorf("%s: got deletedCount : %v, expect deletedCount : %v", tt.name, deletedCount, tt.deletedCount)
			}
		}
	}
}

func TestDeleteAllByAdmin(t *testing.T){
	tests := []struct {
		name string
		deletedCount int64
		isErr bool
	}{
		{"test success", 0, false},
		{"test echec mongo", 0, true},
	}

	for _, tt := range tests {
		switch tt.name {
		case "test echec mongo":
			deletedCount, err := s.DeleteAllByAdmin(cancelCtx)

			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}

			if deletedCount != tt.deletedCount {
				t.Errorf("%s: got deletedCount %v, expect deletedCount %v", tt.name, deletedCount, tt.deletedCount)
			}
			
		default:
			deletedCount, err := s.DeleteAllByAdmin(ctx)

			if (err != nil) != tt.isErr {
				t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
			}

			if deletedCount != tt.deletedCount {
				t.Errorf("%s: got deletedCount %v, expect deletedCount %v", tt.name, deletedCount, tt.deletedCount)
			}
		}
	}
}

//middleware
func TestRequireAuth(t *testing.T) {

	router.GET("/protected", middleware.RequireAuth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "access granted"})
	})

	generateTokenWithInvalidIdUser := func() string {
		claims := &tokenClaims{
			IdUser: "idinvalid",
			Email:  "test@example.com",
			Role:   "user",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(config.Get("JWT_SECRET")))
		return tokenString
	}

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{"test success", "Bearer " + tokenString, http.StatusOK},
		{"test sans Authorization header", "", http.StatusUnauthorized},
		{"test avec Authorization header invalid", "Invalid", http.StatusUnauthorized},
		{"test avec token invalid", "Bearer invalid", http.StatusUnauthorized},
		{"test avec token avec Id invalid", "Bearer " + generateTokenWithInvalidIdUser(), http.StatusUnauthorized},
	}

	for _, tt := range tests {
		req, _ := http.NewRequest("GET", "/protected", nil)
		if tt.authHeader != "" {
			req.Header.Set("Authorization", tt.authHeader)
		}
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != tt.expectedStatus {
			t.Errorf("%s: got status %d, expected %d", tt.name, w.Code, tt.expectedStatus)
		}
	}
}

func TestRequireRole(t *testing.T) {

	router.GET("/admin", func(c *gin.Context) {
		c.Set("role", "admin")
		middleware.RequireRole("admin")(c)
		if !c.IsAborted() {
			c.JSON(http.StatusOK, gin.H{"message": "admin access granted"})
		}
	})

	router.GET("/user", func(c *gin.Context) {
		c.Set("role", "user")
		middleware.RequireRole("admin")(c)
		if !c.IsAborted() {
			c.JSON(http.StatusOK, gin.H{"message": "access denied"})
		}
	})

	router.GET("/invalid", func(c *gin.Context) {
		middleware.RequireRole("admin")(c)
		if !c.IsAborted() {
			c.JSON(http.StatusOK, gin.H{"message": "should not happen"})
		}
	})

	router.GET("/invalid-type", func(c *gin.Context) {
		c.Set("role", 123456789)
		middleware.RequireRole("admin")(c)
		if !c.IsAborted() {
			c.JSON(http.StatusOK, gin.H{"message": "should not happen"})
		}
	})

	tests := []struct {
		name           string
		route          string
		expectedStatus int
	}{
		{"test success pour admin", "/admin", http.StatusOK},
		{"test access refuse", "/user", http.StatusForbidden},
		{"test sans role", "/invalid", http.StatusUnauthorized},
		{"test role mauvais type", "/invalid-type", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		req, _ := http.NewRequest("GET", tt.route, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != tt.expectedStatus {
			t.Errorf("%s: got status %d, expected status %d", tt.name, w.Code, tt.expectedStatus)
		}
	}
}

//usermodel
func TestHashPassword(t *testing.T) {
	tests := []struct {
		name         string
		password  string
		isErr    bool
	}{
		{"test success", "password", false},
		{"test avec password vide", "", false},
	}

	for _, tt := range tests {
		user := &User{
			Password: tt.password,
		}

		err := user.HashPassword()

		if (err != nil) != tt.isErr {
			t.Errorf("%s: got error %v, expect error %v", tt.name, err, tt.isErr)
		}

		if user.Password == tt.password {
			t.Errorf("%s: echec du hashage du password", tt.name)
		}
	}
}

func TestCheckPassword(t *testing.T) {
	hashed, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("erreur de génération du hash pour test: %v", err)
	}
	
	tests := []struct {
		name         string
		password   string
		isMatch  bool
	}{
		{"test success", "password", true},
		{"test avec mauvais password", "incorrect", false},
		{"test avec password vide", "", false},
	}

	for _, tt := range tests {
		user := &User{
			Password: string(hashed),
		}

		match := user.CheckPassword(tt.password)

		if match != tt.isMatch {
			t.Errorf("%s: got match %v, expect match %v", tt.name, match, tt.isMatch)
		}
	}
}
