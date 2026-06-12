import sys
import os

# Add vajraclaw module to path
sys.path.append(os.path.join(os.path.dirname(__file__), '..', 'integrations'))

from vajraclaw.runtime import VajraClaw

pubkey = "3b6a27bcceb6a42d62a3a8d02a6f0d73653215771de243a63ac048a18b59da29"
policy_bin = "policy.bin"

print("[Test] Initializing VajraClaw with V2 Binary Policy...")
try:
    vc = VajraClaw(binary_policy=policy_bin, pubkey=pubkey, license_key="VAJRA-COMMERCIAL-9999")
    print("[Test] Success: Engine initialized!")
except Exception as e:
    print(f"[Test] Failed: {e}")
    sys.exit(1)

# Test db.query (should be ALLOWED per rule)
res1 = vc._lib.evaluate_dynamic_tool_call_with_audit(
    b"db.query", b"{}", b"maintenance-agent", b""
)
import ctypes
print("Evaluate db.query ->", ctypes.string_at(res1).decode('utf-8'))

# Test sys.reboot (should be DENIED per rule)
res2 = vc._lib.evaluate_dynamic_tool_call_with_audit(
    b"sys.reboot", b"{}", b"maintenance-agent", b""
)
print("Evaluate sys.reboot ->", ctypes.string_at(res2).decode('utf-8'))

# Test file.write (not matched, defaults to DENIED)
res3 = vc._lib.evaluate_dynamic_tool_call_with_audit(
    b"file.write", b"{}", b"maintenance-agent", b""
)
print("Evaluate file.write ->", ctypes.string_at(res3).decode('utf-8'))

