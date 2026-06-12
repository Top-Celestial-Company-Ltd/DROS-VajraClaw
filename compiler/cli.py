import argparse
import yaml
import sys
import os
import binascii
from emitter import emit_binary_v3, generate_keypair

COMPILER_VERSION = "1.0.0"

def expand_wildcard(pattern: str, all_tools: list) -> list:
    if not pattern.endswith('.*'):
        return [pattern] if pattern in all_tools else []
    
    prefix = pattern[:-1]
    return [t for t in all_tools if t.startswith(prefix)]

def parse_policy_graph(policy: dict):
    linter_issues = [] # [{"severity": "...", "msg": "..."}]
    
    vajra_version = policy.get("vajra_version")
    if vajra_version != 1:
        linter_issues.append({"severity": "FATAL", "msg": f"Missing or unsupported 'vajra_version'. Expected 1, got {vajra_version}."})
    
    # Deterministic Extraction (Sort everything)
    agents_list = sorted([a["id"] for a in policy.get("agents", []) if "id" in a])
    tools_dict = {t["name"]: t for t in policy.get("tools", []) if "name" in t}
    tools_list = sorted(list(tools_dict.keys()))
    
    capabilities = policy.get("capabilities", {})
    
    flat_rules = []
    all_used_capabilities = set()
    all_granted_caps = {c for caps in capabilities.values() for c in caps}
    
    # Capability Duplicate Meaning Check (Graph Analysis)
    cap_to_tools = {c: set() for c in all_granted_caps}
    for t_name, t_def in tools_dict.items():
        req_caps = t_def.get("requires", [])
        for rc in req_caps:
            if rc in cap_to_tools:
                cap_to_tools[rc].add(t_name)
    
    seen_tool_sets = {}
    for cap, t_set in cap_to_tools.items():
        if not t_set:
            continue
        t_set_froz = frozenset(t_set)
        if t_set_froz in seen_tool_sets:
            linter_issues.append({"severity": "WARN", "msg": f"Capability Duplicate Meaning: '{cap}' and '{seen_tool_sets[t_set_froz]}' map to the exact same set of tools."})
        else:
            seen_tool_sets[t_set_froz] = cap

    # Resolve Capabilities
    for agent_id in agents_list:
        agent_caps = capabilities.get(agent_id, [])
        for t_name in tools_list:
            t_def = tools_dict[t_name]
            req_caps = t_def.get("requires", [])
            for rc in req_caps:
                all_used_capabilities.add(rc)
                
            if req_caps and all(rc in agent_caps for rc in req_caps):
                flat_rules.append({
                    "agent": agent_id,
                    "tool": t_name,
                    "effect": "ALLOW",
                    "source": "CAPABILITY"
                })

    for agent_id, agent_caps in capabilities.items():
        if agent_id not in agents_list:
            linter_issues.append({"severity": "FATAL", "msg": f"Unknown agent '{agent_id}' in capabilities block."})
            continue
        for cap in agent_caps:
            if cap not in all_used_capabilities:
                linter_issues.append({"severity": "INFO", "msg": f"Unused capability '{cap}' assigned to '{agent_id}'."})

    for t_name, t_def in tools_dict.items():
        req_caps = t_def.get("requires", [])
        if req_caps and not all(rc in all_granted_caps for rc in req_caps):
            linter_issues.append({"severity": "WARN", "msg": f"Unreachable Tool '{t_name}'. It requires capabilities nobody has."})

    explicit_rule_pairs = {}
    for r in policy.get("rules", []):
        match = r.get("match", {})
        effect = r.get("effect", "DENY")
        agent_pattern = match.get("agent")
        tool_pattern = match.get("tool")
        
        agents_to_apply = [agent_pattern] if agent_pattern in agents_list else []
        tools_to_apply = expand_wildcard(tool_pattern, tools_list)
        
        if not tools_to_apply:
            linter_issues.append({"severity": "WARN", "msg": f"Wildcard/Tool '{tool_pattern}' matched ZERO tools in the ecosystem."})
            
        for a in agents_to_apply:
            for t in tools_to_apply:
                key = (a, t)
                if key in explicit_rule_pairs and explicit_rule_pairs[key] != effect:
                    linter_issues.append({"severity": "ERROR", "msg": f"Rule Conflict: Multiple explicit rules for ({a}, {t}) resolving to both ALLOW and DENY."})
                explicit_rule_pairs[key] = effect
                
                flat_rules.append({
                    "agent": a,
                    "tool": t,
                    "effect": effect,
                    "source": "EXPLICIT_RULE"
                })
                
                # CRITICAL Check
                if effect == "ALLOW" and t.startswith("admin.") and "admin" not in a.lower() and "maintenance" not in a.lower():
                    linter_issues.append({"severity": "CRITICAL", "msg": f"Dangerous Grant: Non-admin agent '{a}' granted access to '{t}'."})
                
    for r in flat_rules:
        if r["source"] == "CAPABILITY" and explicit_rule_pairs.get((r["agent"], r["tool"])) == "DENY":
            linter_issues.append({"severity": "ERROR", "msg": f"Capability Conflict: Capability granted ALLOW for ({r['agent']}, {r['tool']}), but explicit rule overrides with DENY."})

    return agents_list, tools_list, flat_rules, linter_issues


def report_issues(issues, fail_on_error=False):
    if not issues:
        print("[Vajra Linter] [OK] Policy is clean. No logical conflicts or dead paths found.")
        return True
        
    print(f"[Vajra Linter] Found {len(issues)} potential issues:")
    has_error = False
    for i in issues:
        sev = i["severity"]
        print(f"  [{sev}] {i['msg']}")
        if sev in ["ERROR", "CRITICAL", "FATAL"]:
            has_error = True
            
    if fail_on_error and has_error:
        print("\n[Vajra Compiler] Build aborted due to ERROR/CRITICAL/FATAL issues.")
        sys.exit(1)
        
    return not has_error

def cmd_lint(input_yaml: str):
    print(f"\n[Vajra Linter] Analyzing policy: {input_yaml}")
    if not os.path.exists(input_yaml):
        print(f"Error: {input_yaml} not found.")
        sys.exit(1)
    with open(input_yaml, 'r', encoding='utf-8') as f:
        policy = yaml.safe_load(f)
    _, _, _, issues = parse_policy_graph(policy)
    report_issues(issues, fail_on_error=False)
    print("")

def cmd_doctor(input_yaml: str):
    print(f"\n[Vajra Doctor] Health Check: {input_yaml}")
    if not os.path.exists(input_yaml):
        print(f"Error: {input_yaml} not found.")
        sys.exit(1)
    with open(input_yaml, 'r', encoding='utf-8') as f:
        policy = yaml.safe_load(f)
        
    agents, tools, flat_rules, issues = parse_policy_graph(policy)
    
    total_cells = len(agents) * len(tools)
    # Count unique allows
    allows = set()
    for r in flat_rules:
        if r["effect"] == "ALLOW":
            allows.add((r["agent"], r["tool"]))
    
    # Remove overridden denies
    for r in flat_rules:
        if r["effect"] == "DENY" and r["source"] == "EXPLICIT_RULE":
            allows.discard((r["agent"], r["tool"]))
            
    num_allows = len(allows)
    density = (num_allows / total_cells * 100) if total_cells > 0 else 0
    
    unused_caps = sum(1 for i in issues if i["severity"] == "INFO" and "Unused capability" in i["msg"])
    conflict_risk = "High" if any(i["severity"] in ["ERROR", "CRITICAL"] for i in issues) else "Low"
    
    score = "A"
    if density < 1 or conflict_risk == "High":
        score = "C"
    elif unused_caps > 0:
        score = "B"
        
    print(f"  Policy Complexity Score: {score}")
    print(f"  Unused Capabilities: {unused_caps}")
    print(f"  Sparse Matrix: {'Normal' if density >= 1 else 'Sparse'} (Density: {density:.2f}%)")
    print(f"  Conflict Risk: {conflict_risk}")
    
    if density < 1:
        print("  [WARN] Sparse Matrix Detected. Consider Capability Refactor to avoid IAM explosion.")
    print("")

def cmd_build(input_yaml: str, output_bin: str, key_hex: str = None):
    print(f"\n[Vajra Compiler] Compiling policy from: {input_yaml}")
    if not os.path.exists(input_yaml):
        print(f"Error: {input_yaml} not found.")
        sys.exit(1)

    with open(input_yaml, 'r', encoding='utf-8') as f:
        policy = yaml.safe_load(f)
        
    epoch = str(policy.get("epoch", "UNVERSIONED_EPOCH"))
    dsl_version = policy.get("vajra_version", 0)
    
    agents_list, tools_list, flat_rules, issues = parse_policy_graph(policy)
            
    print(f"[Vajra Compiler] Discovered {len(agents_list)} Agents, {len(tools_list)} Tools.")
    
    # Abort on error
    report_issues(issues, fail_on_error=True)

    if key_hex:
        signing_key_seed = binascii.unhexlify(key_hex)
        import nacl.signing
        signing_key = nacl.signing.SigningKey(signing_key_seed)
        verify_key = signing_key.verify_key.encode()
        print(f"[Vajra Compiler] Using deterministic static key. PubKey: {verify_key.hex()}")
    else:
        signing_key_seed, verify_key = generate_keypair()
        print(f"[Vajra Compiler] Generated ephemeral signing key. PubKey: {verify_key.hex()}")
    
    try:
        binary_payload = emit_binary_v3(epoch, dsl_version, COMPILER_VERSION, agents_list, tools_list, flat_rules, signing_key_seed)
        
        with open(output_bin, 'wb') as f:
            f.write(binary_payload)
        
        print(f"[Vajra Compiler] Successfully emitted V3 execution artifact: {output_bin}")
        print(f"[Vajra Compiler] Binary Size: {len(binary_payload)} bytes.\n")
    except Exception as e:
        print(f"[Vajra Compiler] Compilation failed: {e}")
        sys.exit(1)

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Vajra Execution Policy Toolkit")
    parser.add_argument("command", choices=["build", "lint", "doctor"], help="Command to run")
    parser.add_argument("input", help="Path to Vajra.md (YAML) policy")
    parser.add_argument("-o", "--output", default="policy.bin", help="Output binary file path (for build)")
    parser.add_argument("--key", help="Hex string of 32-byte deterministic seed key (for reproducible builds)")
    
    args = parser.parse_args()
    
    if args.command == "build":
        cmd_build(args.input, args.output, args.key)
    elif args.command == "lint":
        cmd_lint(args.input)
    elif args.command == "doctor":
        cmd_doctor(args.input)
