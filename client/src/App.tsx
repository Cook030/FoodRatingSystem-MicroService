// src/App.tsx
import { BrowserRouter as Router, Routes, Route, Link, useNavigate } from 'react-router-dom';
import Home from './pages/Home';
import Detail from './pages/Detail';
import Login from './pages/Login';
import Register from './pages/Register';
import React, { useState, useEffect } from 'react';

function Header() {
  const [user, setUser] = useState<{ user_name: string } | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    const savedUser = localStorage.getItem('user');
    if (savedUser) {
      setUser(JSON.parse(savedUser));
    }
  }, []);

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    setUser(null);
    navigate('/');
  };

  return (
    <header className="sticky top-0 z-20 border-b border-[#e8eaed] bg-white/90 backdrop-blur">
      <div className="mx-auto flex h-16 w-full max-w-6xl items-center justify-between gap-3 px-4 sm:px-6">
        <Link to="/" className="flex items-center gap-3">
          <div className="grid h-8 w-8 place-items-center rounded-full bg-[#e8f0fe] text-sm text-[#1a73e8]">G</div>
          <h1 className="text-lg font-semibold tracking-tight text-[#202124]">美食评分系统</h1>
        </Link>
        <div className="flex items-center gap-4">
          {user ? (
            <div className="flex items-center gap-3">
              <span className="text-sm text-[#5f6368]">{user.user_name}</span>
              <button
                onClick={handleLogout}
                className="rounded-lg border border-[#dadce0] px-3 py-1.5 text-sm text-[#5f6368] transition hover:bg-[#f1f3f4]"
              >
                退出
              </button>
            </div>
          ) : (
            <Link
              to="/login"
              className="rounded-lg bg-[#1a73e8] px-4 py-2 text-sm font-medium text-white transition hover:bg-[#1557b0]"
            >
              登录
            </Link>
          )}
        </div>
      </div>
    </header>
  );
}

function App() {
  return (
    <Router
      future={{
        v7_startTransition: true,
        v7_relativeSplatPath: true,
      }}
    >
      <div className="min-h-screen bg-[#f6f8fc] text-[#202124]">
        <Header />

        <main className="mx-auto w-full max-w-6xl px-4 py-6 sm:px-6 sm:py-8">
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/restaurant/:id" element={<Detail />} />
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />
          </Routes>
        </main>
      </div>
    </Router>
  );
}

export default App;