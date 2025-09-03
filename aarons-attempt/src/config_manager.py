#!/usr/bin/env python3
"""
Enhanced configuration management system inspired by latex-yearly-planner.
Provides YAML-based configuration with device profiles and template options.
"""

import yaml
import logging
from pathlib import Path
from typing import Dict, Any, Optional, List
from dataclasses import dataclass, field
from .config import config as base_config


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
class ColorScheme:
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
    """Enhanced configuration manager with YAML support and device profiles."""
    
    def __init__(self, config_dir: Optional[Path] = None):
        """Initialize configuration manager."""
        self.logger = logging.getLogger(__name__)
        self.config_dir = config_dir or Path("src/config")
        self.templates: Dict[str, TemplateConfig] = {}
        self.device_profiles: Dict[str, DeviceProfile] = {}
        self.color_schemes: Dict[str, ColorScheme] = {}
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
                self.color_schemes[scheme_id] = ColorScheme(
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
        self.color_schemes['academic'] = ColorScheme(
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
    
    def get_color_scheme(self, scheme_id: Optional[str] = None) -> ColorScheme:
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


# Global configuration manager instance
config_manager = ConfigManager()
