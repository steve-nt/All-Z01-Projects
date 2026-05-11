export const format = {
    formatNumber(value) {
        if (value >= 1000000)
            return (value / 1000000).toFixed(1) + 'M';
        if (value >= 1000)
            return (value / 1000).toFixed(1) + 'K';
        return value.toString();
    },
    
    formatBytes(bytes) {
        if (bytes >= 1000000) {
            return (bytes / 1000000).toFixed(1) + 'MB';
        } else if (bytes >= 1000) {
            return (bytes / 1000).toFixed(1) + 'KB';
        }
        return bytes + 'B';
    }
};