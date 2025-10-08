#!/usr/bin/env python3
"""
Convert PDF pages to PNG images
Requires: pip install pdf2image pillow
"""

import sys
import os
from pathlib import Path

try:
    from pdf2image import convert_from_path
    from PIL import Image
except ImportError as e:
    print(f"ERROR: Missing required package: {e}")
    print("Install with: pip install pdf2image pillow")
    sys.exit(1)

def convert_pdf_to_images(pdf_path, output_dir, num_pages=3, dpi=150):
    """Convert first N pages of PDF to PNG images"""
    
    # Create output directory
    output_path = Path(output_dir)
    output_path.mkdir(parents=True, exist_ok=True)
    
    # Check if PDF exists
    if not Path(pdf_path).exists():
        print(f"ERROR: PDF not found: {pdf_path}")
        return False
    
    try:
        print(f"Converting {pdf_path} to images...")
        print(f"  Pages: {num_pages}")
        print(f"  DPI: {dpi}")
        print(f"  Output: {output_dir}")
        
        # Convert PDF to images
        images = convert_from_path(
            pdf_path,
            dpi=dpi,
            first_page=1,
            last_page=num_pages,
            fmt='png'
        )
        
        # Save images
        for i, image in enumerate(images, start=1):
            output_file = output_path / f"page-{i:02d}.png"
            image.save(output_file, 'PNG')
            size_kb = output_file.stat().st_size / 1024
            print(f"  [OK] Created: {output_file.name} ({size_kb:.1f} KB)")
        
        print(f"\n[OK] Successfully created {len(images)} preview images")
        return True
        
    except Exception as e:
        print(f"ERROR: Failed to convert PDF: {e}")
        print("\nTroubleshooting:")
        print("  1. Install poppler: https://github.com/oschwartz10612/poppler-windows/releases/")
        print("  2. Add poppler/bin to your PATH")
        print("  3. Or use: choco install poppler")
        return False

if __name__ == "__main__":
    if len(sys.argv) < 3:
        print("Usage: python pdf_to_images.py <pdf_file> <output_dir> [num_pages] [dpi]")
        print("Example: python pdf_to_images.py planner.pdf preview 3 150")
        sys.exit(1)
    
    pdf_file = sys.argv[1]
    output_dir = sys.argv[2]
    num_pages = int(sys.argv[3]) if len(sys.argv) > 3 else 3
    dpi = int(sys.argv[4]) if len(sys.argv) > 4 else 150
    
    success = convert_pdf_to_images(pdf_file, output_dir, num_pages, dpi)
    sys.exit(0 if success else 1)
