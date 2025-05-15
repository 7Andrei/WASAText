<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>
<script>
export default {
	data: function() {
		return {
			userId: null
		}
	},
	methods: {
		login() {
			this.userId = sessionStorage.getItem("userId");
		},
		logout()
       	{
           sessionStorage.removeItem("userId")
           this.$router.push("/login")
       	},
	},
	mounted() {
		this.login();
	}
}
</script>

<template>

	<header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
		<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">WASA Texto</a>
		<button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
			<span class="navbar-toggler-icon"></span>
		</button>
	</header>

	<div class="container-fluid">
		<div class="row">
			<nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
				<div class="position-sticky pt-3 sidebar-sticky" v-if="userId">
					<ul class="nav flex-column">
						<!-- <li class="nav-item">
							<RouterLink to="/" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
								Home
							</RouterLink>
						</li> -->
						<!-- <li class="nav-item">
							<RouterLink to="/login" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#layout"/></svg>
								Login
							</RouterLink>
						</li> -->
						<li class="nav-item">
							<RouterLink to="/chats" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#message-circle"/></svg>
								All chats
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink to="/createchat" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#plus-circle"/></svg>
								New group chat
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink to="/settings" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#settings"/></svg>
								Settings
							</RouterLink>
						</li>
						<li class="nav-item">
							<button type="button" class="btn btn-danger ms-2" @click="logout">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#log-out"/></svg>
								Logout
							</button>
						</li>
					</ul>
				</div>
				<div class="position-sticky pt-3 sidebar-sticky" v-else>
					<ul class="nav flex-column">
						<li class="nav-item">
							<RouterLink to="/" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
								Home
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink to="/login" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#layout"/></svg>
								Login
							</RouterLink>
						</li>
					</ul>
				</div>
			</nav>

			<main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
				<RouterView @login="login" />
			</main>
		</div>
	</div>
</template>

<style>
</style>
