export function ionOut(str) {
  return str.match(/\w+tion\b/g)?.map(word => word.slice(0, -3)) || [];
}
