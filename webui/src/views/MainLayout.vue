<script>
export default {
	name: "MainLayout",
};
</script>

<template>
	<!--
		Layout principale della SPA: Sidebar a sinistra, pannello destro (router-view) a destra.
		Utilizzato come componente padre per le rotte autenticate.
	-->
	<div class="main-layout d-flex" style="height: 100vh; overflow: hidden">
		<!-- Barra laterale con la lista delle conversazioni -->
		<Sidebar />

		<!-- Pannello destro: il contenuto cambia in base alla rotta figlia -->
		<div class="flex-grow-1" style="overflow: hidden">
			<router-view />
		</div>
	</div>
</template>

<style scoped>
.main-layout {
	width: 100vw;
}
</style>
