function throttle(fn, wait) {
	let lastCallTime = 0;

	return function(...args) {
		const now = Date.now();

		if (now - lastCallTime >= wait) {
			lastCallTime = now;
			fn.apply(this, args);
		}
	};
}

function opThrottle(fn, wait, options = {}) {
	let timer = null;
	let lastArgs = null;
	let lastThis = null;
	let previous = 0;

	// Handle option defaults correctly
	let leading, trailing;

	if ('leading' in options && !('trailing' in options)) {
		leading = options.leading;
		trailing = false;
	} else if ('trailing' in options && !('leading' in options)) {
		leading = false;
		trailing = options.trailing;
	} else if ('leading' in options && 'trailing' in options) {
		leading = options.leading;
		trailing = options.trailing;
	} else {
		leading = false;
		trailing = true;
	}

	const later = function() {
		previous = leading === false ? 0 : Date.now();
		timer = null;
		if (trailing) {
			fn.apply(lastThis, lastArgs);
		}
		lastArgs = lastThis = null;
	};

	return function(...args) {
		const now = Date.now();
		if (!previous && leading === false) previous = now;
		const remaining = wait - (now - previous);

		lastThis = this;
		lastArgs = args;

		if (remaining <= 0 || remaining > wait) {
			if (timer) {
				clearTimeout(timer);
				timer = null;
			}
			previous = now;
			fn.apply(this, args);
			lastArgs = lastThis = null;
		} else if (!timer && trailing) {
			timer = setTimeout(later, remaining);
		}
	};
}
