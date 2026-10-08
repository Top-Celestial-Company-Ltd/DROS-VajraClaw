"""
DROS-Muse Hacker — deterministic runtime policy guard.

This module is intentionally small and dependency-free.

Epistemic boundary:
- Derived from the Phase 2C policy extraction.
- Independently implemented for Hacker runtime use.
- NOT claimed to be the production Muse policy engine.
- NOT claimed to provide complete mediation.
- NOT claimed to resist raw syscall bypass.
"""

from dataclasses import dataclass
from typing import Any, Dict


@dataclass(frozen=True)
class Decision:
    decision: str
    reason: str


class DrosHackerGuard:
    """
    Minimal deterministic policy guard for DROS-Muse Hacker.

    Policy order is intentional:

        TTL
        -> Runtime Posture
        -> Capability
        -> Payload
        -> Explicit Allowlist
        -> Default DENY
    """

    ALLOWED_CAPABILITIES = frozenset(
        {
            "system.run",
            "file.read",
            "file.write",
        }
    )

    @classmethod
    def evaluate(
        cls,
        principal: str,
        capability: str,
        payload: Dict[str, Any],
        posture: str,
        ttl_valid: bool,
    ) -> Decision:

        # Principal is part of the runtime authorization input.
        # This minimal PoC policy does not yet impose a principal
        # allowlist. It is deliberately retained in the interface
        # for future governance integration.
        _ = principal

        if not ttl_valid:
            return Decision("DENY", "TTL_INVALID")

        if posture != "healthy":
            return Decision("DENY", "POSTURE_NOT_HEALTHY")

        if capability not in cls.ALLOWED_CAPABILITIES:
            return Decision("DENY", "CAPABILITY_NOT_ALLOWED")

        command = str(payload.get("command", "")).strip()

        if capability == "system.run":

            if command.startswith("rm -rf /tmp/dros-deny"):
                return Decision(
                    "DENY",
                    "POLICY_DENY_DESTRUCTIVE_TEST",
                )

            if command.startswith("python -c"):
                return Decision(
                    "ALLOW",
                    "POLICY_ALLOW_PYTHON_TEST",
                )

            if command == "echo dros-test":
                return Decision(
                    "ALLOW",
                    "POLICY_ALLOW_TEST",
                )

            if command.startswith("/usr/bin/systemctl restart"):
                return Decision(
                    "ALLOW",
                    "POLICY_ALLOW_MAINTENANCE_TEST",
                )

            return Decision(
                "DENY",
                "POLICY_DEFAULT_DENY",
            )

        return Decision(
            "ALLOW",
            "POLICY_ALLOW_CAPABILITY",
        )
