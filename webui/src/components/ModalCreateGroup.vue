<script>
import auth from "../services/auth.js";
import axios from "../services/axios.js";

/** Millisecondi di debounce sulla ricerca utenti */
const SEARCH_DEBOUNCE_MS = 420;

export default {
	name: "ModalCreateGroup",

	props: {
		/** Controlla se la modale è visibile */
		visible: {
			type: Boolean,
			required: true,
		},
	},

	emits: ["close", "group-created"],

	data() {
		return {
			//  Campi del form 
			groupName: "",
			groupPhoto: null,        // base64 pura (senza prefisso data URI)
			groupPhotoPreview: null, // data URL per anteprima

			//  Utenti selezionati 
			/** Utenti scelti da aggiungere al gruppo */
			selectedUsers: [],       // array di { userName }

			//  Suggerimenti da conversazioni esistenti 
			conversationSuggestions: [],  // array di { userName } derivati dalle conv.
			loadingSuggestions: false,

			//  Ricerca manuale 
			searchQuery: "",
			searchResults: [],       // array di { userName, photo? }
			searching: false,
			searchError: null,
			_searchTimer: null,

			//  Stato submit 
			submitting: false,
			submitError: null,
			//  Validazione 
			validationError: { visible: false, title: "Input non valido", description: "" },		};
	},

	computed: {
		/** ID dell'utente autenticato (usato come adminID del gruppo) */
		myUserID() {
			return auth.state.userID;
		},

		/** True se il form è in uno stato valido per essere inviato */
		canSubmit() {
			return (
				!this.submitting &&
				this.groupName.trim().length >= 3 &&
				this.groupName.trim().length <= 15
			);
		},

		/**
		 * Suggerimenti filtrati: mostra solo utenti non già selezionati e
		 * diversi dall'utente autenticato.
		 */
		filteredSuggestions() {
			const selectedNames = new Set(this.selectedUsers.map((u) => u.userName));
			return this.conversationSuggestions.filter(
				(u) => !selectedNames.has(u.userName)
			);
		},

		/**
		 * Risultati di ricerca filtrati: esclude gli utenti già selezionati
		 * e l'utente autenticato.
		 */
		filteredSearchResults() {
			const selectedNames = new Set(this.selectedUsers.map((u) => u.userName));
			return this.searchResults.filter((u) => !selectedNames.has(u.userName));
		},
	},

	methods: {
		//  LIFECYCLE DEL MODALE 

		/** Resetta tutto il form e chiude la modale */
		close() {
			this.resetForm();
			this.$emit("close");
		},

		/** Resetta tutti i campi del form allo stato iniziale */
		resetForm() {
			this.groupName = "";
			this.groupPhoto = null;
			this.groupPhotoPreview = null;
			this.selectedUsers = [];
			this.searchQuery = "";
			this.searchResults = [];
			this.searching = false;
			this.searchError = null;
			this.submitting = false;
			this.submitError = null;
			if (this._searchTimer) {
				clearTimeout(this._searchTimer);
				this._searchTimer = null;
			}
		},

		/** Mostra la modale di errore di validazione */
		showValidationError(description) {
			this.validationError = { visible: true, title: "Input non valido", description };
		},

		//  FOTO GRUPPO 

		/** Apre il file picker per la foto del gruppo */
		pickPhoto() {
			this.$refs.photoInput.click();
		},

		/** Gestisce la selezione di un'immagine per il gruppo */
		onPhotoSelected(event) {
			const file = event.target.files[0];
			if (!file) return;
			const reader = new FileReader();
			reader.onload = (e) => {
				this.groupPhotoPreview = e.target.result;
				// Rimuove il prefisso "data:...;base64," per inviare solo i dati raw
				this.groupPhoto = e.target.result.split(",")[1];
			};
			reader.readAsDataURL(file);
			// Reset dell'input per permettere di riselezionare lo stesso file
			event.target.value = "";
		},

		/** Rimuove la foto del gruppo */
		removePhoto() {
			this.groupPhoto = null;
			this.groupPhotoPreview = null;
		},

		//  SUGGERIMENTI DA CONVERSAZIONI 

		/**
		 * Carica i suggerimenti utenti dalle conversazioni esistenti.
		 * GET /api/users/{userID}/conversations
		 * Estrae gli interlocutori dalle conversazioni dirette (type !== "group").
		 */
		async loadSuggestions() {
			if (!this.myUserID) return;
			this.loadingSuggestions = true;
			try {
				const response = await axios.get(
					`/api/users/${this.myUserID}/conversations`
				);
				if (response.status === 200) {
					const conversations = response.data || [];
					// Estrae gli username univoci dalle conversazioni dirette
					const seen = new Set();
					const suggestions = [];
					for (const conv of conversations) {
						if (conv.type === "group") continue;
						const name = conv.name;
						if (name && !seen.has(name)) {
							seen.add(name);
							suggestions.push({ userName: name });
						}
					}
					this.conversationSuggestions = suggestions;
				}
			} catch (e) {
				console.error("Errore caricamento suggerimenti:", e);
			} finally {
				this.loadingSuggestions = false;
			}
		},

		//  RICERCA UTENTE 

		/**
		 * Chiamato ad ogni keystroke nella barra di ricerca.
		 * Implementa un debounce per evitare troppe chiamate API.
		 */
		onSearchInput() {
			if (this._searchTimer) clearTimeout(this._searchTimer);
			const q = this.searchQuery.trim();
			if (!q) {
				this.searchResults = [];
				this.searchError = null;
				this.searching = false;
				return;
			}
			// Mostra subito i placeholder skeleton
			this.searching = true;
			this.searchError = null;
			this._searchTimer = setTimeout(() => this.doSearch(q), SEARCH_DEBOUNCE_MS);
		},

		/**
		 * Esegue la ricerca utente.
		 * GET /api/users/{userID}/other/{userName}
		 */
		async doSearch(query) {
			const vQuery = this.$validator(query, /^[a-z]+[0-9]*$/, 3, 15, "string");
			if (!vQuery.success) {
				this.showValidationError("Il nome utente deve essere di 3-15 caratteri, lettere minuscole seguite da numeri opzionali (es. mario42).");
				this.searching = false;
				return;
			}
			try {
				const response = await axios.get(
					`/api/users/${this.myUserID}/other/${encodeURIComponent(query)}`
				);
				if (response.status === 200) {
					const data = response.data;
					// Il server può restituire un singolo utente o un array
					if (Array.isArray(data)) {
						this.searchResults = data;
					} else if (data && data.userName) {
						this.searchResults = [data];
					} else {
						this.searchResults = [];
					}
					this.searchError = null;
				} else if (response.status === 404) {
					this.searchResults = [];
					this.searchError = null;
				} else {
					this.searchResults = [];
				}
			} catch (e) {
				if (e.response && e.response.status === 404) {
					this.searchResults = [];
					this.searchError = null;
				} else {
					this.searchError = "Errore durante la ricerca.";
					console.error("Errore ricerca utente:", e);
				}
			} finally {
				this.searching = false;
			}
		},

		//  GESTIONE UTENTI SELEZIONATI 

		/**
		 * Aggiunge un utente alla lista dei selezionati (se non già presente).
		 * @param {{ userName: string }} user
		 */
		addUser(user) {
			const already = this.selectedUsers.some((u) => u.userName === user.userName);
			if (!already) {
				this.selectedUsers.push({ userName: user.userName });
			}
			// Pulisce la ricerca dopo la selezione
			this.searchQuery = "";
			this.searchResults = [];
		},

		/**
		 * Rimuove un utente dalla lista dei selezionati.
		 * @param {string} userName
		 */
		removeUser(userName) {
			this.selectedUsers = this.selectedUsers.filter(
				(u) => u.userName !== userName
			);
		},

		//  INVIO 

		/**
		 * Crea il gruppo e aggiunge gli utenti selezionati.
		 *
		 * Flusso:
		 *   1. POST /api/users/{userID}/groups → { groupID }   [201]
		 *   2. Se ci sono utenti selezionati:
		 *      POST /api/groups/{groupID}/users               [204]
		 */
		async submit() {
			if (!this.canSubmit) return;

			// Validazione nome gruppo
			const vName = this.$validator(this.groupName.trim(), null, 3, 15, "string");
			if (!vName.success) {
				this.showValidationError("Il nome del gruppo deve essere tra 3 e 15 caratteri.");
				return;
			}

			// Validazione foto (se presente)
			if (this.groupPhoto) {
				const vPhoto = this.$validator(this.groupPhoto, /^[A-Za-z0-9+/=]+$/, 0, 13981013, "string");
				if (!vPhoto.success) {
					this.showValidationError("La foto del gruppo non è valida o supera la dimensione massima consentita (≈ 10 MB).");
					return;
				}
			}

			this.submitting = true;
			this.submitError = null;

			try {
				//  Passo 1: crea il gruppo 
				const createPayload = {
					name: this.groupName.trim(),
					photo: this.groupPhoto || "",
				};

				const createResp = await axios.post(
					`/api/users/${this.myUserID}/groups`,
					createPayload
				);

				if (createResp.status !== 201) {
					this.submitError = "Errore durante la creazione del gruppo.";
					return;
				}

				const groupID = createResp.data.groupID;

				//  Passo 2: aggiungi utenti (se presenti) 
				if (this.selectedUsers.length > 0) {
					const userNames = this.selectedUsers.map((u) => u.userName);
					const addPayload = {
						userID: this.myUserID,
						userNames,
					};

					const addResp = await axios.post(
						`/api/groups/${groupID}/users`,
						addPayload
					);

					if (addResp.status !== 204) {
						// Il gruppo è stato creato ma l'aggiunta utenti ha fallito;
						// notifica però il successo parziale.
						this.submitError =
							"Gruppo creato, ma non è stato possibile aggiungere tutti gli utenti.";
						this.$emit("group-created", { groupID });
						this.close();
						return;
					}
				}

				//  Tutto ok 
				this.$emit("group-created", { groupID });
				this.close();

			} catch (e) {
				this.submitError =
					e.response?.data || "Errore imprevisto. Riprova.";
				console.error("Errore creazione gruppo:", e);
			} finally {
				this.submitting = false;
			}
		},
	},

	watch: {
		/** Quando la modale viene aperta, carica i suggerimenti */
		visible(newVal) {
			if (newVal) {
				this.resetForm();
				this.loadSuggestions();
			}
		},
	},
};
</script>

<template>
	<Teleport to="body">
		<!--
			Backdrop semitrasparente: clic fuori dalla modale → chiude.
			Coerente con il pattern già usato in MenuMessage.vue.
		-->
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
			aria-labelledby="modal-create-group-title"
		>
			<div
				class="card shadow-lg border rounded-4 overflow-hidden"
				style="width: 100%; max-width: 520px; max-height: 92vh; display: flex; flex-direction: column;"
			>
				<!--  HEADER  -->
				<div
					class="card-header d-flex align-items-center justify-content-between px-4 py-3 bg-primary text-white border-0"
				>
					<h5
						id="modal-create-group-title"
						class="mb-0 fw-semibold fs-6"
					>
						Crea nuovo gruppo
					</h5>
					<button
						type="button"
						class="btn-close btn-close-white"
						aria-label="Chiudi"
						:disabled="submitting"
						@click="close"
					></button>
				</div>

				<!--  BODY (scrollabile)  -->
				<div class="card-body overflow-y-auto px-4 py-4" style="overflow-y: auto; flex: 1 1 auto;">

					<!-- Errore globale del submit -->
					<div
						v-if="submitError"
						class="alert alert-danger py-2 small"
						role="alert"
					>
						{{ submitError }}
					</div>

					<!--  Nome del gruppo  -->
					<div class="mb-4">
						<label for="group-name" class="form-label fw-semibold">
							Nome del gruppo
							<span class="text-danger" aria-hidden="true">*</span>
						</label>
						<input
							id="group-name"
							v-model="groupName"
							type="text"
							class="form-control"
							placeholder="Es. Amici, Progetto…"
							minlength="3"
							maxlength="15"
							:disabled="submitting"
							autocomplete="off"
							aria-describedby="group-name-hint"
						/>
						<div id="group-name-hint" class="form-text">
							Da 3 a 15 caratteri.
							<span
								:class="groupName.trim().length > 15 ? 'text-danger' : 'text-secondary'"
							>{{ groupName.trim().length }}/15</span>
						</div>
					</div>

					<!--  Foto del gruppo (opzionale)  -->
					<div class="mb-4">
						<p class="form-label fw-semibold mb-2">Foto del gruppo <span class="fw-normal text-secondary">(opzionale)</span></p>
						<div class="d-flex align-items-center gap-3">
							<!-- Anteprima o avatar placeholder -->
							<div class="flex-shrink-0">
								<img
									v-if="groupPhotoPreview"
									:src="groupPhotoPreview"
									alt="Anteprima foto gruppo"
									class="rounded-circle border"
									width="64"
									height="64"
									style="object-fit: cover"
								/>
								<div
									v-else
									class="rounded-circle bg-success d-flex align-items-center justify-content-center text-white fw-bold"
									style="width: 64px; height: 64px; font-size: 1.4rem"
									aria-hidden="true"
								>
									{{ groupName.trim().charAt(0).toUpperCase() || "G" }}
								</div>
							</div>

							<!-- Pulsanti scegli/rimuovi foto -->
							<div class="d-flex flex-column gap-2">
								<button
									type="button"
									class="btn btn-outline-secondary btn-sm"
									:disabled="submitting"
									@click="pickPhoto"
								>
									<svg
										xmlns="http://www.w3.org/2000/svg"
										width="14"
										height="14"
										fill="currentColor"
										class="bi bi-camera me-1"
										viewBox="0 0 16 16"
										aria-hidden="true"
									>
										<path
											d="M15 12a1 1 0 0 1-1 1H2a1 1 0 0 1-1-1V6a1 1 0 0 1 1-1h1.172a3 3 0 0 0 2.12-.879l.83-.828A1 1 0 0 1 6.827 3h2.344a1 1 0 0 1 .707.293l.828.828A3 3 0 0 0 12.828 5H14a1 1 0 0 1 1 1zM2 4a2 2 0 0 0-2 2v6a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V6a2 2 0 0 0-2-2h-1.172a2 2 0 0 1-1.414-.586l-.828-.828A2 2 0 0 0 9.172 2H6.828a2 2 0 0 0-1.414.586l-.828.828A2 2 0 0 1 3.172 4z"
										/>
										<path
											d="M8 11a2.5 2.5 0 1 1 0-5 2.5 2.5 0 0 1 0 5m0 1a3.5 3.5 0 1 0 0-7 3.5 3.5 0 0 0 0 7M3 6.5a.5.5 0 1 1-1 0 .5.5 0 0 1 1 0"
										/>
									</svg>
									Scegli foto
								</button>
								<button
									v-if="groupPhotoPreview"
									type="button"
									class="btn btn-outline-danger btn-sm"
									:disabled="submitting"
									@click="removePhoto"
								>
									Rimuovi
								</button>
							</div>
						</div>
						<!-- Input file nascosto -->
						<input
							ref="photoInput"
							type="file"
							accept="image/jpeg,image/png,image/webp"
							class="d-none"
							@change="onPhotoSelected"
						/>
					</div>

					<!--  Aggiungi partecipanti  -->
					<div>
						<p class="form-label fw-semibold mb-2">
							Partecipanti
							<span class="fw-normal text-secondary">(opzionale)</span>
						</p>

						<!-- Pillole utenti selezionati -->
						<div
							v-if="selectedUsers.length > 0"
							class="d-flex flex-wrap gap-2 mb-3"
							aria-label="Utenti selezionati"
						>
							<span
								v-for="user in selectedUsers"
								:key="user.userName"
								class="badge rounded-pill bg-primary d-flex align-items-center gap-1 px-3 py-2"
								style="font-size: 0.8rem; font-weight: 500"
							>
								{{ user.userName }}
								<button
									type="button"
									class="btn-close btn-close-white ms-1"
									style="font-size: 0.55rem"
									:aria-label="`Rimuovi ${user.userName}`"
									:disabled="submitting"
									@click="removeUser(user.userName)"
								></button>
							</span>
						</div>

						<!-- Barra di ricerca utenti -->
						<div class="mb-2">
							<div class="input-group">
								<span class="input-group-text bg-body border-end-0">
									<!-- Spinner durante il fetch -->
									<span
										v-if="searching"
										class="spinner-border spinner-border-sm text-primary"
										role="status"
										aria-label="Ricerca in corso…"
									></span>
									<!-- Icona lente altrimenti -->
									<svg
										v-else
										xmlns="http://www.w3.org/2000/svg"
										width="14"
										height="14"
										fill="currentColor"
										class="bi bi-search text-secondary"
										viewBox="0 0 16 16"
										aria-hidden="true"
									>
										<path
											d="M11.742 10.344a6.5 6.5 0 1 0-1.397 1.398h-.001q.044.06.098.115l3.85 3.85a1 1 0 0 0 1.415-1.414l-3.85-3.85a1 1 0 0 0-.115-.1zM12 6.5a5.5 5.5 0 1 1-11 0 5.5 5.5 0 0 1 11 0"
										/>
									</svg>
								</span>
								<input
									v-model="searchQuery"
									type="search"
									class="form-control border-start-0 ps-0"
									placeholder="Cerca un utente per username…"
									:disabled="submitting"
									autocomplete="off"
									aria-label="Cerca utente"
									@input="onSearchInput"
								/>
							</div>
							<div v-if="searchError" class="form-text text-danger">{{ searchError }}</div>
						</div>

						<!-- Risultati ricerca (skeleton durante il fetch) -->
						<div v-if="searching" class="d-flex flex-column gap-2 mb-3">
							<div
								v-for="i in 3"
								:key="'skel-' + i"
								class="d-flex align-items-center gap-2 px-2 py-2 rounded-3"
								:style="{
									background: 'var(--bs-secondary-bg)',
									opacity: 0.7,
									animation: 'pulse 1.4s ease-in-out infinite',
								}"
							>
								<div
									class="rounded-circle"
									:style="{
										width: '32px', height: '32px',
										background: 'var(--bs-tertiary-bg)',
									}"
								></div>
								<div
									class="rounded-2"
									:style="{
										width: (55 + i * 22) + 'px',
										height: '12px',
										background: 'var(--bs-tertiary-bg)',
									}"
								></div>
							</div>
						</div>
						<div
							v-else-if="filteredSearchResults.length > 0"
							class="list-group mb-3"
						>
							<button
								v-for="user in filteredSearchResults"
								:key="user.userName"
								type="button"
								class="list-group-item list-group-item-action d-flex align-items-center gap-2 py-2"
								:disabled="submitting"
								@click="addUser(user)"
							>
								<!-- Avatar con iniziale -->
								<div
									class="rounded-circle bg-primary d-flex align-items-center justify-content-center text-white fw-semibold flex-shrink-0"
									style="width: 32px; height: 32px; font-size: 0.8rem"
									aria-hidden="true"
								>
									{{ user.userName.charAt(0).toUpperCase() }}
								</div>
								<span class="small">{{ user.userName }}</span>
								<svg
									xmlns="http://www.w3.org/2000/svg"
									width="14"
									height="14"
									fill="currentColor"
									class="bi bi-plus-circle ms-auto text-primary flex-shrink-0"
									viewBox="0 0 16 16"
									aria-hidden="true"
								>
									<path
										d="M8 15A7 7 0 1 1 8 1a7 7 0 0 1 0 14m0 1A8 8 0 1 0 8 0a8 8 0 0 0 0 16"
									/>
									<path
										d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4"
									/>
								</svg>
							</button>
						</div>

						<!-- Suggerimenti da conversazioni esistenti -->
						<div v-if="!searchQuery.trim()">
							<p class="text-body-secondary small mb-2">
								<svg
									xmlns="http://www.w3.org/2000/svg"
									width="12"
									height="12"
									fill="currentColor"
									class="bi bi-clock-history me-1"
									viewBox="0 0 16 16"
									aria-hidden="true"
								>
									<path
										d="M8.515 1.019A7 7 0 0 0 8 1V0a8 8 0 0 1 .589.022zm2.004.45a7 7 0 0 0-.985-.299l.219-.976q.576.129 1.126.342zm1.37.71a7 7 0 0 0-.439-.27l.493-.87a8 8 0 0 1 .979.654l-.615.789a7 7 0 0 0-.418-.302zm1.834 1.79a7 7 0 0 0-.653-.796l.724-.69q.406.429.747.91zm.744 1.352a7 7 0 0 0-.214-.468l.893-.45a8 8 0 0 1 .45 1.088l-.95.313a7 7 0 0 0-.179-.483m.53 2.507a6.991 6.991 0 0 0-.1-1.025l.985-.17q.1.58.116 1.17zm-.131 1.538q.05-.254.081-.51l.993.123a7.957 7.957 0 0 1-.23 1.155l-.964-.267q.069-.247.12-.501m-.952 2.379q.276-.436.486-.908l.914.405q-.24.54-.555 1.038zm-.964 1.205q.183-.183.35-.378l.758.653a8.073 8.073 0 0 1-.401.432z"
									/>
									<path d="M8 1a7 7 0 1 0 4.95 11.95l.707.707A8.001 8.001 0 1 1 8 0z" />
									<path
										d="M7.5 3a.5.5 0 0 1 .5.5v5.21l3.248 1.856a.5.5 0 0 1-.496.868l-3.5-2A.5.5 0 0 1 7 9V3.5a.5.5 0 0 1 .5-.5"
									/>
								</svg>
								Contatti recenti
							</p>

							<!-- Skeleton caricamento suggerimenti -->
							<div v-if="loadingSuggestions" class="d-flex flex-wrap gap-2 mb-1">
								<div
									v-for="i in 4"
									:key="'sugg-skel-' + i"
									class="rounded-pill"
									:style="{
										width: (60 + i * 12) + 'px',
										height: '30px',
										background: 'var(--bs-secondary-bg)',
										animation: 'pulse 1.4s ease-in-out infinite',
									}"
								></div>
							</div>

							<!-- Pillole suggerimento -->
							<div
								v-else-if="filteredSuggestions.length > 0"
								class="d-flex flex-wrap gap-2 mb-1"
							>
								<button
									v-for="user in filteredSuggestions"
									:key="user.userName"
									type="button"
									class="btn btn-sm btn-outline-secondary rounded-pill"
									:disabled="submitting"
									@click="addUser(user)"
								>
									+ {{ user.userName }}
								</button>
							</div>

							<!-- Nessun suggerimento -->
							<p
								v-else-if="!loadingSuggestions"
								class="text-secondary small fst-italic"
							>
								Nessun contatto recente. Usa la ricerca qui sopra.
							</p>
						</div>
					</div>
				</div>

				<!--  FOOTER  -->
				<div class="card-footer d-flex justify-content-end gap-2 px-4 py-3 border-top bg-body-tertiary">
					<button
						type="button"
						class="btn btn-outline-secondary"
						:disabled="submitting"
						@click="close"
					>
						Annulla
					</button>
					<button
						type="button"
						class="btn btn-primary d-flex align-items-center gap-2"
						:disabled="!canSubmit"
						@click="submit"
					>
						<span
							v-if="submitting"
							class="spinner-border spinner-border-sm"
							role="status"
							aria-hidden="true"
						></span>
						<svg
							v-else
							xmlns="http://www.w3.org/2000/svg"
							width="15"
							height="15"
							fill="currentColor"
							class="bi bi-people-fill"
							viewBox="0 0 16 16"
							aria-hidden="true"
						>
							<path
								d="M7 14s-1 0-1-1 1-4 5-4 5 3 5 4-1 1-1 1zm4-6a3 3 0 1 0 0-6 3 3 0 0 0 0 6m-5.784 6A2.238 2.238 0 0 1 5 13c0-1.355.68-2.75 1.936-3.72A6.325 6.325 0 0 0 5 9c-4 0-5 3-5 4s1 1 1 1zM4.5 8a2.5 2.5 0 1 0 0-5 2.5 2.5 0 0 0 0 5"
							/>
						</svg>
						Crea gruppo
					</button>
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

/* Animazione pulse per gli skeleton */
@keyframes pulse {
	0%, 100% { opacity: 0.6; }
	50%       { opacity: 0.25; }
}
</style>
