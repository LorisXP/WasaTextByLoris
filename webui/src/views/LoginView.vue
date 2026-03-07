<script>
import doLogin from "../services/login.js";

export default {
	name: "LoginView",

	data() {
		return {
			userName: "",
			errormsg: null,
			loading: false,
		};
	},

	methods: {
		async doAuth() {
			// Valida l'username prima di inviare la richiesta.
			const v = this.$validator(this.userName, /^[a-zA-Z0-9_]+$/, 3, 15, "string");
			if (!v.success) {
				this.errormsg = "Inserisci uno username valido (3-15 caratteri, solo lettere, numeri e _).";
				this.$nextTick(() => { try { this.$refs.username.focus(); } catch {} });
				return;
			}

			this.errormsg = null;
			this.loading = true;
			try {
				const err = await doLogin(this.userName);
				if (!err) {
					this.$router.push("/chats");
				} else {
					this.errormsg = err;
					this.$nextTick(() => { try { this.$refs.username.focus(); } catch {} });
				}
			} catch (e) {
				this.errormsg = e.toString();
				console.error("Login error:", e);
			} finally {
				this.loading = false;
			}
		},
	},
};
</script>

<template>
	<div class="login-view d-flex align-items-center justify-content-center min-vh-100 bg-body-tertiary">
		<main class="w-100 mx-auto px-3" style="max-width: 360px">
			<form class="text-center" @submit.prevent="doAuth" novalidate>
				<!-- Icona app -->
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="56"
					height="56"
					fill="currentColor"
					class="bi bi-chat-dots-fill text-primary mb-3"
					viewBox="0 0 16 16"
					aria-hidden="true"
				>
					<path
						d="M16 8c0 3.866-3.582 7-8 7a9 9 0 0 1-2.347-.306c-.584.296-1.925.864-4.181 1.234-.2.032-.352-.176-.273-.362.354-.836.674-1.95.77-2.966C.744 11.37 0 9.76 0 8c0-3.866 3.582-7 8-7s8 3.134 8 7M5 8a1 1 0 1 0-2 0 1 1 0 0 0 2 0m4 0a1 1 0 1 0-2 0 1 1 0 0 0 2 0m3 1a1 1 0 1 0 0-2 1 1 0 0 0 0 2"
					/>
				</svg>

				<h1 class="h4 mb-1 fw-semibold">WASAText</h1>
				<p class="text-secondary small mb-4">Accedi con il tuo username</p>

				<!-- Campo username -->
				<div class="form-floating mb-3">
					<input
						type="text"
						id="floatingInput"
						ref="username"
						v-model="userName"
						class="form-control"
						:class="{ 'is-invalid': errormsg }"
						placeholder="Username"
						autocomplete="username"
						:disabled="loading"
						aria-label="Username"
					/>
					<label for="floatingInput">Username</label>
					<div v-if="errormsg" class="invalid-feedback text-start">
						{{ errormsg }}
					</div>
				</div>

				<!-- Bottone accesso -->
				<button
					type="submit"
					class="btn btn-primary w-100 py-2 d-flex align-items-center justify-content-center gap-2"
					:disabled="loading"
				>
					<span
						v-if="loading"
						class="spinner-border spinner-border-sm"
						role="status"
						aria-hidden="true"
					></span>
					{{ loading ? "Accesso in corso…" : "Accedi" }}
				</button>
			</form>
		</main>
	</div>
</template>

<style scoped>
.login-view {
	background: var(--bs-body-bg);
}
</style>
