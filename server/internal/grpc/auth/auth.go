package auth

import (
	"Server_part_finance_control/protos/gen/go/auth"
	models "Server_part_finance_control/server/internal/domains"
	"Server_part_finance_control/server/internal/jwt"
	"Server_part_finance_control/server/internal/repository"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

type VKUserData struct{
	ID string
	Email string
	FirstName string
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

func (s *serverAPI) LoginVK(ctx context.Context, req *auth.VKAuthRequest)(*auth.AuthResponse, error){
	if req.VkToken == ""{
		return nil, fmt.Errorf("VK token is required")
	}

	vkUser, err := s.getVKUser(ctx, req.VkToken)
	if err != nil{
		return nil, fmt.Errorf("VK authentication failed: %v", err)
	}

	existingUser, err := s.userRepo.GetUserByVKID(ctx, vkUser.ID)
	if err == nil {
        token, err := jwt.NewToken(*existingUser, s.app, 24*time.Hour)
        if err != nil {
            return nil, fmt.Errorf("failed to generate token")
        }
        return &auth.AuthResponse{
            Token: token,
            User:  convertToAuthUser(existingUser),
        }, nil
    }

    if vkUser.Email != "" {
        emailUser, err := s.userRepo.GetUserByEmail(ctx, vkUser.Email)
        if err == nil {
            if err := s.userRepo.UpdateUserVKID(ctx, emailUser.ID, vkUser.ID); err != nil {
                return nil, fmt.Errorf("failed to link VK account")
            }
            token, err := jwt.NewToken(*emailUser, s.app, 24*time.Hour)
            if err != nil {
                return nil, fmt.Errorf("failed to generate token")
            }
            return &auth.AuthResponse{
                Token: token,
                User:  convertToAuthUser(emailUser),
            }, nil
        }
    }

    newUser := &models.User{
        ID:        fmt.Sprintf("user_%d", time.Now().UnixNano()),
        Email:     vkUser.Email,
        Username:  fmt.Sprintf("vk_%s", vkUser.FirstName),
        VKID:      &vkUser.ID,
        CreatedAt: time.Now(),
    }

    if err := s.userRepo.CreateUser(ctx, newUser, ""); 
    err != nil {
        return nil, fmt.Errorf("failed to create user")
    }

    token, err := jwt.NewToken(*newUser, s.app, 24*time.Hour)
    if err != nil {
        return nil, fmt.Errorf("failed to generate token")
    }

    return &auth.AuthResponse{
        Token: token,
        User:  convertToAuthUser(newUser),
    }, nil
}

func (s *serverAPI) getVKUser (ctx context.Context, token string) (*VKUserData, error){
	client := &http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("https://api.vk.com/method/users.get?access_token=%s&v=5.131&fields=email,first_name", token)

	resp, err := client.Get(url)
	if err != nil{
		return nil, fmt.Errorf("VK API request failed: %v", err)
	}
	defer resp.Body.Close()

	var result struct{
		Response []struct{
			ID        int    `json:"id"`
            FirstName string `json:"first_name"`
            Email     string `json:"email"`
		} `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil{
		return nil, fmt.Errorf("failed to decode VK response: %v", err)
	}

	if len(result.Response) == 0{
		return nil, fmt.Errorf("empty VK user data")
	}

	return &VKUserData{
        ID:        strconv.Itoa(result.Response[0].ID),
        Email:     result.Response[0].Email,
        FirstName: result.Response[0].FirstName,
    }, nil
}