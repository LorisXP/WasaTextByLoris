<script>
import auth from "../services/auth.js";
import axios from "../services/axios.js";

export default {
	name: "Sidebar",

	data: function () {
		return {
			conversations: [],
			errormsg: null,
			loading: false,
			searchQuery: "",

			/** Controlla la visibilità della modale di creazione gruppo */
			showGroupModal: false,
		};
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
				this.errormsg = e.toString();
				console.error("Errore nel caricamento delle conversazioni:", e);
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
		 * Handler stub per la barra di ricerca utenti.
		 */
		onSearchInput() {
			// TODO: implementare la ricerca utenti via API
			// GET /api/users/{userID}/other/{userName}
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
		 * Naviga alla conversazione selezionata.
		 */
		openConversation(conversation) {
			if (conversation.type === "group") {
				this.$router.push(
					`/conversations/groups/${conversation.conversationID}`
				);
			} else {
				this.$router.push(
					`/conversations/users/${conversation.conversationID}`
				);
			}
		},
	},

	mounted() {
		this.loadConversations();
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

		<!-- Stato di caricamento -->
		<div v-if="loading" class="text-center py-4 text-secondary">
			<LoadingSpinner />
		</div>

		<!-- Messaggio di errore -->
		<ErrorMsg v-if="errormsg" :msg="errormsg" class="mx-3 mt-2" />

		<!-- Lista delle conversazioni -->
		<div
			v-if="!loading"
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
							<small class="text-body-secondary flex-shrink-0 ms-2">
								{{
									conv.lastMessage
										? formatTimestamp(conv.lastMessage.timestamp)
										: ""
								}}
							</small>
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
</style>
