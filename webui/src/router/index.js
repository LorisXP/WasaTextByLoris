import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from '../views/LoginView.vue'
import MainLayout from '../views/MainLayout.vue'
import HomeView from '../views/HomeView.vue'
import StartConversation from '../components/StartConversation.vue'
import ChatsBetweenUsers from '../components/ChatsBetweenUsers.vue'
import ChatsBetweenUserAndGroup from '../components/ChatsBetweenUserAndGroup.vue'
import auth from '../services/auth.js'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		/* Rotta di login (senza sidebar) */
		{path: '/', component: LoginView},

		/* Layout principale: Sidebar a sinistra + pannello destro (rotte figlie) */
		{
			path: '/',
			component: MainLayout,
			children: [
				/* Vista predefinita dopo il login */
				{
					path: 'chats',
					name: 'chats',
					component: HomeView,
				},
				/* Avvia nuova conversazione con un utente */
				{
					path: 'start/:userName',
					name: 'start-conversation',
					component: StartConversation,
					props: route => ({
						userName: route.params.userName,
						userPhoto: route.query.photo || null,
					}),
				},
				/* Chat diretta tra utenti */
				{
					path: 'conversations/users/:conversationID',
					name: 'chat-users',
					component: ChatsBetweenUsers,
					props: route => ({
						conversationID: Number(route.params.conversationID),
						otherUserName: route.query.userName || '',
						otherUserPhoto: route.query.photo || null,
					}),
				},
				/* Chat di gruppo */
				{
					path: 'conversations/groups/:conversationID',
					name: 'chat-groups',
					component: ChatsBetweenUserAndGroup,
					props: route => ({
						conversationID: Number(route.params.conversationID),
						groupName: route.query.groupName || '',
						groupPhoto: route.query.photo || null,
						groupID: route.query.groupID ? Number(route.query.groupID) : 0,
					}),
				},
			],
		},
	]
})

/**
 * Guard globale di navigazione.
 * - Se l'utente non è autenticato e tenta di accedere a una rotta protetta
 *   (qualsiasi rotta al di fuori di '/'), lo reindirizza al login.
 * - Questo gestisce anche il caso di F5 (refresh): poiché l'auth è in-memory,
 *   dopo il refresh l'utente non risulta autenticato e viene riportato alla home.
 * - Se l'utente è già autenticato e naviga verso '/' (login), viene reindirizzato
 *   a '/chats' per evitare di mostrare nuovamente la schermata di login.
 */
router.beforeEach((to, _from, next) => {
	if (to.path === '/') {
		// Pagina di login: se già autenticato, salta direttamente alle chat.
		if (auth.isAuthenticated()) {
			next('/chats');
		} else {
			next();
		}
	} else if (!auth.isAuthenticated()) {
		// Rotta protetta senza autenticazione: torna al login.
		next('/');
	} else {
		next();
	}
});

export default router
