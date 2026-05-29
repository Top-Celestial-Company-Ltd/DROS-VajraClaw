const ffi = require('ffi-napi');
const ref = require('ref-napi');
const path = require('path');

/**
 * VajraClawAdapter - The Official Node.js SDK for VajraClaw
 * Provides O(1) physical interception of LLM streams at the C-FFI layer.
 */
class VajraClawAdapter {
    /**
     * Initializes the C-FFI binding to the VajraClaw binary.
     * @param {string} dllPath - Path to vajra_claw.dll or vajra_claw.so
     * @param {string} matrixPath - Path to Vajra.md (The physical constraints matrix)
     */
    constructor(dllPath, matrixPath) {
        this.dllPath = path.resolve(dllPath);
        this.matrixPath = path.resolve(matrixPath);
        
        // Define C Types
        this.stringPtr = ref.refType(ref.types.CString);
        
        try {
            // Bind to the native C functions exported by VajraClaw
            this.lib = ffi.Library(this.dllPath, {
                'load_vajra_matrix': ['int', ['string']],
                'evaluate_token': ['int', ['string']],
                'inject_ephemeral_bounds': ['int', ['string']],
                'get_last_error': ['string', []]
            });
            
            // Load the static matrix immediately (O(1) initialization)
            const result = this.lib.load_vajra_matrix(this.matrixPath);
            if (result !== 1) {
                const err = this.lib.get_last_error();
                throw new Error(`[VajraClaw] FATAL: Failed to load Static Matrix. Reason: ${err}`);
            }
            
            console.log("🛡️ [VajraClaw-Core] Node.js FFI Adapter Initialized. Matrix Locked.");
        } catch (error) {
            console.error("🛑 [VajraClaw-Exception] Failed to bind to physical DLL:", error.message);
            throw error;
        }
    }

    /**
     * Evaluates a single token from the LLM stream.
     * @param {string} token - The token string output by the LLM
     * @returns {boolean} - True if token is safe, False if fuse blown
     */
    streamMonitor(token) {
        if (!token) return true;
        
        // Physical O(1) evaluation via C-FFI
        const isSafe = this.lib.evaluate_token(token);
        
        if (isSafe !== 1) {
            const err = this.lib.get_last_error();
            console.error(`🛑 [PHYSICAL FUSE BLOWN] Token blocked: "${token}". Reason: ${err}`);
            return false;
        }
        
        return true;
    }
}

module.exports = { VajraClawAdapter };
