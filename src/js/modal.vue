<template>
	<transition name="modal">
		<div class="modal show" @click.self="$emit('close')">
			<div class="modal-dialog">
				<div class="modal-content">
					<div class="modal-header" v-if="title">
						<h5 class="modal-title">{{ title }}</h5>
						<button type="button" class="close" v-if="closable" @click="$emit('close')">&times;</button>
					</div>
					<div class="modal-body">
						<button type="button" class="close" v-if="closable && !title" @click="$emit('close')">&times;</button>
						<slot/>
					</div>
					<slot name="footer"/>
				</div>
			</div>
		</div>
	</transition>
</template>

<script>
module.exports = {
	name: 'modal',
	props: {
		title: String,
		closable: { type: Boolean, default: true },
	},
}
</script>

<style lang="scss">
@import "sass/variables";

.modal.show {
	display: block !important;
}

// Transition
.modal {
	background-color: rgba(0, 0, 0, .3);
	transition: opacity $transition-duration ease;
}
.modal-dialog {
	transition: all $transition-duration ease;
}
.modal-content {
	border: none !important;
	box-shadow: 0 2px 6px 0 rgba(0, 0, 0, .2);
}

.modal-enter, .modal-leave-active {
	opacity: 0;
}
.modal-enter .modal-dialog, .modal-leave-active .modal-dialog {
	transform: scale(1.1) !important;
}
</style>