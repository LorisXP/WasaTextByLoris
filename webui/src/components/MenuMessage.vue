<script>
import auth from "../services/auth.js";
import axios from "../services/axios.js";
import ModalDangerGeneric from "./ModalDangerGeneric.vue";
import ModalForwardMessage from "./ModalForwardMessage.vue";

/** Emoji comuni per le reazioni rapide */
const COMMON_EMOJIS = [
	"👍", "😂", "😮", "😢", "😡",
	"🎉", "🔥", "👏", "🙏", "💯", "✅",
	"😍", "🤔", "😅", "🤣", "😎", "🥺",
	"😊", "🫡", "🫠", "💀", "🤌", "🫶",
];

export default {
	name: "MenuMessage",

	components: {
		ModalDangerGeneric,
		ModalForwardMessage,
	},

	props: {
		/** Controlla se il menu è visibile */
		visible: {
			type: Boolean,
			required: true,
		},
		/** Coordinata X (px) dove posizionare il menu */
		x: {
			type: Number,
			default: 0,
		},
		/** Coordinata Y (px) dove posizionare il menu */
		y: {
			type: Number,
			default: 0,
		},
		/** Oggetto messaggio su cui l'utente ha fatto tasto destro */
		msg: {
			type: Object,
			default: null,
		},
		/** ID della conversazione corrente */
		conversationID: {
			type: Number,
			required: true,
		},
		/**
		 * Tipo di conversazione:
		 * - "users"  → conversazione diretta tra utenti
		 * - "groups" → conversazione di gruppo
		 */
		conversationType: {
			type: String,
			default: "users",
			validator: (v) => ["users", "groups"].includes(v),
		},
		/**
		 * ID del gruppo (obbligatorio solo se conversationType === "groups").
		 */
		groupID: {
			type: Number,
			default: null,
		},
	},

	emits: ["close", "message-deleted", "message-forwarded", "reaction-added", "reaction-removed"],

	data() {
		return {
			showEmojiPicker: false,
			busy: false,
			actionError: null,
			emojiList: COMMON_EMOJIS,
			validationError: { visible: false, title: "Input non valido", description: "" },
			showForwardModal: false,
		};
	},

	computed: {
		/** ID dell'utente autenticato (dal servizio auth) */
		myUserID() {
			return auth.state.userID;
		},

		/** Stile CSS per il posizionamento fixed del menu */
		menuStyle() {
			return {
				position: "fixed",
				left: this.x + "px",
				top: this.y + "px",
				zIndex: 9999,
				minWidth: "170px",
			};
		},

		/** True se il messaggio corrente appartiene all'utente loggato */
		isMyMessage() {
			if (!this.msg || !this.msg.sender) return true;
			if (auth.state.userName) return this.msg.sender.userName === auth.state.userName;
			return false;
		},

		/**
		 * Restituisce il commento (reazione) dell'utente corrente su questo
		 * messaggio, oppure null se non ha reagito.
		 */
		myComment() {
			if (!this.msg?.comments || !auth.state.userName) return null;
			return this.msg.comments.find(
				(c) => c.sender?.userName === auth.state.userName
			) || null;
		},
	},

	methods: {
		//  CONTROLLO MENU 

		/** Chiude il menu e resetta lo stato interno */
		closeMenu() {
			this.showEmojiPicker = false;
			this.actionError = null;
			this.$emit("close");
		},

		/** Mostra la modale di errore di validazione */
		showValidationError(description) {
			this.validationError = { visible: true, title: "Input non valido", description };
		},

		//  AZIONI 

		/** Apre la modale di selezione conversazione per l'inoltro. */
		openForwardModal() {
			if (!this.msg) return;
			this.showForwardModal = true;
		},

		/** Chiamato quando ModalForwardMessage conferma l'inoltro con successo. */
		onForwarded() {
			this.showForwardModal = false;
			this.$emit("message-forwarded");
			this.closeMenu();
		},

		/**
		 * Elimina il messaggio dalla conversazione.
		 * Users:  DELETE /api/users/{userID}/conversations/users/{conversationID}/messages/{messageID}
		 * Groups: DELETE /api/users/{userID}/conversations/groups/{conversationID}/messages/{messageID}
		 */
		async deleteMessage() {
			if (!this.msg) return;
			this.busy = true;
			this.actionError = null;
			try {
				let response;
				if (this.conversationType === "groups") {
					response = await axios.delete(`/api/users/${this.myUserID}/conversations/groups/${this.conversationID}/messages/${this.msg.messageID}`);
				} else {
					response = await axios.delete(`/api/users/${this.myUserID}/conversations/users/${this.conversationID}/messages/${this.msg.messageID}`);
				}
				if (response.status === 200) {
					this.$emit("message-deleted");
					this.closeMenu();
				} else {
					this.actionError = "Errore durante la cancellazione del messaggio.";
				}
			} catch (e) {
				this.actionError = "Errore durante la cancellazione del messaggio.";
				console.error("Errore cancellazione messaggio:", e);
			} finally {
				this.busy = false;
			}
		},

		/** Apre o chiude il pannello emoji */
		toggleEmojiPicker() {
			this.showEmojiPicker = !this.showEmojiPicker;
		},

		/**
		 * Invia una reazione (emoji) al messaggio.
		 * Users:  POST /api/users/{userID}/conversations/users/{conversationID}/messages/{messageID}/comments
		 * Groups: POST /api/users/{userID}/conversations/groups/{conversationID}/messages/{messageID}/comments
		 *
		 * Il body richiede: { reaction: "<emoji>" }
		 */
		async sendReaction(emoji) {
			if (!this.msg) return;
			// Validazione emoji (reaction)
			const vEmoji = this.$validator(emoji, null, 1, 4, "string");
			if (!vEmoji.success) {
				this.showValidationError("La reazione deve essere un'emoji valida (1-4 caratteri).");
				return;
			}
			this.busy = true;
			this.actionError = null;
			try {
				let response;
				if (this.conversationType === "groups") {
					response = await axios.post(`/api/users/${this.myUserID}/conversations/groups/${this.conversationID}/messages/${this.msg.messageID}/comments`, { reaction: emoji });
				} else {
					response = await axios.post(`/api/users/${this.myUserID}/conversations/users/${this.conversationID}/messages/${this.msg.messageID}/comments`, { reaction: emoji });
				}
				if (response.status === 200) {
					this.$emit("reaction-added", emoji);
					this.closeMenu();
				} else {
					this.actionError = "Errore durante l'invio della reazione.";
				}
			} catch (e) {
				this.actionError = "Errore durante l'invio della reazione.";
				console.error("Errore reazione:", e);
			} finally {
				this.busy = false;
			}
		},

		/**
		 * Rimuove la propria reazione dal messaggio.
		 * Users:  DELETE /api/users/{userID}/conversations/users/{conversationID}/messages/{messageID}/comments/{commentID}
		 * Groups: DELETE /api/users/{userID}/conversations/groups/{conversationID}/messages/{messageID}/comments/{commentID}
		 */
		async deleteReaction() {
			const comment = this.myComment;
			if (!comment || !this.msg) return;
			this.busy = true;
			this.actionError = null;
			try {
				let response;
				if (this.conversationType === "groups") {
					response = await axios.delete(`/api/users/${this.myUserID}/conversations/groups/${this.conversationID}/messages/${this.msg.messageID}/comments/${comment.commentID}`);
				} else {
					response = await axios.delete(`/api/users/${this.myUserID}/conversations/users/${this.conversationID}/messages/${this.msg.messageID}/comments/${comment.commentID}`);
				}
				if (response.status === 204) {
					this.$emit("reaction-removed");
					this.closeMenu();
				} else {
					this.actionError = "Errore durante la rimozione della reazione.";
				}
			} catch (e) {
				this.actionError = "Errore durante la rimozione della reazione.";
				console.error("Errore rimozione reazione:", e);
			} finally {
				this.busy = false;
			}
		},
	},

	watch: {
		/** Resetta lo stato interno ogni volta che il menu viene (ri)aperto */
		visible(newVal) {
			if (newVal) {
				this.showEmojiPicker = false;
				this.actionError = null;
				this.showForwardModal = false;
			}
		},
	},
};
</script>

<template>
	<Teleport to="body">
		<!--
			Overlay trasparente: cattura il click al di fuori del menu
			per chiuderlo, senza bloccare la visualizzazione.
		-->
		<div
			v-if="visible"
			class="menu-message-overlay"
			@click.self="closeMenu"
			@contextmenu.prevent
		></div>

		<!-- Menu contestuale -->
		<div
			v-if="visible"
			:style="menuStyle"
			class="menu-message card shadow border rounded-3 overflow-hidden"
			role="menu"
			aria-label="Opzioni messaggio"
		>
			<!-- Errore azione (es. network error) -->
			<div
				v-if="actionError"
				class="px-3 py-2 text-danger small bg-danger-subtle border-bottom"
			>
				{{ actionError }}
			</div>

			<ul class="list-unstyled mb-0">
				<!--  Inoltra  -->
				<li role="none">
					<button
						type="button"
						class="menu-item d-flex align-items-center gap-2 w-100 px-3 py-2 btn btn-link text-body text-decoration-none"
						role="menuitem"
						:disabled="busy"
						@click="openForwardModal"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							width="15"
							height="15"
							fill="currentColor"
							class="bi bi-forward-fill flex-shrink-0"
							viewBox="0 0 16 16"
							aria-hidden="true"
						>
							<path
								d="M9.77 12.11 14.316 8 9.77 3.89a.5.5 0 0 0-.831.374v2.026H2.5a.5.5 0 0 0-.5.5v2.42a.5.5 0 0 0 .5.5h6.439v2.026a.5.5 0 0 0 .831.374z"
							/>
						</svg>
						<span>Inoltra</span>
					</button>
				</li>

				<template v-if="isMyMessage">
				<li class="menu-divider" role="separator"></li>

				<!--  Cancella  -->
				<li role="none">
					<button
						type="button"
						class="menu-item d-flex align-items-center gap-2 w-100 px-3 py-2 btn btn-link text-danger text-decoration-none"
						role="menuitem"
						:disabled="busy"
						@click="deleteMessage"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							width="15"
							height="15"
							fill="currentColor"
							class="bi bi-trash3-fill flex-shrink-0"
							viewBox="0 0 16 16"
							aria-hidden="true"
						>
							<path
								d="M11 1.5v1h3.5a.5.5 0 0 1 0 1h-.538l-.853 10.66A2 2 0 0 1 11.115 16h-6.23a2 2 0 0 1-1.994-1.84L1.038 3.5H.5a.5.5 0 0 1 0-1H4v-1A1.5 1.5 0 0 1 5.5 0h5A1.5 1.5 0 0 1 11 1.5m-5 0v1h4v-1a.5.5 0 0 0-.5-.5h-3a.5.5 0 0 0-.5.5M4.5 5.029l.5 8.5a.5.5 0 1 0 .998-.06l-.5-8.5a.5.5 0 1 0-.998.06m6.53-.528a.5.5 0 0 0-.528.47l-.5 8.5a.5.5 0 0 0 .998.058l.5-8.5a.5.5 0 0 0-.47-.528M8 4.5a.5.5 0 0 0-.5.5v8.5a.5.5 0 0 0 1 0V5a.5.5 0 0 0-.5-.5"
							/>
						</svg>
						<span>Cancella</span>
					</button>
				</li>
				</template>

				<!--  Rimuovi la mia reazione (solo se ho già reagito)  -->
				<template v-if="myComment">
				<li class="menu-divider" role="separator"></li>
				<li role="none">
					<button
						type="button"
						class="menu-item d-flex align-items-center gap-2 w-100 px-3 py-2 btn btn-link text-danger text-decoration-none"
						role="menuitem"
						:disabled="busy"
						@click="deleteReaction"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							width="15"
							height="15"
							fill="currentColor"
							class="bi bi-emoji-frown flex-shrink-0"
							viewBox="0 0 16 16"
							aria-hidden="true"
						>
							<path d="M8 15A7 7 0 1 1 8 1a7 7 0 0 1 0 14m0 1A8 8 0 1 0 8 0a8 8 0 0 0 0 16"/>
							<path d="M4.285 12.433a.5.5 0 0 0 .683-.183A3.498 3.498 0 0 1 8 10.5c1.295 0 2.426.703 3.032 1.75a.5.5 0 0 0 .866-.5A4.498 4.498 0 0 0 8 9.5a4.5 4.5 0 0 0-3.898 2.25.5.5 0 0 0 .183.683M7 6.5C7 7.328 6.552 8 6 8s-1-.672-1-1.5S5.448 5 6 5s1 .672 1 1.5m4 0c0 .828-.448 1.5-1 1.5s-1-.672-1-1.5S9.448 5 10 5s1 .672 1 1.5"/>
						</svg>
						<span>Rimuovi reazione {{ myComment.content }}</span>
					</button>
				</li>
				</template>

				<li class="menu-divider" role="separator"></li>

				<!--  Reagisci  -->
				<li role="none">
					<button
						type="button"
						class="menu-item d-flex align-items-center gap-2 w-100 px-3 py-2 btn btn-link text-body text-decoration-none"
						role="menuitem"
						:disabled="busy"
						:aria-expanded="showEmojiPicker"
						@click="toggleEmojiPicker"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							width="15"
							height="15"
							fill="currentColor"
							class="bi bi-emoji-smile flex-shrink-0"
							viewBox="0 0 16 16"
							aria-hidden="true"
						>
							<path
								d="M8 15A7 7 0 1 1 8 1a7 7 0 0 1 0 14m0 1A8 8 0 1 0 8 0a8 8 0 0 0 0 16"
							/>
							<path
								d="M4.285 9.567a.5.5 0 0 1 .683.183A3.498 3.498 0 0 0 8 11.5a3.498 3.498 0 0 0 3.032-1.75.5.5 0 1 1 .866.5A4.498 4.498 0 0 1 8 12.5a4.498 4.498 0 0 1-3.898-2.25.5.5 0 0 1 .183-.683M7 6.5C7 7.328 6.552 8 6 8s-1-.672-1-1.5S5.448 5 6 5s1 .672 1 1.5m4 0c0 .828-.448 1.5-1 1.5s-1-.672-1-1.5S9.448 5 10 5s1 .672 1 1.5"
							/>
						</svg>
						<span>Reagisci</span>
						<!-- Chevron che ruota quando il picker è aperto -->
						<svg
							xmlns="http://www.w3.org/2000/svg"
							width="11"
							height="11"
							fill="currentColor"
							class="bi bi-chevron-down ms-auto flex-shrink-0"
							:class="{ 'rotate-180': showEmojiPicker }"
							viewBox="0 0 16 16"
							aria-hidden="true"
							style="transition: transform 0.2s"
						>
							<path
								fill-rule="evenodd"
								d="M1.646 4.646a.5.5 0 0 1 .708 0L8 10.293l5.646-5.647a.5.5 0 0 1 .708.708l-6 6a.5.5 0 0 1-.708 0l-6-6a.5.5 0 0 1 0-.708z"
							/>
						</svg>
					</button>

					<!-- Griglia emoji (espandibile sotto la voce Reagisci) -->
					<div v-if="showEmojiPicker" class="emoji-grid px-3 pb-3 pt-1">
						<button
							v-for="emoji in emojiList"
							:key="emoji"
							type="button"
							class="emoji-btn btn btn-link p-0"
							:disabled="busy"
							:title="emoji"
							:aria-label="`Reagisci con ${emoji}`"
							@click="sendReaction(emoji)"
						>
							{{ emoji }}
						</button>
					</div>
				</li>
			</ul>

			<!-- Spinner mentre una richiesta è in corso -->
			<div v-if="busy" class="d-flex justify-content-center py-2 border-top">
				<span
					class="spinner-border spinner-border-sm text-primary"
					role="status"
					aria-label="Caricamento..."
				></span>
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

	<!-- Modale selezione conversazione per l'inoltro (componente dedicato con Teleport interno) -->
	<ModalForwardMessage
		:visible="showForwardModal"
		:msg="msg"
		:conversationID="conversationID"
		:conversationType="conversationType"
		@close="showForwardModal = false"
		@message-forwarded="onForwarded"
	/>
</template>

<style scoped>
/* Overlay trasparente che cattura i click fuori dal menu */
.menu-message-overlay {
	position: fixed;
	inset: 0;
	z-index: 9998;
}

/* Card del menu */
.menu-message {
	background: var(--bs-body-bg);
	border-color: var(--bs-border-color) !important;
}

/* Voce del menu */
.menu-item {
	border-radius: 0;
	font-size: 0.875rem;
	transition: background 0.12s;
}

.menu-item:hover:not(:disabled) {
	background: var(--bs-secondary-bg) !important;
}

.menu-item:disabled {
	opacity: 0.5;
	cursor: not-allowed;
}

/* Separatore sottile tra le voci */
.menu-divider {
	height: 1px;
	background: var(--bs-border-color);
	margin: 0 12px;
}

/* Griglia emoji 6 colonne */
.emoji-grid {
	display: grid;
	grid-template-columns: repeat(6, 1fr);
	gap: 2px;
}

/* Singolo pulsante emoji */
.emoji-btn {
	font-size: 1.2rem;
	line-height: 1.4;
	border-radius: 6px !important;
	transition: background 0.1s, transform 0.1s;
}

.emoji-btn:hover:not(:disabled) {
	background: var(--bs-secondary-bg) !important;
	transform: scale(1.18);
}

.emoji-btn:disabled {
	opacity: 0.5;
	cursor: not-allowed;
}

/* Rotazione del chevron quando il picker è aperto */
.rotate-180 {
	transform: rotate(180deg);
}
</style>
