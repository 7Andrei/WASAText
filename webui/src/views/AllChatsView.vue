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
        }
    },
    methods: 
    {
        async createChat(privateId)
        {   

            let participants=[{userId: privateId, userName: "", userPhoto: null}]
            this.error=null
            let chat = new FormData()

            chat.append('chatType', "private")
            console.log("participants", participants)
            chat.append('chatParticipants', JSON.stringify(participants))
            try 
            {
                let response = await this.$axios.post('/createchat', chat, {headers:{Authorization: this.userId}, contentType: 'multipart/form-data'})
                // console.log(response.data)
                // this.$router.push("/chats")
            } 
            catch (error) 
            {
                console.log("Errore(placeholder)", error)
                console.log(this.chatName)
                console.log(this.chatType)
                this.error = error.response.data
            }
        },
    },
    async mounted()
    {
        this.userId = sessionStorage.getItem("userId")
        console.log(this.userId)
        try 
        {
            let response = await this.$axios.get("/chats", {headers:{Authorization: this.userId}})
            // console.log(response.data)
            this.chats=response.data


        } 
        catch (error) 
        {
            console.log("Errore(DaCambiare)", error)
        }
        try 
        {
            let response = await this.$axios.get("/users", {headers:{Authorization: this.userId}})
            this.users=response.data
            this.users.splice(this.users.indexOf(this.userId), 1)
            console.log(this.users)
        }
        catch (error) 
        {
            console.log("Errore(DaCambiare)", error)
        }
    },
    watch: {
        searchUser: function(){
            if(this.searchUser.length>0)
            {
                this.foundUsers = this.users.filter(user => user.userName.toLowerCase().includes(this.searchUser.toLowerCase()))
                return this.foundUsers
            }
            else
            {
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
                            <input type="text" class="form-control" id="privateChat" placeholder="Search user" v-model="searchUser">
                            <label for="floatingInputGroup1">Username</label>
                        </div>
                    </div>
                </div>
                <div class="col-md-12 mb-4" v-for="user in foundUsers" :key="user.userId">
                    <button @click="createChat(user.userId)" class="btn btn-primary">{{ user.userName }}</button>
                </div>
                <div class="col-md-12" v-for="chat in chats" :key="chat.id">
                    <div class="card mb-4">
                        <div class="card-body">
                            <div class="col-md-12">
                                <div class="row">
                                    <div class="col-10">
                                        <h5 class="card-title"><a :href="`/#/chats/${chat.id}`" class="card-title">{{ chat.chatName }}</a></h5>
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
    </div>

     
</template>