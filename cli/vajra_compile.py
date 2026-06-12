import argparse
import os
import sys

# Add parent dir to path to import compiler
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from compiler.parser import parse_policy
from compiler.transformer import transform_policy
from compiler.emitter import emit_binary, generate_keypair

def main():
    parser = argparse.ArgumentParser(description="VajraClaw Policy Compiler")
    parser.add_argument("--policy", required=True, help="Input policy file (YAML/JSON)")
    parser.add_argument("--out", required=True, help="Output binary file (.bin)")
    parser.add_argument("--key", help="Optional Ed25519 private key seed (hex). If not provided, a new one is generated.")
    parser.add_argument("--epoch", help="Override policy epoch ID")
    
    args = parser.parse_args()
    
    with open(args.policy, 'r', encoding='utf-8') as f:
        content = f.read()
        
    parsed = parse_policy(content)
    if args.epoch:
        parsed['epoch'] = args.epoch
        
    ast = transform_policy(parsed)
    
    if args.key:
        seed = bytes.fromhex(args.key)
    else:
        seed, pubkey = generate_keypair()
        print(f"[!] New keypair generated.")
        print(f"    Private Seed: {seed.hex()}")
        print(f"    Public Key:   {pubkey.hex()}")
        print(f"    KEEP THE PRIVATE SEED SECURE!")
        
    bin_data = emit_binary(ast, seed)
    
    with open(args.out, 'wb') as f:
        f.write(bin_data)
        
    print(f"[SUCCESS] Compiled policy to {args.out} ({len(bin_data)} bytes)")

if __name__ == "__main__":
    main()
