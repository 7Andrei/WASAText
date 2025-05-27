import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import AllChatsView from '../views/AllChatsView.vue'
import ChatView from '../views/ChatView.vue'
import CreateChatView from '../views/CreateChatView.vue'
import UserSettingsView from '../views/UserSettingsView.vue'
import ChatSettingsView from '../views/ChatSettingsView.vue'
import ForwardMessageView from '../views/ForwardMessageView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: HomeView},
		// {path: '/link1', component: HomeView},
		// {path: '/link2', component: HomeView},
		// {path: '/some/:id/link', component: HomeView},

		{path: '/login', component: LoginView},
		{path: '/chats', component: AllChatsView},
		{path: '/chats/:chatId', component: ChatView},
		{path: '/chat', component: CreateChatView},
		{path: '/settings', component: UserSettingsView},
		{path: '/chats/:chatId/settings', component: ChatSettingsView},
		{path: '/chats/:chatId/messages/:messageId', component: ForwardMessageView}
	]
})

export default router
