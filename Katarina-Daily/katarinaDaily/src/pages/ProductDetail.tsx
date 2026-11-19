import { useEffect, useState, useMemo } from 'react';
import { useParams, Link } from 'react-router-dom';
import { motion } from 'framer-motion';
import { ShoppingBag, ArrowLeft, Heart, Plus, Minus, Check, ShieldCheck, Truck } from 'lucide-react';
import { getProductById, formatIDR, type Product, type ProductVariant } from '@/services/catalogService';
import { useCart } from '@/context/CartContext';
import { useWishlist } from '@/context/WishlistContext';

export default function ProductDetail() {
  const { id } = useParams<{ id: string }>();
  const [product, setProduct] = useState<Product | undefined>(undefined);
  const [loading, setLoading] = useState(true);
  const [selectedImage, setSelectedImage] = useState<string>('');
  const [selectedColor, setSelectedColor] = useState<string>('');
  const [selectedSize, setSelectedSize] = useState<string>('');
  const [quantity, setQuantity] = useState(1);
  const [addingToCart, setAddingToCart] = useState(false);

  const { addToCart } = useCart();
  const { isInWishlist, toggleWishlist } = useWishlist();

  useEffect(() => {
    setLoading(true);
    getProductById(id ?? '').then((p) => {
      setProduct(p);
      if (p) {
        setSelectedImage(p.image);
        if (p.variants && p.variants.length > 0) {
          setSelectedColor(p.variants[0].color);
          setSelectedSize(p.variants[0].size);
        }
      }
      setLoading(false);
    });
  }, [id]);

  // Unique Colors & Sizes available
  const availableColors = useMemo(() => {
    if (!product?.variants) return [];
    return Array.from(new Set(product.variants.map((v) => v.color)));
  }, [product]);

  const availableSizes = useMemo(() => {
    if (!product?.variants) return [];
    const colorVariants = selectedColor
      ? product.variants.filter((v) => v.color === selectedColor)
      : product.variants;
    return Array.from(new Set(colorVariants.map((v) => v.size)));
  }, [product, selectedColor]);

  // Selected Variant based on Color + Size
  const currentVariant = useMemo<ProductVariant | undefined>(() => {
    if (!product?.variants || product.variants.length === 0) return undefined;
    return (
      product.variants.find(
        (v) => v.color === selectedColor && v.size === selectedSize
      ) ||
      product.variants.find((v) => v.color === selectedColor) ||
      product.variants[0]
    );
  }, [product, selectedColor, selectedSize]);

  // Handle color change and keep a valid size
  const handleColorSelect = (color: string) => {
    setSelectedColor(color);
    if (product?.variants) {
      const validSizes = product.variants
        .filter((v) => v.color === color)
        .map((v) => v.size);
      if (!validSizes.includes(selectedSize) && validSizes.length > 0) {
        setSelectedSize(validSizes[0]);
      }
    }
  };

  const handleAddToCart = async () => {
    if (!product) return;
    setAddingToCart(true);

    const variantId = currentVariant
      ? currentVariant.id
      : `variant-${product.id}-${selectedColor}-${selectedSize}`;

    const unitPrice =
      currentVariant && typeof currentVariant.price === 'number'
        ? currentVariant.price
        : product.rawPrice || 129000;

    await addToCart(variantId, quantity, {
      productName: product.name,
      productImage: selectedImage || product.image,
      variant: currentVariant,
      unitPrice,
    });

    setAddingToCart(false);
  };

  if (loading) {
    return (
      <div className="flex min-h-[70vh] items-center justify-center bg-primary pt-24">
        <div className="flex flex-col items-center gap-3">
          <div className="h-6 w-6 animate-spin rounded-full border-2 border-secondary border-t-transparent" />
          <p className="text-xs uppercase tracking-editorial text-dark/60">Loading piece details...</p>
        </div>
      </div>
    );
  }

  if (!product) {
    return (
      <div className="flex min-h-[60vh] flex-col items-center justify-center gap-6 bg-primary pt-24">
        <p className="font-heading text-2xl text-dark">Product not found</p>
        <Link
          to="/shop"
          className="text-[11px] uppercase tracking-editorial text-secondary transition-colors hover:text-secondary-700"
        >
          Back to shop
        </Link>
      </div>
    );
  }

  const isWished = isInWishlist(product.id);
  const currentPrice = currentVariant ? formatIDR(currentVariant.price) : product.price;
  const inStock = (currentVariant?.stock ?? product.total_stock ?? 10) > 0;

  return (
    <div className="bg-primary pt-28 lg:pt-36 pb-24 lg:pb-32">
      <div className="container-editorial">
        <Link
          to="/shop"
          className="mb-8 inline-flex items-center gap-2 text-[11px] uppercase tracking-editorial text-dark/60 transition-colors hover:text-secondary"
        >
          <ArrowLeft size={14} /> Back to shop
        </Link>

        <div className="grid grid-cols-1 gap-12 lg:grid-cols-2 lg:gap-16">
          {/* Product Media Gallery */}
          <motion.div
            initial={{ opacity: 0, scale: 1.02 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.8, ease: [0.22, 1, 0.36, 1] }}
            className="flex flex-col gap-4"
          >
            <div className="overflow-hidden bg-primary-100 border border-dark/5">
              <img
                src={selectedImage || product.image}
                alt={product.name}
                className="aspect-[3/4] w-full object-cover transition-transform duration-700"
              />
            </div>

            {/* Thumbnail selector if multiple images exist */}
            {product.images && product.images.length > 1 && (
              <div className="flex gap-3 overflow-x-auto pb-2">
                {product.images.map((img, idx) => (
                  <button
                    key={idx}
                    onClick={() => setSelectedImage(img.image_url)}
                    className={`h-20 w-16 shrink-0 overflow-hidden border transition-all ${
                      selectedImage === img.image_url
                        ? 'border-secondary ring-1 ring-secondary'
                        : 'border-dark/10 opacity-70 hover:opacity-100'
                    }`}
                  >
                    <img
                      src={img.image_url}
                      alt={`${product.name} preview ${idx + 1}`}
                      className="h-full w-full object-cover"
                    />
                  </button>
                ))}
              </div>
            )}
          </motion.div>

          {/* Product Purchase Actions & Info */}
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8, ease: [0.22, 1, 0.36, 1], delay: 0.1 }}
            className="flex flex-col justify-start"
          >
            <div className="flex items-center justify-between">
              <p className="text-[11px] uppercase tracking-editorial text-secondary">
                {product.category}
              </p>
              <button
                onClick={() => toggleWishlist(product.id)}
                className={`flex items-center gap-1.5 text-xs uppercase tracking-editorial transition-colors ${
                  isWished ? 'text-secondary font-medium' : 'text-dark/50 hover:text-secondary'
                }`}
              >
                <Heart size={15} className={isWished ? 'fill-secondary' : ''} />
                <span>{isWished ? 'Saved to Wishlist' : 'Add to Wishlist'}</span>
              </button>
            </div>

            <h1 className="font-heading text-3xl leading-snug text-dark sm:text-4xl mt-3">
              {product.name}
            </h1>

            {/* Price & Stock status */}
            <div className="mt-4 flex items-baseline gap-4">
              <p className="text-2xl font-light text-dark">{currentPrice}</p>
              {inStock ? (
                <span className="flex items-center gap-1 text-[11px] text-emerald-700 font-light bg-emerald-50 px-2 py-0.5 border border-emerald-200">
                  <Check size={12} /> In Stock ({currentVariant?.stock ?? product.total_stock ?? 'Ready'})
                </span>
              ) : (
                <span className="text-[11px] text-red-600 bg-red-50 px-2 py-0.5 border border-red-200">
                  Out of Stock
                </span>
              )}
            </div>

            <p className="mt-6 text-sm font-light leading-relaxed text-dark/70">
              {product.description}
            </p>

            {/* Color Selector */}
            {availableColors.length > 0 && (
              <div className="mt-8">
                <label className="block text-[11px] uppercase tracking-editorial text-dark/80 mb-3">
                  Color: <span className="font-medium text-dark">{selectedColor}</span>
                </label>
                <div className="flex flex-wrap gap-2.5">
                  {availableColors.map((color) => (
                    <button
                      key={color}
                      onClick={() => handleColorSelect(color)}
                      className={`px-4 py-2 text-xs uppercase tracking-editorial transition-all ${
                        selectedColor === color
                          ? 'border-2 border-secondary bg-secondary/10 text-secondary font-medium'
                          : 'border border-dark/15 text-dark hover:border-dark/40'
                      }`}
                    >
                      {color}
                    </button>
                  ))}
                </div>
              </div>
            )}

            {/* Size Selector */}
            {availableSizes.length > 0 && (
              <div className="mt-6">
                <label className="block text-[11px] uppercase tracking-editorial text-dark/80 mb-3">
                  Size: <span className="font-medium text-dark">{selectedSize}</span>
                </label>
                <div className="flex flex-wrap gap-2.5">
                  {availableSizes.map((size) => (
                    <button
                      key={size}
                      onClick={() => setSelectedSize(size)}
                      className={`min-w-[48px] px-3.5 py-2 text-xs uppercase tracking-editorial text-center transition-all ${
                        selectedSize === size
                          ? 'border-2 border-secondary bg-secondary/10 text-secondary font-medium'
                          : 'border border-dark/15 text-dark hover:border-dark/40'
                      }`}
                    >
                      {size}
                    </button>
                  ))}
                </div>
              </div>
            )}

            {/* Quantity Selector & Add to Bag Buttons */}
            <div className="mt-8 flex flex-col sm:flex-row items-stretch sm:items-center gap-4">
              <div className="flex items-center border border-dark/20 h-12 w-fit">
                <button
                  onClick={() => setQuantity((q) => Math.max(1, q - 1))}
                  className="p-3 text-dark/60 hover:text-dark transition-colors"
                  aria-label="Decrease quantity"
                >
                  <Minus size={14} />
                </button>
                <span className="px-4 text-sm font-light text-dark">{quantity}</span>
                <button
                  onClick={() => setQuantity((q) => q + 1)}
                  className="p-3 text-dark/60 hover:text-dark transition-colors"
                  aria-label="Increase quantity"
                >
                  <Plus size={14} />
                </button>
              </div>

              <button
                onClick={handleAddToCart}
                disabled={!inStock || addingToCart}
                className="flex-1 inline-flex items-center justify-center gap-3 bg-secondary px-8 h-12 text-[11px] uppercase tracking-editorial text-primary transition-all duration-300 hover:bg-secondary-600 disabled:opacity-50"
              >
                <ShoppingBag size={16} />
                {addingToCart ? 'Adding to bag...' : 'Add to Shopping Bag'}
              </button>
            </div>

            {/* Shopee Alternative Link */}
            <a
              href={product.shopeeLink || 'https://shopee.co.id/'}
              target="_blank"
              rel="noopener noreferrer"
              className="mt-3 flex items-center justify-center gap-2 border border-dark/15 py-3 text-[11px] uppercase tracking-editorial text-dark transition-all duration-300 hover:border-secondary hover:text-secondary"
            >
              Order directly on Shopee Store
            </a>

            {/* Value Props & Assurances */}
            <div className="mt-10 pt-8 border-t border-dark/10 grid grid-cols-2 gap-4">
              <div className="flex items-start gap-3">
                <ShieldCheck size={18} className="text-secondary shrink-0 mt-0.5" />
                <div>
                  <h4 className="text-[11px] uppercase tracking-editorial text-dark font-medium">
                    Self-Manufactured
                  </h4>
                  <p className="text-[11px] font-light text-dark/60 mt-0.5">
                    Carefully tailored with signature embroidery details.
                  </p>
                </div>
              </div>

              <div className="flex items-start gap-3">
                <Truck size={18} className="text-secondary shrink-0 mt-0.5" />
                <div>
                  <h4 className="text-[11px] uppercase tracking-editorial text-dark font-medium">
                    Fast Delivery
                  </h4>
                  <p className="text-[11px] font-light text-dark/60 mt-0.5">
                    Shipped nationwide with real-time order tracking.
                  </p>
                </div>
              </div>
            </div>
          </motion.div>
        </div>
      </div>
    </div>
  );
}
