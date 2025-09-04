#!/usr/bin/env python3
"""
Demonstration of the Ad-Hoc Debugging Workflow

This script demonstrates the complete APM debugging workflow from issue creation
to resolution, showing how the system handles both simple and complex issues.
"""

import sys
import time
from pathlib import Path

# Add src to path for imports
sys.path.insert(0, str(Path(__file__).parent / "src"))

from src.debug_cli import DebugCLI


def demo_simple_issue():
    """* Demonstrate handling a simple issue."""
    print("🔧 DEMO: Simple Issue Resolution")
    print("=" * 50)
    
    cli = DebugCLI()
    
    # * Create a simple issue
    class SimpleArgs:
        def __init__(self):
            self.issue_id = "demo_simple_001"
            self.title = "Simple Import Error"
            self.description = "Module 'pandas' not found during execution"
            self.error_context = None
            self.verbose = False
            self.quiet = False
    
    args = SimpleArgs()
    
    print("1. Creating debug session for simple issue...")
    result = cli.create_session(args)
    print(f"   Result: {'Success' if result == 0 else 'Failed'}")
    
    print("\n2. Listing current sessions...")
    class ListArgs:
        def __init__(self):
            self.status = None
            self.limit = 5
            self.offset = 0
            self.verbose = False
            self.quiet = False
    
    list_args = ListArgs()
    cli.list_sessions(list_args)
    
    return result == 0


def demo_complex_issue():
    """* Demonstrate handling a complex issue requiring delegation."""
    print("\n🔧 DEMO: Complex Issue Delegation")
    print("=" * 50)
    
    cli = DebugCLI()
    
    # * Create a complex issue
    class ComplexArgs:
        def __init__(self):
            self.issue_id = "demo_complex_001"
            self.title = "Performance Degradation"
            self.description = """
            The LaTeX generation process has become significantly slower over time.
            Processing 1000+ tasks now takes 5+ minutes instead of the usual 30 seconds.
            Memory usage is also increasing linearly with task count.
            Suspected memory leak or inefficient algorithm in the template generation.
            """
            self.error_context = '{"error_type": "PerformanceIssue", "execution_time": "5+ minutes", "memory_usage": "2GB+"}'
            self.save_prompt = True
            self.verbose = False
            self.quiet = False
    
    args = ComplexArgs()
    
    print("1. Handling complex issue through APM workflow...")
    result = cli.handle_issue(args)
    print(f"   Result: {'Success' if result == 0 else 'Requires Ad-Hoc Session' if result == 2 else 'Failed'}")
    
    if result == 2:
        print("\n2. Issue requires Ad-Hoc debugging session")
        print("   Delegation prompt has been generated and saved")
        
        # * Simulate Ad-Hoc session completion
        print("\n3. Simulating Ad-Hoc session completion...")
        
        # * Get the session ID from the last created session
        sessions = cli.session_manager.list_sessions(limit=1)
        if sessions:
            session_id = sessions[0]['session_id']
            
            class CompleteArgs:
                def __init__(self, session_id):
                    self.session_id = session_id
                    self.solution = """
                    Root Cause Analysis:
                    The performance issue was caused by inefficient string concatenation in the template generation loop.
                    
                    Solution:
                    1. Replaced string concatenation with list.join() method
                    2. Added caching for frequently used template fragments
                    3. Optimized the TikZ generation algorithm
                    4. Added memory monitoring and cleanup
                    
                    Testing:
                    - Verified processing time reduced from 5+ minutes to 45 seconds
                    - Memory usage stabilized at 500MB for 1000+ tasks
                    - No regression in output quality
                    
                    Prevention:
                    - Added performance benchmarks to CI/CD pipeline
                    - Implemented memory usage monitoring
                    - Added automated performance regression tests
                    """
                    self.ad_hoc_session_id = "ad_hoc_demo_001"
            
            complete_args = CompleteArgs(session_id)
            complete_result = cli.complete_workflow(complete_args)
            print(f"   Workflow completion: {'Success' if complete_result == 0 else 'Failed'}")
    
    return result == 2


def demo_statistics():
    """* Demonstrate the statistics and reporting features."""
    print("\n📊 DEMO: Statistics and Reporting")
    print("=" * 50)
    
    cli = DebugCLI()
    
    class StatsArgs:
        def __init__(self):
            self.verbose = False
            self.quiet = False
    
    stats_args = StatsArgs()
    
    print("1. Current debug session statistics:")
    cli.show_statistics(stats_args)
    
    print("\n2. Recent sessions:")
    class ListArgs:
        def __init__(self):
            self.status = None
            self.limit = 3
            self.offset = 0
            self.verbose = False
            self.quiet = False
    
    list_args = ListArgs()
    cli.list_sessions(list_args)
    
    return True


def demo_export_functionality():
    """* Demonstrate session export functionality."""
    print("\n📄 DEMO: Session Export")
    print("=" * 50)
    
    cli = DebugCLI()
    
    # * Get the most recent session
    sessions = cli.session_manager.list_sessions(limit=1)
    if not sessions:
        print("No sessions available for export")
        return False
    
    session_id = sessions[0]['session_id']
    
    print(f"1. Exporting session {session_id} as markdown...")
    
    class ExportArgs:
        def __init__(self, session_id):
            self.session_id = session_id
            self.format = "markdown"
            self.output = f"exported_session_{session_id[:8]}.md"
    
    export_args = ExportArgs(session_id)
    result = cli.export_session(export_args)
    
    if result == 0:
        print(f"   Export successful: {export_args.output}")
        
        # * Show a preview of the exported content
        try:
            with open(export_args.output, 'r') as f:
                content = f.read()
                preview = content[:300] + "..." if len(content) > 300 else content
                print(f"\n2. Export preview:")
                print("-" * 30)
                print(preview)
        except Exception as e:
            print(f"   Could not read exported file: {e}")
    
    return result == 0


def main():
    """* Run the complete demonstration."""
    print("🎯 Ad-Hoc Debugging System Demonstration")
    print("=" * 60)
    print("This demo shows the complete APM debugging workflow")
    print("from issue creation to resolution and reporting.")
    print("=" * 60)
    
    demos = [
        ("Simple Issue Resolution", demo_simple_issue),
        ("Complex Issue Delegation", demo_complex_issue),
        ("Statistics and Reporting", demo_statistics),
        ("Session Export", demo_export_functionality)
    ]
    
    results = []
    
    for demo_name, demo_func in demos:
        try:
            print(f"\n🚀 Starting: {demo_name}")
            result = demo_func()
            results.append((demo_name, result))
            print(f"✅ Completed: {demo_name}")
        except Exception as e:
            print(f"❌ Failed: {demo_name} - {e}")
            results.append((demo_name, False))
        
        time.sleep(1)  # * Brief pause between demos
    
    # * Summary
    print("\n" + "=" * 60)
    print("DEMONSTRATION SUMMARY")
    print("=" * 60)
    
    passed = 0
    total = len(results)
    
    for demo_name, result in results:
        status = "✅ PASS" if result else "❌ FAIL"
        print(f"{demo_name:<30} {status}")
        if result:
            passed += 1
    
    print(f"\nOverall: {passed}/{total} demos completed successfully")
    
    if passed == total:
        print("\n🎉 All demonstrations completed successfully!")
        print("The Ad-Hoc debugging system is working correctly.")
    else:
        print(f"\n⚠️ {total - passed} demonstrations failed.")
        print("Please review the implementation.")
    
    print("\n📚 Next Steps:")
    print("- Use 'python main.py debug --help' to see all available commands")
    print("- Check the DEBUG_SYSTEM_README.md for detailed documentation")
    print("- Run 'python test_debug_system.py' for comprehensive testing")
    
    return 0 if passed == total else 1


if __name__ == "__main__":
    sys.exit(main())
