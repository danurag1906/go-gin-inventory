import axios from "axios";

const instance=axios.create({
    baseURL:'http://localhost:8080',
})

//get the token from local storage
const token=localStorage.getItem('token')

//set the authorization header for all requests
instance.defaults.headers.common['Authorization']=`Bearer ${token}`

export default instance