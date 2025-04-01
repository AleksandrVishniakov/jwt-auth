package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/AleksandrVishniakov/jwt-auth/internal/e"
	"github.com/AleksandrVishniakov/jwt-auth/internal/usecases"
)

type contextKey string

const (
	userIDKey         contextKey = "userID"
	roleKey           contextKey = "role"
	permissionMaskKey contextKey = "permissionMask"
)

type Usecase interface {
	Login(
		ctx context.Context,
		req *usecases.LoginRequest,
	) (resp *usecases.LoginResponse, err error)

	Register(
		ctx context.Context,
		req *usecases.RegisterRequest,
	) (resp *usecases.RegisterResponse, err error)

	UpdateRoleById(
		ctx context.Context,
		req *usecases.UpdateUserRoleRequest,
	) (err error)

	GetUserByID(
		ctx context.Context,
		req *usecases.GetUserByIDRequest,
	) (*usecases.GetUserByIDResponse, error)
}

type Handler struct {
	log         *slog.Logger
	usecase     Usecase
	tokenParser TokenParser
}

func New(
	log *slog.Logger,
	usecase Usecase,
	tokenParser TokenParser,
) *Handler {
	return &Handler{
		log:         log,
		usecase:     usecase,
		tokenParser: tokenParser,
	}
}

func (h *Handler) InitRoutes() http.Handler {
	logger := Logger(h.log)
	jwt := JWTAuth(h.log, h.tokenParser)

	mux := http.NewServeMux()
	v1 := http.NewServeMux()

	v1.Handle("GET /ping", Error(h.Ping))
	v1.Handle("POST /login", Error(h.Login))
	v1.Handle("POST /register", Error(h.Register))
	v1.Handle("PUT /change-role", jwt(Error(h.ChangeRole)))
	v1.Handle("GET /user/{id}", jwt(Error(h.GetUser)))

	mux.Handle("/v1/", http.StripPrefix("/v1", v1))

	return http.StripPrefix("/api", logger(CORS(mux)))
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()

	type loginRequset struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	req, err := Decode[loginRequset](r.Body)
	if err != nil {
		return e.BadRequest(e.WithError(err))
	}

	dto, err := usecases.NewLoginRequest(
		req.Login,
		req.Password,
	)
	if err != nil {
		return e.BadRequest(e.WithError(err))
	}

	resp, err := h.usecase.Login(r.Context(), dto)
	if err != nil {
		if errors.Is(err, e.ErrNotFound) {
			return e.NotFound()
		}

		return e.Authorization(e.WithError(err))
	}

	_ = EncodeResponse(w, struct {
		ID    int32  `json:"id"`
		Token string `json:"token"`
	}{
		ID:    resp.ID,
		Token: resp.Token,
	}, http.StatusOK)

	return nil
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()

	type registerRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	req, err := Decode[registerRequest](r.Body)
	if err != nil {
		return e.BadRequest(e.WithError(err))
	}

	dto, err := usecases.NewRegisterRequest(
		req.Login,
		req.Password,
	)
	if err != nil {
		return e.BadRequest(e.WithError(err))
	}

	resp, err := h.usecase.Register(r.Context(), dto)
	if err != nil {
		if errors.Is(err, e.ErrAlreadyExists) {
			return e.BadRequest(e.WithMessage("already exists"))
		}

		return e.Authorization(e.WithError(err))
	}

	_ = EncodeResponse(w, struct {
		ID    int32  `json:"id"`
		Token string `json:"token"`
	}{
		ID:    resp.ID,
		Token: resp.Token,
	}, http.StatusOK)

	return nil
}

func (h *Handler) ChangeRole(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()

	type changeRoleRequest struct {
		UserID  int32  `json:"userID"`
		NewRole string `json:"newRole"`
	}

	req, err := Decode[changeRoleRequest](r.Body)
	if err != nil {
		return e.BadRequest(e.WithError(err))
	}

	mask, err := PermissionMaskFromContext(r.Context())
	if err != nil {
		return e.Authorization()
	}

	dto, err := usecases.NewUpdateUserRoleRequest(
		req.UserID,
		req.NewRole,
		mask,
	)
	if err != nil {
		return e.BadRequest(e.WithError(err))
	}

	err = h.usecase.UpdateRoleById(r.Context(), dto)
	if err != nil {
		if errors.Is(err, e.ErrForbiddenAction) {
			return e.Forbidden()
		}

		return e.Internal(e.WithError(err))
	}

	return nil
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()

	mask, err := PermissionMaskFromContext(r.Context())
	if err != nil {
		return e.Authorization()
	}

	userID, err := UserIDFromContext(r.Context())
	if err != nil {
		return e.Authorization()
	}

	profileID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return e.BadRequest()
	}

	dto, err := usecases.NewGetUserByIDRequest(
		userID,
		mask,
		int32(profileID),
	)
	if err != nil {
		return e.BadRequest(e.WithError(err))
	}

	resp, err := h.usecase.GetUserByID(r.Context(), dto)
	if err != nil {
		if errors.Is(err, e.ErrForbiddenAction) {
			return e.Forbidden()
		}

		if errors.Is(err, e.ErrNotFound) {
			return e.NotFound()
		}

		return e.Internal(e.WithError(err))
	}

	return EncodeResponse(w, &struct {
		ID    int32  `json:"id"`
		Login string `json:"login"`
		Role  string `json:"role"`
	}{
		ID:    resp.ID,
		Login: resp.Login,
		Role:  resp.Role,
	}, http.StatusOK)
}
