package handler

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"

    "golang/myapp/internal/service"
)

type AdminHandler struct{
    adminService service.AdminService
}

func NewAdminHandler(as service.AdminService) *AdminHandler{
    return &AdminHandler{adminService: as}
}

// ListUsers godoc
// @Summary      Получить всех пользователей
// @Description  Возвращает список всех пользователей (только для администратора)
// @Tags         admin
// @Produce      json
// @Security     bearerAuth
// @Success      200  {array}   models.User
// @Failure      401  {object}  ErrorResponse
// @Router       /api/admin/users [get]
func (h *AdminHandler) ListUsers(c *gin.Context){
    users, err := h.adminService.ListUsers()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, users)
}

// GetUser godoc
// @Summary      Получить пользователя по ID
// @Description  Возвращает информацию о пользователе
// @Tags         admin
// @Produce      json
// @Security     bearerAuth
// @Param        id   path      int  true  "ID пользователя"
// @Success      200  {object}  models.User
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /api/admin/users/{id} [get]
func (h *AdminHandler) GetUser(c *gin.Context){
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil{
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
        return
    }
    user, err := h.adminService.GetUserByID(id)
    if err != nil{
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    if user == nil{
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }
    c.JSON(http.StatusOK, user)
}

// UpdateUserRequest модель запроса обновления пользователя
// содержит Email и Role
// swagger:model
// swagger:parameters UpdateUser
type UpdateUserRequest struct{
    // Email пользователя
    // required: true
    Email string `json:"email" binding:"required,email"`
    // Роль пользователя (user или admin)
    // required: true
    Role  string `json:"role" binding:"required,oneof=user admin"`
}

// UpdateUser godoc
// @Summary      Обновить пользователя
// @Description  Изменяет email и роль пользователя
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     bearerAuth
// @Param        id     path      int               true  "ID пользователя"
// @Param        input  body      UpdateUserRequest true  "Данные для обновления"
// @Success      200    {object}  models.User
// @Failure      400    {object}  ErrorResponse
// @Failure      401    {object}  ErrorResponse
// @Failure      404    {object}  ErrorResponse
// @Router       /api/admin/users/{id} [put]
func (h *AdminHandler) UpdateUser(c *gin.Context){
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil{
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
        return
    }
    var req UpdateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil{
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    user, err := h.adminService.UpdateUser(id, req.Email, req.Role)
    if err != nil{
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    if user == nil{
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }
    c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary      Удалить пользователя
// @Description  Удаляет пользователя по ID
// @Tags         admin
// @Produce      json
// @Security     bearerAuth
// @Param        id   path      int  true  "ID пользователя"
// @Success      204  {object}  nil
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /api/admin/users/{id} [delete]
func (h *AdminHandler) DeleteUser(c *gin.Context){
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil{
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
        return
    }
    deleted, err := h.adminService.DeleteUser(id)
    if err != nil{
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    if !deleted{
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }
    c.Status(http.StatusNoContent)
}

// PromoteUser godoc
// @Summary      Повысить пользователя до администратора
// @Description  Назначает роль admin пользователю по ID
// @Tags         admin
// @Produce      json
// @Security     bearerAuth
// @Param        id   path      int  true  "ID пользователя"
// @Success      200  {object}  models.User
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /api/admin/users/{id}/promote [post]
func (h *AdminHandler) PromoteUser(c *gin.Context){
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil{
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
        return
    }
    user, err := h.adminService.PromoteUser(id)
    if err != nil{
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    if user == nil{
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }
    c.JSON(http.StatusOK, user)
}
