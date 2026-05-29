package com.dros.vajraclaw.testapp

import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.widget.TextView
import mobile.Mobile

class MainActivity : AppCompatActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        
        val textView = TextView(this)
        textView.textSize = 16f
        textView.setPadding(32, 32, 32, 32)
        setContentView(textView)

        val log = StringBuilder()
        log.append("Starting VajraClaw Mobile SDK Test...\n\n")

        try {
            // 1. Initialize SDK
            val initResult = Mobile.initStaticVajraFromString("test-rule-1\ntest-rule-2")
            log.append("InitResult: $initResult\n")

            // 2. Test Secure Prompt
            val secureResult = Mobile.matchTokenStream("Hello world")
            log.append("MatchTokenStream('Hello world') (Expected 1): $secureResult\n")

            // 3. Test Blocked Prompt
            val blockedResult = Mobile.matchTokenStream("This contains test-rule-1 inside")
            log.append("MatchTokenStream('This contains test-rule-1 inside') (Expected 0): $blockedResult\n")

            if (initResult == 1L && secureResult == 1L && blockedResult == 0L) {
                log.append("\n✅ INTEGRATION TEST PASSED!")
            } else {
                log.append("\n❌ INTEGRATION TEST FAILED!")
            }

        } catch (e: Exception) {
            log.append("\n❌ Exception occurred: ${e.message}")
            e.printStackTrace()
        }

        textView.text = log.toString()
    }
}
