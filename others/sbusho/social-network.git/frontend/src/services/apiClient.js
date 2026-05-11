// Part 1: shared API client conventions
let unauthorizedHandler = null;
let forbiddenHandler = null;
let networkErrorHandler = null;

// Part 1: 401 handling (clear session + redirect)
export function setUnauthorizedHandler(handler) {
  unauthorizedHandler = handler;
}

// Part 1: 403 handling (show "forbidden" UI). Handler receives (path, payload) for context (e.g. group "members only").
export function setForbiddenHandler(handler) {
  forbiddenHandler = handler;
}

// Part 1: network failure handling (show banner/toast)
export function setNetworkErrorHandler(handler) {
  networkErrorHandler = handler;
}

// Part 1: normalize API errors for callers
export class ApiError extends Error {
  constructor(message, status, data) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.data = data;
  }
}

function isJsonResponse(response) {
  const contentType = response.headers.get("content-type") || "";
  return contentType.includes("application/json");
}

// Part 1: wrapper around fetch with shared defaults + errors
export async function apiRequest(path, options = {}) {
  const config = {
    credentials: "include",
    headers: {
      ...(options.headers || {})
    },
    ...options
  };

  if (config.body && !(config.body instanceof FormData)) {
    config.headers["Content-Type"] =
      config.headers["Content-Type"] || "application/json";
  }

  let response;
  try {
    response = await fetch(path, config);
  } catch (error) {
    if (typeof networkErrorHandler === "function") {
      networkErrorHandler(error);
    }
    throw new ApiError("Network error", 0, null);
  }

  if (response.status === 401 && typeof unauthorizedHandler === "function") {
    unauthorizedHandler();
  }

  const shouldParseJson = isJsonResponse(response);
  let payload = null;

  if (!response.ok) {
    payload = shouldParseJson ? await response.json().catch(() => null) : null;
    if (response.status === 403 && typeof forbiddenHandler === "function") {
      forbiddenHandler(path, payload);
    }
    const message =
      (payload && (payload.message || payload.error)) ||
      response.statusText ||
      "Request failed";
    throw new ApiError(message, response.status, payload);
  }

  payload = shouldParseJson ? await response.json() : null;
  return payload;
}
