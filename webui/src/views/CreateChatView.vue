<script>
export default {
    data: function() {
        return {
            userId: null,
            chatName: "",
            chatType: "",
            chatPhoto: null,
            chatParticipants: {},
        }
    },
    methods: 
    {
        async createChat(event)
        {
            event.preventDefault()
            try 
            {
                let response = await this.$axios.post('/createchat', {
                    chatName: this.chatName,
                    chatType: this.chatType,
                    chatPhoto: this.chatPhoto,
                    chatParticipants: this.chatParticipants},
                    {headers:{Authorization: this.userId}})
                // console.log(response.data)
                this.$router.push("/chats")
            } 
            catch (error) 
            {
                console.log("puppa", error)
            }
        },
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
        getUser(userId)
        {
            const user = this.chat.chatParticipants.find(participant => participant.userId === userId);
            console.log(user.userName)
            return user.userName
        },
        isSender(senderId)
        {
            // console.log("senderId, userId", senderId, this.userId)
            // console.log(senderId == this.userId)
            return senderId == this.userId
        }
    },
    async mounted()
    {
        this.userId = sessionStorage.getItem("userId")
        // try 
        // {
        //     let response = await this.$axios.get(`/chats/${this.chatId}`, {headers:{Authorization: this.userId}})
        //     this.chat=response.data


        // } 
        // catch (error) 
        // {
        //     console.log("puppa", error)
        // }
    }
}
</script>
<template>
    <div class="container mt-4">
        <div v-if="!userId" class="row border rounded">
            <h3>Unauthorized</h3>
            <button @click="$router.push('/login')" class="btn btn-primary">Login</button>
        </div>
        <div v-else>
            <div class="row border rounded mb-4">
                <h3>Create a New Chat</h3>
            </div>
            <form @submit.prevent="createChat">
                <div class="mb-3">
                    <label for="chatName" class="form-label">Chat Name</label>
                    <input type="text" class="form-control" id="chatName" v-model="chatName" required>
                </div>
                <div class="mb-3">
                    <label for="chatPhoto" class="form-label">Chat Photo</label>
                    <input type="file" class="form-control" id="chatPhoto" @change="handleFileUpload">
                </div>
                <div class="mb-3">
                    <label for="chatType" class="form-label">Chat Type</label>
                    <input type="text" class="form-control" id="chatType" v-model="chatType" required>
                </div>
                <div class="mb-3">
                    <label for="participants" class="form-label">Select Participants</label>
                    <select multiple class="form-control" id="participants" v-model="participants">
                        <option v-for="user in allUsers" :key="user.id" :value="user.id">{{ user.name }}</option>
                    </select>
                </div>
                <button type="submit" class="btn btn-primary">Create Chat</button>
            </form>
        </div>
    </div>
</template>
