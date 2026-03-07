import { createApp, reactive } from "vue";
import App from "./App.vue";
import router from "./router";
import axios from "./services/axios.js";
import auth from "./services/auth.js";
import validate from "./services/validator.js";

import MenuMessage from "./components/MenuMessage.vue";
import ModalCreateGroup from "./components/ModalCreateGroup.vue";
import GroupEdit from "./components/GroupEdit.vue";
import ModalDangerGeneric from "./components/ModalDangerGeneric.vue";
import Sidebar from "./components/Sidebar.vue";
import StartConversation from "./components/StartConversation.vue";


import "./assets/dashboard.css";
import "./assets/main.css";

const app = createApp(App);

app.config.globalProperties.$axios = axios;
app.config.globalProperties.$apiDomain =
	import.meta.env.VITE_API_URL || "http://localhost:3000";
app.config.globalProperties.$auth = auth;

app.component("Sidebar", Sidebar);
app.component("StartConversation", StartConversation);
app.component("MenuMessage", MenuMessage);
app.component("ModalCreateGroup", ModalCreateGroup);
app.component("GroupEdit", GroupEdit);
app.component("ModalDangerGeneric", ModalDangerGeneric);

app.config.globalProperties.$validator = validate;

app.use(router);
app.mount("#app");
