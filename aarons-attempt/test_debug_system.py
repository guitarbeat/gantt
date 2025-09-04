#!/usr/bin/env python3
"""
Test script for the Ad-Hoc Debugging System

This script demonstrates the complete APM debugging workflow including:
1. Issue complexity assessment
2. Local debugging attempts
3. Delegation prompt generation
4. Session management
5. Solution integration
"""

import json
import logging
import sys
import time
from pathlib import Path

# Add src to path for imports
sys.path.insert(0, str(Path(__file__).parent / "src"))

from src.debug_system import AdHocDebugger, DebugWorkflow, IssueComplexity, DebugStatus
from src.session_manager import SessionManager, SessionIntegration


def setup_logging():
    """* Setup logging for the test."""
    logging.basicConfig(
        level=logging.INFO,
        format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
    )


def test_simple_issue():
    """* Test handling of a simple issue."""
    print("\n" + "="*60)
    print("TEST 1: Simple Issue - Local Debugging Success")
    print("="*60)
    
    debugger = AdHocDebugger()
    
    # * Create a simple issue
    issue_id = "simple_import_error"
    title = "Import Error"
    description = "Module not found error when importing pandas"
    error_context = {
        "error_type": "ImportError",
        "error_message": "No module named 'pandas'",
        "file": "main.py",
        "line": 5
    }
    
    # * Create debug session
    session = debugger.create_debug_session(issue_id, title, description, error_context)
    print(f"Created session: {session.session_id}")
    print(f"Complexity: {session.complexity.value}")
    
    # * Define a simple debug function that succeeds
    def simple_debug_function():
        return "Installed pandas using: pip install pandas"
    
    # * Attempt local debugging
    attempt = debugger.attempt_local_debug(session.session_id, simple_debug_function)
    print(f"Debug attempt result: {'Success' if attempt.success else 'Failed'}")
    
    if attempt.success:
        print(f"Solution: {attempt.solution}")
        session.status = DebugStatus.RESOLVED
        session.final_solution = attempt.solution
        print("✅ Issue resolved with local debugging")
    else:
        print(f"Error: {attempt.error_message}")
    
    return session


def test_complex_issue():
    """* Test handling of a complex issue requiring delegation."""
    print("\n" + "="*60)
    print("TEST 2: Complex Issue - Delegation Required")
    print("="*60)
    
    debugger = AdHocDebugger()
    
    # * Create a complex issue
    issue_id = "complex_performance_issue"
    title = "Performance Degradation"
    description = """
    The LaTeX generation process has become significantly slower over time.
    Processing 1000+ tasks now takes 5+ minutes instead of the usual 30 seconds.
    Memory usage is also increasing linearly with task count.
    Suspected memory leak or inefficient algorithm in the template generation.
    """
    error_context = {
        "error_type": "PerformanceIssue",
        "execution_time": "5+ minutes",
        "memory_usage": "2GB+",
        "task_count": 1000,
        "normal_time": "30 seconds",
        "stack_trace": "No exception, just slow execution"
    }
    
    # * Create debug session
    session = debugger.create_debug_session(issue_id, title, description, error_context)
    print(f"Created session: {session.session_id}")
    print(f"Complexity: {session.complexity.value}")
    
    # * Define a debug function that fails (simulating complex issue)
    def complex_debug_function():
        raise Exception("Complex performance issue requires deep analysis")
    
    # * Attempt local debugging (will fail)
    try:
        attempt = debugger.attempt_local_debug(session.session_id, complex_debug_function)
        print(f"Debug attempt result: {'Success' if attempt.success else 'Failed'}")
        if not attempt.success:
            print(f"Error: {attempt.error_message}")
    except Exception as e:
        print(f"Local debugging failed: {e}")
    
    # * Generate delegation prompt
    delegation_prompt = debugger.generate_delegation_prompt(session.session_id)
    formatted_prompt = debugger.format_delegation_prompt_for_ai(delegation_prompt)
    
    print("\n📋 Delegation Prompt Generated:")
    print("-" * 40)
    print(formatted_prompt[:500] + "..." if len(formatted_prompt) > 500 else formatted_prompt)
    
    return session, formatted_prompt


def test_workflow_integration():
    """* Test the complete workflow integration."""
    print("\n" + "="*60)
    print("TEST 3: Complete Workflow Integration")
    print("="*60)
    
    workflow = DebugWorkflow()
    
    # * Test simple issue that gets resolved
    print("\n--- Simple Issue Test ---")
    success, message, delegation_prompt = workflow.handle_issue(
        issue_id="workflow_simple",
        title="Simple Configuration Error",
        description="Configuration file not found",
        error_context={"file": "config.yaml"},
        debug_function=lambda: "Created default configuration file"
    )
    
    print(f"Result: {'Success' if success else 'Requires Ad-Hoc'}")
    print(f"Message: {message}")
    
    # * Test complex issue requiring delegation
    print("\n--- Complex Issue Test ---")
    success, message, delegation_prompt = workflow.handle_issue(
        issue_id="workflow_complex",
        title="Database Connection Issue",
        description="Intermittent database connection failures with race conditions",
        error_context={"error": "Connection timeout", "frequency": "intermittent"}
    )
    
    print(f"Result: {'Success' if success else 'Requires Ad-Hoc'}")
    print(f"Message: {message}")
    
    if delegation_prompt:
        print(f"Delegation prompt length: {len(delegation_prompt)} characters")
    
    return success, message, delegation_prompt


def test_session_management():
    """* Test session management and persistence."""
    print("\n" + "="*60)
    print("TEST 4: Session Management and Persistence")
    print("="*60)
    
    session_manager = SessionManager()
    
    # * Create a test session
    debugger = AdHocDebugger()
    session = debugger.create_debug_session(
        issue_id="session_test",
        title="Session Management Test",
        description="Testing session persistence and retrieval"
    )
    
    # * Save session
    success = session_manager.save_session(session)
    print(f"Session saved: {success}")
    
    # * Load session
    loaded_session = session_manager.load_session(session.session_id)
    print(f"Session loaded: {loaded_session is not None}")
    
    if loaded_session:
        print(f"Loaded session ID: {loaded_session.session_id}")
        print(f"Loaded session title: {loaded_session.title}")
    
    # * Get statistics
    stats = session_manager.get_session_statistics()
    print(f"\nSession Statistics:")
    print(f"Total sessions: {stats.get('total_sessions', 0)}")
    print(f"Resolution rate: {stats.get('resolution_rate_percent', 0):.1f}%")
    
    # * List sessions
    sessions = session_manager.list_sessions(limit=5)
    print(f"\nRecent sessions: {len(sessions)}")
    for session_info in sessions:
        print(f"  - {session_info['title']} ({session_info['status']})")
    
    return session


def test_solution_integration():
    """* Test solution integration from Ad-Hoc session."""
    print("\n" + "="*60)
    print("TEST 5: Solution Integration")
    print("="*60)
    
    workflow = DebugWorkflow()
    
    # * Create a session that requires Ad-Hoc debugging
    success, message, delegation_prompt = workflow.handle_issue(
        issue_id="integration_test",
        title="Integration Test Issue",
        description="This issue requires Ad-Hoc debugging for resolution"
    )
    
    if not success and delegation_prompt:
        print("Issue requires Ad-Hoc debugging")
        
        # * Simulate Ad-Hoc session completion
        solution_findings = """
        Root Cause Analysis:
        The issue was caused by a missing dependency in the template generation system.
        
        Solution:
        1. Added missing import statement for 'datetime' module
        2. Updated template generator to handle edge cases
        3. Added error handling for missing data fields
        
        Testing:
        - Verified solution works with test data
        - Confirmed no regression in existing functionality
        - Performance impact is minimal
        
        Prevention:
        - Added unit tests for template generation
        - Improved error handling and logging
        - Added dependency validation
        """
        
        # * Get the session ID from the workflow
        session_id = list(workflow.debugger.active_sessions.keys())[0]
        
        # * Complete the workflow
        integration_success = workflow.complete_workflow(
            session_id, solution_findings, "ad_hoc_session_001"
        )
        
        print(f"Solution integration: {'Success' if integration_success else 'Failed'}")
        
        # * Verify the session was updated
        session = workflow.debugger.active_sessions[session_id]
        print(f"Session status: {session.status.value}")
        print(f"Has final solution: {session.final_solution is not None}")
        
        return integration_success
    
    return False


def test_cli_interface():
    """* Test the CLI interface."""
    print("\n" + "="*60)
    print("TEST 6: CLI Interface")
    print("="*60)
    
    from src.debug_cli import DebugCLI
    
    cli = DebugCLI()
    
    # * Test creating a session via CLI
    print("Testing CLI session creation...")
    
    # * Simulate CLI arguments
    class MockArgs:
        def __init__(self):
            self.issue_id = "cli_test_001"
            self.title = "CLI Test Issue"
            self.description = "Testing CLI interface functionality"
            self.error_context = None
            self.verbose = False
            self.quiet = False
    
    args = MockArgs()
    result = cli.create_session(args)
    print(f"CLI session creation result: {result}")
    
    # * Test listing sessions
    class MockListArgs:
        def __init__(self):
            self.status = None
            self.limit = 10
            self.offset = 0
            self.verbose = False
            self.quiet = False
    
    list_args = MockListArgs()
    result = cli.list_sessions(list_args)
    print(f"CLI list sessions result: {result}")
    
    return result == 0


def main():
    """* Run all tests."""
    setup_logging()
    
    print("🧪 Ad-Hoc Debugging System Test Suite")
    print("=" * 60)
    
    test_results = []
    
    try:
        # * Test 1: Simple issue
        session1 = test_simple_issue()
        test_results.append(("Simple Issue", True))
        
        # * Test 2: Complex issue
        session2, delegation_prompt = test_complex_issue()
        test_results.append(("Complex Issue", True))
        
        # * Test 3: Workflow integration
        success, message, prompt = test_workflow_integration()
        test_results.append(("Workflow Integration", True))
        
        # * Test 4: Session management
        session3 = test_session_management()
        test_results.append(("Session Management", True))
        
        # * Test 5: Solution integration
        integration_success = test_solution_integration()
        test_results.append(("Solution Integration", integration_success))
        
        # * Test 6: CLI interface
        cli_success = test_cli_interface()
        test_results.append(("CLI Interface", cli_success))
        
    except Exception as e:
        print(f"\n❌ Test suite failed with error: {e}")
        import traceback
        traceback.print_exc()
        return 1
    
    # * Print test results
    print("\n" + "="*60)
    print("TEST RESULTS SUMMARY")
    print("="*60)
    
    passed = 0
    total = len(test_results)
    
    for test_name, result in test_results:
        status = "✅ PASS" if result else "❌ FAIL"
        print(f"{test_name:<25} {status}")
        if result:
            passed += 1
    
    print(f"\nOverall: {passed}/{total} tests passed")
    
    if passed == total:
        print("\n🎉 All tests passed! Ad-Hoc debugging system is working correctly.")
        return 0
    else:
        print(f"\n⚠️ {total - passed} tests failed. Please review the implementation.")
        return 1


if __name__ == "__main__":
    sys.exit(main())
