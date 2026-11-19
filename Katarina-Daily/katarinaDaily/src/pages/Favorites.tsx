import { useEffect, useState, useMemo } from 'react';
import { Link } from 'react-router-dom';
import { motion } from 'framer-motion';
import { Heart, Search, SlidersHorizontal, ArrowLeft, X } from 'lucide-react';
import { ProductGrid } from '@/components/ProductGrid';
import { Reveal } from '@/components/Reveal';
import { useWishlist } from '@/context/WishlistContext';
import { getWishlistProducts } from '@/services/wishlistService';
import { getProducts, type Product } from '@/services/catalogService';

export default function Favorites() {
  const { wishlist, wishlistCount } = useWishlist();
  const [favoriteProducts, setFavoriteProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState('');
  const [activeCategory, setActiveCategory] = useState('All');
  const [sortBy, setSortBy] = useState<'featured' | 'price-low' | 'price-high' | 'name-az'>('featured');

  useEffect(() => {
    let isMounted = true;

    async function loadFavorites() {
      setLoading(true);
      try {
        if (wishlist.length === 0) {
          if (isMounted) {
            setFavoriteProducts([]);
            setLoading(false);
          }
          return;
        }

        // Load full product catalog to ensure all wishlist items are matched regardless of source
        const [wishlistProds, allProds] = await Promise.all([
          getWishlistProducts().catch(() => [] as Product[]),
          getProducts().catch(() => [] as Product[]),
        ]);

        const combined = new Map<string, Product>();
        for (const p of wishlistProds) {
          combined.set(String(p.id), p);
        }
        for (const p of allProds) {
          if (!combined.has(String(p.id))) {
            combined.set(String(p.id), p);
          }
        }

        const validProds: Product[] = [];
        const wishlistIdSet = new Set(wishlist.map((id) => String(id).trim()));

        for (const [id, prod] of combined.entries()) {
          if (
            wishlistIdSet.has(id) ||
            (prod.slug && wishlistIdSet.has(prod.slug))
          ) {
            validProds.push(prod);
          }
        }

        if (isMounted) {
          setFavoriteProducts(validProds);
        }
      } catch (err) {
        console.warn('Failed to load favorites:', err);
        if (isMounted) {
          setFavoriteProducts([]);
        }
      } finally {
        if (isMounted) {
          setLoading(false);
        }
      }
    }

    loadFavorites();

    return () => {
      isMounted = false;
    };
  }, [wishlist]);

  // Derive unique categories present in favorites
  const categories = useMemo(() => {
    const set = new Set<string>();
    favoriteProducts.forEach((p) => {
      if (p.category) set.add(p.category);
    });
    return ['All', ...Array.from(set)];
  }, [favoriteProducts]);

  // Filter and sort favorites
  const filteredAndSorted = useMemo(() => {
    const wishlistIdSet = new Set(wishlist.map((id) => String(id).trim()));
    let list = favoriteProducts.filter(
      (p) =>
        wishlistIdSet.has(String(p.id)) ||
        (p.slug && wishlistIdSet.has(p.slug))
    );

    // Category filter
    if (activeCategory !== 'All') {
      list = list.filter((p) => p.category.toLowerCase() === activeCategory.toLowerCase());
    }

    // Search query filter
    if (searchQuery.trim()) {
      const q = searchQuery.toLowerCase().trim();
      list = list.filter(
        (p) =>
          p.name.toLowerCase().includes(q) ||
          p.description.toLowerCase().includes(q) ||
          p.category.toLowerCase().includes(q)
      );
    }

    // Sorting
    if (sortBy === 'price-low') {
      list.sort((a, b) => (a.rawPrice || 0) - (b.rawPrice || 0));
    } else if (sortBy === 'price-high') {
      list.sort((a, b) => (b.rawPrice || 0) - (a.rawPrice || 0));
    } else if (sortBy === 'name-az') {
      list.sort((a, b) => a.name.localeCompare(b.name));
    }

    return list;
  }, [favoriteProducts, wishlist, activeCategory, searchQuery, sortBy]);

  return (
    <div className="bg-primary pt-32 lg:pt-40 pb-24 lg:pb-32 min-h-screen">
      <div className="container-editorial">
        <Link
          to="/shop"
          className="mb-8 inline-flex items-center gap-2 text-[11px] uppercase tracking-editorial text-dark/60 transition-colors hover:text-secondary"
        >
          <ArrowLeft size={14} /> Back to shop
        </Link>

        <Reveal className="mb-10 text-center">
          <div className="mx-auto mb-3 flex h-12 w-12 items-center justify-center rounded-full bg-secondary/10 text-secondary">
            <Heart size={22} className="fill-secondary/20" />
          </div>
          <p className="text-[11px] uppercase tracking-editorial text-secondary mb-2">Saved Pieces</p>
          <h1 className="font-heading text-4xl text-dark sm:text-5xl">My Favorites</h1>
          <p className="mx-auto mt-4 max-w-md text-sm font-light leading-relaxed text-dark/65">
            Your curated collection of loved pieces. Keep track of what you want to wear next.
          </p>
        </Reveal>

        {/* If empty favorites list */}
        {!loading && (wishlistCount === 0 || favoriteProducts.length === 0) ? (
          <div className="py-20 text-center bg-white/50 border border-dark/10 max-w-xl mx-auto p-12">
            <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-dark/5 text-dark/40">
              <Heart size={28} className="stroke-1" />
            </div>
            <h2 className="font-heading text-2xl text-dark">Your favorites list is empty</h2>
            <p className="mx-auto mt-3 max-w-sm text-xs font-light leading-relaxed text-dark/60">
              Save your favorite items by tapping the heart icon on any piece you like while browsing our collection.
            </p>
            <Link
              to="/shop"
              className="mt-8 inline-flex items-center gap-2 bg-secondary px-8 py-3.5 text-xs uppercase tracking-editorial text-primary hover:bg-secondary-600 transition-colors"
            >
              Explore the Collection
            </Link>
          </div>
        ) : (
          <>
            {/* Search, Category & Sort Bar */}
            <div className="mb-10 flex flex-col md:flex-row items-center justify-between gap-4 border-b border-dark/10 pb-6">
              {/* Category Tabs if multiple categories exist */}
              <div className="flex flex-wrap items-center justify-center gap-2">
                {categories.map((cat) => (
                  <button
                    key={cat}
                    onClick={() => setActiveCategory(cat)}
                    className={`px-4 py-1.5 text-[11px] uppercase tracking-editorial transition-all duration-300 ${
                      activeCategory === cat
                        ? 'bg-secondary text-primary font-medium'
                        : 'border border-dark/15 text-dark hover:border-secondary hover:text-secondary'
                    }`}
                  >
                    {cat}
                  </button>
                ))}
              </div>

              {/* Search & Sort */}
              <div className="flex items-center gap-3 w-full md:w-auto">
                <div className="relative flex-1 md:w-60">
                  <Search size={14} className="absolute left-3 top-3 text-dark/40" />
                  <input
                    type="text"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    placeholder="Search favorites..."
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
                    aria-label="Sort favorites"
                    className="appearance-none bg-white/60 border border-dark/15 pl-4 pr-8 py-2 text-xs text-dark focus:border-secondary focus:outline-none cursor-pointer"
                  >
                    <option value="featured">Sort: Featured</option>
                    <option value="price-low">Price: Low to High</option>
                    <option value="price-high">Price: High to Low</option>
                    <option value="name-az">Name: A to Z</option>
                  </select>
                  <SlidersHorizontal
                    size={12}
                    className="pointer-events-none absolute right-3 top-3 text-dark/50"
                  />
                </div>
              </div>
            </div>

            {/* Summary info */}
            <div className="mb-8 flex items-center justify-between text-xs text-dark/60 font-light">
              <p>
                Showing <strong className="text-dark font-medium">{filteredAndSorted.length}</strong> favorited {filteredAndSorted.length === 1 ? 'piece' : 'pieces'}
                {activeCategory !== 'All' && <span> in {activeCategory}</span>}
              </p>
            </div>

            {/* Content / Loading / Empty filtered state */}
            {loading ? (
              <div className="grid grid-cols-1 gap-x-8 gap-y-14 sm:grid-cols-2 lg:grid-cols-3">
                {[1, 2, 3].map((i) => (
                  <div key={i} className="animate-pulse flex flex-col gap-3">
                    <div className="aspect-[3/4] bg-dark/5 w-full" />
                    <div className="h-3 bg-dark/10 w-1/3 mt-2" />
                    <div className="h-4 bg-dark/10 w-3/4" />
                    <div className="h-3 bg-dark/5 w-1/2" />
                  </div>
                ))}
              </div>
            ) : filteredAndSorted.length === 0 ? (
              <div className="py-16 text-center bg-white/40 border border-dark/10 p-8">
                <p className="font-heading text-xl text-dark">No favorites match your search</p>
                <p className="mt-2 text-xs font-light text-dark/60">
                  Try adjusting your search terms or selecting a different category.
                </p>
                <button
                  onClick={() => {
                    setActiveCategory('All');
                    setSearchQuery('');
                  }}
                  className="mt-6 border border-dark/20 px-6 py-2.5 text-[11px] uppercase tracking-editorial text-dark hover:border-secondary hover:text-secondary transition-colors"
                >
                  Reset Filters
                </button>
              </div>
            ) : (
              <motion.div
                key={`${activeCategory}-${searchQuery}-${sortBy}-${filteredAndSorted.length}`}
                initial={{ opacity: 0, y: 15 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.5, ease: [0.22, 1, 0.36, 1] }}
              >
                <ProductGrid products={filteredAndSorted} />
              </motion.div>
            )}
          </>
        )}
      </div>
    </div>
  );
}
