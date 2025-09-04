#!/usr/bin/env python3
"""
Session Management System for Ad-Hoc Debugging

This module provides session management capabilities for the APM Ad-Hoc debugging system,
including session persistence, state tracking, and integration with the main application.
"""

import json
import logging
import os
import sqlite3
from datetime import datetime, timedelta
from pathlib import Path
from typing import Any, Dict, List, Optional, Tuple
from dataclasses import asdict
import threading
import time

from .debug_system import DebugSession, DebugStatus, AdHocDebugger, IssueComplexity


class SessionManager:
    """
    * Session management system for Ad-Hoc debugging sessions.
    
    Provides persistent storage, session tracking, and integration
    with the main application workflow.
    """
    
    def __init__(self, workspace_root: str = None, db_path: str = None):
        """* Initialize the session manager with database and file storage."""
        self.workspace_root = Path(workspace_root or os.getcwd())
        self.sessions_dir = self.workspace_root / "debug_sessions"
        self.sessions_dir.mkdir(exist_ok=True)
        
        # * Database setup
        self.db_path = db_path or str(self.sessions_dir / "sessions.db")
        self.init_database()
        
        self.logger = logging.getLogger(__name__)
        self.lock = threading.Lock()
        
        # * Session cleanup settings
        self.cleanup_interval = 3600  # 1 hour
        self.session_retention_days = 30
        self.last_cleanup = datetime.now()
    
    def init_database(self) -> None:
        """* Initialize SQLite database for session tracking."""
        with sqlite3.connect(self.db_path) as conn:
            conn.execute("""
                CREATE TABLE IF NOT EXISTS sessions (
                    session_id TEXT PRIMARY KEY,
                    issue_id TEXT NOT NULL,
                    title TEXT NOT NULL,
                    description TEXT,
                    complexity TEXT NOT NULL,
                    status TEXT NOT NULL,
                    created_at TIMESTAMP NOT NULL,
                    updated_at TIMESTAMP NOT NULL,
                    delegation_prompt_id TEXT,
                    ad_hoc_session_id TEXT,
                    final_solution TEXT,
                    escalation_reason TEXT,
                    metadata TEXT
                )
            """)
            
            conn.execute("""
                CREATE TABLE IF NOT EXISTS attempts (
                    attempt_id TEXT PRIMARY KEY,
                    session_id TEXT NOT NULL,
                    timestamp TIMESTAMP NOT NULL,
                    description TEXT,
                    success BOOLEAN NOT NULL,
                    error_message TEXT,
                    solution TEXT,
                    duration_seconds REAL,
                    FOREIGN KEY (session_id) REFERENCES sessions (session_id)
                )
            """)
            
            conn.execute("""
                CREATE TABLE IF NOT EXISTS delegation_prompts (
                    prompt_id TEXT PRIMARY KEY,
                    session_id TEXT NOT NULL,
                    issue_context TEXT,
                    specific_instructions TEXT,
                    expected_deliverables TEXT,
                    constraints TEXT,
                    success_criteria TEXT,
                    created_at TIMESTAMP NOT NULL,
                    FOREIGN KEY (session_id) REFERENCES sessions (session_id)
                )
            """)
            
            # * Create indexes for better performance
            conn.execute("CREATE INDEX IF NOT EXISTS idx_sessions_issue_id ON sessions (issue_id)")
            conn.execute("CREATE INDEX IF NOT EXISTS idx_sessions_status ON sessions (status)")
            conn.execute("CREATE INDEX IF NOT EXISTS idx_sessions_created_at ON sessions (created_at)")
            conn.execute("CREATE INDEX IF NOT EXISTS idx_attempts_session_id ON attempts (session_id)")
            conn.execute("CREATE INDEX IF NOT EXISTS idx_prompts_session_id ON delegation_prompts (session_id)")
    
    def save_session(self, session: DebugSession) -> bool:
        """* Save debug session to database and file system."""
        with self.lock:
            try:
                with sqlite3.connect(self.db_path) as conn:
                    # * Save session metadata
                    conn.execute("""
                        INSERT OR REPLACE INTO sessions 
                        (session_id, issue_id, title, description, complexity, status,
                         created_at, updated_at, delegation_prompt_id, ad_hoc_session_id,
                         final_solution, escalation_reason, metadata)
                        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                    """, (
                        session.session_id,
                        session.issue_id,
                        session.title,
                        session.description,
                        session.complexity.value,
                        session.status.value,
                        session.created_at.isoformat(),
                        session.updated_at.isoformat(),
                        session.delegation_prompt,
                        session.ad_hoc_session_id,
                        session.final_solution,
                        session.escalation_reason,
                        json.dumps({})  # * Placeholder for future metadata
                    ))
                    
                    # * Save attempts
                    conn.execute("DELETE FROM attempts WHERE session_id = ?", (session.session_id,))
                    for attempt in session.attempts:
                        conn.execute("""
                            INSERT INTO attempts 
                            (attempt_id, session_id, timestamp, description, success,
                             error_message, solution, duration_seconds)
                            VALUES (?, ?, ?, ?, ?, ?, ?, ?)
                        """, (
                            attempt.attempt_id,
                            session.session_id,
                            attempt.timestamp.isoformat(),
                            attempt.description,
                            attempt.success,
                            attempt.error_message,
                            attempt.solution,
                            attempt.duration_seconds
                        ))
                
                # * Save detailed session to JSON file
                session_file = self.sessions_dir / f"{session.session_id}.json"
                session_data = asdict(session)
                session_data['created_at'] = session.created_at.isoformat()
                session_data['updated_at'] = session.updated_at.isoformat()
                session_data['complexity'] = session.complexity.value
                session_data['status'] = session.status.value
                
                for attempt in session_data['attempts']:
                    attempt['timestamp'] = attempt['timestamp'].isoformat()
                
                with open(session_file, 'w', encoding='utf-8') as f:
                    json.dump(session_data, f, indent=2, ensure_ascii=False)
                
                self.logger.debug(f"Saved session {session.session_id} to database and file system")
                return True
                
            except Exception as e:
                self.logger.error(f"Failed to save session {session.session_id}: {e}")
                return False
    
    def load_session(self, session_id: str) -> Optional[DebugSession]:
        """* Load debug session from database."""
        with self.lock:
            try:
                with sqlite3.connect(self.db_path) as conn:
                    # * Load session metadata
                    session_row = conn.execute("""
                        SELECT session_id, issue_id, title, description, complexity, status,
                               created_at, updated_at, delegation_prompt_id, ad_hoc_session_id,
                               final_solution, escalation_reason
                        FROM sessions WHERE session_id = ?
                    """, (session_id,)).fetchone()
                    
                    if not session_row:
                        return None
                    
                    # * Load attempts
                    attempts_rows = conn.execute("""
                        SELECT attempt_id, timestamp, description, success, error_message,
                               solution, duration_seconds
                        FROM attempts WHERE session_id = ? ORDER BY timestamp
                    """, (session_id,)).fetchall()
                    
                    # * Reconstruct session object
                    from .debug_system import DebugAttempt
                    
                    attempts = []
                    for row in attempts_rows:
                        attempt = DebugAttempt(
                            attempt_id=row[0],
                            timestamp=datetime.fromisoformat(row[1]),
                            description=row[2],
                            success=bool(row[3]),
                            error_message=row[4],
                            solution=row[5],
                            duration_seconds=row[6]
                        )
                        attempts.append(attempt)
                    
                    session = DebugSession(
                        session_id=session_row[0],
                        issue_id=session_row[1],
                        title=session_row[2],
                        description=session_row[3],
                        complexity=IssueComplexity(session_row[4]),
                        status=DebugStatus(session_row[5]),
                        created_at=datetime.fromisoformat(session_row[6]),
                        updated_at=datetime.fromisoformat(session_row[7]),
                        attempts=attempts,
                        delegation_prompt=session_row[8],
                        ad_hoc_session_id=session_row[9],
                        final_solution=session_row[10],
                        escalation_reason=session_row[11]
                    )
                    
                    return session
                    
            except Exception as e:
                self.logger.error(f"Failed to load session {session_id}: {e}")
                return None
    
    def list_sessions(self, status_filter: Optional[DebugStatus] = None,
                     limit: int = 100, offset: int = 0) -> List[Dict[str, Any]]:
        """* List sessions with optional filtering."""
        with self.lock:
            try:
                with sqlite3.connect(self.db_path) as conn:
                    query = """
                        SELECT session_id, issue_id, title, complexity, status,
                               created_at, updated_at, final_solution
                        FROM sessions
                    """
                    params = []
                    
                    if status_filter:
                        query += " WHERE status = ?"
                        params.append(status_filter.value)
                    
                    query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
                    params.extend([limit, offset])
                    
                    rows = conn.execute(query, params).fetchall()
                    
                    sessions = []
                    for row in rows:
                        # * Get attempt count
                        attempt_count = conn.execute(
                            "SELECT COUNT(*) FROM attempts WHERE session_id = ?",
                            (row[0],)
                        ).fetchone()[0]
                        
                        sessions.append({
                            'session_id': row[0],
                            'issue_id': row[1],
                            'title': row[2],
                            'complexity': row[3],
                            'status': row[4],
                            'created_at': row[5],
                            'updated_at': row[6],
                            'has_solution': row[7] is not None,
                            'attempt_count': attempt_count
                        })
                    
                    return sessions
                    
            except Exception as e:
                self.logger.error(f"Failed to list sessions: {e}")
                return []
    
    def get_session_statistics(self) -> Dict[str, Any]:
        """* Get statistics about debug sessions."""
        with self.lock:
            try:
                with sqlite3.connect(self.db_path) as conn:
                    # * Total sessions
                    total_sessions = conn.execute("SELECT COUNT(*) FROM sessions").fetchone()[0]
                    
                    # * Sessions by status
                    status_counts = {}
                    for status in DebugStatus:
                        count = conn.execute(
                            "SELECT COUNT(*) FROM sessions WHERE status = ?",
                            (status.value,)
                        ).fetchone()[0]
                        status_counts[status.value] = count
                    
                    # * Sessions by complexity
                    complexity_counts = {}
                    complexity_rows = conn.execute("""
                        SELECT complexity, COUNT(*) FROM sessions GROUP BY complexity
                    """).fetchall()
                    for row in complexity_rows:
                        complexity_counts[row[0]] = row[1]
                    
                    # * Average attempts per session
                    avg_attempts = conn.execute("""
                        SELECT AVG(attempt_count) FROM (
                            SELECT session_id, COUNT(*) as attempt_count 
                            FROM attempts GROUP BY session_id
                        )
                    """).fetchone()[0] or 0
                    
                    # * Resolution rate
                    resolved_sessions = conn.execute(
                        "SELECT COUNT(*) FROM sessions WHERE status = ?",
                        (DebugStatus.RESOLVED.value,)
                    ).fetchone()[0]
                    resolution_rate = (resolved_sessions / total_sessions * 100) if total_sessions > 0 else 0
                    
                    return {
                        'total_sessions': total_sessions,
                        'status_breakdown': status_counts,
                        'complexity_breakdown': complexity_counts,
                        'average_attempts_per_session': round(avg_attempts, 2),
                        'resolution_rate_percent': round(resolution_rate, 2),
                        'last_updated': datetime.now().isoformat()
                    }
                    
            except Exception as e:
                self.logger.error(f"Failed to get session statistics: {e}")
                return {}
    
    def cleanup_old_sessions(self) -> int:
        """* Clean up old sessions based on retention policy."""
        cutoff_date = datetime.now() - timedelta(days=self.session_retention_days)
        
        with self.lock:
            try:
                with sqlite3.connect(self.db_path) as conn:
                    # * Find old sessions
                    old_sessions = conn.execute("""
                        SELECT session_id FROM sessions 
                        WHERE created_at < ? AND status IN (?, ?)
                    """, (
                        cutoff_date.isoformat(),
                        DebugStatus.RESOLVED.value,
                        DebugStatus.ESCALATED.value
                    )).fetchall()
                    
                    deleted_count = 0
                    for (session_id,) in old_sessions:
                        # * Delete from database
                        conn.execute("DELETE FROM attempts WHERE session_id = ?", (session_id,))
                        conn.execute("DELETE FROM delegation_prompts WHERE session_id = ?", (session_id,))
                        conn.execute("DELETE FROM sessions WHERE session_id = ?", (session_id,))
                        
                        # * Delete JSON file
                        session_file = self.sessions_dir / f"{session_id}.json"
                        if session_file.exists():
                            session_file.unlink()
                        
                        deleted_count += 1
                    
                    conn.commit()
                    
                    if deleted_count > 0:
                        self.logger.info(f"Cleaned up {deleted_count} old sessions")
                    
                    self.last_cleanup = datetime.now()
                    return deleted_count
                    
            except Exception as e:
                self.logger.error(f"Failed to cleanup old sessions: {e}")
                return 0
    
    def auto_cleanup_check(self) -> None:
        """* Check if automatic cleanup should be performed."""
        if (datetime.now() - self.last_cleanup).total_seconds() > self.cleanup_interval:
            self.cleanup_old_sessions()
    
    def export_session_data(self, session_id: str, format: str = "json") -> Optional[str]:
        """* Export session data in specified format."""
        session = self.load_session(session_id)
        if not session:
            return None
        
        if format.lower() == "json":
            session_data = asdict(session)
            session_data['created_at'] = session.created_at.isoformat()
            session_data['updated_at'] = session.updated_at.isoformat()
            
            for attempt in session_data['attempts']:
                attempt['timestamp'] = attempt['timestamp'].isoformat()
            
            return json.dumps(session_data, indent=2, ensure_ascii=False)
        
        elif format.lower() == "markdown":
            return self._format_session_as_markdown(session)
        
        else:
            raise ValueError(f"Unsupported export format: {format}")
    
    def _format_session_as_markdown(self, session: DebugSession) -> str:
        """* Format session as markdown report."""
        md_lines = [
            f"# Debug Session Report",
            f"",
            f"**Session ID**: `{session.session_id}`",
            f"**Issue ID**: `{session.issue_id}`",
            f"**Title**: {session.title}",
            f"**Complexity**: {session.complexity.value.title()}",
            f"**Status**: {session.status.value.title()}",
            f"**Created**: {session.created_at.strftime('%Y-%m-%d %H:%M:%S')}",
            f"**Updated**: {session.updated_at.strftime('%Y-%m-%d %H:%M:%S')}",
            f"",
            f"## Description",
            f"",
            f"{session.description}",
            f"",
            f"## Debug Attempts",
            f""
        ]
        
        for i, attempt in enumerate(session.attempts, 1):
            md_lines.extend([
                f"### Attempt {i}",
                f"",
                f"- **Timestamp**: {attempt.timestamp.strftime('%Y-%m-%d %H:%M:%S')}",
                f"- **Success**: {'✅' if attempt.success else '❌'}",
                f"- **Description**: {attempt.description}",
            ])
            
            if attempt.duration_seconds:
                md_lines.append(f"- **Duration**: {attempt.duration_seconds:.2f} seconds")
            
            if attempt.error_message:
                md_lines.extend([
                    f"- **Error**:",
                    f"```",
                    f"{attempt.error_message}",
                    f"```"
                ])
            
            if attempt.solution:
                md_lines.extend([
                    f"- **Solution**:",
                    f"```",
                    f"{attempt.solution}",
                    f"```"
                ])
            
            md_lines.append("")
        
        if session.final_solution:
            md_lines.extend([
                f"## Final Solution",
                f"",
                f"{session.final_solution}",
                f""
            ])
        
        if session.escalation_reason:
            md_lines.extend([
                f"## Escalation",
                f"",
                f"**Reason**: {session.escalation_reason}",
                f""
            ])
        
        return "\n".join(md_lines)
    
    def search_sessions(self, query: str, limit: int = 50) -> List[Dict[str, Any]]:
        """* Search sessions by title, description, or issue ID."""
        with self.lock:
            try:
                with sqlite3.connect(self.db_path) as conn:
                    search_query = f"%{query}%"
                    rows = conn.execute("""
                        SELECT session_id, issue_id, title, description, complexity, status,
                               created_at, updated_at
                        FROM sessions 
                        WHERE title LIKE ? OR description LIKE ? OR issue_id LIKE ?
                        ORDER BY created_at DESC LIMIT ?
                    """, (search_query, search_query, search_query, limit)).fetchall()
                    
                    sessions = []
                    for row in rows:
                        sessions.append({
                            'session_id': row[0],
                            'issue_id': row[1],
                            'title': row[2],
                            'description': row[3],
                            'complexity': row[4],
                            'status': row[5],
                            'created_at': row[6],
                            'updated_at': row[7]
                        })
                    
                    return sessions
                    
            except Exception as e:
                self.logger.error(f"Failed to search sessions: {e}")
                return []


class SessionIntegration:
    """
    * Integration layer between session management and the main application.
    
    Provides high-level methods for integrating debug sessions with
    the LaTeX Gantt chart generator workflow.
    """
    
    def __init__(self, workspace_root: str = None):
        """* Initialize session integration."""
        self.session_manager = SessionManager(workspace_root)
        self.debugger = AdHocDebugger(workspace_root)
        self.logger = logging.getLogger(__name__)
    
    def handle_application_error(self, error: Exception, context: Dict[str, Any] = None) -> Tuple[bool, str, Optional[str]]:
        """
        * Handle application errors through the APM debugging workflow.
        
        Args:
            error: Exception that occurred
            context: Additional context about the error
            
        Returns:
            Tuple[bool, str, Optional[str]]: (success, message, delegation_prompt)
        """
        # * Generate issue information from error
        issue_id = f"app_error_{int(datetime.now().timestamp())}"
        title = f"Application Error: {type(error).__name__}"
        description = f"An error occurred in the application: {str(error)}"
        
        if context:
            description += f"\n\nContext: {json.dumps(context, indent=2)}"
        
        error_context = {
            "exception_type": type(error).__name__,
            "error_message": str(error),
            "context": context or {}
        }
        
        # * Use the debug workflow
        workflow = self.debugger
        session = workflow.create_debug_session(issue_id, title, description, error_context)
        
        # * Try to handle based on complexity
        if session.complexity.value == "simple":
            # * For simple errors, try basic recovery
            try:
                # * This would be replaced with actual recovery logic
                recovery_result = self._attempt_basic_recovery(error, context)
                if recovery_result:
                    session.status = DebugStatus.RESOLVED
                    session.final_solution = recovery_result
                    self.session_manager.save_session(session)
                    return True, "Error resolved with basic recovery", None
            except Exception as recovery_error:
                self.logger.warning(f"Basic recovery failed: {recovery_error}")
        
        # * Generate delegation prompt for complex or unresolved issues
        delegation_prompt = workflow.generate_delegation_prompt(session.session_id)
        formatted_prompt = workflow.format_delegation_prompt_for_ai(delegation_prompt)
        
        # * Save session
        self.session_manager.save_session(session)
        
        return False, f"Error requires Ad-Hoc debugging session", formatted_prompt
    
    def _attempt_basic_recovery(self, error: Exception, context: Dict[str, Any] = None) -> Optional[str]:
        """* Attempt basic error recovery for simple issues."""
        # * This is a placeholder for basic recovery logic
        # * In a real implementation, this would contain specific recovery strategies
        
        error_type = type(error).__name__
        
        if error_type == "FileNotFoundError":
            return "File not found - check file path and permissions"
        elif error_type == "PermissionError":
            return "Permission denied - check file/directory permissions"
        elif error_type == "ImportError":
            return "Import error - check module installation and Python path"
        elif error_type == "SyntaxError":
            return "Syntax error - check code syntax and indentation"
        else:
            return None  # * No basic recovery available
    
    def get_debug_dashboard_data(self) -> Dict[str, Any]:
        """* Get data for debug dashboard display."""
        stats = self.session_manager.get_session_statistics()
        recent_sessions = self.session_manager.list_sessions(limit=10)
        
        return {
            'statistics': stats,
            'recent_sessions': recent_sessions,
            'active_sessions': len(self.debugger.active_sessions),
            'last_cleanup': self.session_manager.last_cleanup.isoformat()
        }
    
    def export_debug_report(self, session_id: str, format: str = "markdown") -> Optional[str]:
        """* Export debug session as a report."""
        return self.session_manager.export_session_data(session_id, format)


# * Global session manager instance
_session_manager = None
_session_integration = None


def get_session_manager(workspace_root: str = None) -> SessionManager:
    """* Get global session manager instance."""
    global _session_manager
    if _session_manager is None:
        _session_manager = SessionManager(workspace_root)
    return _session_manager


def get_session_integration(workspace_root: str = None) -> SessionIntegration:
    """* Get global session integration instance."""
    global _session_integration
    if _session_integration is None:
        _session_integration = SessionIntegration(workspace_root)
    return _session_integration


def main():
    """* Example usage of the session management system."""
    # * Setup logging
    logging.basicConfig(level=logging.INFO)
    
    # * Create session manager
    session_manager = SessionManager()
    
    # * Example: Create and save a session
    from .debug_system import DebugSession, DebugStatus, IssueComplexity, DebugAttempt
    
    session = DebugSession(
        session_id="test_session_001",
        issue_id="test_issue_001",
        title="Test Session",
        description="This is a test session for the session manager",
        complexity=IssueComplexity.SIMPLE,
        status=DebugStatus.PENDING,
        created_at=datetime.now(),
        updated_at=datetime.now(),
        attempts=[]
    )
    
    # * Add a test attempt
    attempt = DebugAttempt(
        attempt_id="test_attempt_001",
        timestamp=datetime.now(),
        description="Test debug attempt",
        success=True,
        solution="Test solution"
    )
    session.attempts.append(attempt)
    
    # * Save session
    success = session_manager.save_session(session)
    print(f"Session saved: {success}")
    
    # * Load session
    loaded_session = session_manager.load_session("test_session_001")
    print(f"Session loaded: {loaded_session is not None}")
    
    # * Get statistics
    stats = session_manager.get_session_statistics()
    print(f"Statistics: {stats}")
    
    # * List sessions
    sessions = session_manager.list_sessions()
    print(f"Total sessions: {len(sessions)}")


if __name__ == "__main__":
    main()
