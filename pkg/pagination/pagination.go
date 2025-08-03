package pagination

import (
	"math"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// PaginationParams represents pagination parameters from query string
type PaginationParams struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Sort   string `json:"sort"`
	Order  string `json:"order"`
	Search string `json:"search"`
}

// PaginationMeta represents pagination metadata in response
type PaginationMeta struct {
	CurrentPage  int  `json:"current_page"`
	PerPage      int  `json:"per_page"`
	TotalRecords int  `json:"total_records"`
	TotalPages   int  `json:"total_pages"`
	HasNext      bool `json:"has_next"`
	HasPrev      bool `json:"has_prev"`
}

// FilterParams represents filtering parameters for users
type FilterParams struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	CreatedAfter string `json:"created_after"`
	CreatedBefore string `json:"created_before"`
	EmailVerified *bool  `json:"email_verified"`
}

// FileFilterParams represents filtering parameters for files
type FileFilterParams struct {
	FileName      string `json:"file_name"`
	MimeType      string `json:"mime_type"`
	Category      string `json:"category"`
	UploadedBy    *int   `json:"uploaded_by"`
	CreatedAfter  string `json:"created_after"`
	CreatedBefore string `json:"created_before"`
}

// GetPaginationParams extracts pagination parameters from Echo context
func GetPaginationParams(c echo.Context) PaginationParams {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 { // Max 100 records per page
		limit = 10
	}

	sort := c.QueryParam("sort")
	if sort == "" {
		sort = "id" // Default sort by ID
	}

	order := strings.ToUpper(c.QueryParam("order"))
	if order != "ASC" && order != "DESC" {
		order = "DESC" // Default descending order
	}

	search := strings.TrimSpace(c.QueryParam("search"))

	return PaginationParams{
		Page:   page,
		Limit:  limit,
		Sort:   sort,
		Order:  order,
		Search: search,
	}
}

// GetFilterParams extracts filter parameters from Echo context
func GetFilterParams(c echo.Context) FilterParams {
	var emailVerified *bool
	if emailVerifiedStr := c.QueryParam("email_verified"); emailVerifiedStr != "" {
		if val, err := strconv.ParseBool(emailVerifiedStr); err == nil {
			emailVerified = &val
		}
	}

	return FilterParams{
		Name:          strings.TrimSpace(c.QueryParam("name")),
		Email:         strings.TrimSpace(c.QueryParam("email")),
		Role:          strings.TrimSpace(c.QueryParam("role")),
		CreatedAfter:  strings.TrimSpace(c.QueryParam("created_after")),
		CreatedBefore: strings.TrimSpace(c.QueryParam("created_before")),
		EmailVerified: emailVerified,
	}
}

// GetFileFilterParams extracts file filter parameters from Echo context
func GetFileFilterParams(c echo.Context) FileFilterParams {
	var uploadedBy *int
	if uploadedByStr := c.QueryParam("uploaded_by"); uploadedByStr != "" {
		if val, err := strconv.Atoi(uploadedByStr); err == nil {
			uploadedBy = &val
		}
	}

	return FileFilterParams{
		FileName:      strings.TrimSpace(c.QueryParam("file_name")),
		MimeType:      strings.TrimSpace(c.QueryParam("mime_type")),
		Category:      strings.TrimSpace(c.QueryParam("category")),
		UploadedBy:    uploadedBy,
		CreatedAfter:  strings.TrimSpace(c.QueryParam("created_after")),
		CreatedBefore: strings.TrimSpace(c.QueryParam("created_before")),
	}
}

// CalculateOffset calculates the offset for database queries
func (p PaginationParams) CalculateOffset() int {
	return (p.Page - 1) * p.Limit
}

// NewPaginationMeta creates pagination metadata
func NewPaginationMeta(page, limit, totalRecords int) PaginationMeta {
	totalPages := int(math.Ceil(float64(totalRecords) / float64(limit)))
	
	return PaginationMeta{
		CurrentPage:  page,
		PerPage:      limit,
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
		HasNext:      page < totalPages,
		HasPrev:      page > 1,
	}
}

// ValidateSortField validates if the sort field is allowed for security
func ValidateSortField(field string, allowedFields []string) string {
	field = strings.ToLower(field)
	for _, allowed := range allowedFields {
		if field == strings.ToLower(allowed) {
			return allowed
		}
	}
	return "id" // Default fallback
}

// SanitizeOrder ensures order is either ASC or DESC
func SanitizeOrder(order string) string {
	order = strings.ToUpper(order)
	if order == "ASC" || order == "DESC" {
		return order
	}
	return "DESC" // Default fallback
}