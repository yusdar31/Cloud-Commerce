package transport

import (
	"errors"

	"github.com/cloudcommerce/shared-go/response"
	"github.com/cloudcommerce/store-service/internal/application"
	"github.com/cloudcommerce/store-service/internal/domain"
	"github.com/gin-gonic/gin"
)

// StoreHandler handles HTTP requests for the Store/Tenant endpoints.
type StoreHandler struct {
	storeService *application.StoreService
}

// NewStoreHandler creates a new StoreHandler.
func NewStoreHandler(storeService *application.StoreService) *StoreHandler {
	return &StoreHandler{storeService: storeService}
}

// Create handles POST /api/v1/stores.
func (h *StoreHandler) Create(c *gin.Context) {
	ownerID := c.GetString("user_id")
	if ownerID == "" {
		response.Unauthorized(c, "Authentication required")
		return
	}

	var req CreateStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	if err := Validate(req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	tenant, err := h.storeService.CreateStore(c.Request.Context(), application.CreateStoreInput{
		Name:        req.Name,
		Slug:        req.Slug,
		OwnerID:     ownerID,
		Description: req.Description,
		LogoURL:     req.LogoURL,
		BannerURL:   req.BannerURL,
	})
	if err != nil {
		if errors.Is(err, domain.ErrSlugAlreadyExists) {
			response.Conflict(c, "Store slug/subdomain is already taken")
			return
		}
		if errors.Is(err, domain.ErrOwnerAlreadyHasStore) {
			response.Conflict(c, "You already own a store and cannot create another one")
			return
		}
		response.InternalError(c, "Failed to create store: "+err.Error())
		return
	}

	response.Created(c, ToStoreResponse(tenant))
}

// GetByID handles GET /api/v1/stores/:id.
func (h *StoreHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.ValidationError(c, "Store ID is required")
		return
	}

	tenant, err := h.storeService.GetStoreByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrTenantNotFound) {
			response.NotFound(c, "Store not found")
			return
		}
		response.InternalError(c, "Failed to get store: "+err.Error())
		return
	}

	response.OK(c, ToStoreResponse(tenant))
}

// GetBySlug handles GET /api/v1/stores/slug/:slug.
func (h *StoreHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		response.ValidationError(c, "Store slug is required")
		return
	}

	tenant, err := h.storeService.GetStoreBySlug(c.Request.Context(), slug)
	if err != nil {
		if errors.Is(err, domain.ErrTenantNotFound) {
			response.NotFound(c, "Store not found")
			return
		}
		response.InternalError(c, "Failed to get store: "+err.Error())
		return
	}

	response.OK(c, ToStoreResponse(tenant))
}

// GetMe handles GET /api/v1/stores/me.
func (h *StoreHandler) GetMe(c *gin.Context) {
	ownerID := c.GetString("user_id")
	if ownerID == "" {
		response.Unauthorized(c, "Authentication required")
		return
	}

	tenant, err := h.storeService.GetStoreByOwnerID(c.Request.Context(), ownerID)
	if err != nil {
		if errors.Is(err, domain.ErrTenantNotFound) {
			response.NotFound(c, "You do not own any store yet")
			return
		}
		response.InternalError(c, "Failed to get your store: "+err.Error())
		return
	}

	response.OK(c, ToStoreResponse(tenant))
}

// Update handles PUT /api/v1/stores/me.
func (h *StoreHandler) Update(c *gin.Context) {
	ownerID := c.GetString("user_id")
	if ownerID == "" {
		response.Unauthorized(c, "Authentication required")
		return
	}

	tenant, err := h.storeService.GetStoreByOwnerID(c.Request.Context(), ownerID)
	if err != nil {
		if errors.Is(err, domain.ErrTenantNotFound) {
			response.NotFound(c, "Store not found")
			return
		}
		response.InternalError(c, "Failed to locate store: "+err.Error())
		return
	}

	var req UpdateStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	if err := Validate(req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	updatedTenant, err := h.storeService.UpdateStore(c.Request.Context(), tenant.ID, ownerID, application.UpdateStoreInput{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		LogoURL:     req.LogoURL,
		BannerURL:   req.BannerURL,
		Status:      req.Status,
	})
	if err != nil {
		if errors.Is(err, domain.ErrSlugAlreadyExists) {
			response.Conflict(c, "Store slug/subdomain is already taken")
			return
		}
		if errors.Is(err, domain.ErrInvalidTenantStatus) {
			response.ValidationError(c, "Invalid store status")
			return
		}
		response.InternalError(c, "Failed to update store: "+err.Error())
		return
	}

	response.OK(c, ToStoreResponse(updatedTenant))
}
