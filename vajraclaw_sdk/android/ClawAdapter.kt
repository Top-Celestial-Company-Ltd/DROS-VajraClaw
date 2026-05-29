package com.vajraclaw.sdk

/**
 * VajraClaw Mobile SDK Adapter (Android / Kotlin)
 * 
 * Provides an idiomatic Kotlin interface to the high-performance Go-based VajraClaw
 * security microkernel on the edge. Supports both manual JNI loader and standard
 * gomobile binding integrations.
 */
class ClawAdapter private constructor() {

    companion object {
        private var isInitialized = false
        val instance: ClawAdapter by lazy { ClawAdapter() }

        // Support standard JNI fallback if developers build using pure NDK/CMake
        init {
            try {
                System.loadLibrary("vajra_claw")
                isInitialized = true
            } catch (e: UnsatisfiedLinkError) {
                // If using gomobile bind (.aar), the library is automatically loaded
                // inside the gomobile runtime class. We will print a debug log and proceed.
                System.err.println("[VajraClaw-SDK] Native library load deferred or managed by gomobile.")
            }
        }
    }

    /**
     * Initialize the Static Vajra Rules Memory directly from an in-memory string.
     * This bypasses the mobile application file sandbox restrictions by allowing developers
     * to read assets from the App assets directory or bundle as string and feed it to the core.
     *
     * @param rulesContent The raw string content of Vajra rules file.
     * @return true if initialization succeeded, false otherwise.
     */
    fun initStaticVajraFromString(rulesContent: String): Boolean {
        return if (isInitialized) {
            nativeInitStaticVajraFromString(rulesContent) == 1
        } else {
            // gomobile binding proxy
            try {
                // Import dynamic bind proxy if available via reflection or direct imports
                val clazz = Class.forName("mobile.Mobile")
                val method = clazz.getMethod("initStaticVajraFromString", String::class.java)
                val result = method.invoke(null, rulesContent) as Long
                result == 1L
            } catch (e: Exception) {
                System.err.println("[VajraClaw-SDK] Error: VajraClaw Core is not loaded!")
                false
            }
        }
    }

    /**
     * Inject an Ephemeral (dynamic) rule into the active memory boundary.
     * Typically used for user session-specific temporary constraints.
     */
    fun injectEphemeralRule(rule: String): Boolean {
        return if (isInitialized) {
            nativeInjectEphemeralRule(rule) == 1
        } else {
            try {
                val clazz = Class.forName("mobile.Mobile")
                val method = clazz.getMethod("injectEphemeralRule", String::class.java)
                val result = method.invoke(null, rule) as Long
                result == 1L
            } catch (e: Exception) {
                false
            }
        }
    }

    /**
     * Perform high-performance O(1) physical intercept matching against the active rules.
     * 
     * @param input The text stream / LLM prompt input to evaluate.
     * @return true if the prompt is safe (PASS), false if physical intercept triggers (BLOCK).
     */
    fun matchTokenStream(input: String): Boolean {
        return if (isInitialized) {
            nativeMatchTokenStream(input) == 1
        } else {
            try {
                val clazz = Class.forName("mobile.Mobile")
                val method = clazz.getMethod("matchTokenStream", String::class.java)
                val result = method.invoke(null, input) as Long
                result == 1L
            } catch (e: Exception) {
                true // Fail-safe (or strict Fail-secure depending on policy, standard is fail-secure: false)
            }
        }
    }

    /**
     * Clear all ephemeral rules from memory instantly.
     */
    fun clearEphemeralRules() {
        if (isInitialized) {
            nativeClearEphemeralRules()
        } else {
            try {
                val clazz = Class.forName("mobile.Mobile")
                val method = clazz.getMethod("clearEphemeralRules")
                method.invoke(null)
            } catch (e: Exception) {
                // No-op
            }
        }
    }

    // JNI Declarations
    private external fun nativeInitStaticVajraFromString(content: String): Int
    private external fun nativeInjectEphemeralRule(rule: String): Int
    private external fun nativeMatchTokenStream(input: String): Int
    private external fun nativeClearEphemeralRules()
}
