<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
			userId: null
		}
	},
	methods: {
		login() {
			this.userId = sessionStorage.getItem("userId");
		},
		goToProfile() {
            this.$router.push(`/settings`);
        },
		goToLogin() {
            this.$router.push('/login');
        },

	},
	mounted() {
		this.login();

		if (this.userId) {
			this.$router.push('/chats');
		} 
	}
}
</script>

<template>
	<div>
		<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Home page</h1>
			<div class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
                    <button v-if="userId" type="button" class="btn btn-sm btn-outline-primary" @click="goToProfile">
                        Profile
                    </button>
                    <button v-else type="button" class="btn btn-sm btn-outline-primary" @click="goToLogin">
                        Login
                    </button>
                </div>
			</div>
		</div>
		<div class="text-center">
			<button v-if="!userId" type="button" class="btn btn-primary" @click="goToLogin">
				Login
			</button>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>

<style>
</style>
