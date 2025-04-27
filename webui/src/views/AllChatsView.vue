<script>
export default {
    data: function() {
        return {
            userId: null,
            chatId: null,
            chats: [],
            users:[],
            searchUser: "",
            foundUsers:[],
            show: false,
            privateChats: [],
            groupChats: [],
        }
    },
    methods: 
    {
        async createChat(privateId)
        {   

            let participants=[{userId: privateId, userName: "", userPhoto: null}]
            this.error=null
            let chat = new FormData()
            let returnId = null

            chat.append('chatType', "private")
            console.log("participants", participants)
            chat.append('chatParticipants', JSON.stringify(participants))
            try 
            {
                let response = await this.$axios.post('/createchat', chat, {headers:{Authorization: this.userId}, contentType: 'multipart/form-data'})
                console.log(response.data)
                returnId = response.data
                // this.$router.push("/chats")
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
                        this.$router.push(`/chats/${privateChat.id}`)
                        return
                    }
                }
            }
            let chatId = await this.createChat(privateUserId)
            this.$router.push(`/chats/${chatId}`)
        },
    },
    async mounted()
    {
        this.userId = sessionStorage.getItem("userId")
        console.log(this.userId)
        try 
        {
            let response = await this.$axios.get("/chats", {headers:{Authorization: this.userId}})
            this.chats=response.data
            // console.log(this.chats)
        } 
        catch (error) 
        {
            console.log("Errore(DaCambiare)", error)
        }
        

        try 
        {
            let response = await this.$axios.get("/chats", { headers: { Authorization: this.userId } });
            this.chats = response.data;

            this.chats.sort((a, b) => {

                const lastMessageA = a.chatMessages[a.chatMessages.length - 1];
                const lastMessageB = b.chatMessages[b.chatMessages.length - 1];

                const timeA = lastMessageA ? new Date(lastMessageA.dateTime).getTime() : 0;
                const timeB = lastMessageB ? new Date(lastMessageB.dateTime).getTime() : 0;

                return timeB - timeA;
            }
            );

            // console.log("Sorted Chats:", this.chats);
        } 
        catch (error) 
        {
            console.log("Errore(DaCambiare)", error);
        }


        try 
        {
            let response = await this.$axios.get("/users", {headers:{Authorization: this.userId}})
            this.users=response.data
            this.users.splice(this.users.indexOf(this.userId), 1)
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
        console.log(this.privateChats)
        // console.log(this.groupChats)

    },
    watch: {
        searchUser: function(){
            if(this.searchUser.length>3)
            {
                this.foundUsers = this.users.filter(user => user.userName.toLowerCase().includes(this.searchUser.toLowerCase()))
                return this.foundUsers
            }
            else
            {
                this.foundUsers = this.users.filter(user => user.userName.toLowerCase().includes(this.searchUser.toLowerCase()))
                return this.users
            }
        }
    }

}
</script>
<template>

    <div class="container mt-4">
        <div v-if="!userId" class="row rounded">
            <div class="text-center">
                <h3>Unauthorized</h3>
            </div>
            <button @click="$router.push('/login')" class="btn btn-primary">Login</button>
        </div>
        <div v-else>
            <div class="row">
                <div class="col-md-12 mb-4 justify-content-center">
                    <div class="input-group mb-3">
                        <span class="input-group-text">Lente</span>
                        <div class="form-floating">
                            <input 
                                type="text" 
                                class="form-control"   
                                id="privateChat"   
                                placeholder="Search user" 
                                v-model="searchUser" >
                            <label for="floatingInputGroup1">Username</label>
                        </div>
                    </div>
                </div>
                <div class="col-md-12 mb-4" v-for="user in foundUsers" :key="user.userId">
                    <button @click="createOrSend(user.userId)" class="btn btn-primary">{{ user.userName }}</button>
                </div>
                <div class="col-md-12" v-for="chat in chats" :key="chat.id">
                    <div class="card mb-4">
                        <div class="card-body">
                            <div class="col-md-12">
                                <div class="row">
                                    <div class="col-10">
                                        <div v-if="chat.chatType=='group'">
                                        <h5 class="card-title">
                                            <a :href="`/#/chats/${chat.id}`" class="card-title">{{ chat.chatName }}</a>
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
                                        </h5>
                                        <p class="card-text">{{ chat.chatType }}</p>
                                    </div>
                                    <div class="col-2">
                                        <img :src="`data:image/jpeg;base64,${chat.chatPhoto}`" height="128" width="128" alt="Chat Photo">
                                    </div>
                                    <div class="row ms-2">
                                        <div class="col-12" v-for="message in chat.chatMessages" :key="message.id">
                                            {{ message.text }}
                                        </div>
                                    </div>
                                </div>
                            </div>
                            
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>

     
</template>