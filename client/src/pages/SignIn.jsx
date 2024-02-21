// Signin.jsx

import React, { useState } from "react";
import axios from "axios";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext.jsx";

const SignIn = () => {
  const navigate = useNavigate();
  const { dispatch } = useAuth();
  const [user, setUser] = useState({
    username: "",
    password: "",
  });

  const [error, setError] = useState("");

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setUser({ ...user, [name]: value });
  };

  const handleSignin = async () => {
    try {
      const response = await axios.post("http://localhost:8080/signin", user);
      //   console.log(response);

      const { token, username } = response.data;

      dispatch({
        type: "SIGNIN",
        payload: { token, username },
      });

      localStorage.setItem("token", response?.data?.token);
      localStorage.setItem("username", response?.data?.username);

      if (response.status === 200) {
        navigate("/home");
        setError("");
      }
    } catch (error) {
      console.error(error);
      setError("Invalid Credentials");
    }
  };

  return (
    <div className="container mx-auto mt-8">
      <h2 className="text-2xl font-bold mb-4">Signin</h2>
      <form className=" container mx-auto max-w-md">
        <div className="mb-4">
          <label
            htmlFor="username"
            className="block text-sm font-medium text-gray-600"
          >
            Username:
          </label>
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
          <label
            htmlFor="password"
            className="block text-sm font-medium text-gray-600"
          >
            Password:
          </label>
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
          onClick={handleSignin}
          className="bg-blue-500 text-white p-2 rounded-md hover:bg-blue-600"
        >
          Signin
        </button>
        <p className="text-red-700 font-bold">{error}</p>
        <p>
          Dont have account ?{" "}
          <Link to={"/signup"} className="text-blue-700 font-bold">
            SignUp
          </Link>
        </p>
      </form>
    </div>
  );
};

export default SignIn;
