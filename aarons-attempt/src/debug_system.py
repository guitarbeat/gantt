#!/usr/bin/env python3
"""
Ad-Hoc Debugger Chat System for Agentic Project Management (APM)

This module implements the APM Ad-Hoc debugging workflow as described in the framework.
It provides structured debugging capabilities with delegation prompts and session management.
"""

import json
import logging
import os
import sys
from datetime import datetime
from enum import Enum
from pathlib import Path
from typing import Any, Dict, List, Optional, Tuple, Union
from dataclasses import dataclass, asdict
import uuid


class IssueComplexity(Enum):
    """* Issue complexity levels for debugging workflow."""
    SIMPLE = "simple"
    COMPLEX = "complex"


class DebugStatus(Enum):
    """* Debug session status tracking."""
    PENDING = "pending"
    IN_PROGRESS = "in_progress"
    RESOLVED = "resolved"
    FAILED = "failed"
    ESCALATED = "escalated"


@dataclass
class DebugAttempt:
    """* Individual debug attempt record."""
    attempt_id: str
    timestamp: datetime
    description: str
    success: bool
    error_message: Optional[str] = None
    solution: Optional[str] = None
    duration_seconds: Optional[float] = None


@dataclass
class DebugSession:
    """* Complete debug session with all attempts and metadata."""
    session_id: str
    issue_id: str
    title: str
    description: str
    complexity: IssueComplexity
    status: DebugStatus
    created_at: datetime
    updated_at: datetime
    attempts: List[DebugAttempt]
    delegation_prompt: Optional[str] = None
    ad_hoc_session_id: Optional[str] = None
    final_solution: Optional[str] = None
    escalation_reason: Optional[str] = None


@dataclass
class DelegationPrompt:
    """* Structured delegation prompt for Ad-Hoc sessions."""
    prompt_id: str
    session_id: str
    issue_context: str
    specific_instructions: str
    expected_deliverables: List[str]
    constraints: List[str]
    success_criteria: List[str]
    created_at: datetime


class AdHocDebugger:
    """
    * Main Ad-Hoc Debugger class implementing APM workflow.
    
    Handles issue complexity assessment, local debugging attempts,
    delegation prompt generation, and session management.
    """
    
    def __init__(self, workspace_root: str = None):
        """* Initialize the Ad-Hoc Debugger with workspace configuration."""
        self.workspace_root = Path(workspace_root or os.getcwd())
        self.sessions_dir = self.workspace_root / "debug_sessions"
        self.sessions_dir.mkdir(exist_ok=True)
        
        self.logger = logging.getLogger(__name__)
        self.max_local_attempts = 2
        
        # * Session tracking
        self.active_sessions: Dict[str, DebugSession] = {}
        self.session_history: List[str] = []
    
    def assess_issue_complexity(self, issue_description: str, 
                              error_context: Dict[str, Any] = None) -> IssueComplexity:
        """
        * Assess issue complexity based on description and context.
        
        Args:
            issue_description: Human-readable description of the issue
            error_context: Additional context like stack traces, logs, etc.
            
        Returns:
            IssueComplexity: SIMPLE or COMPLEX
        """
        complexity_indicators = [
            # * Simple issue indicators
            "syntax error", "import error", "typo", "missing file",
            "configuration error", "permission denied", "file not found",
            
            # * Complex issue indicators  
            "race condition", "memory leak", "deadlock", "performance issue",
            "integration error", "data corruption", "concurrent access",
            "distributed system", "network timeout", "database deadlock"
        ]
        
        issue_lower = issue_description.lower()
        error_text = ""
        
        if error_context:
            error_text = " ".join(str(v) for v in error_context.values()).lower()
        
        combined_text = f"{issue_lower} {error_text}"
        
        # * Count complexity indicators
        simple_count = sum(1 for indicator in complexity_indicators[:6] 
                          if indicator in combined_text)
        complex_count = sum(1 for indicator in complexity_indicators[6:] 
                           if indicator in combined_text)
        
        # * Additional heuristics
        if len(issue_description) > 200:
            complex_count += 1
        if error_context and len(str(error_context)) > 1000:
            complex_count += 1
        if "traceback" in combined_text or "stack trace" in combined_text:
            complex_count += 1
        
        return IssueComplexity.COMPLEX if complex_count > simple_count else IssueComplexity.SIMPLE
    
    def create_debug_session(self, issue_id: str, title: str, description: str,
                           error_context: Dict[str, Any] = None) -> DebugSession:
        """
        * Create a new debug session with complexity assessment.
        
        Args:
            issue_id: Unique identifier for the issue
            title: Short title for the issue
            description: Detailed description of the issue
            error_context: Additional error context
            
        Returns:
            DebugSession: Newly created debug session
        """
        complexity = self.assess_issue_complexity(description, error_context)
        session_id = str(uuid.uuid4())
        
        session = DebugSession(
            session_id=session_id,
            issue_id=issue_id,
            title=title,
            description=description,
            complexity=complexity,
            status=DebugStatus.PENDING,
            created_at=datetime.now(),
            updated_at=datetime.now(),
            attempts=[]
        )
        
        self.active_sessions[session_id] = session
        self.session_history.append(session_id)
        
        self.logger.info(f"Created debug session {session_id} for issue {issue_id} "
                        f"(complexity: {complexity.value})")
        
        return session
    
    def attempt_local_debug(self, session_id: str, debug_function: callable,
                          *args, **kwargs) -> DebugAttempt:
        """
        * Attempt local debugging with structured attempt tracking.
        
        Args:
            session_id: Debug session ID
            debug_function: Function to call for debugging
            *args, **kwargs: Arguments to pass to debug function
            
        Returns:
            DebugAttempt: Record of this debug attempt
        """
        if session_id not in self.active_sessions:
            raise ValueError(f"Session {session_id} not found")
        
        session = self.active_sessions[session_id]
        
        if len(session.attempts) >= self.max_local_attempts:
            raise ValueError(f"Maximum local attempts ({self.max_local_attempts}) exceeded")
        
        attempt_id = str(uuid.uuid4())
        start_time = datetime.now()
        
        self.logger.info(f"Starting local debug attempt {attempt_id} for session {session_id}")
        
        try:
            # * Execute debug function
            result = debug_function(*args, **kwargs)
            duration = (datetime.now() - start_time).total_seconds()
            
            attempt = DebugAttempt(
                attempt_id=attempt_id,
                timestamp=start_time,
                description=f"Local debug attempt {len(session.attempts) + 1}",
                success=True,
                solution=str(result) if result else "Debug function completed successfully",
                duration_seconds=duration
            )
            
            self.logger.info(f"Local debug attempt {attempt_id} succeeded")
            
        except Exception as e:
            duration = (datetime.now() - start_time).total_seconds()
            
            attempt = DebugAttempt(
                attempt_id=attempt_id,
                timestamp=start_time,
                description=f"Local debug attempt {len(session.attempts) + 1}",
                success=False,
                error_message=str(e),
                duration_seconds=duration
            )
            
            self.logger.warning(f"Local debug attempt {attempt_id} failed: {e}")
        
        # * Update session
        session.attempts.append(attempt)
        session.updated_at = datetime.now()
        
        # * Check if we should escalate
        if not attempt.success and len(session.attempts) >= self.max_local_attempts:
            session.status = DebugStatus.ESCALATED
            self.logger.info(f"Session {session_id} escalated after {len(session.attempts)} failed attempts")
        
        return attempt
    
    def generate_delegation_prompt(self, session_id: str) -> DelegationPrompt:
        """
        * Generate structured delegation prompt for Ad-Hoc session.
        
        Args:
            session_id: Debug session ID
            
        Returns:
            DelegationPrompt: Structured prompt for Ad-Hoc agent
        """
        if session_id not in self.active_sessions:
            raise ValueError(f"Session {session_id} not found")
        
        session = self.active_sessions[session_id]
        
        # * Build context from session and attempts
        context_parts = [
            f"Issue ID: {session.issue_id}",
            f"Title: {session.title}",
            f"Description: {session.description}",
            f"Complexity: {session.complexity.value}",
            f"Created: {session.created_at.isoformat()}"
        ]
        
        if session.attempts:
            context_parts.append("\nPrevious Debug Attempts:")
            for i, attempt in enumerate(session.attempts, 1):
                context_parts.append(f"  Attempt {i} ({attempt.timestamp.isoformat()}):")
                context_parts.append(f"    Success: {attempt.success}")
                if attempt.error_message:
                    context_parts.append(f"    Error: {attempt.error_message}")
                if attempt.solution:
                    context_parts.append(f"    Solution: {attempt.solution}")
        
        issue_context = "\n".join(context_parts)
        
        # * Generate specific instructions based on complexity
        if session.complexity == IssueComplexity.COMPLEX:
            specific_instructions = """
            This is a complex issue requiring deep analysis. Please:
            1. Analyze the root cause systematically
            2. Review related code and dependencies
            3. Check for patterns or similar issues
            4. Provide a comprehensive solution with testing steps
            5. Document any potential side effects or considerations
            """
        else:
            specific_instructions = """
            This is a simple issue that has persisted through local attempts. Please:
            1. Verify the exact error conditions
            2. Check for common solutions or workarounds
            3. Provide a clear, step-by-step fix
            4. Test the solution thoroughly
            """
        
        prompt = DelegationPrompt(
            prompt_id=str(uuid.uuid4()),
            session_id=session_id,
            issue_context=issue_context,
            specific_instructions=specific_instructions,
            expected_deliverables=[
                "Root cause analysis",
                "Step-by-step solution",
                "Testing verification",
                "Prevention recommendations"
            ],
            constraints=[
                "Maintain code quality standards",
                "Preserve existing functionality",
                "Follow project conventions",
                "Document all changes"
            ],
            success_criteria=[
                "Issue is completely resolved",
                "Solution is tested and verified",
                "No new issues introduced",
                "Code follows project standards"
            ],
            created_at=datetime.now()
        )
        
        # * Update session with delegation prompt
        session.delegation_prompt = prompt.prompt_id
        session.status = DebugStatus.IN_PROGRESS
        session.updated_at = datetime.now()
        
        self.logger.info(f"Generated delegation prompt {prompt.prompt_id} for session {session_id}")
        
        return prompt
    
    def format_delegation_prompt_for_ai(self, prompt: DelegationPrompt) -> str:
        """
        * Format delegation prompt for AI assistant consumption.
        
        Args:
            prompt: DelegationPrompt object
            
        Returns:
            str: Formatted prompt string for AI
        """
        formatted_prompt = f"""
# Ad-Hoc Debug Session - Delegation Prompt

## Session Information
- **Prompt ID**: {prompt.prompt_id}
- **Session ID**: {prompt.session_id}
- **Created**: {prompt.created_at.isoformat()}

## Issue Context
```
{prompt.issue_context}
```

## Specific Instructions
{prompt.specific_instructions}

## Expected Deliverables
{chr(10).join(f"- {item}" for item in prompt.expected_deliverables)}

## Constraints
{chr(10).join(f"- {item}" for item in prompt.constraints)}

## Success Criteria
{chr(10).join(f"- {item}" for item in prompt.success_criteria)}

## Instructions for Ad-Hoc Agent
You are now operating as an Ad-Hoc Debug Agent. Your task is to:

1. **Analyze** the issue context thoroughly
2. **Investigate** the root cause using available tools and information
3. **Develop** a comprehensive solution
4. **Test** the solution to ensure it works
5. **Document** your findings and solution

Please provide your analysis and solution in a structured format. When you have completed your work, return to the main session with your findings.

---
*This is an Ad-Hoc debugging session as part of the Agentic Project Management (APM) framework.*
"""
        return formatted_prompt
    
    def integrate_solution(self, session_id: str, solution_findings: str,
                         ad_hoc_session_id: str = None) -> bool:
        """
        * Integrate solution findings from Ad-Hoc session back to main workflow.
        
        Args:
            session_id: Original debug session ID
            solution_findings: Solution findings from Ad-Hoc agent
            ad_hoc_session_id: ID of the Ad-Hoc session (optional)
            
        Returns:
            bool: True if integration was successful
        """
        if session_id not in self.active_sessions:
            self.logger.error(f"Session {session_id} not found for solution integration")
            return False
        
        session = self.active_sessions[session_id]
        
        # * Create solution attempt record
        solution_attempt = DebugAttempt(
            attempt_id=str(uuid.uuid4()),
            timestamp=datetime.now(),
            description="Ad-Hoc agent solution integration",
            success=True,
            solution=solution_findings
        )
        
        session.attempts.append(solution_attempt)
        session.final_solution = solution_findings
        session.ad_hoc_session_id = ad_hoc_session_id
        session.status = DebugStatus.RESOLVED
        session.updated_at = datetime.now()
        
        self.logger.info(f"Successfully integrated solution for session {session_id}")
        
        # * Save session to disk
        self.save_session(session)
        
        # * Also save to database if session manager is available
        try:
            from .session_manager import get_session_manager
            session_manager = get_session_manager()
            session_manager.save_session(session)
        except Exception as e:
            self.logger.warning(f"Failed to save session to database: {e}")
        
        return True
    
    def escalate_to_manager(self, session_id: str, escalation_reason: str) -> bool:
        """
        * Escalate unresolved session to manager level.
        
        Args:
            session_id: Debug session ID
            escalation_reason: Reason for escalation
            
        Returns:
            bool: True if escalation was successful
        """
        if session_id not in self.active_sessions:
            self.logger.error(f"Session {session_id} not found for escalation")
            return False
        
        session = self.active_sessions[session_id]
        session.status = DebugStatus.ESCALATED
        session.escalation_reason = escalation_reason
        session.updated_at = datetime.now()
        
        self.logger.warning(f"Session {session_id} escalated to manager: {escalation_reason}")
        
        # * Save session to disk
        self.save_session(session)
        
        return True
    
    def save_session(self, session: DebugSession) -> None:
        """* Save debug session to disk for persistence."""
        session_file = self.sessions_dir / f"{session.session_id}.json"
        
        # * Convert to serializable format
        session_data = asdict(session)
        session_data['created_at'] = session.created_at.isoformat()
        session_data['updated_at'] = session.updated_at.isoformat()
        session_data['complexity'] = session.complexity.value
        session_data['status'] = session.status.value
        
        for attempt in session_data['attempts']:
            attempt['timestamp'] = attempt['timestamp'].isoformat()
        
        with open(session_file, 'w', encoding='utf-8') as f:
            json.dump(session_data, f, indent=2, ensure_ascii=False)
        
        self.logger.debug(f"Saved session {session.session_id} to {session_file}")
    
    def load_session(self, session_id: str) -> Optional[DebugSession]:
        """* Load debug session from disk."""
        session_file = self.sessions_dir / f"{session_id}.json"
        
        if not session_file.exists():
            return None
        
        try:
            with open(session_file, 'r', encoding='utf-8') as f:
                session_data = json.load(f)
            
            # * Convert back to objects
            session_data['created_at'] = datetime.fromisoformat(session_data['created_at'])
            session_data['updated_at'] = datetime.fromisoformat(session_data['updated_at'])
            session_data['complexity'] = IssueComplexity(session_data['complexity'])
            session_data['status'] = DebugStatus(session_data['status'])
            
            for attempt_data in session_data['attempts']:
                attempt_data['timestamp'] = datetime.fromisoformat(attempt_data['timestamp'])
            
            return DebugSession(**session_data)
            
        except Exception as e:
            self.logger.error(f"Failed to load session {session_id}: {e}")
            return None
    
    def get_session_summary(self, session_id: str) -> Optional[Dict[str, Any]]:
        """* Get summary information for a debug session."""
        if session_id not in self.active_sessions:
            session = self.load_session(session_id)
            if not session:
                return None
        else:
            session = self.active_sessions[session_id]
        
        return {
            'session_id': session.session_id,
            'issue_id': session.issue_id,
            'title': session.title,
            'complexity': session.complexity.value,
            'status': session.status.value,
            'created_at': session.created_at.isoformat(),
            'updated_at': session.updated_at.isoformat(),
            'attempt_count': len(session.attempts),
            'successful_attempts': sum(1 for a in session.attempts if a.success),
            'has_solution': session.final_solution is not None,
            'escalated': session.status == DebugStatus.ESCALATED
        }
    
    def list_active_sessions(self) -> List[Dict[str, Any]]:
        """* List all active debug sessions."""
        return [self.get_session_summary(session_id) 
                for session_id in self.active_sessions.keys()]
    
    def cleanup_completed_sessions(self) -> int:
        """* Clean up completed sessions from memory (keep on disk)."""
        completed_statuses = {DebugStatus.RESOLVED, DebugStatus.ESCALATED}
        to_remove = []
        
        for session_id, session in self.active_sessions.items():
            if session.status in completed_statuses:
                to_remove.append(session_id)
        
        for session_id in to_remove:
            del self.active_sessions[session_id]
        
        self.logger.info(f"Cleaned up {len(to_remove)} completed sessions from memory")
        return len(to_remove)


class DebugWorkflow:
    """
    * High-level workflow orchestrator for the APM debugging process.
    
    Implements the complete workflow from issue assessment to resolution.
    """
    
    def __init__(self, workspace_root: str = None):
        """* Initialize the debug workflow."""
        self.debugger = AdHocDebugger(workspace_root)
        self.logger = logging.getLogger(__name__)
    
    def handle_issue(self, issue_id: str, title: str, description: str,
                    error_context: Dict[str, Any] = None,
                    debug_function: callable = None) -> Tuple[bool, str, Optional[str]]:
        """
        * Handle an issue through the complete APM debugging workflow.
        
        Args:
            issue_id: Unique identifier for the issue
            title: Short title for the issue
            description: Detailed description of the issue
            error_context: Additional error context
            debug_function: Function to call for local debugging attempts
            
        Returns:
            Tuple[bool, str, Optional[str]]: (success, message, delegation_prompt)
        """
        self.logger.info(f"Starting APM debug workflow for issue: {issue_id}")
        
        # * Step 1: Create debug session and assess complexity
        session = self.debugger.create_debug_session(
            issue_id, title, description, error_context
        )
        
        # * Step 2: Handle based on complexity
        if session.complexity == IssueComplexity.SIMPLE and debug_function:
            # * Try local debugging for simple issues
            self.logger.info(f"Issue {issue_id} is simple, attempting local debugging")
            
            for attempt_num in range(self.debugger.max_local_attempts):
                try:
                    attempt = self.debugger.attempt_local_debug(
                        session.session_id, debug_function
                    )
                    
                    if attempt.success:
                        self.logger.info(f"Issue {issue_id} resolved with local debugging")
                        return True, f"Issue resolved in local attempt {attempt_num + 1}", None
                    
                except Exception as e:
                    self.logger.warning(f"Local debug attempt {attempt_num + 1} failed: {e}")
            
            # * Local debugging failed, escalate to delegation
            self.logger.info(f"Local debugging failed for issue {issue_id}, generating delegation prompt")
        
        # * Step 3: Generate delegation prompt for Ad-Hoc session
        delegation_prompt = self.debugger.generate_delegation_prompt(session.session_id)
        formatted_prompt = self.debugger.format_delegation_prompt_for_ai(delegation_prompt)
        
        self.logger.info(f"Issue {issue_id} requires Ad-Hoc debugging session")
        
        return False, f"Issue requires Ad-Hoc debugging session", formatted_prompt
    
    def complete_workflow(self, session_id: str, solution_findings: str,
                         ad_hoc_session_id: str = None) -> bool:
        """
        * Complete the workflow by integrating solution findings.
        
        Args:
            session_id: Original debug session ID
            solution_findings: Solution findings from Ad-Hoc agent
            ad_hoc_session_id: ID of the Ad-Hoc session
            
        Returns:
            bool: True if workflow completion was successful
        """
        success = self.debugger.integrate_solution(
            session_id, solution_findings, ad_hoc_session_id
        )
        
        if success:
            self.logger.info(f"Successfully completed APM debug workflow for session {session_id}")
        else:
            self.logger.error(f"Failed to complete APM debug workflow for session {session_id}")
        
        return success


# * Example usage and testing functions
def example_debug_function():
    """* Example debug function for testing."""
    import random
    if random.random() > 0.3:  # 70% success rate
        return "Debug successful"
    else:
        raise Exception("Simulated debug failure")


def main():
    """* Example usage of the Ad-Hoc Debugger system."""
    # * Setup logging
    logging.basicConfig(level=logging.INFO)
    
    # * Create workflow
    workflow = DebugWorkflow()
    
    # * Example issue
    issue_id = "test_issue_001"
    title = "Test Issue"
    description = "This is a test issue for the Ad-Hoc debugger system"
    error_context = {"error": "Test error message", "stack_trace": "Test stack trace"}
    
    # * Handle the issue
    success, message, delegation_prompt = workflow.handle_issue(
        issue_id, title, description, error_context, example_debug_function
    )
    
    print(f"Workflow result: {success}")
    print(f"Message: {message}")
    
    if delegation_prompt:
        print("\nDelegation Prompt:")
        print(delegation_prompt)
        
        # * Simulate Ad-Hoc session completion
        solution_findings = """
        Root cause analysis: The issue was caused by a configuration mismatch.
        
        Solution:
        1. Update configuration file with correct settings
        2. Restart the service
        3. Verify the fix works
        
        Testing: Solution has been tested and verified.
        Prevention: Added configuration validation to prevent future occurrences.
        """
        
        # * Complete the workflow
        workflow.complete_workflow(workflow.debugger.active_sessions[list(workflow.debugger.active_sessions.keys())[0]].session_id, solution_findings)


if __name__ == "__main__":
    main()
