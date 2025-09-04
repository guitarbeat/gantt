#!/usr/bin/env python3
"""
Unified configuration management for the LaTeX Gantt chart generator.
Combines core configuration classes with YAML-based device profiles and templates.
"""

import yaml
import logging
from dataclasses import dataclass, field
from pathlib import Path
from typing import Dict, List, Tuple, Any, Optional


@dataclass
class ColorScheme:
    """Color definitions for different task categories and statuses."""
    
    # Task category colors
    milestone: Tuple[int, int, int] = (147, 51, 234)
    researchcore: Tuple[int, int, int] = (59, 130, 246)
    researchexp: Tuple[int, int, int] = (16, 185, 129)
    researchout: Tuple[int, int, int] = (245, 158, 11)
    administrative: Tuple[int, int, int] = (107, 114, 128)
    accountability: Tuple[int, int, int] = (139, 92, 246)
    service: Tuple[int, int, int] = (236, 72, 153)
    other: Tuple[int, int, int] = (156, 163, 175)
    
    # Status colors
    completed: Tuple[int, int, int] = (34, 197, 94)
    inprogress: Tuple[int, int, int] = (251, 146, 60)
    blocked: Tuple[int, int, int] = (239, 68, 68)
    planned: Tuple[int, int, int] = (59, 130, 246)
    
    def to_latex_colors(self) -> str:
        """Generate LaTeX color definitions."""
        colors = []
        for attr_name in dir(self):
            if not attr_name.startswith('_') and not callable(getattr(self, attr_name)):
                rgb = getattr(self, attr_name)
                colors.append(f"\\definecolor{{{attr_name}}}{{RGB}}{{{rgb[0]}, {rgb[1]}, {rgb[2]}}}")
        return '\n'.join(colors)


@dataclass
class CalendarConfig:
    """Configuration for calendar layout and styling."""
    
    # Page layout
    page_orientation: str = "landscape"
    page_size: str = "a4paper"
    margin: str = "0.5in"
    
    # Calendar grid
    calendar_scale: float = 1.0
    calendar_width: int = 7
    calendar_height: int = 6
    max_tasks_per_day: int = 2
    max_task_name_length: int = 18
    
    # Typography
    title_font_size: str = "\\LARGE"
    month_font_size: str = "\\large"
    day_font_size: str = "\\large"
    task_font_size: str = "\\tiny"
    
    # Spacing
    title_spacing: str = "1cm"
    month_spacing: str = "0.5cm"
    task_spacing: str = "0.5cm"


@dataclass
class TaskConfig:
    """Configuration for task processing and categorization."""
    
    # Task name cleaning patterns
    name_prefixes_to_remove: List[str] = field(default_factory=lambda: [
        "Milestone: ",
        "Draft ",
        "Complete "
    ])
    
    # Category mapping
    category_keywords: Dict[str, List[str]] = field(default_factory=lambda: {
        "researchcore": ["PROPOSAL"],
        "researchexp": ["LASER", "EXPERIMENTAL"],
        "researchout": ["PUBLICATION", "PRESENTATION"],
        "administrative": ["ADMINISTRATIVE"],
        "accountability": ["ACCOUNTABILITY"],
        "service": ["SERVICE"]
    })
    
    # Date format
    date_format: str = "%Y-%m-%d"
    display_date_format: str = "%m/%d"


@dataclass
class LaTeXConfig:
    """Configuration for LaTeX document generation."""
    
    # Document class and packages
    document_class: str = "article"
    packages: List[str] = field(default_factory=lambda: [
        "[utf8]{inputenc}",
        "[T1]{fontenc}",
        "lmodern",
        "helvet",
        "[landscape,margin=0.5in]{geometry}",
        "tikz",
        "pgfplots",
        "xcolor",
        "enumitem",
        "booktabs",
        "array",
        "longtable",
        "fancyhdr",
        "graphicx",
        "amsmath",
        "amsfonts",
        "amssymb",
        "ragged2e"
    ])
    
    # Document metadata
    default_title: str = "Project Timeline"
    subtitle: str = "PhD Research Calendar"
    
    # * Enhanced TikZ libraries for better functionality
    tikz_libraries: List[str] = field(default_factory=lambda: [
        "arrows.meta",           # Better arrows
        "shapes.geometric",      # More shapes (diamonds, stars, etc.)
        "positioning",           # Relative positioning
        "calc",                  # Coordinate calculations
        "decorations.pathmorphing", # Decorative paths
        "patterns",              # Fill patterns
        "shadows",               # Drop shadows
        "fit",                   # Fit nodes to content
        "backgrounds",           # Background layers
        "matrix",                # Matrix layouts
        "chains",                # Node chains
        "scopes",                # Scoped operations
        "pgfgantt"               # Gantt charts (from awesome-tikz)
    ])
    
    def get_tikz_libraries(self) -> List[str]:
        """Get list of TikZ libraries to load."""
        return self.tikz_libraries


@dataclass
class DeviceProfile:
    """Device-specific configuration profile."""
    
    name: str
    description: str
    device_type: str
    optimizations: List[str] = field(default_factory=list)
    layout: Dict[str, Any] = field(default_factory=dict)
    colors: Dict[str, List[int]] = field(default_factory=dict)
    
    def get_color_rgb(self, color_name: str) -> tuple[int, int, int]:
        """Get RGB color tuple for a color name."""
        if color_name in self.colors:
            rgb = self.colors[color_name]
            return tuple(rgb)
        return (0, 0, 0)  # Default to black
    
    def get_layout_value(self, key: str, default: Any = None) -> Any:
        """Get layout configuration value."""
        return self.layout.get(key, default)


@dataclass
class TemplateConfig:
    """Template configuration."""
    
    name: str
    description: str
    layout: str
    orientation: str
    page_size: str
    margin: str
    features: List[str] = field(default_factory=list)


@dataclass
class ColorSchemeConfig:
    """Color scheme configuration."""
    
    name: str
    description: str
    colors: Dict[str, List[int]] = field(default_factory=dict)
    
    def get_color_rgb(self, color_name: str) -> tuple[int, int, int]:
        """Get RGB color tuple for a color name."""
        if color_name in self.colors:
            rgb = self.colors[color_name]
            return tuple(rgb)
        return (0, 0, 0)  # Default to black


class ConfigManager:
    """Unified configuration manager with YAML support and device profiles."""
    
    def __init__(self, config_dir: Optional[Path] = None):
        """Initialize configuration manager."""
        self.logger = logging.getLogger(__name__)
        self.config_dir = config_dir or Path("src/config")
        self.templates: Dict[str, TemplateConfig] = {}
        self.device_profiles: Dict[str, DeviceProfile] = {}
        self.color_schemes: Dict[str, ColorSchemeConfig] = {}
        self.defaults: Dict[str, str] = {}
        
        # Load configurations
        self._load_configurations()
    
    def _load_configurations(self) -> None:
        """Load all configuration files."""
        try:
            self._load_templates()
            self._load_device_profiles()
            self._load_color_schemes()
            self.logger.info("Configuration files loaded successfully")
        except Exception as e:
            self.logger.error(f"Error loading configurations: {e}")
            self._load_defaults()
    
    def _load_templates(self) -> None:
        """Load template configurations from YAML."""
        templates_file = self.config_dir / "templates.yaml"
        if not templates_file.exists():
            self.logger.warning(f"Templates file not found: {templates_file}")
            return
        
        with open(templates_file, 'r', encoding='utf-8') as f:
            data = yaml.safe_load(f)
        
        if 'templates' in data:
            for template_id, template_data in data['templates'].items():
                self.templates[template_id] = TemplateConfig(
                    name=template_data.get('name', template_id),
                    description=template_data.get('description', ''),
                    layout=template_data.get('layout', 'vertical'),
                    orientation=template_data.get('orientation', 'portrait'),
                    page_size=template_data.get('page_size', 'a4paper'),
                    margin=template_data.get('margin', '0.5in'),
                    features=template_data.get('features', [])
                )
        
        if 'defaults' in data:
            self.defaults = data['defaults']
    
    def _load_device_profiles(self) -> None:
        """Load device profile configurations from YAML."""
        profiles_file = self.config_dir / "device_profiles.yaml"
        if not profiles_file.exists():
            self.logger.warning(f"Device profiles file not found: {profiles_file}")
            return
        
        with open(profiles_file, 'r', encoding='utf-8') as f:
            data = yaml.safe_load(f)
        
        if 'profiles' in data:
            for profile_id, profile_data in data['profiles'].items():
                self.device_profiles[profile_id] = DeviceProfile(
                    name=profile_data.get('name', profile_id),
                    description=profile_data.get('description', ''),
                    device_type=profile_data.get('device_type', 'digital'),
                    optimizations=profile_data.get('optimizations', []),
                    layout=profile_data.get('layout', {}),
                    colors=profile_data.get('colors', {})
                )
    
    def _load_color_schemes(self) -> None:
        """Load color scheme configurations from YAML."""
        templates_file = self.config_dir / "templates.yaml"
        if not templates_file.exists():
            return
        
        with open(templates_file, 'r', encoding='utf-8') as f:
            data = yaml.safe_load(f)
        
        if 'color_schemes' in data:
            for scheme_id, scheme_data in data['color_schemes'].items():
                self.color_schemes[scheme_id] = ColorSchemeConfig(
                    name=scheme_data.get('name', scheme_id),
                    description=scheme_data.get('description', ''),
                    colors=scheme_data.get('colors', {})
                )
    
    def _load_defaults(self) -> None:
        """Load default configurations when YAML files are not available."""
        self.logger.info("Loading default configurations")
        
        # Default template
        self.templates['gantt_timeline'] = TemplateConfig(
            name="Gantt Timeline",
            description="Vertical timeline with task bars and dependencies",
            layout="vertical",
            orientation="portrait",
            page_size="a4paper",
            margin="0.5in",
            features=["task_bars", "dependencies", "milestones"]
        )
        
        # Default device profile
        self.device_profiles['standard_print'] = DeviceProfile(
            name="Standard Print",
            description="Optimized for standard office/home printing",
            device_type="print",
            optimizations=["print_safe_colors", "standard_margins"],
            layout={
                "page_size": "a4paper",
                "orientation": "portrait",
                "margin": "0.5in",
                "font_size": "10pt",
                "line_thickness": "0.8pt"
            },
            colors={
                "background": [255, 255, 255],
                "text": [0, 0, 0],
                "accent": [59, 130, 246],
                "grid": [200, 200, 200]
            }
        )
        
        # Default color scheme
        self.color_schemes['academic'] = ColorSchemeConfig(
            name="Academic",
            description="Professional colors for academic and research use",
            colors={
                "primary": [59, 130, 246],
                "secondary": [16, 185, 129],
                "accent": [245, 158, 11],
                "neutral": [107, 114, 128],
                "success": [34, 197, 94],
                "warning": [251, 146, 60],
                "error": [239, 68, 68],
                "info": [59, 130, 246]
            }
        )
        
        # Default settings
        self.defaults = {
            "template": "gantt_timeline",
            "device": "standard_print",
            "color_scheme": "academic"
        }
    
    def get_template(self, template_id: Optional[str] = None) -> TemplateConfig:
        """Get template configuration."""
        if template_id is None:
            template_id = self.defaults.get('template', 'gantt_timeline')
        
        if template_id in self.templates:
            return self.templates[template_id]
        
        self.logger.warning(f"Template '{template_id}' not found, using default")
        return self.templates.get('gantt_timeline', list(self.templates.values())[0])
    
    def get_device_profile(self, profile_id: Optional[str] = None) -> DeviceProfile:
        """Get device profile configuration."""
        if profile_id is None:
            profile_id = self.defaults.get('device', 'standard_print')
        
        if profile_id in self.device_profiles:
            return self.device_profiles[profile_id]
        
        self.logger.warning(f"Device profile '{profile_id}' not found, using default")
        return self.device_profiles.get('standard_print', list(self.device_profiles.values())[0])
    
    def get_color_scheme(self, scheme_id: Optional[str] = None) -> ColorSchemeConfig:
        """Get color scheme configuration."""
        if scheme_id is None:
            scheme_id = self.defaults.get('color_scheme', 'academic')
        
        if scheme_id in self.color_schemes:
            return self.color_schemes[scheme_id]
        
        self.logger.warning(f"Color scheme '{scheme_id}' not found, using default")
        return self.color_schemes.get('academic', list(self.color_schemes.values())[0])
    
    def list_templates(self) -> List[str]:
        """List available template IDs."""
        return list(self.templates.keys())
    
    def list_device_profiles(self) -> List[str]:
        """List available device profile IDs."""
        return list(self.device_profiles.keys())
    
    def list_color_schemes(self) -> List[str]:
        """List available color scheme IDs."""
        return list(self.color_schemes.keys())
    
    def get_combined_config(self, template_id: Optional[str] = None, 
                          device_profile_id: Optional[str] = None,
                          color_scheme_id: Optional[str] = None) -> Dict[str, Any]:
        """Get combined configuration from template, device profile, and color scheme."""
        template = self.get_template(template_id)
        device_profile = self.get_device_profile(device_profile_id)
        color_scheme = self.get_color_scheme(color_scheme_id)
        
        return {
            'template': template,
            'device_profile': device_profile,
            'color_scheme': color_scheme,
            'layout': {
                'orientation': template.orientation,
                'page_size': template.page_size,
                'margin': template.margin,
                'features': template.features
            },
            'device_optimizations': device_profile.optimizations,
            'colors': {
                'background': device_profile.get_color_rgb('background'),
                'text': device_profile.get_color_rgb('text'),
                'accent': device_profile.get_color_rgb('accent'),
                'grid': device_profile.get_color_rgb('grid'),
                'scheme': color_scheme.colors
            }
        }


@dataclass
class AppConfig:
    """Main application configuration combining all sub-configurations."""
    
    colors: ColorScheme = field(default_factory=ColorScheme)
    calendar: CalendarConfig = field(default_factory=CalendarConfig)
    tasks: TaskConfig = field(default_factory=TaskConfig)
    latex: LaTeXConfig = field(default_factory=LaTeXConfig)
    
    # File paths
    default_input: str = "../input/data.cleaned.csv"
    default_output: str = "output/tex/Calendar_template.tex"


# Global configuration instances
config = AppConfig()
config_manager = ConfigManager()


def get_category_color(category: str) -> str:
    """Get the color name for a given task category."""
    category_upper = category.upper()
    
    for color_name, keywords in config.tasks.category_keywords.items():
        if any(keyword in category_upper for keyword in keywords):
            return color_name
    
    return "other"


def get_task_marker(is_milestone: bool) -> str:
    """Get the LaTeX marker for a task based on its type."""
    return "$\\diamond$" if is_milestone else "$\\bullet$"