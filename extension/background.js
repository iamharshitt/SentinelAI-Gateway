let port = chrome.runtime.connectNative("com.sentinelai.gateway");

// Handle disconnects (e.g., if the Go agent crashes or is closed)
port.onDisconnect.addListener(() => {
    console.error("Disconnected from Go agent:", chrome.runtime.lastError?.message);
    // Optional: Reconnect logic here
});

chrome.runtime.onMessage.addListener((msg, sender, sendResponse) => {
    // 1. Send the prompt to the Go agent
    port.postMessage({ prompt: msg.text });

    // 2. Define a ONE-TIME listener for the response
    const responseHandler = (response) => {
        port.onMessage.removeListener(responseHandler); // Clean up to prevent memory leaks
        sendResponse(response);
    };

    port.onMessage.addListener(responseHandler);

    return true; // Keep the message channel open for sendResponse
});