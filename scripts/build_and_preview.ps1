# Build PDF and generate preview images
param(
    [int]$Pages = 3
)

Write-Host "Building PDF..." -ForegroundColor Cyan

# Build
$env:PLANNER_CSV_FILE = "input_data\research_timeline_v5_comprehensive.csv"
go build -o generated\plannergen.exe .\cmd\planner
if ($LASTEXITCODE -ne 0) { exit 1 }

# Generate LaTeX
.\generated\plannergen.exe --config "src\core\base.yaml,src\core\monthly_calendar.yaml" --outdir generated
if ($LASTEXITCODE -ne 0) { exit 1 }

# Compile PDF
Push-Location generated
pdflatex -interaction=nonstopmode monthly_calendar.tex > monthly_calendar.log 2>&1
Pop-Location

if (-not (Test-Path "generated\monthly_calendar.pdf")) {
    Write-Host "ERROR: PDF not created" -ForegroundColor Red
    exit 1
}

Write-Host "PDF created successfully" -ForegroundColor Green

# Create preview directory
$previewDir = "generated\preview"
if (Test-Path $previewDir) {
    Remove-Item -Recurse -Force $previewDir
}
New-Item -ItemType Directory -Path $previewDir | Out-Null

# Try to convert to PNG using Python
Write-Host "Generating preview images..." -ForegroundColor Cyan

$converted = $false
try {
    python scripts\pdf_to_images.py "generated\monthly_calendar.pdf" $previewDir $Pages 150
    if ($LASTEXITCODE -eq 0 -and (Test-Path "$previewDir\page-01.png")) {
        $converted = $true
    }
} catch {
    Write-Host "Python conversion failed: $_" -ForegroundColor Yellow
}

if ($converted) {
    $count = (Get-ChildItem "$previewDir\*.png").Count
    Write-Host ""
    Write-Host "Created $count preview images in $previewDir" -ForegroundColor Green
    Write-Host ""
    Write-Host "Drag these images into the chat to share them:" -ForegroundColor Yellow
    Get-ChildItem "$previewDir\*.png" | ForEach-Object {
        Write-Host "  $($_.FullName)"
    }
} else {
    Write-Host ""
    Write-Host "Could not generate previews" -ForegroundColor Yellow
    Write-Host "Install poppler: choco install poppler" -ForegroundColor Gray
    Write-Host "Or download: https://github.com/oschwartz10612/poppler-windows/releases/" -ForegroundColor Gray
    Write-Host ""
    Write-Host "PDF is available at: generated\monthly_calendar.pdf" -ForegroundColor White
}
