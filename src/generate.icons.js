const svg2png = require('svg2png')
const toIco = require('to-ico')
const fs = require('fs')

// Promise version of readFile
const readFile = file => new Promise((resolve, reject) => {
	fs.readFile(file, (err, data) => {
		if (err) return reject(err)
		resolve(data)
	})
})
// Promise version of writeFile
const writeFile = (file, data) => new Promise((resolve, reject) => {
	fs.writeFile(file, data, err => {
		if (err) return reject(err)
		resolve(true)
	})
})

module.exports = color => {
	color = color || '#508cb6'
	const changeColor = svg => (svg + '').replace(/class="fill-primary"/g, 'style="fill: ' + color + '"')

	// 128x128 png
	readFile(__dirname + '/assets/logo.svg')
		.then(svg => changeColor(svg))
		.then(svg => svg2png(svg, { width: 196, height: 196 }))
		.then(png => writeFile(__dirname + '/../icon-196.png', png))
		.then(() => console.log('icon-196.png generated'))
		.catch(err => console.error(err))

	// 16, 32, 48 icons
	readFile(__dirname + '/assets/logo.svg')
		.then(svg => changeColor(svg))
		.then(svg => Promise.all([
			svg2png(svg, { width: 16, height: 16 }),
			svg2png(svg, { width: 32, height: 32 }),
			svg2png(svg, { width: 48, height: 48 }),
		]))
		.then(pngs => toIco(pngs))
		.then(ico => writeFile(__dirname + '/../favicon.ico', ico))
		.then(() => console.log('favicon.ico generated'))
		.catch(err => console.error(err))
}

// Cli
if (require.main === module) {
	module.exports(process.argv[2])
}