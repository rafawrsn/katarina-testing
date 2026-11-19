import { useEffect, useState } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { motion, AnimatePresence } from 'framer-motion';
import { Menu, X, ShoppingBag, Heart, User as UserIcon, LogOut, Package, Shield } from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useCart } from '@/context/CartContext';
import { useWishlist } from '@/context/WishlistContext';

const navLinks = [
  { label: 'Home', path: '/' },
  { label: 'Shop', path: '/shop' },
  { label: 'About', path: '/about' },
  { label: 'Contact', path: '/contact' },
];

export function Navbar() {
  const [scrolled, setScrolled] = useState(false);
  const [open, setOpen] = useState(false);
  const [userDropdownOpen, setUserDropdownOpen] = useState(false);
  const location = useLocation();

  const { user, isAuthenticated, openAuthModal, logout } = useAuth();
  const { cart, openCart } = useCart();
  const { wishlistCount } = useWishlist();

  useEffect(() => {
    const onScroll = () => setScrolled(window.scrollY > 24);
    onScroll();
    window.addEventListener('scroll', onScroll, { passive: true });
    return () => window.removeEventListener('scroll', onScroll);
  }, []);

  useEffect(() => {
    setOpen(false);
    setUserDropdownOpen(false);
  }, [location.pathname]);

  return (
    <header
      className={`fixed inset-x-0 top-0 z-50 transition-all duration-700 ease-elegant ${
        scrolled
          ? 'bg-primary/90 backdrop-blur-md border-b border-dark/5 py-4'
          : 'bg-transparent py-6'
      }`}
    >
      <nav className="container-editorial flex items-center justify-between">
        {/* Brand */}
        <Link
          to="/"
          className="font-heading text-2xl tracking-wide text-dark"
          aria-label="Katarina Daily home"
        >
          Katarina Daily
        </Link>

        {/* Center navigation links */}
        <ul className="hidden md:flex items-center gap-10">
          {navLinks.map((link) => {
            const active = location.pathname === link.path;
            return (
              <li key={link.path}>
                <Link
                  to={link.path}
                  className={`relative text-[11px] uppercase tracking-editorial transition-colors duration-500 ${
                    active ? 'text-secondary' : 'text-dark hover:text-secondary'
                  }`}
                >
                  {link.label}
                  <span
                    className={`absolute -bottom-1 left-0 h-px bg-secondary transition-all duration-500 ease-elegant ${
                      active ? 'w-full' : 'w-0'
                    }`}
                  />
                </Link>
              </li>
            );
          })}
        </ul>

        {/* Right action icons */}
        <div className="flex items-center gap-5">
          {/* Wishlist / Favorites indicator */}
          <Link
            to="/favorites"
            className={`relative hidden sm:flex transition-colors ${
              location.pathname === '/favorites' || location.pathname === '/wishlist'
                ? 'text-secondary'
                : 'text-dark/80 hover:text-secondary'
            }`}
            title="Favorites"
            aria-label="Favorites"
          >
            <Heart size={19} className={wishlistCount > 0 ? 'fill-secondary/20' : ''} />
            {wishlistCount > 0 && (
              <span className="absolute -top-1.5 -right-2 flex h-4 min-w-4 items-center justify-center rounded-full bg-secondary px-1 text-[9px] font-medium text-primary">
                {wishlistCount}
              </span>
            )}
          </Link>

          {/* User Account / Auth */}
          <div className="relative">
            {isAuthenticated ? (
              <div className="relative">
                <button
                  onClick={() => setUserDropdownOpen((v) => !v)}
                  className="flex items-center gap-2 text-xs font-light text-dark hover:text-secondary transition-colors"
                  aria-label="User menu"
                >
                  <div className="flex h-7 w-7 items-center justify-center rounded-full bg-secondary/15 text-secondary text-xs font-medium uppercase">
                    {user?.name ? user.name[0] : 'U'}
                  </div>
                  <span className="hidden lg:inline text-[11px] tracking-editorial uppercase truncate max-w-[100px]">
                    {user?.name?.split(' ')[0]}
                  </span>
                </button>

                {/* Dropdown Menu */}
                <AnimatePresence>
                  {userDropdownOpen && (
                    <motion.div
                      initial={{ opacity: 0, y: 10 }}
                      animate={{ opacity: 1, y: 0 }}
                      exit={{ opacity: 0, y: 10 }}
                      transition={{ duration: 0.2 }}
                      className="absolute right-0 mt-2 w-48 bg-primary border border-dark/10 shadow-xl py-2 z-50"
                    >
                      <div className="px-4 py-2 border-b border-dark/5">
                        <p className="text-xs font-medium text-dark truncate">{user?.name}</p>
                        <p className="text-[10px] text-dark/50 truncate">{user?.email}</p>
                      </div>
                      {user?.role === 'admin' && (
                        <Link
                          to="/admin"
                          onClick={() => setUserDropdownOpen(false)}
                          className="flex items-center gap-2 px-4 py-2.5 text-xs text-secondary font-medium hover:bg-secondary/10 transition-colors"
                        >
                          <Shield size={14} /> Admin Dashboard
                        </Link>
                      )}
                      <Link
                        to="/orders"
                        onClick={() => setUserDropdownOpen(false)}
                        className="flex items-center gap-2 px-4 py-2.5 text-xs text-dark hover:bg-secondary/10 hover:text-secondary transition-colors"
                      >
                        <Package size={14} /> My Orders
                      </Link>
                      <button
                        onClick={() => {
                          logout();
                          setUserDropdownOpen(false);
                        }}
                        className="flex w-full items-center gap-2 px-4 py-2.5 text-xs text-red-600 hover:bg-red-50 transition-colors"
                      >
                        <LogOut size={14} /> Sign Out
                      </button>
                    </motion.div>
                  )}
                </AnimatePresence>
              </div>
            ) : (
              <button
                onClick={() => openAuthModal('login')}
                className="hidden sm:flex items-center gap-1.5 text-[11px] uppercase tracking-editorial text-dark/80 hover:text-secondary transition-colors"
              >
                <UserIcon size={18} />
                <span className="hidden lg:inline">Sign In</span>
              </button>
            )}
          </div>

          {/* Cart Bag trigger with live count */}
          <button
            onClick={openCart}
            className="relative flex items-center text-dark/90 hover:text-secondary transition-colors"
            aria-label="Open cart bag"
          >
            <ShoppingBag size={20} />
            {cart.total_items > 0 && (
              <span className="absolute -top-1.5 -right-2 flex h-4 min-w-4 items-center justify-center rounded-full bg-secondary px-1 text-[9px] font-medium text-primary">
                {cart.total_items}
              </span>
            )}
          </button>

          {/* Mobile hamburger */}
          <button
            className="md:hidden text-dark ml-1"
            onClick={() => setOpen((v) => !v)}
            aria-label={open ? 'Close menu' : 'Open menu'}
            aria-expanded={open}
          >
            {open ? <X size={22} /> : <Menu size={22} />}
          </button>
        </div>
      </nav>

      {/* Mobile Drawer */}
      <AnimatePresence>
        {open && (
          <motion.div
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: 'auto' }}
            exit={{ opacity: 0, height: 0 }}
            transition={{ duration: 0.4, ease: [0.22, 1, 0.36, 1] }}
            className="md:hidden overflow-hidden bg-primary/98 backdrop-blur-md border-t border-dark/5"
          >
            <ul className="container-editorial flex flex-col gap-5 py-6">
              {navLinks.map((link) => {
                const active = location.pathname === link.path;
                return (
                  <li key={link.path}>
                    <Link
                      to={link.path}
                      className={`text-xs uppercase tracking-editorial ${
                        active ? 'text-secondary font-medium' : 'text-dark'
                      }`}
                    >
                      {link.label}
                    </Link>
                  </li>
                );
              })}

              <li>
                <Link
                  to="/favorites"
                  className={`flex items-center justify-between text-xs uppercase tracking-editorial ${
                    location.pathname === '/favorites' || location.pathname === '/wishlist'
                      ? 'text-secondary font-medium'
                      : 'text-dark'
                  }`}
                >
                  <span className="flex items-center gap-2">
                    <Heart size={14} className={wishlistCount > 0 ? 'fill-secondary' : ''} />
                    My Favorites
                  </span>
                  {wishlistCount > 0 && (
                    <span className="rounded-full bg-secondary/15 px-2 py-0.5 text-[10px] font-medium text-secondary">
                      {wishlistCount}
                    </span>
                  )}
                </Link>
              </li>

              <li className="pt-4 border-t border-dark/10 flex items-center justify-between">
                {isAuthenticated ? (
                  <div className="flex flex-col w-full gap-2">
                    <div className="flex w-full items-center justify-between">
                      <Link
                        to="/orders"
                        onClick={() => setOpen(false)}
                        className="text-xs uppercase tracking-editorial text-secondary"
                      >
                        My Orders ({user?.name})
                      </Link>
                      <button
                        onClick={logout}
                        className="text-xs uppercase tracking-editorial text-red-600"
                      >
                        Sign Out
                      </button>
                    </div>
                    {user?.role === 'admin' && (
                      <Link
                        to="/admin"
                        onClick={() => setOpen(false)}
                        className="text-xs uppercase tracking-editorial text-secondary font-medium flex items-center gap-1.5"
                      >
                        <Shield size={13} /> Admin Dashboard
                      </Link>
                    )}
                  </div>
                ) : (
                  <button
                    onClick={() => {
                      setOpen(false);
                      openAuthModal('login');
                    }}
                    className="text-xs uppercase tracking-editorial text-secondary"
                  >
                    Sign In / Register
                  </button>
                )}
              </li>
            </ul>
          </motion.div>
        )}
      </AnimatePresence>
    </header>
  );
}
