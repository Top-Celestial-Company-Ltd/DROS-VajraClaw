from typing import Dict, Any

def transform_policy(parsed: Dict[str, Any]) -> Dict[str, Any]:
    """
    Transforms shorthand policies into the canonical VajraClaw DynamicPolicy AST.
    """
    epoch = parsed.get("epoch", "UNVERSIONED")
    tool_policies = []
    
    # Handle shorthand 'allowed_tools' list
    if "allowed_tools" in parsed:
        for tool in parsed["allowed_tools"]:
            tool_policies.append({
                "tool_name": tool,
                "action": "ALLOW",
                "is_write_action": False,
                "conditions": []
            })
            
    # Handle explicit 'tool_policies' AST
    if "tool_policies" in parsed:
        for tp in parsed["tool_policies"]:
            tool_policies.append(tp)
            
    return {
        "epoch": epoch,
        "tool_policies": tool_policies
    }
