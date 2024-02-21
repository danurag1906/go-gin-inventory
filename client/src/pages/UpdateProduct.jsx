import React, { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import axios from "axios";
import { useAuth } from "../context/AuthContext.jsx";

const UpdateProduct = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const { state: user } = useAuth();

  const [editedProduct, setEditedProduct] = useState({
    productName: "",
    units: 0,
    price: 0.0,
  });

  useEffect(() => {
    let token = localStorage.getItem("token");

    if (!user.token) {
      navigate("/signin");
    }

    // Fetch data for the specific item with the given id
    axios
      .get(`http://localhost:8080/auth/products/${id}`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      .then((response) => {
        // Update the state with the fetched data
        setEditedProduct(response.data);
      })
      .catch((error) => console.error(error));
  }, [id, user]);

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    // Convert string values to the correct types
    const typedValue =
      name === "units"
        ? parseInt(value, 10)
        : name === "price"
        ? parseFloat(value)
        : value;

    setEditedProduct({ ...editedProduct, [name]: typedValue });
  };

  const handleUpdateProduct = () => {
    // Update the product with the edited details
    const token = localStorage.getItem("token");
    axios
      .put(`http://localhost:8080/auth/updateProduct/${id}`, editedProduct, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      .then((response) => {
        // Redirect back to the main page after successful update
        if (response.status == 200) {
          navigate("/home");
        }
      })
      .catch((error) => console.error(error));
  };

  return (
    <div className="container mx-auto mt-8">
      <h2 className="text-2xl font-bold mb-4">Edit Product</h2>
      <form className="max-w-md">
        <div className="mb-4">
          <label
            htmlFor="productName"
            className="block text-sm font-medium text-gray-600"
          >
            Product Name:
          </label>
          <input
            type="text"
            id="productName"
            name="productName"
            value={editedProduct.productName}
            onChange={handleInputChange}
            className="mt-1 p-2 border rounded-md w-full"
          />
        </div>
        <div className="mb-4">
          <label
            htmlFor="units"
            className="block text-sm font-medium text-gray-600"
          >
            Units:
          </label>
          <input
            type="number"
            id="units"
            name="units"
            value={editedProduct.units}
            onChange={handleInputChange}
            className="mt-1 p-2 border rounded-md w-full"
          />
        </div>
        <div className="mb-4">
          <label
            htmlFor="price"
            className="block text-sm font-medium text-gray-600"
          >
            Price:
          </label>
          <input
            type="number"
            id="price"
            name="price"
            value={editedProduct.price}
            onChange={handleInputChange}
            className="mt-1 p-2 border rounded-md w-full"
          />
        </div>
        <button
          type="button"
          onClick={handleUpdateProduct}
          className="bg-blue-500 text-white p-2 rounded-md hover:bg-blue-600"
        >
          Update Product
        </button>
      </form>
    </div>
  );
};

export default UpdateProduct;
