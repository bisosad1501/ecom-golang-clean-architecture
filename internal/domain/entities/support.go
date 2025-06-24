package entities

import (
	"time"

	"github.com/google/uuid"
)

// TicketStatus represents the status of a support ticket
type TicketStatus string

const (
	TicketStatusOpen       TicketStatus = "open"
	TicketStatusInProgress TicketStatus = "in_progress"
	TicketStatusPending    TicketStatus = "pending"
	TicketStatusResolved   TicketStatus = "resolved"
	TicketStatusClosed     TicketStatus = "closed"
	TicketStatusCancelled  TicketStatus = "cancelled"
)

// TicketPriority represents the priority of a support ticket
type TicketPriority string

const (
	TicketPriorityLow      TicketPriority = "low"
	TicketPriorityNormal   TicketPriority = "normal"
	TicketPriorityHigh     TicketPriority = "high"
	TicketPriorityCritical TicketPriority = "critical"
	TicketPriorityUrgent   TicketPriority = "urgent"
)

// TicketCategory represents the category of a support ticket
type TicketCategory string

const (
	TicketCategoryGeneral    TicketCategory = "general"
	TicketCategoryOrder      TicketCategory = "order"
	TicketCategoryPayment    TicketCategory = "payment"
	TicketCategoryShipping   TicketCategory = "shipping"
	TicketCategoryProduct    TicketCategory = "product"
	TicketCategoryAccount    TicketCategory = "account"
	TicketCategoryTechnical  TicketCategory = "technical"
	TicketCategoryRefund     TicketCategory = "refund"
	TicketCategoryComplaint  TicketCategory = "complaint"
	TicketCategoryFeedback   TicketCategory = "feedback"
)

// SupportTicket represents a customer support ticket
type SupportTicket struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TicketNumber    string         `json:"ticket_number" gorm:"uniqueIndex;not null"`
	UserID          uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	User            User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	AssignedToID    *uuid.UUID     `json:"assigned_to_id" gorm:"type:uuid;index"`
	AssignedTo      *User          `json:"assigned_to,omitempty" gorm:"foreignKey:AssignedToID"`
	
	// Ticket details
	Subject         string         `json:"subject" gorm:"not null" validate:"required"`
	Description     string         `json:"description" gorm:"type:text;not null" validate:"required"`
	Category        TicketCategory `json:"category" gorm:"not null"`
	Priority        TicketPriority `json:"priority" gorm:"default:'normal'"`
	Status          TicketStatus   `json:"status" gorm:"default:'open'"`
	
	// Reference information
	OrderID         *uuid.UUID     `json:"order_id" gorm:"type:uuid;index"`
	Order           *Order         `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	ProductID       *uuid.UUID     `json:"product_id" gorm:"type:uuid;index"`
	Product         *Product       `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	
	// Contact information
	ContactEmail    string         `json:"contact_email"`
	ContactPhone    string         `json:"contact_phone"`
	
	// Resolution information
	Resolution      string         `json:"resolution" gorm:"type:text"`
	ResolutionTime  *time.Time     `json:"resolution_time"`
	ResolvedBy      *uuid.UUID     `json:"resolved_by" gorm:"type:uuid"`
	
	// SLA tracking
	FirstResponseAt *time.Time     `json:"first_response_at"`
	LastResponseAt  *time.Time     `json:"last_response_at"`
	DueDate         *time.Time     `json:"due_date"`
	
	// Customer satisfaction
	SatisfactionRating *int        `json:"satisfaction_rating"`     // 1-5 scale
	SatisfactionFeedback string    `json:"satisfaction_feedback" gorm:"type:text"`
	
	// Tags and labels
	Tags            string         `json:"tags"`                    // JSON array of tags
	Labels          string         `json:"labels"`                  // JSON array of labels
	
	// Metadata
	CreatedAt       time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	
	// Relationships (removed to avoid circular dependencies during migration)
}

// TableName returns the table name for SupportTicket entity
func (SupportTicket) TableName() string {
	return "support_tickets"
}

// IsOpen checks if ticket is open
func (st *SupportTicket) IsOpen() bool {
	return st.Status == TicketStatusOpen || st.Status == TicketStatusInProgress
}

// IsResolved checks if ticket is resolved
func (st *SupportTicket) IsResolved() bool {
	return st.Status == TicketStatusResolved || st.Status == TicketStatusClosed
}

// GetResponseTime calculates first response time in hours
func (st *SupportTicket) GetResponseTime() float64 {
	if st.FirstResponseAt == nil {
		return 0
	}
	return st.FirstResponseAt.Sub(st.CreatedAt).Hours()
}

// GetResolutionTime calculates resolution time in hours
func (st *SupportTicket) GetResolutionTime() float64 {
	if st.ResolutionTime == nil {
		return 0
	}
	return st.ResolutionTime.Sub(st.CreatedAt).Hours()
}

// IsOverdue checks if ticket is overdue
func (st *SupportTicket) IsOverdue() bool {
	if st.DueDate == nil {
		return false
	}
	return time.Now().After(*st.DueDate) && !st.IsResolved()
}

// TicketMessage represents messages in a support ticket
type TicketMessage struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TicketID    uuid.UUID `json:"ticket_id" gorm:"type:uuid;not null;index"`
	Ticket      SupportTicket `json:"ticket,omitempty" gorm:"foreignKey:TicketID"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	User        User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	
	// Message details
	Message     string    `json:"message" gorm:"type:text;not null" validate:"required"`
	IsInternal  bool      `json:"is_internal" gorm:"default:false"`        // Internal notes vs customer messages
	IsFromStaff bool      `json:"is_from_staff" gorm:"default:false"`      // Message from staff vs customer
	
	// Metadata
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	
	// Relationships (removed to avoid circular dependencies during migration)
}

// TableName returns the table name for TicketMessage entity
func (TicketMessage) TableName() string {
	return "ticket_messages"
}

// TicketAttachment represents file attachments in tickets
type TicketAttachment struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TicketID    uuid.UUID `json:"ticket_id" gorm:"type:uuid;not null;index"`
	Ticket      SupportTicket `json:"ticket,omitempty" gorm:"foreignKey:TicketID"`
	MessageID   *uuid.UUID `json:"message_id" gorm:"type:uuid;index"`
	Message     *TicketMessage `json:"message,omitempty" gorm:"foreignKey:MessageID"`
	
	// File details
	FileName    string    `json:"file_name" gorm:"not null"`
	FileSize    int64     `json:"file_size" gorm:"not null"`
	FileType    string    `json:"file_type" gorm:"not null"`
	FilePath    string    `json:"file_path" gorm:"not null"`
	FileURL     string    `json:"file_url"`
	
	// Upload details
	UploadedBy  uuid.UUID `json:"uploaded_by" gorm:"type:uuid;not null"`
	UploadedAt  time.Time `json:"uploaded_at" gorm:"autoCreateTime"`
}

// TableName returns the table name for TicketAttachment entity
func (TicketAttachment) TableName() string {
	return "ticket_attachments"
}

// FAQ represents frequently asked questions
type FAQ struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Question    string         `json:"question" gorm:"not null" validate:"required"`
	Answer      string         `json:"answer" gorm:"type:text;not null" validate:"required"`
	Category    TicketCategory `json:"category" gorm:"not null"`
	
	// SEO and search
	Slug        string         `json:"slug" gorm:"uniqueIndex"`
	Keywords    string         `json:"keywords"`                    // Comma-separated keywords
	
	// Display settings
	IsPublished bool           `json:"is_published" gorm:"default:true"`
	IsFeatured  bool           `json:"is_featured" gorm:"default:false"`
	SortOrder   int            `json:"sort_order" gorm:"default:0"`
	
	// Usage statistics
	ViewCount   int            `json:"view_count" gorm:"default:0"`
	HelpfulCount int           `json:"helpful_count" gorm:"default:0"`
	NotHelpfulCount int        `json:"not_helpful_count" gorm:"default:0"`
	
	// Metadata
	CreatedBy   uuid.UUID      `json:"created_by" gorm:"type:uuid"`
	UpdatedBy   uuid.UUID      `json:"updated_by" gorm:"type:uuid"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for FAQ entity
func (FAQ) TableName() string {
	return "faqs"
}

// GetHelpfulPercentage calculates helpful percentage
func (faq *FAQ) GetHelpfulPercentage() float64 {
	total := faq.HelpfulCount + faq.NotHelpfulCount
	if total == 0 {
		return 0
	}
	return float64(faq.HelpfulCount) / float64(total) * 100
}

// KnowledgeBase represents knowledge base articles
type KnowledgeBase struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Title       string         `json:"title" gorm:"not null" validate:"required"`
	Content     string         `json:"content" gorm:"type:text;not null" validate:"required"`
	Summary     string         `json:"summary" gorm:"type:text"`
	Category    TicketCategory `json:"category" gorm:"not null"`
	
	// SEO and search
	Slug        string         `json:"slug" gorm:"uniqueIndex"`
	Keywords    string         `json:"keywords"`                    // Comma-separated keywords
	MetaTitle   string         `json:"meta_title"`
	MetaDescription string     `json:"meta_description"`
	
	// Display settings
	IsPublished bool           `json:"is_published" gorm:"default:true"`
	IsFeatured  bool           `json:"is_featured" gorm:"default:false"`
	SortOrder   int            `json:"sort_order" gorm:"default:0"`
	
	// Content details
	ReadingTime int            `json:"reading_time" gorm:"default:0"`   // Estimated reading time in minutes
	Difficulty  string         `json:"difficulty" gorm:"default:'beginner'"` // beginner, intermediate, advanced
	
	// Usage statistics
	ViewCount   int            `json:"view_count" gorm:"default:0"`
	HelpfulCount int           `json:"helpful_count" gorm:"default:0"`
	NotHelpfulCount int        `json:"not_helpful_count" gorm:"default:0"`
	ShareCount  int            `json:"share_count" gorm:"default:0"`
	
	// Related content
	RelatedArticles string     `json:"related_articles" gorm:"type:text"` // JSON array of article IDs
	Tags           string      `json:"tags" gorm:"type:text"`             // JSON array of tags
	
	// Metadata
	CreatedBy   uuid.UUID      `json:"created_by" gorm:"type:uuid"`
	UpdatedBy   uuid.UUID      `json:"updated_by" gorm:"type:uuid"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for KnowledgeBase entity
func (KnowledgeBase) TableName() string {
	return "knowledge_base"
}

// GetHelpfulPercentage calculates helpful percentage
func (kb *KnowledgeBase) GetHelpfulPercentage() float64 {
	total := kb.HelpfulCount + kb.NotHelpfulCount
	if total == 0 {
		return 0
	}
	return float64(kb.HelpfulCount) / float64(total) * 100
}

// LiveChatSession represents live chat sessions
type LiveChatSession struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SessionID       string    `json:"session_id" gorm:"uniqueIndex;not null"`
	UserID          *uuid.UUID `json:"user_id" gorm:"type:uuid;index"`        // null for anonymous users
	User            *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	AgentID         *uuid.UUID `json:"agent_id" gorm:"type:uuid;index"`
	Agent           *User     `json:"agent,omitempty" gorm:"foreignKey:AgentID"`
	
	// Session details
	Status          string    `json:"status" gorm:"default:'waiting'"`        // waiting, active, ended
	Subject         string    `json:"subject"`
	Department      string    `json:"department" gorm:"default:'general'"`
	
	// Contact information (for anonymous users)
	GuestName       string    `json:"guest_name"`
	GuestEmail      string    `json:"guest_email"`
	
	// Session metrics
	StartedAt       time.Time `json:"started_at" gorm:"autoCreateTime"`
	EndedAt         *time.Time `json:"ended_at"`
	FirstResponseAt *time.Time `json:"first_response_at"`
	Duration        int       `json:"duration" gorm:"default:0"`              // Duration in seconds
	MessageCount    int       `json:"message_count" gorm:"default:0"`
	
	// Satisfaction
	Rating          *int      `json:"rating"`                                 // 1-5 scale
	Feedback        string    `json:"feedback" gorm:"type:text"`
	
	// Metadata
	UserAgent       string    `json:"user_agent"`
	IPAddress       string    `json:"ip_address"`
	Referrer        string    `json:"referrer"`
	
	// Relationships (removed to avoid circular dependencies during migration)
}

// TableName returns the table name for LiveChatSession entity
func (LiveChatSession) TableName() string {
	return "live_chat_sessions"
}

// IsActive checks if chat session is active
func (lcs *LiveChatSession) IsActive() bool {
	return lcs.Status == "active"
}

// GetResponseTime calculates first response time in seconds
func (lcs *LiveChatSession) GetResponseTime() int {
	if lcs.FirstResponseAt == nil {
		return 0
	}
	return int(lcs.FirstResponseAt.Sub(lcs.StartedAt).Seconds())
}

// ChatMessage represents messages in live chat
type ChatMessage struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SessionID   uuid.UUID `json:"session_id" gorm:"type:uuid;not null;index"`
	SenderID    *uuid.UUID `json:"sender_id" gorm:"type:uuid;index"`          // null for system messages
	Sender      *User     `json:"sender,omitempty" gorm:"foreignKey:SenderID"`
	
	// Message details
	Message     string    `json:"message" gorm:"type:text;not null"`
	MessageType string    `json:"message_type" gorm:"default:'text'"`         // text, image, file, system
	IsFromAgent bool      `json:"is_from_agent" gorm:"default:false"`
	
	// File attachment (for file messages)
	FileName    string    `json:"file_name"`
	FileURL     string    `json:"file_url"`
	FileSize    int64     `json:"file_size"`
	
	// Message status
	IsRead      bool      `json:"is_read" gorm:"default:false"`
	ReadAt      *time.Time `json:"read_at"`
	
	// Metadata
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName returns the table name for ChatMessage entity
func (ChatMessage) TableName() string {
	return "chat_messages"
}
