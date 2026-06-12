import json
import struct
import hashlib
from typing import Tuple, List, Dict

try:
    import nacl.signing
except ImportError:
    raise ImportError("Please install PyNaCl to use the compiler: pip install pynacl")

def generate_keypair() -> Tuple[bytes, bytes]:
    """Generate Ed25519 signing key and verifying key."""
    signing_key = nacl.signing.SigningKey.generate()
    return signing_key.encode(), signing_key.verify_key.encode()

def build_capability_bitmap(agents: List[str], tools: List[str], rules: List[Dict]) -> bytes:
    """
    Build a flattened 1D capability bitmap.
    matrix[agent_idx][tool_idx] = 1 (ALLOW) or 0 (DENY).
    """
    num_agents = len(agents)
    num_tools = len(tools)
    total_bits = num_agents * num_tools
    total_bytes = (total_bits + 7) // 8
    bitmap = bytearray(total_bytes)
    
    agent_idx_map = {a: i for i, a in enumerate(agents)}
    tool_idx_map = {t: i for i, t in enumerate(tools)}

    for rule in rules:
        agent = rule.get("agent")
        tool = rule.get("tool")
        effect = rule.get("effect", "DENY")
        
        if agent in agent_idx_map and tool in tool_idx_map:
            a_idx = agent_idx_map[agent]
            t_idx = tool_idx_map[tool]
            bit_offset = a_idx * num_tools + t_idx
            byte_idx = bit_offset // 8
            bit_in_byte = bit_offset % 8
            
            if effect == "ALLOW":
                bitmap[byte_idx] |= (1 << bit_in_byte)
            elif effect == "DENY":
                bitmap[byte_idx] &= ~(1 << bit_in_byte)

    return bytes(bitmap)

def pack_string_list(strings: List[str]) -> bytes:
    """Pack a list of strings as null-terminated bytes."""
    return b"".join(s.encode('utf-8') + b'\x00' for s in strings)

def emit_binary_v3(epoch: str, dsl_version: int, compiler_version: str, agents: List[str], tools: List[str], rules: List[Dict], signing_key_seed: bytes) -> bytes:
    """
    Emit the sealed .bin binary package format V3:
     - 6 bytes Header: "VAJRAC"
     - 1 byte Version: 0x03
     - 32 bytes Epoch ID
     - 64 bytes Ed25519 Signature
     - 32 bytes SHA-256 Policy Hash
     - 1 byte DSL Version
     - 16 bytes Compiler Version
     - 4 bytes Payload Length (Big Endian uint32)
     --- Payload Starts Here ---
     - 2 bytes Num Agents (uint16)
     - 2 bytes Num Tools (uint16)
     - Variable length: Null-terminated Agent IDs
     - Variable length: Null-terminated Tool Names
     - Variable length: Bitmap
    """
    signing_key = nacl.signing.SigningKey(signing_key_seed)
    
    epoch_bytes = epoch.encode('utf-8')[:32].ljust(32, b'\x00')
    cv_bytes = compiler_version.encode('utf-8')[:16].ljust(16, b'\x00')
    dv_byte = struct.pack('>B', dsl_version)
    
    # 1. Build Payload (Agents and Tools MUST be sorted alphabetically)
    # Sorting is done in cli.py, but we ensure it here implicitly as they are passed in.
    num_agents = len(agents)
    num_tools = len(tools)
    
    payload = bytearray()
    payload.extend(struct.pack('>H', num_agents))
    payload.extend(struct.pack('>H', num_tools))
    
    payload.extend(pack_string_list(agents))
    payload.extend(pack_string_list(tools))
    
    bitmap = build_capability_bitmap(agents, tools, rules)
    payload.extend(bitmap)
    
    payload_bytes = bytes(payload)
    payload_len = len(payload_bytes)
    
    # 2. Hash Payload
    policy_hash = hashlib.sha256(payload_bytes).digest()
    
    # 3. Sign Payload
    signed = signing_key.sign(payload_bytes)
    signature = signed.signature
    
    # 4. Build Header
    header = b'VAJRAC'
    version = b'\x03'
    len_bytes = struct.pack('>I', payload_len)
    
    return header + version + epoch_bytes + signature + policy_hash + dv_byte + cv_bytes + len_bytes + payload_bytes
