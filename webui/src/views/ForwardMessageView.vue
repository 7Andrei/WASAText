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
        console.log(this.userId)
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
            <div class="col-md-12" v-for="chat in chats" :key="chat.id">
                <div class="card mb-4">
                    <div class="card-body">
                        <div class="col-md-12">
                            <div class="row">
                                <div class="col-10">
                                    <!-- <h5 class="card-title"><a :href="`/#/chats/${chat.id}`" class="card-title">{{ chat.chatName }}</a></h5> -->
                                    <button @click="forwardMessage(chat.id)" class="btn btn-link" >
                                        <h5 class="card-title">{{ chat.chatName }}</h5>
                                    </button>
                                    <p class="card-text">{{ chat.chatType }}</p>
                                </div>
                                <div class="col-2">
                                    <img :src="`data:image/jpeg;base64,${chat.chatPhoto}`" height="128" width="128" alt="Chat Photo">
                                </div>
                            </div>
                        </div>
                        <ul class="list-group list-group-flush">
                            <li class="list-group-item" v-for="participant in chat.chatParticipants" :key="participant.userId">
                                <!-- <img :src="`data:image/jpeg;base64,${participant.userPhoto}`" class="rounded-circle me-2" alt="User Photo" width="30" height="30"> -->
                                {{ participant.userName }}
                            </li>
                        </ul>
                    </div>
                </div>
            </div>
        </div>
    </div>

     
</template>