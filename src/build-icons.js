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

/**
 * @param	{Function}	editSvg		called after the svg is loaded, to replace some text, etc with new_svg = editSvg(svg)
 */
module.exports = (entry, output_path, editSvg) => {
	entry = entry || 'logo.svg'
	output_path = output_path || __dirname
	editSvg = editSvg || (svg => svg)

	// 128x128 png
	readFile(entry)
		.then(svg => editSvg(svg))
		.then(svg => svg2png(svg, { width: 196, height: 196 }))
		.then(png => writeFile(output_path + '/icon-196.png', png))
		.then(() => console.log('icon-196.png generated'))
		.catch(err => console.error(err))

	// 16, 32, 48 icons
	readFile(entry)
		.then(svg => editSvg(svg))
		.then(svg => Promise.all([
			svg2png(svg, { width: 16, height: 16 }),
			svg2png(svg, { width: 32, height: 32 }),
			svg2png(svg, { width: 48, height: 48 }),
		]))
		.then(pngs => toIco(pngs))
		.then(ico => writeFile(output_path + '/favicon.ico', ico))
		.then(() => console.log('favicon.ico generated'))
		.catch(err => console.error(err))
}