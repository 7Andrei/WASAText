<script>
export default {
    data: function() {
        return {
            userId: null,
            userPhoto: null,
            uerName: "",
        }
    },
    methods: 
    {
       async changePhoto(event)
       {
        try 
            {
                let userPhotoForm = new FormData()
                userPhotoForm.append('userPhoto', this.userPhoto)
                let response = await this.$axios.put('/settings/changephoto', userPhotoForm, {headers:{Authorization: this.userId}, contentType: 'multipart/form-data'})
                alert("Photo changed")
                // console.log(response.data)
            } 
            catch (error) 
            {
                console.log("Errore(placeholder)", error)
            }
       },
       async changeName(event)
       {
            event.preventDefault()
            try
            {
                let response = await this.$axios.put(`/settings/changename`, {userName: this.userName}, {headers:{Authorization: this.userId}})
                alert("Name changed")
                // console.log(response.data)
            }
            catch (error)
            {
                console.log("Errore(placeholder)", error)
            }
       },
       handleFileUpload(event)
       {
            this.userPhoto = event.target.files[0]
       },
       logout()
       {
            this.userId = null
            sessionStorage.removeItem("userId")
            this.$emit("Logout")
            this.$router.push("/login")
       }
    },
    async mounted()
    {
        this.userId = sessionStorage.getItem("userId")
    }
}
</script>
<template>
    <div class="row">
        <h2>Settings</h2>
        <div class="row">
            <div class="col-6">
                <form @submit.prevent="changePhoto" class="mt-4">
                    <div class="mb-3">
                        <label for="userPhoto" class="form-label">Change Photo</label>
                        <input type="file" class="form-control" id="userPhoto" @change="handleFileUpload" required>
                    </div>
                    <button type="submit" class="btn btn-primary">Change Photo</button>
                </form>
            </div>
            <div class="col-6">
                <form @submit.prevent="changeName" class="mt-4">
                    <div class="mb-3">
                        <label for="userName" class="form-label">Change Name</label>
                        <input type="text" class="form-control" id="userName" v-model="userName" required>
                    </div>
                    <button type="submit" class="btn btn-primary">Change Name</button>
                </form>
            </div>
        </div>
        <div class="col-4">
            <button @click="logout" class="btn btn-danger mt-4">Logout</button>
        </div>
    </div>
</template>
