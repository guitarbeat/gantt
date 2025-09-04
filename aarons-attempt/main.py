#!/usr/bin/env python3
"""
Main entry point for the LaTeX Gantt Chart Generator.
This provides access to the unified application system.
"""

import sys

def main():
    """Main entry point for the unified application system."""
    from src.core import main as core_main
    return core_main()

if __name__ == "__main__":
    sys.exit(main())
