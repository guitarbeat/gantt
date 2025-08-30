#!/usr/bin/env python3
"""
Simple HTML to PDF Converter using Chrome/Chromium
Converts HTML timeline files to high-quality PDFs
"""

import argparse
import os
import subprocess
import sys


def find_chrome():
    """Find Chrome or Chromium executable"""
    chrome_paths = [
        "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
        "/Applications/Chromium.app/Contents/MacOS/Chromium",
        "/usr/bin/google-chrome",
        "/usr/bin/chromium",
        "/usr/bin/chromium-browser"
    ]
    
    for path in chrome_paths:
        if os.path.exists(path):
            return path
    
    # Try to find in PATH
    try:
        result = subprocess.run(["which", "google-chrome"], 
                              capture_output=True, text=True, check=True)
        return result.stdout.strip()
    except subprocess.CalledProcessError:
        pass
    
    try:
        result = subprocess.run(["which", "chromium"], 
                              capture_output=True, text=True, check=True)
        return result.stdout.strip()
    except subprocess.CalledProcessError:
        pass
    
    return None


def convert_html_to_pdf_chrome(html_path: str, pdf_path: str, chrome_path: str):
    """
    Convert HTML to PDF using Chrome/Chromium headless mode
    """
    try:
        # Ensure input file exists
        if not os.path.exists(html_path):
            raise FileNotFoundError(f"HTML file not found: {html_path}")
        
        # Create output directory if it doesn't exist (only if path has a directory)
        output_dir = os.path.dirname(pdf_path)
        if output_dir and output_dir != ".":
            os.makedirs(output_dir, exist_ok=True)
        
        # Convert to absolute paths
        html_abs = os.path.abspath(html_path)
        pdf_abs = os.path.abspath(pdf_path)
        
        # Create file URL
        file_url = f"file://{html_abs}"
        
        # Build command as a single string to handle paths with spaces
        cmd_str = f'"{chrome_path}" --headless --disable-gpu --no-pdf-header-footer --print-to-pdf-no-header --print-to-pdf="{pdf_abs}" --run-all-compositor-stages-before-draw --disable-background-timer-throttling --disable-backgrounding-occluded-windows --disable-renderer-backgrounding "{file_url}"'
        
        print(f"Converting {html_path} to {pdf_path}...")
        print(f"Using Chrome: {chrome_path}")
        
        # Run Chrome using shell=True to handle paths with spaces
        result = subprocess.run(cmd_str, shell=True, capture_output=True, text=True, timeout=30)
        
        if result.returncode == 0 and os.path.exists(pdf_abs):
            print(f"✅ Successfully created: {pdf_path}")
            return True
        else:
            print(f"❌ Chrome conversion failed:")
            print(f"Return code: {result.returncode}")
            if result.stderr:
                print(f"Error: {result.stderr}")
            return False
            
    except subprocess.TimeoutExpired:
        print("❌ Conversion timed out")
        return False
    except Exception as e:
        print(f"❌ Error converting to PDF: {e}")
        return False


def main():
    parser = argparse.ArgumentParser(
        description="Convert HTML timeline to PDF using Chrome/Chromium"
    )
    parser.add_argument(
        "input", 
        help="Input HTML file path (default: Timeline.html)",
        nargs="?", 
        default="Timeline.html"
    )
    parser.add_argument(
        "-o", "--output",
        help="Output PDF file path (default: Timeline.pdf)",
        default="Timeline.pdf"
    )
    parser.add_argument(
        "--chrome-path",
        help="Path to Chrome/Chromium executable",
        default=None
    )
    
    args = parser.parse_args()
    
    # Find Chrome
    chrome_path = args.chrome_path or find_chrome()
    if not chrome_path:
        print("❌ Chrome/Chromium not found!")
        print("\nPlease install Chrome or Chromium, or specify the path with --chrome-path")
        print("\nInstallation options:")
        print("1. Download Chrome from: https://www.google.com/chrome/")
        print("2. Install Chromium via Homebrew: brew install --cask chromium")
        print("3. Or specify the path manually: --chrome-path /path/to/chrome")
        sys.exit(1)
    
    print(f"🔍 Found Chrome at: {chrome_path}")
    
    # Convert to PDF
    success = convert_html_to_pdf_chrome(args.input, args.output, chrome_path)
    
    if success:
        print(f"\n🎉 PDF created successfully!")
        print(f"📁 Location: {os.path.abspath(args.output)}")
        print(f"📏 File size: {os.path.getsize(args.output) / 1024:.1f} KB")
    else:
        sys.exit(1)


if __name__ == "__main__":
    main()
