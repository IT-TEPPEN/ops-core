import { useState, useEffect } from "react";
import { Routes, Route, Link } from "react-router-dom"; // Import routing components
import BlogPage from "./pages/BlogPage"; // Import the new BlogPage

// Define a simple Home component for the root path
function HomePage() {
  const [count, setCount] = useState(0);
  const [message, setMessage] = useState("Loading...");

  useEffect(() => {
    const apiHost = import.meta.env.VITE_API_HOST;
    const apiUrl = apiHost ? `${apiHost}/api` : "/api";

    console.log(`Fetching data from: ${apiUrl}`);

    fetch(apiUrl)
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
    <div className="text-center p-8">
      {/* Basic Tailwind styling for the home page content */}
      <h1 className="text-3xl font-bold mb-4">Vite + React</h1>
      <p className="mb-4">
        Message from backend: <span className="font-semibold">{message}</span>
      </p>
      <div className="card bg-gray-100 dark:bg-gray-700 p-6 rounded-lg shadow">
        <button
          onClick={() => setCount((count) => count + 1)}
          className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded mb-4 transition duration-150 ease-in-out"
        >
          count is {count}
        </button>
        <p className="text-sm text-gray-600 dark:text-gray-400">
          Edit <code>src/App.tsx</code> and save to test HMR
        </p>
      </div>
      <p className="mt-4 text-sm text-gray-500 dark:text-gray-400">
        Click on the Vite and React logos to learn more
      </p>
    </div>
  );
}

function App() {
  return (
    // Apply base layout and background colors using Tailwind
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 text-gray-900 dark:text-gray-100">
      {/* Navigation Bar */}
      <nav className="bg-white dark:bg-gray-800 shadow-md sticky top-0 z-10">
        <div className="max-w-5xl mx-auto px-4">
          <div className="flex justify-center items-center h-16">
            <ul className="flex space-x-6">
              <li>
                <Link
                  to="/"
                  className="text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 px-3 py-2 rounded-md text-sm font-medium transition duration-150 ease-in-out"
                >
                  Home
                </Link>
              </li>
              <li>
                <Link
                  to="/blog"
                  className="text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 px-3 py-2 rounded-md text-sm font-medium transition duration-150 ease-in-out"
                >
                  Blog
                </Link>
              </li>
            </ul>
          </div>{" "}
          {/* Close inner div */}
        </div>{" "}
        {/* Close outer div */}
      </nav>{" "}
      {/* Close nav */}
      {/* Page Content Area */}
      <main className="max-w-5xl mx-auto p-4">
        {/* Ensure Routes, Route, HomePage, BlogPage are used */}
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/blog" element={<BlogPage />} />
        </Routes>
      </main>{" "}
      {/* Close main */}
    </div> // Close outer div
  );
}

export default App;
