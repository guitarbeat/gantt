#!/usr/bin/env python3
"""
Main entry point for the LaTeX Gantt Chart Generator.
This provides access to both the main application and build system.
"""

import sys
import argparse

def main():
    """Main entry point that routes to appropriate subsystem."""
    parser = argparse.ArgumentParser(
        description="LaTeX Gantt Chart Generator",
        add_help=False
    )
    parser.add_argument('command', nargs='?', default='app', 
                       choices=['app', 'build'],
                       help='Command to run: app (default) or build')
    
    # Parse just the first argument to determine which system to use
    args, remaining_args = parser.parse_known_args()
    
    if args.command == 'build':
        # Route to build system
        from src.build import main as build_main
        # Reconstruct sys.argv for the build system
        sys.argv = ['main.py'] + remaining_args
        return build_main()
    else:
        # Route to main application (default)
        from src.app import main as app_main
        # Reconstruct sys.argv for the app
        sys.argv = ['main.py'] + remaining_args
        return app_main()

if __name__ == "__main__":
    sys.exit(main())
