package calendar

import (
	"fmt"
	"math"
	"sort"
	"time"
)

// MonthBoundaryEngine handles month boundary transitions and grid continuity
type MonthBoundaryEngine struct {
	gridConfig        *GridConfig
	visualConstraints *VisualConstraints
	boundaryRules     []BoundaryRule
	transitionRules   []TransitionRule
	continuityRules   []ContinuityRule
}

// BoundaryRule defines how tasks should behave at month boundaries
type BoundaryRule struct {
	Name        string
	Description string
	Priority    int
	Condition   func(*IntegratedTaskBar, *MonthBoundaryContext) bool
	Action      func(*IntegratedTaskBar, *MonthBoundaryContext) *BoundaryAction
}

// TransitionRule defines how tasks should transition between months
type TransitionRule struct {
	Name        string
	Description string
	Priority    int
	Condition   func(*IntegratedTaskBar, *MonthBoundaryContext) bool
	Action      func(*IntegratedTaskBar, *MonthBoundaryContext) *TransitionAction
}

// ContinuityRule defines how to maintain visual continuity across months
type ContinuityRule struct {
	Name        string
	Description string
	Priority    int
	Condition   func([]*IntegratedTaskBar, *MonthBoundaryContext) bool
	Action      func([]*IntegratedTaskBar, *MonthBoundaryContext) *ContinuityAction
}

// MonthBoundaryContext provides context for month boundary decisions
type MonthBoundaryContext struct {
	CurrentMonth      time.Month
	NextMonth         time.Month
	CurrentYear       int
	NextYear          int
	CalendarStart     time.Time
	CalendarEnd       time.Time
	DayWidth          float64
	DayHeight         float64
	MonthBoundaryGap  float64
	TransitionBuffer  float64
	VisualSettings    *IntegratedVisualSettings
	GridConstraints   *GridConstraints
	TaskDensity       float64
	OverlapCount      int
	ConflictCount     int
}

// BoundaryAction defines how a task should be handled at month boundaries
type BoundaryAction struct {
	SplitTask         bool
	CreateContinuation bool
	AdjustPosition    bool
	AdjustWidth       bool
	AdjustHeight      bool
	NewX              float64
	NewY              float64
	NewWidth          float64
	NewHeight         float64
	ContinuationID    string
	VisualStyle       *BoundaryVisualStyle
	Priority          int
}

// TransitionAction defines how a task should transition between months
type TransitionAction struct {
	SmoothTransition  bool
	FadeIn           bool
	FadeOut          bool
	SlideAnimation   bool
	ScaleAnimation   bool
	Duration         time.Duration
	EasingFunction   EasingFunction
	VisualEffects    []VisualEffect
	Priority         int
}

// ContinuityAction defines how to maintain visual continuity
type ContinuityAction struct {
	MaintainAlignment bool
	PreserveSpacing   bool
	ConsistentColors  bool
	ConsistentSizes   bool
	VisualConnections []VisualConnection
	Priority          int
}

// BoundaryVisualStyle defines visual styling for month boundaries
type BoundaryVisualStyle struct {
	BorderStyle       BorderStyle
	BorderColor       string
	BorderWidth       float64
	BackgroundColor   string
	Opacity           float64
	ShadowEnabled     bool
	ShadowColor       string
	ShadowBlur        float64
	ShadowOffsetX     float64
	ShadowOffsetY     float64
	HighlightEnabled  bool
	HighlightColor    string
	HighlightWidth    float64
}

// BorderStyle defines the style of borders at month boundaries
type BorderStyle string

const (
	BorderSolid      BorderStyle = "SOLID"
	BorderDashed     BorderStyle = "DASHED"
	BorderDotted     BorderStyle = "DOTTED"
	BorderDouble     BorderStyle = "DOUBLE"
	BorderGroove     BorderStyle = "GROOVE"
	BorderRidge      BorderStyle = "RIDGE"
	BorderInset      BorderStyle = "INSET"
	BorderOutset     BorderStyle = "OUTSET"
	BorderNone       BorderStyle = "NONE"
)

// EasingFunction defines easing functions for transitions
type EasingFunction string

const (
	EasingLinear      EasingFunction = "LINEAR"
	EasingEaseIn      EasingFunction = "EASE_IN"
	EasingEaseOut     EasingFunction = "EASE_OUT"
	EasingEaseInOut   EasingFunction = "EASE_IN_OUT"
	EasingBounce      EasingFunction = "BOUNCE"
	EasingElastic     EasingFunction = "ELASTIC"
	EasingBack        EasingFunction = "BACK"
	EasingCubic       EasingFunction = "CUBIC"
	EasingQuart       EasingFunction = "QUART"
	EasingQuint       EasingFunction = "QUINT"
	EasingSine        EasingFunction = "SINE"
	EasingExpo        EasingFunction = "EXPO"
	EasingCirc        EasingFunction = "CIRC"
)

// VisualEffect defines visual effects for transitions
type VisualEffect string

const (
	EffectFadeIn      VisualEffect = "FADE_IN"
	EffectFadeOut     VisualEffect = "FADE_OUT"
	EffectSlideLeft   VisualEffect = "SLIDE_LEFT"
	EffectSlideRight  VisualEffect = "SLIDE_RIGHT"
	EffectSlideUp     VisualEffect = "SLIDE_UP"
	EffectSlideDown   VisualEffect = "SLIDE_DOWN"
	EffectScaleIn     VisualEffect = "SCALE_IN"
	EffectScaleOut    VisualEffect = "SCALE_OUT"
	EffectRotateIn    VisualEffect = "ROTATE_IN"
	EffectRotateOut   VisualEffect = "ROTATE_OUT"
	EffectFlipIn      VisualEffect = "FLIP_IN"
	EffectFlipOut     VisualEffect = "FLIP_OUT"
	EffectZoomIn      VisualEffect = "ZOOM_IN"
	EffectZoomOut     VisualEffect = "ZOOM_OUT"
)

// VisualConnection defines visual connections between months
type VisualConnection struct {
	FromTaskID    string
	ToTaskID      string
	ConnectionType ConnectionType
	LineStyle     LineStyle
	LineColor     string
	LineWidth     float64
	ArrowStyle    ArrowStyle
	Label         string
	Priority      int
}

// ConnectionType defines the type of visual connection
type ConnectionType string

const (
	ConnectionArrow      ConnectionType = "ARROW"
	ConnectionLine       ConnectionType = "LINE"
	ConnectionCurve      ConnectionType = "CURVE"
	ConnectionDashed     ConnectionType = "DASHED"
	ConnectionDotted     ConnectionType = "DOTTED"
	ConnectionThick      ConnectionType = "THICK"
	ConnectionThin       ConnectionType = "THIN"
	ConnectionDouble     ConnectionType = "DOUBLE"
)

// LineStyle defines the style of connection lines
type LineStyle string

const (
	LineSolid      LineStyle = "SOLID"
	LineDashed     LineStyle = "DASHED"
	LineDotted     LineStyle = "DOTTED"
	LineDashDot    LineStyle = "DASH_DOT"
	LineDashDotDot LineStyle = "DASH_DOT_DOT"
)

// ArrowStyle defines the style of arrows
type ArrowStyle string

const (
	ArrowNone       ArrowStyle = "NONE"
	ArrowSimple     ArrowStyle = "SIMPLE"
	ArrowFilled     ArrowStyle = "FILLED"
	ArrowHollow     ArrowStyle = "HOLLOW"
	ArrowDouble     ArrowStyle = "DOUBLE"
	ArrowCurved     ArrowStyle = "CURVED"
	ArrowBarbed     ArrowStyle = "BARBED"
	ArrowFeathered  ArrowStyle = "FEATHERED"
)

// MonthBoundaryResult contains the result of month boundary processing
type MonthBoundaryResult struct {
	ProcessedBars     []*IntegratedTaskBar
	Continuations     []*TaskContinuation
	Transitions       []*TaskTransition
	VisualConnections []*VisualConnection
	BoundaryMetrics   *BoundaryMetrics
	Recommendations   []string
	AnalysisDate      time.Time
}

// TaskContinuation represents a task continuation across month boundaries
type TaskContinuation struct {
	OriginalTaskID    string
	ContinuationID    string
	StartMonth        time.Month
	EndMonth          time.Month
	StartYear         int
	EndYear           int
	ContinuationType  ContinuationType
	VisualStyle       *BoundaryVisualStyle
	ConnectionStyle   *VisualConnection
	Priority          int
}

// TaskTransition represents a task transition between months
type TaskTransition struct {
	TaskID           string
	FromMonth        time.Month
	ToMonth          time.Month
	TransitionType   TransitionType
	Animation        *TransitionAnimation
	VisualEffects    []VisualEffect
	Duration         time.Duration
	EasingFunction   EasingFunction
	Priority         int
}

// ContinuationType defines the type of task continuation
type ContinuationType string

const (
	ContinuationSplit     ContinuationType = "SPLIT"
	ContinuationExtend    ContinuationType = "EXTEND"
	ContinuationWrap      ContinuationType = "WRAP"
	ContinuationOverflow  ContinuationType = "OVERFLOW"
	ContinuationTruncate  ContinuationType = "TRUNCATE"
	ContinuationMinimize  ContinuationType = "MINIMIZE"
	ContinuationCollapse  ContinuationType = "COLLAPSE"
)

// TransitionType defines the type of task transition
type TransitionType string

const (
	TransitionSmooth     TransitionType = "SMOOTH"
	TransitionFade       TransitionType = "FADE"
	TransitionSlide      TransitionType = "SLIDE"
	TransitionScale      TransitionType = "SCALE"
	TransitionRotate     TransitionType = "ROTATE"
	TransitionFlip       TransitionType = "FLIP"
	TransitionZoom       TransitionType = "ZOOM"
	TransitionBounce     TransitionType = "BOUNCE"
	TransitionElastic    TransitionType = "ELASTIC"
)

// TransitionAnimation defines animation properties for transitions
type TransitionAnimation struct {
	Type             TransitionType
	Duration         time.Duration
	EasingFunction   EasingFunction
	Delay            time.Duration
	IterationCount   int
	Direction        AnimationDirection
	FillMode         FillMode
	PlayState        PlayState
}

// AnimationDirection defines the direction of animation
type AnimationDirection string

const (
	DirectionNormal     AnimationDirection = "NORMAL"
	DirectionReverse    AnimationDirection = "REVERSE"
	DirectionAlternate  AnimationDirection = "ALTERNATE"
	DirectionAlternateReverse AnimationDirection = "ALTERNATE_REVERSE"
)

// FillMode defines how animations fill their target values
type FillMode string

const (
	FillNone       FillMode = "NONE"
	FillForwards   FillMode = "FORWARDS"
	FillBackwards  FillMode = "BACKWARDS"
	FillBoth       FillMode = "BOTH"
)

// PlayState defines the play state of animations
type PlayState string

const (
	PlayRunning    PlayState = "RUNNING"
	PlayPaused     PlayState = "PAUSED"
	PlayStopped    PlayState = "STOPPED"
)

// BoundaryMetrics contains metrics about month boundary processing
type BoundaryMetrics struct {
	TotalTasks           int
	ProcessedTasks       int
	ContinuationsCreated int
	TransitionsApplied   int
	VisualConnections    int
	BoundaryConflicts    int
	TransitionErrors     int
	ContinuityScore      float64
	VisualConsistency    float64
	TransitionSmoothness float64
	GridContinuity       float64
	SpaceEfficiency      float64
	VisualBalance        float64
}

// NewMonthBoundaryEngine creates a new month boundary engine
func NewMonthBoundaryEngine(gridConfig *GridConfig) *MonthBoundaryEngine {
	engine := &MonthBoundaryEngine{
		gridConfig:        gridConfig,
		visualConstraints: gridConfig.VisualConstraints,
		boundaryRules:     []BoundaryRule{},
		transitionRules:   []TransitionRule{},
		continuityRules:   []ContinuityRule{},
	}
	
	// Add default boundary rules
	engine.addDefaultBoundaryRules()
	
	// Add default transition rules
	engine.addDefaultTransitionRules()
	
	// Add default continuity rules
	engine.addDefaultContinuityRules()
	
	return engine
}

// ProcessMonthBoundaries processes month boundary transitions for task bars
func (mbe *MonthBoundaryEngine) ProcessMonthBoundaries(taskBars []*IntegratedTaskBar, currentMonth time.Month, currentYear int) (*MonthBoundaryResult, error) {
	// Create month boundary context
	context := &MonthBoundaryContext{
		CurrentMonth:     currentMonth,
		NextMonth:        mbe.getNextMonth(currentMonth),
		CurrentYear:      currentYear,
		NextYear:         mbe.getNextYear(currentMonth, currentYear),
		CalendarStart:    mbe.gridConfig.CalendarStart,
		CalendarEnd:      mbe.gridConfig.CalendarEnd,
		DayWidth:         mbe.gridConfig.DayWidth,
		DayHeight:        mbe.gridConfig.DayHeight,
		MonthBoundaryGap: mbe.gridConfig.MonthBoundaryGap,
		TransitionBuffer: 2.0,
		VisualSettings:   mbe.getDefaultVisualSettings(),
		GridConstraints:  mbe.getDefaultGridConstraints(),
		TaskDensity:      mbe.calculateTaskDensity(taskBars),
		OverlapCount:     mbe.countOverlaps(taskBars),
		ConflictCount:    mbe.countConflicts(taskBars),
	}
	
	// Process boundary rules
	processedBars := mbe.applyBoundaryRules(taskBars, context)
	
	// Process transition rules
	transitions := mbe.applyTransitionRules(processedBars, context)
	
	// Process continuity rules
	continuations := mbe.applyContinuityRules(processedBars, context)
	
	// Create visual connections
	visualConnections := mbe.createVisualConnections(processedBars, continuations, context)
	
	// Calculate metrics
	metrics := mbe.calculateBoundaryMetrics(processedBars, continuations, transitions, context)
	
	// Generate recommendations
	recommendations := mbe.generateBoundaryRecommendations(metrics, context)
	
	return &MonthBoundaryResult{
		ProcessedBars:     processedBars,
		Continuations:     continuations,
		Transitions:       transitions,
		VisualConnections: visualConnections,
		BoundaryMetrics:   metrics,
		Recommendations:   recommendations,
		AnalysisDate:      time.Now(),
	}, nil
}

// applyBoundaryRules applies boundary rules to task bars
func (mbe *MonthBoundaryEngine) applyBoundaryRules(bars []*IntegratedTaskBar, context *MonthBoundaryContext) []*IntegratedTaskBar {
	// Sort rules by priority
	sort.Slice(mbe.boundaryRules, func(i, j int) bool {
		return mbe.boundaryRules[i].Priority > mbe.boundaryRules[j].Priority
	})
	
	// Apply rules to each bar
	for _, bar := range bars {
		for _, rule := range mbe.boundaryRules {
			if rule.Condition(bar, context) {
				action := rule.Action(bar, context)
				mbe.applyBoundaryAction(bar, action, context)
				break // Apply only the first matching rule
			}
		}
	}
	
	return bars
}

// applyTransitionRules applies transition rules to task bars
func (mbe *MonthBoundaryEngine) applyTransitionRules(bars []*IntegratedTaskBar, context *MonthBoundaryContext) []*TaskTransition {
	var transitions []*TaskTransition
	
	// Sort rules by priority
	sort.Slice(mbe.transitionRules, func(i, j int) bool {
		return mbe.transitionRules[i].Priority > mbe.transitionRules[j].Priority
	})
	
	// Apply rules to each bar
	for _, bar := range bars {
		for _, rule := range mbe.transitionRules {
			if rule.Condition(bar, context) {
				action := rule.Action(bar, context)
				transition := mbe.createTransition(bar, action, context)
				transitions = append(transitions, transition)
				break // Apply only the first matching rule
			}
		}
	}
	
	return transitions
}

// applyContinuityRules applies continuity rules to task bars
func (mbe *MonthBoundaryEngine) applyContinuityRules(bars []*IntegratedTaskBar, context *MonthBoundaryContext) []*TaskContinuation {
	var continuations []*TaskContinuation
	
	// Sort rules by priority
	sort.Slice(mbe.continuityRules, func(i, j int) bool {
		return mbe.continuityRules[i].Priority > mbe.continuityRules[j].Priority
	})
	
	// Apply rules to groups of bars
	for _, rule := range mbe.continuityRules {
		if rule.Condition(bars, context) {
			action := rule.Action(bars, context)
			continuation := mbe.createContinuation(bars, action, context)
			continuations = append(continuations, continuation)
		}
	}
	
	return continuations
}

// applyBoundaryAction applies a boundary action to a task bar
func (mbe *MonthBoundaryEngine) applyBoundaryAction(bar *IntegratedTaskBar, action *BoundaryAction, context *MonthBoundaryContext) {
	if action.SplitTask {
		// Split the task at month boundary
		mbe.splitTaskAtBoundary(bar, context)
	}
	
	if action.CreateContinuation {
		// Create a continuation for the next month
		mbe.createTaskContinuation(bar, context)
	}
	
	if action.AdjustPosition {
		bar.StartX = action.NewX
		bar.Y = action.NewY
	}
	
	if action.AdjustWidth {
		bar.Width = action.NewWidth
		bar.EndX = bar.StartX + bar.Width
	}
	
	if action.AdjustHeight {
		bar.Height = action.NewHeight
	}
	
	// Apply visual style
	if action.VisualStyle != nil {
		mbe.applyVisualStyle(bar, action.VisualStyle)
	}
}

// splitTaskAtBoundary splits a task at the month boundary
func (mbe *MonthBoundaryEngine) splitTaskAtBoundary(bar *IntegratedTaskBar, context *MonthBoundaryContext) {
	// Calculate the split point (end of current month)
	splitDate := mbe.getMonthEndDate(context.CurrentMonth, context.CurrentYear)
	splitX := mbe.calculateXPosition(splitDate, context)
	
	// Adjust the current bar to end at the split point
	bar.EndX = splitX
	bar.Width = bar.EndX - bar.StartX
	
	// Mark as month boundary
	bar.MonthBoundary = true
}

// createTaskContinuation creates a continuation for the next month
func (mbe *MonthBoundaryEngine) createTaskContinuation(bar *IntegratedTaskBar, context *MonthBoundaryContext) {
	// This would create a new task bar for the next month
	// Implementation depends on the specific requirements
}

// createTransition creates a transition for a task bar
func (mbe *MonthBoundaryEngine) createTransition(bar *IntegratedTaskBar, action *TransitionAction, context *MonthBoundaryContext) *TaskTransition {
	transition := &TaskTransition{
		TaskID:         bar.TaskID,
		FromMonth:      context.CurrentMonth,
		ToMonth:        context.NextMonth,
		TransitionType: TransitionSmooth,
		Duration:       action.Duration,
		EasingFunction: action.EasingFunction,
		VisualEffects:  action.VisualEffects,
		Priority:       action.Priority,
	}
	
	// Create animation if specified
	if action.SmoothTransition {
		transition.Animation = &TransitionAnimation{
			Type:           TransitionSmooth,
			Duration:       action.Duration,
			EasingFunction: action.EasingFunction,
			Direction:      DirectionNormal,
			FillMode:       FillBoth,
			PlayState:      PlayRunning,
		}
	}
	
	return transition
}

// createContinuation creates a continuation for task bars
func (mbe *MonthBoundaryEngine) createContinuation(bars []*IntegratedTaskBar, action *ContinuityAction, context *MonthBoundaryContext) *TaskContinuation {
	// Find the primary task for continuation
	var primaryTask *IntegratedTaskBar
	for _, bar := range bars {
		if bar.MonthBoundary {
			primaryTask = bar
			break
		}
	}
	
	if primaryTask == nil {
		return nil
	}
	
	continuation := &TaskContinuation{
		OriginalTaskID:   primaryTask.TaskID,
		ContinuationID:   fmt.Sprintf("%s_cont_%d_%d", primaryTask.TaskID, context.NextYear, int(context.NextMonth)),
		StartMonth:       context.NextMonth,
		EndMonth:         context.NextMonth,
		StartYear:        context.NextYear,
		EndYear:          context.NextYear,
		ContinuationType: ContinuationExtend,
		Priority:         action.Priority,
	}
	
	// Create visual connections if specified
	if len(action.VisualConnections) > 0 {
		continuation.ConnectionStyle = &action.VisualConnections[0]
	}
	
	return continuation
}

// createVisualConnections creates visual connections between months
func (mbe *MonthBoundaryEngine) createVisualConnections(bars []*IntegratedTaskBar, continuations []*TaskContinuation, context *MonthBoundaryContext) []*VisualConnection {
	var connections []*VisualConnection
	
	// Create connections for continuations
	for _, continuation := range continuations {
		connection := &VisualConnection{
			FromTaskID:    continuation.OriginalTaskID,
			ToTaskID:      continuation.ContinuationID,
			ConnectionType: ConnectionArrow,
			LineStyle:     LineSolid,
			LineColor:     "#666666",
			LineWidth:     1.0,
			ArrowStyle:    ArrowSimple,
			Label:         "Continues",
			Priority:      continuation.Priority,
		}
		connections = append(connections, connection)
	}
	
	return connections
}

// calculateBoundaryMetrics calculates metrics for month boundary processing
func (mbe *MonthBoundaryEngine) calculateBoundaryMetrics(bars []*IntegratedTaskBar, continuations []*TaskContinuation, transitions []*TaskTransition, context *MonthBoundaryContext) *BoundaryMetrics {
	metrics := &BoundaryMetrics{
		TotalTasks:           len(bars),
		ProcessedTasks:       len(bars),
		ContinuationsCreated: len(continuations),
		TransitionsApplied:   len(transitions),
		VisualConnections:    len(continuations), // Simplified
	}
	
	// Calculate continuity score
	metrics.ContinuityScore = mbe.calculateContinuityScore(bars, continuations, context)
	
	// Calculate visual consistency
	metrics.VisualConsistency = mbe.calculateVisualConsistency(bars, context)
	
	// Calculate transition smoothness
	metrics.TransitionSmoothness = mbe.calculateTransitionSmoothness(transitions, context)
	
	// Calculate grid continuity
	metrics.GridContinuity = mbe.calculateGridContinuity(bars, context)
	
	// Calculate space efficiency
	metrics.SpaceEfficiency = mbe.calculateSpaceEfficiency(bars, context)
	
	// Calculate visual balance
	metrics.VisualBalance = mbe.calculateVisualBalance(bars, context)
	
	return metrics
}

// calculateContinuityScore calculates the continuity score across months
func (mbe *MonthBoundaryEngine) calculateContinuityScore(bars []*IntegratedTaskBar, continuations []*TaskContinuation, context *MonthBoundaryContext) float64 {
	if len(bars) == 0 {
		return 1.0
	}
	
	// Count tasks that have continuations
	continuationCount := 0
	for _, bar := range bars {
		if bar.MonthBoundary {
			continuationCount++
		}
	}
	
	return float64(continuationCount) / float64(len(bars))
}

// calculateVisualConsistency calculates visual consistency across months
func (mbe *MonthBoundaryEngine) calculateVisualConsistency(bars []*IntegratedTaskBar, context *MonthBoundaryContext) float64 {
	if len(bars) == 0 {
		return 1.0
	}
	
	// Calculate consistency based on visual properties
	var consistencyScore float64
	for _, bar := range bars {
		// Check if visual properties are consistent
		if bar.Color != "" && bar.BorderColor != "" {
			consistencyScore += 1.0
		}
	}
	
	return consistencyScore / float64(len(bars))
}

// calculateTransitionSmoothness calculates transition smoothness
func (mbe *MonthBoundaryEngine) calculateTransitionSmoothness(transitions []*TaskTransition, context *MonthBoundaryContext) float64 {
	if len(transitions) == 0 {
		return 1.0
	}
	
	// Calculate smoothness based on transition types
	var smoothnessScore float64
	for _, transition := range transitions {
		switch transition.TransitionType {
		case TransitionSmooth:
			smoothnessScore += 1.0
		case TransitionFade:
			smoothnessScore += 0.8
		case TransitionSlide:
			smoothnessScore += 0.7
		default:
			smoothnessScore += 0.5
		}
	}
	
	return smoothnessScore / float64(len(transitions))
}

// calculateGridContinuity calculates grid continuity across months
func (mbe *MonthBoundaryEngine) calculateGridContinuity(bars []*IntegratedTaskBar, context *MonthBoundaryContext) float64 {
	if len(bars) == 0 {
		return 1.0
	}
	
	// Calculate continuity based on grid alignment
	var continuityScore float64
	for _, bar := range bars {
		// Check if bar is aligned to grid
		if context.GridConstraints.SnapToGrid {
			xAligned := math.Mod(bar.StartX, context.GridConstraints.GridResolution) < context.GridConstraints.AlignmentTolerance
			yAligned := math.Mod(bar.Y, context.GridConstraints.GridResolution) < context.GridConstraints.AlignmentTolerance
			if xAligned && yAligned {
				continuityScore += 1.0
			}
		}
	}
	
	return continuityScore / float64(len(bars))
}

// calculateSpaceEfficiency calculates space efficiency across months
func (mbe *MonthBoundaryEngine) calculateSpaceEfficiency(bars []*IntegratedTaskBar, context *MonthBoundaryContext) float64 {
	if len(bars) == 0 {
		return 0.0
	}
	
	// Calculate used space
	var usedSpace float64
	for _, bar := range bars {
		usedSpace += bar.Width * bar.Height
	}
	
	// Calculate available space
	availableSpace := context.DayWidth * context.DayHeight * 7.0 * 4.0 // 7 days, 4 weeks
	
	return usedSpace / availableSpace
}

// calculateVisualBalance calculates visual balance across months
func (mbe *MonthBoundaryEngine) calculateVisualBalance(bars []*IntegratedTaskBar, context *MonthBoundaryContext) float64 {
	if len(bars) == 0 {
		return 1.0
	}
	
	// Calculate center of mass
	var totalX, totalY, totalWeight float64
	for _, bar := range bars {
		weight := bar.VisualWeight
		centerX := bar.StartX + bar.Width/2
		centerY := bar.Y + bar.Height/2
		
		totalX += centerX * weight
		totalY += centerY * weight
		totalWeight += weight
	}
	
	if totalWeight == 0 {
		return 1.0
	}
	
	centerOfMassX := totalX / totalWeight
	centerOfMassY := totalY / totalWeight
	
	// Calculate distance from center of available space
	availableCenterX := context.DayWidth * 7.0 / 2
	availableCenterY := context.DayHeight * 4.0 / 2
	
	distanceFromCenter := math.Sqrt(
		math.Pow(centerOfMassX-availableCenterX, 2) + 
		math.Pow(centerOfMassY-availableCenterY, 2),
	)
	
	// Normalize to 0-1 scale
	maxDistance := math.Sqrt(
		math.Pow(availableCenterX, 2) + 
		math.Pow(availableCenterY, 2),
	)
	
	return 1.0 - (distanceFromCenter / maxDistance)
}

// generateBoundaryRecommendations generates recommendations based on boundary metrics
func (mbe *MonthBoundaryEngine) generateBoundaryRecommendations(metrics *BoundaryMetrics, context *MonthBoundaryContext) []string {
	var recommendations []string
	
	// Continuity recommendations
	if metrics.ContinuityScore < 0.7 {
		recommendations = append(recommendations, "Consider adding more task continuations to improve month-to-month continuity")
	}
	
	// Visual consistency recommendations
	if metrics.VisualConsistency < 0.8 {
		recommendations = append(recommendations, "Improve visual consistency by standardizing task colors and styles")
	}
	
	// Transition smoothness recommendations
	if metrics.TransitionSmoothness < 0.6 {
		recommendations = append(recommendations, "Use smoother transitions between months to improve user experience")
	}
	
	// Grid continuity recommendations
	if metrics.GridContinuity < 0.8 {
		recommendations = append(recommendations, "Enable grid snapping to improve alignment consistency across months")
	}
	
	// Space efficiency recommendations
	if metrics.SpaceEfficiency < 0.7 {
		recommendations = append(recommendations, "Optimize space usage to improve calendar density")
	}
	
	// Visual balance recommendations
	if metrics.VisualBalance < 0.6 {
		recommendations = append(recommendations, "Redistribute tasks to improve visual balance across months")
	}
	
	return recommendations
}

// Helper methods
func (mbe *MonthBoundaryEngine) getNextMonth(month time.Month) time.Month {
	if month == time.December {
		return time.January
	}
	return month + 1
}

func (mbe *MonthBoundaryEngine) getNextYear(month time.Month, year int) int {
	if month == time.December {
		return year + 1
	}
	return year
}

func (mbe *MonthBoundaryEngine) getMonthEndDate(month time.Month, year int) time.Time {
	// Get the last day of the month
	nextMonth := mbe.getNextMonth(month)
	nextYear := mbe.getNextYear(month, year)
	firstDayNextMonth := time.Date(nextYear, nextMonth, 1, 0, 0, 0, 0, time.UTC)
	return firstDayNextMonth.AddDate(0, 0, -1)
}

func (mbe *MonthBoundaryEngine) calculateXPosition(date time.Time, context *MonthBoundaryContext) float64 {
	daysFromStart := int(date.Sub(context.CalendarStart).Hours() / 24)
	return float64(daysFromStart) * context.DayWidth
}

func (mbe *MonthBoundaryEngine) calculateTaskDensity(bars []*IntegratedTaskBar) float64 {
	if len(bars) == 0 {
		return 0.0
	}
	
	// Calculate total calendar area
	totalArea := mbe.gridConfig.DayWidth * mbe.gridConfig.DayHeight * 7.0 * 4.0
	
	// Calculate average task area
	var totalTaskArea float64
	for _, bar := range bars {
		totalTaskArea += bar.Width * bar.Height
	}
	
	avgTaskArea := totalTaskArea / float64(len(bars))
	
	// Calculate density
	return (avgTaskArea * float64(len(bars))) / totalArea
}

func (mbe *MonthBoundaryEngine) countOverlaps(bars []*IntegratedTaskBar) int {
	count := 0
	for i := 0; i < len(bars); i++ {
		for j := i + 1; j < len(bars); j++ {
			if mbe.barsOverlap(bars[i], bars[j]) {
				count++
			}
		}
	}
	return count
}

func (mbe *MonthBoundaryEngine) countConflicts(bars []*IntegratedTaskBar) int {
	// For now, use overlap count as conflict count
	return mbe.countOverlaps(bars)
}

func (mbe *MonthBoundaryEngine) barsOverlap(bar1, bar2 *IntegratedTaskBar) bool {
	// Check horizontal overlap
	horizontalOverlap := bar1.StartX < bar2.EndX && bar2.StartX < bar1.EndX
	
	// Check vertical overlap
	verticalOverlap := bar1.Y < bar2.Y+bar2.Height && bar2.Y < bar1.Y+bar1.Height
	
	return horizontalOverlap && verticalOverlap
}

func (mbe *MonthBoundaryEngine) applyVisualStyle(bar *IntegratedTaskBar, style *BoundaryVisualStyle) {
	// Apply visual style to task bar
	if style.BorderColor != "" {
		bar.BorderColor = style.BorderColor
	}
	
	if style.BackgroundColor != "" {
		bar.Color = style.BackgroundColor
	}
	
	if style.Opacity > 0 {
		bar.Opacity = style.Opacity
	}
}

func (mbe *MonthBoundaryEngine) getDefaultVisualSettings() *IntegratedVisualSettings {
	return &IntegratedVisualSettings{
		ShowTaskNames:         true,
		ShowTaskDurations:     true,
		ShowTaskPriorities:    true,
		ShowConflictIndicators: true,
		CollapseThreshold:     5,
		AnimationEnabled:      false,
		HighlightConflicts:    true,
		ColorScheme:           "default",
		FontSize:              "small",
		TaskBarOpacity:        0.9,
		BorderWidth:           0.5,
	}
}

func (mbe *MonthBoundaryEngine) getDefaultGridConstraints() *GridConstraints {
	return &GridConstraints{
		MinTaskSpacing:     1.0,
		MaxTaskSpacing:     10.0,
		MinRowHeight:       8.0,
		MaxRowHeight:       20.0,
		MinColumnWidth:     5.0,
		MaxColumnWidth:     50.0,
		SnapToGrid:         true,
		GridResolution:     1.0,
		AlignmentTolerance: 0.5,
		CollisionBuffer:    2.0,
	}
}

// Add default rules
func (mbe *MonthBoundaryEngine) addDefaultBoundaryRules() {
	// High priority tasks boundary rule
	mbe.boundaryRules = append(mbe.boundaryRules, BoundaryRule{
		Name:        "High Priority Boundary",
		Description: "Handle high priority tasks at month boundaries",
		Priority:    10,
		Condition: func(bar *IntegratedTaskBar, context *MonthBoundaryContext) bool {
			return bar.Priority >= 4 && bar.MonthBoundary
		},
		Action: func(bar *IntegratedTaskBar, context *MonthBoundaryContext) *BoundaryAction {
			return &BoundaryAction{
				SplitTask:         true,
				CreateContinuation: true,
				AdjustPosition:    true,
				NewX:              bar.StartX,
				NewY:              context.DayHeight * 0.1,
				VisualStyle: &BoundaryVisualStyle{
					BorderStyle:     BorderSolid,
					BorderColor:     "#FF0000",
					BorderWidth:     2.0,
					BackgroundColor: bar.Color,
					Opacity:         0.9,
					HighlightEnabled: true,
					HighlightColor:  "#FFFF00",
					HighlightWidth:  1.0,
				},
				Priority: 10,
			}
		},
	})
	
	// Milestone tasks boundary rule
	mbe.boundaryRules = append(mbe.boundaryRules, BoundaryRule{
		Name:        "Milestone Boundary",
		Description: "Handle milestone tasks at month boundaries",
		Priority:    8,
		Condition: func(bar *IntegratedTaskBar, context *MonthBoundaryContext) bool {
			return (bar.Category == "MILESTONE" || 
				   (bar.TaskName != "" && len(bar.TaskName) > 10 && 
				    bar.TaskName[:10] == "MILESTONE:")) && bar.MonthBoundary
		},
		Action: func(bar *IntegratedTaskBar, context *MonthBoundaryContext) *BoundaryAction {
			return &BoundaryAction{
				SplitTask:         false,
				CreateContinuation: true,
				AdjustPosition:    true,
				NewX:              bar.StartX,
				NewY:              context.DayHeight * 0.4,
				VisualStyle: &BoundaryVisualStyle{
					BorderStyle:     BorderDouble,
					BorderColor:     "#0000FF",
					BorderWidth:     3.0,
					BackgroundColor: bar.Color,
					Opacity:         1.0,
					ShadowEnabled:   true,
					ShadowColor:     "#000000",
					ShadowBlur:      2.0,
					ShadowOffsetX:   1.0,
					ShadowOffsetY:   1.0,
				},
				Priority: 8,
			}
		},
	})
	
	// Default boundary rule
	mbe.boundaryRules = append(mbe.boundaryRules, BoundaryRule{
		Name:        "Default Boundary",
		Description: "Default handling for tasks at month boundaries",
		Priority:    1,
		Condition: func(bar *IntegratedTaskBar, context *MonthBoundaryContext) bool {
			return bar.MonthBoundary
		},
		Action: func(bar *IntegratedTaskBar, context *MonthBoundaryContext) *BoundaryAction {
			return &BoundaryAction{
				SplitTask:         true,
				CreateContinuation: false,
				AdjustPosition:    false,
				VisualStyle: &BoundaryVisualStyle{
					BorderStyle:     BorderDashed,
					BorderColor:     "#666666",
					BorderWidth:     1.0,
					BackgroundColor: bar.Color,
					Opacity:         0.8,
				},
				Priority: 1,
			}
		},
	})
}

func (mbe *MonthBoundaryEngine) addDefaultTransitionRules() {
	// Smooth transition rule
	mbe.transitionRules = append(mbe.transitionRules, TransitionRule{
		Name:        "Smooth Transition",
		Description: "Apply smooth transitions for high priority tasks",
		Priority:    10,
		Condition: func(bar *IntegratedTaskBar, context *MonthBoundaryContext) bool {
			return bar.Priority >= 4 && bar.MonthBoundary
		},
		Action: func(bar *IntegratedTaskBar, context *MonthBoundaryContext) *TransitionAction {
			return &TransitionAction{
				SmoothTransition: true,
				Duration:         500 * time.Millisecond,
				EasingFunction:   EasingEaseInOut,
				VisualEffects:    []VisualEffect{EffectFadeIn, EffectSlideRight},
				Priority:         10,
			}
		},
	})
	
	// Default transition rule
	mbe.transitionRules = append(mbe.transitionRules, TransitionRule{
		Name:        "Default Transition",
		Description: "Default transition for all tasks",
		Priority:    1,
		Condition: func(bar *IntegratedTaskBar, context *MonthBoundaryContext) bool {
			return bar.MonthBoundary
		},
		Action: func(bar *IntegratedTaskBar, context *MonthBoundaryContext) *TransitionAction {
			return &TransitionAction{
				SmoothTransition: false,
				Duration:         200 * time.Millisecond,
				EasingFunction:   EasingLinear,
				VisualEffects:    []VisualEffect{EffectFadeIn},
				Priority:         1,
			}
		},
	})
}

func (mbe *MonthBoundaryEngine) addDefaultContinuityRules() {
	// Visual continuity rule
	mbe.continuityRules = append(mbe.continuityRules, ContinuityRule{
		Name:        "Visual Continuity",
		Description: "Maintain visual continuity across months",
		Priority:    10,
		Condition: func(bars []*IntegratedTaskBar, context *MonthBoundaryContext) bool {
			// Check if there are tasks that cross month boundaries
			for _, bar := range bars {
				if bar.MonthBoundary {
					return true
				}
			}
			return false
		},
		Action: func(bars []*IntegratedTaskBar, context *MonthBoundaryContext) *ContinuityAction {
			return &ContinuityAction{
				MaintainAlignment: true,
				PreserveSpacing:   true,
				ConsistentColors:  true,
				ConsistentSizes:   true,
				VisualConnections: []VisualConnection{
					{
						ConnectionType: ConnectionArrow,
						LineStyle:     LineSolid,
						LineColor:     "#666666",
						LineWidth:     1.0,
						ArrowStyle:    ArrowSimple,
						Label:         "Continues",
						Priority:      10,
					},
				},
				Priority: 10,
			}
		},
	})
}
