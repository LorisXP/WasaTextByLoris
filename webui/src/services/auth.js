import { reactive } from "vue";
import { setUserID as setAxiosUserID } from "./axios.js";

// In-memory reactive auth state. No persistence (no cookies, no localStorage),
// in accordance with the simplified login spec: the server returns an identifier
// which must be sent in the Authorization header for subsequent calls.
const state = reactive({
	userID: null,
	userName: null,
	userPhoto: null,
});

function init() {
	// No-op: do not read any persisted identifier (no sessions/cookies)
}

function setUserID(id) {
	state.userID = id || null;
	setAxiosUserID(id);
}

function setUserName(name) {
	state.userName = name || null;
}

function setUserPhoto(photo) {
	state.userPhoto = photo || null;
}

function clear() {
	setUserID(null);
	setUserName(null);
	state.userPhoto = null;
}

function isAuthenticated() {
	return !!state.userID;
}

export default {
	state,
	init,
	setUserID,
	setUserName,
	setUserPhoto,
	clear,
	isAuthenticated,
};
