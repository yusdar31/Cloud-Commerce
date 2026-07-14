package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuccessResponse is the standard API success envelope.
type SuccessResponse struct {
	Data  interface{} `json:"data"`
	Meta  *Meta       `json:"meta,omitempty"`
	Links *Links      `json:"links,omitempty"`
}

// ErrorResponse follows RFC 9457 Problem Details for HTTP APIs.
type ErrorResponse struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance,omitempty"`
	TraceID  string `json:"traceId,omitempty"`
}

// Meta holds pagination and other metadata.
type Meta struct {
	NextCursor string `json:"nextCursor,omitempty"`
	HasMore    bool   `json:"hasMore,omitempty"`
	Total      int64  `json:"total,omitempty"`
	Page       int    `json:"page,omitempty"`
	PerPage    int    `json:"perPage,omitempty"`
}

// Links holds pagination links.
type Links struct {
	Self  string `json:"self,omitempty"`
	Next  string `json:"next,omitempty"`
	Prev  string `json:"prev,omitempty"`
	First string `json:"first,omitempty"`
	Last  string `json:"last,omitempty"`
}

// OK sends a 200 with data.
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, SuccessResponse{Data: data})
}

// Created sends a 201 with data.
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, SuccessResponse{Data: data})
}

// OKWithMeta sends a 200 with data and meta.
func OKWithMeta(c *gin.Context, data interface{}, meta *Meta) {
	c.JSON(http.StatusOK, SuccessResponse{Data: data, Meta: meta})
}

// NoContent sends a 204.
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Error sends an RFC 9457 problem detail response.
func Error(c *gin.Context, status int, title, detail string) {
	traceID := c.GetString("request_id")
	c.JSON(status, ErrorResponse{
		Type:     "https://api.cloudcommerce.com/errors/" + title,
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: c.Request.URL.Path,
		TraceID:  traceID,
	})
}

// ValidationError sends a 422 validation error.
func ValidationError(c *gin.Context, detail string) {
	Error(c, http.StatusUnprocessableEntity, "Validation Failed", detail)
}

// NotFound sends a 404.
func NotFound(c *gin.Context, detail string) {
	Error(c, http.StatusNotFound, "Not Found", detail)
}

// Unauthorized sends a 401.
func Unauthorized(c *gin.Context, detail string) {
	Error(c, http.StatusUnauthorized, "Unauthorized", detail)
}

// Forbidden sends a 403.
func Forbidden(c *gin.Context, detail string) {
	Error(c, http.StatusForbidden, "Forbidden", detail)
}

// Conflict sends a 409.
func Conflict(c *gin.Context, detail string) {
	Error(c, http.StatusConflict, "Conflict", detail)
}

// InternalError sends a 500.
func InternalError(c *gin.Context, detail string) {
	Error(c, http.StatusInternalServerError, "Internal Server Error", detail)
}
