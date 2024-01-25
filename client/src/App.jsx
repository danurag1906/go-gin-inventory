
import './App.css'
import { BrowserRouter,Routes,Route } from 'react-router-dom'
import UpdateProduct from './UpdateProduct'
import Home from './components/Home'

function App() {
  return (
    <BrowserRouter>
    <Routes>
      <Route path='/' element={<Home/>} />
      <Route path='/updateProduct/:id' element={<UpdateProduct/>} />
    </Routes>
    </BrowserRouter>
  )
}

export default App
