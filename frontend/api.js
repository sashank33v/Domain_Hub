function getToken() {
    return localStorage.getItem("jwt_token");
}

function getUserEmail() {
    return localStorage.getItem("user_email") || "";
}

function saveAuth(token, email) {
    localStorage.setItem("jwt_token", token);
    localStorage.setItem("user_email", email);
}

function clearAuth() {
    localStorage.removeItem("jwt_token");
    localStorage.removeItem("user_email");
}

function redirectToLogin() {
    clearAuth();
    window.location.href = "login.html?error=401";
}

function requireAuth() {
    const token = getToken();

    if (!token) {
        redirectToLogin();
        return false;
    }

    return true;
}

async function apiRequest(path, options = {}) {

    const headers = {
        ...(options.headers || {})
    };

    if (options.body && !headers["Content-Type"]) {
        headers["Content-Type"] = "application/json";
    }

    const token = getToken();

    if (token) {
        headers["Authorization"] = `Bearer ${token}`;
    }

    const response = await fetch(
        `${API_BASE_URL}${path}`,
        {
            ...options,
            headers
        }
    );

    let data = null;

    const contentType =
        response.headers.get("content-type") || "";

    if (contentType.includes("application/json")) {
        data = await response.json();
    } else {
        const text = await response.text();

        data = text
            ? { message: text }
            : null;
    }

    /*
     * Only redirect to login when an authenticated API
     * returns 401.
     *
     * /auth/login itself can legitimately return 401.
     */
    if (
        response.status === 401 &&
        !path.startsWith("/auth/")
    ) {
        redirectToLogin();
        throw new Error("Session expired. Please login again.");
    }

    if (!response.ok) {
        throw new Error(
            data?.message || "Request failed"
        );
    }

    return data;
}

// -------------------------
// Authentication
// -------------------------

async function loginUser(email, password) {

    const result = await apiRequest(
        "/auth/login",
        {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                email,
                password
            })
        }
    );

    if (!result?.data?.token) {
        throw new Error(
            "Login response did not contain a token"
        );
    }

    saveAuth(
        result.data.token,
        email
    );

    return result;
}


async function registerUser(email, password) {

    return apiRequest(
        "/auth/register",
        {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                email,
                password
            })
        }
    );
}


// -------------------------
// Domains
// -------------------------

async function getDomains({
    page = 1,
    limit = 10,
    status = "",
    registrar = ""
} = {}) {

    const params = new URLSearchParams();

    params.set("page", String(page));
    params.set("limit", String(limit));

    if (status) {
        params.set("status", status);
    }

    if (registrar) {
        params.set("registrar", registrar);
    }

    return apiRequest(
        `/domains?${params.toString()}`,
        {
            method: "GET"
        }
    );
}


async function getDomain(id) {

    return apiRequest(
        `/domains/${id}`,
        {
            method: "GET"
        }
    );
}


async function createDomain(domain) {

    return apiRequest(
        "/domains",
        {
            method: "POST",
            body: JSON.stringify(domain)
        }
    );
}


async function updateDomain(id, domain) {

    return apiRequest(
        `/domains/${id}`,
        {
            method: "PUT",
            body: JSON.stringify(domain)
        }
    );
}


async function deleteDomain(id) {

    return apiRequest(
        `/domains/${id}`,
        {
            method: "DELETE"
        }
    );
}


// -------------------------
// Helpers
// -------------------------

function formatDateForInput(value) {

    if (!value) {
        return "";
    }

    return String(value).slice(0, 10);
}


function formatDisplayDate(value) {

    if (!value) {
        return "-";
    }

    const date = new Date(value);

    if (Number.isNaN(date.getTime())) {
        return value;
    }

    return date.toLocaleDateString(
        "en-IN",
        {
            year: "numeric",
            month: "short",
            day: "numeric"
        }
    );
}


function dateToRFC3339(value) {

    if (!value) {
        return "";
    }

    return `${value}T00:00:00Z`;
}


function logout() {

    clearAuth();

    window.location.href = "login.html";
}