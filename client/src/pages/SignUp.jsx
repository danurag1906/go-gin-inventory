// Signup.jsx

import React, { useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';

const SignUp = () => {
    const navigate=useNavigate()
   const [user, setUser] = useState({
      username: '',
      password: '',
   });

   const handleInputChange = (e) => {
      const { name, value } = e.target;
      setUser({ ...user, [name]: value });
   };

   const handleSignup = () => {
      axios.post('http://localhost:8080/signup', user)
         .then(response => {
            console.log(response);
            if(response.status==201){
                navigate("/signin")
            }
            // Handle successful signup, e.g., redirect to signin page
         })
         .catch(error => console.error(error));
   };

   return (
      <div className="container mx-auto mt-8">
         <h2 className="text-2xl font-bold mb-4">Signup</h2>
         <form className="max-w-md">
            <div className="mb-4">
               <label htmlFor="username" className="block text-sm font-medium text-gray-600">Username:</label>
               <input
                  type="text"
                  id="username"
                  name="username"
                  value={user.username}
                  onChange={handleInputChange}
                  className="mt-1 p-2 border rounded-md w-full"
               />
            </div>
            <div className="mb-4">
               <label htmlFor="password" className="block text-sm font-medium text-gray-600">Password:</label>
               <input
                  type="password"
                  id="password"
                  name="password"
                  value={user.password}
                  onChange={handleInputChange}
                  className="mt-1 p-2 border rounded-md w-full"
               />
            </div>
            <button
               type="button"
               onClick={handleSignup}
               className="bg-blue-500 text-white p-2 rounded-md hover:bg-blue-600"
            >
               Signup
            </button>
         </form>
      </div>
   );
};

export default SignUp;
