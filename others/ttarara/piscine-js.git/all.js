async function all(obj) {
	return new Promise((resolve, reject) => {
		const keys = Object.keys(obj)
		const result = {}
		let remaining = keys.length

		if (remaining === 0) {
			return resolve(result)
		}

		keys.forEach(key => {
			Promise.resolve(obj[key])
				.then(value => {
					result[key] = value
					remaining--
					if (remaining === 0) {
						resolve(result)
					}
				})
				.catch(reject)
		})
	})
}
