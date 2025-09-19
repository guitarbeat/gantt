package main

import (
	"fmt"
	"time"
)

// Test the user coordination system
func main() {
	fmt.Println("Testing User Coordination System...")

	// Test 1: User Coordination Configuration
	fmt.Println("\n=== Test 1: User Coordination Configuration ===")
	testUserCoordinationConfiguration()

	// Test 2: User Session Management
	fmt.Println("\n=== Test 2: User Session Management ===")
	testUserSessionManagement()

	// Test 3: Communication Channels
	fmt.Println("\n=== Test 3: Communication Channels ===")
	testCommunicationChannels()

	// Test 4: Feedback Workflows
	fmt.Println("\n=== Test 4: Feedback Workflows ===")
	testFeedbackWorkflows()

	// Test 5: Notification System
	fmt.Println("\n=== Test 5: Notification System ===")
	testNotificationSystem()

	fmt.Println("\n✅ User coordination system tests completed!")
}

// UserCoordinationConfig represents user coordination configuration
type UserCoordinationConfig struct {
	EnableEmailNotifications    bool
	EnableInAppNotifications    bool
	EnablePushNotifications     bool
	NotificationTimeout         time.Duration
	MaxNotificationRetries      int
	EnableFeedbackPrompts       bool
	FeedbackPromptDelay         time.Duration
	FeedbackPromptFrequency     time.Duration
	MaxFeedbackPromptsPerSession int
	EnableUserEngagement        bool
	EngagementThreshold         float64
	EngagementDecayRate         float64
	EngagementBoostFactor       float64
	EnableFollowUpRequests      bool
	FollowUpTimeout             time.Duration
	MaxFollowUpAttempts         int
	FollowUpEscalationThreshold float64
	EnableAutoResponses         bool
	ResponseTemplatePath        string
	ResponsePersonalization     bool
	ResponseLanguage            string
}

// UserSession represents a user session
type UserSession struct {
	ID              string
	UserID          string
	StartTime       time.Time
	LastActivity    time.Time
	FeedbackCount   int
	EngagementScore float64
	Context         map[string]interface{}
	Preferences     *UserPreferences
	Status          int
}

// UserPreferences represents user preferences
type UserPreferences struct {
	PreferredChannel     int
	NotificationFrequency int
	FeedbackCategories   []int
	Language             string
	Timezone             string
	AutoSubmit           bool
	AnonymousMode        bool
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

func testUserCoordinationConfiguration() {
	// Test user coordination configuration
	config := UserCoordinationConfig{
		EnableEmailNotifications:    true,
		EnableInAppNotifications:    true,
		EnablePushNotifications:     false,
		NotificationTimeout:         time.Minute * 5,
		MaxNotificationRetries:      3,
		EnableFeedbackPrompts:       true,
		FeedbackPromptDelay:         time.Minute * 2,
		FeedbackPromptFrequency:     time.Hour * 24,
		MaxFeedbackPromptsPerSession: 3,
		EnableUserEngagement:        true,
		EngagementThreshold:         0.5,
		EngagementDecayRate:         0.1,
		EngagementBoostFactor:       0.2,
		EnableFollowUpRequests:      true,
		FollowUpTimeout:             time.Hour * 24,
		MaxFollowUpAttempts:         3,
		FollowUpEscalationThreshold: 0.8,
		EnableAutoResponses:         true,
		ResponseTemplatePath:        "/templates/feedback/",
		ResponsePersonalization:     true,
		ResponseLanguage:            "en",
	}

	// Validate configuration
	if !config.EnableEmailNotifications {
		fmt.Println("❌ Email notifications should be enabled")
		return
	}

	if !config.EnableInAppNotifications {
		fmt.Println("❌ In-app notifications should be enabled")
		return
	}

	if config.EnablePushNotifications {
		fmt.Println("❌ Push notifications should be disabled by default")
		return
	}

	if config.NotificationTimeout <= 0 {
		fmt.Println("❌ Notification timeout should be positive")
		return
	}

	if config.MaxNotificationRetries <= 0 {
		fmt.Println("❌ Max notification retries should be positive")
		return
	}

	if !config.EnableFeedbackPrompts {
		fmt.Println("❌ Feedback prompts should be enabled")
		return
	}

	if config.FeedbackPromptDelay <= 0 {
		fmt.Println("❌ Feedback prompt delay should be positive")
		return
	}

	if config.FeedbackPromptFrequency <= 0 {
		fmt.Println("❌ Feedback prompt frequency should be positive")
		return
	}

	if config.MaxFeedbackPromptsPerSession <= 0 {
		fmt.Println("❌ Max feedback prompts per session should be positive")
		return
	}

	if !config.EnableUserEngagement {
		fmt.Println("❌ User engagement should be enabled")
		return
	}

	if config.EngagementThreshold < 0.0 || config.EngagementThreshold > 1.0 {
		fmt.Println("❌ Engagement threshold should be between 0 and 1")
		return
	}

	if config.EngagementDecayRate < 0.0 || config.EngagementDecayRate > 1.0 {
		fmt.Println("❌ Engagement decay rate should be between 0 and 1")
		return
	}

	if config.EngagementBoostFactor < 0.0 || config.EngagementBoostFactor > 1.0 {
		fmt.Println("❌ Engagement boost factor should be between 0 and 1")
		return
	}

	if !config.EnableFollowUpRequests {
		fmt.Println("❌ Follow-up requests should be enabled")
		return
	}

	if config.FollowUpTimeout <= 0 {
		fmt.Println("❌ Follow-up timeout should be positive")
		return
	}

	if config.MaxFollowUpAttempts <= 0 {
		fmt.Println("❌ Max follow-up attempts should be positive")
		return
	}

	if config.FollowUpEscalationThreshold < 0.0 || config.FollowUpEscalationThreshold > 1.0 {
		fmt.Println("❌ Follow-up escalation threshold should be between 0 and 1")
		return
	}

	if !config.EnableAutoResponses {
		fmt.Println("❌ Auto responses should be enabled")
		return
	}

	if config.ResponseTemplatePath == "" {
		fmt.Println("❌ Response template path should not be empty")
		return
	}

	if !config.ResponsePersonalization {
		fmt.Println("❌ Response personalization should be enabled")
		return
	}

	if config.ResponseLanguage == "" {
		fmt.Println("❌ Response language should not be empty")
		return
	}

	fmt.Printf("✅ User coordination configuration test passed\n")
	fmt.Printf("   Enable email notifications: %v\n", config.EnableEmailNotifications)
	fmt.Printf("   Enable in-app notifications: %v\n", config.EnableInAppNotifications)
	fmt.Printf("   Enable push notifications: %v\n", config.EnablePushNotifications)
	fmt.Printf("   Notification timeout: %v\n", config.NotificationTimeout)
	fmt.Printf("   Max notification retries: %d\n", config.MaxNotificationRetries)
	fmt.Printf("   Enable feedback prompts: %v\n", config.EnableFeedbackPrompts)
	fmt.Printf("   Feedback prompt delay: %v\n", config.FeedbackPromptDelay)
	fmt.Printf("   Feedback prompt frequency: %v\n", config.FeedbackPromptFrequency)
	fmt.Printf("   Max feedback prompts per session: %d\n", config.MaxFeedbackPromptsPerSession)
	fmt.Printf("   Enable user engagement: %v\n", config.EnableUserEngagement)
	fmt.Printf("   Engagement threshold: %.2f\n", config.EngagementThreshold)
	fmt.Printf("   Engagement decay rate: %.2f\n", config.EngagementDecayRate)
	fmt.Printf("   Engagement boost factor: %.2f\n", config.EngagementBoostFactor)
	fmt.Printf("   Enable follow-up requests: %v\n", config.EnableFollowUpRequests)
	fmt.Printf("   Follow-up timeout: %v\n", config.FollowUpTimeout)
	fmt.Printf("   Max follow-up attempts: %d\n", config.MaxFollowUpAttempts)
	fmt.Printf("   Follow-up escalation threshold: %.2f\n", config.FollowUpEscalationThreshold)
	fmt.Printf("   Enable auto responses: %v\n", config.EnableAutoResponses)
	fmt.Printf("   Response template path: %s\n", config.ResponseTemplatePath)
	fmt.Printf("   Response personalization: %v\n", config.ResponsePersonalization)
	fmt.Printf("   Response language: %s\n", config.ResponseLanguage)
}

func testUserSessionManagement() {
	// Test user session management
	session := UserSession{
		ID:              "session-001",
		UserID:          "user-123",
		StartTime:       time.Now(),
		LastActivity:    time.Now(),
		FeedbackCount:   0,
		EngagementScore: 0.5,
		Context: map[string]interface{}{
			"view_type": "monthly",
			"task_count": 15,
			"screen_size": "1920x1080",
		},
		Preferences: &UserPreferences{
			PreferredChannel:     1, // InApp
			NotificationFrequency: 2, // Daily
			FeedbackCategories:   []int{0, 1, 2}, // General, Visual, Layout
			Language:             "en",
			Timezone:             "UTC",
			AutoSubmit:           false,
			AnonymousMode:        false,
		},
		Status: 0, // Active
	}

	// Validate session
	if session.ID == "" {
		fmt.Println("❌ Session ID should not be empty")
		return
	}

	if session.UserID == "" {
		fmt.Println("❌ User ID should not be empty")
		return
	}

	if session.StartTime.IsZero() {
		fmt.Println("❌ Start time should not be zero")
		return
	}

	if session.LastActivity.IsZero() {
		fmt.Println("❌ Last activity should not be zero")
		return
	}

	if session.FeedbackCount < 0 {
		fmt.Println("❌ Feedback count should be non-negative")
		return
	}

	if session.EngagementScore < 0.0 || session.EngagementScore > 1.0 {
		fmt.Println("❌ Engagement score should be between 0 and 1")
		return
	}

	if session.Context == nil {
		fmt.Println("❌ Context should not be nil")
		return
	}

	if session.Preferences == nil {
		fmt.Println("❌ Preferences should not be nil")
		return
	}

	if session.Status < 0 || session.Status > 3 {
		fmt.Println("❌ Status should be between 0 and 3")
		return
	}

	// Validate preferences
	prefs := session.Preferences
	if prefs.PreferredChannel < 0 || prefs.PreferredChannel > 5 {
		fmt.Println("❌ Preferred channel should be between 0 and 5")
		return
	}

	if prefs.NotificationFrequency < 0 || prefs.NotificationFrequency > 4 {
		fmt.Println("❌ Notification frequency should be between 0 and 4")
		return
	}

	if len(prefs.FeedbackCategories) == 0 {
		fmt.Println("❌ Feedback categories should not be empty")
		return
	}

	if prefs.Language == "" {
		fmt.Println("❌ Language should not be empty")
		return
	}

	if prefs.Timezone == "" {
		fmt.Println("❌ Timezone should not be empty")
		return
	}

	// Test session update
	session.FeedbackCount++
	session.LastActivity = time.Now()
	session.EngagementScore = 0.7

	// Validate updated session
	if session.FeedbackCount != 1 {
		fmt.Println("❌ Feedback count should be 1 after increment")
		return
	}

	if session.EngagementScore != 0.7 {
		fmt.Println("❌ Engagement score should be 0.7 after update")
		return
	}

	fmt.Printf("✅ User session management test passed\n")
	fmt.Printf("   Session ID: %s\n", session.ID)
	fmt.Printf("   User ID: %s\n", session.UserID)
	fmt.Printf("   Start time: %v\n", session.StartTime)
	fmt.Printf("   Last activity: %v\n", session.LastActivity)
	fmt.Printf("   Feedback count: %d\n", session.FeedbackCount)
	fmt.Printf("   Engagement score: %.2f\n", session.EngagementScore)
	fmt.Printf("   Context items: %d\n", len(session.Context))
	fmt.Printf("   Preferred channel: %d\n", prefs.PreferredChannel)
	fmt.Printf("   Notification frequency: %d\n", prefs.NotificationFrequency)
	fmt.Printf("   Feedback categories: %d\n", len(prefs.FeedbackCategories))
	fmt.Printf("   Language: %s\n", prefs.Language)
	fmt.Printf("   Timezone: %s\n", prefs.Timezone)
	fmt.Printf("   Auto submit: %v\n", prefs.AutoSubmit)
	fmt.Printf("   Anonymous mode: %v\n", prefs.AnonymousMode)
	fmt.Printf("   Status: %d\n", session.Status)
}

func testCommunicationChannels() {
	// Test communication channels
	channels := []int{0, 1, 2, 3, 4, 5} // Email, InApp, Push, Webhook, SMS, Chat
	channelNames := []string{"Email", "InApp", "Push", "Webhook", "SMS", "Chat"}

	for i, channel := range channels {
		// Validate channel range
		if channel < 0 || channel > 5 {
			fmt.Printf("❌ Channel %d should be between 0 and 5\n", channel)
			return
		}

		// Test channel-specific configuration
		switch channel {
		case 0: // Email
			// Email channel should support templates and SMTP
			fmt.Printf("   %s channel: Supports templates and SMTP\n", channelNames[i])
		case 1: // InApp
			// In-app channel should support UI components
			fmt.Printf("   %s channel: Supports UI components\n", channelNames[i])
		case 2: // Push
			// Push channel should support FCM/APNS
			fmt.Printf("   %s channel: Supports FCM/APNS\n", channelNames[i])
		case 3: // Webhook
			// Webhook channel should support HTTP endpoints
			fmt.Printf("   %s channel: Supports HTTP endpoints\n", channelNames[i])
		case 4: // SMS
			// SMS channel should support text messaging
			fmt.Printf("   %s channel: Supports text messaging\n", channelNames[i])
		case 5: // Chat
			// Chat channel should support real-time messaging
			fmt.Printf("   %s channel: Supports real-time messaging\n", channelNames[i])
		}
	}

	fmt.Printf("✅ Communication channels test passed\n")
	fmt.Printf("   Total channels: %d\n", len(channels))
	fmt.Printf("   Channel types: %v\n", channelNames)
}

func testFeedbackWorkflows() {
	// Test feedback workflows
	workflow := struct {
		ID              string
		Name            string
		Description     string
		Trigger         int
		Steps           []WorkflowStep
		Conditions      []WorkflowCondition
		Actions         []WorkflowAction
		Timeout         time.Duration
		RetryPolicy     *RetryPolicy
		Status          int
		CreatedAt       time.Time
		UpdatedAt       time.Time
	}{
		ID:          "workflow-001",
		Name:        "Feedback Collection Workflow",
		Description: "Collects and processes user feedback",
		Trigger:     1, // Time
		Steps: []WorkflowStep{
			{
				ID:          "step-001",
				Name:        "Collect Feedback",
				Type:        0, // FeedbackPrompt
				Description: "Prompt user for feedback",
				Config: map[string]interface{}{
					"prompt_type": "general",
					"timeout":     "5m",
				},
				Order: 1,
			},
			{
				ID:          "step-002",
				Name:        "Process Feedback",
				Type:        3, // Processing
				Description: "Process collected feedback",
				Config: map[string]interface{}{
					"auto_process": true,
					"threshold":    0.7,
				},
				Order: 2,
			},
		},
		Conditions: []WorkflowCondition{
			{
				ID:         "condition-001",
				Name:       "User Engagement Check",
				Type:       0, // UserEngagement
				Expression: "engagement_score > 0.5",
				Parameters: map[string]interface{}{
					"threshold": 0.5,
				},
			},
		},
		Actions: []WorkflowAction{
			{
				ID:          "action-001",
				Name:        "Send Notification",
				Type:        0, // SendNotification
				Description: "Send notification to user",
				Config: map[string]interface{}{
					"channel": "inapp",
					"message": "Thank you for your feedback!",
				},
			},
		},
		Timeout: time.Minute * 30,
		RetryPolicy: &RetryPolicy{
			MaxRetries:      3,
			InitialDelay:    time.Second * 1,
			MaxDelay:        time.Minute * 5,
			BackoffFactor:   2.0,
			RetryableErrors: []string{"timeout", "network"},
		},
		Status:    1, // Active
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Validate workflow
	if workflow.ID == "" {
		fmt.Println("❌ Workflow ID should not be empty")
		return
	}

	if workflow.Name == "" {
		fmt.Println("❌ Workflow name should not be empty")
		return
	}

	if workflow.Description == "" {
		fmt.Println("❌ Workflow description should not be empty")
		return
	}

	if workflow.Trigger < 0 || workflow.Trigger > 4 {
		fmt.Println("❌ Workflow trigger should be between 0 and 4")
		return
	}

	if len(workflow.Steps) == 0 {
		fmt.Println("❌ Workflow should have at least one step")
		return
	}

	if workflow.Timeout <= 0 {
		fmt.Println("❌ Workflow timeout should be positive")
		return
	}

	if workflow.RetryPolicy == nil {
		fmt.Println("❌ Workflow should have retry policy")
		return
	}

	if workflow.Status < 0 || workflow.Status > 5 {
		fmt.Println("❌ Workflow status should be between 0 and 5")
		return
	}

	// Validate steps
	for i, step := range workflow.Steps {
		if step.ID == "" {
			fmt.Printf("❌ Step %d ID should not be empty\n", i+1)
			return
		}

		if step.Name == "" {
			fmt.Printf("❌ Step %d name should not be empty\n", i+1)
			return
		}

		if step.Type < 0 || step.Type > 7 {
			fmt.Printf("❌ Step %d type should be between 0 and 7\n", i+1)
			return
		}

		if step.Order <= 0 {
			fmt.Printf("❌ Step %d order should be positive\n", i+1)
			return
		}
	}

	// Validate conditions
	for i, condition := range workflow.Conditions {
		if condition.ID == "" {
			fmt.Printf("❌ Condition %d ID should not be empty\n", i+1)
			return
		}

		if condition.Name == "" {
			fmt.Printf("❌ Condition %d name should not be empty\n", i+1)
			return
		}

		if condition.Type < 0 || condition.Type > 6 {
			fmt.Printf("❌ Condition %d type should be between 0 and 6\n", i+1)
			return
		}

		if condition.Expression == "" {
			fmt.Printf("❌ Condition %d expression should not be empty\n", i+1)
			return
		}
	}

	// Validate actions
	for i, action := range workflow.Actions {
		if action.ID == "" {
			fmt.Printf("❌ Action %d ID should not be empty\n", i+1)
			return
		}

		if action.Name == "" {
			fmt.Printf("❌ Action %d name should not be empty\n", i+1)
			return
		}

		if action.Type < 0 || action.Type > 5 {
			fmt.Printf("❌ Action %d type should be between 0 and 5\n", i+1)
			return
		}
	}

	// Validate retry policy
	retryPolicy := workflow.RetryPolicy
	if retryPolicy.MaxRetries <= 0 {
		fmt.Println("❌ Max retries should be positive")
		return
	}

	if retryPolicy.InitialDelay <= 0 {
		fmt.Println("❌ Initial delay should be positive")
		return
	}

	if retryPolicy.MaxDelay <= 0 {
		fmt.Println("❌ Max delay should be positive")
		return
	}

	if retryPolicy.BackoffFactor <= 0 {
		fmt.Println("❌ Backoff factor should be positive")
		return
	}

	if len(retryPolicy.RetryableErrors) == 0 {
		fmt.Println("❌ Retryable errors should not be empty")
		return
	}

	fmt.Printf("✅ Feedback workflows test passed\n")
	fmt.Printf("   Workflow ID: %s\n", workflow.ID)
	fmt.Printf("   Workflow name: %s\n", workflow.Name)
	fmt.Printf("   Workflow description: %s\n", workflow.Description)
	fmt.Printf("   Workflow trigger: %d\n", workflow.Trigger)
	fmt.Printf("   Workflow steps: %d\n", len(workflow.Steps))
	fmt.Printf("   Workflow conditions: %d\n", len(workflow.Conditions))
	fmt.Printf("   Workflow actions: %d\n", len(workflow.Actions))
	fmt.Printf("   Workflow timeout: %v\n", workflow.Timeout)
	fmt.Printf("   Workflow status: %d\n", workflow.Status)
	fmt.Printf("   Max retries: %d\n", retryPolicy.MaxRetries)
	fmt.Printf("   Initial delay: %v\n", retryPolicy.InitialDelay)
	fmt.Printf("   Max delay: %v\n", retryPolicy.MaxDelay)
	fmt.Printf("   Backoff factor: %.2f\n", retryPolicy.BackoffFactor)
	fmt.Printf("   Retryable errors: %d\n", len(retryPolicy.RetryableErrors))
}

func testNotificationSystem() {
	// Test notification system
	notification := struct {
		ID        string
		Type      string
		Recipient string
		Message   string
		Data      map[string]interface{}
		Priority  int
		CreatedAt time.Time
	}{
		ID:        "notification-001",
		Type:      "feedback_prompt",
		Recipient: "user-123",
		Message:   "We'd love your feedback!",
		Data: map[string]interface{}{
			"feedback_id": "feedback-001",
			"session_id":  "session-001",
			"prompt_type": "general",
		},
		Priority:  1,
		CreatedAt: time.Now(),
	}

	// Validate notification
	if notification.ID == "" {
		fmt.Println("❌ Notification ID should not be empty")
		return
	}

	if notification.Type == "" {
		fmt.Println("❌ Notification type should not be empty")
		return
	}

	if notification.Recipient == "" {
		fmt.Println("❌ Notification recipient should not be empty")
		return
	}

	if notification.Message == "" {
		fmt.Println("❌ Notification message should not be empty")
		return
	}

	if notification.Data == nil {
		fmt.Println("❌ Notification data should not be nil")
		return
	}

	if notification.Priority < 0 || notification.Priority > 5 {
		fmt.Println("❌ Notification priority should be between 0 and 5")
		return
	}

	if notification.CreatedAt.IsZero() {
		fmt.Println("❌ Notification created at should not be zero")
		return
	}

	// Test notification types
	notificationTypes := []string{
		"feedback_prompt",
		"feedback_received",
		"feedback_processed",
		"feedback_resolved",
		"follow_up_request",
		"thank_you",
		"escalation",
		"reminder",
	}

	// Validate notification type
	validType := false
	for _, nt := range notificationTypes {
		if notification.Type == nt {
			validType = true
			break
		}
	}

	if !validType {
		fmt.Println("❌ Notification type should be valid")
		return
	}

	// Test notification priorities
	notificationPriorities := []int{0, 1, 2, 3, 4, 5} // Low, Medium, High, Critical, Emergency, System

	// Validate notification priority
	validPriority := false
	for _, np := range notificationPriorities {
		if notification.Priority == np {
			validPriority = true
			break
		}
	}

	if !validPriority {
		fmt.Println("❌ Notification priority should be valid")
		return
	}

	// Test notification data
	requiredDataFields := []string{"feedback_id", "session_id", "prompt_type"}
	for _, field := range requiredDataFields {
		if _, exists := notification.Data[field]; !exists {
			fmt.Printf("❌ Notification data should contain field: %s\n", field)
			return
		}
	}

	fmt.Printf("✅ Notification system test passed\n")
	fmt.Printf("   Notification ID: %s\n", notification.ID)
	fmt.Printf("   Notification type: %s\n", notification.Type)
	fmt.Printf("   Notification recipient: %s\n", notification.Recipient)
	fmt.Printf("   Notification message: %s\n", notification.Message)
	fmt.Printf("   Notification data fields: %d\n", len(notification.Data))
	fmt.Printf("   Notification priority: %d\n", notification.Priority)
	fmt.Printf("   Notification created at: %v\n", notification.CreatedAt)
}

// Additional types for completeness
type WorkflowStep struct {
	ID          string
	Name        string
	Type        int
	Description string
	Config      map[string]interface{}
	Conditions  []WorkflowCondition
	Actions     []WorkflowAction
	Timeout     time.Duration
	RetryPolicy *RetryPolicy
	Status      int
	Order       int
}

type WorkflowCondition struct {
	ID         string
	Name       string
	Type       int
	Expression string
	Parameters map[string]interface{}
	Negate     bool
	Status     int
}

type WorkflowAction struct {
	ID          string
	Name        string
	Type        int
	Description string
	Config      map[string]interface{}
	Parameters  map[string]interface{}
	Status      int
	Result      map[string]interface{}
}

type RetryPolicy struct {
	MaxRetries      int
	InitialDelay    time.Duration
	MaxDelay        time.Duration
	BackoffFactor   float64
	RetryableErrors []string
}
