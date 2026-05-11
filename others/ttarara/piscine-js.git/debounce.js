function debounce(fn, wait) {
	let timeout;
	return function(...args) {
		clearTimeout(timeout);
		timeout = setTimeout(() => fn.apply(this, args), wait);
	};
}

function opDebounce(fn, wait, options = {}) {
	let timeout;
	let invoked = false;

	return function(...args) {
		const context = this;

		const callNow = options.leading && !invoked;

		clearTimeout(timeout);

		timeout = setTimeout(() => {
			timeout = null;
			invoked = false;
			if (!options.leading) fn.apply(context, args);
		}, wait);

		if (callNow) {
			invoked = true;
			fn.apply(context, args);
		}
	};
}
