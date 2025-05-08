<script>
export default {
    data: function() {
        return {
            userId: null,
            chatName: "",
            chatType: "",
            chatPhoto: null,
            chatParticipants: {},
            users: {},
            error: null
        }
    },
    methods: 
    {
        async createChat(event)
        {
            this.error=null
            event.preventDefault()
            let chat = new FormData()
            chat.append('chatName', this.chatName)
            chat.append('chatType', "group")
            chat.append('chatPhoto', this.chatPhoto)
            // this.chatParticipants[this.userId] = { userId: this.userId };
            chat.append('chatParticipants', JSON.stringify(this.chatParticipants))
            try 
            {
                let response = await this.$axios.post('/createchat', chat, {headers:{Authorization: this.userId}, contentType: 'multipart/form-data'})
                // console.log(response.data)
                // this.$router.push("/chats")
            } 
            catch (error) 
            {
                console.log("Errore(placeholder)", error)
                console.log(this.chatParticipants)
                console.log(this.chatName)
                console.log(this.chatType)
                console.log(this.chatPhoto)
                this.error = error.response.data
            }
        },
        handleFileUpload(event)
        {
            this.chatPhoto = event.target.files[0]
        }
    },
    async mounted()
    {
        this.userId = sessionStorage.getItem("userId")
        try 
        {
            let response = await this.$axios.get(`/users`, {headers:{Authorization: this.userId}})
            this.users=response.data
            console.log("Prima", this.users)
            this.users = this.users.filter(user => user.userId !== parseInt(this.userId));
            console.log("Dopo", this.users)
        } 
        catch (error) 
        {
            console.log("Errore(placeholder)", error)
            this.error = error
        }

        // this.users = this.users.filter(user => user.id !== this.userId)
        // console.log(this.users)
    }
    
}
</script>
<template>
    <div v-if="error" class="alert alert-danger" role="alert">
        {{ error }}
    </div>
    <div class="container mt-4">
        <div v-if="!userId" class="row rounded">
            <div class="text-center">
                <h3>Unauthorized</h3>
            </div>
            <button @click="$router.push('/login')" class="btn btn-primary">Login</button>
        </div>
        <div v-else>
            <div class="row mb-4">
                <h3>Create a New Chat</h3>
            </div>
            <form @submit.prevent="createChat">
                <div class="row">
                    <div class="mb-3 col-5">
                        <label for="chatName" class="form-label">Chat Name</label>
                        <input type="text" class="form-control" id="chatName" v-model="chatName" required>
                    </div>
                    <div class="mb-3 col-5">
                        <label for="chatPhoto" class="form-label">Chat Photo</label>
                        <input type="file" class="form-control" id="chatPhoto" @change="handleFileUpload">
                    </div>
                    <!-- <div class="mb-3 col-2">
                        <label for="chatType" class="form-label">Chat Type</label>
                        <div class="form-check">
                            <input class="form-check-input" type="radio" id="group" value="group" required v-model="chatType">
                            <label class="form-check-label" for="group">Group</label>
                        </div>
                        <div class="form-check">
                            <input class="form-check-input" type="radio" id="private" value="private" required v-model="chatType">
                            <label class="form-check-label" for="private">Private</label>
                        </div>
                    </div> -->
                </div>
                <div class="row">
                    <div class="mb-3 col-10">
                        <label for="chatParticipants" class="form-label">Select Participants</label>
                        <select multiple class="form-control" id="chatParticipants" v-model="chatParticipants">
                            <option v-for="user in users" :key="user.userId" :value="user">{{ user.userName }}</option>
                        </select>
                    </div>
                </div>
                <button type="submit" class="btn btn-primary">Create Chat</button>
            </form>
        </div>
    </div>
</template>
