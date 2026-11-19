import { Link } from 'react-router-dom';
import { motion } from 'framer-motion';
import { Heart, Plus } from 'lucide-react';
import type { Product } from '@/services/catalogService';
import { useWishlist } from '@/context/WishlistContext';
import { useCart } from '@/context/CartContext';

interface ProductCardProps {
  product: Product;
}

export function ProductCard({ product }: ProductCardProps) {
  const { isInWishlist, toggleWishlist } = useWishlist();
  const { addToCart } = useCart();
  const isWished = isInWishlist(product.id);

  const handleQuickAdd = async (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

    const firstVariant = product.variants?.[0];
    const variantId = firstVariant ? firstVariant.id : `variant-${product.id}`;

    await addToCart(variantId, 1, {
      productName: product.name,
      productImage: product.image,
      variant: firstVariant,
      unitPrice: product.rawPrice || 129000,
    });
  };

  const handleWishlistClick = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    toggleWishlist(product.id);
  };

  return (
    <motion.article
      initial={{ opacity: 0, y: 30 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true, amount: 0.15 }}
      transition={{ duration: 0.8, ease: [0.22, 1, 0.36, 1] }}
      className="group flex flex-col justify-between"
    >
      <div>
        <div className="relative overflow-hidden bg-primary-100">
          <Link to={`/product/${product.slug || product.id}`} className="block">
            <img
              src={product.image}
              alt={product.name}
              className="aspect-[3/4] w-full object-cover transition-transform duration-[1.2s] ease-elegant group-hover:scale-105"
              loading="lazy"
            />
            <div className="absolute inset-0 bg-dark/0 transition-colors duration-500 group-hover:bg-dark/5" />
          </Link>

          {/* Wishlist Button */}
          <button
            onClick={handleWishlistClick}
            className={`absolute top-3.5 right-3.5 flex h-8 w-8 items-center justify-center rounded-full backdrop-blur-md transition-all duration-300 ${
              isWished
                ? 'bg-secondary text-primary'
                : 'bg-primary/80 text-dark/70 hover:bg-primary hover:text-secondary'
            }`}
            aria-label={isWished ? 'Remove from wishlist' : 'Add to wishlist'}
          >
            <Heart size={15} className={isWished ? 'fill-primary' : ''} />
          </button>

          {/* Quick Add Overlay on Hover */}
          <div className="absolute inset-x-0 bottom-0 p-3 translate-y-full opacity-0 transition-all duration-300 ease-elegant group-hover:translate-y-0 group-hover:opacity-100">
            <button
              onClick={handleQuickAdd}
              className="flex w-full items-center justify-center gap-2 bg-primary/95 py-2.5 text-[10px] uppercase tracking-editorial text-dark shadow-sm backdrop-blur-md hover:bg-secondary hover:text-primary transition-colors"
            >
              <Plus size={13} /> Quick Add to Bag
            </button>
          </div>
        </div>

        <div className="mt-4 flex flex-col gap-1">
          <p className="text-[10px] uppercase tracking-editorial text-secondary">
            {product.category}
          </p>
          <Link to={`/product/${product.slug || product.id}`}>
            <h3 className="font-heading text-lg tracking-wide text-dark transition-colors duration-500 hover:text-secondary line-clamp-1">
              {product.name}
            </h3>
          </Link>
          <p className="text-xs font-light text-dark/60 line-clamp-2">{product.description}</p>
        </div>
      </div>

      <div className="mt-3 flex items-center justify-between pt-2 border-t border-dark/5">
        <p className="text-sm font-normal text-dark">{product.price}</p>
        <Link
          to={`/product/${product.slug || product.id}`}
          className="text-[10px] uppercase tracking-editorial text-secondary hover:underline"
        >
          View Details →
        </Link>
      </div>
    </motion.article>
  );
}
