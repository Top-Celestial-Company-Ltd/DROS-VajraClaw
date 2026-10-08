"""
DROS-Muse Hacker runtime hook.

This module intercepts Muse Gadget Executor.run and evaluates
the DROS Hacker policy before the original executor is called.

Epistemic boundary:
- Runtime hook integration for the Hacker package.
- Not a claim of complete mediation.
- Not a claim of raw syscall bypass resistance.
- Not a production security certification.
"""

import json
import os
import time

from dros_guard import DrosHackerGuard


EVENT_FILE = os.environ.get(
    "DROS_MUSE_RUNTIME_EVENTS",
    "/tmp/dros_muse_runtime_events.jsonl",
)


def record(event, **fields):
    payload = {
        "ts_ns": time.time_ns(),
        "event": event,
        **fields,
    }

    try:
        with open(EVENT_FILE, "a", encoding="utf-8") as f:
            f.write(json.dumps(payload, sort_keys=True) + "\n")
    except Exception:
        # Audit recording must never become an execution bypass.
        pass


record("sitecustomize_loaded")


try:
    from musegadget.executor import Executor

    record(
        "executor_imported",
        executor_module="musegadget.executor",
    )

    _original_run = Executor.run

    def dros_intercepted_run(self, *args, **kwargs):
        """
        DROS enforcement point.

        The Hacker PoC expects the following keyword inputs:

            principal
            capability
            payload
            posture
            ttl_valid

        If these are unavailable, the request is denied rather than
        silently bypassing the policy.
        """

        principal = kwargs.get("principal")
        capability = kwargs.get("capability")
        payload = kwargs.get("payload")
        posture = kwargs.get("posture")
        ttl_valid = kwargs.get("ttl_valid")

        if not isinstance(payload, dict):
            payload = {}

        if not isinstance(principal, str):
            principal = ""

        if not isinstance(capability, str):
            capability = ""

        if not isinstance(posture, str):
            posture = ""

        if not isinstance(ttl_valid, bool):
            ttl_valid = False

        decision = DrosHackerGuard.evaluate(
            principal=principal,
            capability=capability,
            payload=payload,
            posture=posture,
            ttl_valid=ttl_valid,
        )

        record(
            "dros_policy_decision",
            decision=decision.decision,
            reason=decision.reason,
            principal=principal,
            capability=capability,
        )

        if decision.decision != "ALLOW":
            record(
                "dros_execution_denied",
                reason=decision.reason,
            )

            raise PermissionError(
                f"DROS-Muse Hacker DENY: {decision.reason}"
            )

        record(
            "dros_execution_allowed",
            reason=decision.reason,
        )

        result = _original_run(self, *args, **kwargs)

        record(
            "executor_run_completed",
        )

        return result

    Executor.run = dros_intercepted_run

    record("executor_run_wrapped")

except Exception as exc:
    record(
        "executor_hook_install_failed",
        error=repr(exc),
    )

    # Fail closed.
    raise
