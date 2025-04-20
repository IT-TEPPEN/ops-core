import { useState, useEffect } from "react";
import "./App.css";

function App() {
  const [count, setCount] = useState(0);
  const [message, setMessage] = useState("Loading...");

  useEffect(() => {
    // 環境変数 VITE_API_HOST を読み込む。未定義の場合は /api をデフォルトとする
    const apiHost = import.meta.env.API_HOST;
    const apiUrl = apiHost ? `http://${apiHost}/api` : "/api"; // VITE_API_HOST があればそれを使い、なければ /api を使う

    console.log(`Fetching data from: ${apiUrl}`); // デバッグ用にURLをログ出力

    fetch(apiUrl) // 修正されたAPI URLを使用
      .then((response) => {
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        return response.json();
      })
      .then((data) => setMessage(data.message))
      .catch((error) => {
        console.error("Error fetching data:", error);
        setMessage("Failed to load message from backend.");
      });
  }, []);

  return (
    <>
      <h1>Vite + React</h1>
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
