<script>
export default {
    data: function() {
        return {
            userId: null,
            chatId: null,
            chat: {},
            message: "",
            messagePhoto: null,
            users: [],
            reply: 0,
            replyMessage: null,
        }
    },
    methods: 
    {
        async sendMessage(event)
        {
            event.preventDefault()
            if((this.messagePhoto==null) && (this.message==""))
            {
                this.replyMessage = null
                this.reply = 0
                alert("Please insert a message or a photo")
                return
            }

            let messagePhotoForm = new FormData()
            messagePhotoForm.append('photo', this.messagePhoto)
            messagePhotoForm.append('text', this.message)
            messagePhotoForm.append('sender', parseInt(this.userId))
            messagePhotoForm.append('receiver', parseInt(this.chatId))
            messagePhotoForm.append('reply', this.reply)

            try 
            {
                await this.$axios.post(`/chats/${this.chatId}`, messagePhotoForm, {headers:{Authorization: this.userId}, contentType: 'multipart/form-data'})
                this.refreshMessages()
                this.message=""
                this.messagePhoto=null
            } 
            catch (error) 
            {
                console.log("Errore(placeholder)", error)
            }
            this.replyMessage = null
            this.reply = 0
        },
        handleFileUpload(event)
       {
           this.messagePhoto = event.target.files[0]
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
        forwardMessage(chatId, messageId)
       {
            console.log("Forwarding message", messageId)
            this.$router.push(`/chats/${chatId}/messages/${messageId}`)
        },
        setReply(message)
        {
            this.reply = message.id
            this.replyMessage = message.text
            console.log("Replying to message", message.id)
        },
        clearReply()
        {
            this.reply = 0
            this.replyMessage = null
        }
    },
    async mounted()
    {
        this.userId = sessionStorage.getItem("userId")
        this.chatId = this.$route.params.chatId
        try 
        {
            let response = await this.$axios.get(`/chats/${this.chatId}`, {headers:{Authorization: this.userId}})
            this.chat=response.data
            console.log(this.chat)
            // return this.chat


        } 
        catch (error) 
        {
            console.log("Errore(placeholder)", error)
        }
        try 
        {
            let response = await this.$axios.get("/users", {headers:{Authorization: this.userId}})
            this.users=response.data
            console.log(this.users)
        }
        catch (error) 
        {
            console.log("Errore(placeholder)", error)
        }
        this.refreshInterval = setInterval(() => {this.refreshMessages()}, 2000);
    },
    unmounted() {
        clearInterval(this.refreshInterval);
    },
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
            <div class="row border-bottom border-primary">
                <h3 v-if="chat.chatType == 'private'">
                    <img :src="`data:image/jpeg;base64,${chat.chatParticipants.find(user => user.userId != userId)?.userPhoto}`" height="64" width="64"  v-if ="chat.chatType=='private' && chat.chatParticipants.find(user => user.userId != userId)?.userPhoto">
                    <img src="https://placehold.co/64x64?text=Placeholder" height="64" width="64" alt="Placeholder" v-else>
                    {{ chat.chatParticipants.find(participant => participant.userId != userId)?.userName || 'Private Chat' }}
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
        <!-- DAFARE: Doppie spunte -->
        <div class="mt-4">
            <div class="row justify-content-start">
                <div class="col-md-10">
                    <!-- DAFARE: Fare card dei messaggi piu' carina -->
                    <div v-for="message in chat.chatMessages" :key="message.id" :class="['d-flex mb-2', isSender(message.sender) ? 'justify-content-end' : 'justify-content-start']">
                        <div class="card" style="max-width: 75%;">
                            <div class="card-body">
                                <div class="row">
                                    <div  class="col-7">    
                                        <h5 class="card-title">
                                            {{ getUser(message.sender) }}
                                            <img
                                                v-if="chat.chatParticipants.find(user => user.userId === message.sender)?.userPhoto"
                                                :src="`data:image/jpeg;base64,${chat.chatParticipants.find(user => user.userId == message.sender)?.userPhoto}`"
                                                height="32"
                                                width="32"
                                            >
                                            <img v-else src="https://placehold.co/32x32?text=User" height="32" width="32">
                                        </h5>
                                    </div>
                                    <div class="col-5">
                                        <button @click="forwardMessage(chatId, message.id)" class="btn btn-link">
                                            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#arrow-right"/></svg>
                                        </button>
                                        <button @click="setReply(message)" class="btn btn-link">
                                            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#arrow-left"/></svg>
                                        </button>
                                    </div>
                                </div>
                                <p v-if="message.forwarded"> Forwarded from {{ users.find(user => user.userId == message.forwarded)?.userName || 'Unknown' }} </p>
                                <img v-if="message.photo" :src="`data:image/jpeg;base64, ${message.photo}`" height="200" width="200" alt="Message Photo" class="mb-2">
                                <p v-if="message.reply" class="text-muted">Replying to: {{ chat.chatMessages.find(msg => msg.id === message.reply)?.text || 'Message not found' }} </p>
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
                                            {{ reaction.reaction }} <small>{{ getUser(reaction.userId) }}</small>
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

            <div class="message-form-container">
                <div class="row" v-if="reply">
                    <div class="alert alert-info text-truncate col-6">
                        Replying to: {{ replyMessage }}
                    </div>
                    <div class="col-6">
                        <button @click="clearReply" class="btn btn-link align-text-center"> 
                            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#x"/></svg>    
                        </button>
                    </div>
                </div>
                <form @submit.prevent="sendMessage">
                    <div class="row">
                        <div class="col-6">
                            <input type="text" class="form-control" v-model="message" placeholder="Your message">
                        </div>
                        <div class="col-4">
                            <input type="file" class="form-control" id="chatPhoto" @change="handleFileUpload">
                        </div>
                        <div class="col-2">
                            <button type="submit" class="btn btn-primary">Send</button>
                        </div>
                    </div>
                </form>
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