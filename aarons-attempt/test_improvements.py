#!/usr/bin/env python3
"""
Test script to demonstrate the improvements made to the Gantt chart generator.
"""

import sys
import os
from pathlib import Path

# Add src to path
sys.path.insert(0, str(Path(__file__).parent / "src"))

from src.data_processor import DataProcessor
from src.export_system import ExportSystem
from src.interactive_generator import EnhancedTemplateGenerator


def test_enhanced_data_processing():
    """Test the enhanced data processing capabilities."""
    print("🧪 Testing Enhanced Data Processing...")
    
    processor = DataProcessor()
    
    # Test with the sample data
    csv_file = "../input/data.cleaned.csv"
    if not Path(csv_file).exists():
        print(f"❌ Test data file not found: {csv_file}")
        return False
    
    try:
        timeline = processor.process_csv_to_timeline(csv_file, "Test Timeline")
        print(f"✅ Successfully processed {len(timeline.tasks)} tasks")
        print(f"   Timeline: {timeline.start_date} to {timeline.end_date}")
        print(f"   Duration: {timeline.total_duration_days} days")
        return True
    except Exception as e:
        print(f"❌ Data processing failed: {e}")
        return False


def test_export_system():
    """Test the new export system."""
    print("\n🧪 Testing Export System...")
    
    # First process the data
    processor = DataProcessor()
    csv_file = "../input/data.cleaned.csv"
    
    if not Path(csv_file).exists():
        print(f"❌ Test data file not found: {csv_file}")
        return False
    
    try:
        timeline = processor.process_csv_to_timeline(csv_file, "Export Test")
        export_system = ExportSystem()
        
        # Test HTML export (doesn't require external tools)
        html_path = "test_output.html"
        success = export_system.export_to_html(timeline, html_path)
        
        if success:
            print(f"✅ Successfully exported HTML: {html_path}")
            return True
        else:
            print("❌ HTML export failed")
            return False
            
    except Exception as e:
        print(f"❌ Export test failed: {e}")
        return False


def test_enhanced_generator():
    """Test the enhanced template generator."""
    print("\n🧪 Testing Enhanced Template Generator...")
    
    try:
        generator = EnhancedTemplateGenerator()
        print("✅ Enhanced template generator initialized successfully")
        
        # Test document generation (without full compilation)
        processor = DataProcessor()
        csv_file = "../input/data.cleaned.csv"
        
        if Path(csv_file).exists():
            timeline = processor.process_csv_to_timeline(csv_file, "Generator Test")
            
            # Generate document content
            content = generator.generate_enhanced_document(timeline)
            
            if content and len(content) > 1000:  # Basic content validation
                print("✅ Enhanced document generation successful")
                print(f"   Generated {len(content)} characters of LaTeX content")
                return True
            else:
                print("❌ Generated content seems too short")
                return False
        else:
            print("⚠️  Skipping full test - data file not found")
            return True
            
    except Exception as e:
        print(f"❌ Enhanced generator test failed: {e}")
        return False


def main():
    """Run all improvement tests."""
    print("🚀 Testing Gantt Chart Generator Improvements")
    print("=" * 50)
    
    tests = [
        test_enhanced_data_processing,
        test_export_system,
        test_enhanced_generator
    ]
    
    results = []
    for test in tests:
        try:
            result = test()
            results.append(result)
        except Exception as e:
            print(f"❌ Test {test.__name__} crashed: {e}")
            results.append(False)
    
    print("\n" + "=" * 50)
    print("📊 Test Results Summary:")
    
    passed = sum(results)
    total = len(results)
    
    for i, (test, result) in enumerate(zip(tests, results)):
        status = "✅ PASS" if result else "❌ FAIL"
        print(f"  {i+1}. {test.__name__}: {status}")
    
    print(f"\nOverall: {passed}/{total} tests passed")
    
    if passed == total:
        print("🎉 All improvements are working correctly!")
        return 0
    else:
        print("⚠️  Some improvements need attention")
        return 1


if __name__ == "__main__":
    sys.exit(main())
