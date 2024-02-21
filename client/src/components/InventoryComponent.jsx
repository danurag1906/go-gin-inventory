// InventoryComponent.js
import React, { useState, useEffect } from "react";
import axios from "axios";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

const InventoryComponent = () => {
  const navigate = useNavigate();
  const { state: user, dispatch } = useAuth();
  // console.log('user',user);

  const [products, setProducts] = useState([]);
  const [newProduct, setNewProduct] = useState({
    productName: "",
    units: 0,
    price: 0.0,
  });

  useEffect(() => {
    const token = localStorage.getItem("token");
    // console.log('use effect',token);
    if (token == null) {
      navigate("/signin");
    }

    // Fetch all products
    axios
      .get("http://localhost:8080/auth/allProducts", {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      .then((response) => {
        // console.log(response.data);
        setProducts(response.data);
      })
      .catch((error) => console.error(error));
  }, [user.username]);

  const handleInputChange = (e) => {
    const { name, value } = e.target;

    // Convert string values to the correct types
    const typedValue =
      name === "units"
        ? parseInt(value, 10)
        : name === "price"
        ? parseFloat(value)
        : value;

    setNewProduct({ ...newProduct, [name]: typedValue });
  };

  const addProduct = () => {
    // Create a new product
    // console.log(newProduct);
    const token = localStorage.getItem("token");
    // console.log('add product',token);
    axios
      .post("http://localhost:8080/auth/createProduct", newProduct, {
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
      })
      .then((response) => {
        // console.log(response);
        setProducts([...products, response.data.result]);
        setNewProduct({ productName: "", units: 0, price: 0.0 });
      })
      .catch((error) => console.error(error));
  };

  const deleteProduct = (id) => {
    // Delete a product
    // console.log(id);
    const token = localStorage.getItem("token");
    axios
      .delete(`http://localhost:8080/auth/deleteProduct/${id}`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      .then(() => setProducts(products.filter((product) => product.id !== id)))
      .catch((error) => console.error(error));
  };

  const logout = () => {
    // Remove the token from local storage
    dispatch({ type: "SIGNOUT" });
    localStorage.removeItem("token");
    localStorage.removeItem("username");
    // Redirect the user to the signing page
    navigate("/signin");
  };

  return (
    <div className="container mx-auto p-8">
      <h1 className="text-3xl font-bold mb-4">Inventory Management</h1>
      <h3 className="text-xl mb-4">welcome {user.username}</h3>
      <button
        onClick={logout}
        className="bg-red-500 text-white p-2 rounded cursor-pointer mb-4"
      >
        Logout
      </button>
      <div className="mb-8">
        <label className="block mb-2" htmlFor="productName">
          Product Name
        </label>
        <input
          type="text"
          name="productName"
          value={newProduct.productName}
          onChange={handleInputChange}
          className="p-2 border rounded w-full"
        />
        <label className="block mb-2" htmlFor="units">
          Units
        </label>
        <input
          type="number"
          name="units"
          value={newProduct.units}
          onChange={handleInputChange}
          className="p-2 border rounded w-full"
        />
        <label className="block mb-2" htmlFor="price">
          Price
        </label>
        <input
          type="number"
          name="price"
          value={newProduct.price}
          onChange={handleInputChange}
          className="p-2 border rounded w-full"
        />
        <button
          onClick={addProduct}
          className="mt-4 bg-green-500 text-white p-2 rounded cursor-pointer"
        >
          Add Product
        </button>
      </div>
      <div>
        <h2 className="text-2xl font-bold mb-4">Product List</h2>
        <ul>
          {products &&
            products.map((product) => (
              <li key={product.id} className="border p-4 mb-4 rounded">
                <div>
                  {product.productName} - {product.units} units - Rs.
                  {product.price}
                </div>
                <button
                  onClick={() => deleteProduct(product.id)}
                  className="mt-2 bg-red-500 text-white p-2 rounded cursor-pointer"
                >
                  Delete
                </button>
                <Link
                  to={`/updateProduct/${product.id}`}
                  className="bg-blue-500 text-white p-2 rounded cursor-pointer mx-2"
                >
                  Edit
                </Link>
              </li>
            ))}
        </ul>
      </div>
    </div>
  );
};

export default InventoryComponent;
