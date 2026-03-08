<script>
import auth from "../services/auth.js";
import axios from "../services/axios.js";
import ModalDangerGeneric from "./ModalDangerGeneric.vue";

export default {
	name: "StartConversation",

	components: {
		ModalDangerGeneric,
	},

	props: {
		/** Username dell'utente con cui iniziare la conversazione */
		userName: {
			type: String,
			required: true,
		},
		/** Foto dell'utente in base64 (opzionale) */
		userPhoto: {
			type: String,
			default: null,
		},
	},

	data() {
		return {
			loading: false,
			errormsg: null,
		};
	},

	computed: {
		/** Iniziale dell'username come fallback avatar */
		initial() {
			return this.userName ? this.userName.charAt(0).toUpperCase() : "?";
		},

		/** Src base64 della foto profilo */
		photoSrc() {
			return this.userPhoto
				? `data:image/jpeg;base64,${this.userPhoto}`
				: null;
		},
	},

	methods: {
		/**
		 * Crea la conversazione con l'utente.
		 * POST /api/users/{userID}/conversations/users
		 * Body: { userName }
		 * On success (201): naviga a /conversations/users/{conversationID}
		 */
		async startConversation() {
			this.errormsg = null;
			this.loading = true;
			try {
				const myUserID = auth.state.userID;
				const response = await axios.post(
					`/api/users/${myUserID}/conversations/users`,
					{ userName: this.userName }
				);
				if (response.status === 201) {
					const { conversationID } = response.data;
					this.$router.push({
						path: `/conversations/users/${conversationID}`,
						query: {
							userName: this.userName,
							...(this.userPhoto ? { photo: this.userPhoto } : {}),
						},
					});
				} else {
					this.errormsg = "Impossibile creare la conversazione.";
				}
			} catch (e) {
				this.errormsg = e.toString();
				console.error("Errore nella creazione della conversazione:", e);
			} finally {
				this.loading = false;
			}
		},
	},
};
</script>

<template>
	<div
		class="start-conversation d-flex flex-column align-items-center justify-content-center h-100 px-4"
	>
		<!-- Avatar circolare -->
		<div class="avatar-wrapper mb-4">
			<img
				v-if="photoSrc"
				:src="photoSrc"
				:alt="userName"
				class="rounded-circle shadow"
				width="120"
				height="120"
				style="object-fit: cover"
			/>
			<div
				v-else
				class="rounded-circle bg-primary d-flex align-items-center justify-content-center text-white fw-bold shadow"
				style="width: 120px; height: 120px; font-size: 2.5rem"
				:aria-label="userName"
			>
				{{ initial }}
			</div>
		</div>

		<!-- Username -->
		<h3 class="mb-1 fw-semibold text-center">{{ userName }}</h3>

		<!-- Label -->
		<p class="text-secondary mb-4 small">è su WASAText</p>

		<!-- Errore -->
		<ModalDangerGeneric :visible="!!errormsg" :description="errormsg || ''" @close="errormsg = null" />

		<!-- Pulsante -->
		<button
			type="button"
			class="btn btn-primary px-4 py-2 d-flex align-items-center gap-2"
			:disabled="loading"
			@click="startConversation"
			aria-label="Inizia una conversazione"
		>
			<span
				v-if="loading"
				class="spinner-border spinner-border-sm"
				role="status"
				aria-hidden="true"
			></span>
			<svg
				v-else
				xmlns="http://www.w3.org/2000/svg"
				width="16"
				height="16"
				fill="currentColor"
				class="bi bi-chat-dots-fill"
				viewBox="0 0 16 16"
			>
				<path
					d="M16 8c0 3.866-3.582 7-8 7a9 9 0 0 1-2.347-.306c-.584.296-1.925.864-4.181 1.234-.2.032-.352-.176-.273-.362.354-.836.674-1.95.77-2.966C.744 11.37 0 9.76 0 8c0-3.866 3.582-7 8-7s8 3.134 8 7M5 8a1 1 0 1 0-2 0 1 1 0 0 0 2 0m4 0a1 1 0 1 0-2 0 1 1 0 0 0 2 0m3 1a1 1 0 1 0 0-2 1 1 0 0 0 0 2"
				/>
			</svg>
			Inizia una conversazione
		</button>
	</div>
</template>

<style scoped>
.start-conversation {
	min-height: 100vh;
	background: var(--bs-body-bg);
}

.avatar-wrapper {
	transition: transform 0.2s;
}

.avatar-wrapper:hover {
	transform: scale(1.04);
}
</style>
