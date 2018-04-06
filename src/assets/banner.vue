<template>
	<div class="banner">
		<div class="color" v-for="n in ['a', 'b', 'c']" :class="n">
			<div class="rect"></div>
			<div class="triangle"></div>
		</div>
		<h1 menu-ignore>super alias your commands<br />in seconds.</h1>
		<a href="#getting-started" class="btn btn-outline-dark">Get started</a>
	</div>
</template>

<style lang="scss">
@import "sass/variables";

.banner {
	width: 100%;
	height: 300px;
	overflow: hidden;
	position: relative;
	background: mix($primary, #FFF, 3%);
	
	.color {
		position: absolute;
		top: 0;
		width: 100%;
		height: 100%;
		
		.rect {
			width: 50%;
			height: 100%;
			position: absolute;
			top: 0;
			left: -50%;
		}
		
		.triangle {
			width: 0;
			height: 0;
			border-style: solid;
			border-width: 304px 0 0 304px;
			border-color: transparent;
			margin-left: -2px;
			margin-top: -2px;
		}
		
		&.a {
			left: calc(30% - 280px);
			.rect {
				background: mix($primary, #FFF, 20%);
			}
			.triangle {
				border-left-color: mix($primary, #FFF, 20%);
			}
		}
		&.b {
			transform: scaleX(-1);
			left: calc(-70% + 280px);
			.rect {
				background: mix($primary, #FFF, 60%);
			}
			.triangle {
				border-left-color: mix($primary, #FFF, 60%);
			}
		}
		&.c {
			transform: scaleX(-1);
			left: calc(-50% + 280px);
			.rect {
				background: $primary;
			}
			.triangle {
				border-left-color: $primary;
			}
		}
	}
	
	h1 {
		&:before {
			content: none;
		}
		position: absolute;
		left: 1em;
		top: 1em;
	}
	a {
		position: absolute;
		left: 2em;
		bottom: 2em;
	}
}
</style>