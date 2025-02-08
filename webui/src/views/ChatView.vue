<script>
export default {
    data: function() {
        return {
            userId: null,
            chatId: null,
            chat: {},
            message: ""
        }
    },
    methods: 
    {
        async sendMessage(event)
        {
            event.preventDefault()
            // console.log(this.message)
            // console.log(this.userId)
            // console.log(this.chatId)
            try 
            {
                let response = await this.$axios.post(`/chats/${this.chatId}`, {
                    text: this.message,
                    photo: null,
                    sender: parseInt(this.userId),
                    receiver: parseInt(this.chatId)},
                    {headers:{Authorization: this.userId}})
                // console.log(response.data)
                this.refreshMessages()
                this.message=""
            } 
            catch (error) 
            {
                console.log("puppa", error)
            }
        },

        async refreshMessages()
        {
            let response = await this.$axios.get(`/chats/${this.chatId}`, {headers:{Authorization: this.userId}})
            this.chat=response.data
        },
        getUser(userId)
        {   
            const user = this.chat.chatParticipants.find(participant => participant.userId === userId);
            console.log(userId)
            // return user.userName
            return user ? user.userName : 'U.n.k.n.o.w.n';
        },
        isSender(senderId)
        {
            return senderId == this.userId
        }
    },
    async mounted()
    {
        this.userId = sessionStorage.getItem("userId")
        // console.log(this.userId)
        this.chatId = this.$route.params.chatId
        // console.log("vue chatId", this.chatId)
        try 
        {
            // let response = await this.$axios.get('/chats/$chatId', {headers:{Authorization: this.userId}})
            let response = await this.$axios.get(`/chats/${this.chatId}`, {headers:{Authorization: this.userId}})
            // console.log(response.data)
            this.chat=response.data
            // console.log(this.chat)


        } 
        catch (error) 
        {
            console.log("puppa", error)
        }
    }
}
</script>
<template>

    <div class="container mt-4 row">
        <div class="row border rounded">
            <h3 v-if="chat">
                <router-link :to="`/chats/${chatId}/settings`">{{ chat.chatName }}</router-link>
            </h3>
            <p v-else>Loading chat</p>
        </div>
    </div>
    <div class="mt-4">
        <div class="row justify-content-start">
            <div class="col-md-8">
                <div v-for="message in chat.chatMessages" :key="message.id" :class="['d-flex mb-2', isSender(message.sender) ? 'justify-content-end' : 'justify-content-start']">
                    <div class="card" style="max-width: 50%;">
                        <div class="card-body">
                            <router-link :to="`/chats/${chatId}/messages/${message.id}`">
                                <h5 class="card-title">{{ getUser(message.sender) }}</h5>
                            </router-link>
                            <p class="card-text">{{ message.text }}</p>
                            <small class="text-muted float-end">{{ message.dateTime }}</small>
                        </div>
                    </div>
                </div>
            </div>
        </div>

    </div>
    <div class="container row">
        <div class="message-form-container">
            <form @submit.prevent="sendMessage">
                <div class="input-group">
                    <input type="text" class="form-control" v-model="message" placeholder="Your message">
                    <button type="submit" class="btn btn-primary">Send</button>
                </div>
            </form>
        </div>
    </div>
</template>

<style scoped>
.message-form-container {
    position: fixed;
    bottom: 0;
    width: 50%;
    padding: 10px;
    }
</style>