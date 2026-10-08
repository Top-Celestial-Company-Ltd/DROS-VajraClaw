# DROS-Muse Hacker Integration

DROS-Muse Hacker is a runtime execution governance layer for Muse Linux Gadgets.

It is installed locally on a Muse Linux Gadget and evaluates execution requests
before they reach the Muse Executor boundary.

## Relationship

This directory is the DROS-VajraClaw ecosystem integration entry for DROS-Muse Hacker.

The runtime package itself remains independently structured under:

DROS-Muse/hacker/

This integration does not modify the upstream Muse SDK executor.py.

## Verified scope

The current evidence establishes:

- Real Muse Linux runtime enforcement was verified on an isolated Linux environment.
- ALLOW execution reached the original Executor.
- DENY cases raised before the original Executor call.
- Upstream executor.py remained unchanged.
- Install and uninstall lifecycle was verified.

## Explicit non-claims

The current evidence does NOT establish:

- complete mediation of all possible execution paths
- raw syscall bypass resistance
- direct Muse Agent reachability to arbitrary Python primitives
- production security
- universal protection against all bypass techniques

See the DROS-Muse evidence package for the corresponding evidence and hashes.

## Package boundary

This integration must not be treated as a replacement for the upstream Muse SDK.

The Hacker package is an external runtime governance layer.
