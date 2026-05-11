// examples.js
// Τρέξιμο:
//   node examples.js --update    # δημιουργεί/ενημερώνει expected αρχεία
//   node examples.js             # συγκρίνει τρέχον output με expected και δείχνει ✓/✗

const fs = require('fs')
const path = require('path')
const crosswordSolver = require('./crosswordSolver.js')

const SNAP_DIR = path.join(__dirname, 'expected')
if (!fs.existsSync(SNAP_DIR)) fs.mkdirSync(SNAP_DIR)

function slugify(name) {
  return name.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/(^-|-$)/g, '')
}

// capture console.log από τη συνάρτηση
function runAndCapture(puzzle, words) {
  const old = console.log
  const logs = []
  console.log = (msg = '') => logs.push(String(msg))
  try {
    crosswordSolver(puzzle, words)
  } finally {
    console.log = old
  }
  return logs.join('\n')
}

function writeExpected(name, content) {
  const file = path.join(SNAP_DIR, `${slugify(name)}.txt`)
  fs.writeFileSync(file, content, 'utf8')
  return file
}

function readExpected(name) {
  const file = path.join(SNAP_DIR, `${slugify(name)}.txt`)
  if (!fs.existsSync(file)) return null
  return fs.readFileSync(file, 'utf8')
}

function diff(a, b) {
  // απλό line-by-line diff μήκους
  if (a === b) return null
  return `--- expected\n${b}\n--- got\n${a}\n`
}

const puzzleBig = `...1...........
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

const CASES = [
  {
    name: 'ex1 basic',
    puzzle: '2001\n0..0\n1000\n0..0',
    words: ['casa', 'alan', 'ciao', 'anta'],
  },
  {
    name: 'ex2 summer set',
    puzzle: puzzleBig,
    words: [
      'sun','sunglasses','suncream','swimming','bikini','beach',
      'icecream','tan','deckchair','sand','seaside','sandals',
    ],
  },
  {
    name: 'ex3 groceries',
    puzzle: `..1.1..1...
10000..1000
..0.0..0...
..1000000..
..0.0..0...
1000..10000
..0.1..0...
....0..0...
..100000...
....0..0...
....0......`,
    words: ['popcorn','fruit','flour','chicken','eggs','vegetables','pasta','pork','steak','cheese'],
  },
  {
    name: 'ex4 reversed still unique',
    puzzle: puzzleBig,
    words: [
      'sun','sunglasses','suncream','swimming','bikini','beach',
      'icecream','tan','deckchair','sand','seaside','sandals',
    ].reverse(),
  },

  // Negative cases
  { name: 'err mismatch words vs slots', puzzle: '2001\n0..0\n2000\n0..0', words: ['casa','alan','ciao','anta'] },
  { name: 'err start greater than 2', puzzle: '0001\n0..0\n3000\n0..0', words: ['casa','alan','ciao','anta'] },
  { name: 'err duplicate word', puzzle: '2001\n0..0\n1000\n0..0', words: ['casa','casa','ciao','anta'] },
  { name: 'err empty puzzle', puzzle: '', words: ['casa','alan','ciao','anta'] },
  { name: 'err wrong puzzle type', puzzle: 123, words: ['casa','alan','ciao','anta'] },
  { name: 'err wrong words type', puzzle: '2001\n0..0\n1000\n0..0', words: 123 },
  { name: 'err multiple solutions', puzzle: '2000\n0...\n0...\n0...', words: ['abba','assa'] },
  { name: 'err no solution', puzzle: '2001\n0..0\n1000\n0..0', words: ['aaab','aaac','aaad','aaae'] },
]

const UPDATE = process.argv.includes('--update')

let failures = 0
for (const tc of CASES) {
  const out = runAndCapture(tc.puzzle, tc.words)

  if (UPDATE) {
    // γράφουμε/ανανεώνουμε το expected σύμφωνα με το τρέχον output
    const file = writeExpected(tc.name, out)
    console.log(`✓ wrote ${path.basename(file)}`)
  } else {
    const expected = readExpected(tc.name)
    if (expected == null) {
      console.log(`✗ ${tc.name} (missing expected; run with --update)`)
      failures++
      continue
    }
    const d = diff(out, expected)
    if (d) {
      console.log(`✗ ${tc.name}`)
      console.log(d)
      failures++
    } else {
      console.log(`✓ ${tc.name}`)
    }
  }
}

if (!UPDATE) {
  if (failures > 0) {
    process.exitCode = 1
  } else {
    console.log('All good ✅')
  }
}
