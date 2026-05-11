// Merges two objects; the second object overrides the first
const defaultCurry = (obj1) => (obj2) => ({ ...obj1, ...obj2 });

// Maps over entries [key, value] and returns a new object
const mapCurry = (fn) => (obj) =>
  Object.fromEntries(Object.entries(obj).map(fn));

// Reduces object entries starting from an initial value
const reduceCurry = (fn) => (obj, initialValue) =>
  Object.entries(obj).reduce(fn, initialValue);

// Filters object entries based on a predicate function
const filterCurry = (fn) => (obj) =>
  Object.fromEntries(Object.entries(obj).filter(fn));

// --- Personnel Specific Logic ---

// Returns total score (shooting + piloting) only for Force users
const reduceScore = (personnel, initial = 0) =>
  reduceCurry((acc, [_, v]) =>
    v.isForceUser ? acc + v.pilotingScore + v.shootingScore : acc
  )(personnel, initial);

// Returns users who use the Force AND have a shootingScore >= 80
const filterForce = (personnel) =>
  filterCurry(([_, v]) => 
    v.isForceUser && v.shootingScore >= 80
  )(personnel);

// Returns a new object with the average score added to each person
const mapAverage = (personnel) =>
  mapCurry(([k, v]) => [
    k,
    { ...v, averageScore: (v.pilotingScore + v.shootingScore) / 2 },
  ])(personnel);