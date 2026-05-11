package server

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/gorilla/websocket"
)

// maxLobbyPlayers is the maximum clients allowed in the waiting lobby before a match.
// Spectators may connect while a game is running; when the match ends, excess clients are disconnected.
const maxLobbyPlayers = 4

const maxLobbyChatMessages = 50
const maxLobbyChatRunes = 160

type lobbyChatMsg struct {
	From string `json:"from"`
	Text string `json:"text"`
}

type client struct {
	conn      *websocket.Conn
	nickname  string
	joinedAt  time.Time
	slot      int
	inGame    bool
	spectator bool
}

type lobby struct {
	mu             sync.Mutex
	clients        map[*client]struct{}
	startDeadline  time.Time
	phase          string
	countdownStart time.Time
	game           *gameSession
	chatLog        []lobbyChatMsg
}

func newLobby() *lobby {
	return &lobby{
		clients: make(map[*client]struct{}),
		phase:   "waiting",
	}
}

func (l *lobby) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	l.mu.Lock()
	lobbyFull := l.game == nil && len(l.clients) >= maxLobbyPlayers
	l.mu.Unlock()
	if lobbyFull {
		http.Error(w, "lobby full: maximum 4 players", http.StatusServiceUnavailable)
		return
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	c := &client{conn: conn, nickname: "Player", joinedAt: time.Now(), slot: -1}
	l.mu.Lock()
	l.clients[c] = struct{}{}
	gs := l.game
	l.mu.Unlock()

	if gs != nil {
		c.spectator = true
		gs.addSpectator(c)
		gs.sendSpectateStart(c)
	} else {
		l.broadcastLobby()
	}

	defer func() {
		l.mu.Lock()
		gs := l.game
		l.mu.Unlock()
		if gs != nil {
			if c.spectator {
				gs.removeSpectator(c)
			} else if c.inGame {
				gs.disconnectPlayer(c)
			}
		}
		l.mu.Lock()
		delete(l.clients, c)
		l.mu.Unlock()
		_ = conn.Close()
		l.broadcastLobby()
	}()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			return
		}
		var msg map[string]any
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}
		typ, _ := msg["type"].(string)

		if typ == "input" {
			l.mu.Lock()
			gs := l.game
			l.mu.Unlock()
			if gs != nil && c.inGame {
				gs.handleInput(c, msg)
			}
			continue
		}

		if typ == "slime_hit" {
			l.mu.Lock()
			gs := l.game
			l.mu.Unlock()
			if gs != nil && c.inGame {
				gs.handleSlimeHit(c)
			}
			continue
		}

		switch typ {
		case "join":
			n, _ := msg["nickname"].(string)
			if n != "" {
				c.nickname = n
			}
			l.broadcastLobby()
		case "chat":
			text, _ := msg["text"].(string)
			trimmed, ok := l.recordChat(c.nickname, text)
			if !ok {
				continue
			}
			l.broadcast(map[string]any{
				"type": "chat",
				"from": c.nickname,
				"text": trimmed,
			})
		}
	}
}

// recordChat stores a trimmed message and returns it for broadcast. Caller must not hold l.mu.
func (l *lobby) recordChat(from, raw string) (trimmed string, ok bool) {
	trimmed = strings.TrimSpace(raw)
	if trimmed == "" {
		return "", false
	}
	if utf8.RuneCountInString(trimmed) > maxLobbyChatRunes {
		trimmed = string([]rune(trimmed)[:maxLobbyChatRunes])
	}
	l.mu.Lock()
	l.chatLog = append(l.chatLog, lobbyChatMsg{From: from, Text: trimmed})
	if len(l.chatLog) > maxLobbyChatMessages {
		l.chatLog = l.chatLog[len(l.chatLog)-maxLobbyChatMessages:]
	}
	l.mu.Unlock()
	return trimmed, true
}

func (l *lobby) snapshotSortedClients() []*client {
	out := make([]*client, 0, len(l.clients))
	for c := range l.clients {
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].joinedAt.Before(out[j].joinedAt)
	})
	return out
}

func (l *lobby) beginGame() {
	l.mu.Lock()
	if l.game != nil {
		l.mu.Unlock()
		return
	}
	participants := l.snapshotSortedClients()
	if len(participants) < 2 {
		l.phase = "waiting"
		l.countdownStart = time.Time{}
		l.mu.Unlock()
		return
	}
	if len(participants) > maxLobbyPlayers {
		participants = participants[:maxLobbyPlayers]
	}
	l.mu.Unlock()

	g := newGameSession(l, participants)

	l.mu.Lock()
	l.game = g
	extras := make([]*client, 0)
	for c := range l.clients {
		inPart := false
		for _, p := range participants {
			if p == c {
				inPart = true
				break
			}
		}
		if !inPart {
			extras = append(extras, c)
		}
	}
	l.mu.Unlock()

	for _, c := range extras {
		c.spectator = true
		g.addSpectator(c)
		g.sendSpectateStart(c)
	}

	g.sendGameStart()
	go g.run()
}

func (l *lobby) monitorStart() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		var shouldBegin bool
		l.mu.Lock()
		if l.game != nil {
			l.mu.Unlock()
			continue
		}

		now := time.Now()
		count := len(l.clients)

		// Fewer than 2: nothing to schedule.
		if count < 2 {
			l.startDeadline = time.Time{}
			l.countdownStart = time.Time{}
			l.phase = "waiting"
		} else if count < maxLobbyPlayers {
			// 2–3 players: 20s grace for more to join before a 10s start countdown can begin.
			if l.startDeadline.IsZero() {
				l.startDeadline = now.Add(20 * time.Second)
			}
		}

		// Full lobby (4): 10s countdown to start immediately (no 20s grace).
		if count >= maxLobbyPlayers {
			if l.countdownStart.IsZero() {
				l.countdownStart = now.Add(10 * time.Second)
			}
		} else if count >= 2 && count < maxLobbyPlayers && !l.startDeadline.IsZero() && now.After(l.startDeadline) {
			// Grace ended with 2–3 players still waiting: 10s countdown, then match starts.
			if l.countdownStart.IsZero() {
				l.countdownStart = now.Add(10 * time.Second)
			}
		}

		if !l.countdownStart.IsZero() && now.After(l.countdownStart) && l.phase != "playing" {
			l.phase = "playing"
			shouldBegin = true
		}
		l.mu.Unlock()

		if shouldBegin {
			l.beginGame()
		}
		l.broadcastLobby()
	}
}

func (l *lobby) broadcastLobby() {
	l.mu.Lock()
	players := make([]string, 0, len(l.clients))
	for c := range l.clients {
		players = append(players, c.nickname)
	}
	startAt := l.countdownStart.UnixMilli()
	if l.countdownStart.IsZero() {
		startAt = 0
	}
	var graceEndsAt int64
	n := len(players)
	if l.countdownStart.IsZero() && n >= 2 && n < maxLobbyPlayers && !l.startDeadline.IsZero() {
		graceEndsAt = l.startDeadline.UnixMilli()
	}
	msgs := make([]map[string]any, 0, len(l.chatLog))
	for _, m := range l.chatLog {
		msgs = append(msgs, map[string]any{"from": m.From, "text": m.Text})
	}
	payload := map[string]any{
		"type":           "lobby",
		"count":          len(players),
		"players":        players,
		"phase":          l.phase,
		"startAtMs":      startAt,
		"graceEndsAtMs":  graceEndsAt,
		"messages":       msgs,
	}
	l.mu.Unlock()
	l.broadcast(payload)
}

func (l *lobby) broadcast(payload map[string]any) {
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}
	l.mu.Lock()
	clients := make([]*client, 0, len(l.clients))
	for c := range l.clients {
		clients = append(clients, c)
	}
	l.mu.Unlock()

	for _, c := range clients {
		if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("websocket write failed: %v", err)
		}
	}
}

// disconnectExcessWaitingClients closes connections over the lobby cap (newest first).
// Call after a match ends when spectators are folded back into the waiting pool.
func (l *lobby) disconnectExcessWaitingClients() {
	l.mu.Lock()
	if l.game != nil || len(l.clients) <= maxLobbyPlayers {
		l.mu.Unlock()
		return
	}
	sorted := make([]*client, 0, len(l.clients))
	for c := range l.clients {
		sorted = append(sorted, c)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].joinedAt.Before(sorted[j].joinedAt)
	})
	toKick := sorted[maxLobbyPlayers:]
	l.mu.Unlock()

	for _, c := range toKick {
		_ = c.conn.Close()
	}
}
