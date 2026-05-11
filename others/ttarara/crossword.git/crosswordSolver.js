function crosswordSolver(puzzle, words) {

  if (typeof puzzle !== 'string' || !puzzle) {
    console.log('Error')
    return
  }
  if (!Array.isArray(words)) {
    console.log('Error')
    return
  }

 
  const seen = new Set()
  for (const w of words) {
    if (typeof w !== 'string' || w.length === 0) {
      console.log('Error')
      return
    }
    const low = w.toLowerCase()
    if (seen.has(low)) {
      console.log('Error')
      return
    }
    seen.add(low)
  }
  const WORDS = words.map(w => w.toLowerCase())

  // ---------- Parse puzzle ----------
  const rows = puzzle.split('\n')
  const H = rows.length
  if (H === 0) {
    console.log('Error')
    return
  }
  const W = rows[0].length
  if (W === 0) {
    console.log('Error')
    return
  }
  for (const r of rows) {
    if (r.length !== W) {
      console.log('Error')
      return
    }
    for (const ch of r) {
      if (ch !== '.' && !(ch >= '0' && ch <= '9')) {
        console.log('Error')
        return
      }
    }
  }

  const gridChar = rows.map(r => r.split('')) // '.' ή '0'..'9'
  const letters = Array.from({ length: H }, () => Array(W).fill(null))
  const refCount = Array.from({ length: H }, () => Array(W).fill(0))

  // ----------Find slots ----------
  function isBlock(r, c) { return gridChar[r][c] === '.' }
  function spanRight(r, c) {
    let len = 0
    for (let j = c; j < W && !isBlock(r, j); j++) len++
    return len
  }
  function spanDown(r, c) {
    let len = 0
    for (let i = r; i < H && !isBlock(i, c); i++) len++
    return len
  }

  const slots = [] // { id, dir:'H'|'V', cells:[{r,c}], len }
  const startsAt = Array.from({ length: H }, () => Array(W).fill(0))

  // Οριζόντια slots (μήκος >= 2)
  for (let i = 0; i < H; i++) {
    for (let j = 0; j < W; j++) {
      if (isBlock(i, j)) continue
      const leftIsBlock = (j === 0) || isBlock(i, j - 1)
      if (leftIsBlock) {
        const len = spanRight(i, j)
        if (len >= 2) {
          const cells = []
          for (let jj = j; jj < j + len; jj++) cells.push({ r: i, c: jj })
          slots.push({ id: slots.length, dir: 'H', cells, len })
          startsAt[i][j] += 1
        }
      }
    }
  }

  //Down slots (μήκος >= 2)
  for (let j = 0; j < W; j++) {
    for (let i = 0; i < H; i++) {
      if (isBlock(i, j)) continue
      const upIsBlock = (i === 0) || isBlock(i - 1, j)
      if (upIsBlock) {
        const len = spanDown(i, j)
        if (len >= 2) {
          const cells = []
          for (let ii = i; ii < i + len; ii++) cells.push({ r: ii, c: j })
          slots.push({ id: slots.length, dir: 'V', cells, len })
          startsAt[i][j] += 1
        }
      }
    }
  }

  // ---------- Check start ----------
  for (let i = 0; i < H; i++) {
    for (let j = 0; j < W; j++) {
      const ch = gridChar[i][j]
      if (ch === '.') continue
      const d = Number(ch)
      if (Number.isNaN(d)) { console.log('Error'); return }
      if (d > 2) { console.log('Error'); return }
      if (startsAt[i][j] !== d) { console.log('Error'); return }
    }
  }

  // ---------- Slots vs Words ----------
  if (slots.length !== WORDS.length) {
    console.log('Error')
    return
  }

  // ---------- Backtracking ----------
  slots.sort((a, b) => a.len - b.len)

  const used = Array(WORDS.length).fill(false)
  const wordsByLen = new Map()
  for (let idx = 0; idx < WORDS.length; idx++) {
    const L = WORDS[idx].length
    if (!wordsByLen.has(L)) wordsByLen.set(L, [])
    wordsByLen.get(L).push(idx)
  }

  let solutions = 0
  let solvedBoardSnapshot = null

  function canPlace(slot, word) {
    for (let k = 0; k < slot.len; k++) {
      const { r, c } = slot.cells[k]
      const ch = letters[r][c]
      if (ch !== null && ch !== word[k]) return false
    }
    return true
  }
  function place(slot, word) {
    for (let k = 0; k < slot.len; k++) {
      const { r, c } = slot.cells[k]
      if (letters[r][c] === null) letters[r][c] = word[k]
      refCount[r][c] += 1
    }
  }
  function unplace(slot) {
    for (let k = 0; k < slot.len; k++) {
      const { r, c } = slot.cells[k]
      refCount[r][c] -= 1
      if (refCount[r][c] === 0) letters[r][c] = null
    }
  }

  function dfs(si) {
    if (solutions > 1) return
    if (si === slots.length) {
      solutions += 1
      if (solutions === 1) {
        solvedBoardSnapshot = letters.map(row => row.slice())
      }
      return
    }
    const slot = slots[si]
    const candIdxs = wordsByLen.get(slot.len) || []
    for (const wi of candIdxs) {
      if (used[wi]) continue
      const w = WORDS[wi]
      if (!canPlace(slot, w)) continue
      used[wi] = true
      place(slot, w)
      dfs(si + 1)
      unplace(slot)
      used[wi] = false
      if (solutions > 1) return
    }
  }

  dfs(0)

  if (solutions !== 1 || !solvedBoardSnapshot) {
    console.log('Error')
    return
  }

  // ---------- Print ----------
  let out = ''
  for (let i = 0; i < H; i++) {
    let line = ''
    for (let j = 0; j < W; j++) {
      if (gridChar[i][j] === '.') {
        line += '.'
      } else {
        const ch = solvedBoardSnapshot[i][j]
        if (!ch || ch.length !== 1) { console.log('Error'); return }
        line += ch
      }
    }
    out += line + (i + 1 < H ? '\n' : '')
  }
  console.log(out)
}

// ---------------
// Για να τρέξεις γρήγορα τα παραδείγματα του ελέγχου, αποσχολίασε ένα block κάθε φορά:

//
// Example 1
//const puzzle1 = '2001\n0..0\n1000\n0..0'
//const words1 = ['casa', 'alan', 'ciao', 'anta']
//crosswordSolver(puzzle1, words1)
//// αναμενόμενο:
// casa
// i..l
// anta
// o..n


/*
// Example 2
const puzzle2 = `...1...........
..1000001000...
...0....0......
.1......0...1..
.0....100000000
100000..0...0..
.0.....1001000.
.0.1....0.0....
.10000000.0....
.0.0......0....
.0.0.....100...
...0......0....
..........0....`
const words2 = [
  'sun','sunglasses','suncream','swimming','bikini','beach',
  'icecream','tan','deckchair','sand','seaside','sandals',
]
crosswordSolver(puzzle2, words2)
*/

/*
// Example 3
const puzzle3 = `..1.1..1...
10000..1000
..0.0..0...
..1000000..
..0.0..0...
1000..10000
..0.1..0...
....0..0...
..100000...
....0..0...
....0......`
const words3 = ['popcorn','fruit','flour','chicken','eggs','vegetables','pasta','pork','steak','cheese']
crosswordSolver(puzzle3, words3)
*/

/*]
// Example 4 (reverse order)
const puzzle4 = `...1...........
..1000001000...
...0....0......
.1......0...1..
.0....100000000
100000..0...0..
.0.....1001000.
.0.1....0.0....
.10000000.0....
.0.0......0....
.0.0.....100...
...0......0....
..........0....`
const words4 = [
  'sun','sunglasses','suncream','swimming','bikini','beach',
  'icecream','tan','deckchair','sand','seaside','sandals',
].reverse()
crosswordSolver(puzzle4, words4)
*/

/*
// Mismatch πλήθους slots-λέξεων
const puzzle5 = '2001\n0..0\n2000\n0..0'
const words5 = ['casa', 'alan', 'ciao', 'anta']
crosswordSolver(puzzle5, words5) // Error
*/

/*
// Αριθμός start > 2 (αδύνατο)
const puzzle6 = '0001\n0..0\n3000\n0..0'
const words6 = ['casa', 'alan', 'ciao', 'anta']
crosswordSolver(puzzle6, words6) // Error
*/

/*
// Διπλή λέξη
const puzzle7 = '2001\n0..0\n1000\n0..0'
const words7 = ['casa', 'casa', 'ciao', 'anta']
crosswordSolver(puzzle7, words7) // Error
*/

/*
// Κενό puzzle
const puzzle8 = ''
const words8 = ['casa', 'alan', 'ciao', 'anta']
crosswordSolver(puzzle8, words8) // Error
*/

/*
// Λανθαστός τύπος puzzle
const puzzle9 = 123
const words9 = ['casa', 'alan', 'ciao', 'anta']
crosswordSolver(puzzle9, words9) // Error
*/

/*
// Λανθαστός τύπος words
const puzzle10 = '2001\n0..0\n1000\n0..0'
const words10 = 123
crosswordSolver(puzzle10, words10) // Error
*/

/*
// Πολλαπλές λύσεις
const puzzle11 = '2000\n0...\n0...\n0...'
const words11 = ['abba', 'assa']
crosswordSolver(puzzle11, words11) // Error
*/

/*
// Καμία λύση
const puzzle12 = '2001\n0..0\n1000\n0..0'
const words12 = ['aaab', 'aaac', 'aaad', 'aaae']
crosswordSolver(puzzle12, words12) // Error
*/

// ---- exports ----
if (typeof module !== 'undefined') {
  module.exports = crosswordSolver
}

// ---- optional: run a single demo only when called directly ----
if (require.main === module) {
  // const p = '2001\n0..0\n1000\n0..0'
  // const w = ['casa', 'alan', 'ciao', 'anta']
  // crosswordSolver(p, w)
}

