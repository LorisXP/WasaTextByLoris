<script>
import auth from "../services/auth.js";
import axios from "../services/axios.js";

export default {
	name: "ModalForwardMessage",

	props: {
		/** Controlla se la modale è visibile */
		visible: {
			type: Boolean,
			required: true,
		},
		/** Messaggio da inoltrare */
		msg: {
			type: Object,
			default: null,
		},
		/** ID della conversazione corrente (da escludere dalla lista) */
		conversationID: {
			type: Number,
			required: true,
		},
		/** "users" | "groups" — tipo della conversazione sorgente */
		conversationType: {
			type: String,
			default: "users",
		},
	},

	emits: ["close", "message-forwarded"],

	data() {
		return {
			conversations: [],
			loading: false,
			error: null,
			busy: false,
		};
	},

	computed: {
		myUserID() {
			return auth.state.userID;
		},
	},

	watch: {
		visible(newVal) {
			if (newVal) {
				this.fetchConversations();
			} else {
				this.conversations = [];
				this.error = null;
				this.busy = false;
			}
		},
	},

	methods: {
		/**
		 * Carica la lista conversazioni dell'utente ed esclude quella corrente.
		 * GET /api/users/{userID}/conversations
		 */
		async fetchConversations() {
			this.error = null;
			this.conversations = [];
			this.loading = true;
			try {
				const response = await axios.get(`/api/users/${this.myUserID}/conversations`);
				if (response.status === 200) {
					this.conversations = (response.data || []).filter(
						c => !(
							c.conversationID === this.conversationID &&
							((c.type === "direct" && this.conversationType === "users") ||
							 (c.type === "group"  && this.conversationType === "groups"))
						)
					);
				} else {
					this.error = "Impossibile caricare le conversazioni.";
				}
			} catch (e) {
				this.error = "Errore di rete nel caricare le conversazioni.";
				console.error("ModalForwardMessage: errore fetch conversazioni", e);
			} finally {
				this.loading = false;
			}
		},

		/**
		 * Invia la richiesta di inoltro verso la conversazione selezionata.
		 * POST /api/users/{userID}/conversations/{targetType}/{targetID}/messages/{messageID}
		 * Body: { source_type: "users"|"groups" }
		 */
		async confirm(targetConv) {
			if (!this.msg || !targetConv) return;
			this.busy = true;
			this.error = null;
			const targetType = targetConv.type === "group" ? "groups" : "users";
			try {
				const url = `/api/users/${this.myUserID}/conversations/${targetType}/${targetConv.conversationID}/messages/${this.msg.messageID}`;
				const response = await axios.post(url, { source_type: this.conversationType });
				if (response.status === 201) {
					this.$emit("message-forwarded");
					this.close();
				} else {
					this.error = "Errore durante l'inoltro del messaggio.";
				}
			} catch (e) {
				this.error = "Errore durante l'inoltro del messaggio.";
				console.error("ModalForwardMessage: errore inoltro", e);
			} finally {
				this.busy = false;
			}
		},

		close() {
			this.$emit("close");
		},
	},
};
</script>

<template>
	<Teleport to="body">
		<div
			v-if="visible"
			class="modal d-block modal-forward-backdrop"
			tabindex="-1"
			@click.self="close"
		>
			<div
				class="modal-dialog modal-dialog-centered modal-dialog-scrollable"
				style="max-width: 420px"
			>
				<div class="modal-content">
					<!-- Header -->
					<div class="modal-header py-2 px-3">
						<h6 class="modal-title fw-semibold mb-0">Inoltra a…</h6>
						<button
							type="button"
							class="btn-close btn-sm"
							aria-label="Chiudi"
							@click="close"
						></button>
					</div>

					<!-- Body -->
					<div class="modal-body p-0" style="max-height: 380px; overflow-y: auto">
						<!-- Caricamento -->
						<div v-if="loading" class="d-flex justify-content-center py-4">
							<span
								class="spinner-border spinner-border-sm text-primary"
								role="status"
								aria-label="Caricamento…"
							></span>
						</div>

						<!-- Errore -->
						<div v-else-if="error" class="px-3 py-3 text-danger small">
							{{ error }}
						</div>

						<!-- Lista vuota -->
						<div
							v-else-if="conversations.length === 0"
							class="text-center text-secondary small py-4"
						>
							Nessun'altra conversazione disponibile.
						</div>

						<!-- Lista conversazioni -->
						<ul v-else class="list-unstyled mb-0">
							<li
								v-for="conv in conversations"
								:key="conv.conversationID + '-' + conv.type"
							>
								<button
									type="button"
									class="forward-conv-item d-flex align-items-center gap-3 w-100 px-3 py-2 btn btn-link text-body text-decoration-none text-start"
									:disabled="busy"
									@click="confirm(conv)"
								>
									<!-- Avatar -->
									<img
										v-if="conv.photo"
										:src="'data:image/jpeg;base64,' + conv.photo"
										:alt="conv.name"
										class="rounded-circle flex-shrink-0"
										width="36"
										height="36"
										style="object-fit: cover"
									/>
									<div
										v-else
										class="rounded-circle flex-shrink-0 d-flex align-items-center justify-content-center text-white fw-semibold"
										:class="conv.type === 'group' ? 'bg-success' : 'bg-primary'"
										style="width: 36px; height: 36px; font-size: 0.9rem"
									>
										{{ conv.name ? conv.name.charAt(0).toUpperCase() : "?" }}
									</div>

									<!-- Info -->
									<div class="d-flex flex-column overflow-hidden">
										<span class="fw-semibold text-truncate" style="font-size: 0.875rem">
											{{ conv.name }}
										</span>
										<span class="text-secondary" style="font-size: 0.72rem">
											{{ conv.type === "group" ? "Gruppo" : "Utente" }}
										</span>
									</div>

									<!-- Spinner inoltro in corso -->
									<span
										v-if="busy"
										class="spinner-border spinner-border-sm text-primary ms-auto flex-shrink-0"
										role="status"
									></span>
								</button>
								<hr class="my-0 mx-3" style="opacity: 0.1" />
							</li>
						</ul>
					</div>
				</div>
			</div>
		</div>
	</Teleport>
</template>

<style scoped>
.modal-forward-backdrop {
	background: rgba(0, 0, 0, 0.45);
	z-index: 10500;
}

.forward-conv-item {
	border-radius: 0;
	font-size: 0.875rem;
	transition: background 0.12s;
}

.forward-conv-item:hover:not(:disabled) {
	background: var(--bs-secondary-bg) !important;
}

.forward-conv-item:disabled {
	opacity: 0.6;
	cursor: not-allowed;
}
</style>
