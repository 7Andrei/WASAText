<script>
export default {
    data: function() {
        return {
            userId: null,
            chatPhoto: null,
            chatName: "",
            chatId: 0,
            // chat: {},
            chatParticipants: {},
            users: {},
            user: {},
            usersToAdd: {},
        }
    },
    methods: 
    {
       async changePhoto(event)
       {
        try 
            {
                let chatPhotoForm = new FormData()
                chatPhotoForm.append('chatPhoto', this.chatPhoto)
                let response = await this.$axios.put(`/chats/${this.chatId}/settings/changephoto`, chatPhotoForm, {headers:{Authorization: this.userId}, contentType: 'multipart/form-data'})
                alert("Photo changed")
            } 
            catch (error) 
            {
                console.log("puppa", error)
            }
       },
       async changeName(event)
       {
            event.preventDefault()
            try
            {
                let response = await this.$axios.put(`/chats/${this.chatId}/settings/changename`, {chatName: this.chatName}, {headers:{Authorization: this.userId}})
                alert("Name changed")
            }
            catch (error)
            {
                console.log("puppa", error)
            }
       },
       handleFileUpload(event)
       {
           this.chatPhoto = event.target.files[0]
       },
       async addUser(event)
       {
           event.preventDefault()
           try 
           {
               let response = await this.$axios.post(`/chats/${this.chatId}/settings/add`, {chatParticipants: this.usersToAdd}, {headers:{Authorization: this.userId}})
               alert("User added")
           } 
           catch (error) 
           {
               console.log("puppa", error)
           }
       },
       async leaveChat()
       {
            try
            {
                let response = await this.$axios.delete(`/chats/${this.chatId}/settings/leave`, {headers:{Authorization: this.userId}})
                alert("Chat left")
                this.$router.push("/chats")
            }
            catch (error)
            {
                console.log("puppa", error)
            }
       }
       
    },
    async mounted()
    {
        this.userId = sessionStorage.getItem("userId")
        this.chatId = this.$route.params.chatId
        try 
        {
            let response = await this.$axios.get(`/users`, {headers:{Authorization: this.userId}})
            this.users=response.data
        } 
        catch (error) 
        {
            console.log("puppa", error)
        }
        try 
        {
            let response = await this.$axios.get(`/chats/${this.chatId}`, {headers:{Authorization: this.userId}})
            this.chat=response.data
            this.chatParticipants = this.chat.chatParticipants
        } 
        catch (error) 
        {
            console.log("puppa", error)
        }
        this.users = this.users.filter(user => !this.chatParticipants.some(participant => participant.userId === user.userId))
    }
}
</script>
<template>
    <div class="row">
        <h2>Chat Settings</h2>
        <div class="row">
            <div class="col-6">
                <form @submit.prevent="changePhoto" class="mt-4">
                    <div class="mb-3">
                        <label for="chatPhoto" class="form-label">Change Photo</label>
                        <input type="file" class="form-control" id="chatPhoto" @change="handleFileUpload" required>
                    </div>
                    <button type="submit" class="btn btn-primary">Change Photo</button>
                </form>
            </div>
            <div class="col-6">
                <form @submit.prevent="changeName" class="mt-4">
                    <div class="mb-3">
                        <label for="chatName" class="form-label">Change Name</label>
                        <input type="text" class="form-control" id="chatName" v-model="chatName" required>
                    </div>
                    <button type="submit" class="btn btn-primary">Change Name</button>
                </form>
            </div>
            <div class="col-6">
                <label for="chatParticipants" class="form-label">Select Participants</label>
                <form @submit.prevent="addUser">
                    <select multiple class="form-control mt-4 mb-3" id="chatParticipants" v-model="usersToAdd">
                        <option v-for="user in users" :key="user.userId" :value="user">{{ user.userName }}</option>
                    </select>
                    <button type="submit" class="btn btn-primary">Add Users</button>
                </form>
            </div>
            <div class="col-6">
                <!-- <button @click="logout" class="btn btn-danger mt-4">Logout</button> -->
                <button @click = "leaveChat" class="btn btn-danger mt-5">Leave Chat</button>
            </div>
        </div>
        <!-- <div class="col-4">
            <button @click="logout" class="btn btn-danger mt-4">Logout</button>
        </div> -->
    </div>
</template>
