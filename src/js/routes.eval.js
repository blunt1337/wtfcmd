const fs = require('fs')
const page_dir = __dirname + '/../pages'

// Convert a filename into a displayed name
const displayify = name => name.replace(/^[0-9]+-/, '').replace(/_/g, ' ')

// Convert a filename into an url
const urlify = name => '/' + name.replace(/^[0-9]+-/, '').replace(/-/g, ' ')

// Build the component import code
const importCode = filename => {
	let mod = JSON.stringify('@/pages/' + filename)
	return `() => {
		setLoading(true)
		return import(${mod}).then(res => {
			setLoading(false)
			return res
		})
	}`
}

export default function () {
	return new Promise((resolve, reject) => {
		// List files from page folder
		fs.readdir(page_dir, (err, files) => {
			if (err) return reject(err)
			
			let deps = []
			let routes = `
				const loaderClasses = document.querySelector('.top-border .loading').classList
				const setLoading = bool => loaderClasses.toggle('d-none', !bool)
				
				export default [
				{
					name: 'home',
					path: '/',
					component: ${importCode('index')}
				},`
			
			// For each md, html, vue files (not index or 404)
			for (let file of files) {
				if (/\.(md|html|vue)$/.test(file)) {
					let basename = file.substr(0, file.lastIndexOf('.'))
					if (basename !== '404' && basename !== 'index') {
						// Don't forget dependencies
						deps.push(page_dir + '/' + file)
						
						// Result code
						let name = JSON.stringify(displayify(basename))
						let path = JSON.stringify(urlify(basename))
						let mod = importCode(file)
						
						routes += `{
							name: ${name},
							path: ${path},
							component: ${mod},
						},`
					}
				}
			}
			
			// 404 page
			routes += `
				{
					path: '*',
					component: ${importCode('404')}
				}
			]`
			
			resolve({
				code: routes,
				dependencies: deps,
				contextDependencies: [page_dir],
				cacheable: true,
			})
		})
	})
}