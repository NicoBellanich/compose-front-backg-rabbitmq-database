import { useState } from "react";

function App() {
  const [msg, setMsg] = useState("");
  const [resp, setResp] = useState("");

  const sendMsg = async () => {
    if (!msg.trim()) return;

    const currentMsg = msg;
    setMsg("");

    try {
      // POST over local /publish
      const res = await fetch(`${process.env.REACT_APP_BACKEND_URL}/publish`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ msg: currentMsg }),
      });


      const text = await res.text();
      setResp(text);

      // clean response after 3 secs
      setTimeout(() => {
        setResp("");
      }, 3000);

    } catch (error) {
      console.error("Error sending the message:", error);
      setResp("Error sending the message");
      setTimeout(() => setResp(""), 3000);
    }
  };

  return (
    <div className="p-6 text-center">
      <h1>RabbitMQ Demo 🐇</h1>
      <input
        placeholder="Escribí un mensaje"
        value={msg}
        onChange={(e) => setMsg(e.target.value)}
      />
      <button onClick={sendMsg}>Enviar</button>
      <p>{resp}</p>
    </div>
  );
}

export default App;
