import { useState } from 'react';
import { userApi } from '../api/index.js';
import { Link } from 'react-router-dom';
import React from 'react';

const Login = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const res = await userApi.login(username, password);
      localStorage.setItem('token', res.data.token);
      localStorage.setItem('user', JSON.stringify(res.data.user));
      window.location.href = '/';
    } catch (err: any) {
      setError(err.message || '登录失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex min-h-[60vh] items-center justify-center">
      <div className="w-full max-w-md rounded-2xl border border-[#e8eaed] bg-white p-8 shadow-sm">
        <h2 className="mb-6 text-2xl font-bold text-[#202124]">登录</h2>
        
        {error && (
          <div className="mb-4 rounded-lg bg-red-50 p-3 text-sm text-red-600">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="mb-1 block text-sm font-medium text-[#5f6368]">
              用户名
            </label>
            <input
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="w-full rounded-lg border border-[#dadce0] px-4 py-2.5 outline-none focus:border-[#1a73e8] focus:ring-1 focus:ring-[#1a73e8]"
              placeholder="请输入用户名"
              required
            />
          </div>

          <div>
            <label className="mb-1 block text-sm font-medium text-[#5f6368]">
              密码
            </label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full rounded-lg border border-[#dadce0] px-4 py-2.5 outline-none focus:border-[#1a73e8] focus:ring-1 focus:ring-[#1a73e8]"
              placeholder="请输入密码"
              required
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full rounded-lg bg-[#1a73e8] py-2.5 font-medium text-white transition hover:bg-[#1557b0] disabled:opacity-50"
          >
            {loading ? '登录中...' : '登录'}
          </button>
        </form>

        <p className="mt-4 text-center text-sm text-[#5f6368]">
          还没有账号？{' '}
          <Link to="/register" className="text-[#1a73e8] hover:underline">
            注册
          </Link>
        </p>
      </div>
    </div>
  );
};

export default Login;