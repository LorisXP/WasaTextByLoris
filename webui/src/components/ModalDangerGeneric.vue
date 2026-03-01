<script>
export default {
	name: "ModalDangerGeneric",

	props: {
		/** Controlla se la modale è visibile */
		visible: {
			type: Boolean,
			required: true,
		},
		/** Titolo della modale */
		title: {
			type: String,
			default: "Errore",
		},
		/** Descrizione / messaggio di errore */
		description: {
			type: String,
			default: "",
		},
	},

	emits: ["close"],

	methods: {
		close() {
			this.$emit("close");
		},
	},
};
</script>

<template>
	<Teleport to="body">
		<!-- Backdrop -->
		<div
			v-if="visible"
			class="modal-danger-backdrop"
			@click.self="close"
		></div>

		<!-- Contenitore centrato -->
		<div
			v-if="visible"
			class="modal-danger-container d-flex align-items-center justify-content-center"
			role="alertdialog"
			aria-modal="true"
			:aria-labelledby="'modal-danger-title'"
			:aria-describedby="'modal-danger-desc'"
		>
			<div
				class="card shadow-lg border border-danger rounded-4 overflow-hidden"
				style="width: 100%; max-width: 420px;"
			>
				<!-- Header -->
				<div class="card-header d-flex align-items-center gap-2 px-4 py-3 bg-danger text-white border-0">
					<!-- Icona pericolo -->
					<svg
						xmlns="http://www.w3.org/2000/svg"
						width="18"
						height="18"
						fill="currentColor"
						class="bi bi-exclamation-triangle-fill flex-shrink-0"
						viewBox="0 0 16 16"
						aria-hidden="true"
					>
						<path d="M8.982 1.566a1.13 1.13 0 0 0-1.96 0L.165 13.233c-.457.778.091 1.767.98 1.767h13.713c.889 0 1.438-.99.98-1.767zM8 5c.535 0 .954.462.9.995l-.35 3.507a.552.552 0 0 1-1.1 0L7.1 5.995A.905.905 0 0 1 8 5m.002 6a1 1 0 1 1 0 2 1 1 0 0 1 0-2"/>
					</svg>
					<h5
						id="modal-danger-title"
						class="mb-0 fw-semibold fs-6"
					>
						{{ title }}
					</h5>
					<button
						type="button"
						class="btn-close btn-close-white ms-auto"
						aria-label="Chiudi"
						@click="close"
					></button>
				</div>

				<!-- Body -->
				<div class="card-body px-4 py-4">
					<p
						id="modal-danger-desc"
						class="mb-0 text-body"
						style="white-space: pre-line;"
					>
						{{ description }}
					</p>
				</div>

				<!-- Footer -->
				<div class="card-footer d-flex justify-content-end px-4 py-3 border-top bg-body-tertiary">
					<button
						type="button"
						class="btn btn-danger"
						@click="close"
					>
						Chiudi
					</button>
				</div>
			</div>
		</div>
	</Teleport>
</template>

<style scoped>
.modal-danger-backdrop {
	position: fixed;
	inset: 0;
	background: rgba(0, 0, 0, 0.5);
	z-index: 10000;
}

.modal-danger-container {
	position: fixed;
	inset: 0;
	z-index: 10001;
	padding: 1rem;
}
</style>
