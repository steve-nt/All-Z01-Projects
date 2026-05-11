export function getURL(dataSet) {
  return dataSet.match(/https?:\/\/[^\s]+/g) || [];
}

export function greedyQuery(dataSet) {
  return dataSet.match(/https?:\/\/[^\s?]*\?[^&\s]+=(?:[^&\s])*&[^&\s]+=(?:[^&\s])*&[^&=\s]+=(?:[^&\s])*(?:&[^&=\s]+=(?:[^&\s])*)*(?=\s|$)/g) || [];
}

export function notSoGreedy(dataSet) {
  return dataSet.match(/https?:\/\/[^\s?]*\?[^&\s]+=(?:[^&\s])*&[^&\s]+=(?:[^&\s])*(?:&[^&=\s]+=(?:[^&\s])*)?(?=\s|$)/g) || [];
}

