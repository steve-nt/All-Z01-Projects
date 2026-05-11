async function series(functions) {
	const results = []

	for (const fn of functions) {
		const result = await fn()
		results.push(result)
	}

	return results
}
