<template>
	<ul class="menu top">
		<!-- Mobile sidebar opener -->
		<li class="sidebar-opener"><img markup-inline @click="$emit('sidebar')" src="../../assets/menu.svg" height="1em"/></li>
		
		<!-- Pages -->
		<li v-for="route in routes" :class="{ selected: $route.path === route.path }" v-if="route.name">			
			<router-link :to="route.path">{{ route.name }}</router-link>
		</li>
		
		<!-- Other -->
		<li><a :href="'https://github.com/' + project_user + '/' + project_repo"><img markup-inline src="../../assets/github.svg" height="1em"/></a></li>
	</ul>
</template>

<script>
import routes from '../routes'

export default {
	data: () => ({
		routes,
		project_user: PROJECT_USER,
		project_repo: PROJECT_REPO,
	}),
}
</script>

<style lang="scss">
@import "sass/variables";

.sidebar-opener {
	@include media-breakpoint-up(md) {
		display: none !important;
	}
}

.menu.top {
	position: sticky;
	top: 0;
	padding-top: 1em;
	margin-top: -1em;
	z-index: 1000;
	background: $body-bg;
	
	@include media-breakpoint-up(md) {
		padding-top: 2em;
		margin-top: -2em;
	}
}
</style>