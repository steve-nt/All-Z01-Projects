function interpolation({ step, start, end, callback, duration }) {
	if (step <= 0 || duration <= 0) return;

	const interval = duration / step;
	let i = 0;

	function tick() {
		const x = start + (i / step) * (end - start);
		const y = ((i + 1) / step) * duration;

		callback([x, y]);

		i++;
		if (i < step) {
			setTimeout(tick, interval);
		}
	}

	setTimeout(tick, interval);
}
