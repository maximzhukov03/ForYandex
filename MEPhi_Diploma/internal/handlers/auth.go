package handler

import (
    "net/http"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"

    "golang/myapp/internal/service"
)

type AuthHandler struct{
    userService service.UserService
}

func NewAuthHandler(us service.UserService) *AuthHandler{
    return &AuthHandler{userService: us}
}

type RegisterRequest struct{
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct{
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type LoginResponse struct{
    Token string `json:"token"`
}

// Register godoc
// @Summary      Регистрация пользователя
// @Description  Создаёт нового пользователя
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        input  body      RegisterRequest  true  "Данные регистрации"
// @Success      201    {object}  models.User
// @Failure      400    {object}  ErrorResponse
// @Failure      409    {object}  ErrorResponse
// @Router       /api/register [post]
func (h *AuthHandler) Register(c *gin.Context){
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil{
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    user, err := h.userService.Register(req.Email, req.Password)
    if err != nil {
        if err.Error() == "user already exists"{
            c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }
    c.JSON(http.StatusCreated, user)
}

// Login godoc
// @Summary      Аутентификация пользователя
// @Description  Возвращает JWT при правильных учётных данных
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        input  body      LoginRequest   true  "Данные для логина"
// @Success      200    {object}  LoginResponse
// @Failure      400    {object}  ErrorResponse
// @Failure      401    {object}  ErrorResponse
// @Router       /api/login [post]⁡
func (h *AuthHandler) Login(c *gin.Context){
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    user, err := h.userService.Authenticate(req.Email, req.Password)
    if err != nil{
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }
    secret := os.Getenv("JWT_SECRET")
    claims := jwt.MapClaims{
        "user_id": user.ID,
        "exp":     time.Now().Add(time.Hour * 72).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenStr, err := token.SignedString([]byte(secret))
    if err != nil{
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, LoginResponse{Token: tokenStr})
}
