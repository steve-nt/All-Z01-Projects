async function getJSON(path, params = {}) {
	const url = path.includes('?') ? path : path + '?';
	const searchParams = new URLSearchParams();

	Object.entries(params).forEach(([key, value]) => {
		searchParams.append(key, value);
	});

	const fullUrl = url + searchParams.toString();

	const response = await fetch(fullUrl);
	if (!response.ok) {
		throw new Error(response.statusText);
	}

	const result = await response.json();

	if (result.error) {
		throw new Error(result.error);
	}

	return result.data ?? result;
}
