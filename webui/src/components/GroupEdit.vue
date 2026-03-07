<script>
import auth from "../services/auth.js";
import axios from "../services/axios.js";
import ModalDangerGeneric from "./ModalDangerGeneric.vue";
import LoadingSpinner from "./LoadingSpinner.vue";

/** Millisecondi di debounce sulla ricerca utenti */
const SEARCH_DEBOUNCE_MS = 420;

export default {
	name: "GroupEdit",

	components: {
		ModalDangerGeneric,
		LoadingSpinner,
	},

	props: {
		/** Controlla se il pannello è visibile */
		visible: {
			type: Boolean,
			required: true,
		},
		/** ID del gruppo da modificare */
		groupID: {
			type: Number,
			required: true,
		},
	},

	emits: ["close", "group-updated", "group-left"],

	data() {
		return {
			// ── Dati del gruppo ──────────────────────────────────────────────
			group: null,           // GroupInfoResponse completo
			loading: false,
			errormsg: null,

			// ── Modifica nome (inline) ──────────────────────────────────────
			editingName: false,
			editedName: "",
			savingName: false,
			nameError: null,

			// ── Modifica foto ───────────────────────────────────────────────
			savingPhoto: false,
			photoError: null,

			// ── Aggiunta membro ─────────────────────────────────────────────
			showAddMember: false,
			searchQuery: "",
			searchResults: [],
			searching: false,
			searchError: null,
			addingMember: false,
			addError: null,
			_searchTimer: null,

			// ── Kick membro ─────────────────────────────────────────────────
			kickingUser: null,      // userName in fase di kick

			// ── Uscita / Eliminazione ───────────────────────────────────────
			leaving: false,
			leaveError: null,

			// ── Validazione ────────────────────────────────────────────────
			validationError: { visible: false, title: "Input non valido", description: "" },
		};
	},

	computed: {
		myUserID() {
			return auth.state.userID;
		},
		myUserName() {
			return auth.state.userName;
		},

		/** True se l'utente corrente è admin del gruppo */
		isAdmin() {
			return this.group && this.group.administrator === this.myUserName;
		},

		/** Sorgente base64 della foto del gruppo */
		groupPhotoSrc() {
			if (!this.group || !this.group.photo) return null;
			return `data:image/jpeg;base64,${this.group.photo}`;
		},

		/** Iniziale del nome per avatar di fallback */
		groupInitial() {
			if (!this.group || !this.group.name) return "G";
			return this.group.name.charAt(0).toUpperCase();
		},

		/** Membri filtrati (escluso se stesso, opzionale) */
		members() {
			return this.group ? this.group.members || [] : [];
		},
	},

	watch: {
		visible(newVal) {
			if (newVal) {
				this.loadGroupInfo();
			} else {
				this.resetState();
			}
		},
	},

	methods: {
		// ─── CARICAMENTO INFO GRUPPO ────────────────────────────────────

		/**
		 * GET /api/users/{userID}/groups/{groupID}
		 * Carica le informazioni complete del gruppo.
		 */
		async loadGroupInfo() {
			this.errormsg = null;
			this.loading = true;
			try {
				const res = await axios.get(
					`/api/users/${this.myUserID}/groups/${this.groupID}`
				);
				if (res.status === 200) {
					this.group = res.data;
				} else {
					this.errormsg = "Impossibile caricare il gruppo";
				}
			} catch (e) {
				this.errormsg = e.toString();
			} finally {
				this.loading = false;
			}
		},

		// ─── MODIFICA NOME ─────────────────────────────────────────────

		/** Entra in modalità di editing del nome */
		startEditName() {
			this.editedName = this.group.name;
			this.editingName = true;
			this.nameError = null;
			this.$nextTick(() => {
				const el = this.$refs.nameInput;
				if (el) el.focus();
			});
		},

		/** Annulla la modifica del nome */
		cancelEditName() {
			this.editingName = false;
			this.nameError = null;
		},

		/**
		 * PATCH /api/groups/{groupID}/name
		 * Body: { userID: int, name: string }
		 */
		async saveName() {
			const trimmed = this.editedName.trim();
			const vName = this.$validator(trimmed, null, 3, 15, "string");
			if (!vName.success) {
				this.showValidationError("Il nome del gruppo deve essere tra 3 e 15 caratteri.");
				return;
			}
			this.savingName = true;
			this.nameError = null;
			try {
				await axios.patch(
					`/api/groups/${this.groupID}/name`,
					{ userID: this.myUserID, name: trimmed }
				);
				this.group.name = trimmed;
				this.editingName = false;
				this.$emit("group-updated");
			} catch (e) {
				this.nameError = e.response?.data || e.toString();
			} finally {
				this.savingName = false;
			}
		},

		// ─── MODIFICA FOTO ─────────────────────────────────────────────

		/** Apre il selettore file per la foto */
		pickPhoto() {
			this.$refs.photoInput.click();
		},

		/**
		 * PUT /api/groups/{groupID}/photo
		 * Multipart form: userID (string) + photo (file)
		 */
		async onPhotoSelected(event) {
			const file = event.target.files[0];
			if (!file) return;
			this.savingPhoto = true;
			this.photoError = null;
			try {
				const formData = new FormData();
				formData.append("userID", String(this.myUserID));
				formData.append("photo", file);

				await axios.put(
					`/api/groups/${this.groupID}/photo`,
					formData
				);

				// Aggiorna la foto in locale leggendo il file come base64
				const reader = new FileReader();
				reader.onload = () => {
					// rimuovi il prefisso "data:image/...;base64,"
					const base64 = reader.result.split(",")[1];
					this.group.photo = base64;
					this.$emit("group-updated");
				};
				reader.readAsDataURL(file);
			} catch (e) {
				this.photoError = e.response?.data || e.toString();
			} finally {
				this.savingPhoto = false;
				// Reset input per consentire la stessa selezione
				event.target.value = "";
			}
		},

		// ─── AGGIUNTA MEMBRO ───────────────────────────────────────────

		/** Toggle pannello di ricerca utenti */
		toggleAddMember() {
			this.showAddMember = !this.showAddMember;
			if (!this.showAddMember) {
				this.searchQuery = "";
				this.searchResults = [];
				this.searchError = null;
			} else {
				this.$nextTick(() => {
					const el = this.$refs.memberSearchInput;
					if (el) el.focus();
				});
			}
		},

		/** Ricerca utenti con debounce */
		onSearchInput() {
			clearTimeout(this._searchTimer);
			const q = this.searchQuery.trim();
			if (q.length < 3) {
				this.searchResults = [];
				return;
			}
			this.searching = true;
			this._searchTimer = setTimeout(() => this.searchUsers(q), SEARCH_DEBOUNCE_MS);
		},

		/**
		 * GET /api/users/{userID}/other/{userName}
		 * Cerca utenti il cui nome corrisponde.
		 */
		async searchUsers(query) {
			this.searchError = null;
			const vQuery = this.$validator(query, /^[a-z]+[0-9]*$/, 3, 15, "string");
			if (!vQuery.success) {
				this.showValidationError("Il nome utente deve essere di 3-15 caratteri, lettere minuscole seguite da numeri opzionali (es. mario42).");
				this.searching = false;
				return;
			}
			try {
				const res = await axios.get(
					`/api/users/${this.myUserID}/other/${encodeURIComponent(query)}`
				);
				if (res.status === 200) {
					// Filtra utenti già membri del gruppo
					const memberNames = new Set(this.members.map(m => m.userName));
					this.searchResults = (res.data || []).filter(
						u => !memberNames.has(u.userName)
					);
				}
			} catch (e) {
				this.searchError = e.response?.data || e.toString();
			} finally {
				this.searching = false;
			}
		},

		/**
		 * POST /api/groups/{groupID}/users
		 * Body: { userID: string, userNames: [string] }
		 */
		async addMember(userName) {
			this.addingMember = true;
			this.addError = null;
			try {
				await axios.post(
					`/api/groups/${this.groupID}/users`,
					{
						userID: String(this.myUserID),
						userNames: [userName],
					}
				);
				// Ricarica le info del gruppo per aggiornare la lista membri
				await this.loadGroupInfo();
				// Rimuovi il risultato dalla ricerca
				this.searchResults = this.searchResults.filter(u => u.userName !== userName);
				this.$emit("group-updated");
			} catch (e) {
				this.addError = e.response?.data || e.toString();
			} finally {
				this.addingMember = false;
			}
		},

		// ─── KICK MEMBRO ───────────────────────────────────────────────

		/**
		 * DELETE /api/users/{userID}/groups/{groupID}/member/{userName}
		 * Rimuove un membro dal gruppo (solo admin).
		 */
		async kickMember(userName) {
			this.kickingUser = userName;
			try {
				await axios.delete(
					`/api/users/${this.myUserID}/groups/${this.groupID}/member/${encodeURIComponent(userName)}`
				);
				// Ricarica per aggiornare la lista membri
				await this.loadGroupInfo();
				this.$emit("group-updated");
			} catch (e) {
				alert("Errore rimozione: " + (e.response?.data || e.toString()));
			} finally {
				this.kickingUser = null;
			}
		},

		// ─── USCITA / ELIMINAZIONE GRUPPO ──────────────────────────────

		/**
		 * Esci dal gruppo (non admin) o elimina il gruppo (admin).
		 * - Non-admin: DELETE /api/users/{userID}/groups/{groupID}
		 * - Admin:     DELETE /api/groups/{groupID}/users/{userID}
		 */
		async leaveOrDelete() {
			const action = this.isAdmin ? "eliminare il gruppo" : "uscire dal gruppo";
			if (!confirm(`Sicuro di voler ${action}?`)) return;

			this.leaving = true;
			this.leaveError = null;
			try {
				if (this.isAdmin) {
					// Elimina il gruppo intero
					await axios.delete(
						`/api/groups/${this.groupID}/users/${this.myUserID}`
					);
				} else {
					// Esce dal gruppo
					await axios.delete(
						`/api/users/${this.myUserID}/groups/${this.groupID}`
					);
				}
				this.$emit("group-left");
			} catch (e) {
				this.leaveError = e.response?.data || e.toString();
			} finally {
				this.leaving = false;
			}
		},

		// ─── UTILITY ───────────────────────────────────────────────────

		close() {
			this.$emit("close");
		},

		resetState() {
			this.group = null;
			this.loading = false;
			this.errormsg = null;
			this.editingName = false;
			this.editedName = "";
			this.savingName = false;
			this.nameError = null;
			this.savingPhoto = false;
			this.photoError = null;
			this.showAddMember = false;
			this.searchQuery = "";
			this.searchResults = [];
			this.searching = false;
			this.searchError = null;
			this.addingMember = false;
			this.addError = null;
			this.kickingUser = null;
			this.leaving = false;
			this.leaveError = null;
			clearTimeout(this._searchTimer);
		},

		/** Mostra la modale di errore di validazione */
		showValidationError(description) {
			this.validationError = { visible: true, title: "Input non valido", description };
		},

		/** Restituisce la sorgente base64 della foto di un membro */
		memberPhotoSrc(member) {
			if (!member.photo) return null;
			return `data:image/jpeg;base64,${member.photo}`;
		},
	},
};
</script>

<template>
	<Teleport to="body">
		<!--
			Backdrop semitrasparente.
			Click su di esso chiude il pannello.
		-->
		<div
			v-if="visible"
			class="modal-backdrop-custom"
			@click.self="close"
		></div>

		<!-- Contenitore centrato -->
		<div
			v-if="visible"
			class="modal-container-custom d-flex align-items-center justify-content-center"
			role="dialog"
			aria-modal="true"
			aria-labelledby="group-edit-title"
		>
			<div
				class="card shadow-lg border rounded-4 overflow-hidden"
				style="width: 100%; max-width: 720px; max-height: 92vh; display: flex; flex-direction: column;"
			>
				<!-- ── HEADER ─────────────────────────────────────────────────── -->
				<div
					class="card-header d-flex align-items-center justify-content-between px-4 py-3 bg-primary text-white border-0"
				>
					<h5 id="group-edit-title" class="mb-0 fw-semibold fs-6">
						Modifica gruppo
					</h5>
					<button
						type="button"
						class="btn-close btn-close-white"
						aria-label="Chiudi"
						@click="close"
					></button>
				</div>

				<!-- ── BODY ───────────────────────────────────────────────────── -->
				<div class="card-body overflow-y-auto px-4 py-4" style="overflow-y: auto; flex: 1 1 auto;">

					<!-- Caricamento in corso -->
					<div v-if="loading" class="text-center py-5">
						<LoadingSpinner />
					</div>

					<!-- Errore di caricamento -->
					<template v-else-if="errormsg">
						<ModalDangerGeneric :visible="true" :description="errormsg" @close="errormsg = null" />
					</template>

					<!-- Contenuto principale: due colonne -->
					<div v-else-if="group" class="row gx-4 gy-4">

						<!-- ═══ COLONNA SINISTRA: Gruppo ═══════════════════════ -->
						<div class="col-12 col-md-5">
							<h6 class="fw-bold text-uppercase text-secondary mb-3" style="font-size: 0.75rem; letter-spacing: 0.08em;">
								Gruppo
							</h6>

							<!-- Foto del gruppo -->
							<div class="text-center mb-4">
								<div class="position-relative d-inline-block">
									<img
										v-if="groupPhotoSrc"
										:src="groupPhotoSrc"
										alt="Foto gruppo"
										class="rounded-circle border"
										width="96"
										height="96"
										style="object-fit: cover"
									/>
									<div
										v-else
										class="rounded-circle bg-success d-flex align-items-center justify-content-center text-white fw-bold mx-auto"
										style="width: 96px; height: 96px; font-size: 2rem"
									>
										{{ groupInitial }}
									</div>

									<!-- Icona matita per cambiare foto (solo admin) -->
									<button
										v-if="isAdmin"
										type="button"
										class="btn btn-sm btn-light rounded-circle position-absolute shadow-sm"
										style="bottom: 0; right: 0; width: 30px; height: 30px; padding: 0;"
										title="Cambia foto"
										:disabled="savingPhoto"
										@click="pickPhoto"
									>
										<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" fill="currentColor" viewBox="0 0 16 16">
											<path d="M15.502 1.94a.5.5 0 0 1 0 .706L14.459 3.69l-2-2L13.502.646a.5.5 0 0 1 .707 0l1.293 1.293zm-1.75 2.456-2-2L4.939 9.21a.5.5 0 0 0-.121.196l-.805 2.414a.25.25 0 0 0 .316.316l2.414-.805a.5.5 0 0 0 .196-.12l6.813-6.814z"/>
											<path fill-rule="evenodd" d="M1 13.5A1.5 1.5 0 0 0 2.5 15h11a1.5 1.5 0 0 0 1.5-1.5v-6a.5.5 0 0 0-1 0v6a.5.5 0 0 1-.5.5h-11a.5.5 0 0 1-.5-.5v-11a.5.5 0 0 1 .5-.5H9a.5.5 0 0 0 0-1H2.5A1.5 1.5 0 0 0 1 2.5z"/>
										</svg>
									</button>
								</div>

								<!-- Spinner / errore foto -->
								<div v-if="savingPhoto" class="text-center mt-2">
									<span class="spinner-border spinner-border-sm text-primary" role="status"></span>
									<span class="ms-1 small text-secondary">Caricamento…</span>
								</div>
								<div v-if="photoError" class="text-danger small mt-2">{{ photoError }}</div>

								<!-- Input file nascosto -->
								<input
									ref="photoInput"
									type="file"
									accept="image/jpeg,image/png,image/webp"
									class="d-none"
									@change="onPhotoSelected"
								/>
							</div>

							<!-- Nome del gruppo -->
							<div class="mb-3">
								<div v-if="!editingName" class="d-flex align-items-center gap-2">
									<span class="fw-semibold fs-5">{{ group.name }}</span>
									<button
										v-if="isAdmin"
										type="button"
										class="btn btn-sm btn-link text-secondary p-0"
										title="Modifica nome"
										@click="startEditName"
									>
										<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" fill="currentColor" viewBox="0 0 16 16">
											<path d="M15.502 1.94a.5.5 0 0 1 0 .706L14.459 3.69l-2-2L13.502.646a.5.5 0 0 1 .707 0l1.293 1.293zm-1.75 2.456-2-2L4.939 9.21a.5.5 0 0 0-.121.196l-.805 2.414a.25.25 0 0 0 .316.316l2.414-.805a.5.5 0 0 0 .196-.12l6.813-6.814z"/>
											<path fill-rule="evenodd" d="M1 13.5A1.5 1.5 0 0 0 2.5 15h11a1.5 1.5 0 0 0 1.5-1.5v-6a.5.5 0 0 0-1 0v6a.5.5 0 0 1-.5.5h-11a.5.5 0 0 1-.5-.5v-11a.5.5 0 0 1 .5-.5H9a.5.5 0 0 0 0-1H2.5A1.5 1.5 0 0 0 1 2.5z"/>
										</svg>
									</button>
								</div>

								<!-- Editing inline del nome -->
								<div v-else class="d-flex align-items-center gap-2">
									<input
										ref="nameInput"
										v-model="editedName"
										type="text"
										class="form-control form-control-sm"
										style="max-width: 200px"
										minlength="3"
										maxlength="15"
										:disabled="savingName"
										@keyup.enter="saveName"
										@keyup.escape="cancelEditName"
									/>
									<button
										type="button"
										class="btn btn-sm btn-success"
										:disabled="savingName"
										title="Salva"
										@click="saveName"
									>
										<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" fill="currentColor" viewBox="0 0 16 16">
											<path d="M13.854 3.646a.5.5 0 0 1 0 .708l-7 7a.5.5 0 0 1-.708 0l-3.5-3.5a.5.5 0 1 1 .708-.708L6.5 10.293l6.646-6.647a.5.5 0 0 1 .708 0"/>
										</svg>
									</button>
									<button
										type="button"
										class="btn btn-sm btn-outline-secondary"
										:disabled="savingName"
										title="Annulla"
										@click="cancelEditName"
									>
										<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" fill="currentColor" viewBox="0 0 16 16">
											<path d="M4.646 4.646a.5.5 0 0 1 .708 0L8 7.293l2.646-2.647a.5.5 0 0 1 .708.708L8.707 8l2.647 2.646a.5.5 0 0 1-.708.708L8 8.707l-2.646 2.647a.5.5 0 0 1-.708-.708L7.293 8 4.646 5.354a.5.5 0 0 1 0-.708"/>
										</svg>
									</button>
								</div>
								<div v-if="nameError" class="text-danger small mt-1">{{ nameError }}</div>
							</div>

							<!-- ID del gruppo -->
							<div class="text-secondary small">
								<span class="text-body-tertiary">#{{ group.groupID }}</span>
							</div>
						</div>

						<!-- ═══ COLONNA DESTRA: Membri ════════════════════════ -->
						<div class="col-12 col-md-7">
							<h6 class="fw-bold text-uppercase text-secondary mb-3" style="font-size: 0.75rem; letter-spacing: 0.08em;">
								Membri
								<span class="badge bg-secondary ms-1">{{ members.length }}</span>
							</h6>

							<!-- Pulsante "Aggiungi membro" (solo admin) -->
							<div v-if="isAdmin" class="mb-3">
								<button
									type="button"
									class="btn btn-outline-primary btn-sm d-flex align-items-center gap-2"
									@click="toggleAddMember"
								>
									<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
										<path d="M12.5 16a3.5 3.5 0 1 0 0-7 3.5 3.5 0 0 0 0 7m.5-5v1h1a.5.5 0 0 1 0 1h-1v1a.5.5 0 0 1-1 0v-1h-1a.5.5 0 0 1 0-1h1v-1a.5.5 0 0 1 1 0m-2-6a3 3 0 1 1-6 0 3 3 0 0 1 6 0M8 7a2 2 0 1 0 0-4 2 2 0 0 0 0 4"/>
										<path d="M8.256 14a4.474 4.474 0 0 1-.229-1.004H3c.001-.246.154-.986.832-1.664C4.484 10.68 5.711 10 8 10c.26 0 .507.009.74.025.226-.341.496-.65.804-.918C9.077 9.038 8.564 9 8 9c-5 0-6 3-6 4s1 1 1 1z"/>
									</svg>
									{{ showAddMember ? "Chiudi ricerca" : "Aggiungi membro" }}
								</button>

								<!-- Campo di ricerca utenti -->
								<div v-if="showAddMember" class="mt-2">
									<input
										ref="memberSearchInput"
										v-model="searchQuery"
										type="text"
										class="form-control form-control-sm"
										placeholder="Cerca utente per nome…"
										autocomplete="off"
										@input="onSearchInput"
									/>

									<!-- Risultati ricerca -->
									<div
										v-if="searchResults.length > 0"
										class="list-group list-group-flush mt-2"
										style="max-height: 150px; overflow-y: auto"
									>
										<button
											v-for="user in searchResults"
											:key="user.userName"
											type="button"
											class="list-group-item list-group-item-action d-flex align-items-center gap-2 py-2"
											:disabled="addingMember"
											@click="addMember(user.userName)"
										>
											<img
												v-if="user.photo"
												:src="`data:image/jpeg;base64,${user.photo}`"
												class="rounded-circle flex-shrink-0"
												width="28"
												height="28"
												style="object-fit: cover"
											/>
											<div
												v-else
												class="rounded-circle bg-secondary d-flex align-items-center justify-content-center text-white flex-shrink-0"
												style="width: 28px; height: 28px; font-size: 0.75rem"
											>
												{{ user.userName.charAt(0).toUpperCase() }}
											</div>
											<span class="small">{{ user.userName }}</span>
											<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" fill="currentColor" class="bi ms-auto text-primary" viewBox="0 0 16 16">
												<path d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4"/>
											</svg>
										</button>
									</div>

									<!-- Stato ricerca -->
									<div v-if="searching" class="text-center mt-2">
										<span class="spinner-border spinner-border-sm text-secondary" role="status"></span>
									</div>
									<div v-if="searchError" class="text-danger small mt-1">{{ searchError }}</div>
									<div v-if="addError" class="text-danger small mt-1">{{ addError }}</div>
									<p
										v-if="!searching && searchQuery.trim().length >= 3 && searchResults.length === 0"
										class="text-secondary small mt-2 mb-0"
									>
										Nessun utente trovato.
									</p>
								</div>
							</div>

							<!-- Lista membri -->
							<div class="list-group list-group-flush" style="max-height: 320px; overflow-y: auto;">
								<div
									v-for="member in members"
									:key="member.userName"
									class="list-group-item d-flex align-items-center gap-3 py-2 px-0 border-0 border-bottom"
								>
									<!-- Avatar membro -->
									<img
										v-if="memberPhotoSrc(member)"
										:src="memberPhotoSrc(member)"
										:alt="member.userName"
										class="rounded-circle flex-shrink-0"
										width="36"
										height="36"
										style="object-fit: cover"
									/>
									<div
										v-else
										class="rounded-circle bg-secondary d-flex align-items-center justify-content-center text-white fw-semibold flex-shrink-0"
										style="width: 36px; height: 36px; font-size: 0.9rem"
									>
										{{ member.userName.charAt(0).toUpperCase() }}
									</div>

									<!-- Nome + badge admin -->
									<div class="d-flex flex-column lh-sm flex-grow-1">
										<span class="fw-medium small">{{ member.userName }}</span>
										<span
											v-if="member.userName === group.administrator"
											class="text-primary"
											style="font-size: 0.68rem"
										>
											Admin
										</span>
									</div>

									<!-- Pulsante kick (solo admin, non su se stesso) -->
									<button
										v-if="isAdmin && member.userName !== myUserName"
										type="button"
										class="btn btn-sm btn-outline-danger rounded-circle d-flex align-items-center justify-content-center"
										style="width: 28px; height: 28px; padding: 0;"
										title="Rimuovi dal gruppo"
										:disabled="kickingUser === member.userName"
										@click="kickMember(member.userName)"
									>
										<span
											v-if="kickingUser === member.userName"
											class="spinner-border spinner-border-sm"
											style="width: 12px; height: 12px"
											role="status"
										></span>
										<svg v-else xmlns="http://www.w3.org/2000/svg" width="12" height="12" fill="currentColor" viewBox="0 0 16 16">
											<path d="M4.646 4.646a.5.5 0 0 1 .708 0L8 7.293l2.646-2.647a.5.5 0 0 1 .708.708L8.707 8l2.647 2.646a.5.5 0 0 1-.708.708L8 8.707l-2.646 2.647a.5.5 0 0 1-.708-.708L7.293 8 4.646 5.354a.5.5 0 0 1 0-.708"/>
										</svg>
									</button>
								</div>
							</div>

							<!-- Lista vuota -->
							<p
								v-if="members.length === 0 && !loading"
								class="text-secondary small mt-2 mb-0"
							>
								Nessun membro nel gruppo.
							</p>
						</div>
					</div>
				</div>

				<!-- ── FOOTER ──────────────────────────────────────────────────── -->
				<div
					v-if="group"
					class="card-footer d-flex justify-content-center px-4 py-3 border-top bg-body-tertiary"
				>
					<button
						type="button"
						class="btn btn-outline-danger d-flex align-items-center gap-2"
						:disabled="leaving"
						@click="leaveOrDelete"
					>
						<span
							v-if="leaving"
							class="spinner-border spinner-border-sm"
							role="status"
						></span>
						<svg v-else xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" viewBox="0 0 16 16">
							<path fill-rule="evenodd" d="M10 12.5a.5.5 0 0 1-.5.5h-8a.5.5 0 0 1-.5-.5v-9a.5.5 0 0 1 .5-.5h8a.5.5 0 0 1 .5.5v2a.5.5 0 0 0 1 0v-2A1.5 1.5 0 0 0 9.5 2h-8A1.5 1.5 0 0 0 0 3.5v9A1.5 1.5 0 0 0 1.5 14h8a1.5 1.5 0 0 0 1.5-1.5v-2a.5.5 0 0 0-1 0z"/>
							<path fill-rule="evenodd" d="M15.854 8.354a.5.5 0 0 0 0-.708l-3-3a.5.5 0 0 0-.708.708L14.293 7.5H5.5a.5.5 0 0 0 0 1h8.793l-2.147 2.146a.5.5 0 0 0 .708.708z"/>
						</svg>
						{{ isAdmin ? "Elimina gruppo" : "Esci dal gruppo" }}
					</button>

					<!-- Errore uscita/eliminazione -->
					<div v-if="leaveError" class="text-danger small ms-3 align-self-center">{{ leaveError }}</div>
				</div>
			</div>
		</div>
	</Teleport>

	<!-- Modale errore di validazione -->
	<ModalDangerGeneric
		:visible="validationError.visible"
		:title="validationError.title"
		:description="validationError.description"
		@close="validationError.visible = false"
	/>
</template>

<style scoped>
/* Backdrop semitrasparente */
.modal-backdrop-custom {
	position: fixed;
	inset: 0;
	background: rgba(0, 0, 0, 0.5);
	z-index: 10000;
}

/* Contenitore centrato fisso sopra al backdrop */
.modal-container-custom {
	position: fixed;
	inset: 0;
	z-index: 10001;
	padding: 1rem;
}
</style>
