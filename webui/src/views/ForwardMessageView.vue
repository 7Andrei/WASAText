<script>
export default {
    data: function() {
        return {
            userId: 0,
            chatId: 0,
            chats: [],
            messageId: 0,
        }
    },
    methods: 
    {
        forwardMessage(chatId)
        {
            this.chatId = chatId
            let response = this.$axios.post(`/chats/${this.chatId}/messages/${this.messageId}`, {}, {headers:{Authorization: this.userId}})
            this.$router.push(`/chats/${this.chatId}`)
        }
    },
    async mounted()
    {
        this.userId = sessionStorage.getItem("userId")
        this.messageId = this.$route.params.messageId
        // console.log(this.userId)
        try 
        {
            let response = await this.$axios.get("/chats", {headers:{Authorization: this.userId}})
            // console.log(response.data)
            this.chats=response.data


        } 
        catch (error) 
        {
            console.log("Errore(placeholder)", error)
        }
    }
}
</script>
<template>
    <div class="container mt-4">
        <div class="row">
            <div class="col-md-12">
                <h1>Forward Message</h1>
                <p>Select a chat to forward the message to:</p>
            </div>
            <div class="col-md-12" v-for="chat in chats" :key="chat.id">
                <div class="card mb-4">
                    <div class="card-body">
                        <div class="col-md-12">
                            <div class="row">
                                <div class="col-10">
                                    <!-- <h5 class="card-title"><a :href="`/#/chats/${chat.id}`" class="card-title">{{ chat.chatName }}</a></h5> -->
                                    <button @click="forwardMessage(chat.id)" class="btn btn-link" >
                                        <div v-if="chat.chatType=='group'">
                                            <h5 class="card-title">
                                                <a :href="`/#/chats/${chat.id}`" class="card-title">
                                                    {{ chat.chatName }}
                                                </a>
                                                <div class="badge bg-primary ms-2">
                                                    {{ chat.chatType }}
                                                </div>
                                            </h5>

                                            <div class="d-flex flex-row flex-wrap">
                                                <span v-for="participant in chat.chatParticipants" :key="participant.userId" class="badge bg-secondary me-2">
                                                    {{ participant.userName }}
                                                </span>
                                            </div>
                                        </div>
                                        <h5 v-else class="card-title">
                                            <a :href="`/#/chats/${chat.id}`" class="card-title">
                                                {{ chat.chatParticipants.find(participant => participant.userId != userId)?.userName || 'Private Chat' }}
                                            </a>
                                            <div class="badge bg-primary ms-2">
                                                {{ chat.chatType }}
                                            </div>
                                        </h5>
                                    </button>
                                </div>
                                <div class="col-2 text-end">
                                        <img :src="`data:image/jpeg;base64,${chat.chatPhoto}`" height="64" width="64" alt="Chat Photo" v-if="chat.chatPhoto && chat.chatType=='group'">
                                        <img :src="`data:image/jpeg;base64,${chat.chatParticipants.find(user => user.userId != userId)?.userPhoto}`" height="64" width="64"  v-else-if ="chat.chatType=='private' && chat.chatParticipants.find(user => user.userId != userId)?.userPhoto">
                                        <img src="https://placehold.co/64x64?text=Placeholder" height="64" width="64" alt="Placeholder" v-else>
                                    </div>
                            </div>
                        </div>
                        <!-- <ul class="list-group list-group-flush">
                            <li class="list-group-item" v-for="participant in chat.chatParticipants" :key="participant.userId">
                                <img :src="`data:image/jpeg;base64,${participant.userPhoto}`" class="rounded-circle me-2" alt="User Photo" width="30" height="30">
                                {{ participant.userName }}
                            </li>
                        </ul> -->
                    </div>
                </div>
            </div>
        </div>
    </div>

     
</template>