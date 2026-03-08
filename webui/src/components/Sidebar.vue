<script>
import auth from "../services/auth.js";
import axios from "../services/axios.js";
import LoadingSpinner from "./LoadingSpinner.vue";
import ModalDangerGeneric from "./ModalDangerGeneric.vue";

export default {
	name: "Sidebar",

	components: {
		LoadingSpinner,
		ModalDangerGeneric,
	},

	data: function () {
		return {
			conversations: [],
			errormsg: null,
			loading: false,
			searchQuery: "",

			/** Risultati della ricerca utenti */
			searchResults: [],
			searchLoading: false,
			searchError: null,
			/** Timer per il debounce della ricerca */
			searchDebounce: null,

			/** Controlla la visibilità della modale di creazione gruppo */
			showGroupModal: false,



			/**
			 * Mappa conversationID → true per le conversazioni con messaggi
			 * non letti (arrivati mentre la chat non era visualizzata).
			 */
			newMessageConvIds: {},
		}
	},

	computed: {
		/**
		 * Filtra le conversazioni in base alla query di ricerca (lato client).
		 */
		filteredConversations() {
			if (!this.searchQuery.trim()) {
				return this.conversations;
			}
			const q = this.searchQuery.trim().toLowerCase();
			return this.conversations.filter((c) =>
				c.name && c.name.toLowerCase().includes(q)
			);
		},
	},

	methods: {
		/**
		 * Recupera le conversazioni dell'utente autenticato.
		 * Implementa: GET /api/users/{userID}/conversations
		 */
		async loadConversations() {
			this.errormsg = null;
			this.loading = true;
			try {
				const userID = auth.state.userID;
				if (!userID) {
					this.errormsg = "Utente non autenticato";
					return;
				}
				const response = await axios.get(
					`/api/users/${userID}/conversations`
				);
				if (response.status === 200) {
					this.conversations = response.data;
				} else {
					this.errormsg = "Impossibile caricare le conversazioni";
				}
			} catch (e) {
				// 404 = nessuna conversazione (utente appena registrato): non è un errore.
				if (e.response && e.response.status === 404) {
					this.conversations = [];
				} else {
					this.errormsg = e.toString();
					console.error("Errore nel caricamento delle conversazioni:", e);
				}
			} finally {
				this.loading = false;
			}
		},

		/**
		 * Restituisce l'URL della foto di una conversazione come data URI base64,
		 * oppure null se non disponibile (mostra le iniziali come fallback).
		 */
		getPhotoSrc(conversation) {
			if (conversation.photo) {
				return `data:image/jpeg;base64,${conversation.photo}`;
			}
			return null;
		},

		/**
		 * Restituisce le iniziali del nome della conversazione da usare
		 * come avatar di fallback quando la foto non è disponibile.
		 */
		getInitials(name) {
			if (!name) return "?";
			return name.charAt(0).toUpperCase();
		},

		/**
		 * Formatta il timestamp dell'ultimo messaggio in modo leggibile.
		 */
		formatTimestamp(timestamp) {
			if (!timestamp) return "";
			const date = new Date(timestamp);
			const now = new Date();
			const isToday =
				date.getDate() === now.getDate() &&
				date.getMonth() === now.getMonth() &&
				date.getFullYear() === now.getFullYear();
			if (isToday) {
				return date.toLocaleTimeString("it-IT", {
					hour: "2-digit",
					minute: "2-digit",
				});
			}
			return date.toLocaleDateString("it-IT", {
				day: "2-digit",
				month: "2-digit",
			});
		},

		/**
		 * Restituisce un'anteprima testuale dell'ultimo messaggio.
		 */
		getLastMessagePreview(lastMessage) {
			if (!lastMessage) return "";
			if (lastMessage.messageType === "text") {
				return lastMessage.preview || "";
			}
			if (lastMessage.messageType === "photo") {
				return "📷 Foto";
			}
			return "";
		},

		/**
		 * Ricerca utenti con debounce.
		 * GET /api/users/{userID}/other/{userName}
		 */
		onSearchInput() {
			// Reset immediato dei risultati se la query è vuota
			if (!this.searchQuery.trim()) {
				this.searchResults = [];
				this.searchError = null;
				clearTimeout(this.searchDebounce);
				return;
			}

			// Attende almeno 3 caratteri prima di avviare la ricerca
			if (this.searchQuery.trim().length < 3) {
				this.searchResults = [];
				this.searchError = null;
				clearTimeout(this.searchDebounce);
				return;
			}

			// Debounce 400ms
			clearTimeout(this.searchDebounce);
			this.searchDebounce = setTimeout(() => {
				this.fetchUsers(this.searchQuery.trim());
			}, 400);
		},

		/**
		 * Chiama l'API di ricerca utenti.
		 */
		async fetchUsers(query) {
			this.searchLoading = true;
			this.searchError = null;
			try {
				const userID = auth.state.userID;
				const response = await axios.get(
					`/api/users/${userID}/other/${encodeURIComponent(query)}`
				);
				if (response.status === 200) {
					this.searchResults = response.data || [];
				} else if (response.status === 404) {
					this.searchResults = [];
				} else {
					this.searchError = "Errore nella ricerca";
				}
			} catch (e) {
				// 404 means no users found – not a real error
				if (e.response && e.response.status === 404) {
					this.searchResults = [];
				} else {
					this.searchError = "Errore nella ricerca utenti";
					console.error("Errore ricerca utenti:", e);
				}
			} finally {
				this.searchLoading = false;
			}
		},

		/**
		 * Naviga al profilo dell'utente trovato per avviare una conversazione.
		 */
		openUserProfile(user) {
			const params = { userName: user.userName };
			const query = user.photo ? { photo: user.photo } : {};
			this.$router.push({ path: `/start/${encodeURIComponent(user.userName)}`, query });
		},

		/**
		 * Apre la modale di creazione gruppo.
		 * POST /api/users/{userID}/groups
		 */
		onCreateGroup() {
			this.showGroupModal = true;
		},

		/**
		 * Chiamato quando la modale conferma la creazione del gruppo.
		 * Ricarica la lista delle conversazioni per mostrare il nuovo gruppo.
		 */
		onGroupCreated() {
			this.showGroupModal = false;
			this.loadConversations();
		},

		/**
		 * Naviga alla conversazione selezionata e azzera il suo badge.
		 */
		openConversation(conversation) {
			// Rimuove il badge di nuovo messaggio per questa conversazione.
			if (this.newMessageConvIds[conversation.conversationID]) {
				const updated = { ...this.newMessageConvIds };
				delete updated[conversation.conversationID];
				this.newMessageConvIds = updated;
			}
			if (conversation.type === "group") {
				this.$router.push({
					path: `/conversations/groups/${conversation.conversationID}`,
					query: {
						groupName: conversation.name || '',
						photo: conversation.photo || undefined,
						groupID: conversation.groupID || undefined,
					},
				});
			} else {
				this.$router.push({
					path: `/conversations/users/${conversation.conversationID}`,
					query: {
						userName: conversation.name || '',
						photo: conversation.photo || undefined,
					},
				});
			}
		},

		// ─ POLLING CONVERSAZIONI ─

		/** Avvia il polling silenzioso delle conversazioni ogni 3 secondi. */
		startPolling() {
			this.stopPolling();
			this.pollTimerID = setInterval(() => this.pollConversations(), 3000);
		},

		/** Ferma il polling. */
		stopPolling() {
			if (this.pollTimerID !== null) {
				clearInterval(this.pollTimerID);
				this.pollTimerID = null;
			}
		},

		/**
		 * Recupera silenziosamente la lista delle conversazioni ogni 3 s.
		 * Confronta con lo snapshot in localStorage per rilevare
		 * nuovi messaggi nelle conversazioni non attualmente aperte.
		 * Se ne trova, accende il badge "N" e riproduce un suono.
		 */
		async pollConversations() {
			const userID = auth.state.userID;
			if (!userID) return;
			try {
				const response = await axios.get(`/api/users/${userID}/conversations`);
				if (response.status !== 200) return;
				const freshConvs = response.data || [];

				// Legge lo snapshot precedente dal localStorage.
				const lsKey = `wasatext_convs_${userID}`;
				let prevState = {};
			try { prevState = JSON.parse(localStorage.getItem(lsKey) || '{}'); } catch (_e) { /* ignore */ }
				const isFirstLoad = Object.keys(prevState).length === 0;

				// ID della conversazione attualmente aperta nel pannello destro.
				const activeID = this.$route?.params?.conversationID
					? Number(this.$route.params.conversationID)
					: null;

				let hasNewAnywhere = false;
				const updatedBadges = { ...this.newMessageConvIds };

				if (!isFirstLoad) {
					for (const conv of freshConvs) {
						const prevMsgID = prevState[conv.conversationID] ?? null;
						const currMsgID = conv.lastMessage?.messageID ?? null;
						// Nuovo messaggio = ID cambiato, non è la conv attiva.
						if (currMsgID !== prevMsgID && conv.conversationID !== activeID) {
							updatedBadges[conv.conversationID] = true;
							hasNewAnywhere = true;
						}
					}
				}

				// Aggiorna le conversazioni e i badge.
				this.conversations = freshConvs;
				this.newMessageConvIds = updatedBadges;

				if (hasNewAnywhere) this.playNotificationSound();

				// Salva il nuovo snapshot.
				const newState = {};
				for (const conv of freshConvs) {
					newState[conv.conversationID] = conv.lastMessage?.messageID ?? null;
				}
				localStorage.setItem(lsKey, JSON.stringify(newState));
			} catch (e) {
				// 404 = nessuna conversazione (utente nuovo): aggiorna a lista vuota.
				if (e.response && e.response.status === 404) {
					this.conversations = [];
				}
				// Altri errori di rete durante il polling vengono ignorati silenziosamente.
			}
		},

		/**
		 * Riproduce un breve beep di notifica tramite Web Audio API.
		 * Non richiede file audio esterni.
		 */
		playNotificationSound() {
			try {
				const ctx = new (window.AudioContext || window.webkitAudioContext)();
				const osc = ctx.createOscillator();
				const gain = ctx.createGain();
				osc.connect(gain);
				gain.connect(ctx.destination);
				osc.type = 'sine';
				osc.frequency.setValueAtTime(880, ctx.currentTime);
				gain.gain.setValueAtTime(0.25, ctx.currentTime);
				gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.35);
				osc.start(ctx.currentTime);
				osc.stop(ctx.currentTime + 0.35);
			} catch (_e) {
				// Notifica sonora non disponibile (es. policy autoplay del browser).
			}
		},
	},

	created() {
		this.pollTimerID = null;
	},

	mounted() {
		this.loadConversations();
		this.startPolling();
	},

	unmounted() {
		this.stopPolling();
	},

	watch: {
		/**
		 * Quando la rotta cambia (es. l'utente naviga a una conversazione
		 * via $router.push):
		 * 1. Pulisce la ricerca nella sidebar;
		 * 2. Ricarica le conversazioni (utile dopo "Inizia conversazione");
		 * 3. Azzera il badge per la conversazione aperta.
		 */
		'$route'(to) {
			// 1. Rimuove la ricerca e i risultati, così la lista conversazioni torna visibile.
			if (this.searchQuery.trim()) {
				this.searchQuery = "";
				this.searchResults = [];
				this.searchError = null;
			}

			// 2. Ricarica le conversazioni per riflettere eventuali nuove conversazioni.
			this.loadConversations();

			// 3. Azzera il badge "nuovo messaggio" per la conversazione appena aperta.
			const id = to?.params?.conversationID ? Number(to.params.conversationID) : null;
			if (id && this.newMessageConvIds[id]) {
				const updated = { ...this.newMessageConvIds };
				delete updated[id];
				this.newMessageConvIds = updated;
			}
		},
	},
};
</script>

<template>
	<div
		class="sidebar d-flex flex-column align-items-stretch flex-shrink-0 bg-body-tertiary"
		style="width: 380px; height: 100vh; position: relative"
	>
		<!-- Header con titolo "Chat" -->
		<div
			class="sidebar-header d-flex align-items-center flex-shrink-0 p-3 link-body-emphasis text-decoration-none border-bottom"
		>
			<svg
				xmlns="http://www.w3.org/2000/svg"
				class="bi pe-none me-2"
				width="30"
				height="24"
				aria-hidden="true"
				viewBox="0 0 16 16"
				fill="currentColor"
			>
				<path
					d="M2.678 11.894a1 1 0 0 1 .287.801 11 11 0 0 1-.398 2c1.395-.323 2.247-.697 2.634-.893a1 1 0 0 1 .71-.074A8 8 0 0 0 8 14c3.996 0 7-2.807 7-6s-3.004-6-7-6-7 2.808-7 6c0 1.468.617 2.83 1.678 3.894m-.493 3.905a22 22 0 0 1-.713.129c-.2.032-.352-.176-.273-.362a10 10 0 0 0 .244-.637l.003-.01c.248-.72.45-1.548.524-2.319C.743 11.37 0 9.76 0 8c0-3.866 3.582-7 8-7s8 3.134 8 7-3.582 7-8 7a9 9 0 0 1-2.347-.306c-.52.263-1.639.742-3.468 1.105"
				/>
			</svg>
			<span class="fs-5 fw-semibold">Chat</span>
		</div>

		<!-- Barra di ricerca utenti -->
		<div class="px-3 pt-2 pb-1 border-bottom">
			<div class="input-group">
				<span class="input-group-text bg-body border-end-0">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						width="14"
						height="14"
						fill="currentColor"
						class="bi bi-search text-secondary"
						viewBox="0 0 16 16"
					>
						<path
							d="M11.742 10.344a6.5 6.5 0 1 0-1.397 1.398h-.001q.044.06.098.115l3.85 3.85a1 1 0 0 0 1.415-1.414l-3.85-3.85a1 1 0 0 0-.115-.1zM12 6.5a5.5 5.5 0 1 1-11 0 5.5 5.5 0 0 1 11 0"
						/>
					</svg>
				</span>
				<input
					type="search"
					class="form-control border-start-0 ps-0"
					placeholder="Cerca utenti…"
					v-model="searchQuery"
					@input="onSearchInput"
					aria-label="Cerca utenti"
				/>
			</div>
		</div>

		<!-- Risultati ricerca utenti (visibili solo quando si sta cercando) -->
		<div
			v-if="searchQuery.trim()"
			class="search-results border-bottom"
			style="overflow-y: auto; max-height: 40vh"
		>
			<!-- Caricamento ricerca -->
			<div v-if="searchLoading" class="text-center py-3 text-secondary small">
				<span class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
				Ricerca in corso…
			</div>

			<!-- Errore ricerca -->
			<div v-else-if="searchError" class="text-danger small px-3 py-2">{{ searchError }}</div>

			<!-- Nessun risultato -->
			<div
				v-else-if="searchResults.length === 0 && !searchLoading"
				class="text-secondary small px-3 py-3 text-center"
			>
				Nessun utente trovato
			</div>

			<!-- Lista utenti trovati -->
			<a
				v-for="user in searchResults"
				:key="user.userName"
				href="#"
				class="list-group-item list-group-item-action py-2 px-3 d-flex align-items-center gap-3"
				@click.prevent="openUserProfile(user)"
			>
				<!-- Avatar utente trovato -->
				<img
					v-if="user.photo"
					:src="'data:image/jpeg;base64,' + user.photo"
					:alt="user.userName"
					class="rounded-circle flex-shrink-0"
					width="38"
					height="38"
					style="object-fit: cover"
				/>
				<div
					v-else
					class="rounded-circle bg-primary d-flex align-items-center justify-content-center text-white fw-semibold flex-shrink-0"
					style="width: 38px; height: 38px; font-size: 1rem"
				>
					{{ user.userName.charAt(0).toUpperCase() }}
				</div>
				<!-- Username -->
				<span class="fw-medium text-truncate">{{ user.userName }}</span>
			</a>
		</div>

		<!-- Stato di caricamento -->
		<div v-if="loading && !searchQuery.trim()" class="text-center py-4 text-secondary">
			<LoadingSpinner />
		</div>

		<!-- Messaggio di errore -->
		<ModalDangerGeneric :visible="!!(errormsg && !searchQuery.trim())" :description="errormsg || ''" @close="errormsg = null" />

		<!-- Lista delle conversazioni -->
		<div
			v-if="!loading && !searchQuery.trim()"
			class="list-group list-group-flush border-bottom scrollarea flex-grow-1"
			style="overflow-y: auto"
		>
			<!-- Stato vuoto -->
			<div
				v-if="filteredConversations.length === 0 && !errormsg"
				class="text-center text-secondary py-5 px-3"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="40"
					height="40"
					fill="currentColor"
					class="bi bi-chat-dots mb-2 text-secondary opacity-50"
					viewBox="0 0 16 16"
				>
					<path
						d="M5 8a1 1 0 1 1-2 0 1 1 0 0 1 2 0m4 0a1 1 0 1 1-2 0 1 1 0 0 1 2 0m3 1a1 1 0 1 0 0-2 1 1 0 0 0 0 2"
					/>
					<path
						d="m2.165 15.803.02-.004c1.83-.363 2.948-.842 3.468-1.105A9 9 0 0 0 8 15c4.418 0 8-3.134 8-7s-3.582-7-8-7-8 3.134-8 7c0 1.76.743 3.37 1.97 4.6a10.4 10.4 0 0 1-.524 2.318l-.003.011a11 11 0 0 1-.244.637c-.079.186.074.394.273.362a22 22 0 0 0 .693-.125m.8-3.108a1 1 0 0 0-.287-.801C1.618 10.83 1 9.468 1 8c0-3.192 3.004-6 7-6s7 2.808 7 6-3.004 6-7 6a8 8 0 0 1-2.088-.272 1 1 0 0 0-.711.074c-.387.196-1.24.57-2.634.893a11 11 0 0 0 .398-2"
					/>
				</svg>
				<p class="mb-0 small">Nessuna conversazione</p>
			</div>

			<!-- Elemento conversazione -->
			<a
				v-for="conv in filteredConversations"
				:key="conv.conversationID"
				href="#"
				class="list-group-item list-group-item-action py-3 lh-sm"
				:class="{ 'conv-unread': newMessageConvIds[conv.conversationID] }"
				@click.prevent="openConversation(conv)"
			>
				<div class="d-flex align-items-center gap-3">
					<!-- Avatar / Foto profilo cerchiata -->
					<div class="flex-shrink-0">
						<img
							v-if="getPhotoSrc(conv)"
							:src="getPhotoSrc(conv)"
							:alt="conv.name"
							class="rounded-circle object-fit-cover"
							width="46"
							height="46"
							style="object-fit: cover"
						/>
						<!-- Fallback con iniziali se la foto non è disponibile -->
						<div
							v-else
							class="rounded-circle d-flex align-items-center justify-content-center fw-semibold text-white"
							:class="conv.type === 'group' ? 'bg-success' : 'bg-primary'"
							style="width: 46px; height: 46px; font-size: 1.1rem"
							:aria-label="conv.name"
						>
							{{ getInitials(conv.name) }}
						</div>
					</div>

					<!-- Contenuto testuale della conversazione -->
					<div class="flex-grow-1 overflow-hidden">
						<div
							class="d-flex w-100 align-items-center justify-content-between"
						>
							<strong class="mb-1 text-truncate">{{ conv.name }}</strong>
							<div class="d-flex align-items-center gap-1 flex-shrink-0 ms-2">
								<!-- Badge "N" per nuovo messaggio non letto -->
								<span
									v-if="newMessageConvIds[conv.conversationID]"
									class="badge bg-success rounded-pill"
									style="font-size:0.6rem; padding: 3px 6px"
									aria-label="Nuovo messaggio"
								>N</span>
								<small class="text-body-secondary">
									{{
										conv.lastMessage
											? formatTimestamp(conv.lastMessage.timestamp)
											: ""
									}}
								</small>
							</div>
						</div>
						<div class="col-12 mb-1 small text-body-secondary text-truncate">
							{{ getLastMessagePreview(conv.lastMessage) }}
						</div>
					</div>
				</div>
			</a>
		</div>

		<!-- FAB pulsante circolare "+" per la creazione di un gruppo (Material Design) -->
		<button
			type="button"
			class="btn btn-primary rounded-circle shadow-lg d-flex align-items-center justify-content-center"
			style="
				position: absolute;
				bottom: 24px;
				right: 24px;
				width: 56px;
				height: 56px;
				font-size: 1.6rem;
				z-index: 100;
				line-height: 1;
				padding: 0;
				transition: box-shadow 0.2s, transform 0.2s;
			"
			@click="onCreateGroup"
			title="Crea nuovo gruppo"
			aria-label="Crea nuovo gruppo"
		>
			+
		</button>
	</div>

	<!--
		Modale di creazione gruppo.
		Usa <Teleport to="body"> internamente, quindi il rendering avviene
		fuori dalla sidebar (nessun overflow:hidden che potrebbe tagliarlo).
	-->
	<ModalCreateGroup
		:visible="showGroupModal"
		@close="showGroupModal = false"
		@group-created="onGroupCreated"
	/>
</template>

<style scoped>
.sidebar {
	min-width: 320px;
	max-width: 380px;
}

.scrollarea {
	overflow-y: auto;
}

/* Hover elevation sul FAB */
button.btn.rounded-circle:hover {
	transform: scale(1.08);
	box-shadow: 0 6px 20px rgba(0, 0, 0, 0.25) !important;
}

/* Animazione sottile sull'elemento della lista al hover */
.list-group-item-action {
	transition: background-color 0.15s ease;
}

/* Evidenziazione conversazione con messaggio non letto */
.conv-unread {
	background-color: color-mix(in srgb, var(--bs-success) 12%, transparent) !important;
	border-left: 3px solid var(--bs-success);
}
.conv-unread:hover {
	background-color: color-mix(in srgb, var(--bs-success) 20%, transparent) !important;
}
</style>
