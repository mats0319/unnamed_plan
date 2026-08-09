<template>
  <div class="stacked-cards">
    <slot />
  </div>
</template>

<script lang="ts" setup>
</script>

<style lang="less" scoped>
.stacked-cards,
.stacked-cards::before,
.stacked-cards::after {
	border-radius: 0.5rem;
	box-shadow: 2px 1px rgba(0, 0, 0, 0.15);
}

.stacked-cards::before,
.stacked-cards::after {
	content: "";
	position: absolute;
	width: 100%;
	height: 100%;
	background-color: rgba(240, 239, 226, 0.5);
}

.stacked-cards::before {
	@media (min-width: 1280px) {
		left: 1rem;
		top: 1rem;
	}
	left: 0.5rem;
	top: 0.5rem;

	z-index: -1;
}

.stacked-cards::after {
	@media (min-width: 1280px) {
		left: 2rem;
		top: 2rem;
	}
	left: 1rem;
	top: 1rem;
	z-index: -2;
}
</style>
