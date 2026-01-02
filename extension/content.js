// content.js - Scrapes text and sends it to the background script
console.log("SentinelAI Content Script Loaded.");

// Example: Analyze text whenever the user stops typing in a textarea
let timeout = null;
document.addEventListener('keyup', (event) => {
    if (event.target.tagName === 'TEXTAREA' || event.target.tagName === 'INPUT') {
        clearTimeout(timeout);
        timeout = setTimeout(() => {
            const textToAnalyze = event.target.value;
            
            if (textToAnalyze.length > 5) {
                chrome.runtime.sendMessage({ text: textToAnalyze }, (response) => {
                    console.log("SentinelAI Analysis Result:", response);
                    
                    // Optional: Visual feedback if a policy is violated
                    if (response && response.violation) {
                        event.target.style.border = "2px solid red";
                    } else {
                        event.target.style.border = "";
                    }
                });
            }
        }, 1000); // Wait for 1 second of inactivity
    }
});