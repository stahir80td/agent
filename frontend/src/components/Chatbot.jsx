import { useState, useRef, useEffect } from 'react'
import axios from 'axios'

export default function Chatbot({ onClose }) {
  const [messages, setMessages] = useState([
    { role: 'assistant', content: 'Hi! Ask me anything about our trading platform.' }
  ])
  const [input, setInput] = useState('')
  const [loading, setLoading] = useState(false)
  const messagesEndRef = useRef(null)

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" })
  }

  useEffect(scrollToBottom, [messages])

  const sendMessage = async () => {
    if (!input.trim()) return

    const userMessage = { role: 'user', content: input }
    setMessages(prev => [...prev, userMessage])
    setInput('')
    setLoading(true)

    try {
      const response = await axios.post('/api/chat', {
        message: input
      })

      setMessages(prev => [...prev, {
        role: 'assistant',
        content: response.data.response
      }])
    } catch (error) {
      setMessages(prev => [...prev, {
        role: 'assistant',
        content: '❌ Sorry, something went wrong. Please try again.'
      }])
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="chatbot-widget">
      <div className="chatbot-header">
        <h3>💬 Trading Assistant</h3>
        <button onClick={onClose}>✕</button>
      </div>

      <div className="chatbot-messages">
        {messages.map((msg, idx) => (
          <div key={idx} className={`message ${msg.role}`}>
            <div className="message-content">{msg.content}</div>
          </div>
        ))}
        {loading && <div className="message assistant typing">Thinking...</div>}
        <div ref={messagesEndRef} />
      </div>

      <div className="chatbot-input">
        <input
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyPress={(e) => e.key === 'Enter' && sendMessage()}
          placeholder="Ask about trading..."
        />
        <button onClick={sendMessage} disabled={loading}>
          Send
        </button>
      </div>
    </div>
  )
}