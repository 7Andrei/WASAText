<script>
export default {
    data: function() {
        return {
            userName: "",
            userId: null
        }
    },
    methods: 
    {
        async login(event)
        {
            event.preventDefault()
            console.log(this.userName)
            try 
            {
                let response = await this.$axios.post("/session", {userName: this.userName})
                this.userId = response.data.userId
                sessionStorage.setItem("userId", this.userId)
                this.$emit("login")
                this.$router.push("/chats")
            } 
            catch (error) 
            {
                console.log("Errore(placeholder)", error)
            }
        }
    },
    mounted()
    {
        console.log(this.userName)
        this.userId = sessionStorage.getItem("userId")
        if(this.userId != null)
        {
            console.log("UserId trovato")
            this.$router.push("/chats")
        }
    }
}
</script>

<template>
    <div class="container d-flex justify-content-center align-items-center mt-5">
        <div class="card p-4 shadow" style="width: 100%; max-width: 400px;">
            <form @submit.prevent="login">
                <div class="mb-3">
                    <input id="username" type="text" v-model="userName" class="form-control" placeholder="Enter your username" required>
                </div>
                <button type="submit" class="btn btn-primary w-100">Login</button>
            </form>

            <h1 v-if="userId" class="mt-4 text-center text-success">
                Welcome  {{ userName }}
            </h1>
        </div>
    </div>
</template>