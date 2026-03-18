import axios from "axios";

const instance = axios.create({
	baseURL: __API_URL__,
	timeout: 1000 * 5,
});
// Set Authorization header to use Bearer <identifier> as required by the spec.
// No cookies or HTTP sessions; no client-side persistence performed here.
export function setUserID(id) {
	if (id) {
		instance.defaults.headers.common["Authorization"] = `Bearer ${id}`;
	} else {
		delete instance.defaults.headers.common["Authorization"];
	}
}

export default instance;
