package feedback

import "time"

// ImprovementAction represents an improvement action
type ImprovementAction struct {
	ID          string
	Type        int
	Description string
	Priority    int
	Effort      string
	Impact      float64
	Status      int
	CreatedAt   time.Time
	CompletedAt *time.Time
	Details     map[string]interface{}
}

// FeedbackConfig represents feedback configuration
type FeedbackConfig struct {
	EnableFeedbackCollection bool
	FeedbackTimeout         time.Duration
	MaxFeedbackItems        int
	FeedbackRetentionDays   int
	EnableAutoProcessing    bool
	ProcessingThreshold     float64
	ImprovementThreshold    float64
	MinFeedbackScore        float64
	MaxFeedbackScore        float64
	QualityWeight           float64
	UsabilityWeight         float64
	AestheticsWeight        float64
	EnableVisualImprovements bool
	EnableLayoutImprovements bool
	EnablePerformanceImprovements bool
}

// FeedbackItem represents a feedback item
type FeedbackItem struct {
	ID              string
	UserID          string
	SessionID       string
	Timestamp       time.Time
	FeedbackType    int
	Category        int
	Priority        int
	Score           float64
	Title           string
	Description     string
	Context         map[string]interface{}
	Attachments     []FeedbackAttachment
	Tags            []string
	Status          int
	ProcessedAt     *time.Time
	Improvements    []ImprovementAction
	FollowUp        *FollowUpRequest
}

// FeedbackAttachment represents a feedback attachment
type FeedbackAttachment struct {
	ID          string
	Type        string
	Filename    string
	Size        int64
	ContentType string
	URL         string
	UploadedAt  time.Time
}

// FollowUpRequest represents a follow-up request
type FollowUpRequest struct {
	ID          string
	Message     string
	RequestedAt time.Time
	DueDate     time.Time
	Status      string
	Response    string
	RespondedAt *time.Time
}

// FeedbackMetrics represents feedback metrics
type FeedbackMetrics struct {
	TotalFeedback       int
	PendingFeedback     int
	ProcessedFeedback   int
	ImplementedFeedback int
	AverageScore        float64
	CategoryBreakdown   map[string]int
	PriorityBreakdown   map[string]int
	StatusBreakdown     map[string]int
	Trends              []FeedbackTrend
	TopIssues           []TopIssue
	ImprovementStats    ImprovementStats
}

// FeedbackTrend represents a feedback trend
type FeedbackTrend struct {
	Period    string
	Count     int
	AvgScore  float64
	Category  string
}

// TopIssue represents a top issue
type TopIssue struct {
	Issue       string
	Count       int
	AvgScore    float64
	Priority    string
	Category    string
}

// ImprovementStats represents improvement statistics
type ImprovementStats struct {
	TotalImprovements    int
	CompletedImprovements int
	InProgressImprovements int
	PlannedImprovements  int
	AverageImpact        float64
	AverageEffort        string
}

// ImprovementConfig represents improvement configuration
type ImprovementConfig struct {
	EnableAutoImprovements    bool
	ImprovementThreshold      float64
	MaxConcurrentImprovements int
	ImprovementTimeout        time.Duration
	EnableVisualImprovements  bool
	VisualImprovementWeight   float64
	EnableLayoutImprovements  bool
	LayoutImprovementWeight   float64
	EnablePerformanceImprovements bool
	PerformanceImprovementWeight  float64
}

// ImprovementResult represents the result of an improvement
type ImprovementResult struct {
	ActionID      string
	Success       bool
	Message       string
	Changes       map[string]interface{}
	Performance   *PerformanceMetrics
	Timestamp     time.Time
}

// PerformanceMetrics represents performance metrics
type PerformanceMetrics struct {
	BeforeScore float64
	AfterScore  float64
	Improvement float64
	Duration    time.Duration
}
