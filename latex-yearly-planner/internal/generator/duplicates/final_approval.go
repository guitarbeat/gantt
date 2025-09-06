package generator

import (
	"fmt"
	"time"
)

// FinalApproval provides comprehensive final approval and quality validation
type FinalApproval struct {
	config        *FinalApprovalConfig
	validator     *ApprovalValidator
	workflow      *ApprovalWorkflow
	notifier      *ApprovalNotifier
	reporter      *ApprovalReporter
	logger        PDFLogger
}

// FinalApprovalConfig defines configuration for final approval
type FinalApprovalConfig struct {
	// Approval settings
	EnableMultiLevelApproval bool          `json:"enable_multi_level_approval"`
	RequireStakeholderApproval bool        `json:"require_stakeholder_approval"`
	RequireTechnicalApproval  bool         `json:"require_technical_approval"`
	RequireUserApproval       bool         `json:"require_user_approval"`
	ApprovalTimeout           time.Duration `json:"approval_timeout"`
	MaxApprovalAttempts       int           `json:"max_approval_attempts"`
	
	// Quality gates
	MinQualityScore           float64 `json:"min_quality_score"`
	MinTestCoverage           float64 `json:"min_test_coverage"`
	MinPerformanceScore       float64 `json:"min_performance_score"`
	MaxCriticalIssues         int     `json:"max_critical_issues"`
	MaxHighIssues             int     `json:"max_high_issues"`
	MinUserSatisfactionScore  float64 `json:"min_user_satisfaction_score"`
	
	// Validation settings
	EnableComprehensiveValidation bool `json:"enable_comprehensive_validation"`
	EnableRegressionTesting       bool `json:"enable_regression_testing"`
	EnablePerformanceValidation   bool `json:"enable_performance_validation"`
	EnableSecurityValidation      bool `json:"enable_security_validation"`
	EnableAccessibilityValidation bool `json:"enable_accessibility_validation"`
	
	// Notification settings
	EnableApprovalNotifications bool     `json:"enable_approval_notifications"`
	NotificationChannels        []string `json:"notification_channels"`
	NotificationTemplate        string   `json:"notification_template"`
	
	// Reporting settings
	EnableDetailedReports       bool   `json:"enable_detailed_reports"`
	ReportFormat                string `json:"report_format"`
	ReportOutputPath            string `json:"report_output_path"`
	IncludeQualityMetrics       bool   `json:"include_quality_metrics"`
	IncludePerformanceMetrics   bool   `json:"include_performance_metrics"`
	IncludeUserFeedback         bool   `json:"include_user_feedback"`
}

// ApprovalValidator provides approval validation
type ApprovalValidator struct {
	config *FinalApprovalConfig
	logger PDFLogger
}

// ApprovalWorkflow manages the approval workflow
type ApprovalWorkflow struct {
	config     *FinalApprovalConfig
	steps      []ApprovalStep
	currentStep int
	status     WorkflowStatus
	logger     PDFLogger
}

// ApprovalStep represents a step in the approval workflow
type ApprovalStep struct {
	ID          string
	Name        string
	Description string
	Type        ApprovalStepType
	Required    bool
	Approvers   []Approver
	Status      StepStatus
	StartTime   *time.Time
	EndTime     *time.Time
	Duration    time.Duration
	Comments    []ApprovalComment
	Attachments []ApprovalAttachment
}

// ApprovalStepType represents the type of approval step
type ApprovalStepType int

const (
	ApprovalStepTypeQualityReview ApprovalStepType = iota
	ApprovalStepTypeTechnicalReview
	ApprovalStepTypeUserAcceptance
	ApprovalStepTypeStakeholderApproval
	ApprovalStepTypeSecurityReview
	ApprovalStepTypePerformanceReview
	ApprovalStepTypeAccessibilityReview
	ApprovalStepTypeFinalSignOff
)

// Approver represents an approver
type Approver struct {
	ID          string
	Name        string
	Email       string
	Role        string
	Department  string
	Level       ApprovalLevel
	Permissions []string
	Active      bool
}

// ApprovalLevel represents the level of approval authority
type ApprovalLevel int

const (
	ApprovalLevelUser ApprovalLevel = iota
	ApprovalLevelTester
	ApprovalLevelDeveloper
	ApprovalLevelManager
	ApprovalLevelDirector
	ApprovalLevelExecutive
	ApprovalLevelStakeholder
)

// StepStatus represents the status of an approval step
type StepStatus int

const (
	StepStatusPending StepStatus = iota
	StepStatusInProgress
	StepStatusApproved
	StepStatusRejected
	StepStatusSkipped
	StepStatusExpired
)

// WorkflowStatus represents the status of the approval workflow
type WorkflowStatus int

const (
	WorkflowStatusNotStarted WorkflowStatus = iota
	WorkflowStatusInProgress
	WorkflowStatusCompleted
	WorkflowStatusFailed
	WorkflowStatusCancelled
	WorkflowStatusExpired
)

// ApprovalComment represents a comment on an approval step
type ApprovalComment struct {
	ID        string
	Author    string
	Content   string
	Type      CommentType
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CommentType represents the type of comment
type CommentType int

const (
	CommentTypeApproval CommentType = iota
	CommentTypeRejection
	CommentTypeQuestion
	CommentTypeSuggestion
	CommentTypeConcern
	CommentTypeClarification
)

// ApprovalAttachment represents an attachment to an approval step
type ApprovalAttachment struct {
	ID          string
	Name        string
	Type        string
	Size        int64
	URL         string
	Description string
	CreatedAt   time.Time
}

// ApprovalNotifier provides approval notifications
type ApprovalNotifier struct {
	config *FinalApprovalConfig
	logger PDFLogger
}

// ApprovalReporter provides approval reporting
type ApprovalReporter struct {
	config *FinalApprovalConfig
	logger PDFLogger
}

// ApprovalRequest represents an approval request
type ApprovalRequest struct {
	ID              string
	Type            ApprovalRequestType
	Title           string
	Description     string
	Priority        ApprovalPriority
	Status          ApprovalRequestStatus
	RequestedBy     string
	RequestedAt     time.Time
	DueDate         *time.Time
	ApprovedBy      string
	ApprovedAt      *time.Time
	RejectedBy      string
	RejectedAt      *time.Time
	Comments        []ApprovalComment
	Attachments     []ApprovalAttachment
	QualityMetrics  *QualityMetrics
	PerformanceMetrics *PerformanceMetrics
	UserFeedback    *UserFeedback
	Metadata        map[string]interface{}
}

// ApprovalRequestType represents the type of approval request
type ApprovalRequestType int

const (
	ApprovalRequestTypeFeature ApprovalRequestType = iota
	ApprovalRequestTypeRelease
	ApprovalRequestTypeQuality
	ApprovalRequestTypeSecurity
	ApprovalRequestTypePerformance
	ApprovalRequestTypeUserAcceptance
	ApprovalRequestTypeFinal
)

// ApprovalPriority represents the priority of an approval request
type ApprovalPriority int

const (
	ApprovalPriorityCritical ApprovalPriority = iota
	ApprovalPriorityHigh
	ApprovalPriorityMedium
	ApprovalPriorityLow
)

// ApprovalRequestStatus represents the status of an approval request
type ApprovalRequestStatus int

const (
	ApprovalRequestStatusPending ApprovalRequestStatus = iota
	ApprovalRequestStatusInProgress
	ApprovalRequestStatusApproved
	ApprovalRequestStatusRejected
	ApprovalRequestStatusExpired
	ApprovalRequestStatusCancelled
)

// QualityMetrics represents quality metrics
type QualityMetrics struct {
	OverallScore     float64 `json:"overall_score"`
	TestCoverage     float64 `json:"test_coverage"`
	PerformanceScore float64 `json:"performance_score"`
	SecurityScore    float64 `json:"security_score"`
	UsabilityScore   float64 `json:"usability_score"`
	AccessibilityScore float64 `json:"accessibility_score"`
	BugCount         int     `json:"bug_count"`
	CriticalIssues   int     `json:"critical_issues"`
	HighIssues       int     `json:"high_issues"`
	MediumIssues     int     `json:"medium_issues"`
	LowIssues        int     `json:"low_issues"`
	ResolutionRate   float64 `json:"resolution_rate"`
	LastUpdated      time.Time `json:"last_updated"`
}

// PerformanceMetrics represents performance metrics
type PerformanceMetrics struct {
	ResponseTime    time.Duration `json:"response_time"`
	Throughput      float64       `json:"throughput"`
	MemoryUsage     int64         `json:"memory_usage"`
	CPUUsage        float64       `json:"cpu_usage"`
	ErrorRate       float64       `json:"error_rate"`
	Availability    float64       `json:"availability"`
	Scalability     float64       `json:"scalability"`
	LastUpdated     time.Time     `json:"last_updated"`
}

// UserFeedback represents user feedback
type UserFeedback struct {
	OverallSatisfaction float64 `json:"overall_satisfaction"`
	UsabilityScore      float64 `json:"usability_score"`
	PerformanceScore    float64 `json:"performance_score"`
	FeatureCompleteness float64 `json:"feature_completeness"`
	RecommendationScore float64 `json:"recommendation_score"`
	Comments            []string `json:"comments"`
	Suggestions         []string `json:"suggestions"`
	Issues              []string `json:"issues"`
	LastUpdated         time.Time `json:"last_updated"`
}

// FinalApprovalResult represents the result of final approval
type FinalApprovalResult struct {
	Status              ApprovalStatus
	OverallScore        float64
	QualityScore        float64
	PerformanceScore    float64
	UserSatisfactionScore float64
	ApprovalSteps       []ApprovalStepResult
	Issues              []ApprovalIssue
	Recommendations     []ApprovalRecommendation
	GeneratedAt         time.Time
}

// ApprovalStatus represents the status of approval
type ApprovalStatus int

const (
	ApprovalStatusApproved ApprovalStatus = iota
	ApprovalStatusRejected
	ApprovalStatusConditional
	ApprovalStatusPending
	ApprovalStatusExpired
)

// ApprovalStepResult represents the result of an approval step
type ApprovalStepResult struct {
	StepID      string
	StepName    string
	Status      StepStatus
	Score       float64
	Comments    []string
	Issues      []string
	Recommendations []string
	CompletedAt time.Time
}

// ApprovalIssue represents an issue found during approval
type ApprovalIssue struct {
	ID          string
	Type        string
	Severity    string
	Description string
	Impact      string
	Resolution  string
	Status      string
	CreatedAt   time.Time
}

// ApprovalRecommendation represents a recommendation from approval
type ApprovalRecommendation struct {
	ID          string
	Type        string
	Priority    int
	Title       string
	Description string
	Benefits    []string
	Effort      string
	Impact      string
	Status      string
	CreatedAt   time.Time
}

// NewFinalApproval creates a new final approval system
func NewFinalApproval() *FinalApproval {
	return &FinalApproval{
		config:    GetDefaultFinalApprovalConfig(),
		validator: NewApprovalValidator(),
		workflow:  NewApprovalWorkflow(),
		notifier:  NewApprovalNotifier(),
		reporter:  NewApprovalReporter(),
		logger:    &FinalApprovalLogger{},
	}
}

// GetDefaultFinalApprovalConfig returns the default final approval configuration
func GetDefaultFinalApprovalConfig() *FinalApprovalConfig {
	return &FinalApprovalConfig{
		EnableMultiLevelApproval:     true,
		RequireStakeholderApproval:   true,
		RequireTechnicalApproval:     true,
		RequireUserApproval:          true,
		ApprovalTimeout:             time.Hour * 72,
		MaxApprovalAttempts:         3,
		MinQualityScore:             0.9,
		MinTestCoverage:             0.85,
		MinPerformanceScore:         0.85,
		MaxCriticalIssues:           0,
		MaxHighIssues:               5,
		MinUserSatisfactionScore:    0.8,
		EnableComprehensiveValidation: true,
		EnableRegressionTesting:      true,
		EnablePerformanceValidation:  true,
		EnableSecurityValidation:     true,
		EnableAccessibilityValidation: true,
		EnableApprovalNotifications:  true,
		NotificationChannels:        []string{"email", "slack", "webhook"},
		NotificationTemplate:        "default",
		EnableDetailedReports:       true,
		ReportFormat:                "json",
		ReportOutputPath:            "./approval_reports/",
		IncludeQualityMetrics:       true,
		IncludePerformanceMetrics:   true,
		IncludeUserFeedback:         true,
	}
}

// SetLogger sets the logger for final approval
func (fa *FinalApproval) SetLogger(logger PDFLogger) {
	fa.logger = logger
	fa.validator.SetLogger(logger)
	fa.workflow.SetLogger(logger)
	fa.notifier.SetLogger(logger)
	fa.reporter.SetLogger(logger)
}

// RequestFinalApproval requests final approval for the system
func (fa *FinalApproval) RequestFinalApproval(request *ApprovalRequest) (*FinalApprovalResult, error) {
	fa.logger.Info("Requesting final approval for: %s", request.Title)
	
	// Validate approval request
	if err := fa.validator.ValidateApprovalRequest(request); err != nil {
		return nil, fmt.Errorf("approval request validation failed: %w", err)
	}
	
	// Initialize approval workflow
	if err := fa.workflow.InitializeWorkflow(request); err != nil {
		return nil, fmt.Errorf("workflow initialization failed: %w", err)
	}
	
	// Execute approval workflow
	result, err := fa.workflow.ExecuteWorkflow()
	if err != nil {
		return nil, fmt.Errorf("workflow execution failed: %w", err)
	}
	
	// Send notifications
	if fa.config.EnableApprovalNotifications {
		if err := fa.notifier.SendApprovalNotifications(request, result); err != nil {
			fa.logger.Error("Failed to send approval notifications: %v", err)
		}
	}
	
	// Generate approval report
	if fa.config.EnableDetailedReports {
		if err := fa.reporter.GenerateApprovalReport(request, result); err != nil {
			fa.logger.Error("Failed to generate approval report: %v", err)
		}
	}
	
	fa.logger.Info("Final approval request completed with status: %v", result.Status)
	return result, nil
}

// ValidateQualityGates validates all quality gates
func (fa *FinalApproval) ValidateQualityGates(metrics *QualityMetrics) (*QualityGateResult, error) {
	fa.logger.Info("Validating quality gates...")
	
	result := &QualityGateResult{
		OverallPassed: true,
		Gates:         []QualityGate{},
		GeneratedAt:   time.Now(),
	}
	
	// Validate overall quality score
	if metrics.OverallScore < fa.config.MinQualityScore {
		result.OverallPassed = false
		result.Gates = append(result.Gates, QualityGate{
			Name:        "Overall Quality Score",
			Passed:      false,
			Expected:    fmt.Sprintf("%.2f%%", fa.config.MinQualityScore*100),
			Actual:      fmt.Sprintf("%.2f%%", metrics.OverallScore*100),
			Description: "Overall quality score below minimum threshold",
		})
	} else {
		result.Gates = append(result.Gates, QualityGate{
			Name:        "Overall Quality Score",
			Passed:      true,
			Expected:    fmt.Sprintf("%.2f%%", fa.config.MinQualityScore*100),
			Actual:      fmt.Sprintf("%.2f%%", metrics.OverallScore*100),
			Description: "Overall quality score meets requirements",
		})
	}
	
	// Validate test coverage
	if metrics.TestCoverage < fa.config.MinTestCoverage {
		result.OverallPassed = false
		result.Gates = append(result.Gates, QualityGate{
			Name:        "Test Coverage",
			Passed:      false,
			Expected:    fmt.Sprintf("%.2f%%", fa.config.MinTestCoverage*100),
			Actual:      fmt.Sprintf("%.2f%%", metrics.TestCoverage*100),
			Description: "Test coverage below minimum threshold",
		})
	} else {
		result.Gates = append(result.Gates, QualityGate{
			Name:        "Test Coverage",
			Passed:      true,
			Expected:    fmt.Sprintf("%.2f%%", fa.config.MinTestCoverage*100),
			Actual:      fmt.Sprintf("%.2f%%", metrics.TestCoverage*100),
			Description: "Test coverage meets requirements",
		})
	}
	
	// Validate performance score
	if metrics.PerformanceScore < fa.config.MinPerformanceScore {
		result.OverallPassed = false
		result.Gates = append(result.Gates, QualityGate{
			Name:        "Performance Score",
			Passed:      false,
			Expected:    fmt.Sprintf("%.2f%%", fa.config.MinPerformanceScore*100),
			Actual:      fmt.Sprintf("%.2f%%", metrics.PerformanceScore*100),
			Description: "Performance score below minimum threshold",
		})
	} else {
		result.Gates = append(result.Gates, QualityGate{
			Name:        "Performance Score",
			Passed:      true,
			Expected:    fmt.Sprintf("%.2f%%", fa.config.MinPerformanceScore*100),
			Actual:      fmt.Sprintf("%.2f%%", metrics.PerformanceScore*100),
			Description: "Performance score meets requirements",
		})
	}
	
	// Validate critical issues
	if metrics.CriticalIssues > fa.config.MaxCriticalIssues {
		result.OverallPassed = false
		result.Gates = append(result.Gates, QualityGate{
			Name:        "Critical Issues",
			Passed:      false,
			Expected:    fmt.Sprintf("%d", fa.config.MaxCriticalIssues),
			Actual:      fmt.Sprintf("%d", metrics.CriticalIssues),
			Description: "Number of critical issues exceeds maximum threshold",
		})
	} else {
		result.Gates = append(result.Gates, QualityGate{
			Name:        "Critical Issues",
			Passed:      true,
			Expected:    fmt.Sprintf("%d", fa.config.MaxCriticalIssues),
			Actual:      fmt.Sprintf("%d", metrics.CriticalIssues),
			Description: "Critical issues within acceptable limits",
		})
	}
	
	// Validate high issues
	if metrics.HighIssues > fa.config.MaxHighIssues {
		result.OverallPassed = false
		result.Gates = append(result.Gates, QualityGate{
			Name:        "High Issues",
			Passed:      false,
			Expected:    fmt.Sprintf("%d", fa.config.MaxHighIssues),
			Actual:      fmt.Sprintf("%d", metrics.HighIssues),
			Description: "Number of high issues exceeds maximum threshold",
		})
	} else {
		result.Gates = append(result.Gates, QualityGate{
			Name:        "High Issues",
			Passed:      true,
			Expected:    fmt.Sprintf("%d", fa.config.MaxHighIssues),
			Actual:      fmt.Sprintf("%d", metrics.HighIssues),
			Description: "High issues within acceptable limits",
		})
	}
	
	fa.logger.Info("Quality gates validation completed: %v", result.OverallPassed)
	return result, nil
}

// GenerateFinalReport generates the final approval report
func (fa *FinalApproval) GenerateFinalReport(result *FinalApprovalResult) (*FinalApprovalReport, error) {
	fa.logger.Info("Generating final approval report...")
	
	report := &FinalApprovalReport{
		Summary: FinalApprovalReportSummary{
			Status:              result.Status,
			OverallScore:        result.OverallScore,
			QualityScore:        result.QualityScore,
			PerformanceScore:    result.PerformanceScore,
			UserSatisfactionScore: result.UserSatisfactionScore,
			ApprovalSteps:       len(result.ApprovalSteps),
			Issues:              len(result.Issues),
			Recommendations:     len(result.Recommendations),
			GeneratedAt:         time.Now(),
		},
		Result:         result,
		QualityGates:  fa.generateQualityGatesSummary(result),
		ApprovalSteps: fa.generateApprovalStepsSummary(result.ApprovalSteps),
		Issues:        result.Issues,
		Recommendations: result.Recommendations,
		Metadata: map[string]interface{}{
			"version":     "1.0.0",
			"environment": "production",
			"generator":   "final_approval",
		},
	}
	
	fa.logger.Info("Final approval report generated successfully")
	return report, nil
}

// Helper methods
func (fa *FinalApproval) generateQualityGatesSummary(result *FinalApprovalResult) []QualityGateSummary {
	gates := []QualityGateSummary{
		{
			Name:        "Overall Quality",
			Passed:      result.QualityScore >= fa.config.MinQualityScore,
			Score:       result.QualityScore,
			Threshold:   fa.config.MinQualityScore,
		},
		{
			Name:        "Performance",
			Passed:      result.PerformanceScore >= fa.config.MinPerformanceScore,
			Score:       result.PerformanceScore,
			Threshold:   fa.config.MinPerformanceScore,
		},
		{
			Name:        "User Satisfaction",
			Passed:      result.UserSatisfactionScore >= fa.config.MinUserSatisfactionScore,
			Score:       result.UserSatisfactionScore,
			Threshold:   fa.config.MinUserSatisfactionScore,
		},
	}
	
	return gates
}

func (fa *FinalApproval) generateApprovalStepsSummary(steps []ApprovalStepResult) []ApprovalStepSummary {
	summaries := []ApprovalStepSummary{}
	
	for _, step := range steps {
		summaries = append(summaries, ApprovalStepSummary{
			StepID:      step.StepID,
			StepName:    step.StepName,
			Status:      step.Status,
			Score:       step.Score,
			Comments:    len(step.Comments),
			Issues:      len(step.Issues),
			Recommendations: len(step.Recommendations),
			CompletedAt: step.CompletedAt,
		})
	}
	
	return summaries
}

// Additional types for completeness
type QualityGateResult struct {
	OverallPassed bool
	Gates         []QualityGate
	GeneratedAt   time.Time
}

type QualityGate struct {
	Name        string
	Passed      bool
	Expected    string
	Actual      string
	Description string
}

type QualityGateSummary struct {
	Name      string
	Passed    bool
	Score     float64
	Threshold float64
}

type ApprovalStepSummary struct {
	StepID         string
	StepName       string
	Status         StepStatus
	Score          float64
	Comments       int
	Issues         int
	Recommendations int
	CompletedAt    time.Time
}

type FinalApprovalReport struct {
	Summary        FinalApprovalReportSummary
	Result         *FinalApprovalResult
	QualityGates   []QualityGateSummary
	ApprovalSteps  []ApprovalStepSummary
	Issues         []ApprovalIssue
	Recommendations []ApprovalRecommendation
	Metadata       map[string]interface{}
}

type FinalApprovalReportSummary struct {
	Status              ApprovalStatus
	OverallScore        float64
	QualityScore        float64
	PerformanceScore    float64
	UserSatisfactionScore float64
	ApprovalSteps       int
	Issues              int
	Recommendations     int
	GeneratedAt         time.Time
}

// Placeholder implementations for missing methods
func NewApprovalValidator() *ApprovalValidator {
	return &ApprovalValidator{
		config: GetDefaultFinalApprovalConfig(),
		logger: &FinalApprovalLogger{},
	}
}

func (av *ApprovalValidator) SetLogger(logger PDFLogger) {
	av.logger = logger
}

func (av *ApprovalValidator) ValidateApprovalRequest(request *ApprovalRequest) error {
	// Placeholder implementation
	return nil
}

func NewApprovalWorkflow() *ApprovalWorkflow {
	return &ApprovalWorkflow{
		config:     GetDefaultFinalApprovalConfig(),
		steps:      []ApprovalStep{},
		currentStep: 0,
		status:     WorkflowStatusNotStarted,
		logger:     &FinalApprovalLogger{},
	}
}

func (aw *ApprovalWorkflow) SetLogger(logger PDFLogger) {
	aw.logger = logger
}

func (aw *ApprovalWorkflow) InitializeWorkflow(request *ApprovalRequest) error {
	// Placeholder implementation
	return nil
}

func (aw *ApprovalWorkflow) ExecuteWorkflow() (*FinalApprovalResult, error) {
	// Placeholder implementation
	return &FinalApprovalResult{
		Status:              ApprovalStatusApproved,
		OverallScore:        0.95,
		QualityScore:        0.92,
		PerformanceScore:    0.88,
		UserSatisfactionScore: 0.90,
		ApprovalSteps:       []ApprovalStepResult{},
		Issues:              []ApprovalIssue{},
		Recommendations:     []ApprovalRecommendation{},
		GeneratedAt:         time.Now(),
	}, nil
}

func NewApprovalNotifier() *ApprovalNotifier {
	return &ApprovalNotifier{
		config: GetDefaultFinalApprovalConfig(),
		logger: &FinalApprovalLogger{},
	}
}

func (an *ApprovalNotifier) SetLogger(logger PDFLogger) {
	an.logger = logger
}

func (an *ApprovalNotifier) SendApprovalNotifications(request *ApprovalRequest, result *FinalApprovalResult) error {
	// Placeholder implementation
	return nil
}

func NewApprovalReporter() *ApprovalReporter {
	return &ApprovalReporter{
		config: GetDefaultFinalApprovalConfig(),
		logger: &FinalApprovalLogger{},
	}
}

func (ar *ApprovalReporter) SetLogger(logger PDFLogger) {
	ar.logger = logger
}

func (ar *ApprovalReporter) GenerateApprovalReport(request *ApprovalRequest, result *FinalApprovalResult) error {
	// Placeholder implementation
	return nil
}

// FinalApprovalLogger provides logging for final approval
type FinalApprovalLogger struct{}

func (l *FinalApprovalLogger) Info(msg string, args ...interface{})  { fmt.Printf("[FA-INFO] "+msg+"\n", args...) }
func (l *FinalApprovalLogger) Error(msg string, args ...interface{}) { fmt.Printf("[FA-ERROR] "+msg+"\n", args...) }
func (l *FinalApprovalLogger) Debug(msg string, args ...interface{}) { fmt.Printf("[FA-DEBUG] "+msg+"\n", args...) }
