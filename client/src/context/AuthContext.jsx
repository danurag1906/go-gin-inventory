import React, { createContext, useContext, useReducer } from "react";

const AuthContext = createContext();

const initialState = {
  token: localStorage.getItem("token") || null,
  username: localStorage.getItem("username") || null,
};

const authReducer = (state, action) => {
  switch (action.type) {
    case "SIGNIN":
      return {
        ...state,
        token: action.payload.token,
        username: action.payload.username,
      };
    case "SIGNOUT":
      return {
        ...state,
        token: null,
        username: null,
      };
    default:
      return state;
  }
};

const AuthProvider=({children})=>{
    const [state,dispatch]=useReducer(authReducer,initialState);

    return (
        <AuthContext.Provider value={{state,dispatch}} >
            {children}
        </AuthContext.Provider>
    )
}

const useAuth=()=>{
    const context=useContext(AuthContext)
    if(!context){
        throw new Error('useAuth must be used within an AuthProvider')
    }
    return context
}

export {AuthProvider,useAuth}