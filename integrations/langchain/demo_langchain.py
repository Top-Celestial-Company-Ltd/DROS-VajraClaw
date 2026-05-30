"""
VajraClaw × LangChain Integration Demo
---------------------------------------
This demo shows how to wrap a LangChain agent's tool calls
with VajraClaw enforcement. When the agent attempts an
unauthorized action, execution is physically blocked before
any system call is made.

Requirements:
    pip install langchain langchain-openai

Run:
    python demo_langchain.py
"""

import os
from typing import Any

# ── VajraClaw ─────────────────────────────────────────────────────────────────
import sys
sys.path.insert(0, os.path.join(os.path.dirname(__file__), ".."))
from vajraclaw import VajraClaw

# ── LangChain ─────────────────────────────────────────────────────────────────
from langchain.tools import BaseTool


# ── Demo: define the policy inline (no file needed to run this demo) ──────────
DEMO_POLICY = """
# VajraClaw Demo Policy
# Only READ_DB is authorized. WRITE_DB is not listed = denied.
allowed_tools:
  - db.read
  - file.read
"""

# Initialize the enforcement kernel
vc = VajraClaw(rules_string=DEMO_POLICY)


# ── Wrap ANY LangChain BaseTool with VajraClaw enforcement ───────────────────
class VajraProtectedTool(BaseTool):
    """
    A LangChain BaseTool wrapper that enforces VajraClaw capability checks
    before executing the underlying tool.

    Wrap any existing tool:
        protected = VajraProtectedTool.wrap(my_tool, agent_id="finance-agent")
    """

    name: str = "vajra_protected_tool"
    description: str = ""
    _inner: Any = None
    _agent_id: str = "default"

    @classmethod
    def wrap(cls, tool: BaseTool, agent_id: str = "default") -> "VajraProtectedTool":
        wrapped = cls()
        wrapped.name = tool.name
        wrapped.description = tool.description
        wrapped._inner = tool
        wrapped._agent_id = agent_id
        return wrapped

    def _run(self, *args, **kwargs) -> str:
        result = vc.evaluate(tool=self.name, agent_id=self._agent_id)

        if not result:
            # ❌ PHYSICAL BLOCK — print full denial report and raise
            print(result)
            raise PermissionError(
                f"[VajraClaw] Execution blocked: agent '{self._agent_id}' "
                f"is not authorized to call '{self.name}'"
            )

        print(result)
        return self._inner._run(*args, **kwargs)

    async def _arun(self, *args, **kwargs):
        return self._run(*args, **kwargs)


# ── Mock tools (simulate a real agent environment) ────────────────────────────
class MockDbReadTool(BaseTool):
    name: str = "db.read"
    description: str = "Read records from the database"

    def _run(self, query: str = "") -> str:
        return f"[DB] Read result for: {query}"

    async def _arun(self, query: str = "") -> str:
        return self._run(query)


class MockDbWriteTool(BaseTool):
    name: str = "db.write"
    description: str = "Write records to the database"

    def _run(self, data: str = "") -> str:
        return f"[DB] Written: {data}"

    async def _arun(self, data: str = "") -> str:
        return self._run(data)


# ── Main demo ─────────────────────────────────────────────────────────────────
if __name__ == "__main__":
    agent_id = "finance-agent"

    # Wrap both tools with VajraClaw enforcement
    protected_read = VajraProtectedTool.wrap(MockDbReadTool(), agent_id=agent_id)
    protected_write = VajraProtectedTool.wrap(MockDbWriteTool(), agent_id=agent_id)

    print("=" * 60)
    print("VajraClaw × LangChain — Enforcement Demo")
    print("=" * 60)

    # ── Test 1: Authorized tool (should ALLOW) ───────────────────────────
    print("\n[Test 1] Agent calls db.read (authorized)")
    try:
        output = protected_read._run(query="SELECT * FROM transactions")
        print(f"  Result: {output}")
    except PermissionError as e:
        print(f"  {e}")

    # ── Test 2: Unauthorized tool (should BLOCK) ─────────────────────────
    print("\n[Test 2] Agent calls db.write (NOT authorized)")
    try:
        output = protected_write._run(data="DROP TABLE users")
        print(f"  Result: {output}")
    except PermissionError as e:
        print(f"  → Correctly blocked. Agent cannot write to DB.")

    print("\n" + "=" * 60)
    print("Demo complete. VajraClaw enforced capability boundaries.")
    print("=" * 60)
