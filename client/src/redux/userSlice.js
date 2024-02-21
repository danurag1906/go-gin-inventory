import {createSlice} from '@reduxjs/toolkit'

export const userSlice=createSlice({
    name:'user',
    initialState:{
        token:localStorage.getItem('token') || null,
        username:localStorage.getItem('username') || null,
    },
    reducers:{
        setUser:(state,action)=>{
            const {token,username}=action.payload;
            state.token=token;
            state.username=username;
        },
        clearUser:(state)=>{
            state.token=null;
            state.username=null;
        },
    },
});

export const {setUser,clearUser}=userSlice.actions;

export const selectUser=(state)=>state.user;

export default userSlice.reducer;