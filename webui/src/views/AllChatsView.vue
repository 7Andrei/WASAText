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
                let response = await this.$axios.post('/chat', chat, {headers:{Authorization: this.userId}, contentType: 'multipart/form-data'})
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
        async refreshChats()
        {
            try 
            {
                let response = await this.$axios.get("/chats", {headers:{Authorization: this.userId }})
                this.chats = response.data;
                if (this.chats==null)
                {
                    this.chats = []
                }

                this.chats.sort((a, b) => {

                    const lastMessageA = a.chatMessages[a.chatMessages.length - 1];
                    const lastMessageB = b.chatMessages[b.chatMessages.length - 1];

                    const timeA = lastMessageA ? new Date(lastMessageA.dateTime).getTime() : 0;
                    const timeB = lastMessageB ? new Date(lastMessageB.dateTime).getTime() : 0;

                    return timeB - timeA;
                }
                );

            } 
            catch (error) 
            {
                console.log("Errore(DaCambiare)", error);
            }


        },
    },
    async mounted()
    {
        this.userId = sessionStorage.getItem("userId")
        console.log(this.userId)

        try 
        {
            let response = await this.$axios.get("/chats", { headers: { Authorization: this.userId } });
            this.chats = response.data;
            if (this.chats==null)
            {
                this.chats = []
            }

            this.chats.sort((a, b) => {

                const lastMessageA = a.chatMessages[a.chatMessages.length - 1];
                const lastMessageB = b.chatMessages[b.chatMessages.length - 1];

                const timeA = lastMessageA ? new Date(lastMessageA.dateTime).getTime() : 0;
                const timeB = lastMessageB ? new Date(lastMessageB.dateTime).getTime() : 0;

                return timeB - timeA;
            }
            );

        } 
        catch (error) 
        {
            console.log("Errore(DaCambiare)", error);
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
        this.refreshInterval = setInterval(() => {this.refreshChats()}, 2000);

    },
    unmounted() {
        clearInterval(this.refreshInterval);
    },
    watch: {
        searchUser: function(){
            if(this.searchUser.length>2)
            {
                this.foundUsers = this.users.filter(user => user.userName.toLowerCase().includes(this.searchUser.toLowerCase()))
                console.log(this.foundUsers)
            }
            else if(this.searchUser == "/")
            {
                this.foundUsers = this.users
            }
            else
            {
                this.foundUsers = []
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
                            <label for="floatingInputGroup1">Username or "/" for all users</label>
                        </div>
                    </div>
                </div>
                <div class="col-md-12 mb-4" v-for="user in foundUsers" :key="user.userId">
                    <button @click="createOrSend(user.userId)" class="btn btn-primary">{{ user.userName }}</button>
                </div>
                <div v-if="!chats" class="col-md-12">
                    <div class="alert alert-danger text-center" role="alert">
                        No chats found
                    </div>
                    <div class="text-center">
                        <button @click="$router.push('/chat')" class="btn btn-primary">
                            Create a new chat
                        </button>
                        <h2 class="mt-4">Or search Username above to create a private chat</h2>
                    </div>
                </div>
                <div class="col-md-12" v-else v-for="chat in chats" :key="chat.id">
                    <div class="card mb-4">
                        <div class="card-body">
                            <div class="col-md-12">
                                <div class="row">
                                    <div class="col-10">
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
                                    </div>
                                    <div class="col-2 text-end">
                                        <img :src="`data:image/jpeg;base64,${chat.chatPhoto}`" height="64" width="64" alt="Chat Photo" v-if="chat.chatPhoto && chat.chatType=='group'">
                                        <img :src="`data:image/jpeg;base64,${chat.chatParticipants.find(user => user.userId != userId)?.userPhoto}`" height="64" width="64"  v-else-if ="chat.chatType=='private' && chat.chatParticipants.find(user => user.userId != userId)?.userPhoto">
                                        <img src="https://placehold.co/64x64?text=Placeholder" height="64" width="64" alt="Placeholder" v-else>
                                    </div>
                                    <div class="row">
                                        <div class="col-12 mt-2 border-top border-primary" v-if="chat.chatMessages.length">
                                            <span class="badge bg-primary mt-2 me-1">
                                                {{ chat.chatParticipants.find(user => user.userId == chat.chatMessages[chat.chatMessages.length - 1].sender)?.userName || 'Unknown User' }}
                                            </span>
                                            <span class="mt-2">
                                                {{ chat.chatMessages[chat.chatMessages.length - 1].text }}
                                            </span>
                                            <svg v-if="chat.chatMessages[chat.chatMessages.length-1].photo" class="feather">
                                                <use href="/feather-sprite-v4.29.0.svg#image"/>
                                            </svg>
                                            <span class="badge bg-secondary float-end mt-2">
                                                {{ chat.chatMessages[chat.chatMessages.length - 1].dateTime }}
                                            </span>
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