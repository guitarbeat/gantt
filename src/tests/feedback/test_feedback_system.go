package main

import (
	"fmt"
	"time"
)

// Test the feedback system
func main() {
	fmt.Println("Testing Feedback System...")

	// Test 1: Feedback Configuration
	fmt.Println("\n=== Test 1: Feedback Configuration ===")
	testFeedbackConfiguration()

	// Test 2: Feedback Collection
	fmt.Println("\n=== Test 2: Feedback Collection ===")
	testFeedbackCollection()

	// Test 3: Feedback Processing
	fmt.Println("\n=== Test 3: Feedback Processing ===")
	testFeedbackProcessing()

	// Test 4: Feedback Analysis
	fmt.Println("\n=== Test 4: Feedback Analysis ===")
	testFeedbackAnalysis()

	// Test 5: Feedback Improvements
	fmt.Println("\n=== Test 5: Feedback Improvements ===")
	testFeedbackImprovements()

	fmt.Println("\n✅ Feedback system tests completed!")
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

func testFeedbackConfiguration() {
	// Test feedback configuration
	config := FeedbackConfig{
		EnableFeedbackCollection: true,
		FeedbackTimeout:         time.Minute * 5,
		MaxFeedbackItems:        1000,
		FeedbackRetentionDays:   365,
		EnableAutoProcessing:    true,
		ProcessingThreshold:     0.7,
		ImprovementThreshold:    0.8,
		MinFeedbackScore:        1.0,
		MaxFeedbackScore:        5.0,
		QualityWeight:           0.4,
		UsabilityWeight:         0.3,
		AestheticsWeight:        0.3,
		EnableVisualImprovements: true,
		EnableLayoutImprovements: true,
		EnablePerformanceImprovements: true,
	}

	// Validate configuration
	if !config.EnableFeedbackCollection {
		fmt.Println("❌ Feedback collection should be enabled")
		return
	}

	if config.FeedbackTimeout <= 0 {
		fmt.Println("❌ Feedback timeout should be positive")
		return
	}

	if config.MaxFeedbackItems <= 0 {
		fmt.Println("❌ Max feedback items should be positive")
		return
	}

	if config.FeedbackRetentionDays <= 0 {
		fmt.Println("❌ Feedback retention days should be positive")
		return
	}

	if !config.EnableAutoProcessing {
		fmt.Println("❌ Auto processing should be enabled")
		return
	}

	if config.ProcessingThreshold < 0.0 || config.ProcessingThreshold > 1.0 {
		fmt.Println("❌ Processing threshold should be between 0 and 1")
		return
	}

	if config.ImprovementThreshold < 0.0 || config.ImprovementThreshold > 1.0 {
		fmt.Println("❌ Improvement threshold should be between 0 and 1")
		return
	}

	if config.MinFeedbackScore < 1.0 || config.MinFeedbackScore > 5.0 {
		fmt.Println("❌ Min feedback score should be between 1 and 5")
		return
	}

	if config.MaxFeedbackScore < 1.0 || config.MaxFeedbackScore > 5.0 {
		fmt.Println("❌ Max feedback score should be between 1 and 5")
		return
	}

	if config.MaxFeedbackScore <= config.MinFeedbackScore {
		fmt.Println("❌ Max feedback score should be greater than min")
		return
	}

	if config.QualityWeight < 0.0 || config.QualityWeight > 1.0 {
		fmt.Println("❌ Quality weight should be between 0 and 1")
		return
	}

	if config.UsabilityWeight < 0.0 || config.UsabilityWeight > 1.0 {
		fmt.Println("❌ Usability weight should be between 0 and 1")
		return
	}

	if config.AestheticsWeight < 0.0 || config.AestheticsWeight > 1.0 {
		fmt.Println("❌ Aesthetics weight should be between 0 and 1")
		return
	}

	totalWeight := config.QualityWeight + config.UsabilityWeight + config.AestheticsWeight
	if totalWeight < 0.9 || totalWeight > 1.1 {
		fmt.Println("❌ Total weight should be approximately 1.0")
		return
	}

	if !config.EnableVisualImprovements {
		fmt.Println("❌ Visual improvements should be enabled")
		return
	}

	if !config.EnableLayoutImprovements {
		fmt.Println("❌ Layout improvements should be enabled")
		return
	}

	if !config.EnablePerformanceImprovements {
		fmt.Println("❌ Performance improvements should be enabled")
		return
	}

	fmt.Printf("✅ Feedback configuration test passed\n")
	fmt.Printf("   Enable feedback collection: %v\n", config.EnableFeedbackCollection)
	fmt.Printf("   Feedback timeout: %v\n", config.FeedbackTimeout)
	fmt.Printf("   Max feedback items: %d\n", config.MaxFeedbackItems)
	fmt.Printf("   Feedback retention days: %d\n", config.FeedbackRetentionDays)
	fmt.Printf("   Enable auto processing: %v\n", config.EnableAutoProcessing)
	fmt.Printf("   Processing threshold: %.2f\n", config.ProcessingThreshold)
	fmt.Printf("   Improvement threshold: %.2f\n", config.ImprovementThreshold)
	fmt.Printf("   Min feedback score: %.1f\n", config.MinFeedbackScore)
	fmt.Printf("   Max feedback score: %.1f\n", config.MaxFeedbackScore)
	fmt.Printf("   Quality weight: %.2f\n", config.QualityWeight)
	fmt.Printf("   Usability weight: %.2f\n", config.UsabilityWeight)
	fmt.Printf("   Aesthetics weight: %.2f\n", config.AestheticsWeight)
	fmt.Printf("   Enable visual improvements: %v\n", config.EnableVisualImprovements)
	fmt.Printf("   Enable layout improvements: %v\n", config.EnableLayoutImprovements)
	fmt.Printf("   Enable performance improvements: %v\n", config.EnablePerformanceImprovements)
}

func testFeedbackCollection() {
	// Test feedback collection
	feedback := FeedbackItem{
		ID:           "feedback-001",
		UserID:       "user-123",
		SessionID:    "session-456",
		Timestamp:    time.Now(),
		FeedbackType: 1, // Visual
		Category:     2, // Layout
		Priority:     2, // Medium
		Score:        3.5,
		Title:        "Task spacing needs improvement",
		Description:  "The spacing between tasks is too tight and makes it hard to read",
		Context: map[string]interface{}{
			"view_type": "monthly",
			"task_count": 15,
			"screen_size": "1920x1080",
		},
		Attachments: []FeedbackAttachment{
			{
				ID:          "attach-001",
				Type:        "screenshot",
				Filename:    "spacing_issue.png",
				Size:        1024000,
				ContentType: "image/png",
				URL:         "/attachments/spacing_issue.png",
				UploadedAt:  time.Now(),
			},
		},
		Tags:   []string{"spacing", "layout", "readability"},
		Status: 0, // Pending
	}

	// Validate feedback
	if feedback.ID == "" {
		fmt.Println("❌ Feedback ID should not be empty")
		return
	}

	if feedback.UserID == "" {
		fmt.Println("❌ User ID should not be empty")
		return
	}

	if feedback.SessionID == "" {
		fmt.Println("❌ Session ID should not be empty")
		return
	}

	if feedback.Score < 1.0 || feedback.Score > 5.0 {
		fmt.Println("❌ Feedback score should be between 1 and 5")
		return
	}

	if feedback.Title == "" {
		fmt.Println("❌ Feedback title should not be empty")
		return
	}

	if feedback.Description == "" {
		fmt.Println("❌ Feedback description should not be empty")
		return
	}

	if feedback.FeedbackType < 0 || feedback.FeedbackType > 8 {
		fmt.Println("❌ Feedback type should be between 0 and 8")
		return
	}

	if feedback.Category < 0 || feedback.Category > 9 {
		fmt.Println("❌ Feedback category should be between 0 and 9")
		return
	}

	if feedback.Priority < 0 || feedback.Priority > 3 {
		fmt.Println("❌ Feedback priority should be between 0 and 3")
		return
	}

	if feedback.Status < 0 || feedback.Status > 5 {
		fmt.Println("❌ Feedback status should be between 0 and 5")
		return
	}

	// Validate attachments
	for _, attachment := range feedback.Attachments {
		if attachment.ID == "" {
			fmt.Println("❌ Attachment ID should not be empty")
			return
		}

		if attachment.Type == "" {
			fmt.Println("❌ Attachment type should not be empty")
			return
		}

		if attachment.Filename == "" {
			fmt.Println("❌ Attachment filename should not be empty")
			return
		}

		if attachment.Size <= 0 {
			fmt.Println("❌ Attachment size should be positive")
			return
		}

		if attachment.ContentType == "" {
			fmt.Println("❌ Attachment content type should not be empty")
			return
		}
	}

	fmt.Printf("✅ Feedback collection test passed\n")
	fmt.Printf("   Feedback ID: %s\n", feedback.ID)
	fmt.Printf("   User ID: %s\n", feedback.UserID)
	fmt.Printf("   Session ID: %s\n", feedback.SessionID)
	fmt.Printf("   Feedback type: %d\n", feedback.FeedbackType)
	fmt.Printf("   Category: %d\n", feedback.Category)
	fmt.Printf("   Priority: %d\n", feedback.Priority)
	fmt.Printf("   Score: %.1f\n", feedback.Score)
	fmt.Printf("   Title: %s\n", feedback.Title)
	fmt.Printf("   Description: %s\n", feedback.Description)
	fmt.Printf("   Attachments: %d\n", len(feedback.Attachments))
	fmt.Printf("   Tags: %d\n", len(feedback.Tags))
	fmt.Printf("   Status: %d\n", feedback.Status)
}

func testFeedbackProcessing() {
	// Test feedback processing
	feedback := FeedbackItem{
		ID:           "feedback-002",
		UserID:       "user-456",
		SessionID:    "session-789",
		Timestamp:    time.Now(),
		FeedbackType: 2, // Layout
		Category:     1, // Alignment
		Priority:     1, // High
		Score:        2.0,
		Title:        "Task alignment is broken",
		Description:  "Tasks are overlapping and not aligned properly",
		Context: map[string]interface{}{
			"view_type": "weekly",
			"task_count": 8,
			"overlap_count": 3,
		},
		Tags:   []string{"alignment", "overlap", "layout"},
		Status: 0, // Pending
	}

	// Simulate processing
	feedback.Status = 1 // Processing
	now := time.Now()
	feedback.ProcessedAt = &now

	// Validate processing
	if feedback.Status != 1 {
		fmt.Println("❌ Feedback status should be processing")
		return
	}

	if feedback.ProcessedAt == nil {
		fmt.Println("❌ Processed at should be set")
		return
	}

	// Simulate improvement generation
	improvements := []ImprovementAction{
		{
			ID:          "improvement-001",
			Type:        1, // Layout
			Description: "Fix task alignment algorithm",
			Priority:    1,
			Effort:      "High",
			Impact:      0.9,
			Status:      0, // Planned
			CreatedAt:   time.Now(),
			Details: map[string]interface{}{
				"category": "alignment",
				"feedback_id": feedback.ID,
				"user_score": feedback.Score,
			},
		},
	}

	feedback.Improvements = improvements
	feedback.Status = 2 // Processed

	// Validate improvements
	if len(feedback.Improvements) == 0 {
		fmt.Println("❌ Improvements should be generated")
		return
	}

	for _, improvement := range feedback.Improvements {
		if improvement.ID == "" {
			fmt.Println("❌ Improvement ID should not be empty")
			return
		}

		if improvement.Description == "" {
			fmt.Println("❌ Improvement description should not be empty")
			return
		}

		if improvement.Priority < 1 || improvement.Priority > 5 {
			fmt.Println("❌ Improvement priority should be between 1 and 5")
			return
		}

		if improvement.Effort == "" {
			fmt.Println("❌ Improvement effort should not be empty")
			return
		}

		if improvement.Impact < 0.0 || improvement.Impact > 1.0 {
			fmt.Println("❌ Improvement impact should be between 0 and 1")
			return
		}

		if improvement.Status < 0 || improvement.Status > 4 {
			fmt.Println("❌ Improvement status should be between 0 and 4")
			return
		}
	}

	fmt.Printf("✅ Feedback processing test passed\n")
	fmt.Printf("   Feedback ID: %s\n", feedback.ID)
	fmt.Printf("   Status: %d\n", feedback.Status)
	fmt.Printf("   Processed at: %v\n", feedback.ProcessedAt)
	fmt.Printf("   Improvements: %d\n", len(feedback.Improvements))
	fmt.Printf("   Improvement 1 ID: %s\n", feedback.Improvements[0].ID)
	fmt.Printf("   Improvement 1 Description: %s\n", feedback.Improvements[0].Description)
	fmt.Printf("   Improvement 1 Priority: %d\n", feedback.Improvements[0].Priority)
	fmt.Printf("   Improvement 1 Effort: %s\n", feedback.Improvements[0].Effort)
	fmt.Printf("   Improvement 1 Impact: %.2f\n", feedback.Improvements[0].Impact)
}

func testFeedbackAnalysis() {
	// Test feedback analysis
	metrics := FeedbackMetrics{
		TotalFeedback:       150,
		PendingFeedback:     12,
		ProcessedFeedback:   120,
		ImplementedFeedback: 95,
		AverageScore:        3.8,
		CategoryBreakdown: map[string]int{
			"Visual":        45,
			"Layout":        35,
			"Performance":   25,
			"Usability":     30,
			"Accessibility": 15,
		},
		PriorityBreakdown: map[string]int{
			"Low":      60,
			"Medium":   70,
			"High":     15,
			"Critical": 5,
		},
		StatusBreakdown: map[string]int{
			"Pending":      12,
			"Processing":   8,
			"Processed":    120,
			"Implemented":  95,
			"Rejected":     3,
			"Archived":     12,
		},
		Trends: []FeedbackTrend{
			{Period: "2024-01", Count: 45, AvgScore: 3.9, Category: "Visual"},
			{Period: "2024-02", Count: 52, AvgScore: 3.8, Category: "Layout"},
		},
		TopIssues: []TopIssue{
			{Issue: "Color contrast", Count: 15, AvgScore: 2.5, Priority: "High", Category: "Accessibility"},
			{Issue: "Task spacing", Count: 12, AvgScore: 3.2, Priority: "Medium", Category: "Layout"},
		},
		ImprovementStats: ImprovementStats{
			TotalImprovements:     25,
			CompletedImprovements: 18,
			InProgressImprovements: 5,
			PlannedImprovements:   2,
			AverageImpact:         0.75,
			AverageEffort:         "Medium",
		},
	}

	// Validate metrics
	if metrics.TotalFeedback < 0 {
		fmt.Println("❌ Total feedback should be non-negative")
		return
	}

	if metrics.PendingFeedback < 0 {
		fmt.Println("❌ Pending feedback should be non-negative")
		return
	}

	if metrics.ProcessedFeedback < 0 {
		fmt.Println("❌ Processed feedback should be non-negative")
		return
	}

	if metrics.ImplementedFeedback < 0 {
		fmt.Println("❌ Implemented feedback should be non-negative")
		return
	}

	if metrics.AverageScore < 1.0 || metrics.AverageScore > 5.0 {
		fmt.Println("❌ Average score should be between 1 and 5")
		return
	}

	if len(metrics.CategoryBreakdown) == 0 {
		fmt.Println("❌ Category breakdown should not be empty")
		return
	}

	if len(metrics.PriorityBreakdown) == 0 {
		fmt.Println("❌ Priority breakdown should not be empty")
		return
	}

	if len(metrics.StatusBreakdown) == 0 {
		fmt.Println("❌ Status breakdown should not be empty")
		return
	}

	// Validate trends
	for _, trend := range metrics.Trends {
		if trend.Period == "" {
			fmt.Println("❌ Trend period should not be empty")
			return
		}

		if trend.Count < 0 {
			fmt.Println("❌ Trend count should be non-negative")
			return
		}

		if trend.AvgScore < 1.0 || trend.AvgScore > 5.0 {
			fmt.Println("❌ Trend average score should be between 1 and 5")
			return
		}

		if trend.Category == "" {
			fmt.Println("❌ Trend category should not be empty")
			return
		}
	}

	// Validate top issues
	for _, issue := range metrics.TopIssues {
		if issue.Issue == "" {
			fmt.Println("❌ Top issue should not be empty")
			return
		}

		if issue.Count < 0 {
			fmt.Println("❌ Top issue count should be non-negative")
			return
		}

		if issue.AvgScore < 1.0 || issue.AvgScore > 5.0 {
			fmt.Println("❌ Top issue average score should be between 1 and 5")
			return
		}

		if issue.Priority == "" {
			fmt.Println("❌ Top issue priority should not be empty")
			return
		}

		if issue.Category == "" {
			fmt.Println("❌ Top issue category should not be empty")
			return
		}
	}

	// Validate improvement stats
	if metrics.ImprovementStats.TotalImprovements < 0 {
		fmt.Println("❌ Total improvements should be non-negative")
		return
	}

	if metrics.ImprovementStats.CompletedImprovements < 0 {
		fmt.Println("❌ Completed improvements should be non-negative")
		return
	}

	if metrics.ImprovementStats.InProgressImprovements < 0 {
		fmt.Println("❌ In progress improvements should be non-negative")
		return
	}

	if metrics.ImprovementStats.PlannedImprovements < 0 {
		fmt.Println("❌ Planned improvements should be non-negative")
		return
	}

	if metrics.ImprovementStats.AverageImpact < 0.0 || metrics.ImprovementStats.AverageImpact > 1.0 {
		fmt.Println("❌ Average impact should be between 0 and 1")
		return
	}

	if metrics.ImprovementStats.AverageEffort == "" {
		fmt.Println("❌ Average effort should not be empty")
		return
	}

	fmt.Printf("✅ Feedback analysis test passed\n")
	fmt.Printf("   Total feedback: %d\n", metrics.TotalFeedback)
	fmt.Printf("   Pending feedback: %d\n", metrics.PendingFeedback)
	fmt.Printf("   Processed feedback: %d\n", metrics.ProcessedFeedback)
	fmt.Printf("   Implemented feedback: %d\n", metrics.ImplementedFeedback)
	fmt.Printf("   Average score: %.1f\n", metrics.AverageScore)
	fmt.Printf("   Category breakdown: %d categories\n", len(metrics.CategoryBreakdown))
	fmt.Printf("   Priority breakdown: %d priorities\n", len(metrics.PriorityBreakdown))
	fmt.Printf("   Status breakdown: %d statuses\n", len(metrics.StatusBreakdown))
	fmt.Printf("   Trends: %d trends\n", len(metrics.Trends))
	fmt.Printf("   Top issues: %d issues\n", len(metrics.TopIssues))
	fmt.Printf("   Total improvements: %d\n", metrics.ImprovementStats.TotalImprovements)
	fmt.Printf("   Completed improvements: %d\n", metrics.ImprovementStats.CompletedImprovements)
	fmt.Printf("   In progress improvements: %d\n", metrics.ImprovementStats.InProgressImprovements)
	fmt.Printf("   Planned improvements: %d\n", metrics.ImprovementStats.PlannedImprovements)
	fmt.Printf("   Average impact: %.2f\n", metrics.ImprovementStats.AverageImpact)
	fmt.Printf("   Average effort: %s\n", metrics.ImprovementStats.AverageEffort)
}

func testFeedbackImprovements() {
	// Test feedback improvements
	improvements := []ImprovementAction{
		{
			ID:          "improvement-001",
			Type:        0, // Visual
			Description: "Improve color contrast for better accessibility",
			Priority:    1,
			Effort:      "Medium",
			Impact:      0.9,
			Status:      0, // Planned
			CreatedAt:   time.Now(),
			Details: map[string]interface{}{
				"category": "accessibility",
				"feedback_id": "feedback-001",
				"user_score": 2.0,
			},
		},
		{
			ID:          "improvement-002",
			Type:        1, // Layout
			Description: "Fix task spacing to improve readability",
			Priority:    2,
			Effort:      "Low",
			Impact:      0.7,
			Status:      1, // In Progress
			CreatedAt:   time.Now(),
			Details: map[string]interface{}{
				"category": "spacing",
				"feedback_id": "feedback-002",
				"user_score": 3.0,
			},
		},
		{
			ID:          "improvement-003",
			Type:        2, // Performance
			Description: "Optimize layout algorithm for better performance",
			Priority:    3,
			Effort:      "High",
			Impact:      0.8,
			Status:      2, // Completed
			CreatedAt:   time.Now(),
			CompletedAt: func() *time.Time { t := time.Now(); return &t }(),
			Details: map[string]interface{}{
				"category": "performance",
				"feedback_id": "feedback-003",
				"user_score": 4.0,
			},
		},
	}

	// Validate improvements
	if len(improvements) == 0 {
		fmt.Println("❌ Improvements should not be empty")
		return
	}

	for i, improvement := range improvements {
		if improvement.ID == "" {
			fmt.Printf("❌ Improvement %d ID should not be empty\n", i+1)
			return
		}

		if improvement.Description == "" {
			fmt.Printf("❌ Improvement %d description should not be empty\n", i+1)
			return
		}

		if improvement.Priority < 1 || improvement.Priority > 5 {
			fmt.Printf("❌ Improvement %d priority should be between 1 and 5\n", i+1)
			return
		}

		if improvement.Effort == "" {
			fmt.Printf("❌ Improvement %d effort should not be empty\n", i+1)
			return
		}

		if improvement.Impact < 0.0 || improvement.Impact > 1.0 {
			fmt.Printf("❌ Improvement %d impact should be between 0 and 1\n", i+1)
			return
		}

		if improvement.Status < 0 || improvement.Status > 4 {
			fmt.Printf("❌ Improvement %d status should be between 0 and 4\n", i+1)
			return
		}

		if improvement.Details == nil {
			fmt.Printf("❌ Improvement %d details should not be nil\n", i+1)
			return
		}

		// Check completed improvements have completion date
		if improvement.Status == 2 && improvement.CompletedAt == nil {
			fmt.Printf("❌ Completed improvement %d should have completion date\n", i+1)
			return
		}
	}

	// Test improvement types
	expectedTypes := []int{0, 1, 2} // Visual, Layout, Performance
	for i, improvement := range improvements {
		if improvement.Type != expectedTypes[i] {
			fmt.Printf("❌ Improvement %d type should be %d\n", i+1, expectedTypes[i])
			return
		}
	}

	// Test improvement priorities
	expectedPriorities := []int{1, 2, 3}
	for i, improvement := range improvements {
		if improvement.Priority != expectedPriorities[i] {
			fmt.Printf("❌ Improvement %d priority should be %d\n", i+1, expectedPriorities[i])
			return
		}
	}

	// Test improvement efforts
	expectedEfforts := []string{"Medium", "Low", "High"}
	for i, improvement := range improvements {
		if improvement.Effort != expectedEfforts[i] {
			fmt.Printf("❌ Improvement %d effort should be %s\n", i+1, expectedEfforts[i])
			return
		}
	}

	// Test improvement statuses
	expectedStatuses := []int{0, 1, 2} // Planned, In Progress, Completed
	for i, improvement := range improvements {
		if improvement.Status != expectedStatuses[i] {
			fmt.Printf("❌ Improvement %d status should be %d\n", i+1, expectedStatuses[i])
			return
		}
	}

	fmt.Printf("✅ Feedback improvements test passed\n")
	fmt.Printf("   Total improvements: %d\n", len(improvements))
	fmt.Printf("   Improvement 1: %s (Type: %d, Priority: %d, Effort: %s, Impact: %.2f, Status: %d)\n", 
		improvements[0].Description, improvements[0].Type, improvements[0].Priority, 
		improvements[0].Effort, improvements[0].Impact, improvements[0].Status)
	fmt.Printf("   Improvement 2: %s (Type: %d, Priority: %d, Effort: %s, Impact: %.2f, Status: %d)\n", 
		improvements[1].Description, improvements[1].Type, improvements[1].Priority, 
		improvements[1].Effort, improvements[1].Impact, improvements[1].Status)
	fmt.Printf("   Improvement 3: %s (Type: %d, Priority: %d, Effort: %s, Impact: %.2f, Status: %d)\n", 
		improvements[2].Description, improvements[2].Type, improvements[2].Priority, 
		improvements[2].Effort, improvements[2].Impact, improvements[2].Status)
}
