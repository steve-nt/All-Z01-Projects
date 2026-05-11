async function queryServers(serverName, q) {
	const urls = [
		`/${serverName}?q=${q}`,
		`/${serverName}_backup?q=${q}`
	]

	return Promise.race(urls.map(url => getJSON(url)))
}

async function gougleSearch(q) {
	const timeoutPromise = new Promise((_, reject) =>
		setTimeout(() => reject(new Error('timeout')), 80)
	)

	const searchPromise = (async () => {
		const servers = ['web', 'image', 'video']
		const results = await Promise.all(
			servers.map(server => queryServers(server, q))
		)
		return {
			web: results[0],
			image: results[1],
			video: results[2]
		}
	})()

	return Promise.race([searchPromise, timeoutPromise])
}
