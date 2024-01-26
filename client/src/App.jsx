
import './App.css'
import { BrowserRouter,Routes,Route } from 'react-router-dom'
import UpdateProduct from './pages/UpdateProduct'
import Home from './components/Home'
import SignUp from './pages/SignUp'
import SignIn from './pages/SignIn'

function App() {
  return (
    <BrowserRouter>
    <Routes>
      <Route path='/home' element={<Home/>} />
      <Route path='/updateProduct/:id' element={<UpdateProduct/>} />
      <Route path='/signup' element={<SignUp/>} />
      <Route path='/signin' element={<SignIn/>} />
    </Routes>
    </BrowserRouter>
  )
}

export default App
