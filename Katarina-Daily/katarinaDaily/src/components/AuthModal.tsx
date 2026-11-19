import React, { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { X, Mail, Lock, User as UserIcon, Phone, AlertCircle, ArrowRight } from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useNavigate } from 'react-router-dom';

export function AuthModal() {
  const {
    isAuthModalOpen,
    authModalMode,
    closeAuthModal,
    setAuthModalMode,
    login,
    register,
  } = useAuth();
  
  const navigate = useNavigate();

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [name, setName] = useState('');
  const [phone, setPhone] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  if (!isAuthModalOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSubmitting(true);

    try {
      if (authModalMode === 'login') {
        const user = await login({ email, password });
        if (user.role === 'admin') {
          navigate('/admin');
        }
      } else if (authModalMode === 'register') {
        await register({ name, email, password, phone: phone || undefined });
      } else if (authModalMode === 'forgot_password') {
        // Mock forgot password api call
        await new Promise(resolve => setTimeout(resolve, 1000));
        alert('Password reset link sent to your email.');
        setAuthModalMode('login');
      }
    } catch (err: any) {
      setError(err.message || 'Authentication failed. Please try again.');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <AnimatePresence>
      <div className="fixed inset-0 z-[100] flex items-center justify-center p-4">
        {/* Backdrop */}
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          onClick={closeAuthModal}
          className="fixed inset-0 bg-dark/40 backdrop-blur-sm"
        />

        {/* Modal Window */}
        <motion.div
          initial={{ opacity: 0, scale: 0.95, y: 15 }}
          animate={{ opacity: 1, scale: 1, y: 0 }}
          exit={{ opacity: 0, scale: 0.95, y: 15 }}
          transition={{ duration: 0.4, ease: [0.22, 1, 0.36, 1] }}
          className="relative w-full max-w-md overflow-hidden bg-primary p-8 shadow-2xl border border-dark/10 z-10"
        >
          {/* Close button */}
          <button
            onClick={closeAuthModal}
            className="absolute right-6 top-6 text-dark/50 transition-colors hover:text-dark"
            aria-label="Close modal"
          >
            <X size={20} />
          </button>

          {/* Header */}
          <div className="text-center mb-8">
            <p className="text-[10px] uppercase tracking-editorial text-secondary mb-2">
              Katarina Daily Account
            </p>
            <h3 className="font-heading text-2xl text-dark">
              {authModalMode === 'login' ? 'Welcome Back' 
                : authModalMode === 'forgot_password' ? 'Reset Password' 
                : 'Create Your Account'}
            </h3>
            <p className="mt-2 text-xs font-light text-dark/60">
              {authModalMode === 'login'
                ? 'Sign in to access your orders, saved addresses, and wishlist'
                : authModalMode === 'forgot_password'
                ? 'Enter your email to receive a password reset link'
                : 'Join Katarina Daily for exclusive collections and order tracking'}
            </p>
          </div>

          {authModalMode !== 'forgot_password' && (
            <div className="mb-6 flex border-b border-dark/10">
              <button
                onClick={() => {
                  setAuthModalMode('login');
                  setError(null);
                }}
                className={`flex-1 pb-3 text-xs uppercase tracking-editorial font-medium transition-colors ${
                  authModalMode === 'login'
                    ? 'border-b-2 border-secondary text-secondary'
                    : 'text-dark/40 hover:text-dark'
                }`}
              >
                Sign In
              </button>
              <button
                onClick={() => {
                  setAuthModalMode('register');
                  setError(null);
                }}
                className={`flex-1 pb-3 text-xs uppercase tracking-editorial font-medium transition-colors ${
                  authModalMode === 'register'
                    ? 'border-b-2 border-secondary text-secondary'
                    : 'text-dark/40 hover:text-dark'
                }`}
              >
                Register
              </button>
            </div>
          )}

          {/* Error Banner */}
          {error && (
            <div className="mb-5 flex items-start gap-2.5 bg-red-50 p-3.5 border border-red-200 text-red-700 text-xs font-light">
              <AlertCircle size={16} className="shrink-0 mt-0.5" />
              <span>{error}</span>
            </div>
          )}

          {/* Form */}
          <form onSubmit={handleSubmit} className="space-y-4">
            {authModalMode === 'register' && (
              <>
                <div>
                  <label className="block text-[10px] uppercase tracking-editorial text-dark/70 mb-1.5">
                    Full Name
                  </label>
                  <div className="relative">
                    <UserIcon size={16} className="absolute left-3.5 top-3.5 text-dark/40" />
                    <input
                      type="text"
                      required
                      value={name}
                      onChange={(e) => setName(e.target.value)}
                      placeholder="e.g. Jane Doe"
                      className="w-full bg-white/70 border border-dark/15 pl-10 pr-4 py-2.5 text-xs text-dark placeholder:text-dark/30 focus:border-secondary focus:outline-none transition-colors"
                    />
                  </div>
                </div>

                <div>
                  <label className="block text-[10px] uppercase tracking-editorial text-dark/70 mb-1.5">
                    Phone Number (Optional)
                  </label>
                  <div className="relative">
                    <Phone size={16} className="absolute left-3.5 top-3.5 text-dark/40" />
                    <input
                      type="tel"
                      value={phone}
                      onChange={(e) => setPhone(e.target.value)}
                      placeholder="081234567890"
                      className="w-full bg-white/70 border border-dark/15 pl-10 pr-4 py-2.5 text-xs text-dark placeholder:text-dark/30 focus:border-secondary focus:outline-none transition-colors"
                    />
                  </div>
                </div>
              </>
            )}

            <div>
              <label className="block text-[10px] uppercase tracking-editorial text-dark/70 mb-1.5">
                Email Address
              </label>
              <div className="relative">
                <Mail size={16} className="absolute left-3.5 top-3.5 text-dark/40" />
                <input
                  type="email"
                  required
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="name@example.com"
                  className="w-full bg-white/70 border border-dark/15 pl-10 pr-4 py-2.5 text-xs text-dark placeholder:text-dark/30 focus:border-secondary focus:outline-none transition-colors"
                />
              </div>
            </div>

            {authModalMode !== 'forgot_password' && (
              <div>
                <label className="block text-[10px] uppercase tracking-editorial text-dark/70 mb-1.5">
                  Password
                </label>
                <div className="relative">
                  <Lock size={16} className="absolute left-3.5 top-3.5 text-dark/40" />
                  <input
                    type="password"
                    required
                    minLength={8}
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder="Minimum 8 characters"
                    className="w-full bg-white/70 border border-dark/15 pl-10 pr-4 py-2.5 text-xs text-dark placeholder:text-dark/30 focus:border-secondary focus:outline-none transition-colors"
                  />
                </div>
                {authModalMode === 'login' && (
                  <div className="mt-2 text-right">
                    <button
                      type="button"
                      onClick={() => setAuthModalMode('forgot_password')}
                      className="text-[10px] text-dark/60 hover:text-secondary underline"
                    >
                      Forgot Password?
                    </button>
                  </div>
                )}
              </div>
            )}

            <button
              type="submit"
              disabled={submitting}
              className="mt-6 flex w-full items-center justify-center gap-2 bg-secondary py-3.5 text-xs uppercase tracking-editorial text-primary transition-all duration-300 hover:bg-secondary-600 disabled:opacity-50"
            >
              {submitting ? (
                'Processing...'
              ) : (
                <>
                  {authModalMode === 'login' 
                    ? 'Sign In' 
                    : authModalMode === 'forgot_password' 
                    ? 'Send Reset Link' 
                    : 'Create Account'}
                  <ArrowRight size={14} />
                </>
              )}
            </button>
            {authModalMode === 'forgot_password' && (
              <div className="mt-4 text-center">
                <button
                  type="button"
                  onClick={() => setAuthModalMode('login')}
                  className="text-[10px] text-dark/60 hover:text-secondary uppercase tracking-editorial font-medium"
                >
                  Back to Sign In
                </button>
              </div>
            )}
          </form>
        </motion.div>
      </div>
    </AnimatePresence>
  );
}
