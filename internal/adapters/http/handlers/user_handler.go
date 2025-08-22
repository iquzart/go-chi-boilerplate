package handlers

import (
	"encoding/json"
	"go-chi-boilerplate/internal/adapters/http/dto"
	"go-chi-boilerplate/internal/app/usecases/user"
	"go-chi-boilerplate/internal/core/entities"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// CreateUserHandler godoc
// @Summary      Create a new user
// @Description  Create a new user with first name, last name, email, role, and password. Role must be one of: admin, user, maker, checker
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user body dto.CreateUserRequest true "User info (role: admin, user, maker, checker)"
// @Success      201 {object} dto.UserResponse
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /api/v1/users [post]
func CreateUserHandler(uc *user.UserUsecase) http.HandlerFunc {
	validRoles := map[string]bool{
		"admin":   true,
		"user":    true,
		"maker":   true,
		"checker": true,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
			return
		}

		// Validate role
		if !validRoles[req.Role] {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid role, must be one of: admin, user, maker, checker"})
			return
		}

		entity := &entities.User{
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Email:     req.Email,
			Role:      req.Role,
			Password:  req.Password,
			Status:    "active",
		}

		created, err := uc.CreateUser(entity)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		resp := dto.UserResponse{
			ID:        created.ID,
			FirstName: created.FirstName,
			LastName:  created.LastName,
			Email:     created.Email,
			Role:      created.Role,
			Status:    created.Status,
			CreatedAt: created.CreatedAt,
			UpdatedAt: created.UpdatedAt,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}

// ListUsersHandler godoc
// @Summary      List users
// @Description  Get a paginated list of users
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        offset query int false "Offset"
// @Param        limit query int false "Limit"
// @Success      200 {object} dto.ListUsersResponse
// @Failure      500 {object} map[string]string
// @Router       /api/v1/users [get]
func ListUsersHandler(uc *user.UserUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		offset := 0
		limit := 10

		users, total, err := uc.ListUsers(offset, limit)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		resp := dto.ListUsersResponse{
			Users: []*dto.UserResponse{},
			Total: total,
		}
		for _, u := range users {
			resp.Users = append(resp.Users, &dto.UserResponse{
				ID:        u.ID,
				FirstName: u.FirstName,
				LastName:  u.LastName,
				Email:     u.Email,
				Role:      u.Role,
				Status:    u.Status,
				CreatedAt: u.CreatedAt,
				UpdatedAt: u.UpdatedAt,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

// GetUserHandler godoc
// @Summary      Get user by ID
// @Description  Get a single user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID"
// @Success      200 {object} dto.UserResponse
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /api/v1/users/{id} [get]
func GetUserHandler(uc *user.UserUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		user, err := uc.GetUserByID(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		resp := dto.UserResponse{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Role:      user.Role,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

// UpdateUserHandler godoc
// @Summary      Update user
// @Description  Update user fields (role must be one of: admin, user, maker, checker)
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID"
// @Param        user body dto.UpdateUserRequest true "Updated user info (role: admin, user, maker, checker)"
// @Success      200 {object} dto.UserResponse
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /api/v1/users/{id} [put]
func UpdateUserHandler(uc *user.UserUsecase) http.HandlerFunc {
	validRoles := map[string]bool{
		"admin":   true,
		"user":    true,
		"maker":   true,
		"checker": true,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var req dto.UpdateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
			return
		}

		userEntity, err := uc.GetUserByID(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
			return
		}

		if req.FirstName != nil {
			userEntity.FirstName = *req.FirstName
		}
		if req.LastName != nil {
			userEntity.LastName = *req.LastName
		}
		if req.Email != nil {
			userEntity.Email = *req.Email
		}
		if req.Role != nil {
			if !validRoles[*req.Role] {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "invalid role, must be one of: admin, user, maker, checker"})
				return
			}
			userEntity.Role = *req.Role
		}

		updated, err := uc.UpdateUser(userEntity)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		resp := dto.UserResponse{
			ID:        updated.ID,
			FirstName: updated.FirstName,
			LastName:  updated.LastName,
			Email:     updated.Email,
			Role:      updated.Role,
			Status:    updated.Status,
			CreatedAt: updated.CreatedAt,
			UpdatedAt: updated.UpdatedAt,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

// ChangeUserStatusHandler godoc
// @Summary      Change user status
// @Description  Activate or deactivate a user
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID"
// @Param        status body dto.ChangeUserStatusRequest true "Status change (active or inactive)"
// @Success      200 {object} dto.UserResponse
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /api/v1/users/{id}/status [patch]
func ChangeUserStatusHandler(uc *user.UserUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var req dto.ChangeUserStatusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
			return
		}

		// Validate status explicitly
		if req.Status != "active" && req.Status != "inactive" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "status must be 'active' or 'inactive'"})
			return
		}

		updated, err := uc.ChangeStatus(id, req.Status)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		resp := dto.UserResponse{
			ID:        updated.ID,
			FirstName: updated.FirstName,
			LastName:  updated.LastName,
			Email:     updated.Email,
			Role:      updated.Role,
			Status:    updated.Status,
			CreatedAt: updated.CreatedAt,
			UpdatedAt: updated.UpdatedAt,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

// DeleteUserHandler godoc
// @Summary      Delete user
// @Description  Delete a user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID"
// @Success      204
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /api/v1/users/{id} [delete]
func DeleteUserHandler(uc *user.UserUsecase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if err := uc.DeleteUser(id); err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
