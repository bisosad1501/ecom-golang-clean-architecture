package usecases

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// CMSUseCase defines content management use cases
type CMSUseCase interface {
	// Pages Management
	GetPages(ctx context.Context, req GetPagesRequest) (*GetPagesResponse, error)
	GetPage(ctx context.Context, id uuid.UUID) (*PageResponse, error)
	GetPageBySlug(ctx context.Context, slug string) (*PageResponse, error)
	CreatePage(ctx context.Context, req CreatePageRequest) (*PageResponse, error)
	UpdatePage(ctx context.Context, id uuid.UUID, req UpdatePageRequest) (*PageResponse, error)
	DeletePage(ctx context.Context, id uuid.UUID) error
	PublishPage(ctx context.Context, id uuid.UUID) error
	UnpublishPage(ctx context.Context, id uuid.UUID) error

	// Blocks Management
	GetBlocks(ctx context.Context, req GetBlocksRequest) (*GetBlocksResponse, error)
	GetBlock(ctx context.Context, id uuid.UUID) (*BlockResponse, error)
	CreateBlock(ctx context.Context, req CreateBlockRequest) (*BlockResponse, error)
	UpdateBlock(ctx context.Context, id uuid.UUID, req UpdateBlockRequest) (*BlockResponse, error)
	DeleteBlock(ctx context.Context, id uuid.UUID) error

	// Menus Management
	GetMenus(ctx context.Context) (*GetMenusResponse, error)
	GetMenu(ctx context.Context, id uuid.UUID) (*MenuResponse, error)
	CreateMenu(ctx context.Context, req CreateMenuRequest) (*MenuResponse, error)
	UpdateMenu(ctx context.Context, id uuid.UUID, req UpdateMenuRequest) (*MenuResponse, error)
	DeleteMenu(ctx context.Context, id uuid.UUID) error

	// Banners Management
	GetBanners(ctx context.Context, req GetBannersRequest) (*GetBannersResponse, error)
	GetBanner(ctx context.Context, id uuid.UUID) (*BannerResponse, error)
	CreateBanner(ctx context.Context, req CreateBannerRequest) (*BannerResponse, error)
	UpdateBanner(ctx context.Context, id uuid.UUID, req UpdateBannerRequest) (*BannerResponse, error)
	DeleteBanner(ctx context.Context, id uuid.UUID) error

	// Sliders Management
	GetSliders(ctx context.Context, req GetSlidersRequest) (*GetSlidersResponse, error)
	GetSlider(ctx context.Context, id uuid.UUID) (*SliderResponse, error)
	CreateSlider(ctx context.Context, req CreateSliderRequest) (*SliderResponse, error)
	UpdateSlider(ctx context.Context, id uuid.UUID, req UpdateSliderRequest) (*SliderResponse, error)
	DeleteSlider(ctx context.Context, id uuid.UUID) error

	// Widgets Management
	GetWidgets(ctx context.Context, req GetWidgetsRequest) (*GetWidgetsResponse, error)
	GetWidget(ctx context.Context, id uuid.UUID) (*WidgetResponse, error)
	CreateWidget(ctx context.Context, req CreateWidgetRequest) (*WidgetResponse, error)
	UpdateWidget(ctx context.Context, id uuid.UUID, req UpdateWidgetRequest) (*WidgetResponse, error)
	DeleteWidget(ctx context.Context, id uuid.UUID) error

	// Templates Management
	GetTemplates(ctx context.Context) (*GetTemplatesResponse, error)
	GetTemplate(ctx context.Context, id uuid.UUID) (*TemplateResponse, error)
	CreateTemplate(ctx context.Context, req CreateTemplateRequest) (*TemplateResponse, error)
	UpdateTemplate(ctx context.Context, id uuid.UUID, req UpdateTemplateRequest) (*TemplateResponse, error)
	DeleteTemplate(ctx context.Context, id uuid.UUID) error
}

// Page types
type GetPagesRequest struct {
	Page     int    `json:"page" validate:"min=1"`
	Limit    int    `json:"limit" validate:"min=1,max=100"`
	Search   string `json:"search"`
	Status   string `json:"status" validate:"omitempty,oneof=draft published archived"`
	Template string `json:"template"`
	SortBy   string `json:"sort_by" validate:"omitempty,oneof=title created_at updated_at"`
	SortOrder string `json:"sort_order" validate:"omitempty,oneof=asc desc"`
}

type CreatePageRequest struct {
	Title            string                 `json:"title" validate:"required,min=1,max=200"`
	Slug             string                 `json:"slug" validate:"required,min=1,max=200"`
	Content          string                 `json:"content"`
	Excerpt          string                 `json:"excerpt" validate:"max=500"`
	Status           string                 `json:"status" validate:"required,oneof=draft published archived"`
	Template         string                 `json:"template" validate:"required"`
	MetaTitle        string                 `json:"meta_title" validate:"max=60"`
	MetaDescription  string                 `json:"meta_description" validate:"max=160"`
	MetaKeywords     []string               `json:"meta_keywords"`
	FeaturedImage    string                 `json:"featured_image"`
	ParentID         *uuid.UUID             `json:"parent_id"`
	MenuOrder        int                    `json:"menu_order"`
	ShowInMenu       bool                   `json:"show_in_menu"`
	RequireAuth      bool                   `json:"require_auth"`
	AllowComments    bool                   `json:"allow_comments"`
	CustomFields     map[string]interface{} `json:"custom_fields"`
	PublishedAt      *time.Time             `json:"published_at"`
	ScheduledAt      *time.Time             `json:"scheduled_at"`
	AuthorID         uuid.UUID              `json:"author_id" validate:"required"`
}

type UpdatePageRequest struct {
	Title            string                 `json:"title" validate:"min=1,max=200"`
	Slug             string                 `json:"slug" validate:"min=1,max=200"`
	Content          string                 `json:"content"`
	Excerpt          string                 `json:"excerpt" validate:"max=500"`
	Status           string                 `json:"status" validate:"omitempty,oneof=draft published archived"`
	Template         string                 `json:"template"`
	MetaTitle        string                 `json:"meta_title" validate:"max=60"`
	MetaDescription  string                 `json:"meta_description" validate:"max=160"`
	MetaKeywords     []string               `json:"meta_keywords"`
	FeaturedImage    string                 `json:"featured_image"`
	ParentID         *uuid.UUID             `json:"parent_id"`
	MenuOrder        int                    `json:"menu_order"`
	ShowInMenu       bool                   `json:"show_in_menu"`
	RequireAuth      bool                   `json:"require_auth"`
	AllowComments    bool                   `json:"allow_comments"`
	CustomFields     map[string]interface{} `json:"custom_fields"`
	PublishedAt      *time.Time             `json:"published_at"`
	ScheduledAt      *time.Time             `json:"scheduled_at"`
}

type PageResponse struct {
	ID               uuid.UUID              `json:"id"`
	Title            string                 `json:"title"`
	Slug             string                 `json:"slug"`
	Content          string                 `json:"content"`
	Excerpt          string                 `json:"excerpt"`
	Status           string                 `json:"status"`
	Template         string                 `json:"template"`
	MetaTitle        string                 `json:"meta_title"`
	MetaDescription  string                 `json:"meta_description"`
	MetaKeywords     []string               `json:"meta_keywords"`
	FeaturedImage    string                 `json:"featured_image"`
	ParentID         *uuid.UUID             `json:"parent_id"`
	MenuOrder        int                    `json:"menu_order"`
	ShowInMenu       bool                   `json:"show_in_menu"`
	RequireAuth      bool                   `json:"require_auth"`
	AllowComments    bool                   `json:"allow_comments"`
	CustomFields     map[string]interface{} `json:"custom_fields"`
	ViewCount        int64                  `json:"view_count"`
	CommentCount     int64                  `json:"comment_count"`
	PublishedAt      *time.Time             `json:"published_at"`
	ScheduledAt      *time.Time             `json:"scheduled_at"`
	AuthorID         uuid.UUID              `json:"author_id"`
	AuthorName       string                 `json:"author_name"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

type GetPagesResponse struct {
	Pages      []PageResponse     `json:"pages"`
	Total      int64              `json:"total"`
	Pagination *PaginationInfo    `json:"pagination"`
}

// Block types
type GetBlocksRequest struct {
	Page      int    `json:"page" validate:"min=1"`
	Limit     int    `json:"limit" validate:"min=1,max=100"`
	Search    string `json:"search"`
	Type      string `json:"type"`
	Position  string `json:"position"`
	IsActive  *bool  `json:"is_active"`
	SortBy    string `json:"sort_by" validate:"omitempty,oneof=name created_at updated_at"`
	SortOrder string `json:"sort_order" validate:"omitempty,oneof=asc desc"`
}

type CreateBlockRequest struct {
	Name        string                 `json:"name" validate:"required,min=1,max=100"`
	Type        string                 `json:"type" validate:"required,oneof=text html image video banner widget"`
	Position    string                 `json:"position" validate:"required"`
	Content     string                 `json:"content"`
	Settings    map[string]interface{} `json:"settings"`
	IsActive    bool                   `json:"is_active"`
	SortOrder   int                    `json:"sort_order"`
	Conditions  map[string]interface{} `json:"conditions"`
	StartDate   *time.Time             `json:"start_date"`
	EndDate     *time.Time             `json:"end_date"`
}

type UpdateBlockRequest struct {
	Name        string                 `json:"name" validate:"min=1,max=100"`
	Type        string                 `json:"type" validate:"omitempty,oneof=text html image video banner widget"`
	Position    string                 `json:"position"`
	Content     string                 `json:"content"`
	Settings    map[string]interface{} `json:"settings"`
	IsActive    bool                   `json:"is_active"`
	SortOrder   int                    `json:"sort_order"`
	Conditions  map[string]interface{} `json:"conditions"`
	StartDate   *time.Time             `json:"start_date"`
	EndDate     *time.Time             `json:"end_date"`
}

type BlockResponse struct {
	ID         uuid.UUID              `json:"id"`
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Position   string                 `json:"position"`
	Content    string                 `json:"content"`
	Settings   map[string]interface{} `json:"settings"`
	IsActive   bool                   `json:"is_active"`
	SortOrder  int                    `json:"sort_order"`
	Conditions map[string]interface{} `json:"conditions"`
	StartDate  *time.Time             `json:"start_date"`
	EndDate    *time.Time             `json:"end_date"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

type GetBlocksResponse struct {
	Blocks     []BlockResponse    `json:"blocks"`
	Total      int64              `json:"total"`
	Pagination *PaginationInfo    `json:"pagination"`
}

// Menu types
type CreateMenuRequest struct {
	Name        string           `json:"name" validate:"required,min=1,max=100"`
	Location    string           `json:"location" validate:"required"`
	Description string           `json:"description" validate:"max=500"`
	IsActive    bool             `json:"is_active"`
	Items       []MenuItemRequest `json:"items"`
}

type UpdateMenuRequest struct {
	Name        string           `json:"name" validate:"min=1,max=100"`
	Location    string           `json:"location"`
	Description string           `json:"description" validate:"max=500"`
	IsActive    bool             `json:"is_active"`
	Items       []MenuItemRequest `json:"items"`
}

type MenuItemRequest struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title" validate:"required,min=1,max=100"`
	URL         string    `json:"url" validate:"required"`
	Target      string    `json:"target" validate:"omitempty,oneof=_self _blank"`
	Icon        string    `json:"icon"`
	CSSClass    string    `json:"css_class"`
	ParentID    *uuid.UUID `json:"parent_id"`
	SortOrder   int       `json:"sort_order"`
	IsActive    bool      `json:"is_active"`
	Conditions  map[string]interface{} `json:"conditions"`
}

type MenuResponse struct {
	ID          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`
	Location    string           `json:"location"`
	Description string           `json:"description"`
	IsActive    bool             `json:"is_active"`
	Items       []MenuItemResponse `json:"items"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type MenuItemResponse struct {
	ID          uuid.UUID              `json:"id"`
	Title       string                 `json:"title"`
	URL         string                 `json:"url"`
	Target      string                 `json:"target"`
	Icon        string                 `json:"icon"`
	CSSClass    string                 `json:"css_class"`
	ParentID    *uuid.UUID             `json:"parent_id"`
	SortOrder   int                    `json:"sort_order"`
	IsActive    bool                   `json:"is_active"`
	Conditions  map[string]interface{} `json:"conditions"`
	Children    []MenuItemResponse     `json:"children"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type GetMenusResponse struct {
	Menus []MenuResponse `json:"menus"`
	Total int64          `json:"total"`
}

// Banner types
type GetBannersRequest struct {
	Page      int    `json:"page" validate:"min=1"`
	Limit     int    `json:"limit" validate:"min=1,max=100"`
	Search    string `json:"search"`
	Position  string `json:"position"`
	IsActive  *bool  `json:"is_active"`
	SortBy    string `json:"sort_by" validate:"omitempty,oneof=title created_at updated_at"`
	SortOrder string `json:"sort_order" validate:"omitempty,oneof=asc desc"`
}

type CreateBannerRequest struct {
	Title       string                 `json:"title" validate:"required,min=1,max=200"`
	Description string                 `json:"description" validate:"max=500"`
	Image       string                 `json:"image" validate:"required"`
	Link        string                 `json:"link"`
	Target      string                 `json:"target" validate:"omitempty,oneof=_self _blank"`
	Position    string                 `json:"position" validate:"required"`
	SortOrder   int                    `json:"sort_order"`
	IsActive    bool                   `json:"is_active"`
	StartDate   *time.Time             `json:"start_date"`
	EndDate     *time.Time             `json:"end_date"`
	Settings    map[string]interface{} `json:"settings"`
	Conditions  map[string]interface{} `json:"conditions"`
}

type UpdateBannerRequest struct {
	Title       string                 `json:"title" validate:"min=1,max=200"`
	Description string                 `json:"description" validate:"max=500"`
	Image       string                 `json:"image"`
	Link        string                 `json:"link"`
	Target      string                 `json:"target" validate:"omitempty,oneof=_self _blank"`
	Position    string                 `json:"position"`
	SortOrder   int                    `json:"sort_order"`
	IsActive    bool                   `json:"is_active"`
	StartDate   *time.Time             `json:"start_date"`
	EndDate     *time.Time             `json:"end_date"`
	Settings    map[string]interface{} `json:"settings"`
	Conditions  map[string]interface{} `json:"conditions"`
}

type BannerResponse struct {
	ID          uuid.UUID              `json:"id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Image       string                 `json:"image"`
	Link        string                 `json:"link"`
	Target      string                 `json:"target"`
	Position    string                 `json:"position"`
	SortOrder   int                    `json:"sort_order"`
	IsActive    bool                   `json:"is_active"`
	StartDate   *time.Time             `json:"start_date"`
	EndDate     *time.Time             `json:"end_date"`
	Settings    map[string]interface{} `json:"settings"`
	Conditions  map[string]interface{} `json:"conditions"`
	ViewCount   int64                  `json:"view_count"`
	ClickCount  int64                  `json:"click_count"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type GetBannersResponse struct {
	Banners    []BannerResponse   `json:"banners"`
	Total      int64              `json:"total"`
	Pagination *PaginationInfo    `json:"pagination"`
}

// Slider types
type GetSlidersRequest struct {
	Page      int    `json:"page" validate:"min=1"`
	Limit     int    `json:"limit" validate:"min=1,max=100"`
	Search    string `json:"search"`
	IsActive  *bool  `json:"is_active"`
	SortBy    string `json:"sort_by" validate:"omitempty,oneof=name created_at updated_at"`
	SortOrder string `json:"sort_order" validate:"omitempty,oneof=asc desc"`
}

type CreateSliderRequest struct {
	Name        string                 `json:"name" validate:"required,min=1,max=100"`
	Description string                 `json:"description" validate:"max=500"`
	IsActive    bool                   `json:"is_active"`
	Settings    map[string]interface{} `json:"settings"`
	Slides      []SlideRequest         `json:"slides"`
}

type UpdateSliderRequest struct {
	Name        string                 `json:"name" validate:"min=1,max=100"`
	Description string                 `json:"description" validate:"max=500"`
	IsActive    bool                   `json:"is_active"`
	Settings    map[string]interface{} `json:"settings"`
	Slides      []SlideRequest         `json:"slides"`
}

type SlideRequest struct {
	ID          uuid.UUID              `json:"id"`
	Title       string                 `json:"title" validate:"required,min=1,max=200"`
	Description string                 `json:"description" validate:"max=500"`
	Image       string                 `json:"image" validate:"required"`
	Link        string                 `json:"link"`
	Target      string                 `json:"target" validate:"omitempty,oneof=_self _blank"`
	SortOrder   int                    `json:"sort_order"`
	IsActive    bool                   `json:"is_active"`
	Settings    map[string]interface{} `json:"settings"`
}

type SliderResponse struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	IsActive    bool                   `json:"is_active"`
	Settings    map[string]interface{} `json:"settings"`
	Slides      []SlideResponse        `json:"slides"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type SlideResponse struct {
	ID          uuid.UUID              `json:"id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Image       string                 `json:"image"`
	Link        string                 `json:"link"`
	Target      string                 `json:"target"`
	SortOrder   int                    `json:"sort_order"`
	IsActive    bool                   `json:"is_active"`
	Settings    map[string]interface{} `json:"settings"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type GetSlidersResponse struct {
	Sliders    []SliderResponse   `json:"sliders"`
	Total      int64              `json:"total"`
	Pagination *PaginationInfo    `json:"pagination"`
}

// Widget types
type GetWidgetsRequest struct {
	Page      int    `json:"page" validate:"min=1"`
	Limit     int    `json:"limit" validate:"min=1,max=100"`
	Search    string `json:"search"`
	Type      string `json:"type"`
	Position  string `json:"position"`
	IsActive  *bool  `json:"is_active"`
	SortBy    string `json:"sort_by" validate:"omitempty,oneof=name created_at updated_at"`
	SortOrder string `json:"sort_order" validate:"omitempty,oneof=asc desc"`
}

type CreateWidgetRequest struct {
	Name        string                 `json:"name" validate:"required,min=1,max=100"`
	Type        string                 `json:"type" validate:"required"`
	Position    string                 `json:"position" validate:"required"`
	Title       string                 `json:"title" validate:"max=200"`
	Content     string                 `json:"content"`
	Settings    map[string]interface{} `json:"settings"`
	IsActive    bool                   `json:"is_active"`
	SortOrder   int                    `json:"sort_order"`
	Conditions  map[string]interface{} `json:"conditions"`
}

type UpdateWidgetRequest struct {
	Name        string                 `json:"name" validate:"min=1,max=100"`
	Type        string                 `json:"type"`
	Position    string                 `json:"position"`
	Title       string                 `json:"title" validate:"max=200"`
	Content     string                 `json:"content"`
	Settings    map[string]interface{} `json:"settings"`
	IsActive    bool                   `json:"is_active"`
	SortOrder   int                    `json:"sort_order"`
	Conditions  map[string]interface{} `json:"conditions"`
}

type WidgetResponse struct {
	ID         uuid.UUID              `json:"id"`
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Position   string                 `json:"position"`
	Title      string                 `json:"title"`
	Content    string                 `json:"content"`
	Settings   map[string]interface{} `json:"settings"`
	IsActive   bool                   `json:"is_active"`
	SortOrder  int                    `json:"sort_order"`
	Conditions map[string]interface{} `json:"conditions"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

type GetWidgetsResponse struct {
	Widgets    []WidgetResponse   `json:"widgets"`
	Total      int64              `json:"total"`
	Pagination *PaginationInfo    `json:"pagination"`
}

// Template types
type CreateTemplateRequest struct {
	Name        string                 `json:"name" validate:"required,min=1,max=100"`
	Type        string                 `json:"type" validate:"required,oneof=page product category email"`
	Content     string                 `json:"content" validate:"required"`
	Description string                 `json:"description" validate:"max=500"`
	IsActive    bool                   `json:"is_active"`
	Settings    map[string]interface{} `json:"settings"`
}

type UpdateTemplateRequest struct {
	Name        string                 `json:"name" validate:"min=1,max=100"`
	Type        string                 `json:"type" validate:"omitempty,oneof=page product category email"`
	Content     string                 `json:"content"`
	Description string                 `json:"description" validate:"max=500"`
	IsActive    bool                   `json:"is_active"`
	Settings    map[string]interface{} `json:"settings"`
}

type TemplateResponse struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Content     string                 `json:"content"`
	Description string                 `json:"description"`
	IsActive    bool                   `json:"is_active"`
	Settings    map[string]interface{} `json:"settings"`
	UsageCount  int64                  `json:"usage_count"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type GetTemplatesResponse struct {
	Templates []TemplateResponse `json:"templates"`
	Total     int64              `json:"total"`
}
