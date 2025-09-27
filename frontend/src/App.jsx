import { useState } from 'react'
import Navbar from './components/Navbar'
import Hero from './components/Hero'
import Features from './components/Features'
import Chatbot from './components/Chatbot'
import './App.css'

function App() {
  const [showChat, setShowChat] = useState(false)

  return (
    <div className="app">
      <Navbar />
      <Hero />
      <Features />
      
      {/* Floating chat button */}
      <button 
        className="chat-fab"
        onClick={() => setShowChat(!showChat)}
      >
        💬
      </button>

      {/* Chat widget */}
      {showChat && <Chatbot onClose={() => setShowChat(false)} />}
    </div>
  )
}

export default App