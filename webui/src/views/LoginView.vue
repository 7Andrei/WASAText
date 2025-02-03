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
            } 
            catch (error) 
            {
                console.log("puppa", error)
            }
        }
    },
    mounted()
    {
        console.log(this.userName)
    }
}
</script>

<template>
    <form @submit.prevent="login">
        <input type="text" v-model="userName">
        <button type="submit">Puppa chi clicca</button>
    </form>

    <h1 v-if = "userId">
        {{ userId }}
    </h1>
</template>