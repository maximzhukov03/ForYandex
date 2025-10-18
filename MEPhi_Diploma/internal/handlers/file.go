package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "golang/myapp/internal/service"
)

type FileHandler struct{
    fileService service.FileService
}

func NewFileHandler(fs service.FileService) *FileHandler{
    return &FileHandler{fileService: fs}
}

// UploadFile godoc
// @Summary      Загрузка файла
// @Description  Загружает png/jpeg файл до 10 МБ
// @Tags         files
// @Accept       multipart/form-data
// @Produce      json
// @Security     bearerAuth
// @Param        Authorization header    string       true  "Bearer JWT"
// @Param        file           formData  file         true  "Файл (png/jpeg)"
// @Success      201            {object}  models.File
// @Failure      400            {object}  ErrorResponse
// @Failure      401            {object}  ErrorResponse
// @Failure      413            {object}  ErrorResponse  "File too large"
// @Router       /api/files [post]
func (h *FileHandler) UploadFile(c *gin.Context){
    userID := c.GetInt64("userID")
    fileHeader, err := c.FormFile("file")
    if err != nil{
        c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
        return
    }
    if fileHeader.Size > 10<<20{
        c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "file too large"})
        return
    }
    f, err := fileHeader.Open()
    if err != nil{
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer f.Close()
    file, err := h.fileService.Upload(userID, fileHeader.Filename, f, fileHeader.Size, fileHeader.Header.Get("Content-Type"))
    if err != nil{
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, file)
}

// ListFiles godoc
// @Summary      Список файлов пользователя
// @Description  Возвращает список загруженных файлов с presigned URL
// @Tags         files
// @Produce      json
// @Security     bearerAuth
// @Param        Authorization header string true "Bearer JWT"
// @Success      200            {array}   models.File
// @Failure      401            {object}  ErrorResponse
// @Router       /api/files [get]
func (h *FileHandler) ListFiles(c *gin.Context){
    userID := c.GetInt64("userID")
    files, err := h.fileService.List(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, files)
}

type ErrorResponse struct{
    Error string `json:"error"`
}
