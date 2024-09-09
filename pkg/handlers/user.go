package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/harshgupta9473/Social/pkg/middleware"
	"github.com/harshgupta9473/Social/pkg/models"
	"github.com/harshgupta9473/Social/pkg/services"
	"github.com/harshgupta9473/Social/pkg/utils"
)

type UserHandler struct {
	AuthService *services.AuthService
}

func NewUserHandler(auth *services.AuthService) *UserHandler {
	return &UserHandler{AuthService: auth}
}
func GenerateToken() (string, error) {
	token := make([]byte, 16)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}

func (handler *UserHandler) HandleRegistration(w http.ResponseWriter, r *http.Request) {
	var usereq models.NewRegRequest
	err := json.NewDecoder(r.Body).Decode(&usereq)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	_, err = handler.AuthService.UserRepo.GetUserByEmail(usereq.Email)
	if err == nil {
		http.Error(w, "user already exist", http.StatusConflict)
		return
	}
	token, err := GenerateToken()
	if err != nil {
		http.Error(w, "internal server error1", http.StatusInternalServerError)
		return
	}
	tmpUser, err := handler.AuthService.UserRepo.GetTempUserByEmail(usereq.Email)
	if err == nil {
		err = handler.AuthService.UserRepo.UpdateTempUser(tmpUser.Email, usereq.Password, token)
		if err != nil {
			http.Error(w, "internal server error2", http.StatusInternalServerError)
			return
		}
	} else {
		err = handler.AuthService.UserRepo.InsertIntoTempUser(usereq.Email, usereq.Password, token)
		if err != nil {
			log.Println(err)
			http.Error(w, "internal server error3", http.StatusInternalServerError)
			return
		}
	}
	err = utils.SendVerificationEmail(usereq.Email, token)
	if err != nil {
		http.Error(w, "internal server error4", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, "Email sent for Email Verification")

}

func (handler *UserHandler) HandleVerification(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("userID")
	token := r.URL.Query().Get("token")
	tmpUser, err := handler.AuthService.UserRepo.GetTempUserByEmail(email)
	if err != nil {
		http.Error(w, "not allowed", http.StatusForbidden)
		return
	}
	if tmpUser.Token != token {
		http.Error(w, "unauthorised", http.StatusUnauthorized)
		return
	}
	if time.Now().After(tmpUser.ExpiresAt) {
		http.Error(w, "link expired", http.StatusForbidden)
		return
	}
	err = handler.AuthService.Register(*tmpUser)
	if err != nil {
		http.Error(w, "error occured", http.StatusInternalServerError)
		return
	}
	userAc, err := handler.AuthService.UserRepo.GetUserByEmail(tmpUser.Email)
	if err != nil {
		http.Error(w, "error occured", http.StatusInternalServerError)
		return
	}
	jwtToken, err := middleware.CreateJWT(userAc.Email, userAc.CreatedAt, userAc.ID)
	if err != nil {
		http.Error(w, "error occured", http.StatusInternalServerError)
		return
	}
	utils.SetTokenInCookies(w, jwtToken)
	utils.WriteJSON(w, http.StatusOK, &models.LogInResponse{
		Email: userAc.Email,
	})
}

func (handler *UserHandler) HandleCreateUserProfile(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
	var userReq models.ProfileReq
	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	emailFromJWT, err := utils.ExtractEmailFromJWT(token)
	if err != nil {
		http.Error(w, "error occured", http.StatusInternalServerError)
		return
	}
	if userReq.Email != emailFromJWT {
		http.Error(w, "unauthorised access", http.StatusUnauthorized)
		return
	}
	_, err = handler.AuthService.UserRepo.GetUserProfileByUserName(userReq.UserName)
	if err == nil {
		http.Error(w, "username already exists", http.StatusConflict)
		return
	}

	err = handler.AuthService.UserRepo.CreateUserProfile(userReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	user, err := handler.AuthService.UserRepo.GetUserProfileByEmail(userReq.Email)
	if err != nil {
		log.Println(err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, user)
}

func (handler *UserHandler) HandleUserProfileUpdation(w http.ResponseWriter, r *http.Request, token *jwt.Token) {
	var userReq models.ProfileReq
	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		http.Error(w, "invalid formate", http.StatusBadRequest)
		return
	}
	emailFromJWT, err := utils.ExtractEmailFromJWT(token)
	if err != nil {
		http.Error(w, "error occured 1", http.StatusInternalServerError)
		return
	}
	if userReq.Email != emailFromJWT {
		http.Error(w, "unauthorised access", http.StatusUnauthorized)
		return
	}
	_, err = handler.AuthService.UserRepo.GetUserProfileByUserName(userReq.UserName)
	if err == nil {
		http.Error(w, "username already exists", http.StatusConflict)
		return
	}
	err = handler.AuthService.UserRepo.UpdateUserProfile(userReq)
	if err != nil {
		log.Println(err)
		http.Error(w, "error occured 2", http.StatusInternalServerError)
		return
	}
	userPro, err := handler.AuthService.UserRepo.GetUserProfileByEmail(userReq.Email)
	if err != nil {
		http.Error(w, "error occured 3", http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, http.StatusOK, userPro)
	// err=handler.AuthService.UserRepo.
}

func (handler *UserHandler) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	var userReq models.NewRegRequest
	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		http.Error(w, "invalid format", http.StatusBadRequest)
		return
	}
	user, err := handler.AuthService.SignIn(userReq.Email, userReq.Password)
	if err != nil {
		// fmt.Println(err)
		http.Error(w, "wrong email or password", http.StatusUnauthorized)
		return
	}
	token, err := middleware.CreateJWT(user.Email, user.CreatedAt, user.ID)
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	utils.SetTokenInCookies(w, token)
	utils.WriteJSON(w, http.StatusOK, &models.LogInResponse{
		Email: user.Email,
	})
}

func (handler *UserHandler) GetUserProfile(w http.ResponseWriter, r *http.Request, jwtToken *jwt.Token) {
	var email string
	err := json.NewDecoder(r.Body).Decode(&email)
	if err != nil {
		http.Error(w, "invalid format", http.StatusBadRequest)
		return
	}
	emailFromJWT, err := utils.ExtractEmailFromJWT(jwtToken)
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}
	if emailFromJWT != email {
		http.Error(w, "unauthorised", http.StatusUnauthorized)
		return
	}
	userProfile, err := handler.AuthService.UserRepo.GetUserProfileByEmail(email)
	if err != nil {
		http.Error(w, "not found", http.StatusNoContent)
		return
	}

	utils.WriteJSON(w, http.StatusOK, userProfile)

}
