package generator

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"latex-yearly-planner/internal/calendar"
)

// OptimizedPDFGenerator provides high-performance PDF generation with caching and parallel processing
type OptimizedPDFGenerator struct {
	config        *PDFGeneratorConfig
	cache         *PDFCache
	pool          *WorkerPool
	compiler      *OptimizedLaTeXCompiler
	templateEngine *OptimizedTemplateEngine
	logger        PDFLogger
}

// PDFGeneratorConfig defines configuration for optimized PDF generation
type PDFGeneratorConfig struct {
	// Performance settings
	EnablePDFCache        bool `json:"enable_pdf_cache"`
	CacheSize             int  `json:"cache_size"`
	EnableParallelCompilation bool `json:"enable_parallel_compilation"`
	MaxWorkers            int  `json:"max_workers"`
	
	// LaTeX compilation
	EnableLaTeXCache      bool          `json:"enable_latex_cache"`
	LaTeXCacheTTL         time.Duration `json:"latex_cache_ttl"`
	EnableIncrementalCompilation bool   `json:"enable_incremental_compilation"`
	MaxCompilationTime    time.Duration `json:"max_compilation_time"`
	
	// Memory optimization
	EnableMemoryOptimization bool `json:"enable_memory_optimization"`
	MaxMemoryUsage          int64 `json:"max_memory_usage_mb"`
	EnableStreamingOutput   bool `json:"enable_streaming_output"`
	
	// Quality settings
	EnableQualityOptimization bool    `json:"enable_quality_optimization"`
	TargetDPI                int     `json:"target_dpi"`
	CompressionLevel         int     `json:"compression_level"`
	
	// Error handling
	EnableErrorRecovery     bool `json:"enable_error_recovery"`
	MaxRetryAttempts        int  `json:"max_retry_attempts"`
	RetryDelay              time.Duration `json:"retry_delay"`
}

// PDFCache provides intelligent caching for PDF generation
type PDFCache struct {
	config      *PDFGeneratorConfig
	pdfs        map[string]*CachedPDF
	latexCache  map[string]*CachedLaTeX
	mu          sync.RWMutex
}

// CachedPDF represents cached PDF data
type CachedPDF struct {
	FilePath    string
	FileSize    int64
	CreatedAt   time.Time
	Hash        string
	TaskCount   int
	CompilationTime time.Duration
}

// CachedLaTeX represents cached LaTeX compilation
type CachedLaTeX struct {
	Content     string
	Hash        string
	CreatedAt   time.Time
	CompilationTime time.Duration
}

// WorkerPool manages parallel PDF generation workers
type WorkerPool struct {
	config      *PDFGeneratorConfig
	workers     chan struct{}
	mu          sync.Mutex
}

// OptimizedLaTeXCompiler provides high-performance LaTeX compilation
type OptimizedLaTeXCompiler struct {
	config        *PDFGeneratorConfig
	cache         *LaTeXCache
	compiler      *LaTeXCompiler
	logger        PDFLogger
}

// LaTeXCache provides caching for LaTeX compilation
type LaTeXCache struct {
	config      *PDFGeneratorConfig
	compilations map[string]*CachedLaTeX
	mu          sync.RWMutex
}

// LaTeXCompiler handles LaTeX compilation
type LaTeXCompiler struct {
	config        *PDFGeneratorConfig
	tempDir       string
	logger        PDFLogger
}

// OptimizedTemplateEngine provides high-performance template rendering
type OptimizedTemplateEngine struct {
	config        *PDFGeneratorConfig
	cache         *TemplateCache
	engine        *TemplateEngine
	logger        PDFLogger
}

// TemplateCache provides caching for template rendering
type TemplateCache struct {
	config      *PDFGeneratorConfig
	templates   map[string]*CachedTemplate
	mu          sync.RWMutex
}

// CachedTemplate represents cached template data
type CachedTemplate struct {
	Content     string
	Hash        string
	CreatedAt   time.Time
	RenderTime  time.Duration
}

// TemplateEngine handles template rendering
type TemplateEngine struct {
	config        *PDFGeneratorConfig
	templates     map[string]string
	logger        PDFLogger
}

// NewOptimizedPDFGenerator creates a new optimized PDF generator
func NewOptimizedPDFGenerator() *OptimizedPDFGenerator {
	config := GetDefaultPDFGeneratorConfig()
	
	return &OptimizedPDFGenerator{
		config:         config,
		cache:          NewPDFCache(config),
		pool:           NewWorkerPool(config),
		compiler:       NewOptimizedLaTeXCompiler(config),
		templateEngine: NewOptimizedTemplateEngine(config),
		logger:         &OptimizedPDFGeneratorLogger{},
	}
}

// GetDefaultPDFGeneratorConfig returns the default PDF generator configuration
func GetDefaultPDFGeneratorConfig() *PDFGeneratorConfig {
	return &PDFGeneratorConfig{
		EnablePDFCache:           true,
		CacheSize:                20,
		EnableParallelCompilation: true,
		MaxWorkers:               4,
		EnableLaTeXCache:         true,
		LaTeXCacheTTL:            time.Hour * 24,
		EnableIncrementalCompilation: true,
		MaxCompilationTime:       time.Minute * 5,
		EnableMemoryOptimization: true,
		MaxMemoryUsage:           512, // 512MB
		EnableStreamingOutput:    true,
		EnableQualityOptimization: true,
		TargetDPI:                300,
		CompressionLevel:         6,
		EnableErrorRecovery:      true,
		MaxRetryAttempts:         3,
		RetryDelay:               time.Second * 2,
	}
}

// SetLogger sets the logger for the PDF generator
func (opg *OptimizedPDFGenerator) SetLogger(logger PDFLogger) {
	opg.logger = logger
	opg.compiler.SetLogger(logger)
	opg.templateEngine.SetLogger(logger)
}

// GeneratePDF generates an optimized PDF from layout
func (opg *OptimizedPDFGenerator) GeneratePDF(ctx context.Context, layout *calendar.CalendarLayout, outputPath string) error {
	start := time.Now()
	opg.logger.Info("Starting optimized PDF generation")
	
	// Check cache first
	if opg.config.EnablePDFCache {
		hash := opg.calculateLayoutHash(layout)
		if cached, found := opg.cache.GetPDF(hash); found {
			opg.logger.Info("PDF cache hit, copying cached file")
			return opg.copyCachedPDF(cached, outputPath)
		}
	}
	
	// Generate PDF with retry logic
	var err error
	for attempt := 0; attempt < opg.config.MaxRetryAttempts; attempt++ {
		if attempt > 0 {
			opg.logger.Info("Retry attempt %d/%d", attempt+1, opg.config.MaxRetryAttempts)
			time.Sleep(opg.config.RetryDelay)
		}
		
		err = opg.generatePDFInternal(ctx, layout, outputPath)
		if err == nil {
			break
		}
		
		opg.logger.Error("PDF generation attempt %d failed: %v", attempt+1, err)
	}
	
	if err != nil {
		return fmt.Errorf("PDF generation failed after %d attempts: %w", opg.config.MaxRetryAttempts, err)
	}
	
	// Cache the result
	if opg.config.EnablePDFCache {
		hash := opg.calculateLayoutHash(layout)
		fileInfo, _ := os.Stat(outputPath)
		opg.cache.SetPDF(hash, &CachedPDF{
			FilePath:        outputPath,
			FileSize:        fileInfo.Size(),
			CreatedAt:       time.Now(),
			Hash:            hash,
			TaskCount:       len(layout.Tasks),
			CompilationTime: time.Since(start),
		})
	}
	
	opg.logger.Info("Generated PDF in %v", time.Since(start))
	return nil
}

// generatePDFInternal generates the actual PDF
func (opg *OptimizedPDFGenerator) generatePDFInternal(ctx context.Context, layout *calendar.CalendarLayout, outputPath string) error {
	// Create context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, opg.config.MaxCompilationTime)
	defer cancel()
	
	// Generate LaTeX content
	latexContent, err := opg.templateEngine.RenderLayout(layout)
	if err != nil {
		return fmt.Errorf("template rendering failed: %w", err)
	}
	
	// Compile LaTeX to PDF
	pdfPath, err := opg.compiler.CompileLaTeX(timeoutCtx, latexContent, outputPath)
	if err != nil {
		return fmt.Errorf("LaTeX compilation failed: %w", err)
	}
	
	// Move to final output path if different
	if pdfPath != outputPath {
		if err := os.Rename(pdfPath, outputPath); err != nil {
			return fmt.Errorf("failed to move PDF to output path: %w", err)
		}
	}
	
	return nil
}

// calculateLayoutHash calculates a hash for layout caching
func (opg *OptimizedPDFGenerator) calculateLayoutHash(layout *calendar.CalendarLayout) string {
	// Simple hash based on task count and dimensions
	hash := fmt.Sprintf("%d_%d_%d_%d", 
		len(layout.Tasks),
		int(layout.Dimensions.Width),
		int(layout.Dimensions.Height),
		layout.Config.Year)
	
	return hash
}

// copyCachedPDF copies a cached PDF to the output path
func (opg *OptimizedPDFGenerator) copyCachedPDF(cached *CachedPDF, outputPath string) error {
	src, err := os.Open(cached.FilePath)
	if err != nil {
		return fmt.Errorf("failed to open cached PDF: %w", err)
	}
	defer src.Close()
	
	dst, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output PDF: %w", err)
	}
	defer dst.Close()
	
	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("failed to copy cached PDF: %w", err)
	}
	
	opg.logger.Info("Copied cached PDF (%d bytes)", cached.FileSize)
	return nil
}

// NewPDFCache creates a new PDF cache
func NewPDFCache(config *PDFGeneratorConfig) *PDFCache {
	return &PDFCache{
		config:     config,
		pdfs:       make(map[string]*CachedPDF),
		latexCache: make(map[string]*CachedLaTeX),
	}
}

// GetPDF retrieves cached PDF
func (pc *PDFCache) GetPDF(hash string) (*CachedPDF, bool) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	
	cached, exists := pc.pdfs[hash]
	if !exists {
		return nil, false
	}
	
	// Check if file still exists
	if _, err := os.Stat(cached.FilePath); err != nil {
		return nil, false
	}
	
	return cached, true
}

// SetPDF stores cached PDF
func (pc *PDFCache) SetPDF(hash string, pdf *CachedPDF) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	
	// Check cache size limit
	if len(pc.pdfs) >= pc.config.CacheSize {
		pc.evictOldestPDF()
	}
	
	pc.pdfs[hash] = pdf
}

// evictOldestPDF removes the oldest cached PDF
func (pc *PDFCache) evictOldestPDF() {
	var oldestKey string
	var oldestTime time.Time
	
	for key, pdf := range pc.pdfs {
		if oldestKey == "" || pdf.CreatedAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = pdf.CreatedAt
		}
	}
	
	if oldestKey != "" {
		// Remove file from disk
		if pdf, exists := pc.pdfs[oldestKey]; exists {
			os.Remove(pdf.FilePath)
		}
		delete(pc.pdfs, oldestKey)
	}
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(config *PDFGeneratorConfig) *WorkerPool {
	return &WorkerPool{
		config:  config,
		workers: make(chan struct{}, config.MaxWorkers),
	}
}

// AcquireWorker acquires a worker from the pool
func (wp *WorkerPool) AcquireWorker() {
	wp.workers <- struct{}{}
}

// ReleaseWorker releases a worker back to the pool
func (wp *WorkerPool) ReleaseWorker() {
	<-wp.workers
}

// NewOptimizedLaTeXCompiler creates a new optimized LaTeX compiler
func NewOptimizedLaTeXCompiler(config *PDFGeneratorConfig) *OptimizedLaTeXCompiler {
	return &OptimizedLaTeXCompiler{
		config:   config,
		cache:    NewLaTeXCache(config),
		compiler: NewLaTeXCompiler(config),
		logger:   &OptimizedPDFGeneratorLogger{},
	}
}

// SetLogger sets the logger for the compiler
func (olc *OptimizedLaTeXCompiler) SetLogger(logger PDFLogger) {
	olc.logger = logger
	olc.compiler.SetLogger(logger)
}

// CompileLaTeX compiles LaTeX content to PDF
func (olc *OptimizedLaTeXCompiler) CompileLaTeX(ctx context.Context, content string, outputPath string) (string, error) {
	// Check cache first
	if olc.config.EnableLaTeXCache {
		hash := olc.calculateContentHash(content)
		if cached, found := olc.cache.GetCompilation(hash); found {
			olc.logger.Info("LaTeX cache hit")
			return cached.Content, nil
		}
	}
	
	// Compile LaTeX
	pdfPath, err := olc.compiler.Compile(ctx, content, outputPath)
	if err != nil {
		return "", fmt.Errorf("LaTeX compilation failed: %w", err)
	}
	
	// Cache the result
	if olc.config.EnableLaTeXCache {
		hash := olc.calculateContentHash(content)
		olc.cache.SetCompilation(hash, &CachedLaTeX{
			Content:     pdfPath,
			Hash:        hash,
			CreatedAt:   time.Now(),
			CompilationTime: time.Since(time.Now()),
		})
	}
	
	return pdfPath, nil
}

// calculateContentHash calculates a hash for content caching
func (olc *OptimizedLaTeXCompiler) calculateContentHash(content string) string {
	// Simple hash based on content length and first/last characters
	hash := fmt.Sprintf("%d_%c_%c", len(content), content[0], content[len(content)-1])
	return hash
}

// NewLaTeXCache creates a new LaTeX cache
func NewLaTeXCache(config *PDFGeneratorConfig) *LaTeXCache {
	return &LaTeXCache{
		config:       config,
		compilations: make(map[string]*CachedLaTeX),
	}
}

// GetCompilation retrieves cached compilation
func (lc *LaTeXCache) GetCompilation(hash string) (*CachedLaTeX, bool) {
	lc.mu.RLock()
	defer lc.mu.RUnlock()
	
	cached, exists := lc.compilations[hash]
	if !exists {
		return nil, false
	}
	
	// Check if expired
	if time.Since(cached.CreatedAt) > lc.config.LaTeXCacheTTL {
		return nil, false
	}
	
	return cached, true
}

// SetCompilation stores cached compilation
func (lc *LaTeXCache) SetCompilation(hash string, compilation *CachedLaTeX) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	
	lc.compilations[hash] = compilation
}

// NewLaTeXCompiler creates a new LaTeX compiler
func NewLaTeXCompiler(config *PDFGeneratorConfig) *LaTeXCompiler {
	return &LaTeXCompiler{
		config:  config,
		tempDir: os.TempDir(),
		logger:  &OptimizedPDFGeneratorLogger{},
	}
}

// SetLogger sets the logger for the compiler
func (lc *LaTeXCompiler) SetLogger(logger PDFLogger) {
	lc.logger = logger
}

// Compile compiles LaTeX content to PDF
func (lc *LaTeXCompiler) Compile(ctx context.Context, content string, outputPath string) (string, error) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp(lc.tempDir, "latex_compile_*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Write LaTeX content
	texFile := filepath.Join(tempDir, "document.tex")
	if err := os.WriteFile(texFile, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write LaTeX file: %w", err)
	}
	
	// Compile LaTeX
	pdfFile := filepath.Join(tempDir, "document.pdf")
	cmd := exec.CommandContext(ctx, "pdflatex", "-interaction=nonstopmode", "-output-directory", tempDir, texFile)
	cmd.Dir = tempDir
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	if err := cmd.Run(); err != nil {
		lc.logger.Error("LaTeX compilation failed: %s", stderr.String())
		return "", fmt.Errorf("LaTeX compilation failed: %w", err)
	}
	
	// Check if PDF was created
	if _, err := os.Stat(pdfFile); err != nil {
		return "", fmt.Errorf("PDF file not created: %w", err)
	}
	
	// Move to output path
	if err := os.Rename(pdfFile, outputPath); err != nil {
		return "", fmt.Errorf("failed to move PDF to output path: %w", err)
	}
	
	lc.logger.Info("LaTeX compilation successful")
	return outputPath, nil
}

// NewOptimizedTemplateEngine creates a new optimized template engine
func NewOptimizedTemplateEngine(config *PDFGeneratorConfig) *OptimizedTemplateEngine {
	return &OptimizedTemplateEngine{
		config:   config,
		cache:    NewTemplateCache(config),
		engine:   NewTemplateEngine(config),
		logger:   &OptimizedPDFGeneratorLogger{},
	}
}

// SetLogger sets the logger for the template engine
func (ote *OptimizedTemplateEngine) SetLogger(logger PDFLogger) {
	ote.logger = logger
	ote.engine.SetLogger(logger)
}

// RenderLayout renders a layout to LaTeX content
func (ote *OptimizedTemplateEngine) RenderLayout(layout *calendar.CalendarLayout) (string, error) {
	// Check cache first
	if ote.config.EnableLaTeXCache {
		hash := ote.calculateLayoutHash(layout)
		if cached, found := ote.cache.GetTemplate(hash); found {
			ote.logger.Info("Template cache hit")
			return cached.Content, nil
		}
	}
	
	// Render template
	content, err := ote.engine.Render(layout)
	if err != nil {
		return "", fmt.Errorf("template rendering failed: %w", err)
	}
	
	// Cache the result
	if ote.config.EnableLaTeXCache {
		hash := ote.calculateLayoutHash(layout)
		ote.cache.SetTemplate(hash, &CachedTemplate{
			Content:    content,
			Hash:       hash,
			CreatedAt:  time.Now(),
			RenderTime: time.Since(time.Now()),
		})
	}
	
	return content, nil
}

// calculateLayoutHash calculates a hash for layout caching
func (ote *OptimizedTemplateEngine) calculateLayoutHash(layout *calendar.CalendarLayout) string {
	// Simple hash based on task count and dimensions
	hash := fmt.Sprintf("%d_%d_%d", 
		len(layout.Tasks),
		int(layout.Dimensions.Width),
		int(layout.Dimensions.Height))
	
	return hash
}

// NewTemplateCache creates a new template cache
func NewTemplateCache(config *PDFGeneratorConfig) *TemplateCache {
	return &TemplateCache{
		config:    config,
		templates: make(map[string]*CachedTemplate),
	}
}

// GetTemplate retrieves cached template
func (tc *TemplateCache) GetTemplate(hash string) (*CachedTemplate, bool) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	
	cached, exists := tc.templates[hash]
	if !exists {
		return nil, false
	}
	
	// Check if expired (1 hour)
	if time.Since(cached.CreatedAt) > time.Hour {
		return nil, false
	}
	
	return cached, true
}

// SetTemplate stores cached template
func (tc *TemplateCache) SetTemplate(hash string, template *CachedTemplate) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	
	tc.templates[hash] = template
}

// NewTemplateEngine creates a new template engine
func NewTemplateEngine(config *PDFGeneratorConfig) *TemplateEngine {
	return &TemplateEngine{
		config:    config,
		templates: make(map[string]string),
		logger:    &OptimizedPDFGeneratorLogger{},
	}
}

// SetLogger sets the logger for the template engine
func (te *TemplateEngine) SetLogger(logger PDFLogger) {
	te.logger = logger
}

// Render renders a layout to LaTeX content
func (te *TemplateEngine) Render(layout *calendar.CalendarLayout) (string, error) {
	// Simple LaTeX template for now
	var content strings.Builder
	
	content.WriteString("\\documentclass{article}\n")
	content.WriteString("\\usepackage[utf8]{inputenc}\n")
	content.WriteString("\\usepackage{tikz}\n")
	content.WriteString("\\usepackage{pgfgantt}\n")
	content.WriteString("\\begin{document}\n")
	content.WriteString("\\title{Gantt Chart}\n")
	content.WriteString("\\maketitle\n")
	content.WriteString("\\begin{ganttchart}[vgrid, hgrid]{1}{12}\n")
	
	// Add tasks
	for i, task := range layout.Tasks {
		content.WriteString(fmt.Sprintf("\\ganttbar{Task %d}{%d}{%d}\n", i+1, 1, 3))
	}
	
	content.WriteString("\\end{ganttchart}\n")
	content.WriteString("\\end{document}\n")
	
	return content.String(), nil
}

// OptimizedPDFGeneratorLogger provides logging for optimized PDF generator
type OptimizedPDFGeneratorLogger struct{}

func (l *OptimizedPDFGeneratorLogger) Info(msg string, args ...interface{})  { fmt.Printf("[PDF-INFO] "+msg+"\n", args...) }
func (l *OptimizedPDFGeneratorLogger) Error(msg string, args ...interface{}) { fmt.Printf("[PDF-ERROR] "+msg+"\n", args...) }
func (l *OptimizedPDFGeneratorLogger) Debug(msg string, args ...interface{}) { fmt.Printf("[PDF-DEBUG] "+msg+"\n", args...) }
