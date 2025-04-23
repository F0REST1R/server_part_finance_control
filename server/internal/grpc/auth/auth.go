package auth

import (
	"Server_part_finance_control/protos/gen/go/auth"
	models "Server_part_finance_control/server/internal/domains"
	"Server_part_finance_control/server/internal/jwt"
	"Server_part_finance_control/server/internal/repository"
	"context"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
)


type serverAPI struct {
	auth.UnimplementedAuthServiceServer
	userRepo *repository.UserRepository
	app models.App
}

func convertToAuthUser(user *models.User) *auth.User{
	return &auth.User{
		Id: user.ID,
		Email: user.Email,
		Username: user.Username,
	}
}

func NewAuthServiceClient(userRepo *repository.UserRepository, app models.App) *serverAPI {
	return &serverAPI{
		userRepo: userRepo,
		app: app,
	}
}


func Register(gRPC *grpc.Server, app models.App, userRepo *repository.UserRepository){
	authService := NewAuthServiceClient(userRepo, app)
	auth.RegisterAuthServiceServer(gRPC, authService)
}

func (s *serverAPI) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.AuthResponse, error){
	if (req.Email == "" || req.Password == "" || req.Username == "") {
		 return nil, fmt.Errorf("Email, Password and Username are required")
	}

	if len(req.Password) < 6 {
		return nil, fmt.Errorf("Password must be at least 6 characters")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil{
		return nil, fmt.Errorf("failed to hash password")
	}

	user := &models.User{
		ID: fmt.Sprintf("user_%d", time.Now().UnixNano()),
		Email: strings.ToLower(req.Email),
		Username: req.Username,
		PasswordHash: string(hashedPassword),
		CreatedAt: time.Now(),
	}

	err = s.userRepo.CreateUser(ctx, user, string(hashedPassword))
	if err != nil {
		return nil, fmt.Errorf("failed to generate user")
	}

	token, err := jwt.NewToken(*user, s.app, 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to genereate token")
	}

	return &auth.AuthResponse{
		Token: token,
		User: convertToAuthUser(user),
	}, nil

}	

func (s *serverAPI) Login(ctx context.Context, req *auth.LoginRequest) (*auth.AuthResponse, error){
	if req.Email == "" || req.Password == ""{
		return nil, fmt.Errorf("email and password are required")
	}

	user, err := s.userRepo.GetUserByEmail(ctx, strings.ToLower(req.Email))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil{
		return nil, fmt.Errorf("invalid credentials")
	}

	token, err := jwt.NewToken(*user, s.app, 24*time.Hour)
	if err != nil{
		return nil, fmt.Errorf("failed to generate token")
	}

	return &auth.AuthResponse{
		Token: token,
		User: convertToAuthUser(user),
	}, nil
}



