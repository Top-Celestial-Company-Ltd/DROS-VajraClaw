import yaml
import json
from typing import Dict, Any

def parse_policy(content: str) -> Dict[str, Any]:
    """Parse YAML or JSON policy file into a dictionary."""
    try:
        return json.loads(content)
    except json.JSONDecodeError:
        try:
            return yaml.safe_load(content)
        except yaml.YAMLError as e:
            raise ValueError(f"Failed to parse policy as JSON or YAML: {e}")
