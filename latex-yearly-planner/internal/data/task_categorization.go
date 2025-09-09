package data

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// TaskCategoryManager manages task categories and their relationships
type TaskCategoryManager struct {
	categories map[string]TaskCategory
	rules      []CategoryRule
}

// CategoryRule defines rules for automatic task categorization
type CategoryRule struct {
	Name        string
	Patterns    []string
	Category    string
	Priority    int
	Description string
}

// NewTaskCategoryManager creates a new category manager with predefined rules
func NewTaskCategoryManager() *TaskCategoryManager {
	tcm := &TaskCategoryManager{
		categories: make(map[string]TaskCategory),
		rules:      make([]CategoryRule, 0),
	}
	
	// Initialize with predefined categories
	for _, category := range GetAllCategories() {
		tcm.categories[category.Name] = category
	}
	
	// Add categorization rules
	tcm.addCategorizationRules()
	
	return tcm
}

// addCategorizationRules adds intelligent categorization rules
func (tcm *TaskCategoryManager) addCategorizationRules() {
	rules := []CategoryRule{
		// Proposal rules
		{
			Name:        "Proposal Keywords",
			Patterns:    []string{"proposal", "thesis proposal", "research proposal", "defense proposal"},
			Category:    "PROPOSAL",
			Priority:    1,
			Description: "Tasks related to PhD proposal",
		},
		{
			Name:        "Proposal Writing",
			Patterns:    []string{"write proposal", "draft proposal", "proposal writing", "proposal draft"},
			Category:    "PROPOSAL",
			Priority:    1,
			Description: "Proposal writing tasks",
		},
		
		// Laser system rules
		{
			Name:        "Laser Keywords",
			Patterns:    []string{"laser", "laser system", "laser setup", "laser maintenance", "laser alignment"},
			Category:    "LASER",
			Priority:    1,
			Description: "Laser system related tasks",
		},
		{
			Name:        "Laser Equipment",
			Patterns:    []string{"laser equipment", "laser hardware", "laser calibration", "laser repair"},
			Category:    "LASER",
			Priority:    1,
			Description: "Laser equipment tasks",
		},
		
		// Imaging rules
		{
			Name:        "Imaging Keywords",
			Patterns:    []string{"imaging", "image", "microscopy", "microscope", "data collection", "imaging experiment"},
			Category:    "IMAGING",
			Priority:    1,
			Description: "Imaging and data collection tasks",
		},
		{
			Name:        "Imaging Equipment",
			Patterns:    []string{"microscope setup", "imaging setup", "camera", "optics", "imaging system"},
			Category:    "IMAGING",
			Priority:    1,
			Description: "Imaging equipment tasks",
		},
		{
			Name:        "Surgery Keywords",
			Patterns:    []string{"surgery", "surgical", "cranial window", "inject", "injection", "mouse surgery"},
			Category:    "IMAGING",
			Priority:    1,
			Description: "Surgical procedures for imaging",
		},
		
		// Administrative rules
		{
			Name:        "Admin Keywords",
			Patterns:    []string{"admin", "administrative", "paperwork", "form", "grant application", "funding application"},
			Category:    "ADMIN",
			Priority:    1,
			Description: "Administrative tasks",
		},
		{
			Name:        "Meeting Keywords",
			Patterns:    []string{"meeting", "lab meeting", "group meeting", "committee", "presentation"},
			Category:    "ADMIN",
			Priority:    2,
			Description: "Meeting and presentation tasks",
		},
		{
			Name:        "Grant Keywords",
			Patterns:    []string{"grant", "funding", "application", "budget", "financial"},
			Category:    "ADMIN",
			Priority:    1,
			Description: "Grant and funding tasks",
		},
		
		// Dissertation rules
		{
			Name:        "Dissertation Keywords",
			Patterns:    []string{"dissertation", "thesis", "defense", "final defense", "thesis defense"},
			Category:    "DISSERTATION",
			Priority:    1,
			Description: "Dissertation related tasks",
		},
		{
			Name:        "Writing Keywords",
			Patterns:    []string{"writing", "write", "draft", "manuscript", "chapter", "section"},
			Category:    "DISSERTATION",
			Priority:    2,
			Description: "Writing tasks",
		},
		
		// Research rules
		{
			Name:        "Research Keywords",
			Patterns:    []string{"research", "study", "investigation", "data analysis", "analyze data"},
			Category:    "RESEARCH",
			Priority:    2,
			Description: "General research tasks",
		},
		{
			Name:        "Literature Keywords",
			Patterns:    []string{"literature", "literature review", "reading", "bibliography"},
			Category:    "RESEARCH",
			Priority:    2,
			Description: "Literature review tasks",
		},
		
		// Publication rules
		{
			Name:        "Publication Keywords",
			Patterns:    []string{"publication", "publish", "journal submission", "conference submission", "paper submission"},
			Category:    "PUBLICATION",
			Priority:    1,
			Description: "Publication tasks",
		},
		{
			Name:        "Manuscript Keywords",
			Patterns:    []string{"manuscript", "article", "paper writing", "journal article", "submit manuscript"},
			Category:    "PUBLICATION",
			Priority:    1,
			Description: "Manuscript writing tasks",
		},
	}
	
	tcm.rules = rules
}

// CategorizeTask automatically categorizes a task based on its name and description
func (tcm *TaskCategoryManager) CategorizeTask(task *Task) string {
	if task == nil {
		return "RESEARCH" // Default category
	}
	
	// Combine name and description for pattern matching
	text := strings.ToLower(task.Name + " " + task.Description)
	
	// Find the best matching rule
	bestMatch := ""
	bestPriority := 999
	
	for _, rule := range tcm.rules {
		for _, pattern := range rule.Patterns {
			if strings.Contains(text, strings.ToLower(pattern)) {
				if rule.Priority < bestPriority {
					bestMatch = rule.Category
					bestPriority = rule.Priority
				}
			}
		}
	}
	
	if bestMatch != "" {
		return bestMatch
	}
	
	// Default categorization based on task properties
	return tcm.getDefaultCategory(task)
}

// getDefaultCategory provides fallback categorization
func (tcm *TaskCategoryManager) getDefaultCategory(task *Task) string {
	// Check if task is a milestone
	if task.IsMilestone {
		return "DISSERTATION"
	}
	
	// Check if task has dependencies (likely research work)
	if len([]string{}) > 0 {
		return "RESEARCH"
	}
	
	// Check if task is very short-term (1 day or less) - likely admin
	if task.GetDuration() <= 1 {
		return "ADMIN"
	}
	
	// Default to research
	return "RESEARCH"
}

// GetCategory returns a category by name
func (tcm *TaskCategoryManager) GetCategory(name string) (TaskCategory, bool) {
	category, exists := tcm.categories[name]
	return category, exists
}

// GetAllCategories returns all available categories
func (tcm *TaskCategoryManager) GetAllCategories() []TaskCategory {
	categories := make([]TaskCategory, 0, len(tcm.categories))
	for _, category := range tcm.categories {
		categories = append(categories, category)
	}
	
	// Sort by priority
	sort.Slice(categories, func(i, j int) bool {
		return categories[i].Priority < categories[j].Priority
	})
	
	return categories
}

// AddCustomCategory adds a custom category
func (tcm *TaskCategoryManager) AddCustomCategory(category TaskCategory) {
	tcm.categories[category.Name] = category
}

// TaskDateCalculator provides advanced date calculations for tasks
type TaskDateCalculator struct {
	workDays    map[time.Weekday]bool
	holidays    []time.Time
	workHours   int
	hoursPerDay int
}

// NewTaskDateCalculator creates a new date calculator with default settings
func NewTaskDateCalculator() *TaskDateCalculator {
	tdc := &TaskDateCalculator{
		workDays:    make(map[time.Weekday]bool),
		holidays:    make([]time.Time, 0),
		workHours:   8,
		hoursPerDay: 8,
	}
	
	// Set default work days (Monday to Friday)
	for i := 1; i <= 5; i++ {
		tdc.workDays[time.Weekday(i)] = true
	}
	
	// Add common holidays (example)
	tdc.addCommonHolidays()
	
	return tdc
}

// addCommonHolidays adds common US holidays
func (tdc *TaskDateCalculator) addCommonHolidays() {
	year := time.Now().Year()
	
	// New Year's Day
	tdc.holidays = append(tdc.holidays, time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC))
	
	// Independence Day
	tdc.holidays = append(tdc.holidays, time.Date(year, 7, 4, 0, 0, 0, 0, time.UTC))
	
	// Christmas Day
	tdc.holidays = append(tdc.holidays, time.Date(year, 12, 25, 0, 0, 0, 0, time.UTC))
}

// AddHoliday adds a holiday to the calculator
func (tdc *TaskDateCalculator) AddHoliday(date time.Time) {
	tdc.holidays = append(tdc.holidays, date)
}

// IsWorkDay checks if a date is a work day
func (tdc *TaskDateCalculator) IsWorkDay(date time.Time) bool {
	// Check if it's a weekend
	if !tdc.workDays[date.Weekday()] {
		return false
	}
	
	// Check if it's a holiday
	for _, holiday := range tdc.holidays {
		if date.Year() == holiday.Year() && date.Month() == holiday.Month() && date.Day() == holiday.Day() {
			return false
		}
	}
	
	return true
}

// GetWorkDaysBetween calculates the number of work days between two dates
func (tdc *TaskDateCalculator) GetWorkDaysBetween(start, end time.Time) int {
	if start.After(end) {
		return 0
	}
	
	workDays := 0
	current := start
	
	for !current.After(end) {
		if tdc.IsWorkDay(current) {
			workDays++
		}
		current = current.AddDate(0, 0, 1)
	}
	
	return workDays
}

// GetNextWorkDay returns the next work day after a given date
func (tdc *TaskDateCalculator) GetNextWorkDay(date time.Time) time.Time {
	next := date.AddDate(0, 0, 1)
	
	for !tdc.IsWorkDay(next) {
		next = next.AddDate(0, 0, 1)
	}
	
	return next
}

// GetPreviousWorkDay returns the previous work day before a given date
func (tdc *TaskDateCalculator) GetPreviousWorkDay(date time.Time) time.Time {
	prev := date.AddDate(0, 0, -1)
	
	for !tdc.IsWorkDay(prev) {
		prev = prev.AddDate(0, 0, -1)
	}
	
	return prev
}

// CalculateTaskEndDate calculates the end date for a task based on work days
func (tdc *TaskDateCalculator) CalculateTaskEndDate(startDate time.Time, workDays int) time.Time {
	if workDays <= 0 {
		return startDate
	}
	
	current := startDate
	daysAdded := 0
	
	for daysAdded < workDays {
		if tdc.IsWorkDay(current) {
			daysAdded++
		}
		if daysAdded < workDays {
			current = current.AddDate(0, 0, 1)
		}
	}
	
	return current
}

// GetTaskWorkDays calculates the number of work days for a task
func (tdc *TaskDateCalculator) GetTaskWorkDays(task *Task) int {
	if task == nil || task.StartDate.IsZero() || task.EndDate.IsZero() {
		return 0
	}
	
	return tdc.GetWorkDaysBetween(task.StartDate, task.EndDate)
}

// TaskTimelineAnalyzer provides timeline analysis for tasks
type TaskTimelineAnalyzer struct {
	calculator *TaskDateCalculator
}

// NewTaskTimelineAnalyzer creates a new timeline analyzer
func NewTaskTimelineAnalyzer() *TaskTimelineAnalyzer {
	return &TaskTimelineAnalyzer{
		calculator: NewTaskDateCalculator(),
	}
}

// AnalyzeTaskTimeline analyzes a task's timeline and provides insights
func (tta *TaskTimelineAnalyzer) AnalyzeTaskTimeline(task *Task) *TaskTimelineAnalysis {
	if task == nil {
		return nil
	}
	
	analysis := &TaskTimelineAnalysis{
		TaskID:           task.Name,
		TaskName:         task.Name,
		StartDate:        task.StartDate,
		EndDate:          task.EndDate,
		TotalDays:        task.GetDuration(),
		WorkDays:         tta.calculator.GetTaskWorkDays(task),
		IsOverdue:        task.IsOverdue(),
		IsUpcoming:       task.IsUpcoming(),
		ProgressPercent:  task.GetProgressPercentage(),
		DaysRemaining:    tta.calculateDaysRemaining(task),
		WorkDaysRemaining: tta.calculateWorkDaysRemaining(task),
		RiskLevel:        tta.assessRiskLevel(task),
		Recommendations:  tta.generateRecommendations(task),
	}
	
	return analysis
}

// TaskTimelineAnalysis contains detailed timeline analysis
type TaskTimelineAnalysis struct {
	TaskID            string
	TaskName          string
	StartDate         time.Time
	EndDate           time.Time
	TotalDays         int
	WorkDays          int
	IsOverdue         bool
	IsUpcoming        bool
	ProgressPercent   float64
	DaysRemaining     int
	WorkDaysRemaining int
	RiskLevel         string
	Recommendations   []string
}

// calculateDaysRemaining calculates days remaining until task completion
func (tta *TaskTimelineAnalyzer) calculateDaysRemaining(task *Task) int {
	if task.EndDate.IsZero() {
		return 0
	}
	
	now := time.Now()
	if now.After(task.EndDate) {
		return 0
	}
	
	return int(task.EndDate.Sub(now).Hours() / 24)
}

// calculateWorkDaysRemaining calculates work days remaining
func (tta *TaskTimelineAnalyzer) calculateWorkDaysRemaining(task *Task) int {
	if task.EndDate.IsZero() {
		return 0
	}
	
	now := time.Now()
	if now.After(task.EndDate) {
		return 0
	}
	
	return tta.calculator.GetWorkDaysBetween(now, task.EndDate)
}

// assessRiskLevel assesses the risk level of a task
func (tta *TaskTimelineAnalyzer) assessRiskLevel(task *Task) string {
	if task.IsOverdue() {
		return "HIGH"
	}
	
	workDaysRemaining := tta.calculateWorkDaysRemaining(task)
	progress := task.GetProgressPercentage()
	
	// High risk if less than 3 work days remaining and less than 50% complete
	if workDaysRemaining <= 3 && progress < 50.0 {
		return "HIGH"
	}
	
	// Medium risk if less than 7 work days remaining and less than 75% complete
	if workDaysRemaining <= 7 && progress < 75.0 {
		return "MEDIUM"
	}
	
	// Low risk otherwise
	return "LOW"
}

// generateRecommendations generates recommendations for task management
func (tta *TaskTimelineAnalyzer) generateRecommendations(task *Task) []string {
	var recommendations []string
	
	if task.IsOverdue() {
		recommendations = append(recommendations, "Task is overdue - consider extending deadline or prioritizing")
	}
	
	if task.IsUpcoming() {
		recommendations = append(recommendations, "Task starting soon - prepare resources and dependencies")
	}
	
	workDaysRemaining := tta.calculateWorkDaysRemaining(task)
	progress := task.GetProgressPercentage()
	
	if workDaysRemaining <= 3 && progress < 50.0 {
		recommendations = append(recommendations, "High risk of delay - consider additional resources or scope reduction")
	}
	
	if len([]string{}) > 0 {
		recommendations = append(recommendations, "Task has dependencies - ensure they are completed on time")
	}
	
	if task.Priority > 5 {
		recommendations = append(recommendations, "High priority task - monitor closely")
	}
	
	return recommendations
}

// String returns a string representation of the timeline analysis
func (tta *TaskTimelineAnalysis) String() string {
	return fmt.Sprintf("TaskTimelineAnalysis[%s: %s, Risk: %s, Progress: %.1f%%, Work Days: %d/%d]",
		tta.TaskID, tta.TaskName, tta.RiskLevel, tta.ProgressPercent, tta.WorkDaysRemaining, tta.WorkDays)
}
