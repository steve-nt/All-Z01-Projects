/**
 * WebSocket client with connection lifecycle and reconnection.
 * Authenticates using the existing session (browser sends session cookie to same origin).
 * Only connect when the app is logged in; disconnect on logout or 401.
 */

const WS_STATE = {
  DISCONNECTED: "disconnected",
  CONNECTING: "connecting",
  OPEN: "open",
  CLOSED: "closed",
  RECONNECTING: "reconnecting"
};

const DEFAULT_RECONNECT_DELAY_MS = 1000;
const MAX_RECONNECT_DELAY_MS = 30000;
const RECONNECT_BACKOFF_MULTIPLIER = 1.5;

function getWsUrl() {
  const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  const host = window.location.host;
  return `${protocol}//${host}/ws`;
}

class WebSocketService {
  constructor() {
    this._ws = null;
    this._state = WS_STATE.DISCONNECTED;
    this._reconnectDelay = DEFAULT_RECONNECT_DELAY_MS;
    this._reconnectTimer = null;
    this._shouldConnect = false; // only true when app is logged in and wants WS
    this._listeners = { open: [], close: [], message: [], error: [] };
  }

  getState() {
    return this._state;
  }

  isOpen() {
    return this._state === WS_STATE.OPEN && this._ws?.readyState === WebSocket.OPEN;
  }

  /** Call when the app wants an active connection (e.g. after login). */
  connect() {
    if (this._shouldConnect && (this._state === WS_STATE.OPEN || this._state === WS_STATE.CONNECTING)) {
      return;
    }
    this._shouldConnect = true;
    this._doConnect();
  }

  /** Call when the app no longer wants a connection (e.g. logout). */
  disconnect() {
    this._shouldConnect = false;
    this._clearReconnectTimer();
    if (this._ws) {
      this._ws.close(1000, "Client disconnect");
      this._ws = null;
    }
    this._setState(WS_STATE.DISCONNECTED);
  }

  _clearReconnectTimer() {
    if (this._reconnectTimer != null) {
      clearTimeout(this._reconnectTimer);
      this._reconnectTimer = null;
    }
  }

  _setState(state) {
    if (this._state === state) return;
    this._state = state;
    this._emit("close", { state: this._state });
  }

  _doConnect() {
    if (!this._shouldConnect) return;
    if (this._ws != null) {
      this._ws.close();
      this._ws = null;
    }
    const url = getWsUrl();
    this._setState(WS_STATE.CONNECTING);
    try {
      this._ws = new WebSocket(url);
    } catch (err) {
      this._emit("error", err);
      this._scheduleReconnect();
      return;
    }

    this._ws.onopen = () => {
      this._reconnectDelay = DEFAULT_RECONNECT_DELAY_MS;
      this._state = WS_STATE.OPEN;
      this._emit("open", {});
    };

    this._ws.onclose = (event) => {
      const wasOpen = this._state === WS_STATE.OPEN;
      this._ws = null;
      if (!this._shouldConnect) {
        this._setState(WS_STATE.DISCONNECTED);
        return;
      }
      // 4401 or 1008: auth failed or policy; don't reconnect
      if (event.code === 4401 || event.code === 1008) {
        this._shouldConnect = false;
        this._setState(WS_STATE.DISCONNECTED);
        this._emit("error", new Error("WebSocket unauthorized"));
        return;
      }
      this._setState(WS_STATE.CLOSED);
      // Only reconnect if we had a successful connection before (avoids reconnect loop on initial auth failure)
      if (wasOpen) {
        this._scheduleReconnect();
      }
    };

    this._ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        this._emit("message", data);
      } catch {
        this._emit("message", { raw: event.data });
      }
    };

    this._ws.onerror = (event) => {
      this._emit("error", event);
    };
  }

  _scheduleReconnect() {
    if (!this._shouldConnect || this._reconnectTimer != null) return;
    this._state = WS_STATE.RECONNECTING;
    this._reconnectTimer = setTimeout(() => {
      this._reconnectTimer = null;
      this._reconnectDelay = Math.min(
        this._reconnectDelay * RECONNECT_BACKOFF_MULTIPLIER,
        MAX_RECONNECT_DELAY_MS
      );
      this._doConnect();
    }, this._reconnectDelay);
  }

  _emit(event, payload) {
    const list = this._listeners[event] || [];
    list.forEach((fn) => {
      try {
        fn(payload);
      } catch (e) {
        console.warn("[WebSocket] listener error:", e);
      }
    });
  }

  on(event, callback) {
    if (!this._listeners[event]) this._listeners[event] = [];
    this._listeners[event].push(callback);
  }

  off(event, callback) {
    const list = this._listeners[event];
    if (!list) return;
    const i = list.indexOf(callback);
    if (i !== -1) list.splice(i, 1);
  }

  send(data) {
    if (!this.isOpen() || !this._ws) {
      console.warn("[WebSocket] send called while not open");
      return;
    }
    const payload = typeof data === "string" ? data : JSON.stringify(data);
    this._ws.send(payload);
  }
}

export const wsService = new WebSocketService();
export { WS_STATE, getWsUrl };
