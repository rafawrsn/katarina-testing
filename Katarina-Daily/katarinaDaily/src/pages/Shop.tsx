import { useEffect, useState, useMemo } from 'react';
import { useSearchParams } from 'react-router-dom';
import { motion } from 'framer-motion';
import { Search, SlidersHorizontal, X, Heart } from 'lucide-react';
import { ProductGrid } from '@/components/ProductGrid';
import { Reveal } from '@/components/Reveal';
import { getProducts, getCategories, type Product } from '@/services/catalogService';
import { useWishlist } from '@/context/WishlistContext';

export default function Shop() {
  const [searchParams, setSearchParams] = useSearchParams();
  const { wishlist, wishlistCount } = useWishlist();
  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<string[]>([]);
  const [activeCategory, setActiveCategory] = useState<string>('All');
  const [onlyFavorites, setOnlyFavorites] = useState<boolean>(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [sortBy, setSortBy] = useState<'featured' | 'price-low' | 'price-high'>('featured');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const favParam = searchParams.get('favorites') === 'true' || searchParams.get('filter') === 'favorites';
    if (favParam) {
      setOnlyFavorites(true);
    }
  }, [searchParams]);

  useEffect(() => {
    async function loadData() {
      setLoading(true);
      try {
        const [prodList, catList] = await Promise.all([
          getProducts(),
          getCategories(),
        ]);
        setProducts(prodList);
        setCategories(catList);
      } catch (e) {
        console.warn('Failed to load shop data:', e);
      } finally {
        setLoading(false);
      }
    }
    loadData();
  }, []);

  const filteredAndSorted = useMemo(() => {
    let list = [...products];

    // Filter by Favorites
    if (onlyFavorites) {
      const wishlistIdSet = new Set(wishlist.map((id) => String(id).trim()));
      list = list.filter(
        (p) =>
          wishlistIdSet.has(String(p.id)) ||
          (p.slug && wishlistIdSet.has(p.slug))
      );
    }

    // Filter by Category
    if (activeCategory !== 'All') {
      list = list.filter(
        (p) => p.category.toLowerCase() === activeCategory.toLowerCase()
      );
    }

    // Filter by Search
    if (searchQuery.trim()) {
      const q = searchQuery.toLowerCase().trim();
      list = list.filter(
        (p) =>
          p.name.toLowerCase().includes(q) ||
          p.description.toLowerCase().includes(q) ||
          p.category.toLowerCase().includes(q)
      );
    }

    // Sort
    if (sortBy === 'price-low') {
      list.sort((a, b) => (a.rawPrice || 0) - (b.rawPrice || 0));
    } else if (sortBy === 'price-high') {
      list.sort((a, b) => (b.rawPrice || 0) - (a.rawPrice || 0));
    }

    return list;
  }, [products, onlyFavorites, wishlist, activeCategory, searchQuery, sortBy]);

  return (
    <div className="bg-primary pt-32 lg:pt-40 pb-24 lg:pb-32 min-h-screen">
      <div className="container-editorial">
        <Reveal className="mb-12 text-center">
          <p className="text-[11px] uppercase tracking-editorial text-secondary mb-5">Shop</p>
          <h1 className="font-heading text-4xl text-dark sm:text-5xl">The Collection</h1>
          <p className="mx-auto mt-6 max-w-md text-base font-light leading-relaxed text-dark/65">
            Minimalist pieces with special embroidery — designed and self-manufactured for your daily
            style.
          </p>
        </Reveal>

        {/* Search & Filter Bar */}
        <div className="mb-10 flex flex-col md:flex-row items-center justify-between gap-4 border-b border-dark/10 pb-6">
          {/* Category Tabs & Favorites Toggle */}
          <div className="flex flex-wrap items-center justify-center gap-2">
            {categories.map((cat) => (
              <button
                key={cat}
                onClick={() => {
                  setActiveCategory(cat);
                  setOnlyFavorites(false);
                  if (searchParams.has('favorites') || searchParams.has('filter')) {
                    setSearchParams({});
                  }
                }}
                className={`px-5 py-2 text-[11px] uppercase tracking-editorial transition-all duration-300 ${
                  activeCategory === cat && !onlyFavorites
                    ? 'bg-secondary text-primary font-medium'
                    : 'border border-dark/15 text-dark hover:border-secondary hover:text-secondary'
                }`}
              >
                {cat}
              </button>
            ))}

            {/* Favorites filter tab */}
            <button
              onClick={() => {
                setOnlyFavorites((prev) => !prev);
              }}
              className={`flex items-center gap-1.5 px-4 py-2 text-[11px] uppercase tracking-editorial transition-all duration-300 ${
                onlyFavorites
                  ? 'bg-secondary text-primary font-medium'
                  : 'border border-dark/15 text-dark hover:border-secondary hover:text-secondary'
              }`}
            >
              <Heart size={13} className={onlyFavorites || wishlistCount > 0 ? 'fill-current' : ''} />
              Favorites ({wishlistCount})
            </button>
          </div>

          {/* Search input & Sort selector */}
          <div className="flex items-center gap-3 w-full md:w-auto">
            <div className="relative flex-1 md:w-64">
              <Search size={14} className="absolute left-3 top-3 text-dark/40" />
              <input
                type="text"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder="Search pieces..."
                className="w-full bg-white/60 border border-dark/15 pl-9 pr-8 py-2 text-xs text-dark placeholder:text-dark/40 focus:border-secondary focus:outline-none"
              />
              {searchQuery && (
                <button
                  onClick={() => setSearchQuery('')}
                  className="absolute right-2.5 top-2.5 text-dark/40 hover:text-dark"
                >
                  <X size={14} />
                </button>
              )}
            </div>

            <div className="relative shrink-0">
              <select
                value={sortBy}
                onChange={(e) => setSortBy(e.target.value as any)}
                aria-label="Sort products"
                className="appearance-none bg-white/60 border border-dark/15 pl-4 pr-8 py-2 text-xs text-dark focus:border-secondary focus:outline-none cursor-pointer"
              >
                <option value="featured">Sort: Featured</option>
                <option value="price-low">Price: Low to High</option>
                <option value="price-high">Price: High to Low</option>
              </select>
              <SlidersHorizontal
                size={12}
                className="pointer-events-none absolute right-3 top-3 text-dark/50"
              />
            </div>
          </div>
        </div>

        {/* Count summary */}
        <div className="mb-8 flex items-center justify-between text-xs text-dark/60 font-light">
          <p>
            Showing <strong className="text-dark font-medium">{filteredAndSorted.length}</strong> pieces
            {onlyFavorites && <span className="text-secondary font-medium"> in Favorites</span>}
            {activeCategory !== 'All' && <span> in {activeCategory}</span>}
          </p>
        </div>

        {/* Products Grid / Loading / Empty */}
        {loading ? (
          <div className="grid grid-cols-1 gap-x-8 gap-y-14 sm:grid-cols-2 lg:grid-cols-3">
            {[1, 2, 3, 4, 5, 6].map((i) => (
              <div key={i} className="animate-pulse flex flex-col gap-3">
                <div className="aspect-[3/4] bg-dark/5 w-full" />
                <div className="h-3 bg-dark/10 w-1/3 mt-2" />
                <div className="h-4 bg-dark/10 w-3/4" />
                <div className="h-3 bg-dark/5 w-1/2" />
              </div>
            ))}
          </div>
        ) : filteredAndSorted.length === 0 ? (
          <div className="py-20 text-center">
            <p className="font-heading text-2xl text-dark">
              {onlyFavorites ? 'No favorite products found' : 'No products found'}
            </p>
            <p className="mt-2 text-xs font-light text-dark/60">
              {onlyFavorites
                ? 'You have not added any products to your favorites yet, or none match the active filters.'
                : 'Try adjusting your search query or selecting a different category.'}
            </p>
            <button
              onClick={() => {
                setActiveCategory('All');
                setOnlyFavorites(false);
                setSearchQuery('');
                setSearchParams({});
              }}
              className="mt-6 border border-dark/20 px-6 py-2.5 text-[11px] uppercase tracking-editorial text-dark hover:border-secondary hover:text-secondary transition-colors"
            >
              Reset Filters
            </button>
          </div>
        ) : (
          <motion.div
            key={`${activeCategory}-${searchQuery}-${sortBy}`}
            initial={{ opacity: 0, y: 15 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, ease: [0.22, 1, 0.36, 1] }}
          >
            <ProductGrid products={filteredAndSorted} />
          </motion.div>
        )}
      </div>
    </div>
  );
}
