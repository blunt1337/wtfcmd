const MarkdownIt = require('markdown-it')
const hljs = require('highlight.js')

module.exports = {
	wrapper: 'div',
	raw: true,
	html: true,
	xhtmlOut: true,
	linkify: true,
	typographer: true,
	breaks: true,
	highlight(code, lang) {
		if (lang && hljs.getLanguage(lang)) {
			try {
				return '<pre class="hljs"><code>' +
				hljs.highlight(lang, code, true).value +
				'</code></pre>';
			} catch (__) {}
		}
		return '<pre class="hljs"><code>' + md.utils.escapeHtml(code) + '</code></pre>';
	},
	preprocess(md, source) {
		// Table class
		md.renderer.rules.table_open = () => '<table class="table table-striped">\n'

		// Fix for router-link
		md.renderer.rules.link_open = function(tokens, idx, options, env, self) {
			let token = tokens[idx]
			
			// Add rel="noopener" for target="_blank"
			if (token.attrGet('target') === '_blank') {
				token.attrSet('rel', 'noopener')
			}
			
			// Convert <a href> to <router-link to>
			let href_index = token.attrIndex('href')
			if (href_index >= 0 && !/:\/\//.test(token.attrs[href_index][1])) {
				token.tag = 'router-link'
				token.attrs[href_index][0] = 'to'
				
				// Find matching close tag
				for (let i = idx + 1, l = tokens.length; i < l; i++) {
					if (tokens[i].tag === 'a' && tokens[i].nesting === -1) {
						tokens[i].tag = 'router-link'
						break
					}
				}
			}
			
			return self.renderToken(tokens, idx, options)
		}

		// Codes
		let old_code_inline = md.renderer.rules.code_inline
		md.renderer.rules.code_inline = function(tokens, idx, options, env, self) {
			tokens[idx].attrJoin('class', 'hljs inline')
			return old_code_inline(tokens, idx, options, env, self)
		}

		let old_code_block = md.renderer.rules.code_block
		md.renderer.rules.code_block = function(tokens, idx, options, env, self) {
			tokens[idx].attrJoin('class', 'hljs')
			return old_code_block(tokens, idx, options, env, self)
		}

		// Tasks plugin
		md.core.ruler.push('task', state => {
			const class_map = {
				'x': 'done',
				'*': 'done',
				'!': 'warn',
				'?': 'test',
				' ': 'todo',
			}
			const title_map = {
				'x': 'Done',
				'*': 'Done',
				'!': 'Urgent!',
				'?': 'To test',
				' ': 'To do',
			}
			
			// Loop blocks (find inline after list_item_open)
			let tokens = state.tokens
			let lg = tokens.length
			for (let i = 0; i < lg; i++) {
				let token = tokens[i]
				
				// Check inline + li before
				if (token.type !== 'inline' || token.hidden) continue
				let li = false
				for (let j = i - 1; j >= 0; j--) {
					if (tokens[j].type === 'list_item_open') {
						li = tokens[j]
						break
					}
					if (!tokens[j].hidden) {
						break
					}
				}
				if (!li) continue
				
				// Start with a text
				let first_child = token.children[0]
				if (!first_child || first_child.type !== 'text') continue
				
				// Text is a task
				let match = first_child.content.match(/^([\t ]*)\[([x*!? ])\][\t ]?/i)
				if (!match) continue
				
				// Edit li
				li.attrSet('class', 'md-task ' + class_map[match[2]])
				li.attrSet('title', title_map[match[2]])
				
				// Remove markup from content
				first_child.content = first_child.content.substr(match[0].length)
			}
		})

		// Anchor plugin
		md.core.ruler.push('anchor', state => {
			const ids = {}
			const tokens = state.tokens
			
			tokens.filter(token => token.type === 'heading_open').forEach(token => {
				// Aggregate the next token children text.
				const title = tokens[tokens.indexOf(token) + 1].children
					.filter(token => token.type === 'text' || token.type === 'code_inline')
					.reduce((acc, t) => acc + t.content, '')
				
				let id = token.attrGet('id')
				if (id == null) {
					id = title.replace(/[\W\s-]+/g, '-').toLowerCase()
					if (ids[id]) {
						id += '-' + (++ids[id])
					} else {
						ids[id] = 1
					}
					token.attrPush(['id', id])
				}
			})
		})
		
		return source
	},
}