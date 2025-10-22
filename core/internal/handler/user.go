package handler

import (
	"fmt"
	"net/http"

	"github.com/pikachu0310/go-backend-template/core/internal/repository"
	"github.com/pikachu0310/go-backend-template/core/internal/services/photoapi"

	vd "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// スキーマ定義
type (
	GetUsersResponse []GetUserResponse

	GetUserResponse struct {
		ID      uuid.UUID `json:"id"`
		Name    string    `json:"name"`
		Email   string    `json:"email"`
		IconURL string    `json:"iconUrl"`
	}

	CreateUserRequest struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	CreateUserResponse struct {
		ID uuid.UUID `json:"id"`
	}

	ErrorResponse struct {
		Message string `json:"message"`
	}
)

// GetUsers godoc
// @Summary ユーザー一覧取得
// @Description 全ユーザーの情報を取得します
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} GetUserResponse "ユーザー一覧"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /users [get]
func (h *Handler) GetUsers(c echo.Context) error {
	users, err := h.repo.GetUsers(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	photo, err := photoapi.GetPhoto()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	res := make(GetUsersResponse, len(users))
	for i, user := range users {
		res[i] = GetUserResponse{
			ID:      user.ID,
			Name:    user.Name,
			Email:   user.Email,
			IconURL: photo.ThumbnailURL,
		}
	}

	return c.JSON(http.StatusOK, res)
}

// CreateUser godoc
// @Summary ユーザー作成
// @Description 新しいユーザーを作成します
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "ユーザー情報"
// @Success 200 {object} CreateUserResponse "作成されたユーザーのID"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /users [post]
func (h *Handler) CreateUser(c echo.Context) error {
	req := new(CreateUserRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").SetInternal(err)
	}

	err := vd.ValidateStruct(
		req,
		vd.Field(&req.Name, vd.Required),
		vd.Field(&req.Email, vd.Required, is.Email),
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err)).SetInternal(err)
	}

	userID, err := h.repo.CreateUser(c.Request().Context(), repository.CreateUserParams{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	res := CreateUserResponse{
		ID: userID,
	}

	return c.JSON(http.StatusOK, res)
}

// GetUser godoc
// @Summary ユーザー情報取得
// @Description 指定したIDのユーザー情報を取得します
// @Tags users
// @Accept json
// @Produce json
// @Param userID path string true "User ID" format(uuid)
// @Success 200 {object} GetUserResponse "ユーザー情報"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /users/{userID} [get]
func (h *Handler) GetUser(c echo.Context) error {
	userID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid userID").SetInternal(err)
	}

	user, err := h.repo.GetUser(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	photo, err := photoapi.GetPhoto()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	res := GetUserResponse{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		IconURL: photo.ThumbnailURL,
	}

	return c.JSON(http.StatusOK, res)
}
