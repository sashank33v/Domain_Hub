const chatToggle = document.getElementById("chat-toggle");
const chatPanel = document.getElementById("chat-panel");
const chatClose = document.getElementById("chat-close");
const chatMessages = document.getElementById("chat-messages");
const chatInput = document.getElementById("chat-input");

chatToggle.addEventListener("click", () => {
    chatPanel.classList.toggle("hidden");
    if (!chatPanel.classList.contains("hidden")) {
        chatInput.focus();
    }
});

chatClose.addEventListener("click", () => {
    chatPanel.classList.add("hidden");
});
function escapeHtml(text) {
    return text
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#39;");
}

function formatReply(text) {
    return escapeHtml(text).replace(/\*\*(.+?)\*\*/g, "<strong>$1</strong>");
}
function addMessage(text, role) {
    const div = document.createElement("div");

    div.className =
        role === "user"
            ? "bg-[#3525cd] text-white rounded-lg px-3 py-2 w-fit max-w-[85%] ml-auto whitespace-pre-wrap"
            : "bg-slate-100 rounded-lg px-3 py-2 w-fit max-w-[85%] whitespace-pre-wrap";

    if (role === "user") {
        div.textContent = text;
    } else {
        div.innerHTML = formatReply(text);
    }
    chatMessages.appendChild(div);
    chatMessages.scrollTop = chatMessages.scrollHeight;

    return div;
}
const chatForm = document.getElementById("chat-form");
const chatSend = document.getElementById("chat-send");

const chatHistory = [];
const MAX_HISTORY = 6;

chatForm.addEventListener("submit", async (e) => {
    e.preventDefault();

    const message = chatInput.value.trim();
    if (!message) return;

    addMessage(message, "user");
    chatInput.value = "";
    chatInput.disabled = true;
    chatSend.disabled = true;

    const loading = addMessage("Typing...", "assistant");

    try {
        const reply = await sendChat(message, chatHistory.slice(-MAX_HISTORY));

        loading.innerHTML = formatReply(reply);

        chatHistory.push(
            { role: "user", content: message },
            { role: "assistant", content: reply }
        );
    } catch (err) {
        loading.textContent =
            err.message || "Something went wrong. Please try again.";
    } finally {
        chatInput.disabled = false;
        chatSend.disabled = false;
        chatInput.focus();
        chatMessages.scrollTop = chatMessages.scrollHeight;
    }
});