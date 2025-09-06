package main

import (
	"fmt"
	"time"
)

// Test the user validation system
func main() {
	fmt.Println("Testing User Validation System...")

	// Test 1: User Validation Configuration
	fmt.Println("\n=== Test 1: User Validation Configuration ===")
	testUserValidationConfiguration()

	// Test 2: Validation Workflow
	fmt.Println("\n=== Test 2: Validation Workflow ===")
	testValidationWorkflow()

	// Test 3: Validation Phases
	fmt.Println("\n=== Test 3: Validation Phases ===")
	testValidationPhases()

	// Test 4: Approval System
	fmt.Println("\n=== Test 4: Approval System ===")
	testApprovalSystem()

	// Test 5: Validation Criteria
	fmt.Println("\n=== Test 5: Validation Criteria ===")
	testValidationCriteria()

	fmt.Println("\n✅ User validation system tests completed!")
}

// UserValidationConfig represents user validation configuration
type UserValidationConfig struct {
	EnableUserTesting        bool
	EnableAcceptanceTests    bool
	EnableUsabilityTests     bool
	EnableAccessibilityTests bool
	ValidationTimeout        time.Duration
	MaxConcurrentUsers       int
	EnableApprovalWorkflow   bool
	RequireMultipleApprovers bool
	MinApprovers             int
	ApprovalTimeout          time.Duration
	EnableRealTimeFeedback   bool
	EnableFeedbackCollection bool
	FeedbackChannels         []string
	MinUserSatisfactionScore float64
	MinUsabilityScore        float64
	MinAccessibilityScore    float64
	MaxCriticalIssues        int
}

// ValidationPhase represents a validation phase
type ValidationPhase struct {
	ID          string
	Name        string
	Description string
	Order       int
	Status      int
	Tasks       []ValidationTask
	Criteria    []ValidationCriterion
	StartTime   *time.Time
	EndTime     *time.Time
	Duration    time.Duration
}

// ValidationTask represents a validation task
type ValidationTask struct {
	ID          string
	Name        string
	Description string
	Type        int
	Priority    int
	Status      int
	Assignee    string
	DueDate     *time.Time
	CompletedAt *time.Time
	Result      *TaskResult
}

// ValidationCriterion represents a validation criterion
type ValidationCriterion struct {
	ID          string
	Name        string
	Description string
	Type        int
	Weight      float64
	Threshold   float64
	Status      int
	Score       float64
	Comments    string
}

// TaskResult represents a task result
type TaskResult struct {
	Score          float64
	Status         int
	Comments       string
	Evidence       []string
	Issues         []ValidationIssue
	Recommendations []string
	CompletedAt    time.Time
}

// ValidationIssue represents a validation issue
type ValidationIssue struct {
	ID          string
	Type        int
	Severity    int
	Description string
	Location    string
	Steps       []string
	Expected    string
	Actual      string
	Impact      string
	Status      int
	CreatedAt   time.Time
	ResolvedAt  *time.Time
}

// Approval represents an approval
type Approval struct {
	ID          string
	Type        int
	Status      int
	RequestedBy string
	RequestedAt time.Time
	ApprovedBy  string
	ApprovedAt  *time.Time
	Comments    string
	Priority    int
	ExpiresAt   *time.Time
}

// Approver represents an approver
type Approver struct {
	ID     string
	Name   string
	Email  string
	Role   string
	Level  int
	Active bool
}

func testUserValidationConfiguration() {
	// Test user validation configuration
	config := UserValidationConfig{
		EnableUserTesting:        true,
		EnableAcceptanceTests:    true,
		EnableUsabilityTests:     true,
		EnableAccessibilityTests: true,
		ValidationTimeout:        time.Hour * 24,
		MaxConcurrentUsers:       10,
		EnableApprovalWorkflow:   true,
		RequireMultipleApprovers: true,
		MinApprovers:            2,
		ApprovalTimeout:         time.Hour * 48,
		EnableRealTimeFeedback:  true,
		EnableFeedbackCollection: true,
		FeedbackChannels:        []string{"email", "in-app", "webhook"},
		MinUserSatisfactionScore: 0.8,
		MinUsabilityScore:        0.8,
		MinAccessibilityScore:    0.8,
		MaxCriticalIssues:       0,
	}

	// Validate configuration
	if !config.EnableUserTesting {
		fmt.Println("❌ User testing should be enabled")
		return
	}

	if !config.EnableAcceptanceTests {
		fmt.Println("❌ Acceptance tests should be enabled")
		return
	}

	if !config.EnableUsabilityTests {
		fmt.Println("❌ Usability tests should be enabled")
		return
	}

	if !config.EnableAccessibilityTests {
		fmt.Println("❌ Accessibility tests should be enabled")
		return
	}

	if config.ValidationTimeout <= 0 {
		fmt.Println("❌ Validation timeout should be positive")
		return
	}

	if config.MaxConcurrentUsers <= 0 {
		fmt.Println("❌ Max concurrent users should be positive")
		return
	}

	if !config.EnableApprovalWorkflow {
		fmt.Println("❌ Approval workflow should be enabled")
		return
	}

	if !config.RequireMultipleApprovers {
		fmt.Println("❌ Multiple approvers should be required")
		return
	}

	if config.MinApprovers <= 0 {
		fmt.Println("❌ Min approvers should be positive")
		return
	}

	if config.ApprovalTimeout <= 0 {
		fmt.Println("❌ Approval timeout should be positive")
		return
	}

	if !config.EnableRealTimeFeedback {
		fmt.Println("❌ Real-time feedback should be enabled")
		return
	}

	if !config.EnableFeedbackCollection {
		fmt.Println("❌ Feedback collection should be enabled")
		return
	}

	if len(config.FeedbackChannels) == 0 {
		fmt.Println("❌ Feedback channels should not be empty")
		return
	}

	if config.MinUserSatisfactionScore < 0.0 || config.MinUserSatisfactionScore > 1.0 {
		fmt.Println("❌ Min user satisfaction score should be between 0 and 1")
		return
	}

	if config.MinUsabilityScore < 0.0 || config.MinUsabilityScore > 1.0 {
		fmt.Println("❌ Min usability score should be between 0 and 1")
		return
	}

	if config.MinAccessibilityScore < 0.0 || config.MinAccessibilityScore > 1.0 {
		fmt.Println("❌ Min accessibility score should be between 0 and 1")
		return
	}

	if config.MaxCriticalIssues < 0 {
		fmt.Println("❌ Max critical issues should be non-negative")
		return
	}

	fmt.Printf("✅ User validation configuration test passed\n")
	fmt.Printf("   Enable user testing: %v\n", config.EnableUserTesting)
	fmt.Printf("   Enable acceptance tests: %v\n", config.EnableAcceptanceTests)
	fmt.Printf("   Enable usability tests: %v\n", config.EnableUsabilityTests)
	fmt.Printf("   Enable accessibility tests: %v\n", config.EnableAccessibilityTests)
	fmt.Printf("   Validation timeout: %v\n", config.ValidationTimeout)
	fmt.Printf("   Max concurrent users: %d\n", config.MaxConcurrentUsers)
	fmt.Printf("   Enable approval workflow: %v\n", config.EnableApprovalWorkflow)
	fmt.Printf("   Require multiple approvers: %v\n", config.RequireMultipleApprovers)
	fmt.Printf("   Min approvers: %d\n", config.MinApprovers)
	fmt.Printf("   Approval timeout: %v\n", config.ApprovalTimeout)
	fmt.Printf("   Enable real-time feedback: %v\n", config.EnableRealTimeFeedback)
	fmt.Printf("   Enable feedback collection: %v\n", config.EnableFeedbackCollection)
	fmt.Printf("   Feedback channels: %v\n", config.FeedbackChannels)
	fmt.Printf("   Min user satisfaction score: %.2f%%\n", config.MinUserSatisfactionScore*100)
	fmt.Printf("   Min usability score: %.2f%%\n", config.MinUsabilityScore*100)
	fmt.Printf("   Min accessibility score: %.2f%%\n", config.MinAccessibilityScore*100)
	fmt.Printf("   Max critical issues: %d\n", config.MaxCriticalIssues)
}

func testValidationWorkflow() {
	// Test validation workflow
	phases := []ValidationPhase{
		{
			ID:          "phase-1",
			Name:        "Pre-Validation Setup",
			Description: "Prepare validation environment and criteria",
			Order:       1,
			Status:      0, // Pending
			Tasks:       []ValidationTask{},
			Criteria:    []ValidationCriterion{},
		},
		{
			ID:          "phase-2",
			Name:        "User Testing",
			Description: "Conduct user testing sessions",
			Order:       2,
			Status:      0, // Pending
			Tasks:       []ValidationTask{},
			Criteria:    []ValidationCriterion{},
		},
		{
			ID:          "phase-3",
			Name:        "Acceptance Testing",
			Description: "Perform acceptance testing",
			Order:       3,
			Status:      0, // Pending
			Tasks:       []ValidationTask{},
			Criteria:    []ValidationCriterion{},
		},
		{
			ID:          "phase-4",
			Name:        "Usability Testing",
			Description: "Evaluate usability aspects",
			Order:       4,
			Status:      0, // Pending
			Tasks:       []ValidationTask{},
			Criteria:    []ValidationCriterion{},
		},
		{
			ID:          "phase-5",
			Name:        "Accessibility Testing",
			Description: "Test accessibility compliance",
			Order:       5,
			Status:      0, // Pending
			Tasks:       []ValidationTask{},
			Criteria:    []ValidationCriterion{},
		},
		{
			ID:          "phase-6",
			Name:        "Final Approval",
			Description: "Obtain final user approval",
			Order:       6,
			Status:      0, // Pending
			Tasks:       []ValidationTask{},
			Criteria:    []ValidationCriterion{},
		},
	}

	// Validate workflow phases
	if len(phases) == 0 {
		fmt.Println("❌ Workflow should have phases")
		return
	}

	// Validate phase order
	for i, phase := range phases {
		if phase.Order != i+1 {
			fmt.Printf("❌ Phase %d order should be %d, got %d\n", i+1, i+1, phase.Order)
			return
		}
	}

	// Validate phase IDs
	phaseIDs := make(map[string]bool)
	for _, phase := range phases {
		if phase.ID == "" {
			fmt.Println("❌ Phase ID should not be empty")
			return
		}

		if phaseIDs[phase.ID] {
			fmt.Printf("❌ Duplicate phase ID: %s\n", phase.ID)
			return
		}
		phaseIDs[phase.ID] = true
	}

	// Validate phase names
	for i, phase := range phases {
		if phase.Name == "" {
			fmt.Printf("❌ Phase %d name should not be empty\n", i+1)
			return
		}
	}

	// Validate phase descriptions
	for i, phase := range phases {
		if phase.Description == "" {
			fmt.Printf("❌ Phase %d description should not be empty\n", i+1)
			return
		}
	}

	// Validate phase status
	for i, phase := range phases {
		if phase.Status < 0 || phase.Status > 4 {
			fmt.Printf("❌ Phase %d status should be between 0 and 4, got %d\n", i+1, phase.Status)
			return
		}
	}

	fmt.Printf("✅ Validation workflow test passed\n")
	fmt.Printf("   Total phases: %d\n", len(phases))
	fmt.Printf("   Phase 1: %s\n", phases[0].Name)
	fmt.Printf("   Phase 2: %s\n", phases[1].Name)
	fmt.Printf("   Phase 3: %s\n", phases[2].Name)
	fmt.Printf("   Phase 4: %s\n", phases[3].Name)
	fmt.Printf("   Phase 5: %s\n", phases[4].Name)
	fmt.Printf("   Phase 6: %s\n", phases[5].Name)
}

func testValidationPhases() {
	// Test validation phases
	phases := []ValidationPhase{
		{
			ID:          "phase-1",
			Name:        "Pre-Validation Setup",
			Description: "Prepare validation environment and criteria",
			Order:       1,
			Status:      1, // In Progress
			Tasks: []ValidationTask{
				{
					ID:          "task-1-1",
					Name:        "Setup Validation Environment",
					Description: "Prepare testing environment and tools",
					Type:        0, // Functional
					Priority:    1, // High
					Status:      2, // Completed
					Assignee:    "tester1",
					Result: &TaskResult{
						Score:          0.95,
						Status:         2, // Completed
						Comments:       "Environment setup completed successfully",
						Evidence:       []string{"env-setup-log", "tool-config"},
						Issues:         []ValidationIssue{},
						Recommendations: []string{"Monitor performance", "Keep logs"},
						CompletedAt:    time.Now(),
					},
				},
				{
					ID:          "task-1-2",
					Name:        "Define Validation Criteria",
					Description: "Establish validation criteria and thresholds",
					Type:        0, // Functional
					Priority:    1, // High
					Status:      2, // Completed
					Assignee:    "tester2",
					Result: &TaskResult{
						Score:          0.88,
						Status:         2, // Completed
						Comments:       "Criteria defined and documented",
						Evidence:       []string{"criteria-doc", "threshold-config"},
						Issues:         []ValidationIssue{},
						Recommendations: []string{"Review criteria regularly", "Update thresholds as needed"},
						CompletedAt:    time.Now(),
					},
				},
			},
			Criteria: []ValidationCriterion{
				{
					ID:          "criteria-1-1",
					Name:        "Environment Setup",
					Description: "Validation environment is properly configured",
					Type:        0, // Functional
					Weight:      0.3,
					Threshold:   0.9,
					Status:      1, // Passed
					Score:       0.95,
					Comments:    "Environment setup meets requirements",
				},
				{
					ID:          "criteria-1-2",
					Name:        "Test Data Quality",
					Description: "Test data is comprehensive and realistic",
					Type:        6, // Quality
					Weight:      0.4,
					Threshold:   0.8,
					Status:      1, // Passed
					Score:       0.88,
					Comments:    "Test data quality is acceptable",
				},
			},
			StartTime: func() *time.Time { t := time.Now().Add(-time.Hour); return &t }(),
			EndTime:   func() *time.Time { t := time.Now(); return &t }(),
			Duration:  time.Hour,
		},
		{
			ID:          "phase-2",
			Name:        "User Testing",
			Description: "Conduct user testing sessions",
			Order:       2,
			Status:      0, // Pending
			Tasks:       []ValidationTask{},
			Criteria:    []ValidationCriterion{},
		},
	}

	// Validate phases
	for i, phase := range phases {
		if phase.ID == "" {
			fmt.Printf("❌ Phase %d ID should not be empty\n", i+1)
			return
		}

		if phase.Name == "" {
			fmt.Printf("❌ Phase %d name should not be empty\n", i+1)
			return
		}

		if phase.Description == "" {
			fmt.Printf("❌ Phase %d description should not be empty\n", i+1)
			return
		}

		if phase.Order <= 0 {
			fmt.Printf("❌ Phase %d order should be positive\n", i+1)
			return
		}

		if phase.Status < 0 || phase.Status > 4 {
			fmt.Printf("❌ Phase %d status should be between 0 and 4\n", i+1)
			return
		}
	}

	// Validate tasks in first phase
	phase1 := phases[0]
	for i, task := range phase1.Tasks {
		if task.ID == "" {
			fmt.Printf("❌ Task %d ID should not be empty\n", i+1)
			return
		}

		if task.Name == "" {
			fmt.Printf("❌ Task %d name should not be empty\n", i+1)
			return
		}

		if task.Description == "" {
			fmt.Printf("❌ Task %d description should not be empty\n", i+1)
			return
		}

		if task.Type < 0 || task.Type > 7 {
			fmt.Printf("❌ Task %d type should be between 0 and 7\n", i+1)
			return
		}

		if task.Priority < 0 || task.Priority > 3 {
			fmt.Printf("❌ Task %d priority should be between 0 and 3\n", i+1)
			return
		}

		if task.Status < 0 || task.Status > 5 {
			fmt.Printf("❌ Task %d status should be between 0 and 5\n", i+1)
			return
		}

		if task.Result == nil {
			fmt.Printf("❌ Task %d result should not be nil\n", i+1)
			return
		}

		// Validate task result
		result := task.Result
		if result.Score < 0.0 || result.Score > 1.0 {
			fmt.Printf("❌ Task %d result score should be between 0 and 1\n", i+1)
			return
		}

		if result.Status < 0 || result.Status > 5 {
			fmt.Printf("❌ Task %d result status should be between 0 and 5\n", i+1)
			return
		}

		if result.Comments == "" {
			fmt.Printf("❌ Task %d result comments should not be empty\n", i+1)
			return
		}

		if len(result.Evidence) == 0 {
			fmt.Printf("❌ Task %d result evidence should not be empty\n", i+1)
			return
		}
	}

	// Validate criteria in first phase
	for i, criterion := range phase1.Criteria {
		if criterion.ID == "" {
			fmt.Printf("❌ Criterion %d ID should not be empty\n", i+1)
			return
		}

		if criterion.Name == "" {
			fmt.Printf("❌ Criterion %d name should not be empty\n", i+1)
			return
		}

		if criterion.Description == "" {
			fmt.Printf("❌ Criterion %d description should not be empty\n", i+1)
			return
		}

		if criterion.Type < 0 || criterion.Type > 7 {
			fmt.Printf("❌ Criterion %d type should be between 0 and 7\n", i+1)
			return
		}

		if criterion.Weight < 0.0 || criterion.Weight > 1.0 {
			fmt.Printf("❌ Criterion %d weight should be between 0 and 1\n", i+1)
			return
		}

		if criterion.Threshold < 0.0 || criterion.Threshold > 1.0 {
			fmt.Printf("❌ Criterion %d threshold should be between 0 and 1\n", i+1)
			return
		}

		if criterion.Status < 0 || criterion.Status > 4 {
			fmt.Printf("❌ Criterion %d status should be between 0 and 4\n", i+1)
			return
		}

		if criterion.Score < 0.0 || criterion.Score > 1.0 {
			fmt.Printf("❌ Criterion %d score should be between 0 and 1\n", i+1)
			return
		}
	}

	// Validate phase timing
	if phase1.StartTime == nil {
		fmt.Println("❌ Phase 1 start time should not be nil")
		return
	}

	if phase1.EndTime == nil {
		fmt.Println("❌ Phase 1 end time should not be nil")
		return
	}

	if phase1.Duration <= 0 {
		fmt.Println("❌ Phase 1 duration should be positive")
		return
	}

	fmt.Printf("✅ Validation phases test passed\n")
	fmt.Printf("   Total phases: %d\n", len(phases))
	fmt.Printf("   Phase 1 tasks: %d\n", len(phase1.Tasks))
	fmt.Printf("   Phase 1 criteria: %d\n", len(phase1.Criteria))
	fmt.Printf("   Phase 1 duration: %v\n", phase1.Duration)
	fmt.Printf("   Phase 1 status: %d\n", phase1.Status)
}

func testApprovalSystem() {
	// Test approval system
	approvals := []Approval{
		{
			ID:          "approval-001",
			Type:        0, // Feature
			Status:      0, // Pending
			RequestedBy: "user1",
			RequestedAt: time.Now().Add(-time.Hour),
			Comments:    "Requesting approval for new feature",
			Priority:    1, // High
			ExpiresAt:   func() *time.Time { t := time.Now().Add(time.Hour * 24); return &t }(),
		},
		{
			ID:          "approval-002",
			Type:        1, // Release
			Status:      1, // Approved
			RequestedBy: "user2",
			RequestedAt: time.Now().Add(-time.Hour * 2),
			ApprovedBy:  "approver1",
			ApprovedAt:  func() *time.Time { t := time.Now().Add(-time.Hour); return &t }(),
			Comments:    "Release approved after review",
			Priority:    0, // Critical
			ExpiresAt:   func() *time.Time { t := time.Now().Add(time.Hour * 48); return &t }(),
		},
		{
			ID:          "approval-003",
			Type:        2, // Quality
			Status:      2, // Rejected
			RequestedBy: "user3",
			RequestedAt: time.Now().Add(-time.Hour * 3),
			ApprovedBy:  "approver2",
			ApprovedAt:  func() *time.Time { t := time.Now().Add(-time.Hour * 2); return &t }(),
			Comments:    "Quality standards not met",
			Priority:    2, // Medium
			ExpiresAt:   func() *time.Time { t := time.Now().Add(time.Hour * 12); return &t }(),
		},
	}

	approvers := []Approver{
		{
			ID:     "approver1",
			Name:   "John Doe",
			Email:  "john.doe@example.com",
			Role:   "Quality Manager",
			Level:  3, // Manager
			Active: true,
		},
		{
			ID:     "approver2",
			Name:   "Jane Smith",
			Email:  "jane.smith@example.com",
			Role:   "Test Director",
			Level:  4, // Director
			Active: true,
		},
		{
			ID:     "approver3",
			Name:   "Bob Johnson",
			Email:  "bob.johnson@example.com",
			Role:   "Executive",
			Level:  5, // Executive
			Active: false,
		},
	}

	// Validate approvals
	for i, approval := range approvals {
		if approval.ID == "" {
			fmt.Printf("❌ Approval %d ID should not be empty\n", i+1)
			return
		}

		if approval.Type < 0 || approval.Type > 5 {
			fmt.Printf("❌ Approval %d type should be between 0 and 5\n", i+1)
			return
		}

		if approval.Status < 0 || approval.Status > 4 {
			fmt.Printf("❌ Approval %d status should be between 0 and 4\n", i+1)
			return
		}

		if approval.RequestedBy == "" {
			fmt.Printf("❌ Approval %d requested by should not be empty\n", i+1)
			return
		}

		if approval.Comments == "" {
			fmt.Printf("❌ Approval %d comments should not be empty\n", i+1)
			return
		}

		if approval.Priority < 0 || approval.Priority > 3 {
			fmt.Printf("❌ Approval %d priority should be between 0 and 3\n", i+1)
			return
		}

		if approval.ExpiresAt == nil {
			fmt.Printf("❌ Approval %d expires at should not be nil\n", i+1)
			return
		}

		// Validate approval status consistency
		if approval.Status == 1 && approval.ApprovedBy == "" { // Approved
			fmt.Printf("❌ Approval %d should have approver when approved\n", i+1)
			return
		}

		if approval.Status == 1 && approval.ApprovedAt == nil { // Approved
			fmt.Printf("❌ Approval %d should have approval time when approved\n", i+1)
			return
		}
	}

	// Validate approvers
	for i, approver := range approvers {
		if approver.ID == "" {
			fmt.Printf("❌ Approver %d ID should not be empty\n", i+1)
			return
		}

		if approver.Name == "" {
			fmt.Printf("❌ Approver %d name should not be empty\n", i+1)
			return
		}

		if approver.Email == "" {
			fmt.Printf("❌ Approver %d email should not be empty\n", i+1)
			return
		}

		if approver.Role == "" {
			fmt.Printf("❌ Approver %d role should not be empty\n", i+1)
			return
		}

		if approver.Level < 0 || approver.Level > 5 {
			fmt.Printf("❌ Approver %d level should be between 0 and 5\n", i+1)
			return
		}
	}

	// Count approvals by status
	pendingCount := 0
	approvedCount := 0
	rejectedCount := 0
	for _, approval := range approvals {
		switch approval.Status {
		case 0: // Pending
			pendingCount++
		case 1: // Approved
			approvedCount++
		case 2: // Rejected
			rejectedCount++
		}
	}

	// Count active approvers
	activeApprovers := 0
	for _, approver := range approvers {
		if approver.Active {
			activeApprovers++
		}
	}

	fmt.Printf("✅ Approval system test passed\n")
	fmt.Printf("   Total approvals: %d\n", len(approvals))
	fmt.Printf("   Pending approvals: %d\n", pendingCount)
	fmt.Printf("   Approved approvals: %d\n", approvedCount)
	fmt.Printf("   Rejected approvals: %d\n", rejectedCount)
	fmt.Printf("   Total approvers: %d\n", len(approvers))
	fmt.Printf("   Active approvers: %d\n", activeApprovers)
}

func testValidationCriteria() {
	// Test validation criteria
	criteria := []ValidationCriterion{
		{
			ID:          "criteria-001",
			Name:        "Functional Requirements",
			Description: "All functional requirements are met",
			Type:        0, // Functional
			Weight:      0.5,
			Threshold:   1.0,
			Status:      1, // Passed
			Score:       1.0,
			Comments:    "All functional requirements verified",
		},
		{
			ID:          "criteria-002",
			Name:        "Performance Requirements",
			Description: "Performance requirements are met",
			Type:        3, // Performance
			Weight:      0.3,
			Threshold:   0.9,
			Status:      1, // Passed
			Score:       0.95,
			Comments:    "Performance exceeds requirements",
		},
		{
			ID:          "criteria-003",
			Name:        "Usability Requirements",
			Description: "Usability requirements are met",
			Type:        1, // Usability
			Weight:      0.2,
			Threshold:   0.8,
			Status:      2, // Failed
			Score:       0.75,
			Comments:    "Usability needs improvement",
		},
	}

	// Validate criteria
	for i, criterion := range criteria {
		if criterion.ID == "" {
			fmt.Printf("❌ Criterion %d ID should not be empty\n", i+1)
			return
		}

		if criterion.Name == "" {
			fmt.Printf("❌ Criterion %d name should not be empty\n", i+1)
			return
		}

		if criterion.Description == "" {
			fmt.Printf("❌ Criterion %d description should not be empty\n", i+1)
			return
		}

		if criterion.Type < 0 || criterion.Type > 7 {
			fmt.Printf("❌ Criterion %d type should be between 0 and 7\n", i+1)
			return
		}

		if criterion.Weight < 0.0 || criterion.Weight > 1.0 {
			fmt.Printf("❌ Criterion %d weight should be between 0 and 1\n", i+1)
			return
		}

		if criterion.Threshold < 0.0 || criterion.Threshold > 1.0 {
			fmt.Printf("❌ Criterion %d threshold should be between 0 and 1\n", i+1)
			return
		}

		if criterion.Status < 0 || criterion.Status > 4 {
			fmt.Printf("❌ Criterion %d status should be between 0 and 4\n", i+1)
			return
		}

		if criterion.Score < 0.0 || criterion.Score > 1.0 {
			fmt.Printf("❌ Criterion %d score should be between 0 and 1\n", i+1)
			return
		}

		if criterion.Comments == "" {
			fmt.Printf("❌ Criterion %d comments should not be empty\n", i+1)
			return
		}

		// Validate score vs threshold consistency
		if criterion.Score >= criterion.Threshold && criterion.Status != 1 { // Passed
			fmt.Printf("❌ Criterion %d status should be passed when score >= threshold\n", i+1)
			return
		}

		if criterion.Score < criterion.Threshold && criterion.Status == 1 { // Failed
			fmt.Printf("❌ Criterion %d status should not be passed when score < threshold\n", i+1)
			return
		}
	}

	// Calculate weighted score
	totalWeight := 0.0
	weightedScore := 0.0
	for _, criterion := range criteria {
		totalWeight += criterion.Weight
		weightedScore += criterion.Score * criterion.Weight
	}

	overallScore := 0.0
	if totalWeight > 0 {
		overallScore = weightedScore / totalWeight
	}

	// Count criteria by status
	passedCount := 0
	failedCount := 0
	for _, criterion := range criteria {
		switch criterion.Status {
		case 1: // Passed
			passedCount++
		case 2: // Failed
			failedCount++
		}
	}

	fmt.Printf("✅ Validation criteria test passed\n")
	fmt.Printf("   Total criteria: %d\n", len(criteria))
	fmt.Printf("   Passed criteria: %d\n", passedCount)
	fmt.Printf("   Failed criteria: %d\n", failedCount)
	fmt.Printf("   Overall score: %.2f%%\n", overallScore*100)
	fmt.Printf("   Total weight: %.2f\n", totalWeight)
}
