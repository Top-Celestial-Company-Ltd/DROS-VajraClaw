"""
vajraclaw.runtime
-----------------
Core Python binding for the VajraClaw enforcement kernel.
Wraps vajra_claw.dll / .so via ctypes — zero external dependencies.
"""

import ctypes
import json
import os
import platform
from dataclasses import dataclass
from enum import Enum
from typing import Optional


class Decision(Enum):
    ALLOW = "ALLOW"
    BLOCK = "BLOCK"


@dataclass
class EvaluationResult:
    decision: Decision
    tool: str
    agent_id: str
    reason: str
    matched_rule: str = ""

    def __bool__(self):
        """Allows: if vc.evaluate(...): to mean 'allowed'"""
        return self.decision == Decision.ALLOW

    def __str__(self):
        icon = "✅" if self.decision == Decision.ALLOW else "❌"
        lines = [
            f"{icon} {self.decision.value}",
            f"  Agent:        {self.agent_id}",
            f"  Tool:         {self.tool}",
            f"  Reason:       {self.reason}",
        ]
        if self.matched_rule:
            lines.append(f"  Matched Rule: {self.matched_rule}")
        return "\n".join(lines)


def _resolve_binary(core_path: Optional[str]) -> str:
    """Auto-resolve the core binary path based on OS if not provided."""
    if core_path:
        return core_path

    system = platform.system()
    base = os.path.join(os.path.dirname(__file__), "..", "..", "core")
    if system == "Windows":
        return os.path.join(base, "vajra_claw.dll")
    elif system == "Darwin":
        return os.path.join(base, "vajra_claw.dylib")
    else:
        return os.path.join(base, "vajra_claw.so")


class VajraClaw:
    """
    VajraClaw Python Runtime

    Thin ctypes wrapper over the VajraClaw C-shared enforcement kernel.
    Provides evaluate() for capability-based tool call enforcement.

    Example:
        vc = VajraClaw(rules="./Vajra.md")
        result = vc.evaluate(tool="db.write", agent_id="finance-agent")
        if not result:
            raise PermissionError(str(result))
    """

    def __init__(
        self,
        core_path: Optional[str] = None,
        rules: Optional[str] = None,
        rules_string: Optional[str] = None,
    ):
        resolved = os.path.abspath(_resolve_binary(core_path))

        try:
            self._lib = ctypes.CDLL(resolved)
        except OSError as e:
            raise RuntimeError(
                f"[VajraClaw] Failed to load core binary at: {resolved}\n"
                f"Ensure the binary is compiled: go build -buildmode=c-shared\n"
                f"Original error: {e}"
            )

        # ── Wire up C-FFI signatures ──────────────────────────────────────
        self._lib.init_static_vajra.argtypes = [ctypes.c_char_p]
        self._lib.init_static_vajra.restype = ctypes.c_int

        self._lib.init_static_vajra_from_string.argtypes = [ctypes.c_char_p]
        self._lib.init_static_vajra_from_string.restype = ctypes.c_int

        self._lib.inject_ephemeral_rule.argtypes = [ctypes.c_char_p]
        self._lib.inject_ephemeral_rule.restype = ctypes.c_int

        self._lib.match_token_stream.argtypes = [ctypes.c_char_p]
        self._lib.match_token_stream.restype = ctypes.c_int

        self._lib.clear_ephemeral_rules.argtypes = []
        self._lib.clear_ephemeral_rules.restype = None

        # ── Load policy ───────────────────────────────────────────────────
        if rules_string:
            self._load_from_string(rules_string)
        elif rules:
            self._load_from_file(rules)
        else:
            raise ValueError(
                "[VajraClaw] Must provide either `rules` (file path) or `rules_string`."
            )

        print(f"[VajraClaw] ✅ Enforcement kernel ready. Binary: {resolved}")

    def _load_from_file(self, path: str):
        result = self._lib.init_static_vajra(os.path.abspath(path).encode("utf-8"))
        if result != 1:
            raise RuntimeError(f"[VajraClaw] Failed to load rules from: {path}")

    def _load_from_string(self, content: str):
        result = self._lib.init_static_vajra_from_string(content.encode("utf-8"))
        if result != 1:
            raise RuntimeError("[VajraClaw] Failed to load rules from string.")

    def evaluate(self, tool: str, agent_id: str = "default") -> EvaluationResult:
        """
        Evaluate whether an agent is allowed to invoke a tool.

        Returns an EvaluationResult which is truthy on ALLOW, falsy on BLOCK.
        """
        payload = json.dumps({"tool": tool, "agent_id": agent_id})
        raw = self._lib.match_token_stream(payload.encode("utf-8"))

        if raw == 1:
            return EvaluationResult(
                decision=Decision.ALLOW,
                tool=tool,
                agent_id=agent_id,
                reason="Capability authorized",
            )
        else:
            return EvaluationResult(
                decision=Decision.BLOCK,
                tool=tool,
                agent_id=agent_id,
                reason="Capability not authorized — tool call physically blocked",
                matched_rule="enforcement_kernel",
            )

    def inject_rule(self, rule: str):
        """Inject a one-shot ephemeral rule (JIT boundary)."""
        self._lib.inject_ephemeral_rule(rule.encode("utf-8"))

    def clear_rules(self):
        """Evaporate all ephemeral rules (call at session end)."""
        self._lib.clear_ephemeral_rules()
