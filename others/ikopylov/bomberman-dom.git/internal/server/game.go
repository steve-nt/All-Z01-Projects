package server

import (
	"encoding/json"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	gridRows = 20
	gridCols = 20
	tileSize = 32
)

type inputState struct {
	Up, Down, Left, Right, Bomb bool
}

func (in inputState) dir() (int, int) {
	h, v := 0, 0
	if in.Left && !in.Right {
		h = -1
	}
	if in.Right && !in.Left {
		h = 1
	}
	if in.Up && !in.Down {
		v = -1
	}
	if in.Down && !in.Up {
		v = 1
	}
	if h != 0 && v != 0 {
		return h, 0
	}
	if h != 0 {
		return h, 0
	}
	if v != 0 {
		return 0, v
	}
	return 0, 0
}

type gPlayer struct {
	Slot       int
	Nickname   string
	Col, Row   int
	Lives      int
	Dead       bool
	MaxBombs   int
	Radius     int
	SpeedLevel int
	BombCount  int
	InvUntil   time.Time
	MoveAccum  float64
	PrevBomb   bool
	client     *client
	in         inputState
}

type gBomb struct {
	ID        int
	Col, Row  int
	Timer     float64
	Radius    int
	OwnerSlot int
}

type gPowerUp struct {
	ID       int
	Col, Row int
	Type     string
}

// tileDelta is a single map cell change (bandwidth-friendly vs full grid).
type tileDelta struct {
	R int `json:"r"`
	C int `json:"c"`
	V int `json:"v"`
}

type gameSession struct {
	lobby *lobby

	mu       sync.Mutex
	tiles    [][]int
	players  []*gPlayer
	bombs    []*gBomb
	powerUps []*gPowerUp
	nextIDs  struct{ bomb, power int }

	tileDeltas []tileDelta
	spectators []*client

	// Goal portal (same cell chosen for all clients in game_start). Revealed when that soft block is destroyed.
	exitRow, exitCol int
	exitRevealed     bool

	running bool
	rng     *rand.Rand
	frame   int
}

func randomizeSoftBlocks(tiles [][]int, rng *rand.Rand) {
	for r := 0; r < gridRows; r++ {
		for c := 0; c < gridCols; c++ {
			if tiles[r][c] != 2 {
				continue
			}
			if inSpawnSafeZone(r, c) {
				tiles[r][c] = 0
				continue
			}
			if rng.Float64() < 0.38 {
				tiles[r][c] = 0
			}
		}
	}
	for i := 0; i < len(spawnCorners); i++ {
		sr, sc := spawnCorners[i][0], spawnCorners[i][1]
		for dr := -2; dr <= 2; dr++ {
			for dc := -2; dc <= 2; dc++ {
				rr, cc := sr+dr, sc+dc
				if rr >= 0 && rr < gridRows && cc >= 0 && cc < gridCols && tiles[rr][cc] == 2 {
					tiles[rr][cc] = 0
				}
			}
		}
	}
}

func newGameSession(l *lobby, participants []*client) *gameSession {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tiles := copyTiles(baseTileTemplate())
	randomizeSoftBlocks(tiles, rng)

	g := &gameSession{
		lobby:   l,
		tiles:   tiles,
		rng:     rng,
		running: true,
	}
	n := len(participants)
	if n > maxLobbyPlayers {
		n = maxLobbyPlayers
	}
	for i := 0; i < n; i++ {
		c := participants[i]
		sr, sc := spawnCorners[i][0], spawnCorners[i][1]
		g.players = append(g.players, &gPlayer{
			Slot:     i,
			Nickname: c.nickname,
			Col:      sc,
			Row:      sr,
			Lives:    3,
			MaxBombs: 1,
			Radius:   1,
			client:   c,
		})
		c.slot = i
		c.inGame = true
	}
	g.pickExitCellUnderSoft()
	return g
}

// minManhattanToNearestSpawn is used to place the goal far from all corner starts.
func minManhattanToNearestSpawn(row, col int) int {
	minD := gridRows + gridCols
	for i := 0; i < len(spawnCorners); i++ {
		sr, sc := spawnCorners[i][0], spawnCorners[i][1]
		d := abs(row-sr) + abs(col-sc)
		if d < minD {
			minD = d
		}
	}
	return minD
}

func (g *gameSession) pickExitCellUnderSoft() {
	g.exitRow, g.exitCol = -1, -1
	g.exitRevealed = false
	var soft [][2]int
	for r := 0; r < gridRows; r++ {
		for c := 0; c < gridCols; c++ {
			if g.tiles[r][c] == 2 {
				soft = append(soft, [2]int{r, c})
			}
		}
	}
	if len(soft) == 0 {
		return
	}
	const minPreferredExitDist = 10
	filtered := soft
	var preferred [][2]int
	for _, p := range soft {
		if minManhattanToNearestSpawn(p[0], p[1]) >= minPreferredExitDist {
			preferred = append(preferred, p)
		}
	}
	if len(preferred) > 0 {
		filtered = preferred
	}
	best := -1
	var candidates [][2]int
	for _, p := range filtered {
		md := minManhattanToNearestSpawn(p[0], p[1])
		if md > best {
			best = md
			candidates = [][2]int{p}
		} else if md == best {
			candidates = append(candidates, p)
		}
	}
	picked := candidates[g.rng.Intn(len(candidates))]
	g.exitRow, g.exitCol = picked[0], picked[1]
}

func (g *gameSession) moveDelay(p *gPlayer) float64 {
	d := 0.16 - float64(p.SpeedLevel)*0.022
	if d < 0.075 {
		d = 0.075
	}
	return d
}

func (g *gameSession) isBlockedForPlayer(row, col int, exclude *gPlayer) bool {
	if row < 0 || col < 0 || row >= gridRows || col >= gridCols {
		return true
	}
	t := g.tiles[row][col]
	if t == 1 || t == 2 {
		return true
	}
	for _, b := range g.bombs {
		if b.Row == row && b.Col == col {
			return true
		}
	}
	for _, pl := range g.players {
		if pl == exclude || pl.Dead {
			continue
		}
		if pl.Row == row && pl.Col == col {
			return true
		}
	}
	return false
}

func (g *gameSession) tryMove(p *gPlayer) {
	dx, dy := p.in.dir()
	if dx == 0 && dy == 0 {
		return
	}
	nr, nc := p.Row+dy, p.Col+dx
	if g.isBlockedForPlayer(nr, nc, p) {
		return
	}
	p.Row, p.Col = nr, nc

	for i := range g.powerUps {
		pu := g.powerUps[i]
		if pu.Row == p.Row && pu.Col == p.Col {
			g.applyPowerUp(p, pu.Type)
			g.powerUps = append(g.powerUps[:i], g.powerUps[i+1:]...)
			break
		}
	}
}

func (g *gameSession) applyPowerUp(p *gPlayer, typ string) {
	switch typ {
	case "bomb":
		p.MaxBombs++
	case "radius":
		p.Radius++
	case "speed":
		p.SpeedLevel++
	}
}

func (g *gameSession) tryPlaceBomb(p *gPlayer) {
	if p.Dead {
		return
	}
	if p.BombCount >= p.MaxBombs {
		return
	}
	cr, cc := p.Row, p.Col
	for _, b := range g.bombs {
		if b.Row == cr && b.Col == cc {
			return
		}
	}
	g.nextIDs.bomb++
	b := &gBomb{
		ID:        g.nextIDs.bomb,
		Col:       cc,
		Row:       cr,
		Timer:     3.0,
		Radius:    p.Radius,
		OwnerSlot: p.Slot,
	}
	g.bombs = append(g.bombs, b)
	p.BombCount++
}

func (g *gameSession) hitPlayerDamage(p *gPlayer) {
	if p.Dead {
		return
	}
	if time.Now().Before(p.InvUntil) {
		return
	}
	p.Lives--
	if p.Lives <= 0 {
		p.Dead = true
		return
	}
	p.InvUntil = time.Now().Add(2 * time.Second)
}

func (g *gameSession) hitPlayersAt(row, col int) {
	for _, p := range g.players {
		if p.Dead {
			continue
		}
		if p.Row != row || p.Col != col {
			continue
		}
		g.hitPlayerDamage(p)
	}
}

func (g *gameSession) removeBombAt(row, col int) *gBomb {
	for j := 0; j < len(g.bombs); j++ {
		b := g.bombs[j]
		if b.Row == row && b.Col == col {
			g.bombs = append(g.bombs[:j], g.bombs[j+1:]...)
			return b
		}
	}
	return nil
}

func (g *gameSession) detonateBomb(b *gBomb) {
	for _, pl := range g.players {
		if pl.Slot == b.OwnerSlot {
			if pl.BombCount > 0 {
				pl.BombCount--
			}
			break
		}
	}

	g.hitPlayersAt(b.Row, b.Col)

	dirs := [][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
	for _, d := range dirs {
		dr, dc := d[0], d[1]
		for step := 1; step <= b.Radius; step++ {
			nr, nc := b.Row+dr*step, b.Col+dc*step
			if nr < 0 || nr >= gridRows || nc < 0 || nc >= gridCols {
				break
			}
			g.hitPlayersAt(nr, nc)

			tid := g.tiles[nr][nc]
			if tid == 1 {
				break
			}
			if tid == 2 {
				g.destroySoft(nr, nc)
				break
			}
			if ob := g.removeBombAt(nr, nc); ob != nil {
				g.detonateBomb(ob)
			}
		}
	}
}

func (g *gameSession) destroySoft(row, col int) {
	if g.tiles[row][col] != 2 {
		return
	}
	g.tiles[row][col] = 0
	if row == g.exitRow && col == g.exitCol {
		g.exitRevealed = true
	}
	g.tileDeltas = append(g.tileDeltas, tileDelta{R: row, C: col, V: 0})
	if g.rng.Float64() < 0.45 {
		types := []string{"bomb", "radius", "speed"}
		t := types[g.rng.Intn(len(types))]
		g.nextIDs.power++
		g.powerUps = append(g.powerUps, &gPowerUp{
			ID:   g.nextIDs.power,
			Row:  row,
			Col:  col,
			Type: t,
		})
	}
}

func (g *gameSession) step(dt float64) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, p := range g.players {
		if p.Dead {
			continue
		}
		p.MoveAccum += dt
		delay := g.moveDelay(p)
		for p.MoveAccum >= delay {
			p.MoveAccum -= delay
			g.tryMove(p)
		}
		bombEdge := p.in.Bomb && !p.PrevBomb
		p.PrevBomb = p.in.Bomb
		if bombEdge {
			g.tryPlaceBomb(p)
		}
	}

	for i := 0; i < len(g.bombs); i++ {
		b := g.bombs[i]
		b.Timer -= dt
		if b.Timer > 0 {
			continue
		}
		g.bombs = append(g.bombs[:i], g.bombs[i+1:]...)
		i--
		g.detonateBomb(b)
	}

	g.checkMatchEndConditions()
}

// Match ends when someone reaches the revealed goal, or no players are left alive.
// Surviving players keep playing after others are eliminated (no instant win on last-man-standing).
func (g *gameSession) checkMatchEndConditions() {
	if !g.running {
		return
	}
	if g.exitRevealed {
		for _, p := range g.players {
			if p.Dead {
				continue
			}
			if p.Row == g.exitRow && p.Col == g.exitCol {
				g.broadcastGameOver(p)
				g.running = false
				return
			}
		}
	}
	alive := 0
	for _, p := range g.players {
		if !p.Dead {
			alive++
		}
	}
	if alive == 0 && len(g.players) > 0 {
		g.broadcastGameOver(nil)
		g.running = false
	}
}

func (g *gameSession) run() {
	ticker := time.NewTicker(time.Second / 60)
	defer ticker.Stop()
	for g.running {
		<-ticker.C
		g.mu.Lock()
		running := g.running
		g.mu.Unlock()
		if !running {
			break
		}
		g.step(1.0 / 60.0)
		g.frame++
		if g.frame%3 == 0 {
			g.broadcastState()
		}
		g.mu.Lock()
		run := g.running
		g.mu.Unlock()
		if !run {
			break
		}
	}
	g.finishLobby()
}

func (g *gameSession) broadcastState() {
	g.mu.Lock()
	payload := g.buildStatePayload()
	clients := g.snapshotAllBroadcastClients()
	g.mu.Unlock()

	data, err := json.Marshal(payload)
	if err != nil {
		return
	}
	for _, c := range clients {
		if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("state write: %v", err)
		}
	}
}

func (g *gameSession) buildStatePayload() map[string]any {
	players := make([]map[string]any, 0, len(g.players))
	for _, p := range g.players {
		players = append(players, map[string]any{
			"slot":      p.Slot,
			"nickname":  p.Nickname,
			"col":       p.Col,
			"row":       p.Row,
			"x":         p.Col * tileSize,
			"y":         p.Row * tileSize,
			"lives":     p.Lives,
			"dead":      p.Dead,
			"maxBombs":  p.MaxBombs,
			"bombCount": p.BombCount,
			"radius":    p.Radius,
			"speed":     p.SpeedLevel,
			"inv":       time.Now().Before(p.InvUntil),
		})
	}
	bombs := make([]map[string]any, 0, len(g.bombs))
	for _, b := range g.bombs {
		bombs = append(bombs, map[string]any{
			"id": b.ID, "col": b.Col, "row": b.Row,
			"x": b.Col * tileSize, "y": b.Row * tileSize,
			"t": b.Timer, "r": b.Radius,
		})
	}
	pups := make([]map[string]any, 0, len(g.powerUps))
	for _, pu := range g.powerUps {
		pups = append(pups, map[string]any{
			"id": pu.ID, "col": pu.Col, "row": pu.Row,
			"x": pu.Col * tileSize, "y": pu.Row * tileSize,
			"type": pu.Type,
		})
	}
	deltas := g.tileDeltas
	g.tileDeltas = nil
	if deltas == nil {
		deltas = []tileDelta{}
	}
	return map[string]any{
		"type":        "state",
		"tileDeltas": deltas,
		"players":     players,
		"bombs":       bombs,
		"powerUps":    pups,
	}
}

// broadcastGameOver must be called with g.mu held (same as checkWinOrEnd).
func (g *gameSession) broadcastGameOver(winner *gPlayer) {
	var w map[string]any
	if winner != nil && !winner.Dead {
		w = map[string]any{"slot": winner.Slot, "nickname": winner.Nickname}
	}
	payload := map[string]any{"type": "game_over", "winner": w}
	data, _ := json.Marshal(payload)
	clients := g.snapshotAllBroadcastClients()
	for _, c := range clients {
		_ = c.conn.WriteMessage(websocket.TextMessage, data)
	}
}

func (g *gameSession) snapshotAllBroadcastClients() []*client {
	out := make([]*client, 0, len(g.players)+len(g.spectators))
	for _, p := range g.players {
		if p.client != nil {
			out = append(out, p.client)
		}
	}
	out = append(out, g.spectators...)
	return out
}

func (g *gameSession) addSpectator(c *client) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, x := range g.spectators {
		if x == c {
			return
		}
	}
	g.spectators = append(g.spectators, c)
}

func (g *gameSession) removeSpectator(c *client) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for i, x := range g.spectators {
		if x == c {
			g.spectators = append(g.spectators[:i], g.spectators[i+1:]...)
			return
		}
	}
}

func (g *gameSession) sendSpectateStart(c *client) {
	g.mu.Lock()
	tilesCopy := copyTiles(g.tiles)
	er, ec := g.exitRow, g.exitCol
	list := make([]map[string]any, 0, len(g.players))
	for _, p := range g.players {
		list = append(list, map[string]any{"slot": p.Slot, "nickname": p.Nickname})
	}
	g.mu.Unlock()
	payload := map[string]any{
		"type":    "spectate_start",
		"slot":    -1,
		"players": list,
		"tiles":   tilesCopy,
		"tilePx":  tileSize,
		"exitRow": er,
		"exitCol": ec,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}
	_ = c.conn.WriteMessage(websocket.TextMessage, data)
}

func (g *gameSession) disconnectPlayer(c *client) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, p := range g.players {
		if p.client == c && !p.Dead {
			p.Dead = true
			p.Lives = 0
		}
	}
	g.checkMatchEndConditions()
}

func (g *gameSession) finishLobby() {
	g.lobby.mu.Lock()
	g.lobby.game = nil
	g.lobby.phase = "waiting"
	g.lobby.countdownStart = time.Time{}
	g.lobby.startDeadline = time.Time{}
	for c := range g.lobby.clients {
		c.slot = -1
		c.inGame = false
		c.spectator = false
	}
	g.lobby.chatLog = nil
	g.lobby.mu.Unlock()
	g.lobby.disconnectExcessWaitingClients()
}

func (g *gameSession) handleInput(c *client, msg map[string]any) {
	g.mu.Lock()
	defer g.mu.Unlock()
	var p *gPlayer
	for _, pl := range g.players {
		if pl.client == c {
			p = pl
			break
		}
	}
	if p == nil || p.Dead {
		return
	}
	// JSON numbers are float64
	p.in.Up = truthy(msg["up"])
	p.in.Down = truthy(msg["down"])
	p.in.Left = truthy(msg["left"])
	p.in.Right = truthy(msg["right"])
	p.in.Bomb = truthy(msg["bomb"])
}

// handleSlimeHit applies client-reported slime contact (slimes are simulated client-side only).
func (g *gameSession) handleSlimeHit(c *client) {
	g.mu.Lock()
	if !g.running {
		g.mu.Unlock()
		return
	}
	var p *gPlayer
	for _, pl := range g.players {
		if pl.client == c {
			p = pl
			break
		}
	}
	if p == nil || p.Dead {
		g.mu.Unlock()
		return
	}
	g.hitPlayerDamage(p)
	g.checkMatchEndConditions()
	stillRunning := g.running
	g.mu.Unlock()
	if stillRunning {
		g.broadcastState()
	}
}

func truthy(v any) bool {
	switch t := v.(type) {
	case bool:
		return t
	case float64:
		return t != 0
	default:
		return false
	}
}

func (g *gameSession) sendGameStart() {
	g.mu.Lock()
	defer g.mu.Unlock()
	list := make([]map[string]any, 0, len(g.players))
	for _, p := range g.players {
		list = append(list, map[string]any{"slot": p.Slot, "nickname": p.Nickname})
	}
	for _, p := range g.players {
		payload := map[string]any{
			"type":    "game_start",
			"slot":    p.Slot,
			"players": list,
			"tiles":   g.tiles,
			"tilePx":  tileSize,
			"exitRow": g.exitRow,
			"exitCol": g.exitCol,
		}
		data, err := json.Marshal(payload)
		if err != nil {
			continue
		}
		_ = p.client.conn.WriteMessage(websocket.TextMessage, data)
	}
}
