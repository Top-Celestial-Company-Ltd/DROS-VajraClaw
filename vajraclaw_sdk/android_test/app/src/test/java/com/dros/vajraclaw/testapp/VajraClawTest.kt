package com.dros.vajraclaw.testapp

import org.junit.Test
import org.junit.Assert.*
import mobile.Mobile

class VajraClawTest {

    @Test
    fun testVajraClawSDK() {
        try {
            // Test 1: Initialize
            val initRes = Mobile.initStaticVajraFromString("test-rule-1\ntest-rule-2")
            assertEquals("Initialization should return 1", 1L, initRes)

            // Test 2: Secure Prompt
            val secureRes = Mobile.matchTokenStream("This is a safe prompt without any rules.")
            assertEquals("Secure prompt should return 1", 1L, secureRes)

            // Test 3: Blocked Prompt
            val blockedRes = Mobile.matchTokenStream("This prompt contains test-rule-1 inside.")
            assertEquals("Blocked prompt should return 0", 0L, blockedRes)
            
            println("All SDK static checks passed.")
        } catch (e: UnsatisfiedLinkError) {
            println("WARNING: Native library loading failed. This is expected if running on Windows without Android emulator, as the .so libraries in AAR are for Linux/Android ELF format.")
            e.printStackTrace()
            fail("Native library linking failed: ${e.message}")
        }
    }
}
