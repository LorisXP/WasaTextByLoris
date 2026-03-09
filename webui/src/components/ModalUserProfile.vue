<script>
import auth from "../services/auth.js";
import axios from "../services/axios.js";

export default {
	name: "ModalUserProfile",

	props: {
		/** Controlla se la modale è visibile */
		visible: {
			type: Boolean,
			required: true,
		},
	},

	emits: ["close", "photo-updated"],

	data() {
		return {
			// ── Anteprima / stato upload ───────────────────────────────────────
			/** File selezionato dall'utente (oggetto File) */
			selectedFile: null,
			/** data URL per l'anteprima dell'immagine scelta */
			previewSrc: null,

			// ── Stato operazione ──────────────────────────────────────────────
			saving: false,
			errorMsg: null,
			successMsg: null,
		};
	},

	computed: {
		myUserID() {
			return auth.state.userID;
		},
		myUserName() {
			return auth.state.userName;
		},
		/** Mostra la foto corrente dell'utente (da auth state) come data URI */
		currentPhotoSrc() {
			if (auth.state.userPhoto) {
				return `data:image/jpeg;base64,${auth.state.userPhoto}`;
			}
			return null;
		},
		/** Iniziale del nome utente, usata come fallback avatar */
		initials() {
			return this.myUserName ? this.myUserName.charAt(0).toUpperCase() : "?";
		},
	},

	methods: {
		// ── LIFECYCLE MODALE ──────────────────────────────────────────────────

		close() {
			this.reset();
			this.$emit("close");
		},

		reset() {
			this.selectedFile = null;
			this.previewSrc = null;
			this.saving = false;
			this.errorMsg = null;
			this.successMsg = null;
			if (this.$refs.fileInput) {
				this.$refs.fileInput.value = "";
			}
		},

		// ── SELEZIONE FOTO ────────────────────────────────────────────────────

		pickPhoto() {
			this.$refs.fileInput.click();
		},

		onFileSelected(event) {
			const file = event.target.files[0];
			if (!file) return;

			// Mostra anteprima locale immediata
			const reader = new FileReader();
			reader.onload = (e) => {
				this.previewSrc = e.target.result;
			};
			reader.readAsDataURL(file);

			this.selectedFile = file;
			this.errorMsg = null;
			this.successMsg = null;

			// Reset input per permettere la rideselezione dello stesso file
			event.target.value = "";
		},

		removePreview() {
			this.selectedFile = null;
			this.previewSrc = null;
			this.errorMsg = null;
		},

		// ── SALVATAGGIO ───────────────────────────────────────────────────────

		/**
		 * Invia la nuova foto al backend.
		 * PUT /api/users/{userID}/me/photo   (multipart/form-data, campo "photo")
		 */
		async savePhoto() {
			if (!this.selectedFile || this.saving) return;

			this.saving = true;
			this.errorMsg = null;
			this.successMsg = null;

			try {
				const formData = new FormData();
				formData.append("photo", this.selectedFile);

				const response = await axios.put(
					`/api/users/${this.myUserID}/me/photo`,
					formData,
					{ headers: { "Content-Type": "multipart/form-data" } }
				);

				if (response.status === 200) {
					// Aggiorna la foto nello stato auth affinché l'avatar della sidebar
					// si aggiorni in tempo reale senza bisogno di refresh.
					// Il backend ha salvato la foto con base64; la leggiamo dall'anteprima.
					if (this.previewSrc) {
						const base64 = this.previewSrc.split(",")[1] || null;
						auth.setUserPhoto(base64);
					}
					this.successMsg = "Foto profilo aggiornata con successo!";
					this.selectedFile = null;
					this.previewSrc = null;
					this.$emit("photo-updated");
				} else {
					this.errorMsg = "Impossibile aggiornare la foto. Riprova.";
				}
			} catch (e) {
				this.errorMsg =
					e.response?.data || "Errore imprevisto durante l'upload. Riprova.";
				console.error("Errore aggiornamento foto profilo:", e);
			} finally {
				this.saving = false;
			}
		},
	},

	watch: {
		visible(newVal) {
			if (newVal) {
				this.reset();
			}
		},
	},
};
</script>

<template>
	<Teleport to="body">
		<!-- Backdrop -->
		<div
			v-if="visible"
			class="modal-backdrop-custom"
			@click.self="close"
		></div>

		<!-- Contenitore modale centrato -->
		<div
			v-if="visible"
			class="modal-container-custom d-flex align-items-center justify-content-center"
			role="dialog"
			aria-modal="true"
			aria-labelledby="modal-user-profile-title"
		>
			<div
				class="card shadow-lg border rounded-4 overflow-hidden"
				style="width: 100%; max-width: 420px"
			>
				<!-- ── HEADER ──────────────────────────────────────────────────── -->
				<div class="card-header d-flex align-items-center justify-content-between px-4 py-3 bg-primary text-white border-0">
					<h5 id="modal-user-profile-title" class="mb-0 fw-semibold fs-6">
						Il mio profilo
					</h5>
					<button
						type="button"
						class="btn-close btn-close-white"
						aria-label="Chiudi"
						:disabled="saving"
						@click="close"
					></button>
				</div>

				<!-- ── BODY ────────────────────────────────────────────────────── -->
				<div class="card-body px-4 py-4">

					<!-- Nome utente -->
					<p class="text-center text-secondary small mb-4">
						Connesso come <strong class="text-body">{{ myUserName }}</strong>
					</p>

					<!-- Foto corrente / anteprima nuova foto -->
					<div class="d-flex flex-column align-items-center gap-3 mb-4">
						<!-- Avatar corrente (o anteprima) -->
						<div class="position-relative">
							<img
								v-if="previewSrc"
								:src="previewSrc"
								alt="Anteprima nuova foto"
								class="rounded-circle border border-2 border-primary"
								style="width: 96px; height: 96px; object-fit: cover"
							/>
							<img
								v-else-if="currentPhotoSrc"
								:src="currentPhotoSrc"
								:alt="myUserName"
								class="rounded-circle border"
								style="width: 96px; height: 96px; object-fit: cover"
							/>
							<div
								v-else
								class="rounded-circle bg-primary d-flex align-items-center justify-content-center text-white fw-bold border"
								style="width: 96px; height: 96px; font-size: 2rem"
								:aria-label="myUserName"
							>
								{{ initials }}
							</div>

							<!-- Badge "anteprima" sovrapposto -->
							<span
								v-if="previewSrc"
								class="badge bg-primary position-absolute bottom-0 end-0 small"
								style="font-size: 0.6rem"
							>
								Anteprima
							</span>
						</div>

						<!-- Bottoni azione foto -->
						<div class="d-flex gap-2">
							<button
								type="button"
								class="btn btn-outline-primary btn-sm d-flex align-items-center gap-1"
								:disabled="saving"
								@click="pickPhoto"
							>
								<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" fill="currentColor" class="bi bi-camera" viewBox="0 0 16 16">
									<path d="M15 12a1 1 0 0 1-1 1H2a1 1 0 0 1-1-1V6a1 1 0 0 1 1-1h1.172a3 3 0 0 0 2.12-.879l.83-.828A1 1 0 0 1 6.827 3h2.344a1 1 0 0 1 .707.293l.828.828A3 3 0 0 0 12.828 5H14a1 1 0 0 1 1 1zM2 4a2 2 0 0 0-2 2v6a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V6a2 2 0 0 0-2-2h-1.172a2 2 0 0 1-1.414-.586l-.828-.828A2 2 0 0 0 9.172 2H6.828a2 2 0 0 0-1.414.586l-.828.828A2 2 0 0 1 3.172 4z"/>
									<path d="M8 11a2.5 2.5 0 1 1 0-5 2.5 2.5 0 0 1 0 5m0 1a3.5 3.5 0 1 0 0-7 3.5 3.5 0 0 0 0 7M3 6.5a.5.5 0 1 1-1 0 .5.5 0 0 1 1 0"/>
								</svg>
								{{ previewSrc ? "Cambia foto" : "Scegli foto" }}
							</button>
							<button
								v-if="previewSrc"
								type="button"
								class="btn btn-outline-secondary btn-sm"
								:disabled="saving"
								@click="removePreview"
							>
								Annulla selezione
							</button>
						</div>

						<!-- Input file nascosto -->
						<input
							ref="fileInput"
							type="file"
							accept="image/*"
							class="d-none"
							@change="onFileSelected"
						/>
					</div>

					<!-- Messaggi di esito -->
					<div
						v-if="errorMsg"
						class="alert alert-danger py-2 small mb-0"
						role="alert"
					>
						{{ errorMsg }}
					</div>
					<div
						v-if="successMsg"
						class="alert alert-success py-2 small mb-0"
						role="alert"
					>
						{{ successMsg }}
					</div>
				</div>

				<!-- ── FOOTER ──────────────────────────────────────────────────── -->
				<div class="card-footer d-flex justify-content-end gap-2 px-4 py-3 border-top bg-body-tertiary">
					<button
						type="button"
						class="btn btn-outline-secondary"
						:disabled="saving"
						@click="close"
					>
						Chiudi
					</button>
					<button
						type="button"
						class="btn btn-primary d-flex align-items-center gap-2"
						:disabled="!selectedFile || saving"
						@click="savePhoto"
					>
						<span
							v-if="saving"
							class="spinner-border spinner-border-sm"
							role="status"
							aria-hidden="true"
						></span>
						<svg
							v-else
							xmlns="http://www.w3.org/2000/svg"
							width="14"
							height="14"
							fill="currentColor"
							class="bi bi-cloud-upload"
							viewBox="0 0 16 16"
							aria-hidden="true"
						>
							<path fill-rule="evenodd" d="M4.406 1.342A5.53 5.53 0 0 1 8 0c2.69 0 4.923 2 5.166 4.579C14.758 4.804 16 6.137 16 7.773 16 9.569 14.502 11 12.687 11H10a.5.5 0 0 1 0-1h2.688C13.979 10 15 8.988 15 7.773c0-1.216-1.02-2.228-2.313-2.228h-.5v-.5C12.188 2.825 10.328 1 8 1a4.53 4.53 0 0 0-2.941 1.1c-.757.652-1.153 1.438-1.153 2.055v.448l-.445.049C2.064 4.805 1 5.952 1 7.318 1 8.785 2.23 10 3.781 10H6a.5.5 0 0 1 0 1H3.781C1.708 11 0 9.366 0 7.318c0-1.763 1.266-3.223 2.942-3.593.143-.863.698-1.723 1.464-2.383"/>
							<path fill-rule="evenodd" d="M7.646 4.146a.5.5 0 0 1 .708 0l3 3a.5.5 0 0 1-.708.708L8.5 5.707V14.5a.5.5 0 0 1-1 0V5.707L5.354 7.854a.5.5 0 1 1-.708-.708z"/>
						</svg>
						Salva foto
					</button>
				</div>
			</div>
		</div>
	</Teleport>
</template>

<style scoped>
.modal-backdrop-custom {
	position: fixed;
	inset: 0;
	background: rgba(0, 0, 0, 0.5);
	z-index: 10000;
}

.modal-container-custom {
	position: fixed;
	inset: 0;
	z-index: 10001;
	padding: 1rem;
}
</style>
