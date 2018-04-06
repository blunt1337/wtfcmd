<template>
	<div class="body-wrapper" :class="{ 'sidebar-toggled': sidebar }">
		<div class="sidebar">
			<div class="container-fluid">
				<div>
					<img markup-inline src="../assets/logo.svg" width="90"/>
					<h1 class="text-uppercase">{{ project_name }}</h1>
					<em class="text-muted">{{ project_description }}</em>
				</div>
				
				<hr align="right"/>
				
				<!-- Menu -->
				<menu-side @sidebar="sidebar = false"/>
			</div>
			<a class="sidebar-closer" @click="sidebar = false">
				<img markup-inline src="../assets/arrow-left.svg" height="1em"/>
			</a>
		</div>
		<div class="content">
			<div class="container-fluid">
				<!-- Menu -->
				<menu-top @sidebar="sidebar = true"/>
				
				<!-- Pages -->
				<router-view></router-view>
				
				<!-- Footer -->
				<hr />
				© 2017-{{ (new Date).getFullYear() }} <a target="_blanc" :href="'https://github.com/' + project_user">{{ project_user }}</a><br /><br />
			</div>
		</div>
		<div class="sidebar-overlay" @click="sidebar = false"></div>
	</div>
</template>

<script>
import MenuTop from './components/menu-top'
import MenuSide from './components/menu-side'

export default {
	data: () => ({
		sidebar: false,
		project_name: process.env.PROJECT_NAME,
		project_user: process.env.PROJECT_USER,
		project_description: process.env.PROJECT_DESCRIPTION
	}),
	components: {
		MenuTop,
		MenuSide,
	},
}
</script>

<style lang="scss">
@import "~@/sass/index";
$sidebar_width: 300px;

.body-wrapper {
	// Large desktop
	@include media-breakpoint-up(xl) {
		// For 1300 + 300px
		max-width: map-get($grid-breakpoints, 'xl');
		margin: auto;
		position: relative;
		
		.sidebar {
			left: calc((100vw - 1300px) / 2);
		}
	}
	// Large large desktop
	@media (min-width: #{map-get($grid-breakpoints, 'xl') + $sidebar_width}) {
		left: -$sidebar_width / 2;
		
		.sidebar {
			left: calc((100vw - #{map-get($grid-breakpoints, 'xl') - $sidebar_width}) / 2 - #{$sidebar_width});
		}
	}
	
	&.sidebar-toggled {
		overflow: hidden;
	}
}

.sidebar {
	position: fixed;
	left: -$sidebar_width;
	
	width: $sidebar_width;
	height: 100%;
	max-width: 100%;
	padding: 2em 0 0;
	
	.container-fluid {
		display: flex;
		height: 100%;
		flex-direction: column;
	}

	z-index: 3000;
	transition: left $transition-duration;
	
	text-align: right;

	.sidebar-toggled & {
		left: 0;
	}
	
	// Desktop
	@include media-breakpoint-up(md) {
		left: 0;
	}
	
	hr {
		margin-left: auto;
	}
}

.content {
	width: 100%;
	padding-top: 1em;

	transition: margin-left $transition-duration;

	.sidebar-toggled & {
		margin-left: $sidebar_width;
	}

	// Desktop
	@include media-breakpoint-up(md) {
		padding-top: 2em;
		
		width: auto;
		margin-left: $sidebar_width;
		padding-left: 2em;
	}
}

.sidebar-overlay {
	display: none;
	
	.sidebar-toggled & {
		display: block;
		
		&:after {
			content: '';
			display: block;
			position: fixed;
			top: 0;
			right: 0;
			bottom: 0;
			left: 0;
			z-index: 2000;
			
			background: #FFF;
			opacity: 0.75;
		}
		
		@include media-breakpoint-up(md) {
			display: none;
		}
	}
}

.sidebar-closer {
	position: absolute;
	top: 1em;
	left: 1em;
	z-index: 99999;
	
	@include media-breakpoint-up(md) {
		display: none;
	}
}
</style>