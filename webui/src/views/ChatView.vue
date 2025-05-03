<script>
export default {
    data: function() {
        return {
            userId: null,
            chatId: null,
            chat: {},
            message: "",
            messagePhoto: null,
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
                console.log("Errore(placeholder)", error)
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
            // console.log(userId)
            // return user.userName
            return user ? user.userName : 'U.n.k.n.o.w.n';
        },
        isSender(senderId)
        {
            return senderId == this.userId
        },
        async deleteMessage(messageId)
        {
            try 
            {
                let response = await this.$axios.delete(`/chats/${this.chatId}/messages/${messageId}/delete`, {headers:{Authorization: this.userId}})
                this.refreshMessages()
            } 
            catch (error) 
            {
                console.log("Errore(placeholder)", error)
            }
        },
        async addReaction(messageId, reaction)
        {
            try 
            {
                const message = this.chat.chatMessages.find(msg => msg.id === messageId);
                console.log(message);
                if (message.reactions.length>0)
                {
                    for (let i=0; i<message.reactions.length; i++)
                    {
                        if (message.reactions[i].userId==this.userId)
                        {
                            await this.deleteReaction(messageId, message.reactions[i].id)
                            console.log("Reazione gia' presente", message.reactions[i].reaction)
                            // break
                        }
                    }
                }

                let response = await this.$axios.post(`/chats/${this.chatId}/messages/${messageId}/reactions`, {reaction: reaction}, {headers:{Authorization: this.userId}})
                console.log("Added reaction")
                this.refreshMessages()
            } 
            catch (error) 
            {
                console.log("Errore(placeholder)", error)
            }
        },
        async deleteReaction(messageId, reactionId)
        {
            try 
            {
                let response = await this.$axios.delete(`/chats/${this.chatId}/messages/${messageId}/reactions/${reactionId}`, {headers:{Authorization: this.userId}})
                this.refreshMessages()
            } 
            catch (error) 
            {
                console.log("Errore(placeholder)", error)
            }
        },
        async sendPhoto(event)
       {
        try 
            {
                let messagePhotoForm = new FormData()
                messagePhotoForm.append('Photo', this.messagePhoto)
                let response = await this.$axios.post(`/chats/${this.chatId}`, 
                messagePhotoForm, 
                {headers:{Authorization: this.userId}, 
                contentType: 'multipart/form-data'})
                this.refreshMessages()
            } 
            catch (error) 
            {
                console.log("Errore(placeholder)", error)
            }
       },
    },
    async mounted()
    {
        this.userId = sessionStorage.getItem("userId")
        this.chatId = this.$route.params.chatId
        try 
        {
            let response = await this.$axios.get(`/chats/${this.chatId}`, {headers:{Authorization: this.userId}})
            this.chat=response.data
            // console.log(this.chat)
            return this.chat


        } 
        catch (error) 
        {
            console.log("Errore(placeholder)", error)
        }
    }
}
</script>
<template>

    <div v-if="!userId" class="row rounded">
        <div class="text-center">
            <h3>Unauthorized</h3>
        </div>
        <button @click="$router.push('/login')" class="btn btn-primary">Login</button>
    </div>
    <div v-else>
        <div class="container mt-4 row">
            <div class="row border rounded">
                <h3 v-if="chat.chatType == 'private'">
                    {{ chat.chatParticipants.find(participant => participant.userId != userId)?.userName || 'Private Chat' }}
                    <!-- DAFARE: Mettere foto chat privata/placeholder -->
                </h3>
                <h3 v-else-if ="chat.chatType == 'group'">
                    <img :src="`data:image/jpeg;base64,${chat.chatPhoto}`" height="64" width="64" alt="Chat Photo" v-if="chat.chatPhoto" class="mt-2 me-2">
                    <img src="https://placehold.co/64x64?text=Placeholder" height="64" width="64" alt="Placeholder" v-else class="mt-2 me-2">
                    <router-link :to="`/chats/${chatId}/settings`">
                        {{ chat.chatName }}
                    </router-link>
                </h3>
                <p v-else>Loading chat</p>
            </div>
        </div>
        <div class="mt-4">
            <div class="row justify-content-start">
                <div class="col-md-10">
                    <!-- DAFARE: Fare card dei messaggi piu' carina -->
                    <div v-for="message in chat.chatMessages" :key="message.id" :class="['d-flex mb-2', isSender(message.sender) ? 'justify-content-end' : 'justify-content-start']">
                        <div class="card" style="max-width: 50%;">
                            <div class="card-body">
                                <router-link :to="`/chats/${chatId}/messages/${message.id}`">
                                    <h5 class="card-title">{{ getUser(message.sender) }}</h5>
                                    <!-- DAFARE: Fixare forward dei messaggi -->
                                </router-link>
                                <p class="card-text">{{ message.text }}</p>
                                <small class="text-muted float-end">{{ message.dateTime }}</small>
                                <div class="d-flex flex-column align-items-end mt-2">
                                    <div class="d-flex">
                                        <button v-if="isSender(message.sender)" @click="deleteMessage(message.id)" class="btn btn-link">
                                            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#trash-2"/></svg>
                                        </button>
                                        <div class="dropdown ms-2">
                                            <button class="btn btn-link dropdown-toggle" type="button" id="emojiMenuButton" data-bs-toggle="dropdown" aria-expanded="false">
                                                😊
                                            </button>
                                            <ul class="dropdown-menu" aria-labelledby="emojiMenuButton">
                                                <li><span class="dropdown-item" @click="addReaction(message.id, '😀')">😀</span></li>
                                                <li><span class="dropdown-item" @click="addReaction(message.id, '😂')">😂</span></li>
                                                <li><span class="dropdown-item" @click="addReaction(message.id, '😍')">😍</span></li>
                                                <li><span class="dropdown-item" @click="addReaction(message.id, '😢')">😢</span></li>
                                                <li><span class="dropdown-item" @click="addReaction(message.id, '👍')">👍</span></li>
                                                <li><span class="dropdown-item" @click="addReaction(message.id, '👎')">👎</span></li>
                                            </ul>
                                        </div>
                                    </div>
                                    <div class="d-flex flex-wrap mt-2">
                                        <button v-for="reaction in message.reactions" :key="reaction.id" @click="deleteReaction(message.id, reaction.id)" class="btn btn-sm btn-outline-secondary me-1">
                                            {{ reaction.reaction }}
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                    <div class="mb-5"></div>
                </div>
            </div>
        </div>

        <div class="container row">
            <div class="message-form-container">
                <div class="row">
                    <div class="col-8">
                        <form @submit.prevent="sendMessage">
                            <div class="input-group">
                                <input type="text" class="form-control" v-model="message" placeholder="Your message">
                                <button type="submit" class="btn btn-primary">Send</button>
                            </div>
                        </form>
                    </div>
                    <!-- DAFARE: Finire upload foto -->
                    <!-- <div class="col-4">
                        <form @submit.prevent="sendPhoto">
                            <div class="input-group">
                                <input type="file" class="form-control" id="chatPhoto" @change="handleFileUpload" required>
                                <button type="submit" class="btn btn-primary">Change Photo</button>
                            </div>
                        </form>
                    </div> -->
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.message-form-container {
    position: fixed;
    bottom: 0;
    width: 75%;
    padding: 10px;
    }
</style>