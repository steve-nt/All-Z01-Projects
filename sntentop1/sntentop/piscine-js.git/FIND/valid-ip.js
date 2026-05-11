export function findIP(str) {
  
  const octet = '(?:25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])';
  
  const port = '(?::(?:6553[0-5]|655[0-2][0-9]|65[0-4][0-9]{2}|6[0-4][0-9]{3}|[1-5][0-9]{4}|[1-9][0-9]{1,3}|[0-9]))?';
  
  const regex = new RegExp(`(?<![\\d\\.])${octet}\\.${octet}\\.${octet}\\.${octet}${port}(?!\\d|:)`, 'g');
  
  return Array.from(str.matchAll(regex), m => m[0]);
}