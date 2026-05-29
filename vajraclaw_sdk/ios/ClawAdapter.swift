import Foundation

/**
 * VajraClaw Mobile SDK Adapter (iOS / Swift)
 * 
 * Provides an elegant, idiomatic Swift wrapper binding directly to the native C-FFI
 * compiled VajraClaw Go microkernel. Utilizes Swift's seamless native C-Interop.
 */
public final class ClawAdapter {
    
    public static let shared = ClawAdapter()
    
    private init() {}
    
    /**
     * Initialize the Static Vajra Rules Memory directly from a Swift string.
     * Bypasses the iOS bundle path sandboxing restrictions by reading rules in-memory.
     *
     * @param rulesContent Raw string content of the Vajra rules.
     * @return true if initialization succeeded.
     */
    @discardableResult
    public func initStaticVajra(from rulesContent: String) -> Bool {
        return rulesContent.withCString { cString in
            // Call C function exported by Go compiler
            return init_static_vajra_from_string(UnsafeMutablePointer(mutating: cString)) == 1
        }
    }
    
    /**
     * Inject an Ephemeral (dynamic) rule JIT.
     */
    @discardableResult
    public func injectEphemeralRule(_ rule: String) -> Bool {
        return rule.withCString { cString in
            return inject_ephemeral_rule(UnsafeMutablePointer(mutating: cString)) == 1
        }
    }
    
    /**
     * Intercept and match prompt input stream.
     * 
     * @return true if prompt is safe (PASS), false if physical intercept triggers (BLOCK).
     */
    public func matchTokenStream(_ input: String) -> Bool {
        return input.withCString { cString in
            return match_token_stream(UnsafeMutablePointer(mutating: cString)) == 1
        }
    }
    
    /**
     * Physical evaporation of ephemeral session boundaries.
     */
    public func clearEphemeralRules() {
        clear_ephemeral_rules()
    }
}
