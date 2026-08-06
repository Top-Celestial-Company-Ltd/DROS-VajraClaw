# DROS VajraClaw Cryptographic Key Management & Backup Guide (KEY_MANAGEMENT_EN.md)

This document defines the key lifecycle management specifications for the DROS / VajraClaw security governance policy binary contract (`policy.bin`) during the development, compilation, and deployment phases. All system administrators and security engineers must follow this guide to ensure that the zero-trust boundary is never compromised.

> [!CAUTION]
> **Core Security Declaration: No Backdoors & No Key Recovery Mechanism**
> * **No Vendor Assistance for Recovery**: DROS / VajraClaw strictly adheres to zero-trust security principles. **No master recovery keys or backdoors are built into the system**. If you lose your private key seed (Seed Hex), **the vendor (DROS) cannot recover it or sign new `policy.bin` files on your behalf**.
> * **Sole Recovery Path**: You must manually execute the "Trust Anchor Reconstruction" procedure (see Section 3, Scenario B). This requires administrative control (SSH/Root access) over your deployed servers to generate a new key pair and manually update the verification public key configuration on all execution nodes. If you lose both your private key and server access, the system will be permanently locked and cannot be upgraded.

---

## 1. Core Concepts: Ephemeral Keys vs. Static Keys

The VajraClaw compiler (`cli.py`) supports two signing modes:

| Mode | Key Generation & Lifecycle | Production Suitability | Security & Operational Risks |
| :--- | :--- | :--- | :--- |
| **Ephemeral Key** | Generated randomly in-memory during compile time, destroyed immediately after. | **🚨 Strictly Forbidden** | Each compilation yields a different public key, causing execution-side GuardVM validation failures (triggering immediate meltdown) unless public key anchors are manually redeployed on all VMs. |
| **Static Key** | Uses a fixed Ed25519 32-byte (256-bit) private key seed (Seed Hex) for signing. | **🟢 Recommended for Prod** | Public key anchors are permanently deployed to VMs. Policy updates are seamless and require only safeguarding and backing up the private key seed. |

---

## 2. Static Key Generation & Backup Strategy

### 2.1 Keypair Generation (Python)
Run the following script on an offline, secure administrator workstation to generate your Ed25519 key pair:
```python
import nacl.signing
import binascii

# Generate a random signing key seed
signing_key = nacl.signing.SigningKey.generate()
seed_hex = binascii.hexlify(signing_key.encode()).decode('utf-8')
pub_hex = binascii.hexlify(signing_key.verify_key.encode()).decode('utf-8')

print(f"[BACK UP THIS] Private Key Seed (32-byte Seed Hex): {seed_hex}")
print(f"[DEPLOY THIS] Verification Public Key (PubKey Hex)  : {pub_hex}")
```

### 2.2 Private Key Seed Backup Best Practices
1. **Never Commit to Git**: Ensure that any configuration files or scripts containing `seed_hex` are listed in `.gitignore`. Never push them to public or private repositories.
2. **Cold Backup Storage**: Store the private key seed as a "Secure Note" in an enterprise-grade password manager (e.g., 1Password, KeePass, or Bitwarden).
3. **Key Separation Deployment**:
   * **Private Key Seed**: Kept strictly on the secure signing/compilation machine, used only to sign `policy.bin`.
   * **Verification Public Key**: Hard-coded or configured on the execution-side VM / edge devices. Losing the public key does not compromise security.

---

## 3. Disaster Recovery SOP

### 🚨 Scenario A: Build Workstation Failure, but Private Key Seed is Backed Up
In this scenario, the public key on the execution side does not need to change, and the policy update is completely seamless.

1. **Environment Setup**: On a new build machine, install the dependencies:
   ```bash
   pip install pynacl pyyaml
   ```
2. **Recompile**: Retrieve the backed-up `Seed Hex` from the password manager and compile deterministically:
   ```bash
   python cli.py build rules/policy_young.yaml -o policy.bin --key <YOUR_BACKED_UP_SEED_HEX>
   ```
3. **Deploy Policy**: Overwrite `policy.bin` on the execution node. The microkernel will seamlessly verify and load the new policy.

### 🚨 Scenario B: Private Key Seed is Completely Lost (Reconstruction)
In this scenario, a new key pair must be generated, and all VM execution nodes must be updated manually.

1. **Generate New Keys**: Follow Section 2.1 to generate a new `NEW_SEED` and `NEW_PUBKEY`.
2. **Compile Policy**: Compile the policy using the new seed:
   ```bash
   python cli.py build rules/policy_young.yaml -o policy.bin --key <NEW_SEED>
   ```
3. **Update Execution Nodes**:
   * Log into your VM, and replace the old public key with `NEW_PUBKEY` in `gemini_proxy.py` or the microkernel config.
   * Securely backup the `NEW_SEED`.
4. **Deploy and Restart**:
   * Upload the new `policy.bin` to the execution node.
   * Restart the proxy/microkernel service (e.g., `sudo systemctl restart dros-proxy.service`) to apply the new public key and policy.
