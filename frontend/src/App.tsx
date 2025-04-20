import { useState, useEffect } from "react";
import reactLogo from "./assets/react.svg";
import viteLogo from "/vite.svg";
import "./App.css";

function App() {
  const [count, setCount] = useState(0);
  const [message, setMessage] = useState("Loading..."); // Add state for the message

  // Fetch data from backend API on component mount
  useEffect(() => {
    fetch("http://localhost:8080/api") // Assuming backend runs on port 8080
      .then((response) => response.json())
      .then((data) => setMessage(data.message))
      .catch((error) => {
        console.error("Error fetching data:", error);
        setMessage("Failed to load message from backend.");
      });
  }, []); // Empty dependency array ensures this runs only once on mount

  return (
    <>
      <div>
        <a href="https://vite.dev" target="_blank">
          <img src={viteLogo} className="logo" alt="Vite logo" />
        </a>
        <a href="https://react.dev" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
      </div>
      <h1 className="text-red-500">Vite + React</h1>
      {/* Display the message from the backend */}
      <p>Message from backend: {message}</p>
      <div className="card">
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
        <p>
          Edit <code>src/App.tsx</code> and save to test HMR
        </p>
      </div>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p>
    </>
  );
}

export default App;
