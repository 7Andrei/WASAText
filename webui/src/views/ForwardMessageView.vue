<script>
export default {
    data: function() {
        return {
            messageId: 0,
            userId: null,
            chatId: null,
            chats: [],
            users:[],
            foundUsers:[],
            show: false,
            privateChats: [],
            groupChats: [],
            
        }
    },
    methods: 
    {
        forwardMessage(chatId)
        {
            this.chatId = chatId
            let response = this.$axios.post(`/chats/${this.chatId}/messages/${this.messageId}`, {}, {headers:{Authorization: this.userId}})
            this.$router.push(`/chats/${this.chatId}`)
        },


        async createChat(privateId)
        {   
            let participants=[{userId: privateId, userName: "", userPhoto: null}]
            this.error=null
            let chat = new FormData()
            let returnId = null

            chat.append('chatType', "private")
            chat.append('chatParticipants', JSON.stringify(participants))
            try 
            {
                let response = await this.$axios.post('/chat', chat, {headers:{Authorization: this.userId}, contentType: 'multipart/form-data'})
                console.log(response.data)
                returnId = response.data
            } 
            catch (error) 
            {
                console.log("Errore(placeholder)", error)
                console.log(this.chatName)
                console.log(this.chatType)
                this.error = error.response.data
            }
            return returnId
        },
        async createOrSend(privateUserId)
        {
            for (let privateChat of this.privateChats)
            {
                for (let user of privateChat.chatParticipants)
                {
                    if (user.userId==privateUserId)
                    {
                        this.forwardMessage(privateChat.id)
                        // this.$router.push(`/chats/${privateChat.id}`)
                        return
                    }
                }
            }
            let chatId = await this.createChat(privateUserId)
            this.forwardMessage(chatId)
            // this.$router.push(`/chats/${chatId}`)
        },
    },
    async mounted()
    {
        this.userId = sessionStorage.getItem("userId")
        this.messageId = this.$route.params.messageId
        try 
        {
            let response = await this.$axios.get("/chats", {headers:{Authorization: this.userId}})
            this.chats=response.data


        } 
        catch (error) 
        {
            console.log("Errore(placeholder)", error)
        }
        try 
        {
            let response = await this.$axios.get("/users", {headers:{Authorization: this.userId}})
            this.users=response.data
            this.users = this.users.filter(user => user.userId != this.userId)
        }
        catch (error) 
        {
            console.log("Errore(DaCambiare)", error)
        }
        for (let chat of this.chats)
        {
            if(chat.chatType == "private")
            {
                this.privateChats.push(chat)
            }
            else
            {
                this.groupChats.push(chat)
            }
        }
    },
}
</script>
<template>
    <div class="container mt-4">
        <div class="row">
            <div class="col-md-12">
                <h1>Forward Message</h1>
                <p>Select a chat to forward the message to:</p>
            </div>
            <div class="col-md-12" v-for="chat in groupChats" :key="chat.id">
                <div class="card mb-4">
                    <div class="card-body">
                        <div class="col-md-12">
                            <div class="row">
                                <div class="col-10">
                                    <button @click="forwardMessage(chat.id)" class="btn btn-link" >
                                        <h5 class="card-title">
                                            <a :href="`/#/chats/${chat.id}`" class="card-title">
                                                {{ chat.chatName }}
                                            </a>
                                            <div class="badge bg-primary ms-2">
                                                <!-- {{ chat.chatType }} -->
                                                group
                                            </div>
                                        </h5>
                                    </button>
                                    <div class="d-flex flex-row flex-wrap">
                                        <span v-for="participant in chat.chatParticipants" :key="participant.userId" class="badge bg-secondary me-2">
                                            {{ participant.userName }}
                                        </span>
                                    </div>
                                </div>
                                <div class="col-2 text-end">
                                    <img :src="`data:image/jpeg;base64,${chat.chatPhoto}`" height="64" width="64" alt="Chat Photo" v-if="chat.chatPhoto">
                                    <img src="https://placehold.co/64x64?text=Placeholder" height="64" width="64" alt="Placeholder" v-else>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>



            <div class="col-md-12" v-for="user in users" :key="user.userId">
                <div class="card mb-4">
                    <div class="card-body">
                        <div class="col-md-12">
                            <div class="row">
                                <div class="col-10">
                                    <button @click="createOrSend(user.userId)" class="btn btn-link" >
                                        <h5 class="card-title">
                                                  {{ user.userName }}
                                            <div class="badge bg-primary ms-2">
                                                private
                                            </div>
                                        </h5>
                                    </button>
                                </div>
                                <div class="col-2 text-end">
                                    <img v-if="user.userPhoto" :src="`data:image/jpeg;base64,${user.userPhoto}`" height="64" width="64" alt="User Photo">
                                    <img v-else src="https://placehold.co/64x64?text=Placeholder" height="64" width="64" alt="Placeholder">
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>



        </div>
    </div>

     
</template>