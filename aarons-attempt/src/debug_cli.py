#!/usr/bin/env python3
"""
Command Line Interface for Ad-Hoc Debugging System

This module provides CLI commands for managing debug sessions, viewing statistics,
and interacting with the APM debugging workflow.
"""

import argparse
import json
import logging
import sys
from datetime import datetime
from pathlib import Path
from typing import Any, Dict, List, Optional

from .debug_system import AdHocDebugger, DebugWorkflow, IssueComplexity, DebugStatus
from .session_manager import SessionManager, SessionIntegration


class DebugCLI:
    """* Command line interface for the Ad-Hoc debugging system."""
    
    def __init__(self, workspace_root: str = None):
        """* Initialize the CLI with workspace configuration."""
        self.workspace_root = Path(workspace_root or ".")
        self.session_manager = SessionManager(str(self.workspace_root))
        self.session_integration = SessionIntegration(str(self.workspace_root))
        self.debugger = AdHocDebugger(str(self.workspace_root))
        self.workflow = DebugWorkflow(str(self.workspace_root))
        
        self.logger = logging.getLogger(__name__)
    
    def setup_logging(self, verbose: bool = False, quiet: bool = False):
        """* Setup logging configuration."""
        if quiet:
            level = logging.ERROR
        elif verbose:
            level = logging.DEBUG
        else:
            level = logging.INFO
        
        logging.basicConfig(
            level=level,
            format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
        )
    
    def create_session(self, args) -> int:
        """* Create a new debug session."""
        try:
            session = self.debugger.create_debug_session(
                issue_id=args.issue_id,
                title=args.title,
                description=args.description,
                error_context=json.loads(args.error_context) if args.error_context else None
            )
            
            print(f"✅ Created debug session: {session.session_id}")
            print(f"   Issue ID: {session.issue_id}")
            print(f"   Complexity: {session.complexity.value}")
            print(f"   Status: {session.status.value}")
            
            # * Save to database
            self.session_manager.save_session(session)
            
            return 0
            
        except Exception as e:
            print(f"❌ Failed to create session: {e}")
            return 1
    
    def list_sessions(self, args) -> int:
        """* List debug sessions."""
        try:
            status_filter = None
            if args.status:
                status_filter = DebugStatus(args.status)
            
            sessions = self.session_manager.list_sessions(
                status_filter=status_filter,
                limit=args.limit,
                offset=args.offset
            )
            
            if not sessions:
                print("No sessions found.")
                return 0
            
            # * Print header
            print(f"{'Session ID':<36} {'Issue ID':<20} {'Title':<30} {'Status':<12} {'Created':<20}")
            print("-" * 120)
            
            # * Print sessions
            for session in sessions:
                created_date = datetime.fromisoformat(session['created_at']).strftime('%Y-%m-%d %H:%M')
                print(f"{session['session_id']:<36} {session['issue_id']:<20} "
                      f"{session['title'][:29]:<30} {session['status']:<12} {created_date:<20}")
            
            print(f"\nTotal: {len(sessions)} sessions")
            return 0
            
        except Exception as e:
            print(f"❌ Failed to list sessions: {e}")
            return 1
    
    def show_session(self, args) -> int:
        """* Show detailed information about a session."""
        try:
            session = self.session_manager.load_session(args.session_id)
            if not session:
                print(f"❌ Session {args.session_id} not found")
                return 1
            
            # * Print session details
            print(f"# Debug Session: {session.title}")
            print(f"")
            print(f"**Session ID**: `{session.session_id}`")
            print(f"**Issue ID**: `{session.issue_id}`")
            print(f"**Complexity**: {session.complexity.value.title()}")
            print(f"**Status**: {session.status.value.title()}")
            print(f"**Created**: {session.created_at.strftime('%Y-%m-%d %H:%M:%S')}")
            print(f"**Updated**: {session.updated_at.strftime('%Y-%m-%d %H:%M:%S')}")
            print(f"")
            print(f"## Description")
            print(f"")
            print(f"{session.description}")
            print(f"")
            
            # * Print attempts
            if session.attempts:
                print(f"## Debug Attempts ({len(session.attempts)})")
                print(f"")
                
                for i, attempt in enumerate(session.attempts, 1):
                    status_icon = "✅" if attempt.success else "❌"
                    print(f"### Attempt {i} {status_icon}")
                    print(f"")
                    print(f"- **Timestamp**: {attempt.timestamp.strftime('%Y-%m-%d %H:%M:%S')}")
                    print(f"- **Description**: {attempt.description}")
                    
                    if attempt.duration_seconds:
                        print(f"- **Duration**: {attempt.duration_seconds:.2f} seconds")
                    
                    if attempt.error_message:
                        print(f"- **Error**:")
                        print(f"```")
                        print(f"{attempt.error_message}")
                        print(f"```")
                    
                    if attempt.solution:
                        print(f"- **Solution**:")
                        print(f"```")
                        print(f"{attempt.solution}")
                        print(f"```")
                    
                    print(f"")
            
            # * Print final solution
            if session.final_solution:
                print(f"## Final Solution")
                print(f"")
                print(f"{session.final_solution}")
                print(f"")
            
            # * Print escalation info
            if session.escalation_reason:
                print(f"## Escalation")
                print(f"")
                print(f"**Reason**: {session.escalation_reason}")
                print(f"")
            
            return 0
            
        except Exception as e:
            print(f"❌ Failed to show session: {e}")
            return 1
    
    def handle_issue(self, args) -> int:
        """* Handle an issue through the APM workflow."""
        try:
            error_context = None
            if args.error_context:
                error_context = json.loads(args.error_context)
            
            success, message, delegation_prompt = self.workflow.handle_issue(
                issue_id=args.issue_id,
                title=args.title,
                description=args.description,
                error_context=error_context
            )
            
            print(f"Workflow Result: {'✅ Success' if success else '⚠️ Requires Ad-Hoc Session'}")
            print(f"Message: {message}")
            
            # * Save the session to database
            if self.workflow.debugger.active_sessions:
                session_id = list(self.workflow.debugger.active_sessions.keys())[0]
                session = self.workflow.debugger.active_sessions[session_id]
                self.session_manager.save_session(session)
                print(f"💾 Session saved: {session_id}")
            
            if delegation_prompt:
                print(f"\n## Delegation Prompt for Ad-Hoc Session")
                print(f"")
                print(delegation_prompt)
                
                # * Save delegation prompt to file
                if args.save_prompt:
                    prompt_file = self.workspace_root / f"delegation_prompt_{args.issue_id}.md"
                    with open(prompt_file, 'w', encoding='utf-8') as f:
                        f.write(delegation_prompt)
                    print(f"\n💾 Delegation prompt saved to: {prompt_file}")
            
            return 0 if success else 2  # * Return 2 for "requires Ad-Hoc session"
            
        except Exception as e:
            print(f"❌ Failed to handle issue: {e}")
            return 1
    
    def complete_workflow(self, args) -> int:
        """* Complete a workflow with solution findings."""
        try:
            # * Load session from database if not in active sessions
            if args.session_id not in self.workflow.debugger.active_sessions:
                session = self.session_manager.load_session(args.session_id)
                if session:
                    self.workflow.debugger.active_sessions[args.session_id] = session
                    print(f"📂 Loaded session from database: {args.session_id}")
                else:
                    print(f"❌ Session {args.session_id} not found in database")
                    return 1
            
            success = self.workflow.complete_workflow(
                session_id=args.session_id,
                solution_findings=args.solution,
                ad_hoc_session_id=args.ad_hoc_session_id
            )
            
            if success:
                print(f"✅ Successfully completed workflow for session {args.session_id}")
            else:
                print(f"❌ Failed to complete workflow for session {args.session_id}")
            
            return 0 if success else 1
            
        except Exception as e:
            print(f"❌ Failed to complete workflow: {e}")
            return 1
    
    def show_statistics(self, args) -> int:
        """* Show debug session statistics."""
        try:
            stats = self.session_manager.get_session_statistics()
            
            print(f"# Debug Session Statistics")
            print(f"")
            print(f"**Total Sessions**: {stats.get('total_sessions', 0)}")
            print(f"**Resolution Rate**: {stats.get('resolution_rate_percent', 0):.1f}%")
            print(f"**Average Attempts per Session**: {stats.get('average_attempts_per_session', 0):.2f}")
            print(f"")
            
            # * Status breakdown
            status_breakdown = stats.get('status_breakdown', {})
            if status_breakdown:
                print(f"## Status Breakdown")
                print(f"")
                for status, count in status_breakdown.items():
                    print(f"- **{status.title()}**: {count}")
                print(f"")
            
            # * Complexity breakdown
            complexity_breakdown = stats.get('complexity_breakdown', {})
            if complexity_breakdown:
                print(f"## Complexity Breakdown")
                print(f"")
                for complexity, count in complexity_breakdown.items():
                    print(f"- **{complexity.title()}**: {count}")
                print(f"")
            
            return 0
            
        except Exception as e:
            print(f"❌ Failed to show statistics: {e}")
            return 1
    
    def export_session(self, args) -> int:
        """* Export session data."""
        try:
            data = self.session_manager.export_session_data(args.session_id, args.format)
            if not data:
                print(f"❌ Session {args.session_id} not found")
                return 1
            
            if args.output:
                output_file = Path(args.output)
                with open(output_file, 'w', encoding='utf-8') as f:
                    f.write(data)
                print(f"✅ Exported session to: {output_file}")
            else:
                print(data)
            
            return 0
            
        except Exception as e:
            print(f"❌ Failed to export session: {e}")
            return 1
    
    def search_sessions(self, args) -> int:
        """* Search sessions."""
        try:
            sessions = self.session_manager.search_sessions(args.query, args.limit)
            
            if not sessions:
                print(f"No sessions found matching '{args.query}'")
                return 0
            
            print(f"Found {len(sessions)} sessions matching '{args.query}':")
            print(f"")
            
            for session in sessions:
                created_date = datetime.fromisoformat(session['created_at']).strftime('%Y-%m-%d %H:%M')
                print(f"- **{session['title']}** ({session['session_id']})")
                print(f"  Issue ID: {session['issue_id']}")
                print(f"  Status: {session['status']}")
                print(f"  Created: {created_date}")
                print(f"")
            
            return 0
            
        except Exception as e:
            print(f"❌ Failed to search sessions: {e}")
            return 1
    
    def cleanup_sessions(self, args) -> int:
        """* Clean up old sessions."""
        try:
            deleted_count = self.session_manager.cleanup_old_sessions()
            print(f"✅ Cleaned up {deleted_count} old sessions")
            return 0
            
        except Exception as e:
            print(f"❌ Failed to cleanup sessions: {e}")
            return 1
    
    def run(self, args) -> int:
        """* Run the CLI command."""
        self.setup_logging(args.verbose, args.quiet)
        
        if args.command == "create":
            return self.create_session(args)
        elif args.command == "list":
            return self.list_sessions(args)
        elif args.command == "show":
            return self.show_session(args)
        elif args.command == "handle":
            return self.handle_issue(args)
        elif args.command == "complete":
            return self.complete_workflow(args)
        elif args.command == "stats":
            return self.show_statistics(args)
        elif args.command == "export":
            return self.export_session(args)
        elif args.command == "search":
            return self.search_sessions(args)
        elif args.command == "cleanup":
            return self.cleanup_sessions(args)
        else:
            print(f"❌ Unknown command: {args.command}")
            return 1


def create_parser() -> argparse.ArgumentParser:
    """* Create the argument parser."""
    parser = argparse.ArgumentParser(
        description="Ad-Hoc Debugging System CLI",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  %(prog)s create --issue-id "bug_001" --title "Import Error" --description "Module not found"
  %(prog)s list --status resolved --limit 10
  %(prog)s show --session-id "session_123"
  %(prog)s handle --issue-id "bug_002" --title "Build Failure" --description "LaTeX compilation failed"
  %(prog)s complete --session-id "session_123" --solution "Fixed import path"
  %(prog)s stats
  %(prog)s export --session-id "session_123" --format markdown --output report.md
  %(prog)s search --query "import error"
  %(prog)s cleanup
        """
    )
    
    # * Global options
    parser.add_argument("--workspace", "-w", help="Workspace root directory")
    parser.add_argument("--verbose", "-v", action="store_true", help="Enable verbose logging")
    parser.add_argument("--quiet", "-q", action="store_true", help="Suppress all output except errors")
    
    # * Subcommands
    subparsers = parser.add_subparsers(dest="command", help="Available commands")
    
    # * Create command
    create_parser = subparsers.add_parser("create", help="Create a new debug session")
    create_parser.add_argument("--issue-id", required=True, help="Issue ID")
    create_parser.add_argument("--title", required=True, help="Session title")
    create_parser.add_argument("--description", required=True, help="Session description")
    create_parser.add_argument("--error-context", help="Error context as JSON")
    
    # * List command
    list_parser = subparsers.add_parser("list", help="List debug sessions")
    list_parser.add_argument("--status", choices=[s.value for s in DebugStatus], help="Filter by status")
    list_parser.add_argument("--limit", type=int, default=50, help="Maximum number of sessions to show")
    list_parser.add_argument("--offset", type=int, default=0, help="Offset for pagination")
    
    # * Show command
    show_parser = subparsers.add_parser("show", help="Show detailed session information")
    show_parser.add_argument("--session-id", required=True, help="Session ID")
    
    # * Handle command
    handle_parser = subparsers.add_parser("handle", help="Handle an issue through APM workflow")
    handle_parser.add_argument("--issue-id", required=True, help="Issue ID")
    handle_parser.add_argument("--title", required=True, help="Issue title")
    handle_parser.add_argument("--description", required=True, help="Issue description")
    handle_parser.add_argument("--error-context", help="Error context as JSON")
    handle_parser.add_argument("--save-prompt", action="store_true", help="Save delegation prompt to file")
    
    # * Complete command
    complete_parser = subparsers.add_parser("complete", help="Complete workflow with solution")
    complete_parser.add_argument("--session-id", required=True, help="Session ID")
    complete_parser.add_argument("--solution", required=True, help="Solution findings")
    complete_parser.add_argument("--ad-hoc-session-id", help="Ad-Hoc session ID")
    
    # * Stats command
    subparsers.add_parser("stats", help="Show debug session statistics")
    
    # * Export command
    export_parser = subparsers.add_parser("export", help="Export session data")
    export_parser.add_argument("--session-id", required=True, help="Session ID")
    export_parser.add_argument("--format", choices=["json", "markdown"], default="markdown", help="Export format")
    export_parser.add_argument("--output", help="Output file (default: stdout)")
    
    # * Search command
    search_parser = subparsers.add_parser("search", help="Search sessions")
    search_parser.add_argument("--query", required=True, help="Search query")
    search_parser.add_argument("--limit", type=int, default=20, help="Maximum number of results")
    
    # * Cleanup command
    subparsers.add_parser("cleanup", help="Clean up old sessions")
    
    return parser


def main() -> int:
    """* Main entry point for the CLI."""
    parser = create_parser()
    args = parser.parse_args()
    
    if not args.command:
        parser.print_help()
        return 1
    
    cli = DebugCLI(args.workspace)
    return cli.run(args)


if __name__ == "__main__":
    sys.exit(main())
