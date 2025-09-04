#!/usr/bin/env python3
"""Test script to isolate syntax issues."""

# Test the problematic line
config_calendar_day_font_size = "\\large"
current_day = 1
day = 0
y_pos = 4.5

# This should work
test_line = f"    \\node[font=\\bfseries{{{config_calendar_day_font_size}}}, anchor=north west] at ({{day+0.05}},{{y_pos+0.4}}) {{{current_day}}};\\n"
print("Test line works:", test_line)

# Test the other problematic line
task_category_color = "blue"
task_name = "Test Task"
start_x = 0
end_x = 1
y_pos = 1

test_line2 = f"    \\draw[fill={task_category_color}, rounded corners=2pt] ({start_x},{y_pos-0.2}) rectangle ({end_x},{y_pos+0.2}) node[midway, white, font=\\small\\bfseries] {{{task_name}}};\\n"
print("Test line 2 works:", test_line2)
