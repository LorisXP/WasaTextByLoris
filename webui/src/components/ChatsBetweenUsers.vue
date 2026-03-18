<script>
	import auth from "../services/auth.js";
	import axios from "../services/axios.js";
	import ModalDangerGeneric from "./ModalDangerGeneric.vue";

	export default {
		name: "ChatsBetweenUsers",

		components: {
			ModalDangerGeneric,
		},

		props: {
			/** ID della conversazione (int), passato dal componente padre */
			conversationID: {
				type: Number,
				required: true,
			},
			/** Nome utente dell'interlocutore */
			otherUserName: {
				type: String,
				default: "",
			},
			/** Foto dell'interlocutore (base64), opzionale */
			otherUserPhoto: {
				type: String,
				default: null,
			},
		},

		data() {
			return {
				messages: [],
				loading: false,
				errormsg: null,
				sendError: null,
				sending: false,

					// Composizione nuovo messaggio
				newMessageText: "",
				attachedFile: null,       // File object scelto dall'utente
				attachedFileType: null,   // "photo" | "gif"
				attachedFilePreview: null,// data URL per anteprima

				// Menu contestuale
				menuVisible: false,
				menuX: 0,
				menuY: 0,
				menuMsg: null,

					// Validazione input
				validationError: { visible: false, title: "Input non valido", description: "" },
			};
		},

		computed: {
			/** ID dell'utente autenticato */
			myUserID() {
				return auth.state.userID;
			},

			/** Iniziale dell'interlocutore per avatar fallback */
			otherInitial() {
				return this.otherUserName ? this.otherUserName.charAt(0).toUpperCase() : "?";
			},

			/** Src base64 della foto dell'interlocutore */
			otherPhotoSrc() {
				return this.otherUserPhoto
					? `data:image/jpeg;base64,${this.otherUserPhoto}`
					: null;
			},

			/** True se il bottone Invia è disabilitato */
			cannotSend() {
				return this.sending || (!this.newMessageText.trim() && !this.attachedFile);
			},

			/**
			 * Restituisce la lista dei messaggi arricchita con separatori di data.
			 * Ogni elemento e { type: 'separator', label } oppure { type: 'message', data: msg }.
			 */
			messagesWithDateSeparators() {
				const result = [];
				let lastDateKey = null;
				for (const msg of this.messages) {
					const dateKey = msg.timestamp ? msg.timestamp.slice(0, 10) : null;
					if (dateKey && dateKey !== lastDateKey) {
						result.push({ type: 'separator', label: this.formatDateLabel(msg.timestamp) });
						lastDateKey = dateKey;
					}
					result.push({ type: 'message', data: msg });
				}
				return result;
			},
		},

		methods: {
			//  LETTURA MESSAGGI 

			/**
			 * Carica i messaggi della conversazione.
			 * GET /api/users/{userID}/conversations/users/{conversationID}/messages
			 */
			async loadMessages() {
				this.errormsg = null;
				this.loading = true;
				try {
					const response = await axios.get(
						`/api/users/${this.myUserID}/conversations/users/${this.conversationID}/messages`
					);
					if (response.status === 200) {
						this.messages = response.data.messages || [];
						this.$nextTick(() => this.scrollToBottom());
					} else {
						this.errormsg = "Impossibile caricare i messaggi";
					}
				} catch (e) {
					this.errormsg = e.toString();
					console.error("Errore caricamento messaggi:", e);
				} finally {
					this.loading = false;
				}
			},

			//  INVIO MESSAGGIO 

			/**
			 * Invia un messaggio di testo, foto o GIF.
			 * POST /api/users/{userID}/conversations/users/{conversationID}/messages
			 */
			async sendMessage() {
				if (this.cannotSend) return;
				this.sendError = null;
				this.sending = true;

				try {
					let content;
					let content_type;

					if (this.attachedFile) {
						// Allegato: leggi come base64
						content = await this.fileToBase64(this.attachedFile);
						content_type = this.attachedFileType; // "photo" | "gif"
					} else {
						content = this.newMessageText.trim();
						content_type = "text";
					}

				// Validazione contenuto messaggio
				const isMedia = content_type === "photo" || content_type === "gif";
				const vContent = isMedia
					? this.$validator(content, /^[A-Za-z0-9+/=]+$/, 0, 13981013, "string")
					: this.$validator(content, null, 1, 4095, "string");
				if (!vContent.success) {
					this.showValidationError(
						isMedia
							? "Il file allegato non è valido o supera la dimensione massima consentita (≈ 10 MB)."
							: "Il messaggio deve contenere tra 1 e 4095 caratteri."
					);
					return;
				}

					const response = await axios.post(
						`/api/users/${this.myUserID}/conversations/users/${this.conversationID}/messages`,
						{ content, content_type }
					);

					if (response.status === 201) {
						// Reset del campo di input
						this.newMessageText = "";
						this.clearAttachment();
						// Ricarica i messaggi per mostrare quello appena inviato
						await this.loadMessages();
					} else {
						this.sendError = "Errore nell'invio del messaggio";
					}
				} catch (e) {
					this.sendError = e.toString();
					console.error("Errore invio messaggio:", e);
				} finally {
					this.sending = false;
				}
			},

			//  ALLEGATI 

			/** Apre il file picker filtrando per immagini (photo) */
			pickPhoto() {
				this.$refs.fileInputPhoto.click();
			},

			/** Apre il file picker filtrando per GIF */
			pickGif() {
				this.$refs.fileInputGif.click();
			},

			/** Gestisce la selezione di un file immagine (photo) */
			onPhotoSelected(event) {
				const file = event.target.files[0];
				if (!file) return;
				this.setAttachment(file, "photo");
				// Reset input per poter riselezionare lo stesso file
				event.target.value = "";
			},

			/** Gestisce la selezione di un file GIF */
			onGifSelected(event) {
				const file = event.target.files[0];
				if (!file) return;
				this.setAttachment(file, "gif");
				event.target.value = "";
			},

			/** Imposta l'allegato corrente e genera l'anteprima */
			setAttachment(file, type) {
				this.attachedFile = file;
				this.attachedFileType = type;
				this.newMessageText = ""; // allegato esclusivo rispetto al testo
				const reader = new FileReader();
				reader.onload = (e) => { this.attachedFilePreview = e.target.result; };
				reader.readAsDataURL(file);
			},

			/** Rimuove l'allegato corrente */
			clearAttachment() {
				this.attachedFile = null;
				this.attachedFileType = null;
				this.attachedFilePreview = null;
			},

			/** Mostra la modale di errore di validazione */
			showValidationError(description) {
				this.validationError = { visible: true, title: "Input non valido", description };
			},

			//  UTILITY 

			/** Converte un File in stringa base64 */
			fileToBase64(file) {
				return new Promise((resolve, reject) => {
					const reader = new FileReader();
					reader.onload = (e) => {
						// Rimuove il prefisso "data:...;base64," per restituire solo i dati
						const base64 = e.target.result.split(",")[1];
						resolve(base64);
					};
					reader.onerror = reject;
					reader.readAsDataURL(file);
				});
			},

			/** Scorre la lista messaggi fino all'ultimo */
			scrollToBottom() {
				const el = this.$refs.messagesList;
				if (el) {
					el.scrollTop = el.scrollHeight;
					// Secondo scroll dopo breve ritardo per immagini caricate in ritardo
					setTimeout(() => { el.scrollTop = el.scrollHeight; }, 150);
				}
			},

			/** Apre un'immagine base64 in una nuova scheda del browser (Blob URL) */
			openImageInNewTab(base64, type) {
				const mimeType = type === 'gif' ? 'image/gif' : 'image/jpeg';
				const byteChars = atob(base64);
				const byteArray = new Uint8Array(byteChars.length);
				for (let i = 0; i < byteChars.length; i++) {
					byteArray[i] = byteChars.charCodeAt(i);
				}
				const blob = new Blob([byteArray], { type: mimeType });
				const url = URL.createObjectURL(blob);
				window.open(url, '_blank');
			},

			/**
			 * Restituisce true se il messaggio appartiene all'utente autenticato,
			 * confrontando il sender.userName con lo userName dell'utente loggato.
			 */
			isSentByMe(msg) {
				if (!msg.sender) return false;
				return msg.sender.userName === auth.state.userName;
			},

			/** Formatta timestamp in ora HH:MM */
			formatTime(timestamp) {
				if (!timestamp) return "";
				return new Date(timestamp).toLocaleTimeString("it-IT", {
					hour: "2-digit",
					minute: "2-digit",
				});
			},

			/**
			 * Restituisce l'etichetta della data per il separatore:
			 * "Oggi", "Ieri", oppure la data completa in italiano
			 * (es. "lunedi 27 settembre 2025").
			 */
			formatDateLabel(timestamp) {
				if (!timestamp) return "";
				const date = new Date(timestamp);
				const now = new Date();
				const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
				const msgDay = new Date(date.getFullYear(), date.getMonth(), date.getDate());
				const diffDays = Math.round((today - msgDay) / (1000 * 60 * 60 * 24));
				if (diffDays === 0) return "Oggi";
				if (diffDays === 1) return "Ieri";
				return date.toLocaleDateString("it-IT", {
					weekday: "long",
					day: "numeric",
					month: "long",
					year: "numeric",
				});
			},

			/**
			 * Raggruppa le reazioni per emoji e restituisce un array
			 * { emoji, count } ordinato per frequenza decrescente.
			 */
			groupedReactions(comments) {
				if (!comments || comments.length === 0) return [];
				const map = {};
				for (const c of comments) {
					const emoji = c.content;
					map[emoji] = (map[emoji] || 0) + 1;
				}
				return Object.entries(map)
					.map(([emoji, count]) => ({ emoji, count }))
					.sort((a, b) => b.count - a.count);
			},

			/** Gestisce invio con tasto Enter (Shift+Enter per a capo) */
			onTextKeydown(event) {
				if (event.key === "Enter" && !event.shiftKey) {
					event.preventDefault();
					this.sendMessage();
				}
			},

		//  MENU CONTESTUALE 

		/**
		 * Apre il menu contestuale sul messaggio.
		 * Calcola la posizione in modo che il menu resti entro i bordi del viewport.
		 */
		openMenu(event, msg) {
			this.menuMsg = msg;
			// Stima altezza menu: ~160px base + ~130px se emoji aperte
			const estimatedHeight = 165;
			const estimatedWidth = 180;
			const x = Math.min(event.clientX, window.innerWidth - estimatedWidth - 8);
			const y = Math.min(event.clientY, window.innerHeight - estimatedHeight - 8);
			this.menuX = x;
			this.menuY = y;
			this.menuVisible = true;
		},

		/** Chiude il menu contestuale */
		onMenuClose() {
			this.menuVisible = false;
			this.menuMsg = null;
		},

		/** Ricarica i messaggi dopo la cancellazione tramite menu */
		onMessageDeleted() {
			this.menuVisible = false;
			this.menuMsg = null;
			this.loadMessages();
		},

		/** Ricarica i messaggi dopo l'inoltro tramite menu */
		onMessageForwarded() {
			this.menuVisible = false;
			this.menuMsg = null;
			this.loadMessages();
		},

		/** Ricarica i messaggi dopo l'aggiunta di una reazione tramite menu */
		onReactionAdded() {
			this.menuVisible = false;
			this.menuMsg = null;
			this.loadMessages();
		},

		//  POLLING 

		/** Avvia il polling silenzioso dei messaggi ogni 3 secondi. */
		startPolling() {
			this.stopPolling();
			this.pollTimerID = setInterval(() => this.pollMessages(), 3000);
		},

		/** Ferma il polling. */
		stopPolling() {
			if (this.pollTimerID !== null) {
				clearInterval(this.pollTimerID);
				this.pollTimerID = null;
			}
		},

		/**
		 * Fetch silenzioso: aggiorna la lista solo se sono arrivati nuovi
		 * messaggi, senza mostrare skeleton né nessun refresh visibile.
		 */
		async pollMessages() {
			if (this.loading || this.sending) return;
			try {
				const response = await axios.get(
					`/api/users/${this.myUserID}/conversations/users/${this.conversationID}/messages`
				);
				if (response.status === 200) {
					const fetched = response.data.messages || [];
					const knownIDs = new Set(this.messages.map(m => m.messageID));
					const fetchedIDs = new Set(fetched.map(m => m.messageID));
					const hasNew     = fetched.some(m => !knownIDs.has(m.messageID));
					const hasDeleted = this.messages.some(m => !fetchedIDs.has(m.messageID));

					// Controlla anche se lo stato di qualche messaggio è cambiato (received → read)
					const statusChanged = fetched.some(fm => {
						const existing = this.messages.find(m => m.messageID === fm.messageID);
						return existing && existing.status !== fm.status;
					});

					// Controlla se le reazioni (commenti) di qualche messaggio sono cambiate
					const reactionsChanged = fetched.some(fm => {
						const existing = this.messages.find(m => m.messageID === fm.messageID);
						if (!existing) return false;
						const oldCount = (existing.comments || []).length;
						const newCount = (fm.comments || []).length;
						return oldCount !== newCount;
					});

					if (hasNew || hasDeleted || statusChanged || reactionsChanged) {
						this.messages = fetched;
						if (hasNew) {
							this.$nextTick(() => this.scrollToBottom());
						}
					}
				}
			} catch (_e) {
				// Errori di rete durante il polling vengono ignorati silenziosamente.
			}
		},
	},

	created() {
		this.pollTimerID = null;
	},

	mounted() {
		this.loadMessages();
		this.startPolling();
	},

	unmounted() {
		this.stopPolling();
	},

	watch: {
		/** Ricarica i messaggi e riavvia il polling se cambia la conversazione. */
		conversationID() {
			this.stopPolling();
			this.messages = [];
			this.clearAttachment();
			this.newMessageText = "";
			this.loadMessages();
			this.startPolling();
		},
	},
};
</script>

<template>
	<div
		class="chats-between-users d-flex flex-column"
		style="height: 100vh; overflow: hidden"
	>
		<!--  HEADER: avatar + userName dell'interlocutore  -->
		<div
			class="chat-header d-flex align-items-center gap-3 px-4 py-3 border-bottom bg-body-tertiary flex-shrink-0"
		>
			<!-- Avatar circolare interlocutore -->
			<img
				v-if="otherPhotoSrc"
				:src="otherPhotoSrc"
				:alt="otherUserName"
				class="rounded-circle"
				width="42"
				height="42"
				style="object-fit: cover"
			/>
			<div
				v-else
				class="rounded-circle bg-primary d-flex align-items-center justify-content-center text-white fw-semibold flex-shrink-0"
				style="width: 42px; height: 42px; font-size: 1.05rem"
			>
				{{ otherInitial }}
			</div>

			<!-- Nome utente interlocutore -->
			<span class="fw-semibold fs-6">{{ otherUserName || "Chat" }}</span>
		</div>

		<!--  AREA MESSAGGI  -->
		<div
			ref="messagesList"
			class="messages-area flex-grow-1 px-3 py-4"
			style="overflow-y: auto; background: var(--bs-body-bg)"
		>
			<!-- Placeholder scheletrico durante il caricamento -->
			<template v-if="loading">
				<div
					v-for="i in 6"
					:key="'skel-' + i"
					class="d-flex mb-3"
					:class="i % 2 === 0 ? 'justify-content-end' : 'justify-content-start'"
				>
					<div
						class="skeleton-bubble rounded-3"
						:style="{
							width: (55 + (i * 17) % 100) + 'px',
							height: '38px',
							background: 'var(--bs-secondary-bg)',
							opacity: 0.6,
							animation: 'pulse 1.4s ease-in-out infinite',
						}"
					></div>
				</div>
			</template>

			<!-- Messaggio di errore caricamento -->
			<template v-else-if="errormsg">
				<ModalDangerGeneric :visible="true" :description="errormsg" @close="errormsg = null" />
			</template>

			<!-- Lista messaggi (dal più vecchio al più recente) -->
			<template v-else>
				<div
					v-if="messages.length === 0"
					class="text-center text-secondary py-5"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						width="40"
						height="40"
						fill="currentColor"
						class="bi bi-chat-square-text mb-2 opacity-50"
						viewBox="0 0 16 16"
					>
						<path
							d="M14 1a1 1 0 0 1 1 1v8a1 1 0 0 1-1 1h-2.5a2 2 0 0 0-1.6.8L8 14.333 6.1 11.8a2 2 0 0 0-1.6-.8H2a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1zM2 0a2 2 0 0 0-2 2v8a2 2 0 0 0 2 2h2.5a1 1 0 0 1 .8.4l1.9 2.533a1 1 0 0 0 1.6 0l1.9-2.533a1 1 0 0 1 .8-.4H14a2 2 0 0 0 2-2V2a2 2 0 0 0-2-2z"
						/>
						<path
							d="M3 3.5a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9a.5.5 0 0 1-.5-.5M3 6a.5.5 0 0 1 .5-.5h9a.5.5 0 0 1 0 1h-9A.5.5 0 0 1 3 6m0 2.5a.5.5 0 0 1 .5-.5h5a.5.5 0 0 1 0 1h-5a.5.5 0 0 1-.5-.5"
						/>
					</svg>
					<p class="mb-0 small">Nessun messaggio. Inizia la conversazione!</p>
				</div>

				<template
					v-for="item in messagesWithDateSeparators"
					:key="item.type === 'separator' ? 'sep-' + item.label : 'msg-' + item.data.messageID"
				>
					<!-- Separatore data (pill centrata) -->
					<div v-if="item.type === 'separator'" class="d-flex justify-content-center my-3">
						<span
							class="badge rounded-pill bg-secondary text-body-secondary px-3 py-1"
							style="font-size: 0.72rem; font-weight: 500"
						>
							{{ item.label }}
						</span>
					</div>

					<!-- Messaggio -->
					<div
						v-else
						class="d-flex mb-2"
						:class="isSentByMe(item.data) ? 'justify-content-end' : 'justify-content-start'"
					>
						<!-- Wrapper column: bubble + reazioni -->
						<div
							class="d-flex flex-column"
							:class="isSentByMe(item.data) ? 'align-items-end' : 'align-items-start'"
							style="max-width: 70%"
						>
							<div
								class="message-bubble px-3 py-2 rounded-3 shadow-sm"
								:class="
									isSentByMe(item.data)
										? 'bg-primary text-white'
										: 'bg-body-secondary text-body'
								"
								style="word-break: break-word; width: 100%"
								@contextmenu.prevent="openMenu($event, item.data)"
							>
								<!-- Contenuto testuale -->
								<template v-if="item.data.content_type === 'text'">
									<span>{{ item.data.content }}</span>
								</template>

								<!-- Contenuto foto -->
								<template v-else-if="item.data.content_type === 'photo'">
									<img
										:src="'data:image/jpeg;base64,' + item.data.content"
										alt="Foto"
										class="rounded-2 d-block"
										style="max-width: 220px; max-height: 220px; object-fit: cover; cursor: pointer"
										@click="openImageInNewTab(item.data.content, 'photo')"
									/>
								</template>

								<!-- Contenuto GIF -->
								<template v-else-if="item.data.content_type === 'gif'">
									<img
										:src="'data:image/gif;base64,' + item.data.content"
										alt="GIF"
										class="rounded-2 d-block"
										style="max-width: 220px; max-height: 220px; cursor: pointer"
										@click="openImageInNewTab(item.data.content, 'gif')"
									/>
								</template>

								<!-- Timestamp + stato messaggio -->
								<div
									class="mt-1 d-flex align-items-center justify-content-end gap-1"
									style="font-size: 0.68rem; opacity: 0.72"
								>
									<span>{{ formatTime(item.data.timestamp) }}</span>
								<span v-if="item.data.status === 'read'" title="Letto" style="letter-spacing: -0.15em; font-size: 0.72rem">✓✓</span>
								<span v-else-if="item.data.status === 'received'" title="Ricevuto" style="font-size: 0.72rem">✓</span>
								</div>
							</div>

							<!-- Strip reazioni stile WhatsApp -->
							<div
								v-if="item.data.comments && item.data.comments.length > 0"
								class="reactions-strip d-flex align-items-center flex-wrap gap-1 mt-1"
							>
								<!-- Emoji raggruppate -->
								<span
									v-for="r in groupedReactions(item.data.comments)"
									:key="r.emoji"
									class="reaction-pill d-inline-flex align-items-center gap-1 rounded-pill border bg-body shadow-sm px-2"
									style="font-size: 1rem; padding-top: 2px; padding-bottom: 2px; cursor: default"
									:title="r.count + ' reazion' + (r.count === 1 ? 'e' : 'i')"
								>
									{{ r.emoji }}
									<span
										v-if="r.count > 1"
										class="text-body-secondary"
										style="font-size: 0.68rem"
									>{{ r.count }}</span>
								</span>
								<!-- Totale reazioni -->
								<span
									v-if="item.data.comments.length > 1"
									class="text-body-secondary"
									style="font-size: 0.72rem; font-weight: 600"
								>
									{{ item.data.comments.length }}
								</span>
							</div>
						</div>
					</div>
				</template>
			</template>


		</div>

		<!--  BARRA DI INVIO  -->
		<div
			class="message-input-bar border-top bg-body-tertiary flex-shrink-0 px-3 py-2"
		>
			<!-- Anteprima allegato (se presente) -->
			<div v-if="attachedFilePreview" class="mb-2 d-flex align-items-center gap-2">
				<img
					:src="attachedFilePreview"
					alt="Anteprima"
					class="rounded-2 border"
					style="max-height: 80px; max-width: 120px; object-fit: cover"
				/>
				<span class="text-secondary small">
					{{ attachedFileType === "gif" ? "GIF" : "Foto" }} selezionata
				</span>
				<button
					type="button"
					class="btn btn-sm btn-outline-danger ms-auto"
					@click="clearAttachment"
					aria-label="Rimuovi allegato"
				>
					&times;
				</button>
			</div>

			<!-- Errore invio -->
			<ModalDangerGeneric :visible="!!sendError" :description="sendError || ''" @close="sendError = null" />

			<div class="d-flex align-items-end gap-2">
				<!-- Input testo (solo se non c'è un allegato) -->
				<textarea
					v-model="newMessageText"
					:disabled="!!attachedFile || sending"
					class="form-control"
					rows="1"
					placeholder="Scrivi un messaggio…"
					style="resize: none; max-height: 120px; overflow-y: auto"
					@keydown="onTextKeydown"
					@input="$event.target.style.height = 'auto'; $event.target.style.height = Math.min($event.target.scrollHeight, 120) + 'px'"
					aria-label="Campo messaggio"
				></textarea>

				<!-- Bottone allega Foto -->
				<button
					type="button"
					class="btn btn-outline-secondary flex-shrink-0"
					@click="pickPhoto"
					:disabled="sending"
					title="Allega foto"
					aria-label="Allega foto"
				>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						width="18"
						height="18"
						fill="currentColor"
						class="bi bi-image"
						viewBox="0 0 16 16"
					>
						<path d="M6.002 5.5a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0" />
						<path
							d="M2.002 1a2 2 0 0 0-2 2v10a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V3a2 2 0 0 0-2-2zm12 1a1 1 0 0 1 1 1v6.5l-3.777-1.947a.5.5 0 0 0-.577.093l-3.71 3.71-2.66-1.772a.5.5 0 0 0-.63.062L1.002 12V3a1 1 0 0 1 1-1z"
						/>
					</svg>
				</button>

				<!-- Bottone allega GIF -->
				<button
					type="button"
					class="btn btn-outline-secondary flex-shrink-0"
					@click="pickGif"
					:disabled="sending"
					title="Allega GIF"
					aria-label="Allega GIF"
				>
					<span
						style="font-size: 0.75rem; font-weight: 700; letter-spacing: 0.03em; line-height: 1"
						aria-hidden="true"
					>GIF</span>
				</button>

				<!-- Bottone Invia -->
				<button
					type="button"
					class="btn btn-primary flex-shrink-0"
					@click="sendMessage"
					:disabled="cannotSend"
					aria-label="Invia messaggio"
					title="Invia"
				>
					<svg
						v-if="!sending"
						xmlns="http://www.w3.org/2000/svg"
						width="18"
						height="18"
						fill="currentColor"
						class="bi bi-send-fill"
						viewBox="0 0 16 16"
					>
						<path
							d="M15.964.686a.5.5 0 0 0-.65-.65L.767 5.855H.766l-.452.18a.5.5 0 0 0-.082.887l.41.26.001.002 4.995 3.178 3.178 4.995.002.002.26.41a.5.5 0 0 0 .886-.083zm-1.833 1.89L6.637 10.07l-.215-.338L13.43 2.43z"
						/>
					</svg>
					<span
						v-else
						class="spinner-border spinner-border-sm"
						role="status"
						aria-hidden="true"
					></span>
				</button>
			</div>
		</div>

		<!-- Input file nascosti -->
		<input
			ref="fileInputPhoto"
			type="file"
			accept="image/jpeg,image/png,image/webp"
			class="d-none"
			@change="onPhotoSelected"
		/>
		<input
			ref="fileInputGif"
			type="file"
			accept="image/gif"
			class="d-none"
			@change="onGifSelected"
		/>
	</div>

	<!--
		Menu contestuale (tasto destro su un messaggio).
		Usa <Teleport to="body"> internamente, quindi il rendering avviene
		fuori dalla gerarchia del componente (nessun overflow:hidden che
		potrebbe tagliare il menu).
	-->
	<MenuMessage
		:visible="menuVisible"
		:x="menuX"
		:y="menuY"
		:msg="menuMsg"
		:conversationID="conversationID"
		conversationType="users"
		@close="onMenuClose"
		@message-deleted="onMessageDeleted"
		@message-forwarded="onMessageForwarded"
		@reaction-added="onReactionAdded"
		@reaction-removed="onReactionAdded"
	/>
	<ModalDangerGeneric
		:visible="validationError.visible"
		:title="validationError.title"
		:description="validationError.description"
		@close="validationError.visible = false"
	/>
</template>

<style scoped>
/* Animazione pulse per le skeleton bubble */
@keyframes pulse {
	0%, 100% { opacity: 0.6; }
	50% { opacity: 0.25; }
}

.message-bubble {
	transition: none;
}

/* Scrollbar sottile per l'area messaggi */
.messages-area::-webkit-scrollbar {
	width: 5px;
}
.messages-area::-webkit-scrollbar-thumb {
	background: var(--bs-secondary-bg);
	border-radius: 10px;
}
</style>
