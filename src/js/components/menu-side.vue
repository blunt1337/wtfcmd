<template>
	<ul class="menu vertical">
		<li :class="{ selected: selected_index === index }" v-for="(heading, index) of headings">
			<a :href="'#' + heading.anchor" @click.prevent="scrollTo(heading.element)">
				<component :is="'h' + heading.size" class="text-truncate">{{ heading.text }}</component>
			</a>
		</li>
	</ul>
</template>

<script>
const offset = 100

export default {
	data: () => ({
		headings: [],
		selected_index: null,
		content_div: null,
	}),
	methods: {
		listHeadings() {
			let headings = this.content_div.querySelectorAll('h1,h2,h3,h4')
			let res = []
			
			for (let i = 0, l = headings.length; i < l; i++) {
				let heading = headings[i]
				if (heading.hasAttribute('menu-ignore')) continue
				let anchor = heading.id = heading.id || ('section_' + i)
				let text = heading.textContent
				res.push({ anchor, text, element: heading, size: heading.tagName.substr(1) * 1 + 2 })
			}
			this.headings = res
		},
		onScroll() {
			let top = this.content_div.scrollTop || document.body.scrollTop || document.documentElement.scrollTop
			top += offset
			
			let headings = this.headings
			let lg = headings.length
			if (lg == 0) return this.selected_index = 0
			
			for (let i = 0; i < lg; i++) {
				if (headings[i].element.offsetTop >= top) {
					if (i === 0) return this.selected_index = 0
					return this.selected_index = i - 1
				}
			}
			return this.selected_index = lg - 1
		},
		scrollTo(element) {
			// Close sidebar
			this.$emit('sidebar')
			
			// Change url
			if (history.pushState) {
				history.pushState(null, null, '#' + element.id)
			} else {
				location.hash = '#' + element.id
			}
			
			// Scroll to
			this.content_div.scrollTop = document.body.scrollTop = document.documentElement.scrollTop = element.offsetTop - offset + 1
		}
	},
	mounted() {
		this.content_div = document.querySelector('.content')
		this.content_div.addEventListener('scroll', this.onScroll)
		this.$router.afterEach((to, from) => {
			setTimeout(() => {
				this.listHeadings()
				this.onScroll()
			}, 10)
		})
	},
	beforeDestroy() {
		this.content_div.removeEventListener('scroll', this.onScroll)
		//TODO: remove afterEach
	},
}
</script>

<style>
/* Fix scrollbar */
.sidebar .menu {
	overflow-y: auto;
	margin-right: -20px;
	padding-right: 20px;
}
</style>